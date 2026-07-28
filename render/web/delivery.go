// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery.go defines the OSS web design-delivery seam. It binds the canonical
// pk-ui wire identity, pk-design component contract, native visual blueprint,
// authored fixtures, and real gomponents renderer in one contribution. Private
// runtimes adapt this contract; they do not reconstruct it by matching names.

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"

	g "maragu.dev/gomponents"

	"github.com/septagon-oss/pk-design/pkg/blueprint"
	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
)

const (
	// DeliveryOwner is recorded on every OSS web design contribution.
	DeliveryOwner = "septagon-oss/pk-ui/render/web"
	// DeliveryVersion is the contract version for the complete OSS web catalog.
	DeliveryVersion = "1.0.0"
)

// DeliveryRenderFunc renders one component from provider-neutral properties.
// Storybook and other downstream preview runtimes call the same renderer used
// by applications.
type DeliveryRenderFunc func(props map[string]any) (g.Node, error)

// DeliveryExample is one stable authored fixture shared by native Figma and
// executable Storybook projections.
type DeliveryExample struct {
	ID          string
	Name        string
	Description string
	Props       map[string]any
	Design      blueprint.ExampleContract
}

// DeliveryDefinition is the complete OSS-owned delivery contribution for one
// concrete web component or solution.
type DeliveryDefinition struct {
	Identity uicomponent.Descriptor
	Contract designcomponent.Descriptor
	Design   blueprint.Definition
	Examples []DeliveryExample
	Render   DeliveryRenderFunc
}

// PropertyContract adds authored design semantics to fields deterministically
// discovered from a typed public Props struct.
type PropertyContract struct {
	Role        designcomponent.PropRole
	Description string
	Enum        []string
	Default     string
	Required    *bool
}

// Validate rejects an incomplete or internally inconsistent delivery
// definition before any adapter can publish it.
func (d DeliveryDefinition) Validate() error {
	if err := d.Identity.Validate(); err != nil {
		return fmt.Errorf("delivery identity: %w", err)
	}
	contract, err := d.Contract.Normalize()
	if err != nil {
		return fmt.Errorf("delivery component %q contract: %w", d.Identity.ID, err)
	}
	if contract.Name != string(d.Identity.ID) {
		return fmt.Errorf(
			"delivery component %q contract name is %q",
			d.Identity.ID,
			contract.Name,
		)
	}
	if want := categoryForTier(d.Identity.Tier); contract.Category != want {
		return fmt.Errorf(
			"delivery component %q category %q does not match tier %q",
			d.Identity.ID,
			contract.Category,
			d.Identity.Tier,
		)
	}
	if strings.TrimSpace(contract.ID) == "" {
		return fmt.Errorf("delivery component %q has no stable design id", d.Identity.ID)
	}
	if d.Design.SourceOfTruth == "" {
		return fmt.Errorf("delivery component %q has no visual authority", d.Identity.ID)
	}
	if d.Design.Root == nil {
		return fmt.Errorf("delivery component %q has no native visual root", d.Identity.ID)
	}
	if d.Render == nil {
		return fmt.Errorf("delivery component %q has no executable renderer", d.Identity.ID)
	}
	if len(d.Examples) == 0 {
		return fmt.Errorf("delivery component %q has no authored fixtures", d.Identity.ID)
	}
	seenExamples := make(map[string]struct{}, len(d.Examples))
	canonical := 0
	for index, example := range d.Examples {
		id := strings.TrimSpace(example.ID)
		if id == "" {
			return fmt.Errorf("delivery component %q example %d has no stable id", d.Identity.ID, index)
		}
		if _, duplicate := seenExamples[id]; duplicate {
			return fmt.Errorf("delivery component %q repeats example id %q", d.Identity.ID, id)
		}
		seenExamples[id] = struct{}{}
		if strings.TrimSpace(example.Name) == "" {
			return fmt.Errorf("delivery component %q example %q has no name", d.Identity.ID, id)
		}
		if example.Design.Canonical {
			canonical++
		}
		if viewport := strings.TrimSpace(example.Design.Viewport); viewport == "" {
			return fmt.Errorf("delivery component %q example %q has no viewport", d.Identity.ID, id)
		}
	}
	if canonical != 1 {
		return fmt.Errorf(
			"delivery component %q has %d canonical examples; expected exactly one",
			d.Identity.ID,
			canonical,
		)
	}
	if strings.TrimSpace(d.Design.CanonicalExample) == "" {
		return fmt.Errorf("delivery component %q has no canonical example name", d.Identity.ID)
	}
	return nil
}

