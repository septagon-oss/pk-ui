// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_composites.go defines the molecule, organism, template, and page
// half of the OSS delivery catalog. The page is a real executable composition
// of catalog components and a native instance graph—not a flattened drawing.

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/layouts"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	"github.com/septagon-oss/pk-ui/contracts/organisms"
	"github.com/septagon-oss/tw"
)

func checkboxGroupDeliveryDefinition() DeliveryDefinition {
	componentType := "CheckboxGroup"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clChoiceGroupRoot.Compile(),
			deliveryText("Legend", clChoiceGroupLegend.Compile(), "Notify me about", "label"),
			deliveryText("Description", clChoiceGroupDescription.Compile(), "Choose every update you want.", "description"),
			deliveryFrame(
				"Options",
				clChoiceGroupOptions.Merge(clChoiceGroupVertical).Compile(),
				deliveryInstance("Option", "Checkbox", "unchecked", ""),
			),
			deliveryText("Error", clChoiceGroupError.Compile(), "Choose at least one option.", "error"),
		),
		"group",
		"multiple-choice field",
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Native multiple-choice fieldset with shared guidance, validation, orientation, and canonical checkbox options.",
		map[string]PropertyContract{
			"name":        contentProperty("Required submitted field name shared by every option."),
			"label":       contentProperty("Visible group legend."),
			"description": contentProperty("Optional supporting guidance for the group."),
			"ariaLabel":   contentProperty("Accessible name used when no visible legend is present."),
			"options":     contentProperty("Required ordered checkbox choices."),
			"selected":    stateProperty("Initially selected option values."),
			"required":    stateProperty("Whether server validation requires at least one selection."),
			"orientation": variantProperty([]string{"vertical", "horizontal"}, "vertical", "Choice layout."),
			"error":       contentProperty("Validation message associated with the fieldset."),
		},
		nil,
		deliveryDesign("notification-preferences", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "notification-preferences", map[string]any{
				"name": "notifications", "label": "Notify me about",
				"description": "Choose every update you want to receive.",
				"selected":    []string{"security"},
				"options": []map[string]any{
					{"label": "Product updates", "value": "product"},
					{"label": "Security alerts", "value": "security", "description": "Critical account notices."},
				},
			}),
			deliveryExample(componentType, "horizontal", map[string]any{
				"name": "channels", "label": "Channels", "orientation": "horizontal",
				"options": []map[string]any{
					{"label": "Email", "value": "email"},
					{"label": "SMS", "value": "sms"},
				},
			}),
			deliveryExample(componentType, "with-error", map[string]any{
				"name": "agreements", "label": "Agreements", "required": true,
				"error":   "Choose at least one agreement.",
				"options": []map[string]any{{"label": "Product terms", "value": "terms"}},
			}),
		},
		CheckboxGroup,
	)
	return withDeliveryTags(definition, "form", "checkbox", "choice", "multiple-selection", "fieldset")
}

func radioGroupDeliveryDefinition() DeliveryDefinition {
	componentType := "RadioGroup"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clChoiceGroupRoot.Compile(),
			deliveryText("Legend", clChoiceGroupLegend.Compile(), "Choose a plan", "label"),
			deliveryText("Description", clChoiceGroupDescription.Compile(), "Select one plan.", "description"),
			deliveryFrame(
				"Options",
				clChoiceGroupOptions.Merge(clChoiceGroupVertical).Compile(),
				deliveryInstance("Option", "Radio", "basic", ""),
			),
			deliveryText("Error", clChoiceGroupError.Compile(), "Choose one plan.", "error"),
		),
		"radiogroup",
		"single-choice field",
	)
	definition := newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Native single-choice fieldset with shared guidance, validation, orientation, and canonical radio options.",
		map[string]PropertyContract{
			"name":        contentProperty("Required submitted field name shared by every option."),
			"label":       contentProperty("Visible group legend."),
			"description": contentProperty("Optional supporting guidance for the group."),
			"ariaLabel":   contentProperty("Accessible name used when no visible legend is present."),
			"options":     contentProperty("Required ordered radio choices."),
			"value":       stateProperty("Initially selected option value."),
			"required":    stateProperty("Whether the browser requires one selection."),
			"orientation": variantProperty([]string{"vertical", "horizontal"}, "vertical", "Choice layout."),
			"error":       contentProperty("Validation message associated with the fieldset."),
		},
		nil,
		deliveryDesign("plan", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "plan", map[string]any{
				"name": "plan", "label": "Choose a plan", "description": "Select one plan.",
				"value": "growth", "required": true,
				"options": []map[string]any{
					{"label": "Starter", "value": "starter"},
					{"label": "Growth", "value": "growth", "description": "For growing teams."},
					{"label": "Enterprise", "value": "enterprise", "disabled": true},
				},
			}),
			deliveryExample(componentType, "horizontal", map[string]any{
				"name": "billing", "label": "Billing cadence", "orientation": "horizontal",
				"options": []map[string]any{
					{"label": "Monthly", "value": "monthly"},
					{"label": "Annual", "value": "annual"},
				},
			}),
			deliveryExample(componentType, "with-error", map[string]any{
				"name": "approval", "label": "Approval", "required": true,
				"error":   "Choose an approval state.",
				"options": []map[string]any{{"label": "Approve", "value": "approve"}},
			}),
		},
		RadioGroup,
	)
	return withDeliveryTags(definition, "form", "radio", "choice", "single-selection", "fieldset")
}

func tableDeliveryDefinition() DeliveryDefinition {
	componentType := "Table"
	root := deliveryFrame(
		componentType,
		clTableWrap.Compile(),
		deliveryFrame(
			"Table",
			clTable.Compile(),
			deliveryFrame(
				"Header",
				clTableHead.Compile(),
				deliveryText("Column", clTableTh.Compile(), "Name"),
				deliveryText("Column", clTableTh.Compile(), "Status"),
				deliveryText("Column", clTableTh.Compile(), "Updated"),
			),
			deliveryFrame(
				"Row",
				clTableRow.Compile(),
				deliveryText("PrimaryCell", clTableTd.Merge(clTableTdStrong).Compile(), "Alpha workspace"),
				deliveryText("Cell", clTableTd.Compile(), "Active"),
				deliveryText("Cell", clTableTd.Compile(), "Today"),
			),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible tabular data presentation with sorting and selection hooks.",
		map[string]PropertyContract{
			"columns":    contentProperty("Ordered column definitions."),
			"rows":       contentProperty("Ordered data rows."),
			"sortable":   stateProperty("Whether sortable columns expose controls."),
			"selectable": stateProperty("Whether row selection is available."),
			"striped":    modifierProperty("Whether alternate rows are tinted."),
			"compact":    modifierProperty("Whether row density is compact."),
			"emptyText":  contentProperty("Empty-state copy."),
		},
		nil,
		deliveryDesign("data", "03 Molecules", "Data", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "data", tableExampleProps()),
			deliveryExample(componentType, "empty", map[string]any{
				"columns": []map[string]any{
					{"key": "name", "label": "Name", "primary": true},
					{"key": "status", "label": "Status"},
				},
				"rows":      []map[string]any{},
				"emptyText": "No workspaces found.",
			}),
			deliveryExample(componentType, "sortable", map[string]any{
				"columns": []map[string]any{
					{"key": "name", "label": "Name", "sortable": true, "primary": true},
					{"key": "status", "label": "Status", "sortable": true},
					{"key": "updated", "label": "Updated", "sortable": true},
				},
				"rows":     tableExampleProps()["rows"],
				"sortable": true,
			}),
		},
		Table,
	)
}

func cardDeliveryDefinition() DeliveryDefinition {
	componentType := "Card"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clCard.Compile(),
			deliveryText("Title", clCardTitle.Compile(), "Quarterly summary", "title"),
			deliveryText("Description", clCardDesc.Compile(), "Performance is tracking above plan.", "description"),
			deliverySlot(
				"Content",
				"children",
				"",
				deliveryInstance("Body", "Text", "body", ""),
			),
		),
		"clickable",
		map[string]string{"false": "", "true": clCardClickable.Compile()},
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Composable content surface with optional navigation behavior.",
		map[string]PropertyContract{
			"title":         contentProperty("Optional card heading."),
			"description":   contentProperty("Optional supporting copy."),
			"image":         contentProperty("Optional media source."),
			"imageAlt":      contentProperty("Alternative text for card media."),
			"imagePosition": variantProperty([]string{"top", "bottom", "left", "right"}, "top", "Media placement."),
			"variant":       variantProperty([]string{"default", "elevated", "outlined", "plain"}, "default", "Surface treatment."),
			"padding":       variantProperty([]string{"none", "small", "medium", "large"}, "medium", "Content spacing."),
			"shadow":        variantProperty([]string{"none", "small", "medium", "large"}, "small", "Surface elevation."),
			"clickable":     stateProperty("Whether the whole card is interactive."),
			"hoverable":     stateProperty("Whether the card lifts on hover."),
			"href":          contentProperty("Navigation target for clickable cards."),
		},
		[]designcomponent.Slot{{
			Name:         "children",
			Description:  "Ordered card content.",
			AllowedTypes: deliveryAtomicContentTypes(),
			Cardinality:  designcomponent.SlotMany,
			Attrs:        deliveryRepeatingSlotAttrs(),
		}},
		deliveryDesign("summary", "03 Molecules", "Surfaces", root),
		[]DeliveryExample{
			cardDeliveryExample(
				canonicalDeliveryExample(componentType, "summary", map[string]any{
					"title": "Quarterly summary", "description": "Performance is tracking above plan.",
				}),
				"Review the current operational details.",
			),
			cardDeliveryExample(
				mobileDeliveryExample(componentType, "mobile-summary", map[string]any{
					"title": "Open reviews", "description": "Three items need attention.",
				}),
				"Open the review queue.",
			),
			cardDeliveryExample(
				deliveryExample(componentType, "basic", map[string]any{
					"title": "Quarterly summary", "description": "Performance is tracking above plan.",
				}),
				"Review the current operational details.",
			),
		},
		func(props molecules.CardProps, slots DeliverySlotChildren) g.Node {
			return Card(props, slots.Nodes("children")...)
		},
	)
}

func cardDeliveryExample(
	example DeliveryExample,
	content string,
) DeliveryExample {
	return withDeliveryExampleSlots(
		example,
		deliveryExampleSlot(
			"children",
			deliveryExampleComponent(
				"body",
				"Text",
				map[string]any{
					"content": content,
					"size":    "sm",
					"color":   "muted",
				},
			),
		),
	)
}

