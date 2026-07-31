// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_catalog.go is the explicit OSS design-delivery catalog. Every
// entry binds a typed public contract to its production renderer, native
// editable blueprint, stable identity, and reviewed fixtures. There is no
// reflection-based renderer lookup, init registration, or proprietary import.

import (
	"fmt"
	"slices"
	"strings"

	g "maragu.dev/gomponents"

	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/tw"
)

var (
	buttonVariants = []string{
		"primary", "secondary", "outline", "ghost", "link",
	}
	controlSizes = []string{"xs", "sm", "md", "lg", "xl", "2xl"}
	gaps         = []string{"0", "1", "2", "3", "4", "5", "6", "8"}
)

// OSSDeliveryCatalog returns the complete renderer-backed OSS web library in
// deterministic atomic-design order. Returned definitions are defensive
// copies; callers can safely adapt them into isolated runtimes.
func OSSDeliveryCatalog() []DeliveryDefinition {
	definitions := []DeliveryDefinition{
		iconDeliveryDefinition(),
		buttonDeliveryDefinition(),
		badgeDeliveryDefinition(),
		alertDeliveryDefinition(),
		inputDeliveryDefinition(),
		selectDeliveryDefinition(),
		textareaDeliveryDefinition(),
		checkboxDeliveryDefinition(),
		labelDeliveryDefinition(),
		textDeliveryDefinition(),
		headingDeliveryDefinition(),
		dividerDeliveryDefinition(),
		spinnerDeliveryDefinition(),
		skeletonDeliveryDefinition(),
		deferredSlotDeliveryDefinition(),
		emptyStateDeliveryDefinition(),
		kbdDeliveryDefinition(),
		linkDeliveryDefinition(),
		tagDeliveryDefinition(),
		tableDeliveryDefinition(),
		tableSkeletonDeliveryDefinition(),
		cardDeliveryDefinition(),
		cardSkeletonDeliveryDefinition(),
		breadcrumbDeliveryDefinition(),
		paginationDeliveryDefinition(),
		searchBarDeliveryDefinition(),
		tabsDeliveryDefinition(),
		dataGridDeliveryDefinition(),
		windowedCollectionDeliveryDefinition(),
		stackDeliveryDefinition(),
		flexDeliveryDefinition(),
		gridDeliveryDefinition(),
		containerDeliveryDefinition(),
		dataManagementPageDeliveryDefinition(),
	}
	if err := validateDeliveryCatalogExamples(definitions); err != nil {
		panic(fmt.Errorf("validate OSS delivery examples: %w", err))
	}
	out := make([]DeliveryDefinition, len(definitions))
	for index, definition := range definitions {
		if err := definition.Validate(); err != nil {
			panic(fmt.Errorf("OSS delivery catalog entry %d: %w", index, err))
		}
		out[index] = definition.Clone()
	}
	return out
}

// RenderDeliveryExample renders one authored example through the exact same
// catalog contributions used by production runtimes. Root prop overrides are
// intended for Storybook controls; the authored child graph remains explicit
// and stable.
func RenderDeliveryExample(
	componentType string,
	exampleID string,
	rootOverrides map[string]any,
) (g.Node, error) {
	catalog := OSSDeliveryCatalog()
	index := make(map[string]DeliveryDefinition, len(catalog))
	for _, definition := range catalog {
		index[string(definition.Identity.ID)] = definition
	}
	definition, exists := index[strings.TrimSpace(componentType)]
	if !exists {
		return nil, fmt.Errorf("OSS delivery component %q is not registered", componentType)
	}
	for _, example := range definition.Examples {
		if example.ID != strings.TrimSpace(exampleID) {
			continue
		}
		return renderDeliveryExampleWithCatalog(
			definition,
			example,
			rootOverrides,
			index,
		)
	}
	return nil, fmt.Errorf(
		"OSS delivery component %q has no example %q",
		componentType,
		exampleID,
	)
}

func validateDeliveryCatalogExamples(definitions []DeliveryDefinition) error {
	index := make(map[string]DeliveryDefinition, len(definitions))
	for _, definition := range definitions {
		componentType := string(definition.Identity.ID)
		if _, duplicate := index[componentType]; duplicate {
			return fmt.Errorf("duplicate delivery component type %q", componentType)
		}
		index[componentType] = definition
	}
	for _, definition := range definitions {
		for _, example := range definition.Examples {
			node, err := renderDeliveryExampleWithCatalog(
				definition,
				example,
				nil,
				index,
			)
			if err != nil {
				return fmt.Errorf(
					"%s/%s: %w",
					definition.Identity.ID,
					example.ID,
					err,
				)
			}
			if node == nil {
				return fmt.Errorf(
					"%s/%s rendered a nil node",
					definition.Identity.ID,
					example.ID,
				)
			}
		}
	}
	return nil
}

func renderDeliveryExampleWithCatalog(
	definition DeliveryDefinition,
	example DeliveryExample,
	rootOverrides map[string]any,
	catalog map[string]DeliveryDefinition,
) (g.Node, error) {
	seenIDs := map[string]struct{}{example.ID: {}}
	slots, err := renderDeliveryExampleSlots(
		example.Slots,
		definition,
		example.ID,
		catalog,
		seenIDs,
	)
	if err != nil {
		return nil, err
	}
	return definition.Render(
		mergeDeliveryExampleProps(example.Props, rootOverrides),
		slots,
	)
}

