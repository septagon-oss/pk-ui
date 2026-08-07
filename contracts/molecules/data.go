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

	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Image         string `json:"image,omitempty"`
	ImageAlt      string `json:"imageAlt,omitempty"`
	ImagePosition string `json:"imagePosition,omitempty"` // top, bottom, left, right
	Variant       string `json:"variant,omitempty"`       // default, elevated, outlined, plain
	Padding       string `json:"padding,omitempty"`       // none, small, medium, large
	Shadow        string `json:"shadow,omitempty"`        // none, small, medium, large
	Clickable     bool   `json:"clickable,omitempty"`
	Hoverable     bool   `json:"hoverable,omitempty"`
	Href          string `json:"href,omitempty"`
}

// ToMap converts CardProps to map[string]any for unified Component construction.
func (p CardProps) ToMap() map[string]any { return propsToMap(p) }

// ModalProps defines platform-agnostic properties for a Modal component.
type ModalProps struct {
	contracts.ComponentProps

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Body        string `json:"body,omitempty"`
	Footer      string `json:"footer,omitempty"`
	AriaLabel   string `json:"ariaLabel,omitempty"`
	CloseLabel  string `json:"closeLabel,omitempty"`
	Size        string `json:"size,omitempty"` // small, medium, large, xl, full

	// Pointer booleans preserve the intended default-true behavior while still
	// allowing portable clients to explicitly disable an affordance.
	Closable       *bool `json:"closable,omitempty"`
	CloseOnOverlay *bool `json:"closeOnOverlay,omitempty"`
	CloseOnEscape  *bool `json:"closeOnEscape,omitempty"`
	ShowClose      *bool `json:"showClose,omitempty"`
	ShowOverlay    *bool `json:"showOverlay,omitempty"`
	Centered       *bool `json:"centered,omitempty"`
	ClearOnClose   *bool `json:"clearOnClose,omitempty"`
	Open           bool  `json:"open,omitempty"`
	OpenOnSwap     bool  `json:"openOnSwap,omitempty"`
	Deferred       bool  `json:"deferred,omitempty"`
}

// ToMap converts ModalProps to map[string]any for unified Component construction.
func (p ModalProps) ToMap() map[string]any { return propsToMap(p) }

// SidebarProps defines properties for a sidebar navigation.
type SidebarProps struct {
	contracts.ComponentProps

	Items           []SidebarItem    `json:"items,omitempty"`
	Sections        []SidebarSection `json:"sections,omitempty"`
	Current         string           `json:"current,omitempty"`
	Flavor          string           `json:"flavor,omitempty"` // admin, content
	Collapsible     bool             `json:"collapsible,omitempty"`
	Collapsed       bool             `json:"collapsed,omitempty"`
	NavigationLabel string           `json:"navigationLabel,omitempty"`
	BrandLabel      string           `json:"brandLabel,omitempty"`
	BrandHref       string           `json:"brandHref,omitempty"`
}

// SidebarItem represents a sidebar navigation item.
type SidebarItem struct {
	ID           string            `json:"id,omitempty"`
	Label        string            `json:"label"`
	Href         string            `json:"href,omitempty"`
	Icon         string            `json:"icon,omitempty"`
	Prefix       string            `json:"prefix,omitempty"`
	Badge        string            `json:"badge,omitempty"`
	BadgeVariant string            `json:"badgeVariant,omitempty"`
	Active       bool              `json:"active,omitempty"`
	Disabled     bool              `json:"disabled,omitempty"`
	SearchText   string            `json:"searchText,omitempty"`
	Attrs        map[string]string `json:"-" delivery:"internal"`
	Children     []SidebarItem     `json:"children,omitempty"`
}

// SidebarSection groups related sidebar items under an optional heading.
type SidebarSection struct {
	ID         string        `json:"id,omitempty"`
	Label      string        `json:"label,omitempty"`
	Glyph      string        `json:"glyph,omitempty"`
	Tone       string        `json:"tone,omitempty"` // neutral, brand, success, warning, danger, info
	SearchText string        `json:"searchText,omitempty"`
	Items      []SidebarItem `json:"items,omitempty"`
}

// ToMap converts SidebarProps to map[string]any for unified Component construction.
func (p SidebarProps) ToMap() map[string]any { return propsToMap(p) }
