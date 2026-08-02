package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// TextProps defines the platform-agnostic properties for a Text component.
type TextProps struct {
	contracts.ComponentProps

	Content   string `json:"content"`
	Element   string `json:"element,omitempty"`   // p, span, div, strong, em, small, mark, del, ins, sub, sup, blockquote, code, pre, kbd, samp, var
	Size      string `json:"size,omitempty"`      // xs, sm, base, lg, xl, 2xl, 3xl, 4xl, 5xl
	Align     string `json:"align,omitempty"`     // left, center, right, justify
	Weight    string `json:"weight,omitempty"`    // thin, extralight, light, normal, medium, semibold, bold, extrabold, black
	Color     string `json:"color,omitempty"`     // primary, secondary, tertiary, muted, brand, success, warning, danger, info
	Transform string `json:"transform,omitempty"` // none, uppercase, lowercase, capitalize
	Truncate  bool   `json:"truncate,omitempty"`  // truncate with ellipsis
	NoWrap    bool   `json:"nowrap,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
	Underline bool   `json:"underline,omitempty"`
	Lines     int    `json:"lines,omitempty"` // line clamp, 1-6
}

// ToMap converts TextProps to map[string]any for unified Component construction.
func (p TextProps) ToMap() map[string]any { return propsToMap(p) }

// HeadingProps defines properties for heading elements (H1-H6).
type HeadingProps struct {
	contracts.ComponentProps

	Text     string `json:"text"`
	Level    int    `json:"level"`            // 1-6
	Anchor   string `json:"anchor,omitempty"` // optional anchor ID
	Truncate bool   `json:"truncate,omitempty"`
}

// ToMap converts HeadingProps to map[string]any for unified Component construction.
func (p HeadingProps) ToMap() map[string]any { return propsToMap(p) }

// LabelProps defines properties for form labels.
type LabelProps struct {
	contracts.ComponentProps

	Text     string `json:"text"`
	For      string `json:"for,omitempty"` // associated input ID
	Required bool   `json:"required,omitempty"`
}

// ToMap converts LabelProps to map[string]any for unified Component construction.
func (p LabelProps) ToMap() map[string]any { return propsToMap(p) }