func renderDeliveryExampleSlots(
	exampleSlots []DeliveryExampleSlot,
	parent DeliveryDefinition,
	path string,
	catalog map[string]DeliveryDefinition,
	seenIDs map[string]struct{},
) (DeliverySlotChildren, error) {
	if err := validateDeliveryExampleSlots(
		string(parent.Identity.ID),
		parent.Contract.Slots,
		exampleSlots,
		path,
	); err != nil {
		return nil, err
	}
	if len(exampleSlots) == 0 {
		return nil, nil
	}
	contracts := make(map[string]designcomponent.Slot, len(parent.Contract.Slots))
	for _, slot := range parent.Contract.Slots {
		contracts[slot.Name] = slot
	}
	out := make(DeliverySlotChildren, len(exampleSlots))
	for _, slot := range exampleSlots {
		contract := contracts[slot.Name]
		for _, component := range slot.Components {
			if _, duplicate := seenIDs[component.ID]; duplicate {
				return nil, fmt.Errorf(
					"example %q repeats component id %q",
					path,
					component.ID,
				)
			}
			seenIDs[component.ID] = struct{}{}
			definition, exists := catalog[component.Type]
			if !exists {
				return nil, fmt.Errorf(
					"example %q component %q references unknown type %q",
					path,
					component.ID,
					component.Type,
				)
			}
			childPath := path + "/" + component.ID
			childSlots, err := renderDeliveryExampleSlots(
				component.Slots,
				definition,
				childPath,
				catalog,
				seenIDs,
			)
			if err != nil {
				return nil, err
			}
			node, err := definition.Render(component.Props, childSlots)
			if err != nil {
				return nil, fmt.Errorf(
					"render example %q component %q: %w",
					path,
					component.ID,
					err,
				)
			}
			attrs := cloneDeliveryValueMap(component.SlotAttrs)
			if deliverySlotDeclaresAttr(contract, "id") {
				if attrs == nil {
					attrs = make(map[string]any)
				}
				if _, exists := attrs["id"]; !exists {
					attrs["id"] = component.ID
				}
			}
			out[slot.Name] = append(out[slot.Name], DeliverySlotInstance{
				Attrs:    attrs,
				Children: []g.Node{node},
			})
		}
	}
	return out, nil
}

func deliverySlotDeclaresAttr(slot designcomponent.Slot, name string) bool {
	for _, attr := range slot.Attrs {
		if attr.Name == name {
			return true
		}
	}
	return false
}

func mergeDeliveryExampleProps(
	base map[string]any,
	overrides map[string]any,
) map[string]any {
	out := cloneDeliveryValueMap(base)
	if len(overrides) == 0 {
		return out
	}
	if out == nil {
		out = make(map[string]any, len(overrides))
	}
	for key, value := range overrides {
		out[key] = value
	}
	return out
}

