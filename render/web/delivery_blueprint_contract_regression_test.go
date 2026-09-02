// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

import (
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-design/pkg/blueprint"
	"github.com/septagon-oss/tw"
)

func TestCoreDeliveryBlueprintClassBindingOrderIsCompleteAndUnique(t *testing.T) {
	t.Parallel()

	for _, definition := range OSSDeliveryCatalog() {
		component := string(definition.Identity.ID)
		blueprintContractWalk(definition.Design.Root, func(path string, node *blueprint.Node) {
			if len(node.ClassBindings) == 0 {
				return
			}
			if len(node.ClassBindingOrder) != len(node.ClassBindings) {
				t.Errorf(
					"%s/%s binding order = %v, want one entry for each of %v",
					component,
					path,
					node.ClassBindingOrder,
					node.ClassBindings,
				)
			}

			seen := make(map[string]struct{}, len(node.ClassBindingOrder))
			for _, property := range node.ClassBindingOrder {
				if _, duplicate := seen[property]; duplicate {
					t.Errorf("%s/%s repeats binding %q in order %v", component, path, property, node.ClassBindingOrder)
					continue
				}
				seen[property] = struct{}{}
				if _, exists := node.ClassBindings[property]; !exists {
					t.Errorf("%s/%s orders absent binding %q", component, path, property)
				}
			}
			for property := range node.ClassBindings {
				if _, ordered := seen[property]; !ordered {
					t.Errorf("%s/%s binding %q is absent from order %v", component, path, property, node.ClassBindingOrder)
				}
			}
		})
	}
}

func TestCoreDeliveryBlueprintIconInstancesUseQualifiedSelectors(t *testing.T) {
	t.Parallel()

	icon := blueprintContractDefinition(t, "Icon")
	selectors := make(map[string]struct{}, len(icon.Examples))
	for _, example := range icon.Examples {
		selectors[example.Name] = struct{}{}
	}

	for _, definition := range OSSDeliveryCatalog() {
		component := string(definition.Identity.ID)
		blueprintContractWalk(definition.Design.Root, func(path string, node *blueprint.Node) {
			if node.Kind != blueprint.NodeInstance || node.Text != "Icon" {
				return
			}
			selector, _ := node.Props["instance_example"].(string)
			selector = strings.TrimSpace(selector)
			if selector == "" {
				boundProperty, _ := node.Props["instance_example_bound_prop"].(string)
				if strings.TrimSpace(boundProperty) == "" {
					t.Errorf("%s/%s has a bare Icon instance selector", component, path)
				}
				return
			}
			if _, exists := selectors[selector]; !exists {
				t.Errorf("%s/%s selects unknown Icon example %q", component, path, selector)
			}
		})
	}
}

func TestButtonBlueprintPreservesCascadeAndLoadingReplacement(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Button")
	root := definition.Design.Root
	blueprintContractAssertOrder(t, root, "variant", "tone", "size", "fullWidth", "iconOnly")

	loading := blueprintContractOnlyNodeNamed(t, root, "LoadingIndicator")
	if loading.Kind != blueprint.NodeInstance || loading.Text != "Spinner" {
		t.Fatalf("Button/LoadingIndicator = %#v, want Spinner instance", loading)
	}
	if got := loading.Props["instance_example"]; got != "small" {
		t.Errorf("Button/LoadingIndicator selector = %#v, want small", got)
	}
	blueprintContractAssertCondition(t, loading, "visible_when", map[string]any{"loading": true})

	label := blueprintContractOnlyNodeNamed(t, root, "Label")
	blueprintContractAssertCondition(t, label, "hidden_when", map[string]any{"iconOnly": true})
	if hidden, _ := label.Props["hidden_when"].(map[string]any); hidden["loading"] != nil {
		t.Errorf("Button/Label is hidden while loading: %v", hidden)
	}
}

func TestBadgeBlueprintBindsDotToneAndKeepsCountOptional(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Badge")
	root := definition.Design.Root
	blueprintContractAssertOrder(t, root, "variant", "tone", "size")

	dot := blueprintContractOnlyNodeNamed(t, root, "StatusDot")
	blueprintContractAssertOrder(t, dot, "tone")
	if got, want := dot.ClassBindings["tone"], compileClassMap(clBadgeDotTone); !maps.Equal(got, want) {
		t.Errorf("Badge/StatusDot tone classes = %v, want %v", got, want)
	}
	blueprintContractAssertCondition(t, dot, "visible_when", map[string]any{"dot": true})

	count := blueprintContractOnlyNodeNamed(t, root, "Count")
	if count.Kind != blueprint.NodeText || count.Text != "" {
		t.Errorf("Badge/Count = %#v, want empty optional text", count)
	}
	if got := blueprintContractStringValues(count.Props["text_props"]); !slices.Equal(got, []string{"count"}) {
		t.Errorf("Badge/Count text props = %v, want [count]", got)
	}
	for _, property := range []string{"clear_when_unbound", "hide_when_empty"} {
		if enabled, _ := count.Props[property].(bool); !enabled {
			t.Errorf("Badge/Count %s = %#v, want true", property, count.Props[property])
		}
	}
}

