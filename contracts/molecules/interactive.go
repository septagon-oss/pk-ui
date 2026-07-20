package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/platformkit-ui/contracts"

// DropdownProps defines properties for a dropdown/select component.
type DropdownProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string   `json:"name"`
	Options     []Option `json:"options"`
	Placeholder string   `json:"placeholder,omitempty"`
	Searchable  bool     `json:"searchable,omitempty"`
	Clearable   bool     `json:"clearable,omitempty"`
	Multiple    bool     `json:"multiple,omitempty"`
	Value       string   `json:"value,omitempty"`
	Size        string   `json:"size,omitempty"` // small, medium, large
	Label       string   `json:"label,omitempty"`
}

// ToMap converts DropdownProps to map[string]any for unified Component construction.
func (p DropdownProps) ToMap() map[string]any { return propsToMap(p) }

// DatePickerProps defines properties for a date picker.
type DatePickerProps struct {
	contracts.ComponentProps

	Name   string `json:"name"`
	Value  string `json:"value,omitempty"` // ISO date string
	Min    string `json:"min,omitempty"`
	Max    string `json:"max,omitempty"`
	Format string `json:"format,omitempty"` // display format
	Label  string `json:"label,omitempty"`
}

// ToMap converts DatePickerProps to map[string]any for unified Component construction.
func (p DatePickerProps) ToMap() map[string]any { return propsToMap(p) }

// FileUploadProps defines properties for a file upload component.
type FileUploadProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Accept   string `json:"accept,omitempty"` // MIME types
	Multiple bool   `json:"multiple,omitempty"`
	MaxSize  int64  `json:"maxSize,omitempty"` // bytes
	DropZone bool   `json:"dropZone,omitempty"`
	Preview  bool   `json:"preview,omitempty"`
	Label    string `json:"label,omitempty"`
}

// ToMap converts FileUploadProps to map[string]any for unified Component construction.
func (p FileUploadProps) ToMap() map[string]any { return propsToMap(p) }

// AutocompleteProps defines properties for an autocomplete/typeahead.
type AutocompleteProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name      string `json:"name"`
	SearchURL string `json:"searchURL"`
	MinChars  int    `json:"minChars,omitempty"` // default 2
	Debounce  int    `json:"debounce,omitempty"` // ms, default 300
	Label     string `json:"label,omitempty"`
}

// ToMap converts AutocompleteProps to map[string]any for unified Component construction.
func (p AutocompleteProps) ToMap() map[string]any { return propsToMap(p) }

// SearchBarProps defines properties for a search bar.
type SearchBarProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	SearchURL    string `json:"searchURL"`
	Placeholder  string `json:"placeholder,omitempty"`
	Instant      bool   `json:"instant,omitempty"` // search-as-you-type
	ShowClear    bool   `json:"showClear,omitempty"`
	ShowShortcut bool   `json:"showShortcut,omitempty"` // show keyboard shortcut hint
}

// ToMap converts SearchBarProps to map[string]any for unified Component construction.
func (p SearchBarProps) ToMap() map[string]any { return propsToMap(p) }