func iconDeliveryDefinition() DeliveryDefinition {
	componentType := "Icon"
	root := deliveryClassBound(
		deliveryClassBound(
			deliverySVGBound(
				componentType,
				clIcon.Compile(),
				"icon-",
				"name",
			),
			"size",
			compileClassMap(clIconSize),
		),
		"tone",
		compileClassMap(clIconTone),
	)
	return withDeliveryTags(newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Reusable system glyph with governed size, semantic tone, and accessibility.",
		map[string]PropertyContract{
			"name": contentProperty("Canonical glyph name."),
			"size": sizeProperty(controlSizes, "md", "Rendered glyph size."),
			"tone": toneProperty(
				[]string{"neutral", "brand", "success", "warning", "danger", "info"},
				"neutral",
				"Semantic foreground tone.",
			),
			"weight":    variantProperty([]string{"outline"}, "outline", "Glyph weight."),
			"ariaLabel": contentProperty("Accessible label; empty marks the glyph decorative."),
		},
		nil,
		deliveryDesign("Search / Fg Primary / 20", "01 Icons", "System Glyphs", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "Search / Fg Primary / 20", map[string]any{
				"name": "search", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Search / Fg Tertiary / 20", map[string]any{
				"name": "search", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "X Mark / Fg Primary / 24", map[string]any{
				"name": "x-mark", "tone": "neutral", "size": "lg", "weight": "outline",
			}),
			deliveryExample(componentType, "Arrow Right / Fg Primary / 24", map[string]any{
				"name": "arrow-right", "tone": "neutral", "size": "lg", "weight": "outline",
			}),
			deliveryExample(componentType, "Calendar / Fg Tertiary / 20", map[string]any{
				"name": "calendar", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Chat Bubble Left Right / Fg On Brand / 24", map[string]any{
				"name": "chat-bubble-left-right", "tone": "neutral", "size": "lg", "weight": "outline",
			}),
			deliveryExample(componentType, "Check / Fg Brand / 16", map[string]any{
				"name": "check", "tone": "brand", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "Check / Fg On Brand / 16", map[string]any{
				"name": "check", "tone": "neutral", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "Check Circle / Fg Primary / 20", map[string]any{
				"name": "check-circle", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Check Circle / Fg Success / 20", map[string]any{
				"name": "check-circle", "tone": "success", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Chevron Down / Fg Tertiary / 16", map[string]any{
				"name": "chevron-down", "tone": "neutral", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "Chevron Down / Fg Tertiary / 20", map[string]any{
				"name": "chevron-down", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Clock / Fg Primary / 20", map[string]any{
				"name": "clock", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Cog / Fg Primary / 24", map[string]any{
				"name": "cog", "tone": "neutral", "size": "lg", "weight": "outline",
			}),
			deliveryExample(componentType, "Document / Fg Tertiary / 40", map[string]any{
				"name": "document", "tone": "neutral", "size": "2xl", "weight": "outline",
			}),
			deliveryExample(componentType, "Envelope / Fg Primary / 20", map[string]any{
				"name": "envelope", "tone": "neutral", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Ellipsis Vertical / Fg Primary / 12", map[string]any{
				"name": "ellipsis-vertical", "tone": "neutral", "size": "xs", "weight": "outline",
			}),
			deliveryExample(componentType, "Exclamation Circle / Fg Warning / 16", map[string]any{
				"name": "exclamation-circle", "tone": "warning", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "Exclamation Triangle / Fg Warning / 20", map[string]any{
				"name": "exclamation-triangle", "tone": "warning", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Information Circle / Fg Info / 20", map[string]any{
				"name": "information-circle", "tone": "info", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "Trash / Fg Danger / 16", map[string]any{
				"name": "trash", "tone": "danger", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "Upload / Fg Tertiary / 40", map[string]any{
				"name": "upload", "tone": "neutral", "size": "2xl", "weight": "outline",
			}),
			deliveryExample(componentType, "User / Fg Tertiary / 16", map[string]any{
				"name": "user", "tone": "neutral", "size": "sm", "weight": "outline",
			}),
			deliveryExample(componentType, "X Circle / Fg Danger / 20", map[string]any{
				"name": "x-circle", "tone": "danger", "size": "md", "weight": "outline",
			}),
			deliveryExample(componentType, "X Circle / Fg Primary / 20", map[string]any{
				"name": "x-circle", "tone": "neutral", "size": "md", "weight": "outline",
			}),
		},
		Icon,
	), "icon", "glyph")
}

func buttonDeliveryDefinition() DeliveryDefinition {
	componentType := "Button"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryClassBound(
				deliveryClassBound(
					deliveryFrame(
						componentType,
						clButtonBase.Compile(),
						deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
						deliveryHiddenWhen(
							deliveryText("Label", "", "Continue", "label"),
							map[string]any{"iconOnly": true},
						),
						deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
					),
					"iconOnly",
					map[string]string{
						"false": "",
						"true":  clButtonIconOnly.Compile(),
					},
				),
				"variant",
				compileClassMap(clButtonVariant),
			),
			"tone",
			compileClassMap(clButtonTone),
		),
		"size",
		compileClassMap(clButtonSize),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible action control rendered by the OSS web runtime.",
		map[string]PropertyContract{
			"label":     contentProperty("Visible button label."),
			"variant":   variantProperty(buttonVariants, "primary", "Visual treatment."),
			"tone":      toneProperty([]string{"neutral", "brand", "success", "warning", "danger", "info"}, "neutral", "Semantic intent."),
			"size":      sizeProperty(controlSizes, "md", "Control size."),
			"type":      enumProperty([]string{"button", "submit", "reset"}, "button", "HTML button type."),
			"loading":   stateProperty("Whether progress replaces the leading content."),
			"fullWidth": modifierProperty("Whether the control fills its container."),
			"iconOnly":  modifierProperty("Whether the visible label is suppressed for an icon-only control."),
			"ariaLabel": contentProperty("Accessible name required for icon-only controls."),
		},
		[]designcomponent.Slot{
			{
				Name:         "iconStart",
				Description:  "Optional leading icon slot.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "iconEnd",
				Description:  "Optional trailing icon slot.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("primary", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "primary", map[string]any{
				"label": "Continue", "variant": "primary", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "secondary", map[string]any{
				"label": "Cancel", "variant": "secondary", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "loading", map[string]any{
				"label": "Saving", "variant": "primary", "tone": "neutral", "size": "md", "loading": true,
			}),
			deliveryExample(componentType, "danger", map[string]any{
				"label": "Sign out", "variant": "primary", "tone": "danger", "size": "md",
			}),
			deliveryExample(componentType, "outline", map[string]any{
				"label": "Details", "variant": "outline", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "ghost", map[string]any{
				"label": "Cancel", "variant": "ghost", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "link", map[string]any{
				"label": "View all", "variant": "link", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "full-width", map[string]any{
				"label": "Continue", "variant": "primary", "tone": "neutral", "size": "md", "fullWidth": true,
			}),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "button-icon-only-ghost", map[string]any{
					"label": "More actions", "variant": "ghost", "tone": "neutral", "size": "xs",
					"iconOnly": true, "ariaLabel": "More actions",
				}),
				deliveryExampleSlot(
					"iconStart",
					deliveryExampleComponent(
						"more-actions-icon",
						"Icon",
						map[string]any{
							"name": "ellipsis-vertical",
							"size": "xs",
							"tone": "neutral",
						},
					),
				),
			),
		},
		func(props atoms.ButtonProps, slots DeliverySlotChildren) g.Node {
			return buttonWithSlots(
				props,
				slots.Nodes("iconStart"),
				slots.Nodes("iconEnd"),
			)
		},
	)
}

func badgeDeliveryDefinition() DeliveryDefinition {
	componentType := "Badge"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryClassBound(
				deliveryFrame(
					componentType,
					clBadgeBase.Compile(),
					deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
					deliveryText("Label", "", "Active", "label"),
					deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
				),
				"variant",
				compileClassMap(clBadgeVariant),
			),
			"tone",
			compileClassMap(clBadgeTone),
		),
		"size",
		compileClassMap(clBadgeSize),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Compact semantic status label.",
		map[string]PropertyContract{
			"label":   contentProperty("Visible badge label."),
			"variant": variantProperty([]string{"primary", "secondary", "outline"}, "primary", "Visual treatment."),
			"tone": toneProperty(
				[]string{"neutral", "brand", "success", "warning", "danger", "info"},
				"neutral",
				"Semantic intent independent of the visual treatment.",
			),
			"size": sizeProperty(controlSizes, "md", "Badge scale."),
			"dot":  modifierProperty("Shows a leading status dot."),
		},
		[]designcomponent.Slot{
			{
				Name:         "iconStart",
				Description:  "Optional leading icon slot.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "iconEnd",
				Description:  "Optional trailing icon slot.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("default", "02 Atoms", "Status", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"label": "Active", "variant": "primary", "tone": "neutral", "size": "md",
			}),
			deliveryExample(componentType, "success", map[string]any{
				"label": "Published", "variant": "primary", "tone": "success", "size": "md", "dot": true,
			}),
			deliveryExample(componentType, "tone-info", map[string]any{
				"label": "Pending", "variant": "secondary", "tone": "info", "size": "md",
			}),
		},
		func(props atoms.BadgeProps, slots DeliverySlotChildren) g.Node {
			return badgeWithSlots(
				props,
				slots.Nodes("iconStart"),
				slots.Nodes("iconEnd"),
			)
		},
	)
}

func alertDeliveryDefinition() DeliveryDefinition {
	componentType := "Alert"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clAlertBase.Compile(),
			deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
			deliveryFrame(
				"Body",
				clAlertBody.Compile(),
				deliveryText("Title", clAlertTitle.Compile(), "Update available", "title"),
				deliveryText("Message", clAlertMessage.Compile(), "Review the latest changes.", "message"),
			),
			deliverySlot("Actions", "actions", ""),
		),
		"tone",
		compileClassMap(clAlertVariant),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Persistent accessible inline status message.",
		map[string]PropertyContract{
			"message":     contentProperty("Required status message."),
			"title":       contentProperty("Optional status title."),
			"tone":        toneProperty([]string{"info", "success", "warning", "danger"}, "info", "Severity."),
			"dismissible": stateProperty("Whether a dismiss affordance is exposed."),
			"bordered":    modifierProperty("Whether the message has an outline."),
			"compact":     modifierProperty("Whether spacing is compact."),
		},
		[]designcomponent.Slot{
			{
				Name:         "iconStart",
				Description:  "Optional leading status glyph.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "actions",
				Description:  "Optional recovery or contextual actions.",
				AllowedTypes: []string{"Button", "Link"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("info", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "info", map[string]any{
				"title": "Update available", "message": "Review the latest changes.", "tone": "info",
			}),
			deliveryExample(componentType, "warning", map[string]any{
				"title": "Check required", "message": "One field needs attention.", "tone": "warning",
			}),
			deliveryExample(componentType, "info-alert", map[string]any{
				"title": "Heads up", "message": "Your session will expire in 10 minutes.", "tone": "info",
			}),
		},
		func(props atoms.AlertProps, slots DeliverySlotChildren) g.Node {
			return alertWithSlots(
				props,
				slots.Nodes("iconStart"),
				slots.Nodes("actions"),
			)
		},
	)
}

func inputDeliveryDefinition() DeliveryDefinition {
	componentType := "Input"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFieldWrap.Compile(),
			deliverySlot(
				"LabelSlot",
				"label",
				"",
				deliveryText("Label", clLabel.Compile(), "Email address", "label"),
			),
			deliveryClassBound(
				deliveryFrame(
					"Control",
					clInput.Compile()+" "+clInputNormal.Compile(),
					deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
					deliveryText("Value", "", "name@example.com", "value", "placeholder"),
					deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
				),
				"error",
				map[string]string{"": clInputNormal.Compile(), "*": clInputError.Compile()},
			),
			deliveryText("Help", clHelp.Compile(), "We never share your email.", "helpText", "error"),
		),
		"field",
		"field",
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Labelled single-line form control with help and error states.",
		map[string]PropertyContract{
			"name": contentProperty("Submitted field name."),
			"type": variantProperty(
				[]string{
					"text", "email", "password", "number", "tel", "url",
					"search", "date", "time", "datetime-local", "month",
					"week", "color", "hidden",
				},
				"text",
				"Input kind.",
			),
			"value":       contentProperty("Current value."),
			"placeholder": contentProperty("Empty-state prompt."),
			"label":       contentProperty("Visible field label."),
			"helpText":    contentProperty("Supporting guidance."),
			"error":       stateProperty("Validation message; non-empty selects error styling."),
			"required":    stateProperty("Whether a value is required."),
			"readOnly":    stateProperty("Whether editing is prevented."),
			"autoFocus":   stateProperty("Whether focus is requested on mount."),
			"size":        sizeProperty([]string{"sm", "md", "lg"}, "md", "Control size."),
		},
		[]designcomponent.Slot{
			{
				Name:         "iconStart",
				Description:  "Optional leading icon or textual adornment slot.",
				AllowedTypes: []string{"Icon", "Text"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "iconEnd",
				Description:  "Optional trailing icon or textual adornment slot.",
				AllowedTypes: []string{"Icon", "Text"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("text", "02 Atoms", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "text", map[string]any{
				"name": "email", "type": "email", "label": "Email address",
				"placeholder": "name@example.com", "helpText": "We never share your email.",
			}),
			deliveryExample(componentType, "error", map[string]any{
				"name": "email", "type": "email", "label": "Email address",
				"value": "invalid", "error": "Enter a valid email address.",
			}),
			deliveryExample(componentType, "email", map[string]any{
				"name": "email", "type": "email", "label": "Email address",
				"placeholder": "user@example.com",
			}),
			deliveryExample(componentType, "password", map[string]any{
				"name": "password", "type": "password", "label": "Password",
				"value": "correct horse battery staple",
			}),
		},
		func(props atoms.InputProps, slots DeliverySlotChildren) g.Node {
			return inputWithSlots(
				props,
				slots.Nodes("iconStart"),
				slots.Nodes("iconEnd"),
			)
		},
	)
}

func selectDeliveryDefinition() DeliveryDefinition {
	componentType := "Select"
	root := deliveryFrame(
		componentType,
		clFieldWrap.Compile(),
		deliveryText("Label", clLabel.Compile(), "Status", "label"),
		deliveryFrame(
			"Control",
			clInput.Compile()+" "+clInputNormal.Compile(),
			deliveryText("Selection", "", "Choose a status", "value", "placeholder"),
		),
		deliveryText("Help", clHelp.Compile(), "Select one option.", "helpText", "error"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Native single- or multiple-value choice control.",
		map[string]PropertyContract{
			"name":        contentProperty("Submitted field name."),
			"label":       contentProperty("Visible field label."),
			"value":       contentProperty("Selected option value."),
			"values":      contentProperty("Selected option values in multiple mode."),
			"placeholder": contentProperty("Unselected prompt."),
			"options":     contentProperty("Available choices."),
			"required":    stateProperty("Whether a choice is required."),
			"multiple":    stateProperty("Whether more than one choice can be selected."),
			"visibleRows": roleProperty(designcomponent.PropRoleSize, "Visible option rows in multiple mode."),
			"helpText":    contentProperty("Supporting guidance."),
			"error":       stateProperty("Validation message."),
		},
		nil,
		deliveryDesign("single-select", "02 Atoms", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "single-select", map[string]any{
				"name": "status", "label": "Status", "placeholder": "Choose a status",
				"options": []map[string]any{
					{"label": "Draft", "value": "draft"},
					{"label": "Published", "value": "published"},
				},
			}),
			mobileDeliveryExample(componentType, "multiple-select", map[string]any{
				"name": "status", "label": "Status", "multiple": true,
				"visibleRows": 3, "values": []string{"draft", "published"},
				"options": []map[string]any{
					{"label": "Draft", "value": "draft"},
					{"label": "Published", "value": "published"},
					{"label": "Archived", "value": "archived"},
				},
			}),
		},
		Select,
	)
}

func textareaDeliveryDefinition() DeliveryDefinition {
	componentType := "Textarea"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFieldWrap.Compile(),
			deliveryText("Label", clLabel.Compile(), "Notes", "label"),
			deliveryFrame(
				"Control",
				clInput.Compile()+" "+clInputNormal.Compile(),
				deliveryText("Value", "", "Add context for reviewers.", "value", "placeholder"),
			),
			deliveryText("Help", clHelp.Compile(), "Keep it concise.", "helperText", "errorMessage"),
		),
		"field",
		"field",
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Labelled multi-line form control.",
		map[string]PropertyContract{
			"name":         contentProperty("Submitted field name."),
			"placeholder":  contentProperty("Empty-state prompt."),
			"value":        contentProperty("Current value."),
			"label":        contentProperty("Visible field label."),
			"helperText":   contentProperty("Supporting guidance."),
			"errorMessage": stateProperty("Validation message."),
			"required":     stateProperty("Whether a value is required."),
			"readOnly":     stateProperty("Whether editing is prevented."),
			"rows":         roleProperty(designcomponent.PropRoleSize, "Visible row count."),
			"minLength":    roleProperty(designcomponent.PropRoleSize, "Minimum accepted character count."),
			"maxLength":    roleProperty(designcomponent.PropRoleSize, "Maximum accepted character count."),
			"fullWidth":    modifierProperty("Whether the field fills its container."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"name": "notes", "label": "Notes", "placeholder": "Add context for reviewers.",
				"helperText": "Keep it concise.", "rows": 4,
			}),
			deliveryExample(componentType, "auto-resize", map[string]any{
				"name": "description", "label": "Description",
				"placeholder": "Enter a detailed description...",
				"maxLength":   500, "fullWidth": true,
			}),
		},
		Textarea,
	)
}

