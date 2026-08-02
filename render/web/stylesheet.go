// Implements: REQ-004, REQ-011.
// Per: ADR-0004, ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// stylesheet.go closes the OSS browser-delivery loop. It emits the canonical
// pk-design DTCG theme together with exactly the utility rules reachable from
// the production pk-ui renderers. Storybook, examples, and downstream hosts
// can therefore render the OSS library without borrowing proprietary CSS.

import (
	"fmt"
	"strings"

	"github.com/septagon-oss/pk-design/pkg/themes"
	"github.com/septagon-oss/styleengine"
	"github.com/septagon-oss/tw/emission"
)

// Stylesheet renders a complete browser stylesheet for theme and the OSS web
// component catalog.
func Stylesheet(theme themes.Theme, options styleengine.RenderOptions) (string, error) {
	tokenCSS, err := themes.CSSVars(theme)
	if err != nil {
		return "", fmt.Errorf("render OSS theme tokens: %w", err)
	}
	utilities, err := emission.For(ClassLists()...)
	if err != nil {
		return "", fmt.Errorf("compile OSS component utility inventory: %w", err)
	}
	sheet := emission.RoleVars().Merge(utilities)
	componentCSS, err := sheet.Render(options)
	if err != nil {
		return "", fmt.Errorf("render OSS component utilities: %w", err)
	}
	return strings.TrimSpace(tokenCSS) + "\n" + componentCSS + "\n" + progressValueCSS(), nil
}

// progressValueCSS owns the finite visual projection of the Progress
// contract's normalized integer percentage. Attribute selectors keep runtime
// markup CSP-safe and avoid requiring 101 proprietary utility classes from a
// downstream host. The renderer clamps every value to this exact 0..100 set.
func progressValueCSS() string {
	var rules strings.Builder
	for percent := 0; percent <= 100; percent++ {
		fmt.Fprintf(
			&rules,
			`[data-component=progress][data-progress-percent="%d"] [data-progress-fill=true]{width:%d%%}`,
			percent,
			percent,
		)
	}
	return rules.String()
}

// DefaultStylesheet renders the canonical OSS theme and catalog stylesheet.
func DefaultStylesheet(options styleengine.RenderOptions) (string, error) {
	return Stylesheet(themes.Default(), options)
}

// DefaultMinifiedStylesheet is the dependency-light delivery entry point for
// adapters that should not need to import styleengine merely to choose the
// canonical production option.
func DefaultMinifiedStylesheet() (string, error) {
	return DefaultStylesheet(styleengine.RenderOptions{Minify: true})
}
