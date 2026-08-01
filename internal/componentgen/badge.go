// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// badge.go is Badge expressed as an authored component — the first component
// to exist as one document rather than as four hand-maintained files. It is
// the round-trip proof: generating from this must reproduce the Badge that
// ships, and the generated code is checked by pk-ui's own gates rather than by
// review.
//
// Two defects in the hand-written Badge that the single source removes:
//
//  1. The Variant doc comment lists six values; the variant class table
//     defines eight ("default" and "danger" are undocumented, and "default" is
//     the renderer's fallback). The generated comment is derived from the same
//     Enum the table is derived from, so the two cannot disagree.
//  2. The delivery contract advertises seven variant values, omitting
//     "danger" that the class table styles. Same single-source fix.
//
// A third defect is deliberately preserved rather than fixed here: Size is
// declared, documented, and advertised in the delivery contract, but no class
// list or renderer branch consumes it. Removing it is a breaking contract
// change and belongs in its own reviewed commit, not smuggled in under a
// codegen proof. It is expressed faithfully so the proof compares like with
// like.

// BadgeSource returns Badge's renderer-neutral document.
func BadgeSource() Source {
	return Source{
		ID:             "Badge",
		Tier:           "Atom",
		Presentational: true,
		Description:    "Compact semantic status label.",
		LibraryPage:    "02 Atoms",
		LibrarySection: "Status",
		Props: []Prop{
			{
				Name: "label", GoName: "Label", Type: "string",
				Role: "content", Required: true,
				Description: "Visible badge label.",
			},
			{
				Name: "variant", GoName: "Variant", Type: "string",
				Role: "tone",
				Enum: []string{
					"default", "primary", "secondary", "success",
					"warning", "error", "danger", "info",
				},
				Default:     "default",
				Description: "Semantic status treatment.",
			},
			{
				Name: "size", GoName: "Size", Type: "string",
				Enum:        []string{"small", "medium", "large"},
				Default:     "medium",
				Description: "Badge scale.",
			},
			{
				Name: "dot", GoName: "Dot", Type: "bool",
				Role:        "modifier",
				Description: "Shows a leading status dot.",
			},
			{
				Name: "icon", GoName: "Icon", Type: "string",
				Role:        "content",
				Description: "Optional icon identity.",
			},
		},
		Examples: []Example{
			{
				Name:      "default",
				Canonical: true,
				Props:     map[string]any{"label": "Active", "variant": "primary"},
			},
			{
				Name:  "success",
				Props: map[string]any{"label": "Published", "variant": "secondary", "dot": true},
			},
		},
	}
}

// BadgeWebStyle returns Badge's web projection.
func BadgeWebStyle() WebStyle {
	return WebStyle{
		ID:      "Badge",
		Element: "span",
		Base: []Utility{
			{Method: "Display", Arg: "DisplayInlineFlex"},
			{Method: "Items", Arg: "ItemsCenter"},
			{Method: "Gap", Arg: "S1"},
			{Method: "Rounded", Arg: "RadiusFull"},
			{Method: "FontWeight", Arg: "FontMedium"},
			{Method: "PaddingX", Arg: "S2_5"},
			{Method: "PaddingY", Arg: "S0_5"},
			{Method: "FontSize", Arg: "TextXS"},
		},
		VariantProp: "variant",
		Variants: map[string][]Utility{
			"primary": {
				{Method: "Bg", Arg: "SurfaceBrandSoft"},
				{Method: "TextColor", Arg: "FgBrand"},
			},
			// NoUnderline is carried by the hand-written secondary variant.
			// It is preserved verbatim: the proof compares against what
			// ships, not against what would be tidier.
			"secondary": {
				{Method: "NoUnderline"},
				{Method: "Bg", Arg: "SurfaceTertiary"},
				{Method: "TextColor", Arg: "FgSecondary"},
			},
			"outline": {
				{Method: "Border", Arg: "Border1"},
				{Method: "BorderColor", Arg: "BorderPrimary"},
				{Method: "TextColor", Arg: "FgPrimary"},
			},
		},
		Parts: map[string][]Utility{
			"dot": {
				{Method: "Width", Arg: "S1_5"},
				{Method: "Height", Arg: "S1_5"},
				{Method: "Rounded", Arg: "RadiusFull"},
				{Method: "Bg", Arg: "FgBrand"},
			},
		},
		Body: []Step{
			{Kind: StepBaseAttrs},
			{Kind: StepClasses},
			{
				Kind: StepPartIfBool, Prop: "Dot", Part: "dot", Element: "span",
				Attrs: [][2]string{{"aria-hidden", "true"}},
			},
			{Kind: StepIconIfSet, Prop: "Icon"},
			{Kind: StepText, Prop: "Label", Part: "Label", Fallback: "Active"},
		},
	}
}
