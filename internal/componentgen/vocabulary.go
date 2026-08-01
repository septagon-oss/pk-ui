// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// vocabulary.go closes the gap between an authored utility and a class that
// actually exists. A Utility names a tw builder method and a Go constant as
// strings, which is what lets the generated code read as typed Go — but
// strings do not compile, so without this file a typo survives generation and
// only surfaces when someone builds the generated tree.
//
// Validation is total: the method must exist on tw.ClassList, the constant
// must exist with the right type, and the class the pair compiles to must
// resolve in tw's enumerable universe. The last check is the important one —
// it is the same oracle (emission.Rules) that pk-ui's own loop-closure test
// uses, so a utility that passes here cannot produce an unstyled element.
//
// The constant table is derived by parsing tw's source rather than being
// hand-listed, because the name-to-value transform is irregular per type
// (SurfaceTertiary -> "surface-tertiary", S2_5 -> "0.5", Radius2XL -> "2xl")
// and a hand-list would drift the moment tw gains a token.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/septagon-oss/tw"
	"github.com/septagon-oss/tw/emission"
)

// Vocabulary is tw's typed utility surface: which builder methods exist, what
// they take, and which constants are legal arguments.
type Vocabulary struct {
	// methods maps a ClassList method name to its single parameter type.
	// A nil type means the method takes no arguments.
	methods map[string]reflect.Type
	// constants maps a type name to that type's constant names and values.
	constants map[string]map[string]string
}

var (
	loadedVocabulary   *Vocabulary
	loadVocabularyOnce sync.Once
	loadVocabularyErr  error
)

// LoadVocabulary reads tw's surface once per process.
func LoadVocabulary() (*Vocabulary, error) {
	loadVocabularyOnce.Do(func() {
		loadedVocabulary, loadVocabularyErr = buildVocabulary()
	})
	return loadedVocabulary, loadVocabularyErr
}

func buildVocabulary() (*Vocabulary, error) {
	vocabulary := &Vocabulary{
		methods:   map[string]reflect.Type{},
		constants: map[string]map[string]string{},
	}

	// Methods come from reflection, so they track tw's actual API rather than
	// a description of it.
	listType := reflect.TypeOf(tw.New())
	for index := range listType.NumMethod() {
		method := listType.Method(index)
		switch method.Type.NumIn() {
		case 1: // receiver only
			vocabulary.methods[method.Name] = nil
		case 2:
			vocabulary.methods[method.Name] = method.Type.In(1)
		default:
			// Multi-argument builders are outside the authored vocabulary;
			// omitting them makes them unauthorable rather than silently wrong.
		}
	}

	directory, err := twSourceDirectory()
	if err != nil {
		return nil, err
	}
	fileSet := token.NewFileSet()
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read tw source at %s: %w", directory, err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fileSet, filepath.Join(directory, name), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		collectConstants(file, vocabulary.constants)
	}
	if len(vocabulary.constants) == 0 {
		return nil, fmt.Errorf("parsed no typed constants from tw source at %s", directory)
	}
	return vocabulary, nil
}

// twSourceDirectory locates tw in the module cache.
func twSourceDirectory() (string, error) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/septagon-oss/tw")
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("locate tw module: %w", err)
	}
	directory := strings.TrimSpace(string(output))
	if directory == "" {
		return "", fmt.Errorf("go list returned no directory for tw")
	}
	return filepath.Clean(directory), nil
}

// collectConstants records every typed string constant, keyed by type name.
func collectConstants(file *ast.File, into map[string]map[string]string) {
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.CONST {
			continue
		}
		// Within one const block the type carries over from the last spec
		// that declared it, which is how tw writes its token tables.
		currentType := ""
		for _, spec := range general.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if identifier, ok := value.Type.(*ast.Ident); ok {
				currentType = identifier.Name
			}
			if currentType == "" || len(value.Values) == 0 {
				continue
			}
			literal, ok := value.Values[0].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				continue
			}
			unquoted, err := strconv.Unquote(literal.Value)
			if err != nil {
				continue
			}
			for _, name := range value.Names {
				if into[currentType] == nil {
					into[currentType] = map[string]string{}
				}
				into[currentType][name.Name] = unquoted
			}
		}
	}
}