func checkboxDeliveryDefinition() DeliveryDefinition {
	componentType := "Checkbox"
	root := deliveryFrame(
		componentType,
		clFieldWrap.Compile(),
		deliveryFrame(
			"ControlRow",
			clCheckRow.Compile(),
			deliveryFrame("Control", clCheckbox.Compile()),
			deliveryText("Label", clLabel.Compile(), "Include archived items", "label"),
		),
		deliveryText("Help", clHelp.Compile(), "Archived rows remain read-only.", "helpText"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible labelled checkbox control.",
		map[string]PropertyContract{
			"name":          contentProperty("Submitted field name."),
			"label":         contentProperty("Visible option label."),
			"checked":       stateProperty("Checked state."),
			"indeterminate": stateProperty("Mixed selection state."),
			"value":         contentProperty("Submitted value."),
			"required":      stateProperty("Whether selection is required."),
			"helpText":      contentProperty("Supporting guidance."),
		},
		nil,
		deliveryDesign("unchecked", "02 Atoms", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "unchecked", map[string]any{
				"name": "archived", "label": "Include archived items",
				"helpText": "Archived rows remain read-only.",
			}),
			deliveryExample(componentType, "checked", map[string]any{
				"name": "archived", "label": "Include archived items", "checked": true,
			}),
			deliveryExample(componentType, "basic", map[string]any{
				"name": "agree", "label": "I agree to the terms",
			}),
			deliveryExample(componentType, "required", map[string]any{
				"name": "terms", "label": "Accept terms of service", "required": true,
			}),
		},
		Checkbox,
	)
}

