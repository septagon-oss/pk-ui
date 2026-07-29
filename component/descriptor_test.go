// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package component_test

import (
	"testing"

	"github.com/septagon-oss/pk-ui/component"
)

func TestFoundationCatalogIsValidUniqueAndDefensive(t *testing.T) {
	t.Parallel()

	catalog := component.FoundationCatalog()
	if len(catalog) != 43 {
		t.Fatalf("FoundationCatalog() has %d descriptors, want 43", len(catalog))
	}

	seen := make(map[component.ID]struct{}, len(catalog))
	for _, descriptor := range catalog {
		if err := descriptor.Validate(); err != nil {
			t.Errorf("descriptor %q is invalid: %v", descriptor.ID, err)
		}
		if _, duplicate := seen[descriptor.ID]; duplicate {
			t.Errorf("descriptor %q is duplicated", descriptor.ID)
		}
		seen[descriptor.ID] = struct{}{}
	}

	catalog[0].Owner = "mutated"
	fresh := component.FoundationCatalog()
	if fresh[0].Owner != component.FoundationOwner {
		t.Fatal("FoundationCatalog returned shared mutable storage")
	}
}

func TestWindowedCollectionIsFoundationOrganism(t *testing.T) {
	t.Parallel()

	descriptor, ok := component.FoundationDescriptor(component.TypeWindowedCollection)
	if !ok {
		t.Fatal("WindowedCollection is missing from the foundation catalog")
	}
	if descriptor.Tier != component.TierOrganism {
		t.Fatalf("WindowedCollection tier = %q, want %q", descriptor.Tier, component.TierOrganism)
	}
	if descriptor.Owner != component.FoundationOwner {
		t.Fatalf("WindowedCollection owner = %q, want %q", descriptor.Owner, component.FoundationOwner)
	}
	if err := descriptor.Validate(); err != nil {
		t.Fatalf("WindowedCollection descriptor is invalid: %v", err)
	}
}

func TestNamespacedExtensionDoesNotRequireFoundationCatalogEdit(t *testing.T) {
	t.Parallel()

	descriptor := component.Descriptor{
		ID:      "Acme/InvoiceTimeline",
		Tier:    component.TierOrganism,
		Owner:   "example.com/acme/invoicing",
		Version: "1.0.0",
	}
	contribution, err := component.NewContribution(descriptor, struct{ Renderer string }{Renderer: "web"})
	if err != nil {
		t.Fatalf("NewContribution() error = %v", err)
	}
	if contribution.Descriptor != descriptor {
		t.Fatalf("contribution descriptor = %#v, want %#v", contribution.Descriptor, descriptor)
	}
	if _, foundation := component.FoundationDescriptor(string(descriptor.ID)); foundation {
		t.Fatal("extension unexpectedly became part of the closed foundation catalog")
	}
}

func TestDescriptorValidationFailsClosed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		descriptor component.Descriptor
	}{
		{name: "missing id", descriptor: component.Descriptor{Tier: component.TierAtom, Owner: "owner"}},
		{name: "legacy snake case", descriptor: component.Descriptor{ID: "invoice_timeline", Tier: component.TierOrganism, Owner: "owner"}},
		{name: "unknown tier", descriptor: component.Descriptor{ID: "InvoiceTimeline", Tier: "widget", Owner: "owner"}},
		{name: "missing owner", descriptor: component.Descriptor{ID: "InvoiceTimeline", Tier: component.TierOrganism}},
		{name: "padded version", descriptor: component.Descriptor{ID: "InvoiceTimeline", Tier: component.TierOrganism, Owner: "owner", Version: " 1.0.0"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := test.descriptor.Validate(); err == nil {
				t.Fatalf("Descriptor.Validate() accepted %#v", test.descriptor)
			}
		})
	}
}