func breadcrumbDeliveryDefinition() DeliveryDefinition {
	componentType := "Breadcrumb"
	root := deliveryFrame(
		componentType,
		clBreadcrumb.Compile(),
		deliveryText("Parent", clLink.Compile(), "Workspaces"),
		deliveryText("Separator", clBreadcrumbSep.Compile(), "/"),
		deliveryText("Current", clBreadcrumbCur.Compile(), "Alpha"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible hierarchical navigation trail.",
		map[string]PropertyContract{
			"items":     contentProperty("Ordered navigation segments."),
			"separator": contentProperty("Visual separator."),
			"maxItems":  modifierProperty("Maximum visible segments before collapsing."),
		},
		nil,
		deliveryDesign("default", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"items": []map[string]any{
					{"label": "Workspaces", "href": "/workspaces"},
					{"label": "Alpha", "current": true},
				},
			}),
			deliveryExample(componentType, "with-chevrons", map[string]any{
				"separator": "›",
				"items": []map[string]any{
					{"label": "Home", "href": "/"},
					{"label": "Products", "href": "/products"},
					{"label": "Widget Pro", "current": true},
				},
			}),
		},
		Breadcrumb,
	)
}

func accordionDeliveryDefinition() DeliveryDefinition {
	componentType := "Accordion"
	root := deliveryFrame(
		componentType,
		clAccordionRoot.Merge(clAccordionBordered).Compile(),
		deliveryFrame(
			"Section",
			"",
			deliveryFrame(
				"Trigger",
				clAccordionTrigger.Compile(),
				deliveryInstance("LeadingIcon", "Icon", "", clAccordionItemIcon.Compile()),
				deliveryFrame(
					"TitleBlock",
					clAccordionTitleBlock.Compile(),
					deliveryText("Title", clAccordionTitle.Compile(), "What is PlatformKit?"),
					deliveryText("Subtitle", clAccordionSubtitle.Compile(), "Expand to learn more"),
				),
				deliveryInstance("Chevron", "Icon", "", clAccordionChevron.Compile()),
			),
			deliverySlot(
				"Content",
				"section",
				clAccordionPanel.Compile(),
				deliveryInstance("Body", "Text", "body", ""),
			),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible collapsible content sections with single- or multi-open controller behavior.",
		map[string]PropertyContract{
			"items":       contentProperty("Ordered string-backed sections for portable direct rendering."),
			"multiple":    stateProperty("Whether multiple sections may stay open."),
			"defaultOpen": stateProperty("Section IDs open on first render."),
			"bordered": {
				Role:        designcomponent.PropRoleModifier,
				Description: "Whether the non-flush surface has an outer border.",
				Default:     "true",
			},
			"flush": modifierProperty("Whether outer surface chrome is removed for embedding."),
		},
		[]designcomponent.Slot{{
			Name:         "section",
			Description:  "Ordered rich panel content with disclosure metadata.",
			AllowedTypes: deliveryMoleculeContentTypes(),
			Cardinality:  designcomponent.SlotMany,
			Attrs:        accordionDeliverySlotAttrs(),
		}},
		deliveryDesign("basic", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "basic", map[string]any{
					"defaultOpen": []string{"platform"},
				}),
				deliveryExampleSlot(
					"section",
					accordionDeliveryExampleSection(
						"platform",
						"What is PlatformKit?",
						"A governed platform for composing production software.",
					),
					accordionDeliveryExampleSection(
						"start",
						"How do I get started?",
						"Declare the product once, then materialize each supported surface.",
					),
				),
			),
			deliveryExample(componentType, "portable-items", map[string]any{
				"multiple": true,
				"items": []map[string]any{
					{"id": "general", "title": "General", "content": "General settings", "open": true},
					{"id": "security", "title": "Security", "content": "Security settings", "disabled": true},
				},
			}),
			deliveryExample(componentType, "flush", map[string]any{
				"flush": true,
				"items": []map[string]any{
					{"id": "details", "title": "Order details", "content": "Order content"},
				},
			}),
		},
		func(props molecules.AccordionProps, slots DeliverySlotChildren) g.Node {
			if len(slots["section"]) == 0 {
				return Accordion(props)
			}
			sections := make([]AccordionSection, 0, len(slots["section"]))
			for _, instance := range slots["section"] {
				sections = append(sections, AccordionSection{
					ID:          deliveryString(instance.Attrs, "id"),
					Title:       deliveryString(instance.Attrs, "title"),
					Subtitle:    deliveryString(instance.Attrs, "subtitle"),
					Icon:        deliveryString(instance.Attrs, "icon"),
					Content:     instance.Children,
					DefaultOpen: deliveryBool(instance.Attrs, "open"),
					Disabled:    deliveryBool(instance.Attrs, "disabled"),
				})
			}
			return AccordionWithSections(props, sections...)
		},
	)
}

func accordionDeliverySlotAttrs() []designcomponent.Prop {
	return []designcomponent.Prop{
		{Name: "id", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Required: true, Description: "Stable section identity."},
		{Name: "title", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Required: true, Description: "Visible section title."},
		{Name: "subtitle", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Description: "Optional supporting title copy."},
		{Name: "icon", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Description: "Optional OSS icon name."},
		{Name: "open", Type: designcomponent.PropBoolean, Role: designcomponent.PropRoleState, Default: "false", Description: "Whether the section opens on first render."},
		{Name: "disabled", Type: designcomponent.PropBoolean, Role: designcomponent.PropRoleState, Default: "false", Description: "Whether the trigger cannot be activated."},
	}
}

func accordionDeliveryExampleSection(id, title, content string) DeliveryExampleComponent {
	component := deliveryExampleComponent(
		id+"-body",
		"Text",
		map[string]any{"content": content, "size": "sm", "color": "muted"},
	)
	component.SlotAttrs = map[string]any{"id": id, "title": title}
	return component
}

func stepperDeliveryDefinition() DeliveryDefinition {
	componentType := "Stepper"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clStepperListHorizontal.Compile(),
			deliveryFrame(
				"CompletedStep",
				clStepperItemHorizontal.Compile(),
				deliveryInstance(
					"Indicator",
					"Icon",
					"",
					clStepperIndicator.
						Merge(clStepperIndicatorRegular).
						Merge(clStepperIndicatorState["completed"]).
						Compile(),
				),
				deliveryText(
					"Label",
					clStepperLabel.Merge(clStepperLabelState["completed"]).Compile(),
					"Account",
				),
				deliveryFrame(
					"Connector",
					clStepperConnectorHorizontal.Merge(clStepperConnectorState["completed"]).Compile(),
				),
			),
			deliveryFrame(
				"ActiveStep",
				clStepperItemHorizontal.Compile(),
				deliveryText(
					"Indicator",
					clStepperIndicator.
						Merge(clStepperIndicatorRegular).
						Merge(clStepperIndicatorState["active"]).
						Merge(clStepperGlyph).
						Compile(),
					"2",
				),
				deliveryText(
					"Label",
					clStepperLabel.Merge(clStepperLabelState["active"]).Compile(),
					"Profile",
				),
			),
		),
		"orientation",
		map[string]string{
			"horizontal": clStepperListHorizontal.Compile(),
			"vertical":   clStepperListVertical.Compile(),
		},
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible horizontal or vertical progress navigation with derived and explicit step states.",
		map[string]PropertyContract{
			"steps":           contentProperty("Ordered steps and their optional stable keys, descriptions, icons, and explicit statuses."),
			"currentStep":     contentProperty("Zero-based active step used to derive omitted statuses."),
			"orientation":     variantProperty([]string{"horizontal", "vertical"}, "horizontal", "Step layout direction."),
			"clickable":       stateProperty("Whether individual indicators are keyboard-operable navigation controls."),
			"compact":         modifierProperty("Whether indicators and labels use compact sizing."),
			"stepAction":      contentProperty("Optional controller action that makes each whole step operable."),
			"navigationLabel": contentProperty("Accessible navigation landmark label."),
		},
		nil,
		deliveryDesign("checkout", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "checkout", map[string]any{
				"currentStep": 1,
				"steps": []map[string]any{
					{"key": "account", "label": "Account"},
					{"key": "profile", "label": "Profile", "description": "Tell us about yourself"},
					{"key": "review", "label": "Review"},
				},
			}),
			deliveryExample(componentType, "vertical", map[string]any{
				"orientation": "vertical",
				"currentStep": 1,
				"steps": []map[string]any{
					{"label": "Create account", "description": "Set up your credentials"},
					{"label": "Profile", "description": "Tell us about yourself"},
					{"label": "Preferences", "description": "Customize your experience"},
				},
			}),
			deliveryExample(componentType, "clickable-error", map[string]any{
				"clickable":   true,
				"compact":     true,
				"currentStep": 2,
				"steps": []map[string]any{
					{"key": "cart", "label": "Cart"},
					{"key": "shipping", "label": "Shipping", "status": "error"},
					{"key": "payment", "label": "Payment"},
				},
			}),
		},
		Stepper,
	)
}

func sidebarDeliveryDefinition() DeliveryDefinition {
	componentType := "Sidebar"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clSidebarRootAdmin.Merge(clSidebarWidthExpanded).Compile(),
				deliverySlot(
					"Brand",
					"brand",
					clSidebarBrandAdmin.Compile(),
					deliveryInstance("BrandLabel", "Heading", "", ""),
				),
				deliveryFrame(
					"Navigation",
					clSidebarNavAdmin.Compile(),
					deliveryFrame(
						"ActiveItem",
						clSidebarLinkAdmin.
							Merge(clSidebarLinkPadExpanded).
							Merge(clSidebarLinkActiveAdmin).
							Compile(),
						deliveryInstance("Icon", "Icon", "", ""),
						deliveryText("Label", clSidebarLabelVisible.Compile(), "Customers"),
						deliveryInstance("Badge", "Badge", "", ""),
					),
					deliveryFrame(
						"IdleItem",
						clSidebarLinkAdmin.
							Merge(clSidebarLinkPadExpanded).
							Merge(clSidebarLinkIdleAdmin).
							Compile(),
						deliveryInstance("Icon", "Icon", "", ""),
						deliveryText("Label", clSidebarLabelVisible.Compile(), "Reports"),
					),
				),
				deliverySlot(
					"Footer",
					"footer",
					clSidebarFooterAdmin.Compile(),
					deliveryInstance("FooterContent", "Text", "", ""),
				),
			),
			"collapsed",
			map[string]string{
				"false": clSidebarRootAdmin.Merge(clSidebarWidthExpanded).Compile(),
				"true":  clSidebarRootAdmin.Merge(clSidebarWidthCollapsed).Compile(),
			},
		),
		"flavor",
		map[string]string{
			"admin":   clSidebarRootAdmin.Merge(clSidebarWidthExpanded).Compile(),
			"content": clSidebarRootContent.Compile(),
		},
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Persistent admin or content navigation with grouped, nested, searchable items and collapsed-rail support.",
		map[string]PropertyContract{
			"items":           contentProperty("Flat or nested portable navigation items."),
			"sections":        contentProperty("Grouped portable navigation sections."),
			"current":         stateProperty("Current route used to derive active items and parent groups."),
			"flavor":          variantProperty([]string{"admin", "content"}, "admin", "Surface treatment and responsive behavior."),
			"collapsible":     stateProperty("Whether the surface participates in collapsed-rail control."),
			"collapsed":       stateProperty("Whether the admin rail starts collapsed."),
			"navigationLabel": contentProperty("Accessible navigation landmark label."),
			"brandLabel":      contentProperty("Portable brand label used when the rich brand slot is empty."),
			"brandHref":       contentProperty("Portable brand navigation target."),
		},
		[]designcomponent.Slot{
			{
				Name:         "brand",
				Description:  "Optional rich brand or documentation heading.",
				AllowedTypes: deliveryMoleculeContentTypes(),
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "footer",
				Description:  "Optional secondary navigation or account content.",
				AllowedTypes: deliveryMoleculeContentTypes(),
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("default", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"current":         "/admin/customers/accounts",
					"navigationLabel": "Workspace navigation",
					"sections": []map[string]any{{
						"id": "operate", "label": "Operate",
						"items": []map[string]any{
							{"id": "dashboard", "label": "Dashboard", "href": "/admin", "icon": "home"},
							{"id": "customers", "label": "Customers", "href": "/admin/customers", "icon": "users", "badge": "24", "children": []map[string]any{
								{"id": "accounts", "label": "Accounts", "href": "/admin/customers/accounts"},
							}},
						},
					}},
				}),
				deliveryExampleSlot("brand", deliveryExampleComponent(
					"workspace-brand", "Heading", map[string]any{"text": "PlatformKit", "level": 2},
				)),
			),
			deliveryExample(componentType, "docs-content", map[string]any{
				"flavor":          "content",
				"current":         "#figma",
				"navigationLabel": "Documentation sections",
				"brandLabel":      "Documentation",
				"brandHref":       "#overview",
				"sections": []map[string]any{
					{"id": "start", "label": "Start", "glyph": "01", "tone": "brand", "items": []map[string]any{
						{"id": "overview", "label": "Overview", "href": "#overview", "prefix": "01"},
						{"id": "figma", "label": "Figma handoff", "href": "#figma", "prefix": "02"},
					}},
				},
			}),
			deliveryExample(componentType, "collapsed-rail", map[string]any{
				"collapsible": true,
				"collapsed":   true,
				"current":     "/admin/workflows",
				"items": []map[string]any{
					{"label": "Home", "href": "/admin", "icon": "home"},
					{"label": "Workflows", "href": "/admin/workflows", "icon": "workflow"},
					{"label": "Reports", "href": "/admin/reports", "icon": "chart"},
				},
			}),
		},
		func(props molecules.SidebarProps, slots DeliverySlotChildren) g.Node {
			return SidebarWithSlots(props, SidebarSlots{
				Brand:  slots.Nodes("brand"),
				Footer: slots.Nodes("footer"),
			})
		},
	)
}

