package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// SkeletonProps defines properties for a loading placeholder. A skeleton is
// the loading rendering of content that has not arrived yet: it holds the
// geometry of the finished component so the layout does not shift when the
// real content swaps in.
type SkeletonProps struct {
	contracts.ComponentProps

	Shape string `json:"shape,omitempty"` // block, text, circle
	Size  string `json:"size,omitempty"`  // sm, md, lg
	Lines int    `json:"lines,omitempty"` // shape=text: placeholder line count (default 1)
}

// ToMap converts SkeletonProps to map[string]any for unified Component construction.
func (p SkeletonProps) ToMap() map[string]any { return propsToMap(p) }

// DeferredSlotProps defines a placeholder region that HTMX replaces with a
// server-rendered fragment. Get names the fragment URL; Trigger defaults to
// "load" and Swap to "outerHTML", so the fragment replaces the slot (and its
// skeleton) wholesale with zero client code.
type DeferredSlotProps struct {
	contracts.ComponentProps
	contracts.HTMXProps
}

// ToMap converts DeferredSlotProps to map[string]any for unified Component construction.
func (p DeferredSlotProps) ToMap() map[string]any { return propsToMap(p) }
