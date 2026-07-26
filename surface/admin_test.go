package surface

import (
	"context"
	"io"
	"strings"
	"testing"
)

type testRenderable string

func (r testRenderable) Render(w io.Writer) error {
	_, err := io.WriteString(w, string(r))
	return err
}

type testRenderer struct{}

func (testRenderer) Render(context.Context, string, string) Renderable {
	return testRenderable("rendered")
}
func (testRenderer) CanRender(sectionID string) bool { return sectionID == "users" }
func (testRenderer) Priority() int                   { return 10 }

func TestSectionRendererAcceptsFrameworkNeutralOutput(t *testing.T) {
	t.Parallel()

	var renderer SectionRenderer = testRenderer{}
	var output strings.Builder
	if err := renderer.Render(context.Background(), "users", "/admin/users").Render(&output); err != nil {
		t.Fatalf("Render: %v", err)
	}
	if output.String() != "rendered" {
		t.Fatalf("output = %q", output.String())
	}
}

func TestMenuItemPreservesHierarchy(t *testing.T) {
	t.Parallel()

	item := MenuItem{ID: "users", SubItems: []MenuItem{{ID: "users-list"}}}
	if got := item.SubItems[0].ID; got != "users-list" {
		t.Fatalf("child ID = %q", got)
	}
}
