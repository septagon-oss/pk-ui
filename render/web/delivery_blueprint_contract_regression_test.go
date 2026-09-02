// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

import (
	"fmt"
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
	if got, want := blueprintContractDirectChildNames(root), []string{
		"StatusDot", "LeadingIcon", "Label", "Count", "TrailingIcon", "Remove",
	}; !slices.Equal(got, want) {
		t.Fatalf("Badge child order = %v, want runtime order %v", got, want)
	}

	dot := blueprintContractOnlyNodeNamed(t, root, "StatusDot")
	blueprintContractAssertOrder(t, dot, "tone")
	if got, want := dot.ClassBindings["tone"], compileClassMap(clBadgeDotTone); !maps.Equal(got, want) {
		t.Errorf("Badge/StatusDot tone classes = %v, want %v", got, want)
	}
	blueprintContractAssertCondition(t, dot, "visible_when", map[string]any{"dot": true})

	count := blueprintContractOnlyNodeNamed(t, root, "Count")
	if count.Kind != blueprint.NodeFrame {
		t.Errorf("Badge/Count = %#v, want an auto-layout spacing frame", count)
	}
	blueprintContractAssertCondition(t, count, "visible_when", map[string]any{"count": "*"})
	if fields := strings.Fields(count.Classes); !slices.Contains(fields, "pl-1") || slices.Contains(fields, "ml-1") || !slices.Contains(fields, "inline-flex") {
		t.Errorf("Badge/Count classes = %q, want parseable inline-flex/pl-1 and no ml-1", count.Classes)
	}
	countValue := blueprintContractOnlyNodeNamed(t, count, "CountValue")
	if countValue.Kind != blueprint.NodeText || countValue.Text != "" {
		t.Errorf("Badge/CountValue = %#v, want empty optional text", countValue)
	}
	if got := blueprintContractStringValues(countValue.Props["text_props"]); !slices.Equal(got, []string{"count"}) {
		t.Errorf("Badge/Count text props = %v, want [count]", got)
	}
	for _, property := range []string{"clear_when_unbound", "hide_when_empty"} {
		if enabled, _ := countValue.Props[property].(bool); !enabled {
			t.Errorf("Badge/CountValue %s = %#v, want true", property, countValue.Props[property])
		}
	}
	for property, want := range map[string]string{
		"integer_value_prop":    "count",
		"integer_max":           "99",
		"integer_overflow_text": "99+",
	} {
		if got := countValue.Props[property]; !reflect.DeepEqual(got, want) {
			if property != "integer_max" || fmt.Sprint(got) != want {
				t.Errorf("Badge/Count %s = %#v, want %q", property, got, want)
			}
		}
	}
	remove := blueprintContractOnlyNodeNamed(t, root, "Remove")
	blueprintContractAssertCondition(t, remove, "visible_when", map[string]any{"removable": true})
	if got := remove.Props["semantic_role"]; got != "button" {
		t.Errorf("Badge/Remove semantic role = %#v, want button", got)
	}
	if got := remove.Props["accessible_name_prop"]; got != "removeLabel" {
		t.Errorf("Badge/Remove accessible name binding = %#v, want removeLabel", got)
	}
	removeIcon := blueprintContractOnlyNodeNamed(t, remove, "RemoveIcon")
	if removeIcon.Kind != blueprint.NodeInstance || removeIcon.Text != "Icon" || removeIcon.Props["instance_example"] != "X Mark / Fg Primary / 12" {
		t.Errorf("Badge/RemoveIcon = %#v, want the runtime x-mark xs instance", removeIcon)
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
	blueprintContractAssertOrder(t, control, "size", "tone", "error", "invalid", "readOnly", "disabled")
	blueprintContractAssertClassBinding(t, control, "error", map[string]string{
		"": "", "*": clInputError.Compile(),
	})
	blueprintContractAssertClassBinding(t, control, "invalid", map[string]string{
		"false": "", "true": clInputError.Compile(),
	})
	blueprintContractAssertClassBinding(t, control, "readOnly", map[string]string{
		"false": "", "true": clInputReadOnly.Compile(),
	})
	blueprintContractAssertClassBinding(t, control, "disabled", map[string]string{
		"false": "", "true": clInputDisabled.Compile(),
	})

	value := blueprintContractOnlyNodeNamed(t, root, "Value")
	if value.Kind != blueprint.NodeFrame {
		t.Errorf("Input/Value = %#v, want a growing auto-layout frame", value)
	}
	for _, class := range []string{"min-w-0", "flex-1"} {
		if !slices.Contains(strings.Fields(value.Classes), class) {
			t.Errorf("Input/Value classes = %q, want %q", value.Classes, class)
		}
	}
	valueText := blueprintContractOnlyNodeNamed(t, value, "ValueText")
	blueprintContractAssertOrder(t, valueText, "value")
	blueprintContractAssertClassBinding(t, valueText, "value", map[string]string{
		"":  tw.New().TextColor(tw.FgPlaceholder).Compile(),
		"*": tw.New().TextColor(tw.FgPrimary).Compile(),
	})
	if got, want := valueText.Props["mask_when"], (map[string]any{"type": "password"}); !reflect.DeepEqual(got, want) {
		t.Errorf("Input/Value mask_when = %#v, want %#v", got, want)
	}
	labelRow := blueprintContractOnlyNodeNamed(t, root, "LabelRow")
	blueprintContractAssertCondition(t, labelRow, "visible_when", map[string]any{"label": "*"})
	requiredMarker := blueprintContractOnlyNodeNamed(t, root, "RequiredMarker")
	blueprintContractAssertCondition(t, requiredMarker, "visible_when", map[string]any{"required": true})
	if !slices.Contains(strings.Fields(requiredMarker.Classes), "text-fg-danger") {
		t.Errorf("Input/RequiredMarker classes = %q, want danger token", requiredMarker.Classes)
	}

	help := blueprintContractOnlyNodeNamed(t, root, "Help")
	blueprintContractAssertOrder(t, help, "error")
	blueprintContractAssertClassBinding(t, help, "error", map[string]string{
		"": "", "*": clFieldErr.Compile(),
	})

	required := deliveryExamplePropsForTest(t, "Input", "required")
	if required["required"] != true || required["label"] == "" {
		t.Errorf("Input required specimen = %#v, want labelled required state", required)
	}
	disabled := deliveryExamplePropsForTest(t, "Input", "disabled")
	if disabled["disabled"] != true || disabled["value"] == "" {
		t.Errorf("Input disabled specimen = %#v, want populated disabled state", disabled)
	}
	adorned := blueprintContractExample(t, definition, "with-adornments")
	if got, want := blueprintContractExampleSlotNames(adorned), []string{"iconStart", "iconEnd"}; !slices.Equal(got, want) {
		t.Errorf("Input adorned specimen slots = %v, want %v", got, want)
	}
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
	blueprintContractAssertCondition(t, fill, "progress_width_unless", map[string]any{"indeterminate": true})

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
	source := definition.Design.Assets[0].Source
	for _, fragment := range []string{`stroke-width="2"`, `stroke-linecap="round"`, `vector-effect="non-scaling-stroke"`} {
		if !strings.Contains(source, fragment) {
			t.Errorf("Spinner arc asset = %q, want %s", source, fragment)
		}
	}
	if strings.Contains(source, `fill-rule="evenodd"`) {
		t.Errorf("Spinner arc still uses a filled annulus instead of a two-pixel stroke: %s", source)
	}
}

func TestTagBlueprintMatchesRemovalPredicateAndCoversSemanticTones(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Tag")
	remove := blueprintContractOnlyNodeNamed(t, definition.Design.Root, "Remove")
	blueprintContractAssertCondition(t, remove, "visible_when", map[string]any{
		"removable": true, "onRemoveURL": "*",
	})
	if got := remove.Props["semantic_role"]; got != "button" {
		t.Errorf("Tag/Remove semantic role = %#v, want button", got)
	}
	removeGlyph := blueprintContractOnlyNodeNamed(t, remove, "RemoveGlyph")
	if removeGlyph.Kind != blueprint.NodeText || removeGlyph.Text != "×" {
		t.Errorf("Tag/RemoveGlyph = %#v, want runtime glyph", removeGlyph)
	}

	toneExamples := make(map[string]bool, len(canonicalTones))
	for _, example := range definition.Examples {
		if tone, _ := example.Props["tone"].(string); tone != "" {
			toneExamples[tone] = true
		}
	}
	for _, tone := range canonicalTones {
		if !toneExamples[tone] {
			t.Errorf("Tag examples do not cover tone %q: %v", tone, toneExamples)
		}
	}
	removable := deliveryExamplePropsForTest(t, "Tag", "removable")
	if removable["removable"] != true || removable["onRemoveURL"] == "" {
		t.Errorf("Tag removable specimen = %#v, want state plus endpoint", removable)
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
	if got, want := blueprintContractDirectChildNames(panel), []string{
		"Header", "HeaderSeparator", "Body", "FooterSeparator", "Footer",
	}; !slices.Equal(got, want) {
		t.Fatalf("Modal/Panel child order = %v, want %v", got, want)
	}
	for _, name := range []string{"HeaderSeparator", "FooterSeparator"} {
		separator := blueprintContractOnlyNodeNamed(t, panel, name)
		for _, class := range []string{"w-full", "h-px", "bg-border-primary", "flex-shrink-0"} {
			if !slices.Contains(strings.Fields(separator.Classes), class) {
				t.Errorf("Modal/%s classes = %q, want %q", name, separator.Classes, class)
			}
		}
		if got := separator.Props["semantic_role"]; got != "separator" {
			t.Errorf("Modal/%s semantic role = %#v, want separator", name, got)
		}
	}
	for _, name := range []string{"Header", "Footer"} {
		node := blueprintContractOnlyNodeNamed(t, panel, name)
		for _, class := range strings.Fields(node.Classes) {
			if class == "border-t" || class == "border-b" {
				t.Errorf("Modal/%s retains target-specific side border %q", name, class)
			}
		}
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

func blueprintContractExample(t *testing.T, definition DeliveryDefinition, name string) DeliveryExample {
	t.Helper()
	wantID := deliveryExampleID(string(definition.Identity.ID), name)
	for _, example := range definition.Examples {
		if example.ID == wantID {
			return example
		}
	}
	t.Fatalf("component %s has no example %s", definition.Identity.ID, name)
	return DeliveryExample{}
}

func blueprintContractExampleSlotNames(example DeliveryExample) []string {
	names := make([]string, 0, len(example.Slots))
	for _, slot := range example.Slots {
		names = append(names, slot.Name)
	}
	return names
}

func blueprintContractDirectChildNames(node *blueprint.Node) []string {
	if node == nil {
		return nil
	}
	names := make([]string, 0, len(node.Children))
	for _, child := range node.Children {
		names = append(names, child.Name)
	}
	return names
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
