// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// emit.go turns an authored component into Go source. Every emitter returns
// gofmt'd source or an error; nothing here writes files, so the caller decides
// whether output is compared, staged, or discarded.

import (
	"fmt"
	"go/format"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/septagon-oss/tw/emission"
)

// Artifacts is one component's generated Go source.
type Artifacts struct {
	// Contract is a complete file for contracts/<tier-plural>/<name>.go.
	Contract string
	// ClassLists is the var-block fragment for render/web/classlists.go.
	// It is a fragment rather than a file because the class lists of every
	// component share one var block and one ClassLists() inventory.
	ClassLists string
	// ClassListNames are the identifiers the fragment declares, in the order
	// ClassLists() should register them.
	ClassListNames []string
	// Renderer is the render function; empty for behavioral components,
	// whose renderer is hand-written.
	Renderer string
	// DeliveryDefinition is the catalog constructor.
	DeliveryDefinition string
}

// Option configures one generation run.
type Option func(*settings)

type settings struct{ oracle ClassOracle }

// WithOracle governs the component with a specific class universe. The default
// is tw's enumerable universe; products that own their CSS supply a
// StylesheetOracle instead, and the guarantees are identical either way.
func WithOracle(oracle ClassOracle) Option {
	return func(s *settings) { s.oracle = oracle }
}

// Generate produces every artifact for one component.
func Generate(source Source, style WebStyle, options ...Option) (Artifacts, error) {
	resolved := settings{}
	for _, option := range options {
		option(&resolved)
	}
	if err := validate(source, style, resolved); err != nil {
		return Artifacts{}, err
	}
	contract, err := emitContract(source)
	if err != nil {
		return Artifacts{}, fmt.Errorf("componentgen %s: contract: %w", source.ID, err)
	}
	lists, names := emitClassLists(source, style)
	definition, err := emitDelivery(source, style)
	if err != nil {
		return Artifacts{}, fmt.Errorf("componentgen %s: delivery: %w", source.ID, err)
	}
	out := Artifacts{
		Contract:           contract,
		ClassLists:         lists,
		ClassListNames:     names,
		DeliveryDefinition: definition,
	}
	if source.Presentational {
		renderer, err := emitRenderer(source, style)
		if err != nil {
			return Artifacts{}, fmt.Errorf("componentgen %s: renderer: %w", source.ID, err)
		}
		out.Renderer = renderer
	}
	return out, nil
}

func validate(source Source, style WebStyle, resolved settings) error {
	if strings.TrimSpace(source.ID) == "" {
		return fmt.Errorf("componentgen: component ID is required")
	}
	if style.ID != source.ID {
		return fmt.Errorf("componentgen %s: web style declares ID %q", source.ID, style.ID)
	}
	canonical := 0
	for _, example := range source.Examples {
		if example.Canonical {
			canonical++
		}
	}
	if canonical != 1 {
		return fmt.Errorf("componentgen %s: exactly one canonical example is required, found %d", source.ID, canonical)
	}
	// The delivery layer derives variants from enum props, so an enum prop
	// with no values would produce an unnormalizable contract.
	for _, prop := range source.Props {
		if len(prop.Enum) > 0 && prop.Default != "" && !slices.Contains(prop.Enum, prop.Default) {
			return fmt.Errorf("componentgen %s: prop %q default %q is not one of its values", source.ID, prop.Name, prop.Default)
		}
	}
	if style.VariantProp != "" && len(style.Variants) == 0 {
		return fmt.Errorf("componentgen %s: variant prop %q has no variant utilities", source.ID, style.VariantProp)
	}
	// Every utility must name a real tw method and constant, and must compile
	// to a class inside tw's enumerable universe. Checking here means a typo
	// fails at authoring time rather than surviving into generated source.
	vocabulary, err := LoadVocabulary()
	if err != nil {
		return fmt.Errorf("componentgen %s: load tw vocabulary: %w", source.ID, err)
	}
	if err := vocabulary.ValidateStyle(style); err != nil {
		return fmt.Errorf("componentgen: %w", err)
	}
	// A component governed by another universe has its compiled classes
	// re-checked there, so the "cannot render unstyled" guarantee holds on
	// whichever substrate the product actually ships.
	if resolved.oracle != nil {
		if err := validateAgainstOracle(vocabulary, resolved.oracle, style); err != nil {
			return fmt.Errorf("componentgen %s: %w", source.ID, err)
		}
	}
	// pk-ui forbids a base fragment from declaring a CSS property that any of
	// its variants also declares: two single-class utilities have equal
	// specificity, so the winner would be decided by alphabetical sheet order.
	if err := checkVariantCollisions(vocabulary, style); err != nil {
		return fmt.Errorf("componentgen %s: %w", source.ID, err)
	}
	return nil
}