func datePickerDeliveryDefinition() DeliveryDefinition {
	componentType := "DatePicker"
	controlClass := clInput.
		Merge(clInputNormal).
		Merge(clInputSize["md"]).
		Merge(clInputPadStart).
		Compile()
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFieldWrap.Compile(),
			deliveryText("Label", clLabel.Compile(), "Start Date", "label"),
			deliveryClassBound(
				deliveryFrame(
					"Control",
					clInputIconWrap.Compile(),
					deliveryInstance("CalendarIcon", "Icon", "", clInputIconStart.Compile()),
					deliveryText("Value", controlClass, "2026-04-06", "value", "placeholder"),
				),
				"error",
				map[string]string{
					"":  clInputIconWrap.Compile(),
					"*": clInputIconWrap.Compile(),
				},
			),
			deliveryText("Supporting", clHelp.Compile(), "Choose an available date.", "helpText", "error"),
		),
		"field",
		"date field",
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Native date field with accessible feedback, constraints, calendar affordance, and optional HTMX change requests.",
		map[string]PropertyContract{
			"name":        contentProperty("Submitted field name."),
			"value":       contentProperty("Selected ISO date value."),
			"min":         contentProperty("Minimum selectable ISO date."),
			"max":         contentProperty("Maximum selectable ISO date."),
			"format":      contentProperty("Optional consumer display-format metadata."),
			"label":       contentProperty("Visible field label."),
			"placeholder": contentProperty("Empty-state prompt for supporting browsers."),
			"helpText":    contentProperty("Supporting guidance."),
			"error":       contentProperty("Validation message; non-empty selects error styling."),
			"invalid":     stateProperty("Whether aria-invalid and error styling are active without an inline message."),
			"required":    stateProperty("Whether a date is required."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"name": "start_date", "label": "Start Date", "required": true,
			}),
			mobileDeliveryExample(componentType, "with-constraints", map[string]any{
				"name": "booking_date", "label": "Booking date",
				"min": "2026-01-01", "max": "2026-12-31",
				"helpText": "Choose a date in 2026.",
			}),
			deliveryExample(componentType, "validation-error", map[string]any{
				"name": "due_date", "label": "Due date", "error": "Choose a due date.",
			}),
			deliveryExample(componentType, "htmx-dynamic", map[string]any{
				"name": "event_date", "label": "Event date",
				"hx-get": "/api/availability", "hx-target": "#availability-results",
			}),
		},
		DatePicker,
	)
}

func fileUploadDeliveryDefinition() DeliveryDefinition {
	componentType := "FileUpload"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFileUploadRoot.Compile(),
			deliveryText("Label", clLabel.Compile(), "Profile photo", "label", "ariaLabel"),
			deliveryFrame(
				"Control",
				"",
				deliveryFrame(
					"DropZone",
					clFileUploadDropZone.Compile(),
					deliverySlot(
						"Prompt", "content", clFileUploadDropZoneInner.Compile(),
						deliveryInstance("UploadIcon", "Icon", "", clFileUploadIcon.Compile()),
						deliveryText("Action", clFileUploadPromptAction.Compile(), "Click to upload", "promptLabel"),
						deliveryText("DropHint", clFileUploadPromptText.Compile(), "or drag and drop", "dropLabel"),
					),
				),
				deliveryFrame(
					"FileList",
					clFileUploadList.Compile(),
					deliveryFrame(
						"FileItem",
						clFileUploadItem.Compile(),
						deliveryText("FileName", clFileUploadItemName.Compile(), "portrait.png", "currentName"),
						deliveryText("Remove", clFileUploadRemove.Compile(), "Remove file", "removeLabel"),
					),
				),
			),
			deliveryText("Supporting", clHelp.Compile(), "PNG or JPEG, up to 5 MB.", "helpText", "maxSizeLabel"),
		),
		"field",
		"file upload field",
	)
	definition := newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible native and remote file upload with drag-and-drop, validation, previews, localized feedback, and URL-backed form submission.",
		map[string]PropertyContract{
			"name":           contentProperty("Submitted field name; required for native files and remote URL values."),
			"accept":         contentProperty("Comma-separated MIME types, wildcards, or file extensions accepted by browser and controller validation."),
			"multiple":       stateProperty("Whether native multipart mode accepts multiple files; remote URL mode is intentionally single-file."),
			"required":       stateProperty("Whether native multipart selection is required."),
			"maxSize":        modifierProperty("Maximum accepted file size in bytes; non-positive values disable the client-side limit."),
			"dropZone":       {Role: designcomponent.PropRoleModifier, Default: "true", Description: "Whether to render the accessible drag-and-drop surface instead of a native file input."},
			"preview":        stateProperty("Whether selected image files receive thumbnail previews."),
			"showList":       {Role: designcomponent.PropRoleModifier, Default: "true", Description: "Whether selected files and validation errors are displayed."},
			"label":          contentProperty("Visible field label."),
			"helpText":       contentProperty("Supporting guidance associated with the file input."),
			"ariaLabel":      contentProperty("Accessible upload name when visible copy is insufficient."),
			"promptLabel":    contentProperty("Localized primary drop-zone action."),
			"dropLabel":      contentProperty("Localized drag-and-drop continuation."),
			"chooseLabel":    contentProperty("Localized file-picker label template; {label} is replaced with the accessible field name."),
			"loadingLabel":   contentProperty("Localized in-progress upload status."),
			"removeLabel":    contentProperty("Localized selected-file removal action."),
			"uploadedLabel":  contentProperty("Localized metadata for a URL-backed uploaded file."),
			"maxSizeLabel":   contentProperty("Localized size-hint template; {size} is replaced with the formatted limit."),
			"missingError":   contentProperty("Localized missing-file validation message."),
			"tooLargeError":  contentProperty("Localized size-error template supporting {name} and {maxSize}."),
			"typeError":      contentProperty("Localized type-error template supporting {name}."),
			"uploadURL":      contentProperty("Optional provider endpoint for immediate raw-body upload; enables URL-backed mode."),
			"uploadCategory": contentProperty("Optional category query value sent to the remote upload endpoint."),
			"value":          stateProperty("Initial URL submitted by remote upload mode."),
			"currentName":    contentProperty("Display name for the initial remote URL."),
		},
		[]designcomponent.Slot{{
			Name: "content", Description: "Optional branded drop-zone prompt replacing the default icon and copy.",
			AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
			Attrs: deliveryRepeatingSlotAttrs(),
		}},
		deliveryDesign("image-upload", "03 Molecules", "Data Entry & Presentation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "image-upload", map[string]any{
				"name": "avatar", "label": "Profile photo", "accept": "image/*",
				"preview": true, "maxSize": 5 * 1024 * 1024,
				"helpText": "PNG or JPEG, up to 5 MB.",
			}),
			mobileDeliveryExample(componentType, "simple-input", map[string]any{
				"name": "attachment", "label": "Attachment",
				"dropZone": false, "showList": false, "required": true,
			}),
			deliveryExample(componentType, "multi-file-upload", map[string]any{
				"name": "documents", "label": "Documents", "multiple": true,
				"accept": ".pdf,.doc,.docx",
			}),
			deliveryExample(componentType, "remote-upload", map[string]any{
				"name": "logo_url", "label": "Company logo", "accept": "image/*",
				"preview": true, "uploadURL": "/api/v1/files/upload", "uploadCategory": "image",
				"value": "/media/logo.png", "currentName": "logo.png",
			}),
			deliveryExample(componentType, "htmx-upload", map[string]any{
				"name": "evidence", "label": "Evidence", "accept": ".pdf",
				"hx-post": "/api/evidence", "hx-target": "#evidence-results",
			}),
		},
		func(props molecules.FileUploadProps, slots DeliverySlotChildren) g.Node {
			return FileUploadWithSlots(props, FileUploadSlots{Content: slots.Nodes("content")})
		},
	)
	return withDeliveryTags(definition, "upload", "file", "form", "attachment", "media")
}