// Clone returns a defensive copy safe for downstream assembly.
func (d DeliveryDefinition) Clone() DeliveryDefinition {
	cloned := d
	cloned.Contract.Props = slices.Clone(d.Contract.Props)
	for index := range cloned.Contract.Props {
		cloned.Contract.Props[index].EnumValues =
			slices.Clone(d.Contract.Props[index].EnumValues)
	}
	cloned.Contract.Slots = cloneDeliverySlots(d.Contract.Slots)
	cloned.Contract.Variants = slices.Clone(d.Contract.Variants)
	for index := range cloned.Contract.Variants {
		cloned.Contract.Variants[index].Values =
			slices.Clone(d.Contract.Variants[index].Values)
	}
	cloned.Contract.RequiredTokens = slices.Clone(d.Contract.RequiredTokens)
	cloned.Contract.Metadata = maps.Clone(d.Contract.Metadata)
	cloned.Design = cloneDeliveryDesign(d.Design)
	cloned.Examples = make([]DeliveryExample, len(d.Examples))
	for index, example := range d.Examples {
		cloned.Examples[index] = example
		cloned.Examples[index].Props = cloneDeliveryValueMap(example.Props)
		cloned.Examples[index].Design.Viewports =
			slices.Clone(example.Design.Viewports)
	}
	return cloned
}

func newDeliveryDefinition[T any](
	componentType string,
	tier uicomponent.Tier,
	stableID string,
	description string,
	properties map[string]PropertyContract,
	slots []designcomponent.Slot,
	design blueprint.Definition,
	examples []DeliveryExample,
	render func(T) g.Node,
) DeliveryDefinition {
	identity := uicomponent.Descriptor{
		ID:      uicomponent.ID(componentType),
		Tier:    tier,
		Owner:   DeliveryOwner,
		Version: DeliveryVersion,
	}
	props, err := propsFromType[T](properties)
	if err != nil {
		panic(fmt.Errorf("build OSS delivery component %q: %w", componentType, err))
	}
	contract := designcomponent.Descriptor{
		ID:            stableID,
		Name:          componentType,
		Category:      categoryForTier(tier),
		SourceOfTruth: designcomponent.SourceDefinition,
		Description:   description,
		Props:         props,
		Slots:         slots,
		Metadata: map[string]string{
			"owner":   DeliveryOwner,
			"version": DeliveryVersion,
		},
	}
	for _, prop := range props {
		if prop.Role != designcomponent.PropRoleVariant &&
			prop.Role != designcomponent.PropRoleTone &&
			prop.Role != designcomponent.PropRoleSize {
			continue
		}
		if prop.Type != designcomponent.PropEnum || len(prop.EnumValues) == 0 {
			continue
		}
		contract.Variants = append(contract.Variants, designcomponent.Variant{
			Name:        prop.Name,
			Values:      slices.Clone(prop.EnumValues),
			Default:     prop.Default,
			Description: prop.Description,
		})
	}
	normalized, err := contract.Normalize()
	if err != nil {
		panic(fmt.Errorf("build OSS delivery component %q: %w", componentType, err))
	}
	definition := DeliveryDefinition{
		Identity: identity,
		Contract: normalized,
		Design:   design,
		Examples: cloneDeliveryExamples(examples),
		Render: func(props map[string]any) (g.Node, error) {
			decoded, err := decodeDeliveryProps[T](props)
			if err != nil {
				return nil, fmt.Errorf("%s props: %w", componentType, err)
			}
			return render(decoded), nil
		},
	}
	if err := definition.Validate(); err != nil {
		panic(err)
	}
	return definition
}

func newLayoutDeliveryDefinition[T any](
	componentType string,
	stableID string,
	description string,
	properties map[string]PropertyContract,
	design blueprint.Definition,
	examples []DeliveryExample,
	render func(T, ...g.Node) g.Node,
) DeliveryDefinition {
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierTemplate,
		stableID,
		description,
		properties,
		[]designcomponent.Slot{{
			Name:         "children",
			Description:  "Ordered child components arranged by this layout.",
			AllowedTypes: deliveryTemplateContentTypes(),
			Cardinality:  designcomponent.SlotMany,
		}},
		design,
		examples,
		func(props T) g.Node {
			return render(
				props,
				sampleDeliveryChild("Primary content"),
				sampleDeliveryChild("Supporting content"),
			)
		},
	)
	return definition
}

