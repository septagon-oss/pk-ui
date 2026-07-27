package surfacetest

import (
	"testing"

	"github.com/septagon-oss/pk-ui/surface"
)

func TestEntityRouteMemPassesConformance(t *testing.T) {
	EntityRouteRegistrarConformance(t, func(t *testing.T) (
		surface.EntityRouteRegistrar,
		EntityRouteSnapshot,
	) {
		t.Helper()
		driver := &EntityRouteMem{}
		return driver, driver.Snapshot
	})
}
