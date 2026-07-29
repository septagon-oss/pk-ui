// Implements: REQ-011.
// Per: ADR-0007, ADR-0031.
// Discipline: C-14.

package surface

import (
	"fmt"
	"strings"
)

const (
	adminPagePresenterContract = "github.com/septagon-oss/pk-ui/surface.PresenterContract"
	adminPageTemplateContract  = "github.com/septagon-oss/pk-ui/surface.TemplateContract"
)

// AdminPageContract materializes the canonical renderer-neutral page contract
// for an admin route. Unknown patterns intentionally remain unmaterialized.
func AdminPageContract(routeID string, pattern PagePattern) *PageContract {
	return canonicalizeAdminPageContract(routeID, pattern, nil)
}

// CanonicalizeAdminRoute fills missing page metadata from PagePattern and
// normalizes an explicit page contract to the shared admin-shell defaults.
func CanonicalizeAdminRoute(route Route) Route {
	route.ID = strings.TrimSpace(route.ID)
	if route.PagePattern == "" && route.Page != nil && route.Page.Pattern != "" {
		route.PagePattern = route.Page.Pattern
	}
	route.Page = canonicalizeAdminPageContract(route.ID, route.PagePattern, route.Page)
	if route.Page != nil {
		route.PagePattern = route.Page.Pattern
	}
	return route
}

// CanonicalizeAdminSurfaceContribution returns an isolated contribution whose
// admin routes carry explicit page contracts wherever their patterns are
// known. Routes with no explicit target retain the historical admin default;
// routes explicitly targeting only other surfaces are left untouched so a
// mixed contribution does not acquire an invalid admin shell contract.
func CanonicalizeAdminSurfaceContribution(contribution Contribution) Contribution {
	contribution.ModuleID = strings.TrimSpace(contribution.ModuleID)
	if len(contribution.Routes) == 0 {
		return contribution
	}

	routes := make([]Route, len(contribution.Routes))
	for index, route := range contribution.Routes {
		if len(route.Targets) == 0 || route.IncludesTarget(TargetAdmin) {
			route = CanonicalizeAdminRoute(route)
		}
		routes[index] = route
	}
	contribution.Routes = routes
	return contribution
}

// NamespacedSurfaceContribution returns an isolated contribution whose route
// IDs, parent references, and page route IDs use "<module>.<route>". It is
// idempotent. Invalid authoring forms panic because contributions are static
// declarations and a malformed declaration is a programming error.
func NamespacedSurfaceContribution(contribution Contribution) Contribution {
	if len(contribution.Routes) == 0 {
		return contribution
	}
	if contribution.ModuleID == "" ||
		contribution.ModuleID != strings.TrimSpace(contribution.ModuleID) ||
		contribution.ModuleID != strings.ToLower(contribution.ModuleID) ||
		strings.Contains(contribution.ModuleID, ".") {
		panic(fmt.Errorf("surface contribution requires a lowercase canonical module ID, got %q", contribution.ModuleID))
	}

	routes := make([]Route, len(contribution.Routes))
	for index, route := range contribution.Routes {
		route.ID = canonicalNamespacedRouteID(contribution.ModuleID, route.ID)
		route.ParentRouteID = canonicalNamespacedRouteID(contribution.ModuleID, route.ParentRouteID)
		if route.Page != nil {
			route.Page = ClonePageContract(route.Page)
			route.Page.RouteID = route.ID
		}
		routes[index] = route
	}
	contribution.Routes = routes
	return contribution
}

func canonicalizeAdminPageContract(routeID string, pattern PagePattern, contract *PageContract) *PageContract {
	if contract == nil && (pattern == "" || pattern == PagePatternUnknown) {
		return nil
	}

	cloned := ClonePageContract(contract)
	if cloned == nil {
		cloned = &PageContract{}
	}

	if strings.TrimSpace(cloned.RouteID) == "" {
		cloned.RouteID = strings.TrimSpace(routeID)
	}
	if strings.TrimSpace(cloned.RouteID) == "" {
		return nil
	}
	if cloned.Pattern == "" || cloned.Pattern == PagePatternUnknown {
		cloned.Pattern = pattern
	}
	if cloned.Pattern == "" || cloned.Pattern == PagePatternUnknown {
		return nil
	}

	cloned.ShellProfile = ShellProfileAdmin
	cloned.Level = ComponentLevelPage
	if strings.TrimSpace(cloned.Presenter.Contract) == "" {
		cloned.Presenter.Contract = adminPagePresenterContract
	}
	if strings.TrimSpace(cloned.Presenter.ViewModel) == "" {
		cloned.Presenter.ViewModel = adminPageViewModelName(cloned.RouteID)
	}
	if strings.TrimSpace(cloned.Template.Contract) == "" {
		cloned.Template.Contract = adminPageTemplateContract
	}
	if strings.TrimSpace(cloned.Template.Name) == "" {
		cloned.Template.Name = adminPageTemplateName(cloned.Pattern)
	}
	cloned.Template.Level = ComponentLevelTemplate
	return cloned
}

func canonicalNamespacedRouteID(moduleID, routeID string) string {
	if routeID == "" {
		return ""
	}
	if routeID != strings.TrimSpace(routeID) || routeID != strings.ToLower(routeID) {
		panic(fmt.Errorf("surface route %q must be a lowercase local or canonical ID", routeID))
	}
	if strings.HasPrefix(routeID, moduleID+"-") {
		panic(fmt.Errorf("surface route %q uses removed hyphen-qualified ID; publish a local route ID or %q", routeID, moduleID+".<route>"))
	}
	if strings.Contains(routeID, ".") {
		localID, ok := strings.CutPrefix(routeID, moduleID+".")
		if !ok || localID == "" || strings.Contains(localID, ".") {
			panic(fmt.Errorf("surface route %q uses a foreign or nested namespace; owner is %q", routeID, moduleID))
		}
	}
	return NamespacedRouteID(moduleID, routeID)
}

func adminPageViewModelName(routeID string) string {
	tokens := strings.FieldsFunc(strings.TrimSpace(routeID), func(r rune) bool {
		switch r {
		case '-', '_', '/', '.', ':', ' ':
			return true
		default:
			return false
		}
	})
	if len(tokens) == 0 {
		return "AdminPageViewModel"
	}

	var builder strings.Builder
	for _, token := range tokens {
		if token == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(token[:1]))
		if len(token) > 1 {
			builder.WriteString(strings.ToLower(token[1:]))
		}
	}
	if builder.Len() == 0 {
		return "AdminPageViewModel"
	}
	builder.WriteString("PageViewModel")
	return builder.String()
}

func adminPageTemplateName(pattern PagePattern) string {
	switch pattern {
	case PagePatternDashboard:
		return "DashboardPage"
	case PagePatternDetail, PagePatternReview:
		return "DetailPage"
	case PagePatternSettings:
		return "SettingsPage"
	case PagePatternWorkflow, PagePatternOnboarding:
		return "WorkflowPage"
	case PagePatternLanding:
		return "LandingPage"
	default:
		return "ListPage"
	}
}
