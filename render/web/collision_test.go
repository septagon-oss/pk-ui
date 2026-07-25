// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// collision_test.go enforces the variant discipline structurally: in any
// class list a renderer composes onto ONE element, no CSS property may be
// declared by two different utility classes under the same variant prefix.
// Two single-class rules have equal specificity, so the emitted sheet's
// order — alphabetical, an implementation detail — would silently pick the
// winner. This is the bug class that once left secondary buttons borderless
// and selected tags unselected; the guard makes reintroducing it a test
// failure instead of a visual regression.

import (
	"strings"
	"testing"

	"github.com/septagon-oss/styleengine"
	"github.com/septagon-oss/tw"
	"github.com/septagon-oss/tw/emission"
)

// composedLists enumerates every class list the renderers and the exported
// Styles surface apply to a single element. New compositions belong here;
// the completeness check below fails if a variant map gains a key this
// table does not cover.
func composedLists(t *testing.T) map[string]tw.ClassList {
	t.Helper()
	out := map[string]tw.ClassList{}
	for v := range clButtonVariant {
		for s := range clButtonSize {
			out["button/"+v+"/"+s] = ButtonClasses(v, s)
			out["button-full/"+v+"/"+s] = ButtonClasses(v, s).Merge(clButtonFull)
		}
	}
	for v := range clBadgeVariant {
		out["badge/"+v] = BadgeClasses(v)
	}
	for v := range clAlertVariant {
		out["alert/"+v] = clAlertBase.Merge(clAlertVariant[v])
	}
	out["input/normal"] = clInput.Merge(clInputNormal)
	out["input/error"] = clInput.Merge(clInputError)
	out["tag/idle"] = TagClasses(false)
	out["tag/selected"] = TagClasses(true)
	out["empty/default"] = clEmpty.Merge(clEmptyPad)
	out["empty/compact"] = clEmpty.Merge(clEmptyCompact)
	out["empty/default-bordered"] = clEmpty.Merge(clEmptyPad).Merge(clEmptyBordered)
	out["empty/compact-bordered"] = clEmpty.Merge(clEmptyCompact).Merge(clEmptyBordered)
	out["tab/idle"] = clTab.Merge(clTabIdle)
	out["tab/active"] = clTab.Merge(clTabActive)
	out["page/idle"] = PaginationClasses().Button
	out["page/current"] = PaginationClasses().Current
	out["table/th-sort"] = TableClasses().ThSort
	out["table/td-primary"] = TableClasses().TdPrimary
	out["table/row-stripe"] = TableClasses().RowStripe
	out["card/clickable"] = clCard.Merge(clCardClickable)
	out["search/wrap"] = clSearchWrap
	out["search/input"] = clSearchInput
	return out
}

// declaredProperties renders one class alone and reports the CSS property
// names it declares, keyed by variant prefix ("" for plain, "hover",
// "focus-visible", ...). Custom properties (--*) are ignored: cooperative
// families like ring-* communicate through them by design.
func declaredProperties(t *testing.T, class string) map[string][]string {
	t.Helper()
	sheet, err := emission.Rules(class)
	if err != nil {
		t.Fatalf("emission.Rules(%q): %v", class, err)
	}
	css, err := sheet.Render(styleengine.RenderOptions{Minify: true})
	if err != nil {
		t.Fatalf("render %q: %v", class, err)
	}
	prefix := ""
	if i := strings.LastIndex(class, ":"); i >= 0 {
		prefix = class[:i]
	}
	props := map[string][]string{}
	for _, block := range strings.Split(css, "}") {
		open := strings.Index(block, "{")
		if open < 0 {
			continue
		}
		for _, decl := range strings.Split(block[open+1:], ";") {
			name, _, ok := strings.Cut(decl, ":")
			name = strings.TrimSpace(name)
			if !ok || name == "" || strings.HasPrefix(name, "--") {
				continue
			}
			props[prefix] = append(props[prefix], name)
		}
	}
	return props
}

func TestComposedListsHaveNoPropertyCollisions(t *testing.T) {
	t.Parallel()
	propCache := map[string]map[string][]string{}
	for name, list := range composedLists(t) {
		classes := strings.Fields(list.Compile())
		// owner[prefix][property] = class that already declared it
		owner := map[string]map[string]string{}
		for _, class := range classes {
			props, ok := propCache[class]
			if !ok {
				props = declaredProperties(t, class)
				propCache[class] = props
			}
			for prefix, names := range props {
				if owner[prefix] == nil {
					owner[prefix] = map[string]string{}
				}
				for _, prop := range names {
					if prior, taken := owner[prefix][prop]; taken && prior != class {
						t.Errorf("%s: %q and %q both declare %q at prefix %q — stylesheet order would decide the winner",
							name, prior, class, prop, prefix)
						continue
					}
					owner[prefix][prop] = class
				}
			}
		}
	}
}
