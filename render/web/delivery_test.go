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
	"encoding/base64"
	"image/png"
	"reflect"
	"slices"
	"strconv"
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
		{component: "Avatar", example: "with-image-circle-md", property: "src"},
		{component: "FileUpload", example: "remote-upload", property: "value"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.example, func(t *testing.T) {
			t.Parallel()
			props := deliveryExamplePropsForTest(t, tt.component, tt.example)
			value, _ := props[tt.property].(string)
			if !strings.HasPrefix(value, "data:image/png;base64,") {
				t.Fatalf("%s fixture %s must use an inline image, got %q", tt.component, tt.property, value)
			}
		})
	}

	uploadProps := deliveryExamplePropsForTest(t, "FileUpload", "remote-upload")
	if got, _ := uploadProps["currentName"].(string); got != "logo.png" {
		t.Fatalf("FileUpload/remote-upload currentName = %q, want logo.png", got)
	}
}

func TestDeliveryFixtureImageIsValidAvatarPNG(t *testing.T) {
	t.Parallel()

	const prefix = "data:image/png;base64,"
	payload := strings.TrimPrefix(deliveryFixtureImageDataURL, prefix)
	if payload == deliveryFixtureImageDataURL {
		t.Fatalf("delivery fixture must use a PNG data URI, got %q", deliveryFixtureImageDataURL)
	}
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("decode delivery fixture: %v", err)
	}
	config, err := png.DecodeConfig(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("decode delivery fixture PNG: %v", err)
	}
	if config.Width != 64 || config.Height != 64 {
		t.Fatalf("delivery fixture PNG dimensions = %dx%d, want 64x64", config.Width, config.Height)
	}
}

func TestAvatarDeliveryPublishesImageCircleMediumSpecimen(t *testing.T) {
	t.Parallel()

	const example = "with-image-circle-md"
	props := deliveryExamplePropsForTest(t, "Avatar", example)
	want := map[string]any{
		"src":   deliveryFixtureImageDataURL,
		"alt":   "Assistant avatar",
		"shape": "circle",
		"size":  "md",
	}
	if !reflect.DeepEqual(props, want) {
		t.Fatalf("Avatar/%s props = %#v, want %#v", example, props, want)
	}

	node, err := RenderDeliveryExample("Avatar", deliveryExampleID("Avatar", example), nil)
	if err != nil {
		t.Fatal(err)
	}
	var rendered strings.Builder
	if err := node.Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-component="avatar"`,
		`data-avatar-size="md"`,
		`data-avatar-shape="circle"`,
		`h-10`,
		`w-10`,
		`rounded-full`,
		`<img`,
		`alt="Assistant avatar"`,
		`src="data:image/png;base64,`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Avatar/%s is missing %q: %s", example, fragment, html)
		}
	}
	if strings.Contains(html, ">AL<") {
		t.Errorf("Avatar/%s rendered the design fallback instead of its image: %s", example, html)
	}
}

func TestAvatarImageCircleMediumSpecimenSelectsExactDesignGeometry(t *testing.T) {
	t.Parallel()

	props := deliveryExamplePropsForTest(t, "Avatar", "with-image-circle-md")
	definition := blueprintContractDefinition(t, "Avatar")
	root := definition.Design.Root

	size, _ := props["size"].(string)
	sizeClasses := strings.Fields(root.ClassBindings["size"][size])
	for _, class := range []string{"h-10", "w-10"} {
		if !slices.Contains(sizeClasses, class) {
			t.Errorf("Avatar size=%q design classes = %v, want %q", size, sizeClasses, class)
		}
	}

	shape, _ := props["shape"].(string)
	shapeClasses := strings.Fields(root.ClassBindings["shape"][shape])
	if !slices.Contains(shapeClasses, "rounded-full") {
		t.Errorf("Avatar shape=%q design classes = %v, want rounded-full", shape, shapeClasses)
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

func TestIconDeliveryPublishesEditorAndCompletionSpecimens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		example string
		name    string
		size    string
		tone    string
		classes []string
	}{
		{example: "Album Stack / Fg Primary / 16", name: "album-stack", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Text / Fg Primary / 16", name: "text", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Photo / Fg Primary / 16", name: "photo", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Graphics / Fg Primary / 16", name: "graphics", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Shape / Fg Primary / 16", name: "shape", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Design / Fg Primary / 16", name: "design", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Uploads / Fg Primary / 16", name: "uploads", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "User / Fg Primary / 16", name: "user", size: "sm", tone: "neutral", classes: []string{"h-4", "w-4"}},
		{example: "Check / Fg Success / 24", name: "check", size: "lg", tone: "success", classes: []string{"h-6", "w-6", "text-fg-success"}},
		{example: "Check Circle / Fg Success / 48", name: "check-circle", size: "2xl", tone: "success", classes: []string{"h-12", "w-12", "text-fg-success"}},
	}

	for _, test := range tests {
		t.Run(test.example, func(t *testing.T) {
			props := deliveryExamplePropsForTest(t, "Icon", test.example)
			for prop, want := range map[string]string{
				"name": test.name, "size": test.size, "tone": test.tone, "weight": "outline",
			} {
				if got, _ := props[prop].(string); got != want {
					t.Errorf("Icon %s = %q, want %q", prop, got, want)
				}
			}

			node, err := RenderDeliveryExample("Icon", deliveryExampleID("Icon", test.example), nil)
			if err != nil {
				t.Fatal(err)
			}
			var rendered strings.Builder
			if err := node.Render(&rendered); err != nil {
				t.Fatal(err)
			}
			html := rendered.String()
			for _, fragment := range append([]string{`data-pk-icon="` + test.name + `"`}, test.classes...) {
				if !strings.Contains(html, fragment) {
					t.Errorf("%s is missing %q: %s", test.example, fragment, html)
				}
			}
			if strings.Contains(html, `data-pk-icon-fallback="true"`) {
				t.Errorf("%s rendered a fallback glyph: %s", test.example, html)
			}
		})
	}
}

