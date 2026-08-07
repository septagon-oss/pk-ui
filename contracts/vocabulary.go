// Implements: REQ-011 (one canonical design vocabulary for every renderer).
// Per: ADR-0031.
// Discipline: C-14.

package contracts

import "strings"

// vocabulary.go declares the canonical design vocabulary shared by every
// PlatformKit renderer and projection. The Props schemas in this package are
// the delivery source of truth for the web renderer, the A2UI/native
// projections, the Storybook generator, and the Figma adapter; this file pins
// the values those schemas may use on their shared axes (size, tone, state,
// status) so web and mobile can never drift into parallel vocabularies.
//
// vocabulary_test.go enforces the bindings below against the actual contract
// source: every size/tone/state/status enum documented on a Props field must
// be declared in an axis table here, must be a subset of its family's
// canonical values, and must not mix alias spellings (small/md) with
// canonical ones. New components that add an axis prop must add a binding
// here first — that is the ratchet.

// ScaleFamily identifies a deliberate size axis. Families are disjoint on
// purpose: a component never offers both families on one prop.
type ScaleFamily string

const (
	// ScaleFamilyControl is the canonical control scale shared by interactive
	// and display primitives. Components with a narrow ergonomic band (form
	// controls, progress) use a contiguous lower subset of the same scale.
	ScaleFamilyControl ScaleFamily = "control"
	// ScaleFamilyOverlay is the container-extent scale for overlay surfaces
	// (Modal, Drawer). "full" is a layout extent, not a size; it belongs to
	// this family only.
	ScaleFamilyOverlay ScaleFamily = "overlay"
	// ScaleFamilyType is the typographic scale for text atoms.
	ScaleFamilyType ScaleFamily = "type"
)

// Canonical values per scale family, ordered.
var (
	ControlSizes = []string{"xs", "sm", "md", "lg", "xl", "2xl"}
	OverlaySizes = []string{"small", "medium", "large", "xl", "full"}
	TypeSizes    = []string{"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl"}
)

// ControlSizeAliases maps accepted legacy spellings onto the canonical
// control scale. Renderers accept aliases on intake for backward
// compatibility; authored contracts, manifests, and documentation always use
// canonical values.
var ControlSizeAliases = map[string]string{
	"extra-small":       "xs",
	"small":             "sm",
	"medium":            "md",
	"large":             "lg",
	"extra-large":       "xl",
	"xxl":               "2xl",
	"extra-extra-large": "2xl",
}

// CanonicalTones is the shared semantic tone vocabulary.
var CanonicalTones = []string{"neutral", "brand", "success", "warning", "danger", "info"}

// FeedbackTones is the restricted set for inline feedback surfaces (Alert,
// Toast): neutral plus the four status tones.
var FeedbackTones = []string{"neutral", "info", "success", "warning", "danger"}

// ToneAliases maps legacy / cross-platform spellings onto canonical tones.
// Native renderers historically used "accent" for the brand tone; delivery
// uses "error" for danger. Renderers resolve aliases on intake.
var ToneAliases = map[string]string{
	"accent":      "brand",
	"error":       "danger",
	"destructive": "danger",
}

// CollectionStates is the state axis for data-bearing surfaces
// (WindowedCollection). "empty" is derived from itemCount, not a state value.
var CollectionStates = []string{"ready", "loading", "error"}

// PresenceStates is the status axis for identity marks (Avatar).
var PresenceStates = []string{"online", "offline", "busy", "away"}

// StepStates is the status axis for step indicators (Stepper items).
var StepStates = []string{"pending", "active", "completed", "error"}

// ScaleAxisBinding binds a component's size prop to its scale family. The
// key is "<StructName>.<jsonName>"; every documented size enum in this
// package must appear here (enforced by vocabulary_test.go).
var ScaleAxisBindings = map[string]ScaleFamily{
	// Full control scale.
	"ButtonProps.size":  ScaleFamilyControl,
	"BadgeProps.size":   ScaleFamilyControl,
	"IconProps.size":    ScaleFamilyControl,
	"SpinnerProps.size": ScaleFamilyControl,
	"KbdProps.size":     ScaleFamilyControl,
	"AvatarProps.size":  ScaleFamilyControl,
	// Restricted control band (sm–lg).
	"ToggleProps.size":   ScaleFamilyControl,
	"InputProps.size":    ScaleFamilyControl,
	"ProgressProps.size": ScaleFamilyControl,
	"SkeletonProps.size": ScaleFamilyControl,
	"DropdownProps.size": ScaleFamilyControl,
	// Overlay extent scale.
	"ModalProps.size":  ScaleFamilyOverlay,
	"DrawerProps.size": ScaleFamilyOverlay,
	// Typographic scale.
	"TextProps.size": ScaleFamilyType,
}

// ToneAxisBindings records every component prop that carries the canonical
// tone vocabulary (or a declared subset of it). Subsets are allowed;
// divergence from the canonical spellings is not.
var ToneAxisBindings = map[string]bool{
	"ButtonProps.tone":    true,
	"BadgeProps.tone":     true,
	"IconProps.tone":      true,
	"AvatarProps.tone":    true,
	"TagProps.tone":       true,
	"InputProps.tone":     true,
	"SliderProps.tone":    true,
	"SpinnerProps.tone":   true,
	"ProgressProps.tone":  true,
	"AlertProps.tone":     true,
	"ToastProps.tone":     true,
	"ActionMenuItem.tone": true,
	"SidebarSection.tone": true,
}

// StateAxisBindings records the component props carrying a declared state or
// status axis, keyed like the scale bindings. Values must match the listed
// canonical set.
var StateAxisBindings = map[string][]string{
	"WindowedCollectionProps.state": CollectionStates,
	"AvatarProps.status":            PresenceStates,
	"StepItem.status":               StepStates,
}

// ScaleValues returns the canonical values for a family.
func ScaleValues(family ScaleFamily) []string {
	switch family {
	case ScaleFamilyControl:
		return ControlSizes
	case ScaleFamilyOverlay:
		return OverlaySizes
	case ScaleFamilyType:
		return TypeSizes
	default:
		return nil
	}
}

// NormalizeControlSize maps an accepted spelling (canonical or legacy alias)
// onto the canonical control scale, or returns "" for unknown values.
func NormalizeControlSize(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return ""
	}
	for _, size := range ControlSizes {
		if v == size {
			return v
		}
	}
	if canonical, ok := ControlSizeAliases[v]; ok {
		return canonical
	}
	return ""
}

// CanonicalTone resolves a tone spelling (canonical or alias) to its
// canonical value, or returns "" for unknown values.
func CanonicalTone(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return ""
	}
	for _, tone := range CanonicalTones {
		if v == tone {
			return v
		}
	}
	if canonical, ok := ToneAliases[v]; ok {
		return canonical
	}
	return ""
}
