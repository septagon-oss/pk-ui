// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_blueprint.go provides small provider-neutral constructors for the
// native design trees owned by the OSS web catalog. The trees use the exact
// utility classes composed by the production renderers; downstream Figma
// tooling resolves those utilities to native auto-layout and variable-bound
// properties rather than embedding browser images.

import (
	"maps"
	"strings"

	"github.com/septagon-oss/pk-design/pkg/blueprint"
)

func deliveryDesign(
	canonicalExample string,
	libraryPage string,
	librarySection string,
	root blueprint.Node,
) blueprint.Definition {
	return blueprint.Definition{
		SourceOfTruth:    blueprint.SourceRuntime,
		CanonicalExample: canonicalExample,
		ExampleVariants:  true,
		ExampleDefaults: &blueprint.ExampleContract{
			Viewport:  "desktop",
			Viewports: []string{"desktop", "mobile"},
			ColorMode: "light",
			Density:   "medium",
		},
		Taxonomy: &blueprint.Taxonomy{
			LibraryPage:    libraryPage,
			LibrarySection: librarySection,
		},
		Root: &root,
	}
}

func deliveryFrame(name, classes string, children ...blueprint.Node) blueprint.Node {
	return blueprint.Node{
		Kind:     blueprint.NodeFrame,
		Name:     name,
		Classes:  classes,
		Children: children,
	}
}

func deliverySemantics(node blueprint.Node, role, purpose string) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any, 2)
	}
	if role = strings.TrimSpace(role); role != "" {
		node.Props["semantic_role"] = role
	}
	if purpose = strings.TrimSpace(purpose); purpose != "" {
		node.Props["semantic_purpose"] = purpose
	}
	return node
}

func deliveryText(name, classes, fallback string, props ...string) blueprint.Node {
	node := blueprint.Node{
		Kind:    blueprint.NodeText,
		Name:    name,
		Classes: classes,
		Text:    fallback,
	}
	if len(props) > 0 {
		node.Props = map[string]any{"text_props": props}
	}
	return node
}

func deliverySlot(name, slot, classes string, children ...blueprint.Node) blueprint.Node {
	return blueprint.Node{
		Kind:     blueprint.NodeSlot,
		Name:     name,
		Slot:     slot,
		Classes:  classes,
		Children: children,
	}
}

// deliveryPreserveSlotFallback marks prop-bound fallback content that remains
// authoritative when an executable composition does not occupy the slot. A
// Modal title, for example, must not disappear merely because body/actions
// are supplied as rich child components.
func deliveryPreserveSlotFallback(node blueprint.Node) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["preserve_fallback_when_empty"] = true
	return node
}

func deliveryInstance(
	name string,
	componentType string,
	example string,
	classes string,
) blueprint.Node {
	node := blueprint.Node{
		Kind:    blueprint.NodeInstance,
		Name:    name,
		Text:    componentType,
		Classes: classes,
	}
	if example != "" {
		node.Props = map[string]any{"instance_example": example}
	}
	return node
}

func deliveryInstanceProps(
	name string,
	componentType string,
	example string,
	classes string,
	props map[string]any,
) blueprint.Node {
	node := deliveryInstance(name, componentType, example, classes)
	if len(props) == 0 {
		return node
	}
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["instance_props"] = maps.Clone(props)
	return node
}

func deliverySVGBound(
	name string,
	classes string,
	assetPrefix string,
	prop string,
) blueprint.Node {
	return blueprint.Node{
		Kind:     blueprint.NodeSVG,
		Name:     name,
		Classes:  classes,
		AssetRef: "bound:" + assetPrefix,
		Props: map[string]any{
			"svg_bound_prefix": assetPrefix,
			"bind_" + prop:     true,
		},
	}
}

func deliveryClassBound(
	node blueprint.Node,
	prop string,
	values map[string]string,
) blueprint.Node {
	node.ClassBindings = maps.Clone(node.ClassBindings)
	if node.ClassBindings == nil {
		node.ClassBindings = make(map[string]map[string]string)
	}
	node.ClassBindings[prop] = maps.Clone(values)
	order := make([]string, 0, len(node.ClassBindingOrder)+1)
	for _, existing := range node.ClassBindingOrder {
		if existing != prop {
			order = append(order, existing)
		}
	}
	node.ClassBindingOrder = append(order, prop)
	return node
}

func deliveryHiddenWhen(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["hidden_when"] = maps.Clone(conditions)
	return node
}

func deliveryHiddenWhenAny(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["hidden_when_any"] = maps.Clone(conditions)
	return node
}

func deliveryVisibleWhen(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["visible_when"] = maps.Clone(conditions)
	return node
}

func deliveryVisibleWhenAny(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["visible_when_any"] = maps.Clone(conditions)
	return node
}

// deliveryOptionalText keeps one stable native layer across variants, while
// clearing and hiding it when none of its authored content properties exist.
func deliveryOptionalText(node blueprint.Node) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["clear_when_unbound"] = true
	node.Props["hide_when_empty"] = true
	return node
}

// deliveryProgressBound identifies geometry or text derived from the same
// normalized value/max pair as the production Progress renderer.
func deliveryProgressBound(node blueprint.Node, formatText bool) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["progress_value_prop"] = "value"
	node.Props["progress_max_prop"] = "max"
	if formatText {
		node.Props["progress_format"] = "percent"
	}
	return node
}

// deliveryProgressWidthUnless declares the state in which percentage-derived
// fill width must yield to authored layout (for example, an indeterminate
// progress fill that occupies the complete track).
func deliveryProgressWidthUnless(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["progress_width_unless"] = maps.Clone(conditions)
	return node
}

// deliveryIntegerCapped keeps compact numeric UI faithful to renderer-side
// overflow formatting without replacing the underlying numeric property.
func deliveryIntegerCapped(node blueprint.Node, prop string, maximum int, overflow string) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["integer_value_prop"] = strings.TrimSpace(prop)
	node.Props["integer_max"] = maximum
	node.Props["integer_overflow_text"] = overflow
	return node
}
