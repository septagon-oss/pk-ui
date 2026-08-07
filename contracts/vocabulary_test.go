// Implements: REQ-011 (canonical design vocabulary ratchet).
// Per: ADR-0031.
// Discipline: C-14.

package contracts

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// axisField describes one documented enum found on a contract struct field.
type axisField struct {
	strct string
	json  string
	kind  string // size, tone, state, status
	seen  []string
}

func TestVocabulary_PureVocabularyInvariants(t *testing.T) {
	t.Run("feedback tones stay a canonical subset", func(t *testing.T) {
		for _, tone := range FeedbackTones {
			if !contains(CanonicalTones, tone) {
				t.Errorf("feedback tone %q is not a canonical tone", tone)
			}
		}
	})

	t.Run("scale bindings reference real families", func(t *testing.T) {
		for key, family := range ScaleAxisBindings {
			if len(ScaleValues(family)) == 0 {
				t.Errorf("%s bound to unknown scale family %q", key, family)
			}
		}
	})
}

func TestVocabulary_ContractEnumsMatchCanonicalBindings(t *testing.T) {
	fields := collectAxisFields(t)
	var violations []string

	seen := make(map[string]bool, len(fields))
	for _, f := range fields {
		key := f.strct + "." + f.json
		seen[key] = true
		switch f.kind {
		case "size":
			family, bound := ScaleAxisBindings[key]
			if !bound {
				violations = append(violations,
					key+" documents a size enum but has no ScaleAxisBindings entry in vocabulary.go")
				continue
			}
			checkScaleEnum(&violations, key, family, f.seen)
		case "tone":
			bound, declared := ToneAxisBindings[key]
			if !declared {
				violations = append(violations,
					key+" documents a tone enum but has no ToneAxisBindings entry in vocabulary.go")
				continue
			}
			if bound {
				for _, v := range f.seen {
					if !contains(CanonicalTones, v) {
						violations = append(violations,
							key+" uses tone "+strconv(v)+" outside the canonical tone set")
					}
				}
			}
		case "state", "status":
			values, bound := StateAxisBindings[key]
			if !bound {
				violations = append(violations,
					key+" documents a "+f.kind+" enum but has no StateAxisBindings entry in vocabulary.go")
				continue
			}
			for _, v := range f.seen {
				if !contains(values, v) {
					violations = append(violations,
						key+" uses "+f.kind+" "+strconv(v)+" outside its canonical set")
				}
			}
		}
	}

	for key, family := range ScaleAxisBindings {
		if !seen[key] {
			violations = append(violations,
				key+" is bound to scale family "+string(family)+" but no such field exists in the contracts")
		}
	}
	for key, bound := range ToneAxisBindings {
		if bound && !seen[key] {
			violations = append(violations,
				key+" is bound to the canonical tone set but no such field exists in the contracts")
		}
	}
	for key := range StateAxisBindings {
		if !seen[key] {
			violations = append(violations,
				key+" is bound to a canonical state axis but no such field exists in the contracts")
		}
	}

	sort.Strings(violations)
	for _, v := range violations {
		t.Error(v)
	}
}

// checkScaleEnum validates one documented size enum against its scale family:
// every value must be a canonical member of the declared family.
func checkScaleEnum(violations *[]string, key string, family ScaleFamily, seen []string) {
	canonical := ScaleValues(family)
	for _, v := range seen {
		if !contains(canonical, v) {
			*violations = append(*violations,
				key+" uses size "+strconv(v)+" outside the "+string(family)+" scale ["+strings.Join(canonical, ", ")+"]")
		}
	}
}

// collectAxisFields walks this package's source (including subpackages) and
// returns every struct field whose JSON name is a shared axis (size, tone,
// state, status) and whose field comment documents an enum ("// a, b, c").
func collectAxisFields(t *testing.T) []axisField {
	t.Helper()
	root := packageRoot(t)
	var out []axisField

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		file, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			t.Fatalf("parse %s: %v", path, parseErr)
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				strct, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				out = append(out, axisFieldsOf(typeSpec.Name.Name, strct)...)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk contracts: %v", err)
	}
	return out
}

func axisFieldsOf(structName string, strct *ast.StructType) []axisField {
	var out []axisField
	for _, field := range strct.Fields.List {
		jsonName := jsonTagName(field.Tag)
		kind := axisKind(jsonName)
		if kind == "" || field.Comment == nil {
			continue
		}
		if values := enumValues(field.Comment.Text()); len(values) > 0 {
			out = append(out, axisField{
				strct: structName,
				json:  jsonName,
				kind:  kind,
				seen:  values,
			})
		}
	}
	return out
}

func jsonTagName(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	raw := strings.Trim(tag.Value, "`")
	for _, part := range strings.Split(raw, " ") {
		if !strings.HasPrefix(part, `json:"`) {
			continue
		}
		name := strings.TrimPrefix(part, `json:"`)
		name = strings.TrimSuffix(name, `"`)
		return strings.SplitN(name, ",", 2)[0]
	}
	return ""
}

func axisKind(jsonName string) string {
	switch jsonName {
	case "size":
		return "size"
	case "tone":
		return "tone"
	case "state":
		return "state"
	case "status":
		return "status"
	default:
		return ""
	}
}

// enumValues extracts comma-separated enum values from a field comment of the
// form "// xs, sm, md (note)". Parenthetical asides and prose comments yield
// no values.
func enumValues(comment string) []string {
	comment = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(comment), "//"))
	if comment == "" {
		return nil
	}
	if idx := strings.Index(comment, "("); idx >= 0 {
		comment = comment[:idx]
	}
	if idx := strings.Index(comment, ";"); idx >= 0 {
		comment = comment[:idx]
	}
	parts := strings.Split(comment, ",")
	if len(parts) == 0 {
		return nil
	}
	var values []string
	for _, part := range parts {
		p := strings.ToLower(strings.TrimSpace(part))
		if p == "" {
			continue
		}
		// Prose ("grid columns: 1, 2, 3", "1-6") is not an enum.
		if strings.ContainsAny(p, " :\t") {
			return nil
		}
		values = append(values, p)
	}
	if len(values) < 2 {
		return nil
	}
	return values
}

func packageRoot(t *testing.T) string {
	t.Helper()
	// This file lives at the contracts package root; walking "." covers the
	// root file plus the tier subpackages (atoms, molecules, ...).
	abs, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("resolve contracts dir: %v", err)
	}
	return abs
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func strconv(v string) string { return `"` + v + `"` }