func autocompleteDeliveryDefinition() DeliveryDefinition {
	componentType := "Autocomplete"
	controlClass := clInput.
		Merge(clInputNormal).
		Merge(clInputSize["md"]).
		Merge(clInputPadStart).
		Compile()
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clFieldWrap.Compile(),
			deliveryText("Label", clLabel.Compile(), "Assign To", "label"),
			deliveryFrame(
				"Control",
				clAutocompleteControl.Compile(),
				deliveryInstance("SearchIcon", "Icon", "", clInputIconStart.Compile()),
				deliveryText("Query", controlClass, "Search users...", "displayValue", "placeholder"),
				deliveryFrame(
					"Results",
					clAutocompletePanel.Compile(),
					deliveryText("Option", clAutocompleteOption.Compile(), "Ada Lovelace"),
					deliveryText("Option", clAutocompleteOption.Compile(), "Grace Hopper"),
				),
			),
			deliveryText("Supporting", clHelp.Compile(), "Choose a teammate.", "helpText", "error"),
		),
		"field",
		"combobox field",
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible autocomplete with separate display and submitted values, static filtering, and server-backed HTMX suggestions.",
		map[string]PropertyContract{
			"name":         contentProperty("Submitted selected-value field name."),
			"searchURL":    contentProperty("Optional endpoint returning suggestion option fragments."),
			"queryName":    contentProperty("Server-search query parameter; defaults to q."),
			"options":      contentProperty("Deterministic local suggestions with labels and values."),
			"minChars":     modifierProperty("Minimum query length before suggestions open; defaults to 2."),
			"debounce":     modifierProperty("Server-search debounce in milliseconds; defaults to 300."),
			"value":        contentProperty("Initial submitted value."),
			"displayValue": contentProperty("Initial human-readable query value."),
			"placeholder":  contentProperty("Empty-state query prompt."),
			"label":        contentProperty("Visible field label."),
			"helpText":     contentProperty("Supporting guidance."),
			"error":        contentProperty("Validation message; non-empty selects error styling."),
			"invalid":      stateProperty("Whether the field is invalid without an inline error."),
			"required":     stateProperty("Whether a selection is required."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Selection & Filtering", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"name": "user_id", "label": "Assign To", "placeholder": "Search users...",
				"options": []map[string]any{
					{"value": "ada", "label": "Ada Lovelace"},
					{"value": "grace", "label": "Grace Hopper"},
					{"value": "katherine", "label": "Katherine Johnson"},
				},
			}),
			mobileDeliveryExample(componentType, "with-value", map[string]any{
				"name": "category_id", "label": "Category",
				"value": "cat-123", "displayValue": "Electronics",
				"options": []map[string]any{
					{"value": "cat-123", "label": "Electronics"},
					{"value": "cat-456", "label": "Services"},
				},
			}),
			deliveryExample(componentType, "keyboard-nav", map[string]any{
				"name": "product", "placeholder": "Search products...", "minChars": 1, "debounce": 200,
				"options": []map[string]any{
					{"value": "pkg-growth", "label": "Growth package"},
					{"value": "pkg-team", "label": "Team package"},
				},
			}),
			deliveryExample(componentType, "server-backed", map[string]any{
				"name": "company_id", "label": "Company", "searchURL": "/api/companies/search",
				"queryName": "q", "helpText": "Type at least two characters.",
			}),
			deliveryExample(componentType, "validation-error", map[string]any{
				"name": "owner_id", "label": "Owner", "error": "Choose a valid owner.",
			}),
		},
		Autocomplete,
	)
}

func dropdownDeliveryDefinition() DeliveryDefinition {
	componentType := "Dropdown"
	root := deliverySemantics(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clDropdownRoot.Compile(),
				deliveryText("Label", clLabel.Compile(), "Country", "label", "ariaLabel"),
				deliveryFrame(
					"Trigger",
					clDropdownTrigger.Merge(clDropdownTriggerSize["md"]).Compile(),
					deliveryText("Selection", clDropdownTriggerLabel.Compile(), "Select a country", "placeholder", "value"),
					deliveryFrame(
						"Actions",
						clDropdownTriggerActions.Compile(),
						deliveryInstance("Clear", "Icon", "", clDropdownIconButton.Compile()),
						deliveryInstance("Chevron", "Icon", "", clDropdownChevron.Compile()),
					),
				),
				deliveryFrame(
					"Panel",
					clDropdownPanel.Compile(),
					deliveryText("Search", clDropdownSearch.Compile(), "Search options", "searchLabel"),
					deliveryFrame(
						"Options",
						clDropdownOptions.Compile(),
						deliveryText("Option", clDropdownOption.Compile(), "Portugal"),
						deliveryText("SelectedOption", clDropdownOption.Merge(clDropdownOptionSelected).Compile(), "United States"),
					),
				),
			),
			"size",
			compileClassMap(clDropdownTriggerSize),
		),
		"field",
		"select-only listbox field",
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible select-only listbox with optional filtering, clearing, grouping, icons, and multiple submitted values.",
		map[string]PropertyContract{
			"name":        contentProperty("Submitted hidden-input field name."),
			"options":     contentProperty("Ordered options with value, label, icon, group, and disabled state."),
			"placeholder": contentProperty("Text shown while no option is selected."),
			"ariaLabel":   contentProperty("Accessible trigger name when no visible label is present."),
			"searchLabel": contentProperty("Localized search prompt and accessible name."),
			"clearLabel":  contentProperty("Localized clear-action label."),
			"openLabel":   contentProperty("Localized disclosure-action label."),
			"searchable":  stateProperty("Whether the options panel includes client-side filtering."),
			"clearable":   stateProperty("Whether selected values can be cleared explicitly."),
			"multiple":    stateProperty("Whether multiple values may be submitted."),
			"value":       stateProperty("Initial single selected value."),
			"selected":    stateProperty("Initial selected values; used by multiple mode."),
			"size":        sizeProperty([]string{"sm", "md", "lg"}, "md", "Control height and text scale."),
			"label":       contentProperty("Optional visible field label."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Selection & Filtering", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"name": "country", "label": "Country", "placeholder": "Select a country",
				"options": []map[string]any{
					{"value": "us", "label": "United States"},
					{"value": "pt", "label": "Portugal"},
					{"value": "de", "label": "Germany"},
				},
			}),
			mobileDeliveryExample(componentType, "searchable", map[string]any{
				"name": "assignee", "label": "Assignee", "placeholder": "Choose a teammate",
				"searchable": true, "clearable": true, "value": "ada",
				"options": []map[string]any{
					{"value": "ada", "label": "Ada Lovelace", "icon": "user", "group": "Workspace"},
					{"value": "grace", "label": "Grace Hopper", "icon": "user", "group": "Workspace"},
				},
			}),
			deliveryExample(componentType, "multi-select", map[string]any{
				"name": "topics", "ariaLabel": "Topics", "multiple": true, "clearable": true,
				"selected": []string{"alerts", "release"},
				"options": []map[string]any{
					{"value": "alerts", "label": "Alerts"},
					{"value": "release", "label": "Releases"},
					{"value": "billing", "label": "Billing", "disabled": true},
				},
			}),
			deliveryExample(componentType, "htmx-filter", map[string]any{
				"name": "status", "ariaLabel": "Status", "placeholder": "All statuses",
				"hx-get": "/workspaces", "hx-target": "#workspace-grid", "hx-swap": "outerHTML",
				"options": []map[string]any{
					{"value": "active", "label": "Active"},
					{"value": "paused", "label": "Paused"},
				},
			}),
		},
		Dropdown,
	)
}

func actionMenuDeliveryDefinition() DeliveryDefinition {
	componentType := "ActionMenu"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clActionMenuRoot.Compile(),
				deliveryFrame(
					"Trigger",
					clActionMenuTrigger.Compile(),
					deliveryInstance("TriggerIcon", "Icon", "", ""),
					deliveryText("TriggerLabel", "", "Actions", "triggerLabel"),
				),
				deliveryFrame(
					"Panel",
					clActionMenuPanel.Merge(clActionMenuAlign["end"]).Merge(clActionMenuWidth["md"]).Compile(),
					deliveryText("SectionLabel", clActionMenuSectionLabel.Compile(), "Account"),
					deliveryText("Item", clActionMenuItem.Merge(clActionMenuItemTone["neutral"]).Compile(), "View details"),
					deliveryText("SecondaryItem", clActionMenuItem.Merge(clActionMenuItemTone["neutral"]).Compile(), "Edit"),
					deliveryText("DangerItem", clActionMenuItem.Merge(clActionMenuItemTone["danger"]).Compile(), "Delete"),
				),
			),
			"align",
			compileClassMap(clActionMenuAlign),
		),
		"width",
		compileClassMap(clActionMenuWidth),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible menu button for navigation and server-driven row actions with complete keyboard and focus behavior.",
		map[string]PropertyContract{
			"triggerLabel": contentProperty("Visible trigger text; an icon-only trigger remains named Actions."),
			"triggerIcon":  contentProperty("Trigger icon name; defaults to ellipsis-vertical."),
			"align":        variantProperty([]string{"start", "end"}, "end", "Panel alignment relative to the trigger."),
			"width":        variantProperty([]string{"sm", "md", "lg"}, "md", "Governed panel width preset."),
			"items":        contentProperty("Ungrouped action or navigation items."),
			"sections":     contentProperty("Ordered labeled or unlabeled item groups."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"items": []map[string]any{
					{"label": "View details", "icon": "eye", "href": "/items/42"},
					{"label": "Edit", "icon": "pencil", "href": "/items/42/edit"},
					{"label": "Delete", "icon": "trash", "danger": true, "action": "delete-item"},
				},
			}),
			mobileDeliveryExample(componentType, "labeled", map[string]any{
				"triggerLabel": "Manage", "triggerIcon": "chevron-down", "align": "start", "width": "lg",
				"sections": []map[string]any{{
					"label": "Account",
					"items": []map[string]any{{"label": "Profile", "icon": "user", "href": "/account"}},
				}},
			}),
			deliveryExample(componentType, "server-actions", map[string]any{
				"items": []map[string]any{
					{"label": "Retry", "icon": "arrow-path", "hxPost": "/jobs/42/retry", "hxTarget": "#jobs", "hxSwap": "innerHTML"},
					{"label": "Cancel", "icon": "x-circle", "danger": true, "hxPost": "/jobs/42/cancel", "hxConfirm": "Cancel this job?"},
				},
			}),
			deliveryExample(componentType, "disabled", map[string]any{
				"items": []map[string]any{{"label": "Unavailable", "disabled": true}},
			}),
		},
		ActionMenu,
	)
}

