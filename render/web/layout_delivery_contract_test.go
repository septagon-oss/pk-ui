package web

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestLayoutBlueprintsMaterializeEveryPublishedAxisValueIndependently(t *testing.T) {
	t.Parallel()

	tests := []struct {
		component string
		order     []string
		baseline  string
		values    map[string][]string
		count     int
	}{
		{
			component: "Flex",
			order:     []string{"direction", "wrap", "gap", "align", "justify"},
			baseline:  "row",
			values: map[string][]string{
				"direction": {"column", "row"},
				"wrap":      {"false", "true"},
				"gap":       {"0", "1", "2", "3", "4", "5", "6", "8"},
				"align":     {"center", "end", "start", "stretch"},
				"justify":   {"between", "center", "end", "start"},
			},
			count: 16,
		},
		{
			component: "Grid",
			order:     []string{"columns", "gap"},
			baseline:  "one-column",
			values: map[string][]string{
				"columns": {"1", "12", "2", "3", "4", "6"},
				"gap":     {"0", "1", "2", "3", "4", "5", "6", "8"},
			},
			count: 13,
		},
		{
			component: "Container",
			order:     []string{"maxWidth", "padding"},
			baseline:  "default",
			values: map[string][]string{
				"maxWidth": {"2xl", "4xl", "7xl", "full", "lg", "md", "sm", "xl"},
				"padding":  {"0", "1", "2", "3", "4", "5", "6", "8"},
			},
			count: 16,
		},
	}

	for _, test := range tests {
		t.Run(test.component, func(t *testing.T) {
			t.Parallel()

			definition := blueprintContractDefinition(t, test.component)
			root := definition.Design.Root
			blueprintContractAssertOrder(t, root, test.order...)
			if got := len(definition.Examples); got != test.count {
				t.Fatalf("%s examples = %d, want %d non-cartesian axis specimens", test.component, got, test.count)
			}
			baseline := blueprintContractExample(t, definition, test.baseline)

			for _, property := range test.order {
				bindings := root.ClassBindings[property]
				values := test.values[property]
				if got := slices.Sorted(maps.Keys(bindings)); !slices.Equal(got, values) {
					t.Fatalf("%s.%s binding values = %v, want %v", test.component, property, got, values)
				}

				for _, value := range values {
					example, found := layoutIndependentAxisExample(
						definition,
						baseline,
						test.order,
						property,
						value,
					)
					if !found {
						t.Errorf("%s.%s=%s has no example that changes only that axis from %s", test.component, property, value, test.baseline)
						continue
					}

					rootClasses := renderDeliveryExampleRootClasses(t, test.component, example.ID)
					gotGeometry := layoutAxisClasses(property, rootClasses)
					wantGeometry := strings.Fields(bindings[value])
					if !slices.Equal(gotGeometry, wantGeometry) {
						t.Errorf(
							"%s/%s native %s=%s geometry classes = %v, want exactly %v (root %v)",
							test.component,
							example.Name,
							property,
							value,
							gotGeometry,
							wantGeometry,
							rootClasses,
						)
					}
				}
			}
		})
	}
}

func layoutIndependentAxisExample(
	definition DeliveryDefinition,
	baseline DeliveryExample,
	properties []string,
	targetProperty string,
	targetValue string,
) (DeliveryExample, bool) {
	for _, example := range definition.Examples {
		matches := true
		for _, property := range properties {
			want := fmt.Sprint(baseline.Props[property])
			if property == targetProperty {
				want = targetValue
			}
			value, exists := example.Props[property]
			if !exists || fmt.Sprint(value) != want {
				matches = false
				break
			}
		}
		if matches {
			return example, true
		}
	}
	return DeliveryExample{}, false
}

func layoutAxisClasses(property string, classes []string) []string {
	axes := map[string]func(string) bool{
		"direction": func(class string) bool {
			return class == "flex-row" || class == "flex-col"
		},
		"wrap": func(class string) bool {
			return class == "flex-wrap"
		},
		"gap": func(class string) bool {
			return strings.HasPrefix(class, "gap-")
		},
		"align": func(class string) bool {
			return strings.HasPrefix(class, "items-")
		},
		"justify": func(class string) bool {
			return strings.HasPrefix(class, "justify-")
		},
		"columns": func(class string) bool {
			return strings.HasPrefix(class, "grid-cols-")
		},
		"maxWidth": func(class string) bool {
			return strings.HasPrefix(class, "max-w-")
		},
		"padding": func(class string) bool {
			return strings.HasPrefix(class, "px-")
		},
	}

	matching := make([]string, 0, 1)
	for _, class := range classes {
		if axes[property](class) {
			matching = append(matching, class)
		}
	}
	return matching
}

func renderDeliveryExampleRootClasses(
	t *testing.T,
	component string,
	exampleID string,
) []string {
	t.Helper()

	node, err := RenderDeliveryExample(component, exampleID, nil)
	if err != nil {
		t.Fatalf("render %s/%s: %v", component, exampleID, err)
	}
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatalf("render %s/%s node: %v", component, exampleID, err)
	}
	openingTag, _, found := strings.Cut(output.String(), ">")
	if !found {
		t.Fatalf("render %s/%s has no opening tag: %s", component, exampleID, output.String())
	}
	_, classSuffix, found := strings.Cut(openingTag, `class="`)
	if !found {
		t.Fatalf("render %s/%s root has no class attribute: %s", component, exampleID, openingTag)
	}
	classes, _, found := strings.Cut(classSuffix, `"`)
	if !found {
		t.Fatalf("render %s/%s root has an unterminated class attribute: %s", component, exampleID, openingTag)
	}
	return strings.Fields(classes)
}

func TestLayoutBlueprintBindingsUseExecutableUtilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		component string
		property  string
		value     string
		class     string
	}{
		{component: "Flex", property: "direction", value: "column", class: "flex-col"},
		{component: "Flex", property: "wrap", value: "true", class: "flex-wrap"},
		{component: "Flex", property: "gap", value: "6", class: "gap-6"},
		{component: "Flex", property: "align", value: "stretch", class: "items-stretch"},
		{component: "Flex", property: "justify", value: "between", class: "justify-between"},
		{component: "Grid", property: "columns", value: "12", class: "grid-cols-12"},
		{component: "Grid", property: "gap", value: "3", class: "gap-3"},
		{component: "Container", property: "maxWidth", value: "sm", class: "max-w-sm"},
		{component: "Container", property: "padding", value: "2", class: "px-2"},
	}

	for _, test := range tests {
		t.Run(test.component+"/"+test.property, func(t *testing.T) {
			t.Parallel()

			root := blueprintContractDefinition(t, test.component).Design.Root
			classes := strings.Fields(root.ClassBindings[test.property][test.value])
			if !slices.Contains(classes, test.class) {
				t.Errorf("%s.%s[%s] classes = %v, want %q", test.component, test.property, test.value, classes, test.class)
			}
		})
	}
}

func TestTextBlueprintGivesEveryPublishedVariantAnExplicitFamily(t *testing.T) {
	t.Parallel()

	definition := blueprintContractDefinition(t, "Text")
	if classes := strings.Fields(definition.Design.Root.Classes); !slices.Contains(classes, "font-sans") {
		t.Fatalf("Text root classes = %v, want explicit font-sans", classes)
	}

	want := []string{"body", "muted", "muted-body", "emphasis", "clamped"}
	got := make([]string, 0, len(definition.Examples))
	for _, example := range definition.Examples {
		got = append(got, example.Name)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Text examples = %v, want %v", got, want)
	}
}
