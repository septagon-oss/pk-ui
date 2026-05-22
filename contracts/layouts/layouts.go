package layouts

import "github.com/septagon-oss/platformkit-ui/contracts"

// GridProps defines properties for a CSS Grid layout.
type GridProps struct {
	contracts.ComponentProps

	Columns string `json:"columns,omitempty"` // Tailwind grid-cols value
	Gap     string `json:"gap,omitempty"`
}

// StackProps defines properties for a vertical stack layout.
type StackProps struct {
	contracts.ComponentProps

	Gap   string `json:"gap,omitempty"`
	Align string `json:"align,omitempty"` // start, center, end, stretch
}

// FlexProps defines properties for a flexbox layout.
type FlexProps struct {
	contracts.ComponentProps

	Direction string `json:"direction,omitempty"` // row, column
	Wrap      bool   `json:"wrap,omitempty"`
	Gap       string `json:"gap,omitempty"`
	Align     string `json:"align,omitempty"`   // start, center, end, stretch
	Justify   string `json:"justify,omitempty"` // start, center, end, between, around
}

// ContainerProps defines properties for a centered container.
type ContainerProps struct {
	contracts.ComponentProps

	MaxWidth string `json:"maxWidth,omitempty"` // sm, md, lg, xl, 2xl, full
	Padding  string `json:"padding,omitempty"`
}