func drawerDeliveryDefinition() DeliveryDefinition {
	componentType := "Drawer"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clDrawerRoot.Compile(),
				deliveryFrame("Overlay", clDrawerOverlay.Compile()),
				deliveryFrame(
					"Panel",
					clDrawerPanel.Merge(clDrawerPosition["right"]).Merge(clDrawerWidth["medium"]).Compile(),
					deliveryFrame(
						"Header",
						clDrawerHeader.Compile(),
						deliverySlot("HeaderContent", "header", "",
							deliveryText("Title", clDrawerTitle.Compile(), "Account settings", "title"),
						),
						deliveryText("Description", clDrawerDescription.Compile(), "Manage your profile and preferences.", "description"),
						deliveryInstance("Close", "Icon", "", clDrawerClose.Compile()),
					),
					deliverySlot("Body", "content", clDrawerBody.Compile(),
						deliveryInstance("BodyContent", "Text", "body", ""),
					),
					deliverySlot("Footer", "actions", clDrawerFooter.Compile(),
						deliveryInstance("FooterAction", "Button", "primary", ""),
					),
				),
			),
			"position",
			compileClassMap(clDrawerPosition),
		),
		"size",
		compileClassMap(clDrawerWidth),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible modal edge panel with governed focus containment, focus return, scroll locking, and server-open behavior.",
		map[string]PropertyContract{
			"title":       contentProperty("Default drawer heading when the rich header slot is empty."),
			"description": contentProperty("Optional supporting text beneath the default heading."),
			"body":        contentProperty("Portable body copy used when the rich body slot is empty."),
			"footer":      contentProperty("Portable footer copy used when the rich footer slot is empty."),
			"ariaLabel":   contentProperty("Accessible name override; defaults to the title."),
			"closeLabel":  contentProperty("Localized close-button label; defaults to Close."),
			"position":    variantProperty([]string{"left", "right", "bottom"}, "right", "Viewport edge from which the panel opens."),
			"size":        variantProperty([]string{"small", "medium", "large", "xl", "full"}, "medium", "Governed width or bottom-sheet height."),
			"closable": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the close button is shown.",
			},
			"closeOnOverlay": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether backdrop activation dismisses the drawer.",
			},
			"closeOnEscape": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether Escape dismisses the drawer.",
			},
			"showOverlay": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the modal backdrop is visible.",
			},
			"open":       stateProperty("Whether the drawer is open on first render."),
			"openOnSwap": stateProperty("Whether an HTMX content swap opens the drawer."),
		},
		[]designcomponent.Slot{
			{
				Name: "header", Description: "Optional rich header replacing the default title and description.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
			{
				Name: "content", Description: "Scrollable drawer content.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
			{
				Name: "actions", Description: "Optional persistent drawer actions.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("right-panel", "03 Molecules", "Overlays", root),
		[]DeliveryExample{
			drawerDeliveryExample(canonicalDeliveryExample(componentType, "right-panel", map[string]any{
				"title": "Account settings", "description": "Manage your profile and preferences.",
				"position": "right", "size": "medium", "open": true,
			})),
			drawerDeliveryExample(mobileDeliveryExample(componentType, "bottom-sheet", map[string]any{
				"title": "Filters", "position": "bottom", "size": "large", "open": true,
			})),
			drawerDeliveryExample(deliveryExample(componentType, "server-loaded", map[string]any{
				"title": "Order details", "position": "right", "size": "large", "open": true, "openOnSwap": true,
			})),
			drawerDeliveryExample(deliveryExample(componentType, "required-decision", map[string]any{
				"title": "Decision required", "ariaLabel": "Required decision", "position": "bottom", "size": "medium", "open": true,
				"closable": false, "closeOnOverlay": false, "closeOnEscape": false,
			})),
		},
		func(props molecules.DrawerProps, slots DeliverySlotChildren) g.Node {
			return DrawerWithSlots(props, DrawerSlots{
				Header: slots.Nodes("header"),
				Body:   slots.Nodes("content"),
				Footer: slots.Nodes("actions"),
			})
		},
	)
}

func drawerDeliveryExample(example DeliveryExample) DeliveryExample {
	return withDeliveryExampleSlots(
		example,
		// The nested Text instance may only set props the Figma ComponentSet
		// represents natively; size is a class-bound design axis, so an
		// instance-level override cannot round-trip and fails delivery
		// validation. Body copy inherits the authored baseline instead.
		deliveryExampleSlot("content", deliveryExampleComponent("body-copy", "Text", map[string]any{
			"content": "Review and update the information in this panel.",
		})),
		deliveryExampleSlot("actions", deliveryExampleComponent("save", "Button", map[string]any{
			"label": "Save changes", "variant": "primary",
		})),
	)
}

func modalDeliveryDefinition() DeliveryDefinition {
	componentType := "Modal"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clModalRoot.Merge(clModalCentered).Compile(),
			deliveryFrame("Backdrop", clModalOverlay.Compile()),
			deliveryFrame(
				"Panel",
				clModalPanel.Merge(clModalPanelSize["medium"]).Compile(),
				deliveryFrame(
					"Header",
					clModalHeader.Compile(),
					deliverySlot("HeaderContent", "header", "",
						deliveryText("Title", clModalTitle.Compile(), "Archive project", "title"),
					),
					deliveryText("Description", clModalDescription.Compile(), "Confirm this irreversible action.", "description"),
					deliveryInstance("Close", "Icon", "", clModalClose.Compile()),
				),
				deliverySlot("Body", "content", clModalBody.Compile(),
					deliveryInstance("BodyContent", "Text", "body", ""),
				),
				deliverySlot("Footer", "actions", clModalFooter.Compile(),
					deliveryInstance("FooterAction", "Button", "primary", ""),
				),
			),
		),
		"size",
		compileClassMap(clModalPanelSize),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible modal dialog and mobile sheet with governed dismissal, focus containment, focus return, scroll locking, and server-loaded content.",
		map[string]PropertyContract{
			"title":       contentProperty("Default dialog heading when the rich header slot is empty."),
			"description": contentProperty("Optional supporting text beneath the default heading."),
			"body":        contentProperty("Portable body copy used when the rich content slot is empty."),
			"footer":      contentProperty("Portable footer copy used when the rich actions slot is empty."),
			"ariaLabel":   contentProperty("Accessible name override; defaults to the title."),
			"closeLabel":  contentProperty("Localized close-button label; defaults to Close."),
			"size":        variantProperty([]string{"small", "medium", "large", "xl", "full"}, "medium", "Governed dialog width preset."),
			"closable": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the dialog can expose dismissal controls.",
			},
			"closeOnOverlay": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether backdrop activation dismisses the dialog.",
			},
			"closeOnEscape": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether Escape dismisses the dialog.",
			},
			"showClose": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the header close affordance is visible.",
			},
			"showOverlay": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the modal backdrop is visible.",
			},
			"centered": {
				Role: designcomponent.PropRoleModifier, Default: "true",
				Description: "Whether the panel is vertically centered instead of using the mobile-sheet alignment.",
			},
			"clearOnClose": stateProperty("Whether closing removes server-swapped content."),
			"open":         stateProperty("Whether the dialog is open on first render."),
			"openOnSwap":   stateProperty("Whether an HTMX content swap opens the dialog."),
			"deferred":     stateProperty("Whether this is an initially empty server-content swap target."),
		},
		[]designcomponent.Slot{
			{
				Name: "header", Description: "Optional rich header replacing the default title and description.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
			{
				Name: "content", Description: "Scrollable dialog content.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
			{
				Name: "actions", Description: "Optional persistent dialog actions.",
				AllowedTypes: deliveryMoleculeContentTypes(), Cardinality: designcomponent.SlotMany,
				Attrs: deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("confirm-dialog", "03 Molecules", "Overlays", root),
		[]DeliveryExample{
			modalDeliveryExample(canonicalDeliveryExample(componentType, "confirm-dialog", map[string]any{
				"title": "Archive project", "description": "This action cannot be undone.",
				"size": "medium", "open": true,
			})),
			modalDeliveryExample(mobileDeliveryExample(componentType, "mobile-sheet", map[string]any{
				"title": "Filter results", "centered": false, "size": "small", "open": true,
			})),
			deliveryExample(componentType, "server-loaded", map[string]any{
				"ariaLabel": "Loading dialog", "deferred": true, "openOnSwap": true, "clearOnClose": true,
			}),
			modalDeliveryExample(deliveryExample(componentType, "required-decision", map[string]any{
				"title": "Decision required", "size": "large", "open": true,
				"closable": false, "showClose": false, "closeOnOverlay": false, "closeOnEscape": false,
			})),
		},
		func(props molecules.ModalProps, slots DeliverySlotChildren) g.Node {
			return ModalWithSlots(props, ModalSlots{
				Header: slots.Nodes("header"),
				Body:   slots.Nodes("content"),
				Footer: slots.Nodes("actions"),
			})
		},
	)
}

func modalDeliveryExample(example DeliveryExample) DeliveryExample {
	return withDeliveryExampleSlots(
		example,
		deliveryExampleSlot("content", deliveryExampleComponent("body-copy", "Text", map[string]any{
			"content": "Review the consequences before continuing.",
		})),
		deliveryExampleSlot("actions", deliveryExampleComponent("confirm", "Button", map[string]any{
			"label": "Continue", "variant": "primary",
		})),
	)
}

func paginationDeliveryDefinition() DeliveryDefinition {
	componentType := "Pagination"
	numberedTotals := []any{4, 7}
	pageOne := deliveryClassBound(
		deliveryVisibleWhen(
			deliveryText("PageOne", clPageBtn.Merge(clPageIdle).Compile(), "1"),
			map[string]any{"totalPages": numberedTotals},
		),
		"currentPage",
		map[string]string{
			"1": clPageBtn.Merge(clPageCur).Compile(),
			"2": clPageBtn.Merge(clPageIdle).Compile(),
		},
	)
	pageTwo := deliveryClassBound(
		deliveryVisibleWhen(
			deliveryText("PageTwo", clPageBtn.Merge(clPageIdle).Compile(), "2"),
			map[string]any{"totalPages": numberedTotals},
		),
		"currentPage",
		map[string]string{
			"1": clPageBtn.Merge(clPageIdle).Compile(),
			"2": clPageBtn.Merge(clPageCur).Compile(),
		},
	)
	root := deliveryFrame(
		componentType,
		clPagination.Compile(),
		deliveryVisibleWhen(
			deliveryText("Previous", clPageBtn.Merge(clPageIdle).Compile(), "‹"),
			map[string]any{"currentPage": 2},
		),
		pageOne,
		pageTwo,
		deliveryVisibleWhen(
			deliveryText("PageThree", clPageBtn.Merge(clPageIdle).Compile(), "3"),
			map[string]any{"totalPages": 7},
		),
		deliveryVisibleWhen(
			deliveryText("Ellipsis", clBreadcrumbSep.Compile(), "…"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("LastPage", clPageBtn.Merge(clPageIdle).Compile(), "7", "totalPages"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("Next", clPageBtn.Merge(clPageIdle).Compile(), "›"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("CursorPrevious", clPageBtn.Merge(clPageIdle).Compile(), "← Previous"),
			map[string]any{"totalPages": 0},
		),
		deliveryVisibleWhen(
			deliveryText("CursorLabel", clPageLabel.Compile(), "Page 1"),
			map[string]any{"totalPages": 0},
		),
		deliveryVisibleWhen(
			deliveryText("CursorNext", clPageBtn.Merge(clPageIdle).Compile(), "Next →"),
			map[string]any{"totalPages": 0},
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible numbered or cursor pagination.",
		map[string]PropertyContract{
			"currentPage":     contentProperty("Current one-based page."),
			"totalPages":      contentProperty("Known total; zero selects cursor mode."),
			"perPage":         modifierProperty("Rows per page."),
			"siblings":        modifierProperty("Visible siblings around the current page."),
			"baseURL":         contentProperty("Numbered navigation base URL."),
			"cursorMode":      variantProperty([]string{"previous-next", "load-more"}, "previous-next", "Cursor navigation presentation."),
			"previousCursor":  contentProperty("Opaque backward continuation cursor."),
			"nextCursor":      contentProperty("Opaque forward continuation cursor."),
			"beforeParameter": contentProperty("Backward cursor query parameter."),
			"afterParameter":  contentProperty("Forward cursor query parameter."),
			"previousURL":     contentProperty("Opaque backward continuation URL."),
			"nextURL":         contentProperty("Opaque forward continuation URL."),
			"previousLabel":   contentProperty("Localized previous-window label."),
			"nextLabel":       contentProperty("Localized next-window label."),
			"loadMoreLabel":   contentProperty("Localized forward-only label."),
			"navigationLabel": contentProperty("Accessible navigation label."),
		},
		nil,
		deliveryDesign("numbered", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "numbered", map[string]any{
				"currentPage": 2, "totalPages": 7, "siblings": 1, "baseURL": "/workspaces",
			}),
			deliveryExample(componentType, "workspace-page-one", map[string]any{
				"currentPage": 1, "totalPages": 4, "baseURL": "/workspaces",
			}),
			deliveryExample(componentType, "cursor", map[string]any{
				"currentPage": 1, "totalPages": 0,
				"baseURL": "/workspaces?sort=recent", "nextCursor": "opaque",
			}),
		},
		Pagination,
	)
}

func searchBarDeliveryDefinition() DeliveryDefinition {
	componentType := "SearchBar"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clSearchWrap.Compile(),
			deliveryInstance("Icon", "Icon", "Search / Fg Tertiary / 20", ""),
			deliveryText("Input", clSearchInput.Compile(), "Search workspaces", "placeholder", "label"),
			deliveryHiddenWhen(
				deliveryInstance("Clear", "Button", "button-icon-only-ghost", ""),
				map[string]any{"showClear": false},
			),
			deliveryHiddenWhen(
				deliveryInstance("Shortcut", "Kbd", "shortcut", ""),
				map[string]any{"showShortcut": false},
			),
		),
		"field",
		"field",
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Search field shell with progressive-enhancement hooks.",
		map[string]PropertyContract{
			"searchURL":    contentProperty("Search endpoint."),
			"label":        contentProperty("Accessible control name."),
			"placeholder":  contentProperty("Empty-state prompt."),
			"name":         contentProperty("Submitted query field name."),
			"value":        contentProperty("Initial query value."),
			"instant":      stateProperty("Whether typing submits progressively."),
			"showClear":    modifierProperty("Whether a clear affordance is shown."),
			"showShortcut": modifierProperty("Whether a keyboard hint is shown."),
			"debounceMs":   modifierProperty("Typing debounce in milliseconds."),
			"minChars":     modifierProperty("Minimum query length exposed to enhancement controllers."),
			"clearLabel":   contentProperty("Localized clear-action label."),
			"shortcutKey":  contentProperty("Displayed keyboard shortcut."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"searchURL": "/workspaces/search", "label": "Search workspaces",
				"placeholder": "Search workspaces", "instant": true,
			}),
			deliveryExample(componentType, "with-shortcut", map[string]any{
				"searchURL": "/search", "label": "Search everything",
				"placeholder": "Search everything...", "instant": true, "showShortcut": true,
			}),
		},
		SearchBar,
	)
}

func tabsDeliveryDefinition() DeliveryDefinition {
	componentType := "Tabs"
	root := deliveryFrame(
		componentType,
		clTabsRoot.Merge(clTabsRootHorizontal).Compile(),
		deliveryFrame(
			"TabList",
			clTabsListBase.Merge(clTabsListHorizontal).Merge(clTabsListUnderlineHorizontal).Compile(),
			deliveryText("Overview", clTabsButtonBase.Merge(clTabsButtonUnderlineHorizontal).Merge(clTabsUnderlineActive).Compile(), "Overview"),
			deliveryText("Activity", clTabsButtonBase.Merge(clTabsButtonUnderlineHorizontal).Merge(clTabsUnderlineIdle).Compile(), "Activity"),
		),
		deliverySlot(
			"Panel",
			"tab",
			clTabsPanels.Compile(),
			deliveryInstance("OverviewBody", "Text", "body", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible tab navigation and controller-backed panel composition.",
		map[string]PropertyContract{
			"items":        contentProperty("Ordered navigation-only tabs."),
			"activeTab":    stateProperty("Active tab key."),
			"orientation":  variantProperty([]string{"horizontal", "vertical"}, "horizontal", "Tab direction."),
			"variant":      variantProperty([]string{"underline", "pills"}, "underline", "Visual treatment."),
			"hxGet":        contentProperty("Default lazy-panel endpoint."),
			"loadingLabel": contentProperty("Localized lazy-panel loading label."),
		},
		[]designcomponent.Slot{{
			Name:         "tab",
			Required:     true,
			Description:  "Ordered tab panels with per-tab navigation metadata.",
			AllowedTypes: deliveryMoleculeContentTypes(),
			Cardinality:  designcomponent.SlotMany,
			Attrs:        tabsDeliverySlotAttrs(),
		}},
		deliveryDesign("underline", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "underline", map[string]any{
					"activeTab": "overview", "orientation": "horizontal", "variant": "underline",
				}),
				deliveryExampleSlot("tab",
					tabsDeliveryExampleComponent("overview", "Overview", "Current workspace summary."),
					tabsDeliveryExampleComponent("activity", "Activity", "Recent workspace activity."),
					tabsDeliveryExampleComponent("settings", "Settings", "Workspace preferences."),
				),
			),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "vertical-pills", map[string]any{
					"activeTab": "profile", "orientation": "vertical", "variant": "pills",
				}),
				deliveryExampleSlot("tab",
					tabsDeliveryExampleComponent("profile", "Profile", "Update profile information."),
					tabsDeliveryExampleComponent("security", "Security", "Review security settings."),
				),
			),
		},
		func(props molecules.TabsProps, slots DeliverySlotChildren) g.Node {
			tabs := make([]TabSlot, 0, len(slots["tab"]))
			for _, instance := range slots["tab"] {
				tabs = append(tabs, TabSlot{
					ID:       deliveryString(instance.Attrs, "id"),
					Label:    deliveryString(instance.Attrs, "label"),
					Icon:     deliveryString(instance.Attrs, "icon"),
					Badge:    deliveryString(instance.Attrs, "badge"),
					Disabled: deliveryBool(instance.Attrs, "disabled"),
					HxGet:    deliveryString(instance.Attrs, "hxGet"),
					Content:  instance.Children,
				})
			}
			return TabsWithPanels(props, tabs...)
		},
	)
}