func TestInputBlueprintPreservesBindingCascadeAndConditionalSources(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Input")
	root := definition.Design.Root
	if got, want := root.Classes, clFieldWrap.Merge(clFieldWrapFull).Compile(); got != want {
		t.Errorf("Input root classes = %q, want full-width field source %q", got, want)
	}

	control := blueprintContractOnlyNodeNamed(t, root, "Control")
	for _, class := range []string{"flex", "items-center", "gap-2"} {
		if !slices.Contains(strings.Fields(control.Classes), class) {
			t.Errorf("Input/Control classes = %q, want %q", control.Classes, class)
		}
	}
	blueprintContractAssertOrder(t, control, "size", "tone", "error", "invalid", "readOnly")
	blueprintContractAssertClassBinding(t, control, "error", map[string]string{
		"": "", "*": clInputError.Compile(),
	})
	blueprintContractAssertClassBinding(t, control, "invalid", map[string]string{
		"false": "", "true": clInputError.Compile(),
	})
	blueprintContractAssertClassBinding(t, control, "readOnly", map[string]string{
		"false": "", "true": clInputReadOnly.Compile(),
	})

	value := blueprintContractOnlyNodeNamed(t, root, "Value")
	blueprintContractAssertOrder(t, value, "value")
	blueprintContractAssertClassBinding(t, value, "value", map[string]string{
		"":  tw.New().TextColor(tw.FgPlaceholder).Compile(),
		"*": tw.New().TextColor(tw.FgPrimary).Compile(),
	})
	if got, want := value.Props["mask_when"], (map[string]any{"type": "password"}); !reflect.DeepEqual(got, want) {
		t.Errorf("Input/Value mask_when = %#v, want %#v", got, want)
	}

	help := blueprintContractOnlyNodeNamed(t, root, "Help")
	blueprintContractAssertOrder(t, help, "error")
	blueprintContractAssertClassBinding(t, help, "error", map[string]string{
		"": "", "*": clFieldErr.Compile(),
	})
}

func TestDividerBlueprintOwnsThreeLayerLabelledStructure(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Divider")
	root := definition.Design.Root
	blueprintContractAssertOrder(t, root, "orientation", "text")
	if got := len(root.Children); got != 3 {
		t.Fatalf("Divider direct child count = %d, want three labelled layers", got)
	}
	names := make([]string, 0, len(root.Children))
	for index := range root.Children {
		names = append(names, root.Children[index].Name)
	}
	if want := []string{"LineStart", "Label", "LineEnd"}; !slices.Equal(names, want) {
		t.Fatalf("Divider direct children = %v, want %v", names, want)
	}
	for index := range root.Children {
		blueprintContractAssertCondition(t, &root.Children[index], "visible_when", map[string]any{"text": "*"})
	}
	for _, name := range []string{"LineStart", "LineEnd"} {
		line := blueprintContractOnlyNodeNamed(t, root, name)
		for _, class := range []string{"h-px", "bg-border-primary"} {
			if !slices.Contains(strings.Fields(line.Classes), class) {
				t.Errorf("Divider/%s classes = %q, want %q", name, line.Classes, class)
			}
		}
	}

	label := blueprintContractOnlyNodeNamed(t, root, "Label")
	if label.Kind != blueprint.NodeText || label.Text != "" {
		t.Errorf("Divider/Label = %#v, want empty optional text", label)
	}
	if got := blueprintContractStringValues(label.Props["text_props"]); !slices.Equal(got, []string{"text"}) {
		t.Errorf("Divider/Label text props = %v, want [text]", got)
	}
	for _, property := range []string{"clear_when_unbound", "hide_when_empty"} {
		if enabled, _ := label.Props[property].(bool); !enabled {
			t.Errorf("Divider/Label %s = %#v, want true", property, label.Props[property])
		}
	}
}

