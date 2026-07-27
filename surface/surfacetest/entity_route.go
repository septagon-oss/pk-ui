// Package surfacetest provides reusable conformance suites and reference
// adapters for the public surface contracts.
package surfacetest

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/septagon-oss/pk-ui/surface"
)

// EntityRouteSnapshot returns the routes currently indexed by a driver.
type EntityRouteSnapshot func() []surface.EntityRoute

// EntityRouteRegistrarConformance runs the behavior required from every
// EntityRouteRegistrar adapter.
func EntityRouteRegistrarConformance(
	t *testing.T,
	newDriver func(t *testing.T) (surface.EntityRouteRegistrar, EntityRouteSnapshot),
) {
	t.Helper()

	t.Run("preserves the complete route contract", func(t *testing.T) {
		registrar, snapshot := newDriver(t)
		want := surface.EntityRoute{
			ModuleID:      "catalog_management",
			EntityName:    "catalog_item",
			TableName:     "catalog_items",
			CollectionURL: "/api/v1/crud/catalog_items",
		}
		registrar.RegisterEntityRoute(want)
		got := snapshot()
		if len(got) != 1 || got[0] != want {
			t.Fatalf("snapshot = %#v, want %#v", got, []surface.EntityRoute{want})
		}
	})

	t.Run("repeated identity is an idempotent upsert", func(t *testing.T) {
		registrar, snapshot := newDriver(t)
		first := surface.EntityRoute{
			ModuleID:      "catalog_management",
			EntityName:    "catalog_item",
			TableName:     "catalog_items",
			CollectionURL: "/api/v1/crud/catalog_items",
		}
		updated := first
		updated.CollectionURL = "/api/v2/catalog/items"
		registrar.RegisterEntityRoute(first)
		registrar.RegisterEntityRoute(updated)
		got := snapshot()
		if len(got) != 1 || got[0] != updated {
			t.Fatalf("snapshot = %#v, want one updated route %#v", got, updated)
		}
	})

	t.Run("supports concurrent independent registrations", func(t *testing.T) {
		registrar, snapshot := newDriver(t)
		var wait sync.WaitGroup
		for index := range 8 {
			wait.Go(func() {
				registrar.RegisterEntityRoute(surface.EntityRoute{
					ModuleID:      "catalog_management",
					EntityName:    fmt.Sprintf("entity_%d", index),
					TableName:     fmt.Sprintf("entities_%d", index),
					CollectionURL: fmt.Sprintf("/api/v1/crud/entities_%d", index),
				})
			})
		}
		wait.Wait()
		if got := snapshot(); len(got) != 8 {
			t.Fatalf("concurrent snapshot len = %d, want 8: %#v", len(got), got)
		}
	})
}

// EntityRouteMem is an in-memory reference EntityRouteRegistrar.
type EntityRouteMem struct {
	mu     sync.RWMutex
	routes map[string]surface.EntityRoute
}

// RegisterEntityRoute implements surface.EntityRouteRegistrar.
func (m *EntityRouteMem) RegisterEntityRoute(route surface.EntityRoute) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.routes == nil {
		m.routes = make(map[string]surface.EntityRoute)
	}
	m.routes[route.ModuleID+"\x00"+route.EntityName] = route
}

// Snapshot returns a deterministic copy of registered routes.
func (m *EntityRouteMem) Snapshot() []surface.EntityRoute {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]surface.EntityRoute, 0, len(m.routes))
	for _, route := range m.routes {
		result = append(result, route)
	}
	slices.SortFunc(result, func(left, right surface.EntityRoute) int {
		if left.ModuleID != right.ModuleID {
			return compare(left.ModuleID, right.ModuleID)
		}
		return compare(left.EntityName, right.EntityName)
	})
	return result
}

func compare(left, right string) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}
