// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// oracle.go generalizes the one guarantee that makes generation safe: a
// generated component may never emit a class that renders nothing.
//
// For pk-ui that guarantee comes from tw's enumerable universe, checked with
// emission.Rules. But the estate's largest product surface is not on that
// substrate — Collect carries 979 bespoke class literals over ~13k lines of
// hand-authored CSS and zero tw.New() calls. Validating its components against
// tw would reject every one of them, and dropping validation would give up the
// only property that makes generated markup trustworthy.
//
// The resolution is that the *universe* is a parameter, not a constant. tw is
// one oracle; a product's own stylesheet is another. Both answer the same
// question — "is this class real?" — so a component declares which universe
// governs it and the generator's guarantees are unchanged.

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// ClassOracle decides whether a class is real and what it declares.
//
// Resolve reports an error for any class the universe does not define.
// Properties returns the CSS property names a class declares, which the
// base/variant collision check needs; an oracle that cannot determine them
// returns nil, and collision checking degrades to a no-op rather than
// producing false failures.
type ClassOracle interface {
	// Name identifies the universe in error messages.
	Name() string
	// Resolve rejects a class the universe does not define.
	Resolve(class string) error
	// Properties lists the CSS properties a class declares.
	Properties(class string) ([]string, error)
}

// StylesheetOracle validates classes against the selectors a stylesheet
// actually defines. It is the oracle for products that own their CSS rather
// than composing utilities.
type StylesheetOracle struct {
	name    string
	classes map[string][]string
}

// classSelector matches a class selector at the start of a compound selector,
// tolerating the escapes and pseudo-suffixes real stylesheets contain.
var classSelector = regexp.MustCompile(`\.(-?[_a-zA-Z]+[_a-zA-Z0-9-]*)`)

// declaration matches "property: value" inside a rule body.
var declaration = regexp.MustCompile(`^\s*([a-zA-Z-]+)\s*:`)

// LoadStylesheetOracle reads one or more CSS files and indexes every class
// selector they define, together with the properties each declares.
//
// Parsing is deliberately shallow: it records which classes exist and what
// they set, which is all the two guarantees need. It does not model the
// cascade, specificity, or media-query scoping — a generator that needed
// those would be deciding visual outcomes, which is not its job.
func LoadStylesheetOracle(name string, paths ...string) (*StylesheetOracle, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("stylesheet oracle %q: at least one stylesheet is required", name)
	}
	oracle := &StylesheetOracle{name: name, classes: map[string][]string{}}
	for _, path := range paths {
		if err := oracle.indexFile(path); err != nil {
			return nil, err
		}
	}
	if len(oracle.classes) == 0 {
		return nil, fmt.Errorf("stylesheet oracle %q: no class selectors found in %s", name, strings.Join(paths, ", "))
	}
	return oracle, nil
}

func (o *StylesheetOracle) indexFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("stylesheet oracle %q: open %s: %w", o.name, path, err)
	}
	defer func() {
		// justified: read-only handle; a close failure cannot affect the
		// index that was already built.
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var pending []string // classes whose rule body is currently open
	inBody := false
	inComment := false

	for scanner.Scan() {
		line := scanner.Text()
		if inComment {
			if index := strings.Index(line, "*/"); index >= 0 {
				line = line[index+2:]
				inComment = false
			} else {
				continue
			}
		}
		for {
			start := strings.Index(line, "/*")
			if start < 0 {
				break
			}
			end := strings.Index(line[start:], "*/")
			if end < 0 {
				line = line[:start]
				inComment = true
				break
			}
			line = line[:start] + line[start+end+2:]
		}

		if !inBody {
			if brace := strings.Index(line, "{"); brace >= 0 {
				pending = classesIn(line[:brace])
				for _, class := range pending {
					if _, seen := o.classes[class]; !seen {
						o.classes[class] = nil
					}
				}
				inBody = true
				line = line[brace+1:]
			} else {
				continue
			}
		}
		if closing := strings.Index(line, "}"); closing >= 0 {
			o.recordDeclarations(pending, line[:closing])
			pending = nil
			inBody = false
			continue
		}
		o.recordDeclarations(pending, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stylesheet oracle %q: read %s: %w", o.name, path, err)
	}
	return nil
}

func (o *StylesheetOracle) recordDeclarations(classes []string, body string) {
	if len(classes) == 0 {
		return
	}
	for _, statement := range strings.Split(body, ";") {
		match := declaration.FindStringSubmatch(statement)
		if match == nil {
			continue
		}
		property := match[1]
		// Custom properties compose rather than conflict, matching the tw
		// oracle's treatment.
		if strings.HasPrefix(property, "--") {
			continue
		}
		for _, class := range classes {
			if !slices.Contains(o.classes[class], property) {
				o.classes[class] = append(o.classes[class], property)
			}
		}
	}
}

// classesIn extracts the class names a selector list defines.
func classesIn(selectors string) []string {
	var out []string
	for _, match := range classSelector.FindAllStringSubmatch(selectors, -1) {
		if !slices.Contains(out, match[1]) {
			out = append(out, match[1])
		}
	}
	return out
}

// Name identifies the universe.
func (o *StylesheetOracle) Name() string { return o.name }

// Resolve rejects a class the stylesheet does not define.
func (o *StylesheetOracle) Resolve(class string) error {
	if _, ok := o.classes[class]; ok {
		return nil
	}
	return fmt.Errorf("%s does not define %q — it would render unstyled%s",
		o.name, class, suggest(class, o.Classes()))
}

// Properties lists what a class declares.
func (o *StylesheetOracle) Properties(class string) ([]string, error) {
	properties, ok := o.classes[class]
	if !ok {
		return nil, o.Resolve(class)
	}
	return slices.Clone(properties), nil
}

// Classes returns every class the stylesheet defines, sorted.
func (o *StylesheetOracle) Classes() []string {
	out := make([]string, 0, len(o.classes))
	for class := range o.classes {
		out = append(out, class)
	}
	sort.Strings(out)
	return out
}

// suggest offers the closest defined class, because the common failure is a
// typo and naming an alternative turns a rejection into a fix.
func suggest(class string, known []string) string {
	best, bestDistance := "", len(class)/2+1
	for _, candidate := range known {
		distance := editDistance(class, candidate)
		if distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	if best == "" {
		return ""
	}
	return fmt.Sprintf(" (did you mean %q?)", best)
}

func editDistance(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	previous := make([]int, len(a)+1)
	current := make([]int, len(a)+1)
	for i := range previous {
		previous[i] = i
	}
	for j := 1; j <= len(b); j++ {
		current[0] = j
		for i := 1; i <= len(a); i++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[i] = min(min(current[i-1]+1, previous[i]+1), previous[i-1]+cost)
		}
		copy(previous, current)
	}
	return previous[len(a)]
}