func propsFromType[T any](
	overrides map[string]PropertyContract,
) ([]designcomponent.Prop, error) {
	typ := reflect.TypeFor[T]()
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("props type %s is not a struct", typ)
	}
	fields := make(map[string]designcomponent.Prop)
	if err := collectDeliveryProps(typ, fields); err != nil {
		return nil, err
	}
	for name, override := range overrides {
		prop, ok := fields[name]
		if !ok {
			return nil, fmt.Errorf("property contract %q is not present in %s", name, typ)
		}
		prop.Role = override.Role
		if description := strings.TrimSpace(override.Description); description != "" {
			prop.Description = description
		}
		if len(override.Enum) > 0 {
			prop.Type = designcomponent.PropEnum
			prop.EnumValues = slices.Clone(override.Enum)
		}
		prop.Default = strings.TrimSpace(override.Default)
		if override.Required != nil {
			prop.Required = *override.Required
		}
		fields[name] = prop
	}
	out := make([]designcomponent.Prop, 0, len(fields))
	for _, prop := range fields {
		out = append(out, prop)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func collectDeliveryProps(
	typ reflect.Type,
	props map[string]designcomponent.Prop,
) error {
	for index := range typ.NumField() {
		field := typ.Field(index)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("json")
		name, options := parseDeliveryJSONTag(tag)
		if name == "-" {
			continue
		}
		if field.Anonymous && name == "" {
			embedded := field.Type
			for embedded.Kind() == reflect.Pointer {
				embedded = embedded.Elem()
			}
			if embedded.Kind() == reflect.Struct {
				if err := collectDeliveryProps(embedded, props); err != nil {
					return err
				}
				continue
			}
		}
		if name == "" {
			name = field.Name
		}
		if _, duplicate := props[name]; duplicate {
			return fmt.Errorf("typed props repeat JSON field %q", name)
		}
		propType, err := deliveryPropType(field.Type)
		if err != nil {
			return fmt.Errorf("typed prop %q: %w", name, err)
		}
		props[name] = designcomponent.Prop{
			Name:        name,
			Type:        propType,
			Required:    !options["omitempty"] && !options["omitzero"],
			Description: "Typed " + name + " property.",
		}
	}
	return nil
}

func parseDeliveryJSONTag(tag string) (string, map[string]bool) {
	parts := strings.Split(tag, ",")
	options := make(map[string]bool, len(parts))
	for _, option := range parts[1:] {
		options[strings.TrimSpace(option)] = true
	}
	return strings.TrimSpace(parts[0]), options
}

func deliveryPropType(typ reflect.Type) (designcomponent.PropType, error) {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.String:
		return designcomponent.PropString, nil
	case reflect.Bool:
		return designcomponent.PropBoolean, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return designcomponent.PropNumber, nil
	case reflect.Array, reflect.Slice:
		return designcomponent.PropArray, nil
	case reflect.Map, reflect.Struct:
		return designcomponent.PropObject, nil
	default:
		return "", fmt.Errorf("Go type %s is unsupported", typ)
	}
}

func decodeDeliveryProps[T any](props map[string]any) (T, error) {
	var target T
	data, err := json.Marshal(props)
	if err != nil {
		return target, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&target); err != nil {
		return target, err
	}
	return target, nil
}

func categoryForTier(tier uicomponent.Tier) designcomponent.Category {
	switch tier {
	case uicomponent.TierAtom:
		return designcomponent.CategoryAtom
	case uicomponent.TierMolecule:
		return designcomponent.CategoryMolecule
	case uicomponent.TierOrganism:
		return designcomponent.CategoryOrganism
	case uicomponent.TierTemplate:
		return designcomponent.CategoryTemplate
	case uicomponent.TierPage:
		return designcomponent.CategoryPage
	default:
		return ""
	}
}

func canonicalDeliveryExample(
	componentType string,
	name string,
	props map[string]any,
) DeliveryExample {
	return DeliveryExample{
		ID:    deliveryExampleID(componentType, name),
		Name:  name,
		Props: cloneDeliveryValueMap(props),
		Design: blueprint.ExampleContract{
			Canonical: true,
			Viewport:  "desktop",
			Viewports: []string{"desktop", "mobile"},
			ColorMode: "light",
			Density:   "medium",
		},
	}
}