func TestDeliveryPublishesOnboardingActionAndProgressSpecimens(t *testing.T) {
	t.Parallel()

	buttonProps := deliveryExamplePropsForTest(t, "Button", "primary-success-sm")
	wantButtonProps := map[string]any{
		"label": "Complete", "variant": "primary", "tone": "success", "size": "sm",
	}
	if !reflect.DeepEqual(buttonProps, wantButtonProps) {
		t.Fatalf("Button/primary-success-sm props = %#v, want %#v", buttonProps, wantButtonProps)
	}
	button, err := RenderDeliveryExample("Button", deliveryExampleID("Button", "primary-success-sm"), nil)
	if err != nil {
		t.Fatal(err)
	}
	var renderedButton strings.Builder
	if err := button.Render(&renderedButton); err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{
		`data-variant="primary"`, `data-tone="success"`, "Complete",
		"bg-surface-success", "text-sm", "px-3", "py-1.5",
	} {
		if !strings.Contains(renderedButton.String(), fragment) {
			t.Errorf("Button/primary-success-sm is missing %q: %s", fragment, renderedButton.String())
		}
	}

	for _, test := range []struct {
		example string
		value   int
	}{
		{example: "compact-brand-sm", value: 33},
		{example: "compact-brand-sm-67", value: 67},
		{example: "compact-brand-sm-100", value: 100},
	} {
		t.Run("Progress/"+test.example, func(t *testing.T) {
			progressProps := deliveryExamplePropsForTest(t, "Progress", test.example)
			for property, want := range map[string]string{
				"ariaLabel": "Onboarding progress", "tone": "brand", "size": "sm",
			} {
				if got, _ := progressProps[property].(string); got != want {
					t.Errorf("%s = %q, want %q", property, got, want)
				}
			}
			if got, _ := progressProps["showText"].(bool); got {
				t.Errorf("showText = true, want false")
			}
			for property, want := range map[string]int{"value": test.value, "max": 100} {
				got := progressProps[property]
				if got != want && got != float64(want) {
					t.Errorf("%s = %#v, want %d", property, got, want)
				}
			}

			progress, err := RenderDeliveryExample("Progress", deliveryExampleID("Progress", test.example), nil)
			if err != nil {
				t.Fatal(err)
			}
			var renderedProgress strings.Builder
			if err := progress.Render(&renderedProgress); err != nil {
				t.Fatal(err)
			}
			progressHTML := renderedProgress.String()
			value := strconv.Itoa(test.value)
			for _, fragment := range []string{
				`aria-label="Onboarding progress"`, `aria-valuenow="` + value + `"`,
				`aria-valuemax="100"`, `data-progress-percent="` + value + `"`,
				`data-size="sm"`, `data-tone="brand"`, "h-1.5", "bg-surface-brand",
			} {
				if !strings.Contains(progressHTML, fragment) {
					t.Errorf("rendered output is missing %q: %s", fragment, progressHTML)
				}
			}
			if strings.Contains(progressHTML, `data-progress-value="true"`) || strings.Contains(progressHTML, value+"%") {
				t.Errorf("unexpectedly renders percentage text: %s", progressHTML)
			}
		})
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
