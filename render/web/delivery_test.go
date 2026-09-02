// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_test.go proves the OSS catalog is closed over real renderers,
// stable identities, native semantic blueprints, and executable fixtures.
// Browser pixels remain covered by Storybook visual baselines; this boundary
// deliberately rejects raster snapshots as design structure.

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-design/pkg/blueprint"
	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
)

func TestOSSDeliveryCatalogIsCompleteNativeAndExecutable(t *testing.T) {
	t.Parallel()

	catalog := OSSDeliveryCatalog()
	if got, want := len(catalog), 55; got != want {
		t.Fatalf("catalog has %d entries, want %d", got, want)
	}

	identities := make(map[string]struct{}, len(catalog))
	stableIDs := make(map[string]struct{}, len(catalog))
	for _, definition := range catalog {
		if err := definition.Validate(); err != nil {
			t.Fatalf("%s: invalid definition: %v", definition.Identity.ID, err)
		}
		identity := string(definition.Identity.ID)
		if _, duplicate := identities[identity]; duplicate {
			t.Fatalf("duplicate wire identity %q", identity)
		}
		identities[identity] = struct{}{}
		if _, duplicate := stableIDs[definition.Contract.ID]; duplicate {
			t.Fatalf("duplicate stable design id %q", definition.Contract.ID)
		}
		stableIDs[definition.Contract.ID] = struct{}{}
		assertNativeDeliveryNode(t, identity, definition.Design.Root)

		for _, example := range definition.Examples {
			node, err := RenderDeliveryExample(identity, example.ID, nil)
			if err != nil {
				t.Fatalf("%s/%s: render: %v", identity, example.Name, err)
			}
			var rendered bytes.Buffer
			if err := node.Render(&rendered); err != nil {
				t.Fatalf("%s/%s: write HTML: %v", identity, example.Name, err)
			}
			if strings.TrimSpace(rendered.String()) == "" {
				t.Fatalf("%s/%s: rendered empty HTML", identity, example.Name)
			}
		}
	}

	if _, ok := identities["DataManagementPage"]; !ok {
		t.Fatal("catalog has no complete OSS solution page")
	}
	for componentType, expectedTag := range map[string]string{
		"Icon": "icon",
		"Text": "typography",
	} {
		var found bool
		for _, definition := range catalog {
			if string(definition.Identity.ID) != componentType {
				continue
			}
			found = slices.Contains(definition.Contract.Tags, expectedTag)
			break
		}
		if !found {
			t.Errorf("%s does not publish semantic tag %q", componentType, expectedTag)
		}
	}
}

func TestAdaptiveDetailAndOfflineContractsArePublished(t *testing.T) {
	t.Parallel()

	definitions := make(map[string]DeliveryDefinition)
	for _, definition := range OSSDeliveryCatalog() {
		definitions[string(definition.Identity.ID)] = definition
	}

	detail, ok := definitions["DetailList"]
	if !ok {
		t.Fatal("catalog has no DetailList delivery contract")
	}
	detailProps := make(map[string]bool)
	for _, prop := range detail.Contract.Props {
		detailProps[prop.Name] = true
	}
	for _, name := range []string{"title", "description", "semanticRole", "items"} {
		if !detailProps[name] {
			t.Errorf("DetailList contract is missing %q", name)
		}
	}

	window, ok := definitions["WindowedCollection"]
	if !ok {
		t.Fatal("catalog has no WindowedCollection delivery contract")
	}
	windowProps := make(map[string]bool)
	var states []string
	for _, prop := range window.Contract.Props {
		windowProps[prop.Name] = true
		if prop.Name == "state" {
			states = prop.EnumValues
		}
	}
	if !slices.Contains(states, "offline") {
		t.Errorf("WindowedCollection state enum = %v, want offline", states)
	}
	for _, name := range []string{"offlineTitle", "offlineDescription"} {
		if !windowProps[name] {
			t.Errorf("WindowedCollection contract is missing %q", name)
		}
	}
}

func TestDataGridDeliveryExampleUsesRealNamedSlotComposition(t *testing.T) {
	t.Parallel()

	node, err := RenderDeliveryExample(
		"DataGrid",
		deliveryExampleID("DataGrid", "workspace-list"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	var rendered bytes.Buffer
	if err := node.Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`role="search"`,
		`aria-label="Search workspaces"`,
		`>New workspace</a>`,
		`<table`,
		`Alpha workspace`,
		`aria-label="Pagination"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("composed DataGrid example is missing %q:\n%s", fragment, html)
		}
	}
}

func TestDashboardWidgetDeliveryUsesCanonicalRendererAndSlots(t *testing.T) {
	t.Parallel()

	node, err := RenderDeliveryExample(
		"DashboardWidget",
		deliveryExampleID("DashboardWidget", "activity-list"),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	var rendered bytes.Buffer
	if err := node.Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-component="dashboard-widget"`, `data-widget-type="list"`,
		`Three workspace events need review.`, `>View activity</a>`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("composed DashboardWidget example is missing %q:\n%s", fragment, html)
		}
	}
	if strings.Contains(html, `hx-get=`) {
		t.Fatalf("static activity fixture must not auto-request an application route:\n%s", html)
	}
}

