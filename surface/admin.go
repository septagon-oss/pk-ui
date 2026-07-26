// Package surface defines renderer-neutral UI surface contracts.
package surface

import (
	"context"
	"io"
)

// Renderable is the minimal output contract shared by server-side renderers.
// gomponents.Node and template-backed renderers satisfy it without adapters.
type Renderable interface {
	Render(io.Writer) error
}

// SectionRenderer renders one named admin section.
type SectionRenderer interface {
	Render(ctx context.Context, sectionID string, requestPath string) Renderable
	CanRender(sectionID string) bool
	Priority() int
}

// MenuItem is the renderer-neutral navigation description for an admin
// surface. Recursive sub-items preserve hierarchy without prescribing a shell.
type MenuItem struct {
	ID          string     `json:"id"`
	ModuleID    string     `json:"module_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Icon        string     `json:"icon"`
	Route       string     `json:"route"`
	Order       int        `json:"order"`
	Permissions []string   `json:"permissions"`
	SubItems    []MenuItem `json:"subItems,omitempty"`
}

// AdminPreviewProvider is the smallest contract required to project a
// module-owned admin surface into a preview or design tool.
type AdminPreviewProvider interface {
	GetMenuItems() []MenuItem
	GetSectionRenderers() []SectionRenderer
}
