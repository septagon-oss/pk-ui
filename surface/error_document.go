package surface

import (
	"context"
)

// ErrorDocument is the renderer-neutral presentation state for a complete
// browser error response. Copy is resolved before crossing this boundary;
// renderer implementations own visual composition and contextual theming.
type ErrorDocument struct {
	StatusCode   int    `json:"statusCode"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ErrorDetails string `json:"errorDetails,omitempty"`
	HomeURL      string `json:"homeUrl,omitempty"`
	HomeLabel    string `json:"homeLabel,omitempty"`
	BackURL      string `json:"backUrl,omitempty"`
	BackLabel    string `json:"backLabel,omitempty"`
	LoginURL     string `json:"loginUrl,omitempty"`
	LoginLabel   string `json:"loginLabel,omitempty"`
	SupportURL   string `json:"supportUrl,omitempty"`
	SupportLabel string `json:"supportLabel,omitempty"`
	StatusURL    string `json:"statusUrl,omitempty"`
	StatusLabel  string `json:"statusLabel,omitempty"`
	RetryURL     string `json:"retryUrl,omitempty"`
	RetryLabel   string `json:"retryLabel,omitempty"`
	BrandName    string `json:"brandName,omitempty"`
	// SignedInLabel, when set, reports that the visitor already holds an
	// authenticated session (copy is pre-resolved, e.g. "Signed in as
	// marta@example.com"). Renderers must surface it in the page context and
	// suppress the sign-in action — offering a login to a signed-in visitor
	// is a dead end.
	SignedInLabel string `json:"signedInLabel,omitempty"`
	ReferenceID   string `json:"referenceId,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Direction     string `json:"direction,omitempty"`
	Nonce         string `json:"nonce,omitempty"`
}

// ErrorDocumentRenderer builds a complete HTML document. Implementations may
// read presentation-specific values such as theme selection from ctx, keeping
// those implementation types out of capability and backend packages.
type ErrorDocumentRenderer interface {
	RenderErrorDocument(ctx context.Context, document ErrorDocument) Renderable
}
