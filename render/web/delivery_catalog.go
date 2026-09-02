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

	"github.com/septagon-oss/pk-design/pkg/blueprint"
	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/tw"
)

// deliveryFixtureImageDataURL keeps executable browser examples independent
// from application routes and external hosts. Product consumers still provide
// ordinary image URLs through the same public contracts.
const deliveryFixtureImageDataURL = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Crect width='64' height='64' rx='12' fill='%23e2e8f0'/%3E%3Ccircle cx='32' cy='24' r='12' fill='%2364758b'/%3E%3Cpath d='M12 58c2-13 10-20 20-20s18 7 20 20' fill='%2364758b'/%3E%3C/svg%3E"

const deliverySpinnerArcSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="#000" fill-rule="evenodd" d="M12 2a10 10 0 0 1 10 10h-4a6 6 0 0 0-6-6V2Z"/></svg>`

var (
	buttonVariants = []string{
		"primary", "secondary", "outline", "ghost", "link",
	}
	// Scale and tone lists reference the canonical contract vocabulary so the
	// delivery schema can never drift from the Props contracts it serves.
	controlSizes   = contracts.ControlSizes
	canonicalTones = contracts.CanonicalTones
	feedbackTones  = contracts.FeedbackTones
	gaps           = []string{"0", "1", "2", "3", "4", "5", "6", "8"}
)

// OSSDeliveryCatalog returns the complete renderer-backed OSS web library in
// deterministic atomic-design order. Returned definitions are defensive
// copies; callers can safely adapt them into isolated runtimes.
func OSSDeliveryCatalog() []DeliveryDefinition {
	definitions := []DeliveryDefinition{
		iconDeliveryDefinition(),
		avatarDeliveryDefinition(),
		buttonDeliveryDefinition(),
		badgeDeliveryDefinition(),
		alertDeliveryDefinition(),
		toastDeliveryDefinition(),
		inputDeliveryDefinition(),
		selectDeliveryDefinition(),
		textareaDeliveryDefinition(),
		checkboxDeliveryDefinition(),
		radioDeliveryDefinition(),
		checkboxGroupDeliveryDefinition(),
		radioGroupDeliveryDefinition(),
		sliderDeliveryDefinition(),
		toggleDeliveryDefinition(),
		labelDeliveryDefinition(),
		textDeliveryDefinition(),
		headingDeliveryDefinition(),
		dividerDeliveryDefinition(),
		spinnerDeliveryDefinition(),
		progressDeliveryDefinition(),
		tooltipDeliveryDefinition(),
		skeletonDeliveryDefinition(),
		deferredSlotDeliveryDefinition(),
		emptyStateDeliveryDefinition(),
		kbdDeliveryDefinition(),
		linkDeliveryDefinition(),
		tagDeliveryDefinition(),
		tableDeliveryDefinition(),
		detailListDeliveryDefinition(),
		tableSkeletonDeliveryDefinition(),
		cardDeliveryDefinition(),
		cardSkeletonDeliveryDefinition(),
		breadcrumbDeliveryDefinition(),
		accordionDeliveryDefinition(),
		stepperDeliveryDefinition(),
		sidebarDeliveryDefinition(),
		datePickerDeliveryDefinition(),
		fileUploadDeliveryDefinition(),
		autocompleteDeliveryDefinition(),
		dropdownDeliveryDefinition(),
		actionMenuDeliveryDefinition(),
		drawerDeliveryDefinition(),
		modalDeliveryDefinition(),
		paginationDeliveryDefinition(),
		searchBarDeliveryDefinition(),
		tabsDeliveryDefinition(),
		dataGridDeliveryDefinition(),
		windowedCollectionDeliveryDefinition(),
		dashboardWidgetDeliveryDefinition(),
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

func avatarDeliveryDefinition() DeliveryDefinition {
	componentType := "Avatar"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryClassBound(
				deliveryFrame(
					componentType,
					clAvatarBase.
						Merge(clAvatarSize["md"]).
						Merge(clAvatarShape["circle"]).
						Merge(clAvatarTone["neutral"]).
						Compile(),
					deliveryText("Initials", clAvatarInitials.Compile(), "AL", "initials", "name"),
				),
				"size",
				compileClassMap(clAvatarSize),
			),
			"shape",
			compileClassMap(clAvatarShape),
		),
		"tone",
		compileClassMap(clAvatarTone),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible image, initials, or icon identity mark with governed shape, size, tone, and presence status.",
		map[string]PropertyContract{
			"src":            contentProperty("Optional avatar image URL."),
			"alt":            contentProperty("Accessible image or fallback label."),
			"name":           contentProperty("Name used to derive initials and accessible identity."),
			"initials":       contentProperty("Explicit initials override."),
			"fallbackIcon":   contentProperty("Governed icon name used without an image or initials."),
			"size":           sizeProperty(controlSizes, "md", "Avatar dimensions."),
			"shape":          variantProperty([]string{"circle", "rounded", "square", "pill"}, "circle", "Avatar silhouette."),
			"tone":           toneProperty(canonicalTones, "neutral", "Initials or icon surface tone."),
			"status":         variantProperty([]string{"none", "online", "offline", "busy", "away"}, "none", "Optional presence indicator."),
			"statusPosition": variantProperty([]string{"top-right", "bottom-right", "top-left", "bottom-left"}, "bottom-right", "Presence indicator placement."),
			"statusLabel":    contentProperty("Localized accessible presence label."),
		},
		nil,
		deliveryDesign("initials", "02 Atoms", "Media", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "initials", map[string]any{
				"name": "Ada Lovelace", "size": "lg", "shape": "circle", "tone": "neutral",
			}),
			deliveryExample(componentType, "with-status", map[string]any{
				"name": "Grace Hopper", "status": "online", "statusLabel": "Status: online",
				"statusPosition": "bottom-right", "tone": "brand",
			}),
			deliveryExample(componentType, "with-image", map[string]any{
				"src": deliveryFixtureImageDataURL, "alt": "Profile portrait", "shape": "rounded",
			}),
			deliveryExample(componentType, "fallback-icon", map[string]any{
				"fallbackIcon": "user", "alt": "Account", "size": "sm",
			}),
		},
		Avatar,
	)
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
				canonicalTones,
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
			deliveryExample(componentType, "X Mark / Fg Primary / 20", map[string]any{
				"name": "x-mark", "tone": "neutral", "size": "md", "weight": "outline",
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
			deliveryExample(componentType, "Document / Fg Primary / 20", map[string]any{
				"name": "document", "tone": "neutral", "size": "md", "weight": "outline",
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
			deliveryExample(componentType, "User / Fg Primary / 20", map[string]any{
				"name": "user", "tone": "neutral", "size": "md", "weight": "outline",
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
	root := deliveryFrame(
		componentType,
		clButtonBase.Compile(),
		deliveryHiddenWhen(
			deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
			map[string]any{"loading": true},
		),
		deliveryVisibleWhen(
			deliveryInstance("LoadingIndicator", "Spinner", "small", ""),
			map[string]any{"loading": true},
		),
		deliveryHiddenWhen(
			deliveryText("Label", "", "Continue", "label"),
			map[string]any{"iconOnly": true},
		),
		deliveryHiddenWhen(
			deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
			map[string]any{"loading": true},
		),
	)
	root = deliveryClassBound(root, "variant", compileClassMap(clButtonVariant))
	root = deliveryClassBound(root, "tone", compileClassMap(clButtonTone))
	root = deliveryClassBound(root, "size", compileClassMap(clButtonSize))
	root = deliveryClassBound(root, "fullWidth", map[string]string{
		"false": "",
		"true":  clButtonFull.Compile(),
	})
	root = deliveryClassBound(root, "iconOnly", map[string]string{
		"false": "",
		"true":  clButtonIconOnly.Compile(),
	})
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible action control rendered by the OSS web runtime.",
		map[string]PropertyContract{
			"label":     contentProperty("Visible button label."),
			"href":      contentProperty("Optional navigation target; renders link semantics when set."),
			"variant":   variantProperty(buttonVariants, "primary", "Visual treatment."),
			"tone":      toneProperty(canonicalTones, "neutral", "Semantic intent."),
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
			return ButtonWithSlots(props, ButtonSlots{
				IconStart: slots.Nodes("iconStart"),
				IconEnd:   slots.Nodes("iconEnd"),
			})
		},
	)
}

func badgeDeliveryDefinition() DeliveryDefinition {
	componentType := "Badge"
	statusDot := deliveryClassBound(
		deliveryFrame("StatusDot", clBadgeDot.Compile()),
		"tone",
		compileClassMap(clBadgeDotTone),
	)
	toneClasses := compileClassMap(clBadgeTone)
	// Neutral preserves the selected visual variant, matching the renderer.
	toneClasses["neutral"] = ""
	root := deliveryFrame(
		componentType,
		clBadgeBase.Compile(),
		deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
		deliveryVisibleWhen(statusDot, map[string]any{"dot": true}),
		deliveryText("Label", "", "Active", "label"),
		deliveryOptionalText(deliveryText("Count", clBadgeCount.Compile(), "", "count")),
		deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
		deliveryVisibleWhen(
			deliveryText("Remove", clBadgeRemove.Compile(), "×"),
			map[string]any{"removable": true},
		),
	)
	root = deliveryClassBound(root, "variant", compileClassMap(clBadgeVariant))
	root = deliveryClassBound(root, "tone", toneClasses)
	root = deliveryClassBound(root, "size", compileClassMap(clBadgeSize))
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Compact semantic status label.",
		map[string]PropertyContract{
			"label":   contentProperty("Visible badge label."),
			"variant": variantProperty([]string{"primary", "secondary", "outline"}, "primary", "Visual treatment."),
			"tone": toneProperty(
				canonicalTones,
				"neutral",
				"Semantic intent independent of the visual treatment.",
			),
			"size":        sizeProperty(controlSizes, "md", "Badge scale."),
			"dot":         modifierProperty("Shows a leading status dot."),
			"count":       roleProperty(designcomponent.PropRoleContent, "Positive numeric indicator, visually capped to 99+."),
			"removable":   modifierProperty("Shows an accessible remove affordance."),
			"removeLabel": contentProperty("Localized remove-button label."),
			"live":        stateProperty("Announces dynamic badge content as a polite status."),
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
			deliveryExample(componentType, "count", map[string]any{
				"label": "Messages", "count": 125, "variant": "secondary", "tone": "neutral", "size": "sm",
			}),
			deliveryExample(componentType, "removable", map[string]any{
				"label": "Filter", "removable": true, "removeLabel": "Remove Filter", "tone": "brand", "size": "md",
			}),
		},
		func(props atoms.BadgeProps, slots DeliverySlotChildren) g.Node {
			return BadgeWithSlots(props, BadgeSlots{
				IconStart: slots.Nodes("iconStart"),
				IconEnd:   slots.Nodes("iconEnd"),
			})
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
			"tone":        toneProperty(feedbackTones, "info", "Severity."),
			"dismissible": stateProperty("Whether a dismiss affordance is exposed."),
			"bordered":    modifierProperty("Whether the message has a leading accent border."),
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

func toastDeliveryDefinition() DeliveryDefinition {
	componentType := "Toast"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clToastBase.Compile(),
			deliverySlot("LeadingIcon", "iconStart", clToastIcon.Compile()),
			deliveryFrame(
				"Body",
				clToastBody.Compile(),
				deliveryText("Title", clToastTitle.Compile(), "Saved", "title"),
				deliveryText("Message", clToastMessage.Compile(), "Settings updated", "message"),
			),
			deliveryText("Close", clToastClose.Compile(), "×", "closeLabel"),
		),
		"tone",
		compileClassMap(clToastTone),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Transient accessible feedback with governed urgency, dismissal, and placement metadata.",
		map[string]PropertyContract{
			"message": contentProperty("Required notification message."),
			"title":   contentProperty("Optional notification title."),
			"tone":    toneProperty(feedbackTones, "info", "Notification severity."),
			"duration": {
				Role: designcomponent.PropRoleState, Description: "Auto-dismiss delay in milliseconds.", Default: "5000",
			},
			"persistent": stateProperty("Whether the notification remains until explicitly dismissed."),
			"position":   enumProperty([]string{"top-right", "top-left", "bottom-right", "bottom-left"}, "top-right", "Viewport placement metadata."),
			"closable": {
				Role: designcomponent.PropRoleModifier, Description: "Whether an accessible dismiss affordance is shown.", Default: "true",
			},
			"closeLabel": contentProperty("Localized label for the dismiss affordance."),
		},
		[]designcomponent.Slot{{
			Name: "iconStart", Description: "Optional leading status glyph.",
			AllowedTypes: []string{"Icon"}, Cardinality: designcomponent.SlotOne,
		}},
		deliveryDesign("success", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "success", map[string]any{
				"title": "Saved", "message": "Settings updated", "tone": "success", "duration": 5000,
			}),
			deliveryExample(componentType, "persistent-warning", map[string]any{
				"message": "Connection interrupted", "tone": "warning", "persistent": true,
			}),
		},
		func(props atoms.ToastProps, slots DeliverySlotChildren) g.Node {
			return toastWithSlots(props, slots.Nodes("iconStart"))
		},
	)
}

func inputDeliveryDefinition() DeliveryDefinition {
	componentType := "Input"
	value := deliveryOptionalText(deliveryText("Value", "", "", "value", "placeholder"))
	value = deliveryClassBound(value, "value", map[string]string{
		"":  tw.New().TextColor(tw.FgPlaceholder).Compile(),
		"*": tw.New().TextColor(tw.FgPrimary).Compile(),
	})
	value.Props["mask_when"] = map[string]any{"type": "password"}
	control := deliveryFrame(
		"Control",
		clInput.Compile()+" "+tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).Compile(),
		deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
		value,
		deliverySlot("TrailingIcon", "iconEnd", clIcon.Compile()),
	)
	control = deliveryClassBound(control, "size", compileClassMap(clInputSize))
	control = deliveryClassBound(control, "tone", compileClassMap(clInputTone))
	control = deliveryClassBound(control, "error", map[string]string{
		"": "", "*": clInputError.Compile(),
	})
	control = deliveryClassBound(control, "invalid", map[string]string{
		"false": "", "true": clInputError.Compile(),
	})
	control = deliveryClassBound(control, "readOnly", map[string]string{
		"false": "", "true": clInputReadOnly.Compile(),
	})
	help := deliveryOptionalText(
		deliveryText("Help", clHelp.Compile(), "", "error", "helpText"),
	)
	help = deliveryClassBound(help, "error", map[string]string{
		"": "", "*": clFieldErr.Compile(),
	})
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFieldWrap.Merge(clFieldWrapFull).Compile(),
			deliveryOptionalText(deliveryText("Label", clLabel.Compile(), "", "label")),
			control,
			help,
		),
		"field",
		"field",
	)
	definition := newSlottedDeliveryDefinition(
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
			"error":       contentProperty("Validation message; non-empty selects error styling."),
			"invalid":     stateProperty("Whether invalid styling and aria-invalid are active without an inline message."),
			"required":    stateProperty("Whether a value is required."),
			"readOnly":    stateProperty("Whether editing is prevented."),
			"autoFocus":   stateProperty("Whether focus is requested on mount."),
			"size":        sizeProperty([]string{"sm", "md", "lg"}, "md", "Control size."),
			"tone":        toneProperty([]string{"neutral", "success", "warning", "danger"}, "neutral", "Semantic validation or review state."),
			"fullWidth":   modifierProperty("Whether the field fills its available width."),
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
			deliveryExample(componentType, "validated", map[string]any{
				"name": "email", "type": "email", "label": "Email address",
				"value": "name@example.com", "tone": "success",
			}),
			deliveryExample(componentType, "locked-value", map[string]any{
				"name": "workspace", "label": "Workspace", "value": "platform",
				"readOnly": true,
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
	baseRender := definition.Render
	definition.Render = func(props map[string]any, slots DeliverySlotChildren) (g.Node, error) {
		if rawType, exists := props["type"]; exists {
			if typed, ok := rawType.(string); ok {
				canonical, valid := canonicalInputType(typed)
				if !valid {
					return nil, fmt.Errorf("Input props: unsupported type %q", typed)
				}
				if canonical != typed {
					props = cloneDeliveryValueMap(props)
					props["type"] = canonical
				}
			}
		}
		return baseRender(props, slots)
	}
	return definition
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
		uicomponent.TierAtom,
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
			"fullWidth":   modifierProperty("Whether the field fills its available width."),
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
					{"label": "Draft", "value": "draft", "group": "Current"},
					{"label": "Published", "value": "published", "group": "Current"},
					{"label": "Archived", "value": "archived", "group": "History"},
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
			"rows": {
				Role: designcomponent.PropRoleSize, Description: "Visible row count when auto-resize is disabled.", Default: "4",
			},
			"minRows": {
				Role: designcomponent.PropRoleSize, Description: "Minimum auto-resize row count.", Default: "2",
			},
			"maxRows": {
				Role: designcomponent.PropRoleSize, Description: "Maximum auto-resize row count.", Default: "10",
			},
			"minLength":  roleProperty(designcomponent.PropRoleSize, "Minimum accepted character count."),
			"maxLength":  roleProperty(designcomponent.PropRoleSize, "Maximum accepted character count."),
			"showCount":  modifierProperty("Whether a live character count is shown."),
			"autoResize": modifierProperty("Whether the control grows between minRows and maxRows."),
			"fullWidth":  modifierProperty("Whether the field fills its container."),
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
				"maxLength":   500, "fullWidth": true, "autoResize": true,
				"minRows": 3, "maxRows": 15,
			}),
			deliveryExample(componentType, "with-error", map[string]any{
				"name": "notes", "label": "Notes", "value": "too short",
				"errorMessage": "Must be at least 20 characters.",
				"helperText":   "Explain enough for another person to act.",
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
			clCheckboxRoot.Compile(),
			deliveryClassBound(
				deliveryFrame("Indicator", clCheckboxIndicator.Compile()),
				"checked",
				map[string]string{
					"false": clCheckboxIndicatorIdle.Compile(),
					"true":  clCheckboxIndicatorActive.Compile(),
				},
			),
			deliveryText("Label", clCheckboxLabel.Compile(), "Include archived items", "label"),
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
			"label":         contentProperty("Optional visible label; name supplies the accessible label when omitted."),
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

func radioDeliveryDefinition() DeliveryDefinition {
	componentType := "Radio"
	root := deliveryClassBound(
		deliverySemantics(
			deliveryFrame(
				componentType,
				clRadioRoot.Compile(),
				deliveryFrame(
					"NativeControl",
					clRadioInput.Compile(),
					deliveryVisibleWhen(
						deliveryFrame("SelectedDot", clRadioDot.Compile()),
						map[string]any{"checked": true},
					),
				),
				deliveryText("Label", clRadioLabel.Compile(), "Standard plan", "label", "value"),
			),
			"radio",
			"single-selection field",
		),
		"disabled",
		map[string]string{
			"false": clRadioRoot.Compile(),
			"true":  clRadioRoot.Merge(clRadioRootDisabled).Compile(),
		},
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Native single-selection control with governed labels, states, focus, and semantic accent.",
		map[string]PropertyContract{
			"name":     contentProperty("Required field name shared by one radio group."),
			"label":    contentProperty("Optional visible choice label."),
			"helpText": contentProperty("Optional supporting guidance for this choice."),
			"value":    contentProperty("Required submitted choice value."),
			"checked":  stateProperty("Whether this choice is selected."),
			"required": stateProperty("Whether the group requires a selection."),
			"disabled": stateProperty("Whether this choice is unavailable."),
		},
		nil,
		deliveryDesign("basic", "02 Atoms", "Form Controls", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"name": "plan", "value": "standard", "label": "Standard plan",
			}),
			deliveryExample(componentType, "checked", map[string]any{
				"name": "plan", "value": "pro", "label": "Pro plan", "checked": true,
			}),
			deliveryExample(componentType, "disabled", map[string]any{
				"name": "plan", "value": "enterprise", "label": "Enterprise", "disabled": true,
			}),
		},
		Radio,
	)
	return withDeliveryTags(definition, "form", "radio", "choice", "single-selection")
}

func sliderDeliveryDefinition() DeliveryDefinition {
	componentType := "Slider"
	input := deliveryClassBound(
		deliveryFrame("NativeRange", clSliderInput.Merge(clSliderTone["brand"]).Compile()),
		"tone",
		compileClassMap(clSliderTone),
	)
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clSliderRoot.Compile(),
			deliveryText("Label", clSliderLabel.Compile(), "Opacity", "label", "name"),
			deliveryFrame(
				"ControlRow",
				clSliderRow.Compile(),
				input,
				deliveryVisibleWhen(
					deliveryText("Value", clSliderValue.Compile(), "75", "value"),
					map[string]any{"showValue": true},
				),
			),
		),
		"group",
		"numeric range field",
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Native numeric range control with governed bounds, semantic tone, accessible value text, and live value readout.",
		map[string]PropertyContract{
			"name":          contentProperty("Required form field name and fallback accessible label."),
			"label":         contentProperty("Optional visible field label."),
			"min":           roleProperty(designcomponent.PropRoleDefault, "Required minimum allowed value."),
			"max":           roleProperty(designcomponent.PropRoleDefault, "Required maximum allowed value."),
			"step":          roleProperty(designcomponent.PropRoleDefault, "Step increment; non-positive values fall back to one."),
			"value":         contentProperty("Current numeric value."),
			"showValue":     modifierProperty("Whether to render a live numeric value mirror."),
			"tone":          toneProperty(canonicalTones, "brand", "Semantic accent for the native range control."),
			"ariaValueText": contentProperty("Optional human-readable value including units."),
			"disabled":      stateProperty("Whether the range control is unavailable."),
		},
		nil,
		deliveryDesign("with-value", "02 Atoms", "Form Controls", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "with-value", map[string]any{
				"name": "opacity", "label": "Opacity", "min": 0, "max": 100,
				"step": 1, "value": 75, "showValue": true, "tone": "brand",
			}),
			deliveryExample(componentType, "custom-range", map[string]any{
				"name": "cpu_limit", "label": "CPU limit", "min": 10, "max": 500,
				"step": 5, "value": 100, "showValue": true, "tone": "warning",
				"ariaValueText": "100 millicores",
			}),
			deliveryExample(componentType, "disabled", map[string]any{
				"name": "volume", "label": "Volume", "min": 0, "max": 100,
				"step": 1, "value": 40, "disabled": true, "tone": "brand",
			}),
		},
		Slider,
	)
	return withDeliveryTags(definition, "range", "slider", "input", "form", "numeric", "selection")
}

func toggleDeliveryDefinition() DeliveryDefinition {
	componentType := "Toggle"
	knob := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame("Knob", clToggleKnob.Merge(clToggleKnobSize["md"]).Compile()),
			"size",
			compileClassMap(clToggleKnobSize),
		),
		"checked",
		map[string]string{
			"false": clToggleKnobUnchecked.Compile(),
			"true":  clToggleKnobChecked["md"].Compile(),
		},
	)
	track := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame("Track", clToggleTrack.Merge(clToggleTrackSize["md"]).Compile(), knob),
			"size",
			compileClassMap(clToggleTrackSize),
		),
		"checked",
		map[string]string{
			"false": clToggleTrackState["unchecked"].Compile(),
			"true":  clToggleTrackState["checked"].Compile(),
		},
	)
	root := deliveryClassBound(
		deliverySemantics(
			deliveryFrame(
				componentType,
				clToggleRoot.Compile(),
				track,
				deliveryVisibleWhen(
					deliveryText("Label", clToggleLabel.Compile(), "Enable notifications", "label", "name"),
					map[string]any{"hideLabel": false},
				),
			),
			"switch",
			"binary on/off field",
		),
		"disabled",
		map[string]string{
			"false": clToggleRoot.Compile(),
			"true":  clToggleRoot.Merge(clToggleRootDisabled).Compile(),
		},
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible binary switch with one interactive surface, hidden form value, governed sizes, and synchronized track and knob states.",
		map[string]PropertyContract{
			"name":      contentProperty("Required form input name."),
			"label":     contentProperty("Optional visible switch label."),
			"ariaLabel": contentProperty("Optional accessible-name override."),
			"hideLabel": modifierProperty("Whether to omit the visible label while retaining its accessible name."),
			"checked":   stateProperty("Whether the switch starts on."),
			"size":      sizeProperty([]string{"sm", "md", "lg"}, "md", "Track and knob dimensions."),
			"disabled":  stateProperty("Whether the switch is unavailable."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Form Controls", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"name": "notifications", "label": "Enable notifications", "size": "md",
			}),
			deliveryExample(componentType, "checked", map[string]any{
				"name": "dark_mode", "label": "Dark mode", "checked": true, "size": "lg",
			}),
			deliveryExample(componentType, "hidden-label", map[string]any{
				"name": "marketing", "label": "Marketing cookies", "hideLabel": true,
				"checked": true, "size": "sm",
			}),
			deliveryExample(componentType, "disabled", map[string]any{
				"name": "feature", "label": "Feature flag", "disabled": true, "size": "md",
			}),
		},
		Toggle,
	)
	return withDeliveryTags(definition, "switch", "toggle", "boolean", "form", "input", "on-off")
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
				deliveryClassBound(
					deliveryClassBound(
						deliveryText("Text", "", "A clear supporting sentence.", "content"),
						"size",
						textSizeClasses(),
					),
					"align",
					textAlignClasses(),
				),
				"weight",
				textWeightClasses(),
			),
			"color",
			textColorClasses(),
		),
		"transform",
		compileClassMap(clTextTransform),
	)
	return withDeliveryTags(newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Semantic body-copy primitive.",
		map[string]PropertyContract{
			"content": contentProperty("Visible copy."),
			"element": variantProperty(
				[]string{"p", "span", "div", "strong", "em", "small", "mark", "del", "ins", "sub", "sup", "blockquote", "code", "pre", "kbd", "samp", "var"},
				"p", "Allow-listed non-heading semantic element.",
			),
			"size":      sizeProperty([]string{"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl"}, "base", "Text scale."),
			"align":     variantProperty([]string{"left", "center", "right", "justify"}, "left", "Text alignment."),
			"weight":    variantProperty([]string{"thin", "extralight", "light", "normal", "medium", "semibold", "bold", "extrabold", "black"}, "normal", "Font weight."),
			"color":     toneProperty([]string{"primary", "secondary", "tertiary", "muted", "brand", "success", "warning", "danger", "info"}, "primary", "Semantic foreground."),
			"transform": variantProperty([]string{"none", "uppercase", "lowercase", "capitalize"}, "none", "Text casing transform."),
			"truncate":  modifierProperty("Whether single-line overflow is truncated."),
			"nowrap":    modifierProperty("Whether wrapping is disabled."),
			"italic":    modifierProperty("Whether copy is italicized."),
			"underline": modifierProperty("Whether copy is underlined."),
			"lines":     modifierProperty("Optional line-clamp count from one through six."),
		},
		nil,
		deliveryDesign("body", "02 Atoms", "Typography", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "body", map[string]any{
				"content": "A clear supporting sentence.", "element": "p", "size": "base",
				"weight": "normal", "color": "primary",
			}),
			deliveryExample(componentType, "muted", map[string]any{
				"content": "Supporting metadata", "size": "sm", "color": "muted",
			}),
			deliveryExample(componentType, "muted-body", map[string]any{
				"content": "Muted supporting copy", "size": "base", "color": "muted",
			}),
			deliveryExample(componentType, "emphasis", map[string]any{
				"content": "Important", "element": "strong", "weight": "bold", "italic": true,
			}),
			deliveryExample(componentType, "clamped", map[string]any{
				"content": "Long supporting copy", "lines": 3, "color": "secondary",
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
	horizontal := tw.New().Width(tw.S80).Height(tw.SPX).Bg(tw.BorderPrimary).Merge(clDividerH).Compile()
	vertical := tw.New().Width(tw.SPX).Height(tw.S8).Bg(tw.BorderPrimary).Merge(clDividerV).Compile()
	labelled := tw.New().Width(tw.S80).Height(tw.SAuto).Bg(tw.ColorTransparent).Merge(clDividerText).Compile()
	line := tw.New().Height(tw.SPX).Bg(tw.BorderPrimary).Merge(clDividerTextLine).Compile()
	visibleWithText := map[string]any{"text": "*"}
	root := deliveryFrame(
		componentType,
		horizontal,
		deliveryVisibleWhen(deliveryFrame("LineStart", line), visibleWithText),
		deliveryVisibleWhen(
			deliveryOptionalText(deliveryText("Label", clDividerTextLabel.Compile(), "", "text")),
			visibleWithText,
		),
		deliveryVisibleWhen(deliveryFrame("LineEnd", line), visibleWithText),
	)
	root = deliveryClassBound(root, "orientation", map[string]string{
		"horizontal": horizontal,
		"vertical":   vertical,
	})
	root = deliveryClassBound(root, "text", map[string]string{
		"": "", "*": labelled,
	})
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
	track := deliveryClassBound(
		deliveryFrame(
			"Track",
			tw.New().Position(tw.PositionAbsolute).Inset(tw.S0).Rounded(tw.RadiusFull).
				Border(tw.Border2).BorderColor(tw.BorderSecondary).Compile(),
		),
		"size",
		compileClassMap(clSpinnerSize),
	)
	arc := blueprint.Node{
		Kind:     blueprint.NodeSVG,
		Name:     "Arc",
		Classes:  tw.New().Position(tw.PositionAbsolute).Inset(tw.S0).Compile(),
		AssetRef: "asset:spinner-arc",
	}
	arc = deliveryClassBound(arc, "size", compileClassMap(clSpinnerSize))
	arc = deliveryClassBound(arc, "tone", map[string]string{
		"neutral": tw.New().TextColor(tw.FgSecondary).Compile(),
		"brand":   tw.New().TextColor(tw.FgBrand).Compile(),
		"success": tw.New().TextColor(tw.FgSuccess).Compile(),
		"warning": tw.New().TextColor(tw.FgWarning).Compile(),
		"danger":  tw.New().TextColor(tw.FgDanger).Compile(),
		"info":    tw.New().TextColor(tw.FgInfo).Compile(),
	})
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			tw.New().Position(tw.PositionRelative).Display(tw.DisplayInlineBlock).Compile(),
			track,
			arc,
		),
		"size",
		compileClassMap(clSpinnerSize),
	)
	design := deliveryDesign("medium", "02 Atoms", "Feedback", root)
	design.Assets = []blueprint.Asset{{
		Name: "spinner-arc", Kind: blueprint.AssetSVG, Source: deliverySpinnerArcSVG,
		Description: "Token-tinted quarter arc layered over the neutral loading track.",
	}}
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible indeterminate progress indicator.",
		map[string]PropertyContract{
			"label": contentProperty("Screen-reader status label."),
			"size":  sizeProperty(controlSizes, "md", "Indicator size."),
			"tone":  toneProperty(canonicalTones, "brand", "Semantic progress tone."),
		},
		nil,
		design,
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "medium", map[string]any{
				"label": "Loading results", "size": "md", "tone": "brand",
			}),
			deliveryExample(componentType, "small", map[string]any{
				"label": "Loading", "size": "sm", "tone": "brand",
			}),
		},
		Spinner,
	)
}

func progressDeliveryDefinition() DeliveryDefinition {
	componentType := "Progress"
	toneClasses := compileClassMap(clProgressTone)
	fill := deliveryProgressBound(
		deliveryFrame("Fill", clProgressFill.Merge(clProgressTone["brand"]).Compile()),
		false,
	)
	fill = deliveryClassBound(fill, "tone", toneClasses)
	fill = deliveryClassBound(fill, "indeterminate", map[string]string{
		"false": "", "true": clProgressIndeterminate.Compile(),
	})
	track := deliveryClassBound(
		deliveryFrame(
			"Track",
			clProgressTrack.Merge(clProgressTrackSize["md"]).Compile(),
			fill,
		),
		"size",
		compileClassMap(clProgressTrackSize),
	)
	label := deliveryOptionalText(deliveryText("Label", clProgressLabel.Compile(), "", "label"))
	value := deliveryProgressBound(
		deliveryText("Value", clProgressValue.Compile(), "0%"),
		true,
	)
	value = deliveryVisibleWhen(value, map[string]any{"showText": true})
	value = deliveryHiddenWhen(value, map[string]any{"indeterminate": true})
	header := deliveryVisibleWhenAny(
		deliveryFrame("Header", clProgressHeader.Compile(), label, value),
		map[string]any{"label": "*", "showText": true},
	)
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clProgressRoot.Compile(),
			header,
			track,
		),
		"progressbar",
		"determinate or indeterminate task completion status",
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible progress status with governed bounds, labels, semantic tones, sizes, percentage text, and indeterminate state.",
		map[string]PropertyContract{
			"value":         roleProperty(designcomponent.PropRoleDefault, "Current value, normalized into the zero-to-maximum range."),
			"max":           roleProperty(designcomponent.PropRoleDefault, "Positive maximum value; non-positive values fall back to 100."),
			"label":         contentProperty("Optional visible progress label."),
			"ariaLabel":     contentProperty("Accessible name used when no visible label is rendered."),
			"showText":      modifierProperty("Whether to show the normalized integer percentage."),
			"tone":          toneProperty(canonicalTones, "brand", "Semantic fill tone."),
			"size":          sizeProperty([]string{"sm", "md", "lg"}, "md", "Progress track height."),
			"indeterminate": stateProperty("Whether progress has no currently knowable value."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"value": 75, "max": 100, "ariaLabel": "Upload progress",
				"showText": true, "tone": "brand", "size": "md",
			}),
			deliveryExample(componentType, "with-label", map[string]any{
				"value": 45, "max": 100, "label": "Upload progress",
				"showText": true, "tone": "warning", "size": "lg",
			}),
			deliveryExample(componentType, "complete", map[string]any{
				"value": 100, "max": 100, "ariaLabel": "Import completion",
				"showText": true, "tone": "success", "size": "sm",
			}),
			deliveryExample(componentType, "indeterminate", map[string]any{
				"value": 0, "label": "Preparing export", "indeterminate": true,
				"tone": "info", "size": "md",
			}),
		},
		Progress,
	)
	return withDeliveryTags(definition, "progress", "status", "loading", "completion", "percentage", "feedback")
}

func tooltipDeliveryDefinition() DeliveryDefinition {
	componentType := "Tooltip"
	root := deliveryFrame(
		componentType,
		clTooltipContainer.Compile(),
		deliverySlot(
			"Trigger",
			"trigger",
			clTooltipTrigger.Compile(),
			deliveryInstance("TriggerButton", "Button", "primary", ""),
		),
		deliveryClassBound(
			deliveryText(
				"Popup",
				clTooltipPopup.Compile()+" "+clTooltipPosition["top"].Compile(),
				"Helpful context",
				"content",
			),
			"position",
			compileClassMap(clTooltipPosition),
		),
	)
	definition := newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible contextual help shown from one authored trigger on hover or focus.",
		map[string]PropertyContract{
			"content":  contentProperty("Required contextual help text."),
			"position": variantProperty([]string{"top", "bottom", "left", "right"}, "top", "Popup placement relative to the trigger."),
			"delay": {
				Role:        designcomponent.PropRoleDefault,
				Description: "Show delay in milliseconds.",
				Default:     "200",
			},
		},
		[]designcomponent.Slot{{
			Name:         "trigger",
			Required:     true,
			Description:  "The single interactive or informational element that owns the tooltip.",
			AllowedTypes: deliveryAtomicContentTypes(),
			Cardinality:  designcomponent.SlotOne,
		}},
		deliveryDesign("default", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"content": "Helpful context", "position": "top", "delay": 200,
				}),
				deliveryExampleSlot(
					"trigger",
					deliveryExampleComponent(
						"tooltip-trigger",
						"Button",
						map[string]any{
							"label": "More information", "variant": "secondary", "size": "md",
						},
					),
				),
			),
		},
		func(props atoms.TooltipProps, slots DeliverySlotChildren) g.Node {
			return Tooltip(props, slots.Nodes("trigger")...)
		},
	)
	return withDeliveryTags(definition, "tooltip", "contextual-help", "accessibility")
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
			"size":  sizeProperty([]string{"sm", "md", "lg"}, "md", "Placeholder scale."),
			"lines": contentProperty("Line count for the text shape."),
		},
		nil,
		deliveryDesign("block", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "block", map[string]any{
				"shape": "block", "size": "md",
			}),
			deliveryExample(componentType, "text", map[string]any{
				"shape": "text", "lines": 3,
			}),
			deliveryExample(componentType, "circle", map[string]any{
				"shape": "circle", "size": "md",
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
			// A placeholder may be several skeleton pieces standing in for the
			// fragment's shape, so each occupancy carries its own identity.
			Attrs: deliveryRepeatingSlotAttrs(),
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
	toneClasses := compileClassMap(clTagTone)
	// Neutral preserves the selected state; semantic tones intentionally win.
	toneClasses["neutral"] = ""
	root := deliveryFrame(
		componentType,
		clTagBase.Compile()+" "+clTagIdle.Compile(),
		deliverySlot("LeadingIcon", "iconStart", clIcon.Compile()),
		deliveryText("Label", "", "Design", "label"),
		deliveryVisibleWhen(
			deliveryText("Remove", clLink.Compile(), "×"),
			map[string]any{"removable": true},
		),
	)
	root = deliveryClassBound(root, "selected", map[string]string{
		"false": clTagIdle.Compile(),
		"true":  clTagSelected.Compile(),
	})
	root = deliveryClassBound(root, "tone", toneClasses)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Selectable or removable metadata chip.",
		map[string]PropertyContract{
			"label":       contentProperty("Visible tag label."),
			"tone":        toneProperty(canonicalTones, "neutral", "Semantic intent."),
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

func textAlignClasses() map[string]string {
	out := make(map[string]string, len(clTextAlign))
	for key, value := range clTextAlign {
		out[key] = tw.New().TextAlign(value).Compile()
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
