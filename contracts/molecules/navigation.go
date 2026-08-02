package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// TabsProps defines properties for a tabbed interface.
type TabsProps struct {
	contracts.ComponentProps

	Items        []TabItem `json:"items,omitempty"`
	ActiveTab    string    `json:"activeTab,omitempty"`
	Orientation  string    `json:"orientation,omitempty"` // horizontal, vertical
	Variant      string    `json:"variant,omitempty"`     // underline, pills
	HxGet        string    `json:"hxGet,omitempty"`       // default lazy-panel endpoint
	LoadingLabel string    `json:"loadingLabel,omitempty"`
}

// ToMap converts TabsProps to map[string]any for unified Component construction.
func (p TabsProps) ToMap() map[string]any { return propsToMap(p) }

// TabItem represents a tab.
type TabItem struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Icon     string `json:"icon,omitempty"`
	Badge    string `json:"badge,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Content  string `json:"content,omitempty"` // static panel content for direct rendering
	URL      string `json:"url,omitempty"`     // navigation target in item mode
	HxGet    string `json:"hxGet,omitempty"`   // lazy-panel endpoint in panel mode
}

// AccordionProps defines properties for an accordion.
type AccordionProps struct {
	contracts.ComponentProps

	Items       []AccordionItem `json:"items,omitempty"`
	Multiple    bool            `json:"multiple,omitempty"`    // allow multiple open
	DefaultOpen []string        `json:"defaultOpen,omitempty"` // item IDs open on first render
	Bordered    *bool           `json:"bordered,omitempty"`    // nil defaults to true
	Flush       bool            `json:"flush,omitempty"`       // remove outer surface chrome
}

// ToMap converts AccordionProps to map[string]any for unified Component construction.
func (p AccordionProps) ToMap() map[string]any { return propsToMap(p) }

// AccordionItem represents an accordion section.
type AccordionItem struct {
	ID          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"` // compatibility alias for portable clients
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Content     string `json:"content,omitempty"`
	DefaultOpen bool   `json:"defaultOpen,omitempty"`
	Open        bool   `json:"open,omitempty"` // compatibility alias for JSON-driven surfaces
	Disabled    bool   `json:"disabled,omitempty"`
}

// BreadcrumbProps defines properties for breadcrumb navigation.
type BreadcrumbProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Items     []BreadcrumbItem `json:"items"`
	Separator string           `json:"separator,omitempty"` // default "/"
	MaxItems  int              `json:"maxItems,omitempty"`  // collapse middle items
}

// ToMap converts BreadcrumbProps to map[string]any for unified Component construction.
func (p BreadcrumbProps) ToMap() map[string]any { return propsToMap(p) }

// BreadcrumbItem represents a breadcrumb segment.
type BreadcrumbItem struct {
	Label   string `json:"label"`
	Href    string `json:"href,omitempty"` // empty = current page
	Icon    string `json:"icon,omitempty"`
	Current bool   `json:"current,omitempty"`
}

// PaginationProps defines properties for pagination controls.
type PaginationProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	CurrentPage     int    `json:"currentPage"`
	TotalPages      int    `json:"totalPages"`
	PerPage         int    `json:"perPage,omitempty"`
	Siblings        int    `json:"siblings,omitempty"` // pages shown around current
	BaseURL         string `json:"baseURL,omitempty"`
	CursorMode      string `json:"cursorMode,omitempty"` // previous-next, load-more
	PreviousCursor  string `json:"previousCursor,omitempty"`
	NextCursor      string `json:"nextCursor,omitempty"`
	BeforeParameter string `json:"beforeParameter,omitempty"`
	AfterParameter  string `json:"afterParameter,omitempty"`
	PreviousURL     string `json:"previousURL,omitempty"`
	NextURL         string `json:"nextURL,omitempty"`
	PreviousLabel   string `json:"previousLabel,omitempty"`
	NextLabel       string `json:"nextLabel,omitempty"`
	LoadMoreLabel   string `json:"loadMoreLabel,omitempty"`
	NavigationLabel string `json:"navigationLabel,omitempty"`
}

// ToMap converts PaginationProps to map[string]any for unified Component construction.
func (p PaginationProps) ToMap() map[string]any { return propsToMap(p) }

// StepperProps defines properties for a multi-step indicator.
type StepperProps struct {
	contracts.ComponentProps

	Steps           []StepItem `json:"steps"`
	CurrentStep     int        `json:"currentStep"`
	Orientation     string     `json:"orientation,omitempty"` // horizontal, vertical
	Clickable       bool       `json:"clickable,omitempty"`
	Compact         bool       `json:"compact,omitempty"`
	StepAction      string     `json:"stepAction,omitempty"`      // controller action for whole-step activation
	NavigationLabel string     `json:"navigationLabel,omitempty"` // accessible nav label
}

// StepItem represents a step in a stepper.
type StepItem struct {
	Key         string `json:"key,omitempty"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Status      string `json:"status,omitempty"` // pending, active, completed, error
}

// ToMap converts StepperProps to map[string]any for unified Component construction.
func (p StepperProps) ToMap() map[string]any { return propsToMap(p) }