func tabsDeliverySlotAttrs() []designcomponent.Prop {
	return []designcomponent.Prop{
		{Name: "id", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Required: true, Description: "Stable tab identity."},
		{Name: "label", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Required: true, Description: "Visible tab label."},
		{Name: "icon", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Description: "Optional OSS icon name."},
		{Name: "badge", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Description: "Optional status or count badge."},
		{Name: "disabled", Type: designcomponent.PropBoolean, Role: designcomponent.PropRoleState, Default: "false", Description: "Whether the tab cannot be activated."},
		{Name: "hxGet", Type: designcomponent.PropString, Role: designcomponent.PropRoleContent, Description: "Optional lazy-panel endpoint."},
	}
}

func tabsDeliveryExampleComponent(id, label, content string) DeliveryExampleComponent {
	component := deliveryExampleComponent(
		id+"-body",
		"Text",
		map[string]any{"content": content, "size": "sm", "color": "muted"},
	)
	component.SlotAttrs = map[string]any{"id": id, "label": label}
	return component
}

func deliveryString(attrs map[string]any, name string) string {
	value, _ := attrs[name].(string)
	return value
}

func deliveryBool(attrs map[string]any, name string) bool {
	value, _ := attrs[name].(bool)
	return value
}

func dataGridDeliveryDefinition() DeliveryDefinition {
	componentType := "DataGrid"
	root := deliveryFrame(
		componentType,
		clGridSection.Compile(),
		deliveryFrame(
			"Toolbar",
			clGridToolbar.Compile(),
			deliverySlot(
				"Search",
				"search",
				"flex-1",
				deliveryInstance("SearchFixture", "SearchBar", "basic", ""),
			),
			deliverySlot("Filters", "filters", ""),
			deliverySlot(
				"Actions",
				"actions",
				clGridActions.Compile(),
				deliveryInstance("CreateFixture", "Link", "default", ""),
			),
		),
		deliverySlot(
			"Table",
			"table",
			"",
			deliveryInstance("TableFixture", "Table", "data", ""),
		),
		deliverySlot(
			"Status",
			"status",
			"",
			deliveryInstance("Empty", "EmptyState", "compact", ""),
		),
		deliverySlot(
			"Pagination",
			"pagination",
			"",
			deliveryInstance("PaginationFixture", "Pagination", "numbered", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Search, filter, action, table, status, and pagination composition.",
		map[string]PropertyContract{},
		[]designcomponent.Slot{
			{
				Name:         "search",
				Description:  "Optional search control.",
				AllowedTypes: []string{"SearchBar"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "filters",
				Description:  "Ordered filter controls.",
				AllowedTypes: []string{"Input", "Select"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "actions",
				Description:  "Toolbar actions.",
				AllowedTypes: []string{"Button", "Link"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "table",
				Description:  "Primary tabular result view.",
				AllowedTypes: []string{"Table"},
				Required:     true,
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "status",
				Description:  "Loading, error, or empty status between data and pagination.",
				AllowedTypes: []string{"Alert", "EmptyState", "Spinner", "Text"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "pagination",
				Description:  "Optional pagination control.",
				AllowedTypes: []string{"Pagination"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("workspace-list", "04 Organisms", "Data", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "workspace-list", nil),
				dataGridDeliveryExampleSlots()...,
			),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "server-side-grid", map[string]any{
					"hx-get":     "/workspaces",
					"hx-target":  "#workspace-grid",
					"hx-trigger": "workspace:changed from:body",
				}),
				dataGridDeliveryExampleSlots()...,
			),
		},
		func(props organisms.DataGridProps, slots DeliverySlotChildren) g.Node {
			return DataGrid(props, DataGridSlots{
				Search:     firstDeliverySlotNode(slots, "search"),
				Filters:    slots.Nodes("filters"),
				Actions:    slots.Nodes("actions"),
				Table:      firstDeliverySlotNode(slots, "table"),
				Status:     slots.Nodes("status"),
				Pagination: firstDeliverySlotNode(slots, "pagination"),
			})
		},
	)
}

func dataGridDeliveryExampleSlots() []DeliveryExampleSlot {
	return []DeliveryExampleSlot{
		deliveryExampleSlot(
			"search",
			deliveryExampleComponent(
				"workspace-search",
				"SearchBar",
				map[string]any{
					"searchURL":   "/workspaces/search",
					"label":       "Search workspaces",
					"placeholder": "Search workspaces",
					"instant":     true,
				},
			),
		),
		deliveryExampleSlot(
			"actions",
			deliveryExampleComponent(
				"create-workspace",
				"Link",
				map[string]any{
					"label": "New workspace",
					"href":  "/workspaces/new",
				},
			),
		),
		deliveryExampleSlot(
			"table",
			deliveryExampleComponent(
				"workspace-table",
				"Table",
				tableExampleProps(),
			),
		),
		deliveryExampleSlot(
			"pagination",
			deliveryExampleComponent(
				"workspace-pagination",
				"Pagination",
				map[string]any{
					"currentPage": 1,
					"totalPages":  4,
					"baseURL":     "/workspaces",
				},
			),
		),
	}
}

func dashboardWidgetDeliveryDefinition() DeliveryDefinition {
	componentType := "DashboardWidget"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clDashboardWidget.Merge(clDashboardWidgetSpan["1"]).Compile(),
			deliveryFrame(
				"Header",
				clDashboardWidgetHeader.Compile(),
				deliveryText("Title", clDashboardWidgetTitle.Compile(), "Active seats", "title"),
				deliveryText("Subtitle", clDashboardWidgetSubtitle.Compile(), "Current workspace posture", "subtitle"),
			),
			deliveryFrame(
				"Body",
				clDashboardWidgetBody.Compile(),
				deliveryText("Value", clDashboardWidgetValue.Compile(), "1,234", "value"),
				deliveryText("Change", clDashboardWidgetChange.Merge(clDashboardWidgetTrendTone["up"]).Compile(), "+12%", "change"),
			),
			deliverySlot("Content", "content", clDashboardWidgetBody.Compile()),
			deliverySlot("Footer", "footer", clDashboardWidgetFooter.Compile()),
		),
		"span",
		compileClassMap(clDashboardWidgetSpan),
	)
	definition := newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Dashboard metric or rich-content card with governed refresh behavior and explicit content/footer composition.",
		map[string]PropertyContract{
			"title":          contentProperty("Visible widget title."),
			"subtitle":       contentProperty("Optional supporting title copy."),
			"type":           variantProperty([]string{"stat", "chart", "list", "empty"}, "stat", "Widget content mode."),
			"value":          contentProperty("Primary stat value."),
			"previousValue":  contentProperty("Previous stat value for comparison."),
			"change":         contentProperty("Localized trend change label."),
			"trend":          variantProperty([]string{"up", "down", "flat"}, "flat", "Trend direction."),
			"icon":           roleProperty(designcomponent.PropRoleSlot, "Optional governed icon name."),
			"refreshURL":     contentProperty("Optional HTMX refresh endpoint."),
			"refreshOn":      contentProperty("Optional body event that triggers refresh."),
			"refreshSeconds": modifierProperty("Optional bounded polling interval in seconds."),
			"detailURL":      contentProperty("Optional destination linked from the title."),
			"span":           variantProperty([]string{"1", "2", "3", "4"}, "1", "Grid column span."),
		},
		[]designcomponent.Slot{
			{
				Name:         "content",
				Description:  "Optional rich body replacing the scalar stat body.",
				AllowedTypes: deliveryMoleculeContentTypes(),
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "footer",
				Description:  "Optional widget actions or supporting status.",
				AllowedTypes: deliveryAtomicContentTypes(),
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
		},
		deliveryDesign("active-seats", "04 Organisms", "Operations & Admin", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "active-seats", map[string]any{
				"title": "Active seats", "subtitle": "Current workspace posture",
				"type": "stat", "value": "1,234", "previousValue": "1,102",
				"change": "+12%", "trend": "up", "icon": "activity", "span": "1",
			}),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "activity-list", map[string]any{
					"title": "Recent activity", "type": "list", "span": "2",
					"refreshURL": "/dashboard/activity", "refreshOn": "activity:changed",
				}),
				deliveryExampleSlot("content", deliveryExampleComponent(
					"activity-copy", "Text", map[string]any{
						"content": "Three workspace events need review.", "size": "sm", "color": "muted",
					},
				)),
				deliveryExampleSlot("footer", deliveryExampleComponent(
					"activity-link", "Link", map[string]any{"label": "View activity", "href": "/activity"},
				)),
			),
		},
		func(props organisms.DashboardWidgetProps, slots DeliverySlotChildren) g.Node {
			return DashboardWidget(props, DashboardWidgetSlots{
				Content: slots.Nodes("content"),
				Footer:  slots.Nodes("footer"),
			})
		},
	)
	return withDeliveryTags(definition, "dashboard", "widget", "metrics", "admin")
}

