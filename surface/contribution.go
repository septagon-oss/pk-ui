// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

// Package surface defines the renderer-neutral UI contribution seam: the
// structural vocabulary an extension uses to describe routes, navigation,
// widgets, settings, and page composition to any rendering shell.
// Contributions are pure data — no handlers and no product renderer types —
// so applications can extend the OSS contract without reversing ownership.
//
// PlatformKit collects [Contribution] values through its application
// composition layer; other consumers may use any discovery mechanism.
package surface

// contribution.go — surface contribution vocabulary and validation.
//

import (
	"fmt"
	"slices"
	"strings"
)

// Target identifies which product surface a contribution is meant for.
type Target string

// Surface targets.
const (
	TargetAdmin    Target = "admin"
	TargetApp      Target = "app"
	TargetAuth     Target = "auth"
	TargetOperator Target = "operator"
	TargetPublic   Target = "public"
)

// PagePattern gives consuming shells a transport-safe rendering hint.
type PagePattern string

// Page pattern hints.
const (
	PagePatternUnknown    PagePattern = "unknown"
	PagePatternLanding    PagePattern = "landing"
	PagePatternList       PagePattern = "list"
	PagePatternDetail     PagePattern = "detail"
	PagePatternSettings   PagePattern = "settings"
	PagePatternWorkflow   PagePattern = "workflow"
	PagePatternReview     PagePattern = "review"
	PagePatternDashboard  PagePattern = "dashboard"
	PagePatternOnboarding PagePattern = "onboarding"
)

// ComponentLevel is the transport-side component-depth indicator a module
// emits alongside a PageContract. Shells translate it into their own
// component vocabulary; it stays a plain string so no shell's enum leaks
// across the seam.
type ComponentLevel string

// Component depth hints.
const (
	ComponentLevelUnknown  ComponentLevel = "unknown"
	ComponentLevelAtom     ComponentLevel = "atom"
	ComponentLevelMolecule ComponentLevel = "molecule"
	ComponentLevelOrganism ComponentLevel = "organism"
	ComponentLevelTemplate ComponentLevel = "template"
	ComponentLevelPage     ComponentLevel = "page"
)

// Text carries copy metadata without leaking frontend implementation.
type Text struct {
	Key      string `json:"key,omitempty"`
	Fallback string `json:"fallback"`
}

// Effective returns the best available copy value for transport consumers.
func (t Text) Effective() string {
	if strings.TrimSpace(t.Fallback) != "" {
		return strings.TrimSpace(t.Fallback)
	}
	return strings.TrimSpace(t.Key)
}

// NavSection describes a shell-neutral navigation grouping.
type NavSection struct {
	ID    string `json:"id"`
	Label Text   `json:"label"`
	Icon  string `json:"icon,omitempty"`
	// Order is the sort position within the navigation.
	//port:allow-noun order — sort position, not commerce (ADR-0001)
	Order int `json:"order,omitempty"`
}

// Route is the route-level contribution emitted by a module.
//
// The optional Page field carries a richer page-system contract — a
// module that wants to own its page composition (template + slots +
// presenter + storybook + design artifacts) publishes it here. Shells
// that only care about routing and navigation ignore Page entirely;
// shells that render the page translate it into their own page contract
// at the boundary. Keeping Page optional means the contract stays
// backward-compatible for modules that only contribute routes and
// navigation metadata.
type Route struct {
	ID            string      `json:"id"`
	Path          string      `json:"path"`
	Title         Text        `json:"title"`
	NavLabel      Text        `json:"navLabel"`
	Icon          string      `json:"icon,omitempty"`
	ParentRouteID string      `json:"parentRouteId,omitempty"`
	Section       *NavSection `json:"section,omitempty"`
	// Order is the sort position within the navigation.
	//port:allow-noun order — sort position, not commerce (ADR-0001)
	Order          int           `json:"order,omitempty"`
	PagePattern    PagePattern   `json:"pagePattern,omitempty"`
	Page           *PageContract `json:"page,omitempty"`
	CapabilityTags []string      `json:"capabilityTags,omitempty"`
	Targets        []Target      `json:"targets,omitempty"`
	Hidden         bool          `json:"hidden,omitempty"`

	// RendersEntity declares that this route renders rows for a platform
	// entity. When set, it is the contract a route makes with the surface
	// renderer: "I will render entity X via the row-source registry;
	// please verify a matching row source is registered". The boot
	// validator cancels startup when a route with this field set has no
	// matching registered row source. Empty = the route does not render
	// an entity table; the validator skips it.
	RendersEntity *RenderedEntity `json:"rendersEntity,omitempty"`
}

