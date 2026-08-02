// Implements: REQ-011.
// Per: ADR-0029.
// Discipline: C-14.

package organisms

import (
	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
)

// NavigationProps defines properties for the main navigation.
type NavigationProps struct {
	contracts.ComponentProps

	Brand    string                  `json:"brand,omitempty"`
	LogoURL  string                  `json:"logoURL,omitempty"`
	Items    []molecules.SidebarItem `json:"items"`
	UserMenu UserMenuProps           `json:"userMenu"`
}

// UserMenuProps defines the user menu in navigation.
type UserMenuProps struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatarURL,omitempty"`
	LogoutURL string `json:"logoutURL,omitempty"`
}

// DataGridProps defines the scalar and transport properties of a rich data
// grid. Search, filters, actions, table, status, and pagination are composed
// through the renderer's named slots rather than nested component props.
type DataGridProps struct {
	contracts.ComponentProps
	contracts.HTMXProps
}

// WindowedCollectionProps defines the bounded collection shell shared by web,
// adaptive, and native projections. Item content and navigation are slots.
type WindowedCollectionProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	State                 string `json:"state,omitempty"`
	ItemCount             int    `json:"itemCount,omitempty"`
	MaxItems              int    `json:"maxItems,omitempty"`
	CollectionLabel       string `json:"collectionLabel,omitempty"`
	EmptyTitle            string `json:"emptyTitle,omitempty"`
	EmptyDescription      string `json:"emptyDescription,omitempty"`
	LoadingLabel          string `json:"loadingLabel,omitempty"`
	ErrorTitle            string `json:"errorTitle,omitempty"`
	ErrorDescription      string `json:"errorDescription,omitempty"`
	RetryLabel            string `json:"retryLabel,omitempty"`
	RetryURL              string `json:"retryURL,omitempty"`
	NavigationUnavailable string `json:"navigationUnavailable,omitempty"`
	ContractError         bool   `json:"contractError,omitempty" delivery:"internal"`
}

// CommandPaletteProps defines properties for a Cmd+K command palette.
type CommandPaletteProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	SearchURL   string `json:"searchURL"`
	Placeholder string `json:"placeholder,omitempty"`
}

// DashboardWidgetProps defines properties for a dashboard stat/widget card.
type DashboardWidgetProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Title          string `json:"title"`
	Subtitle       string `json:"subtitle,omitempty"`
	Type           string `json:"type,omitempty"` // stat, chart, list, empty
	Value          string `json:"value,omitempty"`
	PreviousValue  string `json:"previousValue,omitempty"`
	Change         string `json:"change,omitempty"` // e.g., "+12%"
	Trend          string `json:"trend,omitempty"`  // up, down, flat
	Icon           string `json:"icon,omitempty"`
	RefreshURL     string `json:"refreshURL,omitempty"`
	RefreshOn      string `json:"refreshOn,omitempty"` // HTMX trigger event
	RefreshSeconds int    `json:"refreshSeconds,omitempty"`
	DetailURL      string `json:"detailURL,omitempty"`
	Span           string `json:"span,omitempty"` // grid columns: 1, 2, 3, 4
}

// WizardProps defines properties for a multi-step wizard/form.
type WizardProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Steps       []WizardStep `json:"steps"`
	CurrentStep int          `json:"currentStep"`
}

// WizardStep represents a wizard step.
type WizardStep struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}
