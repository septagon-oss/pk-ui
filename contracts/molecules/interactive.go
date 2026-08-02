package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// DropdownProps defines properties for a dropdown/select component.
type DropdownProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string   `json:"name"`
	Options     []Option `json:"options"`
	Placeholder string   `json:"placeholder,omitempty"`
	AriaLabel   string   `json:"ariaLabel,omitempty"`
	SearchLabel string   `json:"searchLabel,omitempty"`
	ClearLabel  string   `json:"clearLabel,omitempty"`
	OpenLabel   string   `json:"openLabel,omitempty"`
	Searchable  bool     `json:"searchable,omitempty"`
	Clearable   bool     `json:"clearable,omitempty"`
	Multiple    bool     `json:"multiple,omitempty"`
	Value       string   `json:"value,omitempty"`
	Selected    []string `json:"selected,omitempty"`
	Size        string   `json:"size,omitempty"` // sm, md, lg
	Label       string   `json:"label,omitempty"`
}

// ToMap converts DropdownProps to map[string]any for unified Component construction.
func (p DropdownProps) ToMap() map[string]any { return propsToMap(p) }

// DatePickerProps defines properties for a date picker.
type DatePickerProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string `json:"name"`
	Value       string `json:"value,omitempty"` // ISO date string
	Min         string `json:"min,omitempty"`
	Max         string `json:"max,omitempty"`
	Format      string `json:"format,omitempty"` // optional consumer display metadata
	Label       string `json:"label,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	HelpText    string `json:"helpText,omitempty"`
	Error       string `json:"error,omitempty"`
	Invalid     bool   `json:"invalid,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ToMap converts DatePickerProps to map[string]any for unified Component construction.
func (p DatePickerProps) ToMap() map[string]any { return propsToMap(p) }

// FileUploadProps defines properties for a file upload component.
type FileUploadProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name           string `json:"name"`
	Accept         string `json:"accept,omitempty"` // MIME types and extensions
	Multiple       bool   `json:"multiple,omitempty"`
	Required       bool   `json:"required,omitempty"`
	MaxSize        int64  `json:"maxSize,omitempty"` // bytes
	DropZone       *bool  `json:"dropZone,omitempty"`
	Preview        bool   `json:"preview,omitempty"`
	ShowList       *bool  `json:"showList,omitempty"`
	Label          string `json:"label,omitempty"`
	HelpText       string `json:"helpText,omitempty"`
	AriaLabel      string `json:"ariaLabel,omitempty"`
	PromptLabel    string `json:"promptLabel,omitempty"`
	DropLabel      string `json:"dropLabel,omitempty"`
	ChooseLabel    string `json:"chooseLabel,omitempty"`
	LoadingLabel   string `json:"loadingLabel,omitempty"`
	RemoveLabel    string `json:"removeLabel,omitempty"`
	UploadedLabel  string `json:"uploadedLabel,omitempty"`
	MaxSizeLabel   string `json:"maxSizeLabel,omitempty"`
	MissingError   string `json:"missingError,omitempty"`
	TooLargeError  string `json:"tooLargeError,omitempty"`
	TypeError      string `json:"typeError,omitempty"`
	UploadURL      string `json:"uploadURL,omitempty"`
	UploadCategory string `json:"uploadCategory,omitempty"`
	Value          string `json:"value,omitempty"`
	CurrentName    string `json:"currentName,omitempty"`
}

// ToMap converts FileUploadProps to map[string]any for unified Component construction.
func (p FileUploadProps) ToMap() map[string]any { return propsToMap(p) }

// AutocompleteProps defines properties for an autocomplete/typeahead.
type AutocompleteProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name         string   `json:"name"`
	SearchURL    string   `json:"searchURL,omitempty"`
	QueryName    string   `json:"queryName,omitempty"` // server-search query parameter, default q
	Options      []Option `json:"options,omitempty"`
	MinChars     int      `json:"minChars,omitempty"` // default 2
	Debounce     int      `json:"debounce,omitempty"` // ms, default 300
	Value        string   `json:"value,omitempty"`
	DisplayValue string   `json:"displayValue,omitempty"`
	Placeholder  string   `json:"placeholder,omitempty"`
	Label        string   `json:"label,omitempty"`
	HelpText     string   `json:"helpText,omitempty"`
	Error        string   `json:"error,omitempty"`
	Invalid      bool     `json:"invalid,omitempty"`
	Required     bool     `json:"required,omitempty"`
}

// ToMap converts AutocompleteProps to map[string]any for unified Component construction.
func (p AutocompleteProps) ToMap() map[string]any { return propsToMap(p) }

// SearchBarProps defines properties for a search bar.
type SearchBarProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	SearchURL    string `json:"searchURL"`
	Label        string `json:"label,omitempty"` // accessible name for the input
	Placeholder  string `json:"placeholder,omitempty"`
	Name         string `json:"name,omitempty"`
	Value        string `json:"value,omitempty"`
	Instant      bool   `json:"instant,omitempty"` // search-as-you-type
	ShowClear    bool   `json:"showClear,omitempty"`
	ShowShortcut bool   `json:"showShortcut,omitempty"` // show keyboard shortcut hint
	DebounceMS   int    `json:"debounceMs,omitempty"`
	MinChars     int    `json:"minChars,omitempty"`
	ClearLabel   string `json:"clearLabel,omitempty"`
	ShortcutKey  string `json:"shortcutKey,omitempty"`
}

// ToMap converts SearchBarProps to map[string]any for unified Component construction.
func (p SearchBarProps) ToMap() map[string]any { return propsToMap(p) }
