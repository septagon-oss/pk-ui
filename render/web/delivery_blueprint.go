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
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
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

func deliveryVisibleWhen(node blueprint.Node, conditions map[string]any) blueprint.Node {
	node.Props = maps.Clone(node.Props)
	if node.Props == nil {
		node.Props = make(map[string]any)
	}
	node.Props["visible_when"] = maps.Clone(conditions)
	return node
}

func cardProps(title string) molecules.CardProps {
	return molecules.CardProps{Title: title}
}

func textProps(content, color string) atoms.TextProps {
	return atoms.TextProps{
		Content: content,
		Size:    "sm",
		Color:   color,
	}
}
