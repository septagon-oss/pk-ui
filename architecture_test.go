// Validates: REQ-011 (the OSS UI pillar never depends on PlatformKit).
// Per: ADR-0031.
// Discipline: C-14.

package pkui_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestOSSPillarHasNoProprietaryDependency(t *testing.T) {
	goMod, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if strings.Contains(string(goMod), "github.com/septagon-dev/") {
		t.Fatal("OSS pk-ui go.mod depends on proprietary PlatformKit code")
	}
	if strings.Contains(string(goMod), "\nreplace ") || strings.Contains(string(goMod), "\nreplace(") {
		t.Fatal("OSS pk-ui go.mod contains a local replace directive")
	}

	err = filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "vendor":
				return filepath.SkipDir
			}
			if strings.HasPrefix(entry.Name(), ".tmp") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			t.Errorf("parse %s: %v", path, err)
			return nil
		}
		for _, spec := range parsed.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				t.Errorf("unquote import in %s: %v", path, err)
				continue
			}
			if strings.HasPrefix(importPath, "github.com/septagon-dev/") {
				t.Errorf("OSS source %s imports proprietary package %q", path, importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk OSS source: %v", err)
	}
}