func labelDeliveryDefinition() DeliveryDefinition {
	componentType := "Label"
	root := deliveryFrame(
		componentType,
		clLabel.Compile(),
		deliveryText("Text", "", "Display name", "text"),
		deliveryText("Required", clRequired.Compile(), " *"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Form-label text with a native required marker.",
		map[string]PropertyContract{
			"text":     contentProperty("Visible label text."),
			"for":      contentProperty("Associated field identity."),
			"required": stateProperty("Whether the required marker is shown."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Typography", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"text": "Display name", "for": "display-name",
			}),
			deliveryExample(componentType, "required", map[string]any{
				"text": "Display name", "for": "display-name", "required": true,
			}),
		},
		Label,
	)
}

func textDeliveryDefinition() DeliveryDefinition {
	componentType := "Text"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryClassBound(
				deliveryText("Text", "", "A clear supporting sentence.", "content"),
				"size",
				textSizeClasses(),
			),
			"weight",
			textWeightClasses(),
		),
		"color",
		textColorClasses(),
	)
	return withDeliveryTags(newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Semantic body-copy primitive.",
		map[string]PropertyContract{
			"content":  contentProperty("Visible copy."),
			"size":     sizeProperty([]string{"xs", "sm", "base", "lg", "xl", "2xl"}, "base", "Text scale."),
			"weight":   variantProperty([]string{"normal", "medium", "semibold", "bold"}, "normal", "Font weight."),
			"color":    toneProperty([]string{"primary", "secondary", "muted", "brand", "success", "warning", "danger", "info"}, "primary", "Semantic foreground."),
			"truncate": modifierProperty("Whether overflow is truncated."),
			"lines":    modifierProperty("Optional line-clamp count."),
		},
		nil,
		deliveryDesign("body", "02 Atoms", "Typography", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "body", map[string]any{
				"content": "A clear supporting sentence.", "size": "base",
				"weight": "normal", "color": "primary",
			}),
			deliveryExample(componentType, "muted", map[string]any{
				"content": "Supporting metadata", "size": "sm", "color": "muted",
			}),
			deliveryExample(componentType, "muted-body", map[string]any{
				"content": "Muted supporting copy", "size": "base", "color": "muted",
			}),
		},
		Text,
	), "typography")
}