func TestImageDeliveryExamplesAreBrowserSelfContained(t *testing.T) {
	t.Parallel()

	tests := []struct {
		component string
		example   string
		property  string
	}{
		{component: "Avatar", example: "with-image", property: "src"},
		{component: "FileUpload", example: "remote-upload", property: "value"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.example, func(t *testing.T) {
			t.Parallel()
			props := deliveryExamplePropsForTest(t, tt.component, tt.example)
			value, _ := props[tt.property].(string)
			if !strings.HasPrefix(value, "data:image/svg+xml,") {
				t.Fatalf("%s fixture %s must use an inline image, got %q", tt.component, tt.property, value)
			}
		})
	}
}

func deliveryExamplePropsForTest(t *testing.T, componentType, exampleName string) map[string]any {
	t.Helper()
	for _, definition := range OSSDeliveryCatalog() {
		if string(definition.Identity.ID) != componentType {
			continue
		}
		wantID := deliveryExampleID(componentType, exampleName)
		for _, example := range definition.Examples {
			if example.ID == wantID {
				return example.Props
			}
		}
		t.Fatalf("component %s has no example %s", componentType, exampleName)
	}
	t.Fatalf("delivery catalog has no component %s", componentType)
	return nil
}

func TestIconDeliveryBlueprintIsBoundEditableVector(t *testing.T) {
	t.Parallel()

	var icon *DeliveryDefinition
	catalog := OSSDeliveryCatalog()
	for index := range catalog {
		if catalog[index].Identity.ID == "Icon" {
			icon = &catalog[index]
			break
		}
	}
	if icon == nil {
		t.Fatal("catalog has no Icon definition")
	}
	root := icon.Design.Root
	if root == nil || root.Kind != blueprint.NodeSVG {
		t.Fatalf("Icon design root = %#v, want native SVG", root)
	}
	if root.AssetRef != "bound:icon-" {
		t.Errorf("Icon asset ref = %q, want bound:icon-", root.AssetRef)
	}
	if root.Text != "" || len(root.Children) != 0 {
		t.Errorf("Icon design contains placeholder content: %#v", root)
	}
	if bound, _ := root.Props["bind_name"].(bool); !bound {
		t.Errorf("Icon design does not bind SVG asset to name: %#v", root.Props)
	}
}

func TestIconDeliveryPublishesExtraSmallBrandCheckSpecimen(t *testing.T) {
	t.Parallel()

	const example = "Check / Fg Brand / 12"
	props := deliveryExamplePropsForTest(t, "Icon", example)
	for name, want := range map[string]string{
		"name":   "check",
		"size":   "xs",
		"tone":   "brand",
		"weight": "outline",
	} {
		if got, _ := props[name].(string); got != want {
			t.Errorf("Icon %s = %q, want %q", name, got, want)
		}
	}

	node, err := RenderDeliveryExample("Icon", deliveryExampleID("Icon", example), nil)
	if err != nil {
		t.Fatal(err)
	}
	var rendered strings.Builder
	if err := node.Render(&rendered); err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{
		`data-pk-icon="check"`,
		`data-pk-icon-weight="outline"`,
		"h-3",
		"w-3",
		"text-fg-brand",
	} {
		if !strings.Contains(rendered.String(), fragment) {
			t.Errorf("extra-small brand Check is missing %q: %s", fragment, rendered.String())
		}
	}
}

func TestDeliveryRuntimeNeverInjectsExampleContentWithoutSlots(t *testing.T) {
	t.Parallel()

	catalog := OSSDeliveryCatalog()
	for _, componentType := range []string{"Card", "Stack", "Grid", "Container"} {
		var definition *DeliveryDefinition
		for index := range catalog {
			if string(catalog[index].Identity.ID) == componentType {
				definition = &catalog[index]
				break
			}
		}
		if definition == nil {
			t.Fatalf("missing %s definition", componentType)
		}
		node, err := definition.Render(nil, nil)
		if err != nil {
			t.Fatalf("%s: %v", componentType, err)
		}
		var rendered bytes.Buffer
		if err := node.Render(&rendered); err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{
			"Editable child slot",
			"Primary content",
			"Supporting content",
			"Review the current operational details.",
		} {
			if strings.Contains(rendered.String(), forbidden) {
				t.Errorf(
					"%s injected example content %q without an authored slot",
					componentType,
					forbidden,
				)
			}
		}
	}
}

