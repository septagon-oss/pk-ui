// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package surface

import (
	"fmt"
	"strings"
	"testing"
)

func TestRouteOwnersForTargetSelectsConcreteOwnersWithoutMutatingNavigation(t *testing.T) {
	contribution := Contribution{
		ModuleID: "content_management",
		Routes: []Route{
			{
				ID:             "content",
				Path:           "/admin/content/docs",
				Title:          Text{Fallback: "Content"},
				CapabilityTags: nil,
				Targets:        []Target{TargetAdmin},
			},
			{
				ID:             "docs",
				Path:           "/admin/content/docs",
				Title:          Text{Fallback: "Docs"},
				ParentRouteID:  "content",
				CapabilityTags: []string{"content:read"},
				Targets:        []Target{TargetAdmin},
			},
			{
				ID:             "arc42",
				Path:           "/admin/content/docs?collection=arc42",
				Title:          Text{Fallback: "Arc42"},
				ParentRouteID:  "content",
				CapabilityTags: []string{"content:read"},
				Targets:        []Target{TargetAdmin},
			},
			{
				ID:      "articles",
				Path:    "/admin/content/articles/",
				Title:   Text{Fallback: "Articles"},
				Targets: []Target{TargetAdmin},
			},
			{
				ID:      "public-docs",
				Path:    "/docs",
				Title:   Text{Fallback: "Public docs"},
				Targets: []Target{TargetPublic},
			},
		},
	}

	owners, err := RouteOwnersForTarget(contribution, TargetAdmin)
	if err != nil {
		t.Fatalf("RouteOwnersForTarget() error = %v", err)
	}
	if len(owners) != 2 {
		t.Fatalf("RouteOwnersForTarget() returned %d owners, want 2: %#v", len(owners), owners)
	}
	if owners[0].ID != "docs" || owners[0].Path != "/admin/content/docs" {
		t.Fatalf("canonical docs owner = %#v, want deepest query-free child", owners[0])
	}
	if owners[1].ID != "articles" {
		t.Fatalf("second owner = %#v, want articles", owners[1])
	}
	if len(contribution.Routes) != 5 || contribution.Routes[2].Path != "/admin/content/docs?collection=arc42" {
		t.Fatalf("RouteOwnersForTarget() mutated contribution: %#v", contribution.Routes)
	}
}

func TestRouteOwnersForTargetRejectsSiblingAliases(t *testing.T) {
	contribution := Contribution{
		ModuleID: "content_management",
		Routes: []Route{
			{ID: "documentation", Path: "/admin/content/docs", Title: Text{Fallback: "Documentation"}, Targets: []Target{TargetAdmin}},
			{ID: "knowledge", Path: "/admin/content/docs", Title: Text{Fallback: "Knowledge"}, Targets: []Target{TargetAdmin}},
			{ID: "runbooks", Path: "/admin/content/docs?collection=runbooks", Title: Text{Fallback: "Runbooks"}, Targets: []Target{TargetAdmin}},
		},
	}

	owners, err := RouteOwnersForTarget(contribution, TargetAdmin)
	if err == nil || owners != nil {
		t.Fatalf("RouteOwnersForTarget() = %#v, %v; want fail-closed sibling rejection", owners, err)
	}
}