func headingDeliveryDefinition() DeliveryDefinition {
	componentType := "Heading"
	root := deliveryClassBound(
		deliveryText(
			"Heading",
			clHeadingBase.Compile(),
			"Workspace overview",
			"text",
		),
		"level",
		headingLevelClasses(),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Semantic heading primitive from level one through six.",
		map[string]PropertyContract{
			"text":     contentProperty("Visible heading text."),
			"level":    roleProperty(designcomponent.PropRoleSize, "Semantic heading level."),
			"anchor":   contentProperty("Optional anchor identity."),
			"truncate": modifierProperty("Whether overflow is truncated."),
		},
		nil,
		deliveryDesign("level-two", "02 Atoms", "Typography", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "level-two", map[string]any{
				"text": "Workspace overview", "level": 2,
			}),
			deliveryExample(componentType, "level-one", map[string]any{
				"text": "Operations", "level": 1,
			}),
		},
		Heading,
	)
}

func dividerDeliveryDefinition() DeliveryDefinition {
	componentType := "Divider"
	root := deliveryClassBound(
		deliveryFrame(componentType, clDividerH.Compile()),
		"orientation",
		map[string]string{
			"horizontal": clDividerH.Compile(),
			"vertical":   clDividerV.Compile(),
		},
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Horizontal or vertical semantic separator.",
		map[string]PropertyContract{
			"orientation": variantProperty([]string{"horizontal", "vertical"}, "horizontal", "Separator direction."),
			"text":        contentProperty("Optional separator label."),
		},
		nil,
		deliveryDesign("horizontal", "02 Atoms", "Structure", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "horizontal", map[string]any{
				"orientation": "horizontal",
			}),
			deliveryExample(componentType, "with-text", map[string]any{
				"orientation": "horizontal",
				"text":        "Or continue with",
			}),
			deliveryExample(componentType, "vertical", map[string]any{
				"orientation": "vertical",
			}),
		},
		Divider,
	)
}

