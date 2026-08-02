// Package contracts defines platform-agnostic Props schemas for all UI components.
//
// Contracts are data schemas (Props structs), not behavior interfaces. Existing
// builders gain an optional Props() method. Build() remains the primary API.
// This package is purely additive — zero breaking changes.
//
// Consumers:
//   - Go developers: import Props types for type-safe construction
//   - registry.ComponentDefinition: LLM-facing JSON schema (A2UI/MCP)
//   - Future renderers (iOS, Android): implement rendering from Props
package contracts

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
// ComponentProps is the base set of properties shared by all components.
type ComponentProps struct {
	// ID is supplied by the component-tree transport, not authored as a
	// component property in design manifests or A2UI schemas.
	ID string `json:"id,omitempty" delivery:"internal"`
	// Class is a renderer escape hatch for trusted Go composition. It is not a
	// portable design-system property and is therefore excluded from delivery.
	Class    string `json:"class,omitempty" delivery:"internal"`
	Disabled bool   `json:"disabled,omitempty"`
	Hidden   bool   `json:"hidden,omitempty"`
	// Attrs is restricted to trusted direct-render callers; portable contracts
	// expose explicit typed properties instead of arbitrary HTML attributes.
	Attrs map[string]string `json:"attrs,omitempty" delivery:"internal"`
}

// HTMXProps contains HTMX-specific properties for server-driven interactions.
type HTMXProps struct {
	Get         string `json:"hx-get,omitempty"`
	Post        string `json:"hx-post,omitempty"`
	Put         string `json:"hx-put,omitempty"`
	Patch       string `json:"hx-patch,omitempty"`
	Delete      string `json:"hx-delete,omitempty"`
	Target      string `json:"hx-target,omitempty"`
	Swap        string `json:"hx-swap,omitempty"`
	Trigger     string `json:"hx-trigger,omitempty"`
	Include     string `json:"hx-include,omitempty"`
	Confirm     string `json:"hx-confirm,omitempty"`
	Ext         string `json:"hx-ext,omitempty"`
	Indicator   string `json:"hx-indicator,omitempty"`
	DisabledElt string `json:"hx-disabled-elt,omitempty"`
	Vals        string `json:"hx-vals,omitempty"`
	PushURL     string `json:"hx-push-url,omitempty"`
	Select      string `json:"hx-select,omitempty"`
	Boost       bool   `json:"hx-boost,omitempty"`
	Disable     bool   `json:"hx-disable,omitempty"`
}

// ARIAProps contains accessibility attributes.
type ARIAProps struct {
	Role        string `json:"role,omitempty"`
	Label       string `json:"aria-label,omitempty"`
	LabelledBy  string `json:"aria-labelledby,omitempty"`
	DescribedBy string `json:"aria-describedby,omitempty"`
	Expanded    *bool  `json:"aria-expanded,omitempty"`
	Controls    string `json:"aria-controls,omitempty"`
	Live        string `json:"aria-live,omitempty"`
}

// Buildable is the interface for all component builders that can produce a rendered node.
// The concrete return type depends on the rendering backend (gomponents Node, UIKit View, etc.).
type Buildable interface {
	Build() any
}

// PropsExtractable is an optional interface that builders can implement to expose
// their accumulated state as a Props struct. This enables inspection, serialization,
// and cross-platform rendering from the same builder configuration.
type PropsExtractable[P any] interface {
	Buildable
	Props() P
}
