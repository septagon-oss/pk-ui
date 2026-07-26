package surfacetest

// surfacetest_test.go — the conformance suite validated against Static.
//
// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

import (
	"testing"

	"github.com/septagon-oss/pk-ui/surface"
)

func sample() surface.Contribution {
	return surface.Contribution{
		ModuleID: "blog_management",
		Routes: []surface.Route{{
			ID:             "overview",
			Path:           "/admin/blog",
			Title:          surface.Text{Fallback: "Blog"},
			Order:          100,
			PagePattern:    surface.PagePatternDashboard,
			CapabilityTags: []string{"blog_management:view"},
			Targets:        []surface.Target{surface.TargetAdmin},
		}},
		Widgets: []surface.Widget{{
			ID: "recent-posts", Title: surface.Text{Fallback: "Recent Posts"},
			Kind: "list", Size: "medium", Order: 1,
		}},
		Settings: []surface.SettingsGroup{{
			ID: "general", Title: surface.Text{Fallback: "General"},
			Settings: []surface.Setting{{
				Key: "blog_management.enabled", Title: surface.Text{Fallback: "Enabled"},
				Kind: "bool", Default: "true",
			}},
		}},
	}
}

func TestStaticPassesProviderConformance(t *testing.T) {
	ProviderConformance(t, func(t *testing.T) surface.Provider {
		return Static{Contribution: sample()}
	})
}

func TestValidateContributionCatchesWidgetAndSettingShape(t *testing.T) {
	c := sample()
	c.Widgets = append(c.Widgets, surface.Widget{ID: "recent-posts", Title: surface.Text{Fallback: "Dup"}, Kind: "list"})
	if err := surface.ValidateContribution(c); err == nil {
		t.Fatal("duplicate widget id must be invalid")
	}
	c = sample()
	c.Settings[0].Settings[0].Kind = ""
	if err := surface.ValidateContribution(c); err == nil {
		t.Fatal("setting without kind must be invalid")
	}
}
