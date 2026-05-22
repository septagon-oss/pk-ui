package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// SpinnerProps defines properties for a loading spinner.
type SpinnerProps struct {
	contracts.ComponentProps

	Label string `json:"label,omitempty"` // sr-only text
	Size  string `json:"size,omitempty"`  // small, medium, large
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

// SkeletonProps defines properties for a loading placeholder.
type SkeletonProps struct {
	contracts.ComponentProps

	Type    string `json:"type,omitempty"` // text, circle, rect, card
	Lines   int    `json:"lines,omitempty"`
	Animate bool   `json:"animate,omitempty"`
	Width   string `json:"width,omitempty"`
	Height  string `json:"height,omitempty"`
}

// ToMap converts SkeletonProps to map[string]any for unified Component construction.
func (p SkeletonProps) ToMap() map[string]any { return propsToMap(p) }

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