// checkVariantCollisions reproduces pk-ui's collision invariant at generation
// time, so a base extracted from shared design values cannot silently shadow
// the variants that are supposed to override it.
func checkVariantCollisions(vocabulary *Vocabulary, style WebStyle) error {
	baseProperties, err := declaredProperties(vocabulary, style.Base)
	if err != nil {
		return err
	}
	for _, name := range slices.Sorted(maps.Keys(style.Variants)) {
		variantProperties, err := declaredProperties(vocabulary, style.Variants[name])
		if err != nil {
			return err
		}
		for property := range variantProperties {
			if _, collides := baseProperties[property]; collides {
				return fmt.Errorf(
					"base and variant %q both declare %q — hoist it out of the base into every variant",
					name, property,
				)
			}
		}
	}
	return nil
}

// declaredProperties returns the CSS property names a utility chain declares.
func declaredProperties(vocabulary *Vocabulary, utilities []Utility) (map[string]struct{}, error) {
	out := map[string]struct{}{}
	for _, utility := range utilities {
		compiled, err := vocabulary.Compile([]Utility{utility})
		if err != nil {
			return nil, err
		}
		for _, class := range strings.Fields(compiled) {
			sheet, err := emission.Rules(class)
			if err != nil {
				return nil, err
			}
			for _, property := range declarationNames(sheet.RenderPretty()) {
				out[property] = struct{}{}
			}
		}
	}
	return out, nil
}

// declarationNames pulls property names out of rendered CSS. Custom
// properties are ignored: two utilities may both set a --pk-* variable
// without conflicting, which is how tw composes ring and shadow layers.
func declarationNames(css string) []string {
	var out []string
	for _, line := range strings.Split(css, "\n") {
		line = strings.TrimSpace(line)
		colon := strings.Index(line, ":")
		if colon <= 0 || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "@") {
			continue
		}
		name := strings.TrimSpace(line[:colon])
		if name == "" || strings.ContainsAny(name, "{}.#") {
			continue
		}
		out = append(out, name)
	}
	return out
}

// emitContract writes the typed Props contract.
func emitContract(source Source) (string, error) {
	var b strings.Builder
	b.WriteString("package atoms\n\n")
	b.WriteString("// Code generated by pk-ui componentgen. DO NOT EDIT.\n\n")
	b.WriteString("// Implements: REQ-011.\n// Per: ADR-0031.\n// Discipline: C-14.\n")
	b.WriteString("import \"github.com/septagon-oss/pk-ui/contracts\"\n\n")
	fmt.Fprintf(&b, "// %sProps defines the platform-agnostic properties for a %s component.\n", source.ID, source.ID)
	fmt.Fprintf(&b, "type %sProps struct {\n", source.ID)
	b.WriteString("\tcontracts.ComponentProps\n\n")
	for _, prop := range source.Props {
		tag := prop.Name
		if !prop.Required {
			tag += ",omitempty"
		}
		fmt.Fprintf(&b, "\t%s %s `json:%q`", prop.GoName, prop.Type, tag)
		// The enum comment is generated from the authored values, which is
		// what keeps it from drifting out of step with the variant table.
		if len(prop.Enum) > 0 {
			fmt.Fprintf(&b, " // %s", strings.Join(prop.Enum, ", "))
		} else if prop.Description != "" {
			fmt.Fprintf(&b, " // %s", strings.ToLower(prop.Description[:1])+prop.Description[1:])
		}
		b.WriteString("\n")
	}
	b.WriteString("}\n\n")
	fmt.Fprintf(&b, "// ToMap converts %sProps to map[string]any for unified Component construction.\n", source.ID)
	fmt.Fprintf(&b, "func (p %sProps) ToMap() map[string]any { return propsToMap(p) }\n", source.ID)
	return gofmt(b.String())
}

// emitClassLists writes the var-block fragment and the identifiers it declares.
func emitClassLists(source Source, style WebStyle) (string, []string) {
	var b strings.Builder
	var names []string
	prefix := "cl" + source.ID

	baseName := prefix + "Base"
	if len(style.Variants) == 0 && len(style.Parts) == 0 {
		baseName = prefix
	}
	fmt.Fprintf(&b, "\t%s = %s\n", baseName, utilityChain(style.Base))
	names = append(names, baseName)

	if len(style.Variants) > 0 {
		variantName := prefix + capitalize(style.VariantProp)
		fmt.Fprintf(&b, "\t%s = map[string]tw.ClassList{\n", variantName)
		for _, key := range slices.Sorted(maps.Keys(style.Variants)) {
			fmt.Fprintf(&b, "\t\t%q: %s,\n", key, utilityChain(style.Variants[key]))
		}
		b.WriteString("\t}\n")
		names = append(names, variantName)
	}

	for _, part := range slices.Sorted(maps.Keys(style.Parts)) {
		partName := prefix + capitalize(part)
		fmt.Fprintf(&b, "\t%s = %s\n", partName, utilityChain(style.Parts[part]))
		names = append(names, partName)
	}
	return b.String(), names
}

