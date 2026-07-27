package surfacetest

import (
	"context"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-ui/surface"
)

// ErrorDocumentRendererConformance runs the public safety and semantic
// behavior required from every complete error-document renderer.
func ErrorDocumentRendererConformance(
	t *testing.T,
	newRenderer func(t *testing.T) surface.ErrorDocumentRenderer,
) {
	t.Helper()

	t.Run("renders a complete semantic document", func(t *testing.T) {
		renderer := newRenderer(t)
		document := surface.ErrorDocument{
			StatusCode:  httpStatusNotFound,
			Title:       "Page not found",
			Description: "The requested page moved.",
			HomeURL:     "/workspace",
			HomeLabel:   "Go Home",
			ReferenceID: "req-42",
			Locale:      "en",
			Direction:   "ltr",
			Nonce:       "nonce-42",
		}
		var output strings.Builder
		node := renderer.RenderErrorDocument(context.Background(), document)
		if node == nil {
			t.Fatal("renderer returned a nil document")
		}
		if err := node.Render(&output); err != nil {
			t.Fatalf("render document: %v", err)
		}
		html := output.String()
		for _, expected := range []string{
			"<!DOCTYPE html>",
			`data-component="error-page"`,
			"Page not found",
			"The requested page moved.",
			`href="/workspace"`,
			"req-42",
		} {
			if !strings.Contains(html, expected) {
				t.Errorf("document missing %q: %s", expected, html)
			}
		}
	})

	t.Run("escapes untrusted copy", func(t *testing.T) {
		renderer := newRenderer(t)
		var output strings.Builder
		node := renderer.RenderErrorDocument(context.Background(), surface.ErrorDocument{
			StatusCode:  httpStatusNotFound,
			Title:       `<script>alert("title")</script>`,
			Description: `<img src=x onerror=alert("description")>`,
			HomeURL:     "/",
			HomeLabel:   "Go Home",
		})
		if node == nil {
			t.Fatal("renderer returned a nil document")
		}
		if err := node.Render(&output); err != nil {
			t.Fatalf("render document: %v", err)
		}
		html := output.String()
		if strings.Contains(html, `<script>alert("title")</script>`) ||
			strings.Contains(html, `<img src=x onerror=alert("description")>`) {
			t.Fatalf("renderer emitted unescaped copy: %s", html)
		}
	})
}

const httpStatusNotFound = 404