func TestProgressBlueprintUsesOneDerivedFill(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Progress")
	root := definition.Design.Root
	fills := blueprintContractNodesNamed(root, "Fill")
	if len(fills) != 1 {
		t.Fatalf("Progress Fill count = %d, want exactly one", len(fills))
	}
	if stale := blueprintContractNodesNamed(root, "IndeterminateFill"); len(stale) != 0 {
		t.Fatalf("Progress has %d stale IndeterminateFill nodes", len(stale))
	}
	fill := fills[0]
	blueprintContractAssertOrder(t, fill, "tone", "indeterminate")
	blueprintContractAssertDerivedProgress(t, fill, false)

	value := blueprintContractOnlyNodeNamed(t, root, "Value")
	blueprintContractAssertDerivedProgress(t, value, true)
	blueprintContractAssertCondition(t, value, "visible_when", map[string]any{"showText": true})
	blueprintContractAssertCondition(t, value, "hidden_when", map[string]any{"indeterminate": true})
}

func TestSpinnerBlueprintOwnsTrackAndTintableArc(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Spinner")
	root := definition.Design.Root
	blueprintContractAssertOrder(t, root, "size")
	track := blueprintContractOnlyNodeNamed(t, root, "Track")
	arc := blueprintContractOnlyNodeNamed(t, root, "Arc")
	if arc.Kind != blueprint.NodeSVG || arc.AssetRef != "asset:spinner-arc" {
		t.Fatalf("Spinner/Arc = %#v, want spinner SVG asset", arc)
	}
	blueprintContractAssertOrder(t, track, "size")
	blueprintContractAssertOrder(t, arc, "size", "tone")
	if len(definition.Design.Assets) != 1 || definition.Design.Assets[0].Name != "spinner-arc" {
		t.Fatalf("Spinner assets = %#v, want spinner-arc", definition.Design.Assets)
	}
}

func TestAccordionAndModalBlueprintInstanceAndBindingPlacement(t *testing.T) {
	t.Parallel()

	accordion := blueprintContractDefinition(t, "Accordion")
	chevron := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Chevron")
	if chevron.Kind != blueprint.NodeInstance || chevron.Text != "Icon" {
		t.Fatalf("Accordion/Chevron = %#v, want Icon instance", chevron)
	}
	if got := chevron.Props["instance_example"]; got != "Chevron Down / Fg Tertiary / 20" {
		t.Errorf("Accordion/Chevron selector = %#v, want qualified Chevron Down selector", got)
	}
	sections := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Sections")
	if sections.Kind != blueprint.NodeSlot || sections.Slot != "section" {
		t.Fatalf("Accordion/Sections = %#v, want public section slot", sections)
	}
	section := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Section")
	if got := section.Props["repeat_for_slot"]; got != "section" {
		t.Errorf("Accordion/Section repeat_for_slot = %#v, want section", got)
	}
	if got := section.Props["repeat_items_prop"]; got != "items" {
		t.Errorf("Accordion/Section repeat_items_prop = %#v, want items", got)
	}
	separator := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Separator")
	blueprintContractAssertCondition(t, separator, "visible_when", map[string]any{"separator": true})
	for _, class := range []string{"w-full", "h-px", "bg-border-primary"} {
		if !slices.Contains(strings.Fields(separator.Classes), class) {
			t.Errorf("Accordion/Separator classes = %q, want %q", separator.Classes, class)
		}
	}
	content := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Content")
	if content.Kind != blueprint.NodeFrame {
		t.Fatalf("Accordion/Content = %#v, want internal frame inside repeated section", content)
	}
	blueprintContractAssertCondition(t, content, "visible_when", map[string]any{"open": true})
	trigger := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Trigger")
	blueprintContractAssertOrder(t, trigger, "disabled")
	blueprintContractAssertOrder(t, chevron, "open")
	title := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Title")
	if got := blueprintContractStringValues(title.Props["text_props"]); !slices.Equal(got, []string{"title"}) {
		t.Errorf("Accordion/Title text props = %v, want [title]", got)
	}
	subtitle := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "Subtitle")
	for _, property := range []string{"clear_when_unbound", "hide_when_empty"} {
		if enabled, _ := subtitle.Props[property].(bool); !enabled {
			t.Errorf("Accordion/Subtitle %s = %#v, want true", property, subtitle.Props[property])
		}
	}
	leadingIcon := blueprintContractOnlyNodeNamed(t, accordion.Design.Root, "LeadingIcon")
	if got := leadingIcon.Props["instance_example_bound_prop"]; got != "icon" {
		t.Errorf("Accordion/LeadingIcon instance bound prop = %#v, want icon", got)
	}

	modal := blueprintContractDefinition(t, "Modal")
	root := modal.Design.Root
	blueprintContractAssertOrder(t, root, "centered")
	blueprintContractAssertClassBinding(t, root, "centered", map[string]string{
		"false": clModalBottomSheet.Compile(),
		"true":  clModalCentered.Compile(),
	})
	if _, misplaced := root.ClassBindings["size"]; misplaced {
		t.Error("Modal root owns size binding; want it on Panel")
	}

	panel := blueprintContractOnlyNodeNamed(t, root, "Panel")
	blueprintContractAssertOrder(t, panel, "size")
	blueprintContractAssertClassBinding(t, panel, "size", compileClassMap(clModalPanelSize))
	if _, misplaced := panel.ClassBindings["centered"]; misplaced {
		t.Error("Modal/Panel owns centered binding; want it on root")
	}
	headerContent := blueprintContractOnlyNodeNamed(t, root, "HeaderContent")
	if preserved, _ := headerContent.Props["preserve_fallback_when_empty"].(bool); !preserved {
		t.Errorf("Modal/HeaderContent preserve_fallback_when_empty = %#v, want true", headerContent.Props["preserve_fallback_when_empty"])
	}
}

