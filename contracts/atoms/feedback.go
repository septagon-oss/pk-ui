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
	Tone  string `json:"tone,omitempty"`  // brand, success, warning, danger, info
}

// ToMap converts SpinnerProps to map[string]any for unified Component construction.
func (p SpinnerProps) ToMap() map[string]any { return propsToMap(p) }

// ProgressProps defines properties for a progress indicator.
type ProgressProps struct {
	contracts.ComponentProps

	Value         int    `json:"value"`         // 0-100
	Max           int    `json:"max,omitempty"` // default 100
	Label         string `json:"label,omitempty"`
	ShowText      bool   `json:"showText,omitempty"` // show percentage
	Variant       string `json:"variant,omitempty"`  // default, success, warning, error
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

	Message  string `json:"message"`
	Variant  string `json:"variant,omitempty"` // success, error, warning, info
	Duration int    `json:"duration,omitempty"`
}

// ToMap converts ToastProps to map[string]any for unified Component construction.
func (p ToastProps) ToMap() map[string]any { return propsToMap(p) }
