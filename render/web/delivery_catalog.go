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

	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
	"github.com/septagon-oss/tw"
)

var (
	buttonVariants = []string{
		"primary", "secondary", "success", "warning", "error", "info",
		"outline", "ghost", "link",
	}
	controlSizes = []string{"xs", "small", "medium", "large", "xl", "2xl"}
	gaps         = []string{"0", "1", "2", "3", "4", "5", "6", "8"}
)

// OSSDeliveryCatalog returns the complete renderer-backed OSS web library in
// deterministic atomic-design order. Returned definitions are defensive
// copies; callers can safely adapt them into isolated runtimes.
func OSSDeliveryCatalog() []DeliveryDefinition {
	definitions := []DeliveryDefinition{
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
		emptyStateDeliveryDefinition(),
		kbdDeliveryDefinition(),
		linkDeliveryDefinition(),
		tagDeliveryDefinition(),
		tableDeliveryDefinition(),
		cardDeliveryDefinition(),
		breadcrumbDeliveryDefinition(),
		paginationDeliveryDefinition(),
		searchBarDeliveryDefinition(),
		tabsDeliveryDefinition(),
		dataGridDeliveryDefinition(),
		stackDeliveryDefinition(),
		flexDeliveryDefinition(),
		gridDeliveryDefinition(),
		containerDeliveryDefinition(),
		dataManagementPageDeliveryDefinition(),
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

func buttonDeliveryDefinition() DeliveryDefinition {
	componentType := "Button"
	root := deliveryClassBound(
		deliveryClassBound(
			deliveryFrame(
				componentType,
				clButtonBase.Compile(),
				deliveryText("Label", "", "Continue", "text"),
			),
			"variant",
			compileClassMap(clButtonVariant),
		),
		"size",
		compileClassMap(clButtonSize),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible action control rendered by the OSS web runtime.",
		map[string]PropertyContract{
			"text":         contentProperty("Visible button label."),
			"variant":      variantProperty(buttonVariants, "primary", "Visual treatment."),
			"size":         sizeProperty(controlSizes, "medium", "Control size."),
			"type":         enumProperty([]string{"button", "submit", "reset"}, "button", "HTML button type."),
			"loading":      stateProperty("Whether progress replaces the leading content."),
			"fullWidth":    modifierProperty("Whether the control fills its container."),
			"icon":         contentProperty("Optional icon identity."),
			"iconPosition": enumProperty([]string{"left", "right"}, "left", "Icon placement."),
		},
		nil,
		deliveryDesign("primary", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "primary", map[string]any{
				"text": "Continue", "variant": "primary", "size": "medium",
			}),
			deliveryExample(componentType, "secondary", map[string]any{
				"text": "Cancel", "variant": "secondary", "size": "medium",
			}),
			deliveryExample(componentType, "loading", map[string]any{
				"text": "Saving", "variant": "primary", "size": "medium", "loading": true,
			}),
		},
		Button,
	)
}

func badgeDeliveryDefinition() DeliveryDefinition {
	componentType := "Badge"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clBadgeBase.Compile(),
			deliveryText("Label", "", "Active", "text"),
		),
		"variant",
		compileClassMap(clBadgeVariant),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Compact semantic status label.",
		map[string]PropertyContract{
			"text": contentProperty("Visible badge text."),
			"variant": toneProperty(
				[]string{"default", "primary", "secondary", "success", "warning", "error", "info"},
				"default",
				"Semantic status treatment.",
			),
			"size": enumProperty([]string{"small", "medium", "large"}, "medium", "Badge scale."),
			"dot":  modifierProperty("Shows a leading status dot."),
			"icon": contentProperty("Optional icon identity."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Status", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"text": "Active", "variant": "default",
			}),
			deliveryExample(componentType, "success", map[string]any{
				"text": "Published", "variant": "success", "dot": true,
			}),
		},
		Badge,
	)
}