// RenderedEntity declares that a surface route renders rows for a
// specific platform entity. ModuleID is the producer module owning the
// entity; EntityType is the entity name as the renderer uses it. Both
// are required — leaving either side empty is treated as "route does not
// render an entity" and the validator skips the route.
type RenderedEntity struct {
	ModuleID   string `json:"moduleId"`
	EntityType string `json:"entityType"`
}

// IncludesTarget reports whether the route is intended for the provided
// surface.
func (r Route) IncludesTarget(target Target) bool {
	return slices.Contains(r.Targets, target)
}

// EffectiveNavLabel returns the preferred navigation label for the route.
func (r Route) EffectiveNavLabel() string {
	if label := r.NavLabel.Effective(); label != "" {
		return label
	}
	return r.Title.Effective()
}

// Widget is a structural dashboard-widget contribution: data only, no
// render funcs — the shell decides how (and whether) to render each Kind.
type Widget struct {
	ID    string `json:"id"`
	Title Text   `json:"title"`
	// Kind names the widget archetype ("stats", "chart", "table",
	// "list", ...); shells map it onto their component vocabulary.
	Kind string `json:"kind"`
	// Size is a shell-neutral sizing hint ("small", "medium", "large").
	Size string `json:"size,omitempty"`
	// Order is the sort position on the dashboard.
	//port:allow-noun order — sort position, not commerce (ADR-0001)
	Order          int      `json:"order,omitempty"`
	CapabilityTags []string `json:"capabilityTags,omitempty"`
	// DataRef names the data source the shell should bind ("widget id"
	// convention of the consuming dashboard; empty = shell default).
	DataRef string `json:"dataRef,omitempty"`
}

// Setting is one structural setting declaration.
type Setting struct {
	Key   string `json:"key"`
	Title Text   `json:"title"`
	// Kind is the value type ("bool", "string", "int", "select",
	// "password", ...).
	Kind string `json:"kind"`
	// Default is the string form of the default value.
	Default        string   `json:"default,omitempty"`
	CapabilityTags []string `json:"capabilityTags,omitempty"`
	// Sensitive marks values that must be masked in UIs and excluded
	// from exports/logs (secrets, credentials). Shells must fail closed:
	// a sensitive setting is never rendered or serialized in clear.
	Sensitive bool `json:"sensitive,omitempty"`
}

// SettingsGroup is a structural group of module settings.
type SettingsGroup struct {
	ID       string    `json:"id"`
	Title    Text      `json:"title"`
	Settings []Setting `json:"settings,omitempty"`
}

// Contribution is the surface metadata published by a module.
type Contribution struct {
	ModuleID string  `json:"moduleId"`
	Routes   []Route `json:"routes"`
	// Widgets are structural dashboard contributions (optional).
	Widgets []Widget `json:"widgets,omitempty"`
	// Settings are structural settings groups (optional).
	Settings []SettingsGroup `json:"settings,omitempty"`
}

// ShellProfile identifies the shell mechanics a page targets. It is an alias
// of Target because a route and its page must use the same canonical surface
// vocabulary.
type ShellProfile = Target

const (
	ShellProfileAdmin    = TargetAdmin
	ShellProfileApp      = TargetApp
	ShellProfileAuth     = TargetAuth
	ShellProfileOperator = TargetOperator
	ShellProfilePublic   = TargetPublic
)

// PresenterContract describes the presenter and view-model boundary for a
// governed page.
type PresenterContract struct {
	Contract    string `json:"contract"`
	ViewModel   string `json:"viewModel"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TemplateContract describes the template boundary for a governed page.
type TemplateContract struct {
	Contract    string         `json:"contract"`
	Name        string         `json:"name"`
	Level       ComponentLevel `json:"level"`
	Description string         `json:"description,omitempty"`
}

// SlotContract describes a governed slot inside a page template.
type SlotContract struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Level       ComponentLevel `json:"level"`
	Contract    string         `json:"contract"`
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required,omitempty"`
}

// StorybookReference points the docs / preview harness at the canonical
// Storybook destination for the page. All fields optional, but at least
// one of Package, Story, or URL must be set for the reference to be
// meaningful.
type StorybookReference struct {
	Package string `json:"package,omitempty"`
	Story   string `json:"story,omitempty"`
	URL     string `json:"url,omitempty"`
	Kind    string `json:"kind,omitempty"`
}

