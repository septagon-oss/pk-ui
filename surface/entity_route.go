// Package surface defines framework-neutral UI contracts shared by capability
// producers, presentation adapters, and composition tooling.
package surface

// entity_route.go — neutral entity-route publication contract.

// EntityRoute describes the facts a presentation adapter needs to make a CRUD
// collection discoverable to pickers and admin shells. It contains no renderer,
// framework, transport, or process-global registry type.
type EntityRoute struct {
	ModuleID string `json:"moduleId"`
	// EntityName is the canonical snake_case entity identity.
	EntityName    string `json:"entityName"`
	TableName     string `json:"tableName"`
	CollectionURL string `json:"collectionUrl"`
}

// EntityRouteRegistrar is implemented by a presentation adapter. Backend
// entity composition publishes intent through this seam; the selected
// frontend decides how and where to index it. Implementations must be safe for
// concurrent use and treat repeated registration of the same module/entity as
// an idempotent upsert.
type EntityRouteRegistrar interface {
	RegisterEntityRoute(route EntityRoute)
}