func alertDeliveryDefinition() DeliveryDefinition {
	componentType := "Alert"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clAlertBase.Compile(),
			deliveryFrame(
				"Body",
				clAlertBody.Compile(),
				deliveryText("Title", clAlertTitle.Compile(), "Update available", "title"),
				deliveryText("Message", clAlertMessage.Compile(), "Review the latest changes.", "message"),
			),
		),
		"variant",
		compileClassMap(clAlertVariant),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Persistent accessible inline status message.",
		map[string]PropertyContract{
			"message":     contentProperty("Required status message."),
			"title":       contentProperty("Optional status title."),
			"variant":     toneProperty([]string{"info", "success", "warning", "error"}, "info", "Severity."),
			"icon":        contentProperty("Optional icon identity."),
			"dismissible": stateProperty("Whether a dismiss affordance is exposed."),
			"bordered":    modifierProperty("Whether the message has an outline."),
			"compact":     modifierProperty("Whether spacing is compact."),
		},
		nil,
		deliveryDesign("info", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "info", map[string]any{
				"title": "Update available", "message": "Review the latest changes.", "variant": "info",
			}),
			deliveryExample(componentType, "warning", map[string]any{
				"title": "Check required", "message": "One field needs attention.", "variant": "warning",
			}),
		},
		Alert,
	)
}

func inputDeliveryDefinition() DeliveryDefinition {
	componentType := "Input"
	root := deliveryFrame(
		componentType,
		clFieldWrap.Compile(),
		deliveryText("Label", clLabel.Compile(), "Email address", "label"),
		deliveryClassBound(
			deliveryFrame(
				"Control",
				clInput.Compile()+" "+clInputNormal.Compile(),
				deliveryText("Value", "", "name@example.com", "value", "placeholder"),
			),
			"error",
			map[string]string{"": clInputNormal.Compile(), "*": clInputError.Compile()},
		),
		deliveryText("Help", clHelp.Compile(), "We never share your email.", "helpText", "error"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Labelled single-line form control with help and error states.",
		map[string]PropertyContract{
			"name":        contentProperty("Submitted field name."),
			"type":        variantProperty([]string{"text", "email", "password", "number", "tel", "url", "search", "date", "time"}, "text", "Input kind."),
			"value":       contentProperty("Current value."),
			"placeholder": contentProperty("Empty-state prompt."),
			"label":       contentProperty("Visible field label."),
			"helpText":    contentProperty("Supporting guidance."),
			"error":       stateProperty("Validation message; non-empty selects error styling."),
			"required":    stateProperty("Whether a value is required."),
			"readOnly":    stateProperty("Whether editing is prevented."),
			"autoFocus":   stateProperty("Whether focus is requested on mount."),
			"size":        sizeProperty([]string{"small", "medium", "large"}, "medium", "Control size."),
		},
		nil,
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
		},
		Input,
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
		"Native single-value choice control.",
		map[string]PropertyContract{
			"name":        contentProperty("Submitted field name."),
			"label":       contentProperty("Visible field label."),
			"value":       contentProperty("Selected option value."),
			"placeholder": contentProperty("Unselected prompt."),
			"options":     contentProperty("Available choices."),
			"required":    stateProperty("Whether a choice is required."),
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
		},
		Select,
	)
}

func textareaDeliveryDefinition() DeliveryDefinition {
	componentType := "Textarea"
	root := deliveryFrame(
		componentType,
		clFieldWrap.Compile(),
		deliveryText("Label", clLabel.Compile(), "Notes", "label"),
		deliveryFrame(
			"Control",
			clInput.Compile()+" "+clInputNormal.Compile(),
			deliveryText("Value", "", "Add context for reviewers.", "value", "placeholder"),
		),
		deliveryText("Help", clHelp.Compile(), "Keep it concise.", "helperText", "errorMessage"),
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
			"fullWidth":    modifierProperty("Whether the field fills its container."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"name": "notes", "label": "Notes", "placeholder": "Add context for reviewers.",
				"helperText": "Keep it concise.", "rows": 4,
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
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Semantic body-copy primitive.",
		map[string]PropertyContract{
			"content":  contentProperty("Visible copy."),
			"size":     sizeProperty([]string{"xs", "sm", "base", "md", "lg", "xl", "2xl"}, "base", "Text scale."),
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
		},
		Text,
	)
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
		deliveryFrame(componentType, clSpinner.Compile()),
		"size",
		compileClassMap(clSpinnerSize),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible indeterminate progress indicator.",
		map[string]PropertyContract{
			"label": contentProperty("Screen-reader status label."),
			"size":  sizeProperty([]string{"small", "medium", "large"}, "medium", "Indicator size."),
		},
		nil,
		deliveryDesign("medium", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "medium", map[string]any{
				"label": "Loading results", "size": "medium",
			}),
		},
		Spinner,
	)
}