func windowedCollectionDeliveryDefinition() DeliveryDefinition {
	componentType := "WindowedCollection"
	items := deliveryClassBound(
		deliverySlot(
			"Items",
			"items",
			clWindowItems.Compile(),
			deliveryInstance("FirstItem", "Card", "summary", ""),
			deliveryInstance("SecondItem", "Card", "summary", ""),
		),
		"layout",
		map[string]string{
			"list": "",
			"grid": clWindowItemsGrid.Compile(),
		},
	)
	root := deliveryFrame(
		componentType,
		clWindow.Compile(),
		deliveryFrame(
			"Header",
			clWindowHeader.Compile(),
			deliveryFrame("Accent", clWindowAccent.Compile()),
			deliveryText("Title", clWindowHeading.Compile(), "Recent work", "title"),
			deliveryText(
				"Description",
				clWindowIntro.Compile(),
				"A bounded window of recent work.",
				"description",
			),
		),
		items,
		deliveryVisibleWhen(
			deliveryText(
				"Loading",
				clWindowLoading.Compile(),
				"Loading results",
				"loadingLabel",
			),
			map[string]any{"state": "loading"},
		),
		deliveryVisibleWhen(
			deliveryText(
				"Error",
				clWindowError.Compile(),
				"Unable to load results",
				"errorTitle",
			),
			map[string]any{"state": "error"},
		),
		deliveryVisibleWhen(
			deliveryText(
				"Empty",
				clWindowEmpty.Compile(),
				"No results",
				"emptyTitle",
			),
			map[string]any{"itemCount": 0},
		),
		deliverySlot(
			"Controls",
			"controls",
			clWindowFooter.Compile(),
			deliveryInstance("Pagination", "Pagination", "cursor", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Bounded collection shell with explicit item, state, and cursor-control composition.",
		map[string]PropertyContract{
			"layout": variantProperty(
				[]string{"list", "grid"},
				"list",
				"Semantic item arrangement shared by adaptive, web, and native renderers.",
			),
			"state": variantProperty(
				[]string{"ready", "loading", "error"},
				"ready",
				"Current collection state.",
			),
			"itemCount":             contentProperty("Number of bounded item roots."),
			"maxItems":              modifierProperty("Hard render-node budget."),
			"collectionLabel":       contentProperty("Accessible collection label."),
			"title":                 contentProperty("Optional visible collection heading."),
			"description":           contentProperty("Optional supporting collection copy."),
			"emptyTitle":            contentProperty("Localized empty-state title."),
			"emptyDescription":      contentProperty("Localized empty-state description."),
			"loadingLabel":          contentProperty("Localized loading status."),
			"errorTitle":            contentProperty("Localized error title."),
			"errorDescription":      contentProperty("Localized error description."),
			"retryLabel":            contentProperty("Localized retry label."),
			"retryURL":              contentProperty("Progressive retry URL."),
			"navigationUnavailable": contentProperty("Cursor-navigation failure message."),
		},
		[]designcomponent.Slot{
			{
				Name:         "items",
				Description:  "Bounded item roots in stable source order.",
				AllowedTypes: []string{"Card", "DataGrid", "Table", "Text"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "controls",
				Description:  "Optional bounded cursor navigation.",
				AllowedTypes: []string{"Pagination"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("cards", "04 Organisms", "Data", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "cards", map[string]any{
					"layout": "grid", "state": "ready", "itemCount": 2, "maxItems": 100,
					"title": "Recent work", "description": "A bounded window of recent work.",
				}),
				deliveryExampleSlot(
					"items",
					deliveryExampleComponent(
						"window-card-one",
						"Card",
						map[string]any{"title": "First result"},
					),
					deliveryExampleComponent(
						"window-card-two",
						"Card",
						map[string]any{"title": "Second result"},
					),
				),
				deliveryExampleSlot(
					"controls",
					// A nested instance is only natively representable when its
					// props match the authored variant it targets; the cursor
					// variant expresses forward navigation as baseURL+nextCursor,
					// and nextURL has no native representation in the design.
					deliveryExampleComponent(
						"window-pagination",
						"Pagination",
						map[string]any{
							"currentPage": 1,
							"totalPages":  0,
							"baseURL":     "/workspaces?sort=recent",
							"nextCursor":  "opaque",
						},
					),
				),
			),
			deliveryExample(componentType, "empty", map[string]any{
				"state": "ready", "itemCount": 0, "maxItems": 100,
			}),
		},
		func(
			props organisms.WindowedCollectionProps,
			slots DeliverySlotChildren,
		) g.Node {
			return WindowedCollection(props, WindowedCollectionSlots{
				Items:    slots.Nodes("items"),
				Controls: firstDeliverySlotNode(slots, "controls"),
			})
		},
	)
}

func stackDeliveryDefinition() DeliveryDefinition {
	componentType := "Stack"
	root := deliverySlot(
		componentType,
		"children",
		clStack.Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Card", "summary", ""),
		deliveryInstance("Secondary", "Card", "mobile-summary", ""),
	)
	root = deliveryClassBound(root, "gap", deliveryGapClassMap())
	root = deliveryClassBound(root, "align", deliveryAlignClassMap())
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Vertical auto-layout primitive.",
		map[string]PropertyContract{
			"gap":   sizeProperty(gaps, "4", "Space between children."),
			"align": variantProperty([]string{"start", "center", "end", "stretch"}, "stretch", "Cross-axis alignment."),
		},
		deliveryDesign("default", "05 Templates", "Layout", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"gap": "4", "align": "stretch",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"primary-card",
						"Quarterly summary",
						"Performance is tracking above plan.",
					),
					deliveryExampleCard(
						"secondary-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "compact", map[string]any{
					"gap": "2", "align": "stretch",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"compact-heading",
						"Heading",
						map[string]any{"level": 2, "text": "Compact stack"},
					),
					deliveryExampleComponent(
						"compact-copy",
						"Text",
						map[string]any{
							"color":   "muted",
							"content": "Closely related supporting content.",
							"size":    "base",
						},
					),
				),
			),
		},
		Stack,
	)
}

func deliveryGapClassMap() map[string]string {
	out := make(map[string]string, len(clGapScale))
	for name, gap := range clGapScale {
		out[name] = tw.New().Gap(gap).Compile()
	}
	return out
}

func deliveryAlignClassMap() map[string]string {
	out := make(map[string]string, len(alignItems))
	for name, align := range alignItems {
		out[name] = tw.New().Items(align).Compile()
	}
	return out
}

