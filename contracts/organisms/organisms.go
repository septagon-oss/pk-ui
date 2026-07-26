// Implements: REQ-011.
// Per: ADR-0029.
// Discipline: C-14.

package organisms

import (
	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
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

// DataGridProps defines properties for a rich data grid (table + filters + pagination).
type DataGridProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Table      molecules.TableProps      `json:"table"`
	Pagination molecules.PaginationProps `json:"pagination"`
	Search     molecules.SearchBarProps  `json:"search,omitzero"`
	Actions    []atoms.ButtonProps       `json:"actions,omitempty"`
	SearchURL  string                    `json:"searchURL,omitempty"` // shorthand when Search is zero
	CreateURL  string                    `json:"createURL,omitempty"`
	CreateText string                    `json:"createText,omitempty"`
	Filters    []DataGridFilter          `json:"filters,omitempty"`
}

// DataGridFilter represents a filterable column.
type DataGridFilter struct {
	Key     string             `json:"key"`
	Label   string             `json:"label"`
	Type    string             `json:"type"` // text, select, date
	Options []molecules.Option `json:"options,omitempty"`
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

	Title     string `json:"title"`
	Value     string `json:"value,omitempty"`
	Change    string `json:"change,omitempty"` // e.g., "+12%"
	Trend     string `json:"trend,omitempty"`  // up, down, flat
	Icon      string `json:"icon,omitempty"`
	RefreshOn string `json:"refreshOn,omitempty"` // HTMX trigger event
	DetailURL string `json:"detailURL,omitempty"`
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