// utilityChain renders a typed tw builder chain.
func utilityChain(utilities []Utility) string {
	var b strings.Builder
	b.WriteString("tw.New()")
	for _, utility := range utilities {
		if utility.Arg == "" {
			fmt.Fprintf(&b, ".%s()", utility.Method)
			continue
		}
		fmt.Fprintf(&b, ".%s(tw.%s)", utility.Method, utility.Arg)
	}
	return b.String()
}

// emitRenderer writes the render function for a presentational component.
func emitRenderer(source Source, style WebStyle) (string, error) {
	var b strings.Builder
	prefix := "cl" + source.ID
	fmt.Fprintf(&b, "// %s renders atoms.%sProps.\n", source.ID, source.ID)
	fmt.Fprintf(&b, "func %s(p atoms.%sProps) g.Node {\n", source.ID, source.ID)

	if style.VariantProp != "" {
		prop := propByName(source, style.VariantProp)
		fmt.Fprintf(
			&b,
			"\tcl := %sBase.Merge(variantOr(%s%s, p.%s, %q))\n",
			prefix, prefix, capitalize(style.VariantProp), prop.GoName, prop.Default,
		)
	} else {
		fmt.Fprintf(&b, "\tcl := %s\n", prefix)
	}
	b.WriteString("\tvar children []g.Node\n")

	for _, step := range style.Body {
		switch step.Kind {
		case StepBaseAttrs:
			b.WriteString("\tchildren = append(children, baseAttrs(p.ComponentProps)...)\n")
		case StepClasses:
			b.WriteString("\tchildren = append(children, classes(cl.Compile(), p.Class))\n")
		case StepPartIfBool:
			fmt.Fprintf(&b, "\tif p.%s {\n", step.Prop)
			fmt.Fprintf(&b, "\t\tchildren = append(children, h.%s(h.Class(%s%s.Compile())",
				capitalize(step.Element), prefix, capitalize(step.Part))
			for _, attr := range step.Attrs {
				fmt.Fprintf(&b, ", g.Attr(%q, %q)", attr[0], attr[1])
			}
			b.WriteString("))\n\t}\n")
		case StepIconIfSet:
			fmt.Fprintf(&b, "\tif p.%s != \"\" {\n", step.Prop)
			fmt.Fprintf(&b, "\t\tchildren = append(children, icon(p.%s))\n\t}\n", step.Prop)
		case StepText:
			fmt.Fprintf(&b, "\tchildren = append(children, g.Text(p.%s))\n", step.Prop)
		default:
			return "", fmt.Errorf("unknown render step %q", step.Kind)
		}
	}
	fmt.Fprintf(&b, "\treturn h.%s(children...)\n}\n", capitalize(style.Element))
	return gofmt(b.String())
}

// emitDelivery writes the delivery-catalog constructor.
func emitDelivery(source Source, style WebStyle) (string, error) {
	var b strings.Builder
	prefix := "cl" + source.ID
	fmt.Fprintf(&b, "func %sDeliveryDefinition() DeliveryDefinition {\n", lowerFirst(source.ID))
	fmt.Fprintf(&b, "\tcomponentType := %q\n", source.ID)

	frame := fmt.Sprintf("deliveryFrame(componentType, %sBase.Compile()", prefix)
	for _, step := range style.Body {
		if step.Kind != StepText {
			continue
		}
		prop := propByGoName(source, step.Prop)
		fmt.Fprintf(&frameBuilder{&frame}, ",\n\t\tdeliveryText(%q, \"\", %q, %q)", step.Part, step.Fallback, prop.Name)
	}
	frame += ",\n\t)"

	if style.VariantProp != "" {
		fmt.Fprintf(&b, "\troot := deliveryClassBound(\n\t\t%s,\n\t\t%q,\n\t\tcompileClassMap(%s%s),\n\t)\n",
			frame, style.VariantProp, prefix, capitalize(style.VariantProp))
	} else {
		fmt.Fprintf(&b, "\troot := %s\n", frame)
	}

	b.WriteString("\treturn newDeliveryDefinition(\n\t\tcomponentType,\n")
	fmt.Fprintf(&b, "\t\tuicomponent.Tier%s,\n", capitalize(source.Tier))
	b.WriteString("\t\tstableDeliveryID(componentType),\n")
	fmt.Fprintf(&b, "\t\t%q,\n", source.Description)
	b.WriteString("\t\tmap[string]PropertyContract{\n")
	for _, prop := range source.Props {
		fmt.Fprintf(&b, "\t\t\t%q: %s,\n", prop.Name, propertyContract(prop))
	}
	b.WriteString("\t\t},\n\t\tnil,\n")

	canonical := canonicalExample(source)
	fmt.Fprintf(&b, "\t\tdeliveryDesign(%q, %q, %q, root),\n", canonical.Name, source.LibraryPage, source.LibrarySection)
	b.WriteString("\t\t[]DeliveryExample{\n")
	for _, example := range source.Examples {
		constructor := "deliveryExample"
		if example.Canonical {
			constructor = "canonicalDeliveryExample"
		}
		fmt.Fprintf(&b, "\t\t\t%s(componentType, %q, map[string]any{\n", constructor, example.Name)
		for _, key := range slices.Sorted(maps.Keys(example.Props)) {
			fmt.Fprintf(&b, "\t\t\t\t%q: %s,\n", key, literal(example.Props[key]))
		}
		b.WriteString("\t\t\t}),\n")
	}
	b.WriteString("\t\t},\n")
	fmt.Fprintf(&b, "\t\t%s,\n\t)\n}\n", source.ID)
	return gofmt(b.String())
}

