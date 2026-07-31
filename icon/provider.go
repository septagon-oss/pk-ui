// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

// Package icon owns pk-ui's provider-neutral vector-glyph contract and its
// open-source default provider. PlatformKit and client layers may prepend a
// provider for additional glyphs, but the OSS catalog always remains the
// fallback pillar.
package icon

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

const defaultViewBox = "0 0 256 256"

const (
	maxGlyphBodyBytes = 64 << 10
	maxGlyphElements  = 256
	maxGlyphDepth     = 16

	// ExtensionPriorityProduct is the normal priority for a product layer.
	ExtensionPriorityProduct = 100
	// ExtensionPriorityClient lets a client override product and OSS glyphs.
	ExtensionPriorityClient = 200
)

var allowedSVGElements = map[string]struct{}{
	"circle":   {},
	"ellipse":  {},
	"g":        {},
	"line":     {},
	"path":     {},
	"polygon":  {},
	"polyline": {},
	"rect":     {},
}

var allowedSVGAttributes = map[string]struct{}{
	"clip-rule":           {},
	"cx":                  {},
	"cy":                  {},
	"d":                   {},
	"fill":                {},
	"fill-opacity":        {},
	"fill-rule":           {},
	"height":              {},
	"opacity":             {},
	"pathLength":          {},
	"points":              {},
	"preserveAspectRatio": {},
	"r":                   {},
	"rx":                  {},
	"ry":                  {},
	"stroke":              {},
	"stroke-dasharray":    {},
	"stroke-dashoffset":   {},
	"stroke-linecap":      {},
	"stroke-linejoin":     {},
	"stroke-miterlimit":   {},
	"stroke-opacity":      {},
	"stroke-width":        {},
	"transform":           {},
	"vector-effect":       {},
	"width":               {},
	"x":                   {},
	"x1":                  {},
	"x2":                  {},
	"y":                   {},
	"y1":                  {},
	"y2":                  {},
}

// Glyph is trusted vector content returned by an application-installed icon
// provider. Body contains SVG child elements, never an outer svg element.
type Glyph struct {
	Name    string
	ViewBox string
	Body    string
}

// Validate rejects malformed or unsafe provider output at the extension seam.
func (g Glyph) Validate() error {
	if strings.TrimSpace(g.Name) == "" {
		return fmt.Errorf("icon glyph has no name")
	}
	if strings.TrimSpace(g.ViewBox) == "" {
		return fmt.Errorf("icon glyph %q has no viewBox", g.Name)
	}
	body := strings.TrimSpace(g.Body)
	if body == "" {
		return fmt.Errorf("icon glyph %q has no vector body", g.Name)
	}
	if len(body) > maxGlyphBodyBytes {
		return fmt.Errorf(
			"icon glyph %q exceeds %d bytes",
			g.Name,
			maxGlyphBodyBytes,
		)
	}
	return validateSVGFragment(g.Name, body)
}

func validateSVGFragment(name, body string) error {
	decoder := xml.NewDecoder(strings.NewReader("<svg>" + body + "</svg>"))
	depth := 0
	elements := 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("icon glyph %q is not valid XML: %w", name, err)
		}
		switch typed := token.(type) {
		case xml.StartElement:
			if depth == 0 {
				if typed.Name.Space != "" || typed.Name.Local != "svg" {
					return fmt.Errorf("icon glyph %q has an invalid wrapper", name)
				}
				depth++
				continue
			}
			if typed.Name.Space != "" {
				return fmt.Errorf(
					"icon glyph %q uses namespaced element %q",
					name,
					typed.Name.Local,
				)
			}
			if _, allowed := allowedSVGElements[typed.Name.Local]; !allowed {
				return fmt.Errorf(
					"icon glyph %q uses forbidden element %q",
					name,
					typed.Name.Local,
				)
			}
			elements++
			if elements > maxGlyphElements {
				return fmt.Errorf(
					"icon glyph %q exceeds %d vector elements",
					name,
					maxGlyphElements,
				)
			}
			depth++
			if depth > maxGlyphDepth {
				return fmt.Errorf(
					"icon glyph %q exceeds nesting depth %d",
					name,
					maxGlyphDepth,
				)
			}
			for _, attr := range typed.Attr {
				if err := validateSVGAttribute(name, attr); err != nil {
					return err
				}
			}
		case xml.EndElement:
			depth--
			if depth < 0 {
				return fmt.Errorf("icon glyph %q has unbalanced elements", name)
			}
		case xml.CharData:
			if strings.TrimSpace(string(typed)) != "" {
				return fmt.Errorf(
					"icon glyph %q contains non-vector text",
					name,
				)
			}
		case xml.Comment:
			// Comments cannot affect rendering and are safe to preserve.
		case xml.Directive, xml.ProcInst:
			return fmt.Errorf(
				"icon glyph %q contains forbidden XML instructions",
				name,
			)
		}
	}
}

