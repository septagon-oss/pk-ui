// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_test.go proves the OSS catalog is closed over real renderers,
// stable identities, native semantic blueprints, and executable fixtures.
// Browser pixels remain covered by Storybook visual baselines; this boundary
// deliberately rejects raster snapshots as design structure.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-design/pkg/blueprint"
	uicomponent "github.com/septagon-oss/pk-ui/component"
)

func TestOSSDeliveryCatalogIsCompleteNativeAndExecutable(t *testing.T) {
	t.Parallel()

	catalog := OSSDeliveryCatalog()
	if got, want := len(catalog), 32; got != want {
		t.Fatalf("catalog has %d entries, want %d", got, want)
	}

	identities := make(map[string]struct{}, len(catalog))
	stableIDs := make(map[string]struct{}, len(catalog))
	for _, definition := range catalog {
		if err := definition.Validate(); err != nil {
			t.Fatalf("%s: invalid definition: %v", definition.Identity.ID, err)
		}
		identity := string(definition.Identity.ID)
		if _, duplicate := identities[identity]; duplicate {
			t.Fatalf("duplicate wire identity %q", identity)
		}
		identities[identity] = struct{}{}
		if _, duplicate := stableIDs[definition.Contract.ID]; duplicate {
			t.Fatalf("duplicate stable design id %q", definition.Contract.ID)
		}
		stableIDs[definition.Contract.ID] = struct{}{}
		assertNativeDeliveryNode(t, identity, definition.Design.Root)

		for _, example := range definition.Examples {
			node, err := definition.Render(example.Props)
			if err != nil {
				t.Fatalf("%s/%s: render: %v", identity, example.Name, err)
			}
			var rendered bytes.Buffer
			if err := node.Render(&rendered); err != nil {
				t.Fatalf("%s/%s: write HTML: %v", identity, example.Name, err)
			}
			if strings.TrimSpace(rendered.String()) == "" {
				t.Fatalf("%s/%s: rendered empty HTML", identity, example.Name)
			}
		}
	}

	if _, ok := identities["DataManagementPage"]; !ok {
		t.Fatal("catalog has no complete OSS solution page")
	}
}

func TestOSSDeliveryCatalogComposesDownTheAtomicGraph(t *testing.T) {
	t.Parallel()

	catalog := OSSDeliveryCatalog()
	tiers := make(map[string]uicomponent.Tier, len(catalog))
	for _, definition := range catalog {
		tiers[string(definition.Identity.ID)] = definition.Identity.Tier
	}
	rank := map[uicomponent.Tier]int{
		uicomponent.TierAtom:     1,
		uicomponent.TierMolecule: 2,
		uicomponent.TierOrganism: 3,
		uicomponent.TierTemplate: 4,
		uicomponent.TierPage:     5,
	}
	for _, definition := range catalog {
		source := string(definition.Identity.ID)
		for _, slot := range definition.Contract.Slots {
			for _, allowedType := range slot.AllowedTypes {
				targetTier, exists := tiers[allowedType]
				if !exists {
					t.Errorf("%s slot %s allows unknown type %q", source, slot.Name, allowedType)
					continue
				}
				if rank[targetTier] > rank[definition.Identity.Tier] {
					t.Errorf(
						"%s (%s) slot %s allows higher-tier %s (%s)",
						source,
						definition.Identity.Tier,
						slot.Name,
						allowedType,
						targetTier,
					)
				}
			}
		}
		walkDeliveryNodes(definition.Design.Root, func(node *blueprint.Node) {
			if node.Kind != blueprint.NodeInstance || strings.HasPrefix(node.Text, "icon-") {
				return
			}
			targetTier, exists := tiers[node.Text]
			if !exists {
				t.Errorf("%s references unknown native instance %q", source, node.Text)
				return
			}
			if rank[targetTier] > rank[definition.Identity.Tier] {
				t.Errorf(
					"%s (%s) references higher-tier %s (%s)",
					source,
					definition.Identity.Tier,
					node.Text,
					targetTier,
				)
			}
		})
	}
}

func assertNativeDeliveryNode(t *testing.T, component string, node *blueprint.Node) {
	t.Helper()
	if node == nil {
		t.Fatalf("%s: nil native root", component)
	}
	walkDeliveryNodes(node, func(current *blueprint.Node) {
		if current.Name == "" {
			t.Errorf("%s: native %s node has no stable name", component, current.Kind)
		}
		for key, value := range current.Props {
			if strings.Contains(strings.ToLower(key), "screenshot") ||
				strings.Contains(strings.ToLower(key), "snapshot") {
				t.Errorf("%s/%s: forbidden raster-state property %q", component, current.Name, key)
			}
			if text, ok := value.(string); ok &&
				strings.HasPrefix(strings.ToLower(strings.TrimSpace(text)), "data:image/") {
				t.Errorf("%s/%s: embedded raster payload", component, current.Name)
			}
		}
	})
}

func walkDeliveryNodes(root *blueprint.Node, visit func(*blueprint.Node)) {
	if root == nil {
		return
	}
	visit(root)
	for index := range root.Children {
		walkDeliveryNodes(&root.Children[index], visit)
	}
}
