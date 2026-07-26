// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package surface

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidatePageContractAcceptsCanonicalRichContract(t *testing.T) {
	if err := ValidatePageContract(validPageContract()); err != nil {
		t.Fatalf("ValidatePageContract() error = %v", err)
	}
}

func TestValidatePageContractRejectsIncompleteContracts(t *testing.T) {
	tests := map[string]struct {
		mutate func(*PageContract)
		want   string
	}{
		"route id": {
			mutate: func(page *PageContract) { page.RouteID = " " },
			want:   "route id is required",
		},
		"shell profile": {
			mutate: func(page *PageContract) { page.ShellProfile = "desktop" },
			want:   "unsupported shell profile",
		},
		"pattern": {
			mutate: func(page *PageContract) { page.Pattern = PagePatternUnknown },
			want:   "page pattern is required",
		},
		"page level": {
			mutate: func(page *PageContract) { page.Level = ComponentLevelTemplate },
			want:   "page level must be",
		},
		"presenter contract": {
			mutate: func(page *PageContract) { page.Presenter.Contract = "" },
			want:   "presenter contract is required",
		},
		"presenter view model": {
			mutate: func(page *PageContract) { page.Presenter.ViewModel = "" },
			want:   "presenter view-model is required",
		},
		"template contract": {
			mutate: func(page *PageContract) { page.Template.Contract = "" },
			want:   "template contract is required",
		},
		"template name": {
			mutate: func(page *PageContract) { page.Template.Name = "" },
			want:   "template name is required",
		},
		"template level": {
			mutate: func(page *PageContract) { page.Template.Level = ComponentLevelOrganism },
			want:   "template level must be",
		},
		"duplicate slot": {
			mutate: func(page *PageContract) { page.Slots = append(page.Slots, page.Slots[0]) },
			want:   "duplicate slot id",
		},
		"slot contract": {
			mutate: func(page *PageContract) { page.Slots[0].Contract = "" },
			want:   "slot contract is required",
		},
		"slot level": {
			mutate: func(page *PageContract) { page.Slots[0].Level = ComponentLevelUnknown },
			want:   "unsupported component level",
		},
		"empty storybook": {
			mutate: func(page *PageContract) { page.Storybook = &StorybookReference{} },
			want:   "storybook reference needs",
		},
		"empty design artifacts": {
			mutate: func(page *PageContract) { page.DesignArtifacts = &DesignArtifactSet{} },
			want:   "design artifacts are empty",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			page := validPageContract()
			test.mutate(&page)
			err := ValidatePageContract(page)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidatePageContract() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestValidateContributionRejectsPageRouteDrift(t *testing.T) {
	tests := map[string]struct {
		mutate func(*Route)
		want   string
	}{
		"shell is not a route target": {
			mutate: func(route *Route) {
				route.Page.ShellProfile = ShellProfilePublic
			},
			want: "is not declared in route targets",
		},
		"page pattern disagrees with route hint": {
			mutate: func(route *Route) {
				route.Page.Pattern = PagePatternDetail
			},
			want: "does not match page pattern",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			route := validPageRoute()
			test.mutate(&route)
			err := ValidateContribution(Contribution{
				ModuleID: "order_management",
				Routes:   []Route{route},
			})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateContribution() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestClonePageContractIsolatesNestedMetadata(t *testing.T) {
	page := validPageContract()
	cloned := ClonePageContract(&page)
	if cloned == nil {
		t.Fatal("ClonePageContract() returned nil")
	}

	cloned.Slots[0].ID = "mutated"
	cloned.Storybook.Story = "mutated"
	cloned.DesignArtifacts.Manifest = "mutated"

	if page.Slots[0].ID != "orders" ||
		page.Storybook.Story != "Orders/List" ||
		page.DesignArtifacts.Manifest != "design/orders.json" {
		t.Fatalf("ClonePageContract() exposed input metadata: %#v", page)
	}
}

func TestPageContractJSONUsesNestedCanonicalShape(t *testing.T) {
	content, err := json.Marshal(validPageContract())
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var shape map[string]json.RawMessage
	if err := json.Unmarshal(content, &shape); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	for _, required := range []string{"presenter", "template"} {
		if _, ok := shape[required]; !ok {
			t.Fatalf("canonical page JSON missing nested %q object: %s", required, content)
		}
	}
	for _, retired := range []string{
		"presenterContract",
		"presenterViewModel",
		"templateContract",
		"templateName",
		"templateLevel",
	} {
		if _, ok := shape[retired]; ok {
			t.Fatalf("canonical page JSON contains retired flattened field %q: %s", retired, content)
		}
	}
}

func validPageRoute() Route {
	page := validPageContract()
	return Route{
		ID:          "orders",
		Path:        "/admin/orders",
		Title:       Text{Fallback: "Orders"},
		NavLabel:    Text{Fallback: "Orders"},
		PagePattern: PagePatternList,
		Page:        &page,
		Targets:     []Target{TargetAdmin},
	}
}

func validPageContract() PageContract {
	return PageContract{
		Level:        ComponentLevelPage,
		RouteID:      "orders",
		ShellProfile: ShellProfileAdmin,
		Pattern:      PagePatternList,
		Presenter: PresenterContract{
			Contract:  "example.OrderPresenter",
			ViewModel: "OrdersPageViewModel",
		},
		Template: TemplateContract{
			Contract: "example.PageTemplate",
			Name:     "ListPage",
			Level:    ComponentLevelTemplate,
		},
		Slots: []SlotContract{{
			ID:       "orders",
			Name:     "Orders table",
			Level:    ComponentLevelOrganism,
			Contract: "example.OrdersTable",
			Required: true,
		}},
		Storybook: &StorybookReference{
			Package: "orders",
			Story:   "Orders/List",
		},
		DesignArtifacts: &DesignArtifactSet{
			Manifest: "design/orders.json",
		},
	}
}
