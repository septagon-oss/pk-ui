package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// TableProps defines platform-agnostic properties for a Table component.
type TableProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Columns    []TableColumn `json:"columns"`
	Rows       []TableRow    `json:"rows,omitempty"`
	Sortable   bool          `json:"sortable,omitempty"`
	Selectable bool          `json:"selectable,omitempty"`
	Striped    bool          `json:"striped,omitempty"`
	Compact    bool          `json:"compact,omitempty"`
	EmptyText  string        `json:"emptyText,omitempty"`
}

// ToMap converts TableProps to map[string]any for unified Component construction.
func (p TableProps) ToMap() map[string]any { return propsToMap(p) }

// TableColumn defines a table column.
type TableColumn struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Sortable bool   `json:"sortable,omitempty"`
	Primary  bool   `json:"primary,omitempty"` // emphasized identity cell
	Width    string `json:"width,omitempty"`
	Align    string `json:"align,omitempty"` // left, center, right
}

// TableRow represents a table data row.
type TableRow struct {
	ID    string         `json:"id,omitempty"`
	Cells map[string]any `json:"cells"`
}

// CardProps defines platform-agnostic properties for a Card component.
type CardProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	Variant     string `json:"variant,omitempty"` // default, elevated, outlined
	Clickable   bool   `json:"clickable,omitempty"`
	Href        string `json:"href,omitempty"`
}

// ToMap converts CardProps to map[string]any for unified Component construction.
func (p CardProps) ToMap() map[string]any { return propsToMap(p) }

// ModalProps defines platform-agnostic properties for a Modal component.
type ModalProps struct {
	contracts.ComponentProps

	Title        string `json:"title,omitempty"`
	Size         string `json:"size,omitempty"` // small, medium, large, fullscreen
	Closable     bool   `json:"closable,omitempty"`
	CloseOnClick bool   `json:"closeOnClick,omitempty"` // close on backdrop click
}

// ToMap converts ModalProps to map[string]any for unified Component construction.
func (p ModalProps) ToMap() map[string]any { return propsToMap(p) }

// SidebarProps defines properties for a sidebar navigation.
type SidebarProps struct {
	contracts.ComponentProps

	Items     []SidebarItem `json:"items"`
	Collapsed bool          `json:"collapsed,omitempty"`
}

// SidebarItem represents a sidebar navigation item.
type SidebarItem struct {
	Label    string        `json:"label"`
	Href     string        `json:"href,omitempty"`
	Icon     string        `json:"icon,omitempty"`
	Active   bool          `json:"active,omitempty"`
	Children []SidebarItem `json:"children,omitempty"`
}

// ToMap converts SidebarProps to map[string]any for unified Component construction.
func (p SidebarProps) ToMap() map[string]any { return propsToMap(p) }
