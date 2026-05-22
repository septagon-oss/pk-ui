package molecules

import "github.com/septagon-oss/platformkit-ui/contracts"

// TabsProps defines properties for a tabbed interface.
type TabsProps struct {
	contracts.ComponentProps

	Items       []TabItem `json:"items"`
	ActiveTab   string    `json:"activeTab,omitempty"`
	Orientation string    `json:"orientation,omitempty"` // horizontal, vertical
	Variant     string    `json:"variant,omitempty"`     // underline, pills
}

// ToMap converts TabsProps to map[string]any for unified Component construction.
func (p TabsProps) ToMap() map[string]any { return propsToMap(p) }

// TabItem represents a tab.
type TabItem struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Icon    string `json:"icon,omitempty"`
	Content string `json:"content,omitempty"` // for static content
	URL     string `json:"url,omitempty"`     // for lazy-loaded content (hx-get)
}

// AccordionProps defines properties for an accordion.
type AccordionProps struct {
	contracts.ComponentProps

	Items       []AccordionItem `json:"items"`
	Multiple    bool            `json:"multiple,omitempty"` // allow multiple open
	DefaultOpen []string        `json:"defaultOpen,omitempty"`
}

// ToMap converts AccordionProps to map[string]any for unified Component construction.
func (p AccordionProps) ToMap() map[string]any { return propsToMap(p) }

// AccordionItem represents an accordion section.
type AccordionItem struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

// BreadcrumbProps defines properties for breadcrumb navigation.
type BreadcrumbProps struct {
	contracts.ComponentProps

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

	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
	PerPage     int    `json:"perPage,omitempty"`
	Siblings    int    `json:"siblings,omitempty"` // pages shown around current
	BaseURL     string `json:"baseURL,omitempty"`
}

// ToMap converts PaginationProps to map[string]any for unified Component construction.
func (p PaginationProps) ToMap() map[string]any { return propsToMap(p) }

// StepperProps defines properties for a multi-step indicator.
type StepperProps struct {
	contracts.ComponentProps

	Steps       []StepItem `json:"steps"`
	CurrentStep int        `json:"currentStep"`
	Orientation string     `json:"orientation,omitempty"` // horizontal, vertical
	Clickable   bool       `json:"clickable,omitempty"`
}

// StepItem represents a step in a stepper.
type StepItem struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Status      string `json:"status,omitempty"` // pending, active, completed, error
}

// ToMap converts StepperProps to map[string]any for unified Component construction.
func (p StepperProps) ToMap() map[string]any { return propsToMap(p) }