func spinnerDeliveryDefinition() DeliveryDefinition {
	componentType := "Spinner"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(componentType, clSpinner.Compile()),
			"size",
			compileClassMap(clSpinnerSize),
		),
		"tone",
		compileClassMap(clSpinnerTone),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible indeterminate progress indicator.",
		map[string]PropertyContract{
			"label": contentProperty("Screen-reader status label."),
			"size":  sizeProperty(controlSizes, "md", "Indicator size."),
			"tone":  toneProperty([]string{"brand", "success", "warning", "danger", "info"}, "brand", "Semantic progress tone."),
		},
		nil,
		deliveryDesign("medium", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "medium", map[string]any{
				"label": "Loading results", "size": "md", "tone": "brand",
			}),
		},
		Spinner,
	)
}

func skeletonDeliveryDefinition() DeliveryDefinition {
	componentType := "Skeleton"
	root := deliveryClassBound(
		deliveryFrame(componentType, clSkeleton.Compile()),
		"size",
		compileClassMap(clSkeletonBlockSize),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Pulsing placeholder that holds the geometry of content that has not arrived yet.",
		map[string]PropertyContract{
			"shape": variantProperty([]string{"block", "text", "circle"}, "block", "Placeholder geometry."),
			"size":  sizeProperty([]string{"small", "medium", "large"}, "medium", "Placeholder scale."),
			"lines": contentProperty("Line count for the text shape."),
		},
		nil,
		deliveryDesign("block", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "block", map[string]any{
				"shape": "block", "size": "medium",
			}),
			deliveryExample(componentType, "text", map[string]any{
				"shape": "text", "lines": 3,
			}),
			deliveryExample(componentType, "circle", map[string]any{
				"shape": "circle", "size": "medium",
			}),
		},
		Skeleton,
	)
}

func deferredSlotDeliveryDefinition() DeliveryDefinition {
	componentType := "DeferredSlot"
	root := deliveryFrame(
		componentType,
		"",
		deliverySlot(
			"Placeholder",
			"placeholder",
			"",
			deliveryInstance("Pending", "Skeleton", "block", ""),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"HTMX seam that swaps a server-rendered fragment over its skeleton placeholder on load.",
		map[string]PropertyContract{},
		[]designcomponent.Slot{{
			Name:         "placeholder",
			Description:  "Skeleton shown until the fragment arrives.",
			AllowedTypes: deliveryAtomicContentTypes(),
			Cardinality:  designcomponent.SlotMany,
		}},
		deliveryDesign("load", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "load", map[string]any{}),
		},
		func(props atoms.DeferredSlotProps) g.Node {
			return DeferredSlot(props, Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3}))
		},
	)
}

