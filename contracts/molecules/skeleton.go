package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// TableSkeletonProps defines the loading rendering of a Table: the same wrap,
// header, and cell classes with pulsing placeholders where data will land.
type TableSkeletonProps struct {
	contracts.ComponentProps

	Columns int  `json:"columns,omitempty"` // header/cell count (default 4)
	Rows    int  `json:"rows,omitempty"`    // placeholder row count (default 3)
	Compact bool `json:"compact,omitempty"`
}

// ToMap converts TableSkeletonProps to map[string]any for unified Component construction.
func (p TableSkeletonProps) ToMap() map[string]any { return propsToMap(p) }

// CardSkeletonProps defines the loading rendering of a Card: title and body
// placeholders inside the real card frame.
type CardSkeletonProps struct {
	contracts.ComponentProps

	Lines int `json:"lines,omitempty"` // body line count (default 3)
}

// ToMap converts CardSkeletonProps to map[string]any for unified Component construction.
func (p CardSkeletonProps) ToMap() map[string]any { return propsToMap(p) }
