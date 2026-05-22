package contracts

// Platform identifies the rendering target.
type Platform string

const (
	PlatformWeb     Platform = "web"
	PlatformMobile  Platform = "mobile"
	PlatformDesktop Platform = "desktop"
	PlatformA2UI    Platform = "a2ui"
)

// ComponentRenderer renders a Props struct for a specific platform.
// Each platform implements its own set of renderers.
type ComponentRenderer[P any] interface {
	Render(props P) any
	Platform() Platform
}