func emptyStateDeliveryDefinition() DeliveryDefinition {
	componentType := "EmptyState"
	root := deliveryFrame(
		componentType,
		clEmpty.Merge(clEmptyPad).Compile(),
		deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
		deliveryText("Title", clEmptyTitle.Compile(), "No projects yet", "title"),
		deliveryText("Description", clEmptyDesc.Compile(), "Create a project to begin.", "description"),
		deliverySlot(
			"Actions",
			"actions",
			clCheckRow.Compile(),
			deliveryInstance("PrimaryAction", "Button", "primary", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Empty-data explanation with optional actions.",
		map[string]PropertyContract{
			"title":       contentProperty("Primary empty-state message."),
			"description": contentProperty("Supporting guidance."),
			"compact":     modifierProperty("Whether spacing is compact."),
			"bordered":    modifierProperty("Whether a dashed outline is shown."),
		},
		[]designcomponent.Slot{
			{
				Name:         "iconStart",
				Description:  "Optional leading illustration or icon.",
				AllowedTypes: []string{"Icon"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "actions",
				Description:  "Optional call-to-action controls.",
				AllowedTypes: []string{"Button", "Link"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("default", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"title": "No projects yet", "description": "Create a project to begin.",
					"bordered": true,
				}),
				deliveryExampleSlot(
					"actions",
					deliveryExampleComponent(
						"create-project",
						"Button",
						map[string]any{
							"label":   "Create project",
							"variant": "primary",
							"tone":    "neutral",
							"size":    "md",
						},
					),
				),
			),
			mobileDeliveryExample(componentType, "compact", map[string]any{
				"title": "No results", "description": "Try a different filter.", "compact": true,
			}),
		},
		func(props atoms.EmptyStateProps, slots DeliverySlotChildren) g.Node {
			return emptyStateWithSlots(
				props,
				slots.Nodes("iconStart"),
				slots.Nodes("actions"),
			)
		},
	)
}

func kbdDeliveryDefinition() DeliveryDefinition {
	componentType := "Kbd"
	root := deliveryFrame(
		componentType,
		clCheckRow.Compile(),
		deliveryText("Key", clKbd.Compile(), "⌘"),
		deliveryText("Key", clKbd.Compile(), "K"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Keyboard shortcut key sequence.",
		map[string]PropertyContract{
			"keys": contentProperty("Ordered key labels."),
			"size": sizeProperty([]string{"xs", "sm", "md", "lg"}, "md", "Key cap size."),
		},
		nil,
		deliveryDesign("shortcut", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "shortcut", map[string]any{
				"keys": []string{"⌘", "K"}, "size": "md",
			}),
		},
		Kbd,
	)
}

func linkDeliveryDefinition() DeliveryDefinition {
	componentType := "Link"
	root := deliveryFrame(
		componentType,
		clLink.Compile(),
		deliveryText("Label", "", "View documentation", "label"),
		deliverySlot("TrailingAdornment", "trailingAdornment", clIcon.Compile()),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible navigation link with safe external behavior.",
		map[string]PropertyContract{
			"label":    contentProperty("Visible link label."),
			"href":     contentProperty("Navigation target."),
			"external": stateProperty("Whether the target opens externally."),
			"variant":  variantProperty([]string{"primary", "secondary", "text", "underline"}, "primary", "Visual treatment."),
			"target":   contentProperty("Optional browsing context."),
			"rel":      contentProperty("Link relationship tokens."),
		},
		[]designcomponent.Slot{{
			Name:         "trailingAdornment",
			Description:  "Optional trailing icon or token.",
			AllowedTypes: []string{"Icon"},
			Cardinality:  designcomponent.SlotOne,
		}},
		deliveryDesign("default", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"label": "View documentation", "href": "/docs",
			}),
			deliveryExample(componentType, "external", map[string]any{
				"label": "Open reference", "href": "https://example.com", "external": true,
			}),
		},
		func(props atoms.LinkProps, slots DeliverySlotChildren) g.Node {
			return linkWithSlots(props, slots.Nodes("trailingAdornment"))
		},
	)
}

func tagDeliveryDefinition() DeliveryDefinition {
	componentType := "Tag"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clTagBase.Compile()+" "+clTagIdle.Compile(),
				deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
				deliveryText("Label", "", "Design", "label"),
			),
			"selected",
			map[string]string{
				"false": clTagIdle.Compile(),
				"true":  clTagSelected.Compile(),
			},
		),
		"tone",
		compileClassMap(clTagTone),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Selectable or removable metadata chip.",
		map[string]PropertyContract{
			"label":       contentProperty("Visible tag label."),
			"tone":        toneProperty([]string{"neutral", "brand", "success", "warning", "danger", "info"}, "neutral", "Semantic intent."),
			"removable":   stateProperty("Whether removal is available."),
			"selected":    stateProperty("Whether the tag is selected."),
			"onRemoveURL": contentProperty("Optional progressive-enhancement removal endpoint."),
		},
		[]designcomponent.Slot{{
			Name:         "iconStart",
			Description:  "Optional leading icon slot.",
			AllowedTypes: []string{"Icon"},
			Cardinality:  designcomponent.SlotOne,
		}},
		deliveryDesign("default", "02 Atoms", "Status", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"label": "Design", "tone": "neutral",
			}),
			deliveryExample(componentType, "selected", map[string]any{
				"label": "Design", "tone": "neutral", "selected": true,
			}),
		},
		func(props atoms.TagProps, slots DeliverySlotChildren) g.Node {
			return tagWithSlots(props, slots.Nodes("iconStart"))
		},
	)
}

func compileClassMap[T interface{ Compile() string }](values map[string]T) map[string]string {
	out := make(map[string]string, len(values))
	for key, value := range values {
		out[key] = value.Compile()
	}
	return out
}

func contentProperty(description string) PropertyContract {
	return PropertyContract{Role: designcomponent.PropRoleContent, Description: description}
}

func stateProperty(description string) PropertyContract {
	return PropertyContract{Role: designcomponent.PropRoleState, Description: description}
}

func modifierProperty(description string) PropertyContract {
	return PropertyContract{Role: designcomponent.PropRoleModifier, Description: description}
}

func roleProperty(role designcomponent.PropRole, description string) PropertyContract {
	return PropertyContract{Role: role, Description: description}
}

func enumProperty(values []string, defaultValue, description string) PropertyContract {
	return PropertyContract{
		Role:        designcomponent.PropRoleDefault,
		Description: description,
		Enum:        slices.Clone(values),
		Default:     defaultValue,
	}
}

func variantProperty(values []string, defaultValue, description string) PropertyContract {
	property := enumProperty(values, defaultValue, description)
	property.Role = designcomponent.PropRoleVariant
	return property
}

func toneProperty(values []string, defaultValue, description string) PropertyContract {
	property := enumProperty(values, defaultValue, description)
	property.Role = designcomponent.PropRoleTone
	return property
}

func sizeProperty(values []string, defaultValue, description string) PropertyContract {
	property := enumProperty(values, defaultValue, description)
	property.Role = designcomponent.PropRoleSize
	return property
}

func textSizeClasses() map[string]string {
	out := make(map[string]string, len(clTextSize))
	for key, value := range clTextSize {
		out[key] = tw.New().FontSize(value).Compile()
	}
	return out
}

func textWeightClasses() map[string]string {
	out := make(map[string]string, len(clTextWeight))
	for key, value := range clTextWeight {
		out[key] = tw.New().FontWeight(value).Compile()
	}
	return out
}

func textColorClasses() map[string]string {
	out := make(map[string]string, len(clTextColor))
	for key, value := range clTextColor {
		out[key] = tw.New().TextColor(value).Compile()
	}
	return out
}

func headingLevelClasses() map[string]string {
	out := make(map[string]string, len(clHeadingLevel))
	for level, value := range clHeadingLevel {
		out[fmt.Sprint(level)] = value.Compile()
	}
	return out
}
