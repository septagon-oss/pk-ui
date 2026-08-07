package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// SpinnerProps defines properties for a loading spinner.
type SpinnerProps struct {
	contracts.ComponentProps

	Label string `json:"label,omitempty"` // sr-only text
	Size  string `json:"size,omitempty"`  // xs, sm, md, lg, xl, 2xl
	Tone  string `json:"tone,omitempty"`  // neutral, brand, success, warning, danger, info (default brand)
}

// ToMap converts SpinnerProps to map[string]any for unified Component construction.
func (p SpinnerProps) ToMap() map[string]any { return propsToMap(p) }

// ProgressProps defines properties for a progress indicator.
type ProgressProps struct {
	contracts.ComponentProps

	Value         int    `json:"value"`         // clamped to 0-Max
	Max           int    `json:"max,omitempty"` // positive maximum; default 100
	Label         string `json:"label,omitempty"`
	AriaLabel     string `json:"ariaLabel,omitempty"` // accessible name without a visible label
	ShowText      bool   `json:"showText,omitempty"`  // show normalized percentage
	Tone          string `json:"tone,omitempty"`      // neutral, brand, success, warning, danger, info (default brand)
	Size          string `json:"size,omitempty"`      // sm, md, lg
	Indeterminate bool   `json:"indeterminate,omitempty"`
}

// ToMap converts ProgressProps to map[string]any for unified Component construction.
func (p ProgressProps) ToMap() map[string]any { return propsToMap(p) }

// SkeletonProps lives in skeleton.go alongside DeferredSlotProps: the earlier
// contract-only draft here (free-string width/height) predated the audited
// class pipeline and had no renderer or consumers.

// TooltipProps defines properties for a tooltip.
type TooltipProps struct {
	contracts.ComponentProps

	Content  string `json:"content"`
	Position string `json:"position,omitempty"` // top, bottom, left, right
	Delay    int    `json:"delay,omitempty"`    // ms before showing
}

// ToMap converts TooltipProps to map[string]any for unified Component construction.
func (p TooltipProps) ToMap() map[string]any { return propsToMap(p) }

// ToastProps defines properties for a toast notification.
type ToastProps struct {
	contracts.ComponentProps

	Message    string `json:"message"`
	Title      string `json:"title,omitempty"`
	Tone       string `json:"tone,omitempty"`       // neutral, info, success, warning, danger (default info)
	Duration   int    `json:"duration,omitempty"`   // auto-dismiss delay in milliseconds
	Persistent bool   `json:"persistent,omitempty"` // remains until explicitly dismissed
	Position   string `json:"position,omitempty"`   // top-right, top-left, bottom-right, bottom-left
	Closable   *bool  `json:"closable,omitempty"`   // defaults to true when omitted
	CloseLabel string `json:"closeLabel,omitempty"`
}

// ToMap converts ToastProps to map[string]any for unified Component construction.
func (p ToastProps) ToMap() map[string]any { return propsToMap(p) }
