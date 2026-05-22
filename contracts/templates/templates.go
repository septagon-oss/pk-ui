package templates

import "github.com/septagon-oss/platformkit-ui/contracts"

// ListTemplateProps defines properties for an entity list page.
type ListTemplateProps struct {
	contracts.ComponentProps

	Title      string `json:"title"`
	CreateURL  string `json:"createURL,omitempty"`
	CreateText string `json:"createText,omitempty"`
	SearchURL  string `json:"searchURL,omitempty"`
}

// DetailTemplateProps defines properties for an entity detail page.
type DetailTemplateProps struct {
	contracts.ComponentProps

	Title     string   `json:"title"`
	BackURL   string   `json:"backURL,omitempty"`
	EditURL   string   `json:"editURL,omitempty"`
	DeleteURL string   `json:"deleteURL,omitempty"`
	Actions   []Action `json:"actions,omitempty"`
}

// Action represents a page-level action button.
type Action struct {
	Label   string `json:"label"`
	URL     string `json:"url"`
	Variant string `json:"variant,omitempty"` // primary, danger, etc.
	Icon    string `json:"icon,omitempty"`
	Confirm string `json:"confirm,omitempty"` // confirmation message
}

// SettingsTemplateProps defines properties for a settings page.
type SettingsTemplateProps struct {
	contracts.ComponentProps

	Title    string            `json:"title"`
	Sections []SettingsSection `json:"sections"`
}

// SettingsSection represents a section in settings.
type SettingsSection struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Icon  string `json:"icon,omitempty"`
}

// ErrorTemplateProps defines properties for an error page.
type ErrorTemplateProps struct {
	contracts.ComponentProps

	Code       int    `json:"code"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	ActionURL  string `json:"actionURL,omitempty"`
	ActionText string `json:"actionText,omitempty"`
}