// Compile resolves one utility chain to its class string, rejecting anything
// that is not a real, resolvable tw utility.
func (v *Vocabulary) Compile(utilities []Utility) (string, error) {
	list := tw.New()
	for _, utility := range utilities {
		next, err := v.apply(list, utility)
		if err != nil {
			return "", err
		}
		list = next
	}
	compiled := list.Compile()
	// The enumerable-universe check: the same oracle pk-ui's loop-closure
	// test uses. A class that fails here would render unstyled.
	for _, class := range strings.Fields(compiled) {
		if _, err := emission.Rules(class); err != nil {
			return "", fmt.Errorf("utility produces unusable class %q: %w", class, err)
		}
	}
	return compiled, nil
}

func (v *Vocabulary) apply(list tw.ClassList, utility Utility) (tw.ClassList, error) {
	parameter, known := v.methods[utility.Method]
	if !known {
		return list, fmt.Errorf("tw.ClassList has no method %q", utility.Method)
	}
	receiver := reflect.ValueOf(list)
	method := receiver.MethodByName(utility.Method)

	if parameter == nil {
		if utility.Arg != "" {
			return list, fmt.Errorf("tw.ClassList.%s takes no argument, got %q", utility.Method, utility.Arg)
		}
		return method.Call(nil)[0].Interface().(tw.ClassList), nil
	}
	if utility.Arg == "" {
		return list, fmt.Errorf("tw.ClassList.%s requires a %s argument", utility.Method, parameter.Name())
	}

	typeName := parameter.Name()
	values, ok := v.constants[typeName]
	if !ok {
		return list, fmt.Errorf(
			"tw.ClassList.%s takes %s, which declares no constants — it is not authorable",
			utility.Method, typeName,
		)
	}
	value, ok := values[utility.Arg]
	if !ok {
		return list, fmt.Errorf("tw.%s is not a constant of type %s (required by %s)",
			utility.Arg, typeName, utility.Method)
	}
	argument := reflect.ValueOf(value).Convert(parameter)
	return method.Call([]reflect.Value{argument})[0].Interface().(tw.ClassList), nil
}

// ValidateStyle checks every utility a component declares.
func (v *Vocabulary) ValidateStyle(style WebStyle) error {
	if _, err := v.Compile(style.Base); err != nil {
		return fmt.Errorf("%s base: %w", style.ID, err)
	}
	for name, utilities := range style.Variants {
		if _, err := v.Compile(utilities); err != nil {
			return fmt.Errorf("%s variant %q: %w", style.ID, name, err)
		}
	}
	for name, utilities := range style.Parts {
		if _, err := v.Compile(utilities); err != nil {
			return fmt.Errorf("%s part %q: %w", style.ID, name, err)
		}
	}
	return nil
}

// Name identifies tw's universe in error messages, satisfying ClassOracle.
func (v *Vocabulary) Name() string { return "tw's enumerable universe" }

// Resolve rejects a class tw cannot resolve, satisfying ClassOracle. It is the
// same oracle pk-ui's loop-closure test uses.
func (v *Vocabulary) Resolve(class string) error {
	if _, err := emission.Rules(class); err != nil {
		return err
	}
	return nil
}

// Properties lists the CSS properties a tw class declares, satisfying
// ClassOracle.
func (v *Vocabulary) Properties(class string) ([]string, error) {
	sheet, err := emission.Rules(class)
	if err != nil {
		return nil, err
	}
	return declarationNames(sheet.RenderPretty()), nil
}

var _ ClassOracle = (*Vocabulary)(nil)
var _ ClassOracle = (*StylesheetOracle)(nil)