func flexDeliveryDefinition() DeliveryDefinition {
	componentType := "Flex"
	root := deliverySlot(
		componentType,
		"children",
		clFlex.Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Button", "primary", ""),
		deliveryInstance("Secondary", "Button", "secondary", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"One-dimensional flexible layout primitive.",
		map[string]PropertyContract{
			"direction": variantProperty([]string{"row", "column"}, "row", "Main-axis direction."),
			"wrap":      modifierProperty("Whether children may wrap."),
			"gap":       sizeProperty(gaps, "4", "Space between children."),
			"align":     variantProperty([]string{"start", "center", "end", "stretch"}, "center", "Cross-axis alignment."),
			"justify":   variantProperty([]string{"start", "center", "end", "between"}, "start", "Main-axis distribution."),
		},
		deliveryDesign("row", "05 Templates", "Layout", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "row", map[string]any{
					"direction": "row", "gap": "4", "align": "center", "justify": "start",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"primary-action",
						"Button",
						map[string]any{
							"label": "Continue", "variant": "primary",
							"tone": "neutral", "size": "md",
						},
					),
					deliveryExampleComponent(
						"secondary-action",
						"Button",
						map[string]any{
							"label": "Cancel", "variant": "secondary",
							"tone": "neutral", "size": "md",
						},
					),
				),
			),
		},
		Flex,
	)
}

func gridDeliveryDefinition() DeliveryDefinition {
	componentType := "Grid"
	root := deliverySlot(
		componentType,
		"children",
		clGrid.GridCols(2).Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Card", "summary", ""),
		deliveryInstance("Secondary", "Card", "mobile-summary", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Two-dimensional responsive grid primitive.",
		map[string]PropertyContract{
			"columns": variantProperty([]string{"1", "2", "3", "4", "6", "12"}, "1", "Column count."),
			"gap":     sizeProperty(gaps, "4", "Space between cells."),
		},
		deliveryDesign("two-columns", "05 Templates", "Layout", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "two-columns", map[string]any{
					"columns": "2", "gap": "4",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"primary-card",
						"Quarterly summary",
						"Performance is tracking above plan.",
					),
					deliveryExampleCard(
						"secondary-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
			withDeliveryExampleSlots(
				mobileDeliveryExample(componentType, "single-column", map[string]any{
					"columns": "1", "gap": "3",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"mobile-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
		},
		Grid,
	)
}

func containerDeliveryDefinition() DeliveryDefinition {
	componentType := "Container"
	root := deliverySlot(
		componentType,
		"children",
		clContainer.MaxWScaled(tw.MaxW7XL).Compile(),
		deliveryInstance("Content", "Stack", "default", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Centered maximum-width page container.",
		map[string]PropertyContract{
			"maxWidth": sizeProperty([]string{"sm", "md", "lg", "xl", "2xl", "4xl", "7xl", "full"}, "7xl", "Maximum content width."),
			"padding":  sizeProperty(gaps, "4", "Horizontal padding."),
		},
		deliveryDesign("default", "05 Templates", "Layout", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"maxWidth": "7xl", "padding": "4",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"content-stack",
						"Stack",
						map[string]any{"gap": "2", "align": "stretch"},
						deliveryExampleSlot(
							"children",
							deliveryExampleComponent(
								"content-heading",
								"Heading",
								map[string]any{
									"text":  "Workspace overview",
									"level": 2,
								},
							),
							deliveryExampleComponent(
								"content-copy",
								"Text",
								map[string]any{
									"content": "Review current operational details.",
									"size":    "base", "color": "muted",
								},
							),
						),
					),
				),
			),
		},
		Container,
	)
}

func deliveryExampleCard(
	id string,
	title string,
	description string,
) DeliveryExampleComponent {
	return deliveryExampleComponent(
		id,
		"Card",
		map[string]any{
			"title":       title,
			"description": description,
		},
		deliveryExampleSlot(
			"children",
			deliveryExampleComponent(
				id+"-body",
				"Text",
				map[string]any{
					"content": description,
					"size":    "sm",
					"color":   "muted",
				},
			),
		),
	)
}

// DataManagementPageProps is the typed contract for the canonical OSS
// solution page. It deliberately models product-neutral data management so
// proprietary modules and clients can extend the composition without forking
// its atoms, molecules, organism, or layout primitives.
type DataManagementPageProps struct {
	contracts.ComponentProps

	Title       string               `json:"title"`
	Description string               `json:"description,omitempty"`
	Rows        []molecules.TableRow `json:"rows,omitempty"`
}

func dataManagementPageDeliveryDefinition() DeliveryDefinition {
	componentType := "DataManagementPage"
	root := deliveryFrame(
		componentType,
		clDataManagementPage.Compile(),
		deliveryFrame(
			"Container",
			clContainer.MaxWScaled(tw.MaxW7XL).Compile(),
			deliveryFrame(
				"Stack",
				clStack.Gap(tw.S6).Compile(),
				deliveryInstance("Breadcrumb", "Breadcrumb", "default", ""),
				deliveryFrame(
					"Header",
					clStack.Gap(tw.S2).Compile(),
					deliveryInstance("Title", "Heading", "level-one", ""),
					deliveryInstance("Description", "Text", "muted", ""),
				),
				deliveryInstance("Tabs", "Tabs", "underline", ""),
				deliveryInstance("Results", "DataGrid", "workspace-list", ""),
			),
		),
	)
	design := deliveryDesign(
		"workspace-management",
		"06 OSS Solutions",
		"Data Management",
		root,
	)
	design.Taxonomy.FlowGroup = "data-management"
	design.Taxonomy.FlowOrder = 10
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierPage,
		stableDeliveryID(componentType),
		"Complete OSS data-management solution composed exclusively from shared catalog building blocks.",
		map[string]PropertyContract{
			"title":       contentProperty("Page heading."),
			"description": contentProperty("Page purpose."),
			"rows":        contentProperty("Rows displayed by the shared data grid."),
		},
		nil,
		design,
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "workspace-management", map[string]any{
				"title":       "Workspaces",
				"description": "Manage access, lifecycle, and operational state.",
				"rows": []map[string]any{
					{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
					{"id": "bravo", "cells": map[string]any{"name": "Bravo workspace", "status": "Review", "updated": "Yesterday"}},
					{"id": "charlie", "cells": map[string]any{"name": "Charlie workspace", "status": "Paused", "updated": "3 days ago"}},
				},
			}),
			mobileDeliveryExample(componentType, "workspace-management-mobile", map[string]any{
				"title":       "Workspaces",
				"description": "Manage access and lifecycle.",
				"rows": []map[string]any{
					{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
				},
			}),
		},
		renderDataManagementPage,
	)
}

func renderDataManagementPage(props DataManagementPageProps) g.Node {
	rows := props.Rows
	if len(rows) == 0 {
		rows = []molecules.TableRow{{
			ID: "alpha",
			Cells: map[string]any{
				"name": "Alpha workspace", "status": "Active", "updated": "Today",
			},
		}}
	}
	title := props.Title
	if title == "" {
		title = "Workspaces"
	}
	description := props.Description
	if description == "" {
		description = "Manage access, lifecycle, and operational state."
	}
	main := baseAttrs(props.ComponentProps)
	main = append(main,
		classes(clDataManagementPage.Compile(), props.Class),
		Container(
			layouts.ContainerProps{MaxWidth: "7xl"},
			Stack(
				layouts.StackProps{Gap: "6"},
				Breadcrumb(molecules.BreadcrumbProps{
					Items: []molecules.BreadcrumbItem{
						{Label: "Operations", Href: "/operations"},
						{Label: title, Current: true},
					},
				}),
				Stack(
					layouts.StackProps{Gap: "2"},
					Heading(atoms.HeadingProps{Text: title, Level: 1}),
					Text(atoms.TextProps{Content: description, Size: "sm", Color: "muted"}),
				),
				Tabs(molecules.TabsProps{
					ActiveTab: "all",
					Items: []molecules.TabItem{
						{Key: "all", Label: "All", URL: "/workspaces"},
						{Key: "active", Label: "Active", URL: "/workspaces?status=active"},
						{Key: "archived", Label: "Archived", URL: "/workspaces?status=archived"},
					},
				}),
				DataGrid(
					organisms.DataGridProps{},
					dataGridSlots(rows),
				),
			),
		),
	)
	return h.Main(main...)
}

func tableExampleProps() map[string]any {
	return map[string]any{
		"columns": []map[string]any{
			{"key": "name", "label": "Name", "sortable": true, "primary": true},
			{"key": "status", "label": "Status", "sortable": true},
			{"key": "updated", "label": "Updated"},
		},
		"rows": []map[string]any{
			{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
			{"id": "bravo", "cells": map[string]any{"name": "Bravo workspace", "status": "Review", "updated": "Yesterday"}},
		},
		"sortable": true,
		"striped":  true,
	}
}

func dataGridSlots(rows []molecules.TableRow) DataGridSlots {
	return DataGridSlots{
		Search: SearchBar(molecules.SearchBarProps{
			SearchURL:   "/workspaces/search",
			Label:       "Search workspaces",
			Placeholder: "Search workspaces",
			Instant:     true,
		}),
		Actions: []g.Node{
			Link(atoms.LinkProps{
				Label: "New workspace",
				Href:  "/workspaces/new",
			}),
		},
		Table: Table(molecules.TableProps{
			Columns: []molecules.TableColumn{
				{Key: "name", Label: "Name", Sortable: true, Primary: true},
				{Key: "status", Label: "Status", Sortable: true},
				{Key: "updated", Label: "Updated"},
			},
			Rows:     rows,
			Sortable: true,
			Striped:  true,
		}),
		Pagination: Pagination(molecules.PaginationProps{
			CurrentPage: 1,
			TotalPages:  4,
			BaseURL:     "/workspaces",
		}),
	}
}

func tableSkeletonDeliveryDefinition() DeliveryDefinition {
	componentType := "TableSkeleton"
	line := skeletonLine("sm", false).Compile()
	root := deliveryFrame(
		componentType,
		clTableWrap.Compile(),
		deliveryFrame(
			"Table",
			clTable.Compile(),
			deliveryFrame(
				"Header",
				clTableHead.Compile(),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
			),
			deliveryFrame(
				"Row",
				clTableRow.Compile(),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
			),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Loading rendering of Table: identical wrap, header, and cell classes with pulsing placeholders.",
		map[string]PropertyContract{
			"columns": contentProperty("Placeholder column count."),
			"rows":    contentProperty("Placeholder row count."),
			"compact": modifierProperty("Whether row density is compact."),
		},
		nil,
		deliveryDesign("loading", "03 Molecules", "Data", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "loading", map[string]any{
				"columns": 3, "rows": 3,
			}),
		},
		TableSkeleton,
	)
}

func cardSkeletonDeliveryDefinition() DeliveryDefinition {
	componentType := "CardSkeleton"
	root := deliveryFrame(
		componentType,
		clCard.Compile(),
		deliveryFrame("Title", skeletonLine("lg", true).Compile()),
		deliveryFrame(
			"Body",
			clSkeletonText.Compile(),
			deliveryFrame("Line", skeletonLine("md", false).Compile()),
			deliveryFrame("Line", skeletonLine("md", false).Compile()),
			deliveryFrame("Line", skeletonLine("md", true).Compile()),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Loading rendering of Card: title and body placeholders inside the real card frame.",
		map[string]PropertyContract{
			"lines": contentProperty("Body placeholder line count."),
		},
		nil,
		deliveryDesign("loading", "03 Molecules", "Surfaces", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "loading", map[string]any{
				"lines": 3,
			}),
		},
		CardSkeleton,
	)
}
