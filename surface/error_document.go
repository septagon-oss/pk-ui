package surface

import (
	"context"
)

// ErrorDocument is the renderer-neutral presentation state for a complete
// browser error response. Copy is resolved before crossing this boundary;
// renderer implementations own visual composition and contextual theming.
type ErrorDocument struct {
	StatusCode   int
	Title        string
	Description  string
	ErrorDetails string
	HomeURL      string
	HomeLabel    string
	BackURL      string
	BackLabel    string
	LoginURL     string
	LoginLabel   string
	SupportURL   string
	SupportLabel string
	StatusURL    string
	StatusLabel  string
	RetryURL     string
	RetryLabel   string
	BrandName    string
	ReferenceID  string
	Locale       string
	Direction    string
	Nonce        string
}

// ErrorDocumentRenderer builds a complete HTML document. Implementations may
// read presentation-specific values such as theme selection from ctx, keeping
// those implementation types out of capability and backend packages.
type ErrorDocumentRenderer interface {
	RenderErrorDocument(ctx context.Context, document ErrorDocument) Renderable
}