func deliveryExample(componentType, name string, props map[string]any) DeliveryExample {
	example := canonicalDeliveryExample(componentType, name, props)
	example.Design.Canonical = false
	return example
}

func mobileDeliveryExample(componentType, name string, props map[string]any) DeliveryExample {
	example := deliveryExample(componentType, name, props)
	example.Design.Viewport = "mobile"
	example.Design.Viewports = []string{"mobile"}
	return example
}

func deliveryExampleID(componentType, name string) string {
	return "pk-ui.example/" + strings.ToLower(componentType) + "/" +
		strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(strings.ToLower(name))
}

func stableDeliveryID(componentType string) string {
	return "pk-ui.component." + strings.ToLower(componentType)
}

func sampleDeliveryChild(label string) g.Node {
	return Card(
		cardProps(label),
		Text(textProps("Editable child slot", "muted")),
	)
}

func cloneDeliverySlots(values []designcomponent.Slot) []designcomponent.Slot {
	if len(values) == 0 {
		return nil
	}
	out := make([]designcomponent.Slot, len(values))
	for index, value := range values {
		out[index] = value
		out[index].AllowedTypes = slices.Clone(value.AllowedTypes)
		out[index].Attrs = slices.Clone(value.Attrs)
		for attrIndex := range out[index].Attrs {
			out[index].Attrs[attrIndex].EnumValues =
				slices.Clone(value.Attrs[attrIndex].EnumValues)
		}
		if value.Scope != nil {
			out[index].Scope = &designcomponent.SlotScope{
				Fields: slices.Clone(value.Scope.Fields),
			}
		}
	}
	return out
}

func cloneDeliveryDesign(value blueprint.Definition) blueprint.Definition {
	out := value
	out.VariantExamples = slices.Clone(value.VariantExamples)
	out.VariantMatrix = slices.Clone(value.VariantMatrix)
	if value.ExampleDefaults != nil {
		copy := *value.ExampleDefaults
		copy.Viewports = slices.Clone(value.ExampleDefaults.Viewports)
		out.ExampleDefaults = &copy
	}
	if value.Taxonomy != nil {
		copy := *value.Taxonomy
		out.Taxonomy = &copy
	}
	out.Root = cloneDeliveryNode(value.Root)
	out.Assets = slices.Clone(value.Assets)
	out.InteractiveStates = slices.Clone(value.InteractiveStates)
	return out
}

func cloneDeliveryNode(value *blueprint.Node) *blueprint.Node {
	if value == nil {
		return nil
	}
	out := *value
	out.ClassBindings = make(map[string]map[string]string, len(value.ClassBindings))
	for key, bindings := range value.ClassBindings {
		out.ClassBindings[key] = maps.Clone(bindings)
	}
	out.Props = cloneDeliveryValueMap(value.Props)
	out.Children = make([]blueprint.Node, len(value.Children))
	for index := range value.Children {
		child := cloneDeliveryNode(&value.Children[index])
		out.Children[index] = *child
	}
	return &out
}

func cloneDeliveryExamples(values []DeliveryExample) []DeliveryExample {
	out := make([]DeliveryExample, len(values))
	for index, value := range values {
		out[index] = value
		out[index].Props = cloneDeliveryValueMap(value.Props)
		out[index].Design.Viewports = slices.Clone(value.Design.Viewports)
	}
	return out
}

func cloneDeliveryValueMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return nil
	}
	data, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		panic(err)
	}
	return out
}

func deliveryAtomicContentTypes() []string {
	return []string{
		"Badge", "Button", "Divider", "Heading", "Kbd", "Label", "Link",
		"Spinner", "Tag", "Text",
	}
}

func deliveryTemplateContentTypes() []string {
	return []string{
		"Alert", "Badge", "Breadcrumb", "Button", "Card", "Checkbox",
		"Container", "DataGrid", "Divider", "EmptyState", "Flex", "Grid",
		"Heading", "Input", "Kbd", "Label", "Link", "Pagination", "SearchBar",
		"Select", "Spinner", "Stack", "Table", "Tabs", "Tag", "Text", "Textarea",
	}
}