// DesignArtifactSet points to the compiled design artifacts associated
// with a page. All fields optional; an all-empty set is treated as "no
// artifacts".
type DesignArtifactSet struct {
	Manifest     string `json:"manifest,omitempty"`
	Preview      string `json:"preview,omitempty"`
	Bundle       string `json:"bundle,omitempty"`
	Figma        string `json:"figma,omitempty"`
	ClaudeDesign string `json:"claudeDesign,omitempty"`
}

// PageContract is the canonical renderer-neutral page-system contract. Shells
// consume this exact type and add only their own chrome/runtime state; they do
// not translate it into a parallel page model.
type PageContract struct {
	Level           ComponentLevel      `json:"level"`
	RouteID         string              `json:"routeId"`
	ShellProfile    ShellProfile        `json:"shellProfile"`
	Pattern         PagePattern         `json:"pattern"`
	Presenter       PresenterContract   `json:"presenter"`
	Template        TemplateContract    `json:"template"`
	Slots           []SlotContract      `json:"slots,omitempty"`
	Storybook       *StorybookReference `json:"storybook,omitempty"`
	DesignArtifacts *DesignArtifactSet  `json:"designArtifacts,omitempty"`
}

// NamespacedRouteID returns the globally-namespaced form
// "<module>.<route>" for a route id. Idempotent: an already-prefixed
// route id is returned unchanged; empty inputs are returned as-is. This
// is the canonical helper every surface host uses to key renderers so two
// modules can reuse a local route name (e.g. "settings") without
// colliding globally.
func NamespacedRouteID(moduleID, routeID string) string {
	module := strings.TrimSpace(moduleID)
	route := strings.TrimSpace(routeID)
	if module == "" || route == "" {
		return route
	}
	if strings.HasPrefix(route, module+".") {
		return route
	}
	return module + "." + route
}

// ValidateContributionsGlobal validates a full set of contributions
// across modules: each contribution's own shape plus global route-id
// uniqueness — the surface host keys renderers by route id, so a global
// collision would make one route shadow the other. Modules avoid
// collisions by namespacing their route ids (see NamespacedRouteID).
func ValidateContributionsGlobal(contributions []Contribution) error {
	owners := make(map[string]string)
	allRoutes := make([]Route, 0)
	for _, contribution := range contributions {
		if err := ValidateContribution(contribution); err != nil {
			return err
		}
		moduleID := strings.TrimSpace(contribution.ModuleID)
		for _, route := range contribution.Routes {
			routeID := strings.TrimSpace(route.ID)
			if prevOwner, exists := owners[routeID]; exists {
				return fmt.Errorf(
					"surface contributions invalid: duplicate global route id %q (declared by %q and %q); namespace route ids as \"<module>.<route>\"",
					routeID, prevOwner, moduleID,
				)
			}
			owners[routeID] = moduleID
			allRoutes = append(allRoutes, route)
		}
	}
	global := Contribution{
		ModuleID: "global surface graph",
		Routes:   allRoutes,
	}
	if err := validateRouteHierarchy(global, true); err != nil {
		return err
	}
	if err := ValidateRouteOwnership(global); err != nil {
		return err
	}
	return nil
}

