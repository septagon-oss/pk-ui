// Validates: REQ-011.
// Per: ADR-0007, ADR-0031.
// Discipline: C-14.

package surface

import (
	"reflect"
	"testing"
)

func TestCanonicalizeAdminSurfaceContributionMaterializesPageContract(t *testing.T) {
	original := Contribution{
		ModuleID: " example ",
		Routes: []Route{{
			ID:          "overview",
			PagePattern: PagePatternDashboard,
		}},
	}

	got := CanonicalizeAdminSurfaceContribution(original)
	if got.ModuleID != "example" {
		t.Fatalf("ModuleID = %q, want example", got.ModuleID)
	}
	if original.Routes[0].Page != nil {
		t.Fatal("canonicalization mutated the caller")
	}
	page := got.Routes[0].Page
	if page == nil {
		t.Fatal("expected a materialized page contract")
	}
	if page.RouteID != "overview" ||
		page.ShellProfile != ShellProfileAdmin ||
		page.Pattern != PagePatternDashboard ||
		page.Level != ComponentLevelPage ||
		page.Presenter.Contract != adminPagePresenterContract ||
		page.Presenter.ViewModel != "OverviewPageViewModel" ||
		page.Template.Contract != adminPageTemplateContract ||
		page.Template.Name != "DashboardPage" ||
		page.Template.Level != ComponentLevelTemplate {
		t.Fatalf("unexpected page contract: %#v", page)
	}
}

func TestCanonicalizeAdminSurfaceContributionLeavesAppRoutesUntouched(t *testing.T) {
	original := Contribution{
		ModuleID: "mixed",
		Routes: []Route{
			{
				ID:          "admin",
				PagePattern: PagePatternDashboard,
				Targets:     []Target{TargetAdmin},
			},
			{
				ID:          "account",
				PagePattern: PagePatternSettings,
				Targets:     []Target{TargetApp},
			},
		},
	}

	got := CanonicalizeAdminSurfaceContribution(original)
	if got.Routes[0].Page == nil || got.Routes[0].Page.ShellProfile != ShellProfileAdmin {
		t.Fatalf("admin route was not canonicalized: %#v", got.Routes[0])
	}
	if got.Routes[1].Page != nil {
		t.Fatalf("app route acquired an admin page contract: %#v", got.Routes[1].Page)
	}
	if original.Routes[0].Page != nil || original.Routes[1].Page != nil {
		t.Fatal("canonicalization mutated the caller")
	}
}

func TestCanonicalizeAdminRoutePreservesExplicitMetadata(t *testing.T) {
	route := Route{
		ID: "review",
		Page: &PageContract{
			Pattern:   PagePatternReview,
			Presenter: PresenterContract{ViewModel: "CustomReview"},
			Template:  TemplateContract{Name: "CustomTemplate"},
		},
	}

	got := CanonicalizeAdminRoute(route)
	if got.PagePattern != PagePatternReview {
		t.Fatalf("PagePattern = %q, want review", got.PagePattern)
	}
	if got.Page.Presenter.ViewModel != "CustomReview" || got.Page.Template.Name != "CustomTemplate" {
		t.Fatalf("explicit metadata was overwritten: %#v", got.Page)
	}
}

func TestNamespacedSurfaceContributionIsIdempotentAndIsolated(t *testing.T) {
	original := Contribution{
		ModuleID: "catalog",
		Routes: []Route{
			{
				ID: "home",
				Page: &PageContract{
					RouteID: "home",
					Slots:   []SlotContract{{ID: "body"}},
				},
			},
			{ID: "detail", ParentRouteID: "home"},
		},
	}

	got := NamespacedSurfaceContribution(original)
	twice := NamespacedSurfaceContribution(got)
	if !reflect.DeepEqual(got, twice) {
		t.Fatalf("namespacing is not idempotent:\nfirst %#v\nsecond %#v", got, twice)
	}
	if original.Routes[0].ID != "home" || original.Routes[0].Page.RouteID != "home" {
		t.Fatalf("namespacing mutated caller: %#v", original)
	}
	if got.Routes[0].ID != "catalog.home" ||
		got.Routes[0].Page.RouteID != "catalog.home" ||
		got.Routes[1].ParentRouteID != "catalog.home" {
		t.Fatalf("unexpected namespaced graph: %#v", got.Routes)
	}
}

func TestNamespacedSurfaceContributionRejectsForeignNamespace(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected foreign route namespace to panic")
		}
	}()
	NamespacedSurfaceContribution(Contribution{
		ModuleID: "catalog",
		Routes:   []Route{{ID: "foreign.home"}},
	})
}