func emptyStateDeliveryDefinition() DeliveryDefinition {
	componentType := "EmptyState"
	root := deliveryFrame(
		componentType,
		clEmpty.Merge(clEmptyPad).Compile(),
		deliveryText("Title", clEmptyTitle.Compile(), "No projects yet", "title"),
		deliveryText("Description", clEmptyDesc.Compile(), "Create a project to begin.", "description"),
		deliverySlot(
			"Actions",
			"actions",
			clCheckRow.Compile(),
			deliveryInstance("PrimaryAction", "Button", "primary", ""),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Empty-data explanation with optional actions.",
		map[string]PropertyContract{
			"icon":        contentProperty("Optional illustration identity."),
			"title":       contentProperty("Primary empty-state message."),
			"description": contentProperty("Supporting guidance."),
			"compact":     modifierProperty("Whether spacing is compact."),
			"bordered":    modifierProperty("Whether a dashed outline is shown."),
			"actions":     contentProperty("Call-to-action definitions."),
		},
		[]designcomponent.Slot{{
			Name:         "actions",
			Description:  "Optional call-to-action controls.",
			AllowedTypes: []string{"Button", "Link"},
			Cardinality:  designcomponent.SlotMany,
		}},
		deliveryDesign("default", "02 Atoms", "Feedback", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"title": "No projects yet", "description": "Create a project to begin.",
				"bordered": true,
				"actions": []map[string]any{{
					"label": "Create project", "href": "/projects/new", "variant": "primary",
				}},
			}),
			mobileDeliveryExample(componentType, "compact", map[string]any{
				"title": "No results", "description": "Try a different filter.", "compact": true,
			}),
		},
		EmptyState,
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
			"size": sizeProperty([]string{"xs", "small", "medium", "large"}, "medium", "Key cap size."),
		},
		nil,
		deliveryDesign("shortcut", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "shortcut", map[string]any{
				"keys": []string{"⌘", "K"}, "size": "medium",
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
		deliveryText("Label", "", "View documentation", "text"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Accessible navigation link with safe external behavior.",
		map[string]PropertyContract{
			"text":     contentProperty("Visible link label."),
			"href":     contentProperty("Navigation target."),
			"external": stateProperty("Whether the target opens externally."),
			"variant":  variantProperty([]string{"primary", "secondary", "text", "underline"}, "primary", "Visual treatment."),
			"target":   contentProperty("Optional browsing context."),
			"rel":      contentProperty("Link relationship tokens."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Actions", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"text": "View documentation", "href": "/docs",
			}),
			deliveryExample(componentType, "external", map[string]any{
				"text": "Open reference", "href": "https://example.com", "external": true,
			}),
		},
		Link,
	)
}

func tagDeliveryDefinition() DeliveryDefinition {
	componentType := "Tag"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clTagBase.Compile()+" "+clTagIdle.Compile(),
			deliveryText("Label", "", "Design", "text"),
		),
		"selected",
		map[string]string{
			"false": clTagIdle.Compile(),
			"true":  clTagSelected.Compile(),
		},
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierAtom,
		stableDeliveryID(componentType),
		"Selectable or removable metadata chip.",
		map[string]PropertyContract{
			"text":        contentProperty("Visible tag label."),
			"icon":        contentProperty("Optional icon identity."),
			"removable":   stateProperty("Whether removal is available."),
			"selected":    stateProperty("Whether the tag is selected."),
			"variant":     toneProperty([]string{"default", "primary", "success", "warning", "error"}, "default", "Semantic treatment."),
			"onRemoveURL": contentProperty("Optional progressive-enhancement removal endpoint."),
		},
		nil,
		deliveryDesign("default", "02 Atoms", "Status", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"text": "Design",
			}),
			deliveryExample(componentType, "selected", map[string]any{
				"text": "Design", "selected": true,
			}),
		},
		Tag,
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