// frameBuilder lets the frame expression be appended to with Fprintf.
type frameBuilder struct{ target *string }

func (f *frameBuilder) Write(p []byte) (int, error) {
	*f.target += string(p)
	return len(p), nil
}

// propertyContract picks the delivery helper matching the authored role, so
// role stays a single authored fact rather than being restated per artifact.
func propertyContract(prop Prop) string {
	description := fmt.Sprintf("%q", prop.Description)
	if len(prop.Enum) > 0 {
		values := make([]string, 0, len(prop.Enum))
		for _, value := range prop.Enum {
			values = append(values, fmt.Sprintf("%q", value))
		}
		list := "[]string{" + strings.Join(values, ", ") + "}"
		helper := "enumProperty"
		switch prop.Role {
		case "variant":
			helper = "variantProperty"
		case "tone":
			helper = "toneProperty"
		case "size":
			helper = "sizeProperty"
		}
		return fmt.Sprintf("%s(%s, %q, %s)", helper, list, prop.Default, description)
	}
	switch prop.Role {
	case "content":
		return fmt.Sprintf("contentProperty(%s)", description)
	case "state":
		return fmt.Sprintf("stateProperty(%s)", description)
	case "modifier":
		return fmt.Sprintf("modifierProperty(%s)", description)
	}
	return fmt.Sprintf("contentProperty(%s)", description)
}

func canonicalExample(source Source) Example {
	for _, example := range source.Examples {
		if example.Canonical {
			return example
		}
	}
	return Example{}
}

func propByName(source Source, name string) Prop {
	for _, prop := range source.Props {
		if prop.Name == name {
			return prop
		}
	}
	return Prop{}
}

func propByGoName(source Source, goName string) Prop {
	for _, prop := range source.Props {
		if prop.GoName == goName {
			return prop
		}
	}
	return Prop{}
}

func literal(value any) string {
	switch typed := value.(type) {
	case string:
		return fmt.Sprintf("%q", typed)
	case bool:
		return fmt.Sprintf("%t", typed)
	default:
		return fmt.Sprintf("%#v", typed)
	}
}

func capitalize(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

// gofmt formats a fragment. Fragments that are not whole files (the class-list
// var block) are formatted by the caller-supplied wrapper instead.
func gofmt(source string) (string, error) {
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return "", fmt.Errorf("format generated source: %w\n%s", err, numberLines(source))
	}
	return string(formatted), nil
}

func numberLines(source string) string {
	lines := strings.Split(source, "\n")
	var b strings.Builder
	for index, line := range lines {
		fmt.Fprintf(&b, "%4d | %s\n", index+1, line)
	}
	return b.String()
}

var _ = sort.Strings

// validateAgainstOracle re-checks every class a component emits against a
// second universe. It is what lets one generator serve both the utility
// substrate and a product that owns its stylesheet.
func validateAgainstOracle(vocabulary *Vocabulary, oracle ClassOracle, style WebStyle) error {
	groups := [][]Utility{style.Base}
	for _, name := range slices.Sorted(maps.Keys(style.Variants)) {
		groups = append(groups, style.Variants[name])
	}
	for _, name := range slices.Sorted(maps.Keys(style.Parts)) {
		groups = append(groups, style.Parts[name])
	}
	for _, group := range groups {
		compiled, err := vocabulary.Compile(group)
		if err != nil {
			return err
		}
		for _, class := range strings.Fields(compiled) {
			if err := oracle.Resolve(class); err != nil {
				return err
			}
		}
	}
	return nil
}
