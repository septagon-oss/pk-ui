package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// EmptyStateProps defines properties for an empty data state placeholder.
type EmptyStateProps struct {
	contracts.ComponentProps

	Icon        string             `json:"icon,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Compact     bool               `json:"compact,omitempty"`
	Bordered    bool               `json:"bordered,omitempty"`
	Actions     []EmptyStateAction `json:"actions,omitempty"`
}

// EmptyStateAction describes a CTA button within the empty state.
type EmptyStateAction struct {
	Label   string `json:"label"`
	Href    string `json:"href"`
	Icon    string `json:"icon,omitempty"`
	Variant string `json:"variant,omitempty"` // primary, secondary
}

// ToMap converts EmptyStateProps to map[string]any for unified Component construction.
func (p EmptyStateProps) ToMap() map[string]any { return propsToMap(p) }