// ValidateContribution enforces the supported contract shape.
func ValidateContribution(contribution Contribution) error {
	moduleID := strings.TrimSpace(contribution.ModuleID)
	if moduleID == "" {
		return fmt.Errorf("surface contribution invalid: module_id is required")
	}

	seenRoutes := make(map[string]struct{}, len(contribution.Routes))
	for _, route := range contribution.Routes {
		routeID := strings.TrimSpace(route.ID)
		if routeID == "" {
			return fmt.Errorf("surface contribution invalid for module %q: route id is required", moduleID)
		}
		if _, exists := seenRoutes[routeID]; exists {
			return fmt.Errorf("surface contribution invalid for module %q: duplicate route id %q", moduleID, routeID)
		}
		seenRoutes[routeID] = struct{}{}

		if _, err := parseRouteLocation(route.Path); err != nil {
			return fmt.Errorf("surface contribution invalid for module %q route %q: %w", moduleID, routeID, err)
		}
		if route.Title.Effective() == "" {
			return fmt.Errorf("surface contribution invalid for module %q route %q: title is required", moduleID, routeID)
		}
		if len(route.Targets) == 0 {
			return fmt.Errorf("surface contribution invalid for module %q route %q: at least one target is required", moduleID, routeID)
		}
		for _, target := range route.Targets {
			switch target {
			case TargetAdmin, TargetApp, TargetAuth, TargetOperator, TargetPublic:
			default:
				return fmt.Errorf("surface contribution invalid for module %q route %q: unsupported target %q", moduleID, routeID, target)
			}
		}
		if route.Section != nil {
			if strings.TrimSpace(route.Section.ID) == "" {
				return fmt.Errorf("surface contribution invalid for module %q route %q: section id is required when section is set", moduleID, routeID)
			}
			if route.Section.Label.Effective() == "" {
				return fmt.Errorf("surface contribution invalid for module %q route %q: section label is required when section is set", moduleID, routeID)
			}
		}

		if route.Page != nil {
			if err := validatePageContract(moduleID, routeID, *route.Page); err != nil {
				return err
			}
			if !slices.Contains(route.Targets, route.Page.ShellProfile) {
				return fmt.Errorf(
					"surface contribution invalid for module %q route %q: page shell profile %q is not declared in route targets",
					moduleID,
					routeID,
					route.Page.ShellProfile,
				)
			}
			if route.PagePattern != "" &&
				route.PagePattern != PagePatternUnknown &&
				route.PagePattern != route.Page.Pattern {
				return fmt.Errorf(
					"surface contribution invalid for module %q route %q: route page pattern %q does not match page pattern %q",
					moduleID,
					routeID,
					route.PagePattern,
					route.Page.Pattern,
				)
			}
		}
	}
	if err := ValidateRouteOwnership(contribution); err != nil {
		return err
	}

	seenWidgets := make(map[string]struct{}, len(contribution.Widgets))
	for _, widget := range contribution.Widgets {
		widgetID := strings.TrimSpace(widget.ID)
		if widgetID == "" {
			return fmt.Errorf("surface contribution invalid for module %q: widget id is required", moduleID)
		}
		if _, exists := seenWidgets[widgetID]; exists {
			return fmt.Errorf("surface contribution invalid for module %q: duplicate widget id %q", moduleID, widgetID)
		}
		seenWidgets[widgetID] = struct{}{}
		if widget.Title.Effective() == "" {
			return fmt.Errorf("surface contribution invalid for module %q widget %q: title is required", moduleID, widgetID)
		}
		if strings.TrimSpace(widget.Kind) == "" {
			return fmt.Errorf("surface contribution invalid for module %q widget %q: kind is required", moduleID, widgetID)
		}
	}

	seenGroups := make(map[string]struct{}, len(contribution.Settings))
	for _, group := range contribution.Settings {
		groupID := strings.TrimSpace(group.ID)
		if groupID == "" {
			return fmt.Errorf("surface contribution invalid for module %q: settings group id is required", moduleID)
		}
		if _, exists := seenGroups[groupID]; exists {
			return fmt.Errorf("surface contribution invalid for module %q: duplicate settings group id %q", moduleID, groupID)
		}
		seenGroups[groupID] = struct{}{}
		seenKeys := make(map[string]struct{}, len(group.Settings))
		for _, setting := range group.Settings {
			key := strings.TrimSpace(setting.Key)
			if key == "" {
				return fmt.Errorf("surface contribution invalid for module %q settings group %q: setting key is required", moduleID, groupID)
			}
			if _, exists := seenKeys[key]; exists {
				return fmt.Errorf("surface contribution invalid for module %q settings group %q: duplicate setting key %q", moduleID, groupID, key)
			}
			seenKeys[key] = struct{}{}
			if strings.TrimSpace(setting.Kind) == "" {
				return fmt.Errorf("surface contribution invalid for module %q settings group %q setting %q: kind is required", moduleID, groupID, key)
			}
		}
	}

	return nil
}