func TestValidateRouteOwnershipRejectsSiblingConcreteOwners(t *testing.T) {
	contribution := Contribution{
		ModuleID: "content_management",
		Routes: []Route{
			{ID: "documentation", Path: "/admin/content/docs", CapabilityTags: []string{"content:read"}, Targets: []Target{TargetAdmin}},
			{ID: "knowledge", Path: "/admin/content/docs", CapabilityTags: []string{"content:read"}, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted sibling query-free owners")
	}
}

func TestValidateRouteOwnershipAllowsDefaultChildAndCompatibleStateAliases(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{
			{ID: "orders", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
			{ID: "orders-list", ParentRouteID: "orders", Path: "/admin/orders", CapabilityTags: []string{"orders:read"}, Targets: []Target{TargetAdmin}},
			{ID: "open-orders", ParentRouteID: "orders", Path: "/admin/orders?status=open", CapabilityTags: []string{"orders:read"}, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err != nil {
		t.Fatalf("ValidateRouteOwnership() rejected valid aliases: %v", err)
	}
}

func TestValidateRouteOwnershipRejectsNonEmptyAncestorCapabilityDrift(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{
			{ID: "orders", Path: "/admin/orders", CapabilityTags: []string{"orders:legacy-root"}, Targets: []Target{TargetAdmin}},
			{ID: "orders-list", ParentRouteID: "orders", Path: "/admin/orders", CapabilityTags: []string{"orders:read"}, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted misleading non-empty ancestor capabilities")
	}
}

func TestValidateRouteOwnershipRejectsAliasCapabilityDrift(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{
			{ID: "orders", Path: "/admin/orders", CapabilityTags: []string{"orders:read"}, Targets: []Target{TargetAdmin}},
			{ID: "open-orders", ParentRouteID: "orders", Path: "/admin/orders?status=open", CapabilityTags: []string{"orders:manage"}, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted capability drift")
	}
}

func TestValidateRouteOwnershipRejectsAncestorStateAliasCapabilityDrift(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{
			{ID: "open-orders", Path: "/admin/orders?status=open", CapabilityTags: []string{"orders:manage"}, Targets: []Target{TargetAdmin}},
			{ID: "orders-list", ParentRouteID: "open-orders", Path: "/admin/orders", CapabilityTags: []string{"orders:read"}, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted capability drift on an ancestor state alias")
	}
}

func TestValidateRouteOwnershipRejectsLoneStateAlias(t *testing.T) {
	contribution := Contribution{
		ModuleID: "content_management",
		Routes: []Route{
			{ID: "runbooks", Path: "/admin/content/docs?collection=runbooks", Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted a lone query alias")
	}
}

func TestValidateRouteOwnershipRejectsHiddenDefaultOwnerForVisibleParent(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{
			{ID: "orders", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
			{ID: "orders-list", ParentRouteID: "orders", Path: "/admin/orders", Hidden: true, Targets: []Target{TargetAdmin}},
		},
	}

	if err := ValidateRouteOwnership(contribution); err == nil {
		t.Fatal("ValidateRouteOwnership() accepted a hidden owner shadowing visible navigation")
	}
}

func TestValidateRouteOwnershipRejectsInvalidParentGraphs(t *testing.T) {
	tests := []struct {
		name   string
		routes []Route
		want   string
	}{
		{
			name: "self parent",
			routes: []Route{
				{ID: "orders", ParentRouteID: "orders", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
			},
			want: "cannot parent itself",
		},
		{
			name: "cross target parent",
			routes: []Route{
				{ID: "orders", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
				{ID: "orders-app", ParentRouteID: "orders", Path: "/portal/orders", Targets: []Target{TargetApp}},
			},
			want: "is not declared by parent",
		},
		{
			name: "parent cycle",
			routes: []Route{
				{ID: "orders", ParentRouteID: "orders-list", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
				{ID: "orders-list", ParentRouteID: "orders", Path: "/admin/orders", Targets: []Target{TargetAdmin}},
			},
			want: "parent route cycle",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contribution := Contribution{ModuleID: "order_management", Routes: test.routes}
			err := ValidateRouteOwnership(contribution)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateRouteOwnership() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestValidateRouteOwnershipAllowsExternalParentReference(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{{
			ID:            "order_management.orders",
			ParentRouteID: "commerce.shell",
			Path:          "/admin/orders",
			Targets:       []Target{TargetAdmin},
		}},
	}

	if err := ValidateRouteOwnership(contribution); err != nil {
		t.Fatalf("ValidateRouteOwnership() rejected external parent before global composition: %v", err)
	}
}

func TestValidateContributionsGlobalResolvesExternalParents(t *testing.T) {
	parent := Contribution{
		ModuleID: "commerce",
		Routes: []Route{{
			ID:      "commerce.shell",
			Path:    "/admin/commerce",
			Title:   Text{Fallback: "Commerce"},
			Targets: []Target{TargetAdmin},
		}},
	}
	child := Contribution{
		ModuleID: "order_management",
		Routes: []Route{{
			ID:            "order_management.orders",
			ParentRouteID: "commerce.shell",
			Path:          "/admin/orders",
			Title:         Text{Fallback: "Orders"},
			Targets:       []Target{TargetAdmin},
		}},
	}

	if err := ValidateContributionsGlobal([]Contribution{parent, child}); err != nil {
		t.Fatalf("ValidateContributionsGlobal() rejected valid external parent: %v", err)
	}
}

func TestValidateContributionsGlobalRejectsCrossModuleConcretePathCollision(t *testing.T) {
	contributions := []Contribution{
		{ModuleID: "user_management", Routes: []Route{{
			ID: "user_management.users", Path: "/admin/users", Title: Text{Fallback: "Users"},
			CapabilityTags: []string{"users:read"}, Targets: []Target{TargetAdmin},
		}}},
		{ModuleID: "support_management", Routes: []Route{{
			ID: "support_management.users", Path: "/admin/users", Title: Text{Fallback: "Support users"},
			CapabilityTags: []string{"users:manage"}, Targets: []Target{TargetAdmin},
		}}},
	}

	err := ValidateContributionsGlobal(contributions)
	if err == nil || !strings.Contains(err.Error(), "siblings") {
		t.Fatalf("ValidateContributionsGlobal() error = %v, want cross-module concrete-owner collision", err)
	}
}

func TestValidateContributionsGlobalRejectsInvalidExternalParentGraphs(t *testing.T) {
	tests := []struct {
		name          string
		contributions []Contribution
		want          string
	}{
		{
			name: "missing external parent",
			contributions: []Contribution{{
				ModuleID: "order_management",
				Routes: []Route{{
					ID: "order_management.orders", ParentRouteID: "missing.shell", Path: "/admin/orders",
					Title: Text{Fallback: "Orders"}, Targets: []Target{TargetAdmin},
				}},
			}},
			want: "references missing parent",
		},
		{
			name: "cross target external parent",
			contributions: []Contribution{
				{ModuleID: "commerce", Routes: []Route{{
					ID: "commerce.shell", Path: "/admin/commerce", Title: Text{Fallback: "Commerce"}, Targets: []Target{TargetAdmin},
				}}},
				{ModuleID: "order_management", Routes: []Route{{
					ID: "order_management.orders", ParentRouteID: "commerce.shell", Path: "/portal/orders",
					Title: Text{Fallback: "Orders"}, Targets: []Target{TargetApp},
				}}},
			},
			want: "is not declared by parent",
		},
		{
			name: "cross contribution cycle",
			contributions: []Contribution{
				{ModuleID: "commerce", Routes: []Route{{
					ID: "commerce.shell", ParentRouteID: "order_management.orders", Path: "/admin/commerce",
					Title: Text{Fallback: "Commerce"}, Targets: []Target{TargetAdmin},
				}}},
				{ModuleID: "order_management", Routes: []Route{{
					ID: "order_management.orders", ParentRouteID: "commerce.shell", Path: "/admin/orders",
					Title: Text{Fallback: "Orders"}, Targets: []Target{TargetAdmin},
				}}},
			},
			want: "parent route cycle",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateContributionsGlobal(test.contributions)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateContributionsGlobal() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestValidateRouteOwnershipRejectsExcessiveHierarchyDepth(t *testing.T) {
	routes := make([]Route, 0, maxRouteHierarchyDepth+2)
	for index := range maxRouteHierarchyDepth + 2 {
		parentID := ""
		if index > 0 {
			parentID = fmt.Sprintf("route-%d", index-1)
		}
		routes = append(routes, Route{
			ID:            fmt.Sprintf("route-%d", index),
			ParentRouteID: parentID,
			Path:          fmt.Sprintf("/admin/routes/%d", index),
			Targets:       []Target{TargetAdmin},
		})
	}

	err := ValidateRouteOwnership(Contribution{ModuleID: "deep_management", Routes: routes})
	if err == nil || !strings.Contains(err.Error(), "exceeds maximum depth") {
		t.Fatalf("ValidateRouteOwnership() error = %v, want hierarchy depth rejection", err)
	}
}

func TestValidateRouteOwnershipReportsTargetsAndPathsDeterministically(t *testing.T) {
	contribution := Contribution{
		ModuleID: "content_management",
		Routes: []Route{
			{ID: "operator-z", Path: "/operator/z?view=all", Targets: []Target{TargetOperator}},
			{ID: "admin-z", Path: "/admin/z?view=all", Targets: []Target{TargetAdmin}},
			{ID: "admin-a", Path: "/admin/a?view=all", Targets: []Target{TargetAdmin}},
		},
	}

	err := ValidateRouteOwnership(contribution)
	if err == nil || !strings.Contains(err.Error(), `target "admin" path "/admin/a"`) {
		t.Fatalf("ValidateRouteOwnership() error = %v, want sorted admin /admin/a failure first", err)
	}
}

func TestRouteOwnersForTargetReturnsIsolatedMetadata(t *testing.T) {
	contribution := Contribution{
		ModuleID: "order_management",
		Routes: []Route{{
			ID:             "orders",
			Path:           "/admin/orders",
			Title:          Text{Fallback: "Orders"},
			CapabilityTags: []string{"orders:read"},
			Targets:        []Target{TargetAdmin},
			Section:        &NavSection{ID: "commerce", Label: Text{Fallback: "Commerce"}},
			RendersEntity:  &RenderedEntity{ModuleID: "order_management", EntityType: "order"},
			Page: &PageContract{
				Level:        ComponentLevelPage,
				RouteID:      "orders",
				ShellProfile: ShellProfileAdmin,
				Pattern:      PagePatternList,
				Presenter: PresenterContract{
					Contract:  "example.Presenter",
					ViewModel: "OrdersPageViewModel",
				},
				Template: TemplateContract{
					Contract: "example.Template",
					Name:     "ListPage",
					Level:    ComponentLevelTemplate,
				},
				Slots: []SlotContract{{
					ID:       "orders",
					Name:     "Orders",
					Level:    ComponentLevelOrganism,
					Contract: "example.OrdersTable",
				}},
				Storybook: &StorybookReference{
					Story: "Orders",
				},
				DesignArtifacts: &DesignArtifactSet{Manifest: "orders.json"},
			},
		}},
	}

	owners, err := RouteOwnersForTarget(contribution, TargetAdmin)
	if err != nil {
		t.Fatalf("RouteOwnersForTarget() error = %v", err)
	}
	owners[0].CapabilityTags[0] = "mutated"
	owners[0].Targets[0] = TargetPublic
	owners[0].Section.ID = "mutated"
	owners[0].RendersEntity.EntityType = "mutated"
	owners[0].Page.Slots[0].ID = "mutated"
	owners[0].Page.Storybook.Story = "mutated"
	owners[0].Page.DesignArtifacts.Manifest = "mutated"

	route := contribution.Routes[0]
	if route.CapabilityTags[0] != "orders:read" || route.Targets[0] != TargetAdmin ||
		route.Section.ID != "commerce" || route.RendersEntity.EntityType != "order" ||
		route.Page.Slots[0].ID != "orders" || route.Page.Storybook.Story != "Orders" ||
		route.Page.DesignArtifacts.Manifest != "orders.json" {
		t.Fatalf("RouteOwnersForTarget() exposed mutable input metadata: %#v", route)
	}
}

func TestValidateContributionRejectsUnsafeRouteLocations(t *testing.T) {
	tests := map[string]string{
		"scheme":            "javascript:alert(1)",
		"absolute URL":      "https://evil.example/admin",
		"network path":      "//evil.example/admin",
		"query only":        "?status=open",
		"encoded separator": "/admin/%2fusers",
		"dot segment":       "/admin/../users",
		"repeated slash":    "/admin//users",
		"backslash":         "/admin\\users",
		"surrounding space": " /admin/users",
		"control character": "/admin/users\n",
	}
	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			contribution := Contribution{ModuleID: "user_management", Routes: []Route{{
				ID: "users", Path: path, Title: Text{Fallback: "Users"}, Targets: []Target{TargetAdmin},
			}}}
			if err := ValidateContribution(contribution); err == nil {
				t.Fatalf("ValidateContribution() accepted unsafe path %q", path)
			}
		})
	}
}

func TestCanonicalNavigationLocationFailsClosedForUnsafeLocations(t *testing.T) {
	for _, location := range []string{"javascript:alert(1)", "//evil.example", "/admin/%2fusers", "?state=open"} {
		if canonical := CanonicalNavigationLocation(location); canonical != "" {
			t.Fatalf("CanonicalNavigationLocation(%q) = %q, want empty invalid sentinel", location, canonical)
		}
	}
	if canonical := CanonicalNavigationLocation("/admin/users/?state=open#details"); canonical != "/admin/users?state=open#details" {
		t.Fatalf("canonical safe location = %q", canonical)
	}
	if SameNavigationLocation("javascript:alert(1)", "javascript:alert(1)") {
		t.Fatal("invalid navigation locations must not compare equal")
	}
	if !SameNavigationLocation("/admin/users/", "/admin/users") {
		t.Fatal("canonical trailing-slash variants must compare equal")
	}
}