func validateSVGAttribute(name string, attr xml.Attr) error {
	if attr.Name.Space != "" {
		return fmt.Errorf(
			"icon glyph %q uses namespaced attribute %q",
			name,
			attr.Name.Local,
		)
	}
	if _, allowed := allowedSVGAttributes[attr.Name.Local]; !allowed {
		return fmt.Errorf(
			"icon glyph %q uses forbidden attribute %q",
			name,
			attr.Name.Local,
		)
	}
	value := strings.ToLower(strings.TrimSpace(attr.Value))
	for _, forbidden := range []string{
		"javascript:",
		"data:",
		"url(",
		"<",
		">",
	} {
		if strings.Contains(value, forbidden) {
			return fmt.Errorf(
				"icon glyph %q attribute %q contains forbidden content",
				name,
				attr.Name.Local,
			)
		}
	}
	if (attr.Name.Local == "fill" || attr.Name.Local == "stroke") &&
		value != "" &&
		value != "none" &&
		value != "inherit" &&
		value != "currentcolor" {
		return fmt.Errorf(
			"icon glyph %q attribute %q must inherit currentColor",
			name,
			attr.Name.Local,
		)
	}
	return nil
}

// Provider resolves a canonical icon name and visual weight to trusted vector
// content. Providers must be safe for concurrent use.
type Provider interface {
	Resolve(name string, weight string) (Glyph, bool)
	Name() string
}

type providerChain struct {
	providers []Provider
}

func (chain providerChain) Name() string {
	names := make([]string, 0, len(chain.providers))
	for _, provider := range chain.providers {
		if provider != nil {
			names = append(names, provider.Name())
		}
	}
	return strings.Join(names, "+")
}

func (chain providerChain) Resolve(name string, weight string) (Glyph, bool) {
	for _, provider := range chain.providers {
		if provider == nil {
			continue
		}
		glyph, found := provider.Resolve(name, weight)
		if !found || glyph.Validate() != nil {
			continue
		}
		return glyph, true
	}
	return Glyph{}, false
}

// Chain returns one deterministic provider that resolves from left to right.
func Chain(providers ...Provider) Provider {
	filtered := make([]Provider, 0, len(providers))
	for _, provider := range providers {
		if provider != nil {
			filtered = append(filtered, provider)
		}
	}
	return providerChain{providers: filtered}
}

type providerHolder struct {
	provider Provider
}

type extensionRegistration struct {
	priority int
	provider Provider
	token    uint64
}

var (
	activeProvider    atomic.Pointer[providerHolder]
	extensionSequence atomic.Uint64
	extensionMu       sync.Mutex
	extensions        = make(map[string]extensionRegistration)
)

func init() {
	activeProvider.Store(&providerHolder{provider: BuiltinProvider()})
}

// DefaultProvider returns the process-wide provider chain used by the
// convenience renderers. The returned provider is immutable and concurrent.
func DefaultProvider() Provider {
	holder := activeProvider.Load()
	if holder == nil || holder.provider == nil {
		return BuiltinProvider()
	}
	return holder.provider
}