func TestInputDeliveryRejectsUnsupportedNativeType(t *testing.T) {
	t.Parallel()

	var input *DeliveryDefinition
	catalog := OSSDeliveryCatalog()
	for index := range catalog {
		definition := catalog[index]
		if definition.Identity.ID == "Input" {
			input = &definition
			break
		}
	}
	if input == nil {
		t.Fatal("catalog has no Input definition")
	}
	if _, err := input.Render(map[string]any{
		"name": "attachment",
		"type": "file",
	}, nil); err == nil || !strings.Contains(err.Error(), "unsupported type") {
		t.Fatalf("Input delivery error = %v, want unsupported type", err)
	}
}

func TestOSSDeliveryCatalogPublishesBooleanRestingStates(t *testing.T) {
	t.Parallel()

	for _, definition := range OSSDeliveryCatalog() {
		for _, prop := range definition.Contract.Props {
			if prop.Type != designcomponent.PropBoolean || prop.Required {
				continue
			}
			if strings.TrimSpace(prop.Default) == "" {
				t.Errorf(
					"%s.%s has no explicit boolean resting state",
					definition.Identity.ID,
					prop.Name,
				)
			}
		}
	}
}

func TestOSSDeliveryCatalogComposesDownTheAtomicGraph(t *testing.T) {
	t.Parallel()

	catalog := OSSDeliveryCatalog()
	tiers := make(map[string]uicomponent.Tier, len(catalog))
	for _, definition := range catalog {
		tiers[string(definition.Identity.ID)] = definition.Identity.Tier
	}
	rank := map[uicomponent.Tier]int{
		uicomponent.TierAtom:     1,
		uicomponent.TierMolecule: 2,
		uicomponent.TierOrganism: 3,
		uicomponent.TierTemplate: 4,
		uicomponent.TierPage:     5,
	}
	for _, definition := range catalog {
		source := string(definition.Identity.ID)
		for _, slot := range definition.Contract.Slots {
			for _, allowedType := range slot.AllowedTypes {
				targetTier, exists := tiers[allowedType]
				if !exists {
					t.Errorf("%s slot %s allows unknown type %q", source, slot.Name, allowedType)
					continue
				}
				if rank[targetTier] > rank[definition.Identity.Tier] {
					t.Errorf(
						"%s (%s) slot %s allows higher-tier %s (%s)",
						source,
						definition.Identity.Tier,
						slot.Name,
						allowedType,
						targetTier,
					)
				}
			}
		}
		walkDeliveryNodes(definition.Design.Root, func(node *blueprint.Node) {
			if node.Kind != blueprint.NodeInstance || strings.HasPrefix(node.Text, "icon-") {
				return
			}
			targetTier, exists := tiers[node.Text]
			if !exists {
				t.Errorf("%s references unknown native instance %q", source, node.Text)
				return
			}
			if rank[targetTier] > rank[definition.Identity.Tier] {
				t.Errorf(
					"%s (%s) references higher-tier %s (%s)",
					source,
					definition.Identity.Tier,
					node.Text,
					targetTier,
				)
			}
		})
	}
}

func TestDataManagementSolutionHasOneResponsivePrimaryAction(t *testing.T) {
	t.Parallel()

	var rendered bytes.Buffer
	if err := renderDataManagementPage(DataManagementPageProps{}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()

	if got := strings.Count(html, "New workspace"); got != 1 {
		t.Fatalf("primary action rendered %d times, want exactly once:\n%s", got, html)
	}
	for _, requiredClass := range []string{
		"flex-col",
		"items-stretch",
		"sm:flex-row",
		"sm:items-center",
		"w-full",
		"sm:ml-auto",
	} {
		if !strings.Contains(html, requiredClass) {
			t.Errorf("responsive data toolbar is missing %q:\n%s", requiredClass, html)
		}
	}
}

func assertNativeDeliveryNode(t *testing.T, component string, node *blueprint.Node) {
	t.Helper()
	if node == nil {
		t.Fatalf("%s: nil native root", component)
	}
	walkDeliveryNodes(node, func(current *blueprint.Node) {
		if current.Name == "" {
			t.Errorf("%s: native %s node has no stable name", component, current.Kind)
		}
		for key, value := range current.Props {
			if strings.Contains(strings.ToLower(key), "screenshot") ||
				strings.Contains(strings.ToLower(key), "snapshot") {
				t.Errorf("%s/%s: forbidden raster-state property %q", component, current.Name, key)
			}
			if text, ok := value.(string); ok &&
				strings.HasPrefix(strings.ToLower(strings.TrimSpace(text)), "data:image/") {
				t.Errorf("%s/%s: embedded raster payload", component, current.Name)
			}
		}
	})
}

func walkDeliveryNodes(root *blueprint.Node, visit func(*blueprint.Node)) {
	if root == nil {
		return
	}
	visit(root)
	for index := range root.Children {
		walkDeliveryNodes(&root.Children[index], visit)
	}
}