// ValidatePageContract enforces the renderer-neutral governed page shape.
// Runtime-specific presenter/template existence checks remain downstream, but
// every shell shares these structural invariants.
func ValidatePageContract(page PageContract) error {
	routeID := strings.TrimSpace(page.RouteID)
	if routeID == "" {
		return fmt.Errorf("page contract invalid: route id is required")
	}
	switch page.ShellProfile {
	case ShellProfileAdmin, ShellProfileApp, ShellProfileAuth, ShellProfileOperator, ShellProfilePublic:
	default:
		return fmt.Errorf("page contract invalid for %q: unsupported shell profile %q", routeID, page.ShellProfile)
	}
	switch page.Pattern {
	case PagePatternLanding, PagePatternDashboard, PagePatternList, PagePatternDetail, PagePatternSettings, PagePatternWorkflow, PagePatternReview, PagePatternOnboarding:
	default:
		return fmt.Errorf("page contract invalid for %q: page pattern is required and must be supported", routeID)
	}
	if page.Level != ComponentLevelPage {
		return fmt.Errorf("page contract invalid for %q: page level must be %q", routeID, ComponentLevelPage)
	}
	if strings.TrimSpace(page.Presenter.Contract) == "" {
		return fmt.Errorf("page contract invalid for %q: presenter contract is required", routeID)
	}
	if strings.TrimSpace(page.Presenter.ViewModel) == "" {
		return fmt.Errorf("page contract invalid for %q: presenter view-model is required", routeID)
	}
	if strings.TrimSpace(page.Template.Contract) == "" {
		return fmt.Errorf("page contract invalid for %q: template contract is required", routeID)
	}
	if strings.TrimSpace(page.Template.Name) == "" {
		return fmt.Errorf("page contract invalid for %q: template name is required", routeID)
	}
	if page.Template.Level != ComponentLevelTemplate {
		return fmt.Errorf("page contract invalid for %q: template level must be %q", routeID, ComponentLevelTemplate)
	}

	seenSlots := make(map[string]struct{}, len(page.Slots))
	for _, slot := range page.Slots {
		slotID := strings.TrimSpace(slot.ID)
		if slotID == "" {
			return fmt.Errorf("page contract invalid for %q: slot id is required", routeID)
		}
		if _, duplicate := seenSlots[slotID]; duplicate {
			return fmt.Errorf("page contract invalid for %q: duplicate slot id %q", routeID, slotID)
		}
		seenSlots[slotID] = struct{}{}
		if strings.TrimSpace(slot.Name) == "" {
			return fmt.Errorf("page contract invalid for %q slot %q: slot name is required", routeID, slotID)
		}
		if strings.TrimSpace(slot.Contract) == "" {
			return fmt.Errorf("page contract invalid for %q slot %q: slot contract is required", routeID, slotID)
		}
		switch slot.Level {
		case ComponentLevelAtom, ComponentLevelMolecule, ComponentLevelOrganism, ComponentLevelTemplate, ComponentLevelPage:
		default:
			return fmt.Errorf("page contract invalid for %q slot %q: unsupported component level %q", routeID, slotID, slot.Level)
		}
	}

	if page.Storybook != nil &&
		strings.TrimSpace(page.Storybook.Package) == "" &&
		strings.TrimSpace(page.Storybook.Story) == "" &&
		strings.TrimSpace(page.Storybook.URL) == "" {
		return fmt.Errorf("page contract invalid for %q: storybook reference needs at least one of package, story, or url", routeID)
	}
	if page.DesignArtifacts != nil &&
		strings.TrimSpace(page.DesignArtifacts.Manifest) == "" &&
		strings.TrimSpace(page.DesignArtifacts.Preview) == "" &&
		strings.TrimSpace(page.DesignArtifacts.Bundle) == "" &&
		strings.TrimSpace(page.DesignArtifacts.Figma) == "" &&
		strings.TrimSpace(page.DesignArtifacts.ClaudeDesign) == "" {
		return fmt.Errorf("page contract invalid for %q: design artifacts are empty", routeID)
	}
	return nil
}

// ClonePageContract returns an isolated copy suitable for composition and
// projection.
func ClonePageContract(page *PageContract) *PageContract {
	if page == nil {
		return nil
	}
	cloned := *page
	cloned.Slots = append([]SlotContract(nil), page.Slots...)
	if page.Storybook != nil {
		storybook := *page.Storybook
		cloned.Storybook = &storybook
	}
	if page.DesignArtifacts != nil {
		artifacts := *page.DesignArtifacts
		cloned.DesignArtifacts = &artifacts
	}
	return &cloned
}

func validatePageContract(moduleID, routeID string, page PageContract) error {
	pageRouteID := strings.TrimSpace(page.RouteID)
	if pageRouteID == "" {
		return fmt.Errorf("surface contribution invalid for module %q route %q: page route id is required when page is set", moduleID, routeID)
	}
	if pageRouteID != routeID {
		return fmt.Errorf("surface contribution invalid for module %q route %q: page route id %q does not match route id", moduleID, routeID, pageRouteID)
	}
	if err := ValidatePageContract(page); err != nil {
		return fmt.Errorf("surface contribution invalid for module %q route %q: %w", moduleID, routeID, err)
	}
	return nil
}