// RegisterExtension contributes one named product or client provider. Higher
// priorities resolve first and the OSS provider always remains the final
// fallback. Re-registering a layer is atomic; the returned function restores
// the previous value only while its own registration is still current.
func RegisterExtension(
	layer string,
	priority int,
	provider Provider,
) (restore func(), err error) {
	layer = strings.TrimSpace(layer)
	if layer == "" {
		return nil, fmt.Errorf("icon extension layer is empty")
	}
	if provider == nil {
		return nil, fmt.Errorf("icon extension %q has no provider", layer)
	}
	if strings.TrimSpace(provider.Name()) == "" {
		return nil, fmt.Errorf("icon extension %q has an unnamed provider", layer)
	}

	token := extensionSequence.Add(1)
	extensionMu.Lock()
	previous, hadPrevious := extensions[layer]
	extensions[layer] = extensionRegistration{
		priority: priority,
		provider: provider,
		token:    token,
	}
	rebuildProviderChainLocked()
	extensionMu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			extensionMu.Lock()
			defer extensionMu.Unlock()
			current, exists := extensions[layer]
			if !exists || current.token != token {
				return
			}
			if hadPrevious {
				extensions[layer] = previous
			} else {
				delete(extensions, layer)
			}
			rebuildProviderChainLocked()
		})
	}, nil
}

// InstallExtension is the compatibility helper for short-lived tools and
// tests. Applications should use RegisterExtension with an explicit layer and
// priority so OSS → product → client ordering remains inspectable.
func InstallExtension(extension Provider) (restore func()) {
	if extension == nil {
		return func() {}
	}
	sequence := extensionSequence.Add(1)
	restore, err := RegisterExtension(
		fmt.Sprintf("temporary/%020d", sequence),
		ExtensionPriorityClient+1,
		extension,
	)
	if err != nil {
		panic(err)
	}
	return restore
}

func rebuildProviderChainLocked() {
	type namedRegistration struct {
		layer string
		extensionRegistration
	}
	ordered := make([]namedRegistration, 0, len(extensions))
	for layer, registration := range extensions {
		ordered = append(ordered, namedRegistration{
			layer:                 layer,
			extensionRegistration: registration,
		})
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].priority != ordered[j].priority {
			return ordered[i].priority > ordered[j].priority
		}
		return ordered[i].layer < ordered[j].layer
	})
	providers := make([]Provider, 0, len(ordered)+1)
	for _, registration := range ordered {
		providers = append(providers, registration.provider)
	}
	providers = append(providers, BuiltinProvider())
	activeProvider.Store(&providerHolder{provider: Chain(providers...)})
}

// Resolve returns a known glyph from the active extension chain. Unknown names
// visibly degrade to the OSS question glyph and report known=false.
func Resolve(name string, weight string) (glyph Glyph, known bool) {
	glyph, known = DefaultProvider().Resolve(
		canonicalName(name),
		normalizeWeight(weight),
	)
	if known {
		return glyph, true
	}
	fallback, found := BuiltinProvider().Resolve(
		"question-mark-circle",
		"outline",
	)
	if !found {
		panic("pk-ui icon fallback glyph is unavailable")
	}
	return fallback, false
}

// RenderSVG returns standalone editable SVG markup for design delivery.
func RenderSVG(name string, weight string) string {
	glyph, _ := Resolve(name, weight)
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="` +
		html.EscapeString(glyph.ViewBox) +
		`" fill="currentColor" aria-hidden="true" data-pk-icon="` +
		html.EscapeString(canonicalName(name)) +
		`">` +
		glyph.Body +
		`</svg>`
}

func canonicalName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func normalizeWeight(value string) string {
	switch canonicalName(value) {
	case "", "regular", "outline":
		return "outline"
	default:
		return canonicalName(value)
	}
}