func blueprintContractDefinition(t *testing.T, component string) DeliveryDefinition {
	t.Helper()
	for _, definition := range OSSDeliveryCatalog() {
		if string(definition.Identity.ID) == component {
			return definition
		}
	}
	t.Fatalf("delivery catalog has no %s definition", component)
	return DeliveryDefinition{}
}

func blueprintContractWalk(root *blueprint.Node, visit func(string, *blueprint.Node)) {
	var walk func(*blueprint.Node, string)
	walk = func(node *blueprint.Node, parent string) {
		if node == nil {
			return
		}
		path := node.Name
		if parent != "" {
			path = parent + "/" + node.Name
		}
		visit(path, node)
		for index := range node.Children {
			walk(&node.Children[index], path)
		}
	}
	walk(root, "")
}

func blueprintContractNodesNamed(root *blueprint.Node, name string) []*blueprint.Node {
	var matches []*blueprint.Node
	blueprintContractWalk(root, func(_ string, node *blueprint.Node) {
		if node.Name == name {
			matches = append(matches, node)
		}
	})
	return matches
}

func blueprintContractOnlyNodeNamed(t *testing.T, root *blueprint.Node, name string) *blueprint.Node {
	t.Helper()
	matches := blueprintContractNodesNamed(root, name)
	if len(matches) != 1 {
		t.Fatalf("blueprint node %q count = %d, want exactly one", name, len(matches))
	}
	return matches[0]
}

func blueprintContractAssertOrder(t *testing.T, node *blueprint.Node, want ...string) {
	t.Helper()
	if !slices.Equal(node.ClassBindingOrder, want) {
		t.Errorf("%s binding order = %v, want %v", node.Name, node.ClassBindingOrder, want)
	}
}

func blueprintContractAssertClassBinding(
	t *testing.T,
	node *blueprint.Node,
	property string,
	want map[string]string,
) {
	t.Helper()
	if got := node.ClassBindings[property]; !maps.Equal(got, want) {
		t.Errorf("%s.%s class binding = %v, want %v", node.Name, property, got, want)
	}
}

func blueprintContractAssertCondition(
	t *testing.T,
	node *blueprint.Node,
	property string,
	want map[string]any,
) {
	t.Helper()
	got, ok := node.Props[property].(map[string]any)
	if !ok || !reflect.DeepEqual(got, want) {
		t.Errorf("%s.%s = %#v, want %#v", node.Name, property, node.Props[property], want)
	}
}

func blueprintContractAssertDerivedProgress(t *testing.T, node *blueprint.Node, formatted bool) {
	t.Helper()
	for property, want := range map[string]string{
		"progress_value_prop": "value",
		"progress_max_prop":   "max",
	} {
		if got := node.Props[property]; got != want {
			t.Errorf("Progress/%s %s = %#v, want %q", node.Name, property, got, want)
		}
	}
	format, hasFormat := node.Props["progress_format"]
	if formatted {
		if !hasFormat || format != "percent" {
			t.Errorf("Progress/%s progress_format = %#v, want percent", node.Name, format)
		}
	} else if hasFormat {
		t.Errorf("Progress/%s has unexpected progress_format %#v", node.Name, format)
	}
}

func blueprintContractStringValues(value any) []string {
	switch values := value.(type) {
	case []string:
		return slices.Clone(values)
	case []any:
		out := make([]string, 0, len(values))
		for _, value := range values {
			text, ok := value.(string)
			if !ok {
				return nil
			}
			out = append(out, text)
		}
		return out
	default:
		return nil
	}
}
