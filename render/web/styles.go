// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// styles.go is the exported class surface for consumers that build component
// markup OUTSIDE these renderers — above all, scripts that create rows,
// pills, and controls in the browser and must style them with the SAME
// compiled class lists the Go renderers use. An application embeds these
// compiled strings (PlatformKit's admin does it through a JSON bridge) so
// the design system stays declared exactly once, in Go.
//
// Every list returned here is composed from the declarations in
// classlists.go, so the stylesheet derived from ClassLists() always carries
// their rules, and the collision guard covers them like any renderer path.

import "github.com/septagon-oss/tw"

// ButtonClasses returns the complete class list for a Button of the given
// variant and size — the same composition the Button renderer applies.
// Unknown values fall back to primary/md, mirroring the renderer.
func ButtonClasses(variant, size string) tw.ClassList {
	return clButtonBase.
		Merge(variantOr(clButtonVariant, variant, "primary")).
		Merge(variantOr(clButtonSize, size, "md"))
}

// BadgeClasses returns the complete class list for a Badge of the given
// variant; unknown variants fall back to primary.
func BadgeClasses(variant string) tw.ClassList {
	return clBadgeBase.
		Merge(variantOr(clBadgeVariant, variant, "primary")).
		Merge(variantOr(clBadgeSize, "md", "md"))
}

// TagClasses returns the complete class list for a Tag in its idle or
// selected state.
func TagClasses(selected bool) tw.ClassList {
	if selected {
		return clTagBase.Merge(clTagSelected)
	}
	return clTagBase.Merge(clTagIdle)
}

// TextClasses returns a text utility list by semantic color, size, and
// weight, using the same vocabularies the Text renderer accepts. Empty
// selectors contribute nothing.
func TextClasses(color, size, weight string) tw.ClassList {
	cl := tw.New()
	if c, ok := clTextColor[color]; ok {
		cl = cl.TextColor(c)
	}
	if s, ok := clTextSize[size]; ok {
		cl = cl.FontSize(s)
	}
	if w, ok := clTextWeight[weight]; ok {
		cl = cl.FontWeight(w)
	}
	return cl
}

// InlineCodeClasses returns the mono chip list for short inline values —
// identifiers, secrets, keystrokes — matching the Kbd renderer's look.
func InlineCodeClasses() tw.ClassList { return clKbd.Merge(variantOr(clKbdSize, "md", "md")) }

// HelpTextClasses returns the muted help-line list form controls use.
func HelpTextClasses() tw.ClassList { return clHelp }

// CardClasses returns the Card surface list, for elements that need the
// card treatment but different semantics than the Card renderer's <div> —
// a form's <fieldset>, for one.
func CardClasses() tw.ClassList { return clCard }

// TableClassSet names every class list the Table renderer composes, so a
// script that builds rows at runtime styles them identically.
type TableClassSet struct {
	Wrap        tw.ClassList // scrollable shell around the table
	Table       tw.ClassList // the <table> element
	Head        tw.ClassList // <thead>
	Th          tw.ClassList // header cell
	ThSort      tw.ClassList // header cell hosting a sort button
	SortButton  tw.ClassList // the sort button inside a sortable th
	Td          tw.ClassList // body cell
	TdCompact   tw.ClassList // body cell, compact tables
	TdPrimary   tw.ClassList // emphasized identity cell
	Row         tw.ClassList // body row
	RowStripe   tw.ClassList // alternate-row background for striped tables
	ActionsCell tw.ClassList // right-aligned action cluster inside a cell
	CellNote    tw.ClassList // muted annotation inside a cell
}

// TableClasses returns the Table renderer's class lists.
func TableClasses() TableClassSet {
	return TableClassSet{
		Wrap:        clTableWrap,
		Table:       clTable,
		Head:        clTableHead,
		Th:          clTableTh,
		ThSort:      clTableThSort,
		SortButton:  clTableSortBtn,
		Td:          clTableTd,
		TdCompact:   clTableTdC,
		TdPrimary:   clTableTd.Merge(clTableTdStrong),
		Row:         clTableRow,
		RowStripe:   clTableRow.Merge(clTableRowAlt),
		ActionsCell: clTableActions,
		CellNote:    clTableCellNote,
	}
}

// PaginationClassSet names the class lists both pagination modes compose.
type PaginationClassSet struct {
	Nav     tw.ClassList // the <nav> container
	Button  tw.ClassList // an idle page control
	Current tw.ClassList // the current-page control
	Label   tw.ClassList // the cursor-mode page label
}

// PaginationClasses returns the Pagination renderer's class lists.
func PaginationClasses() PaginationClassSet {
	return PaginationClassSet{
		Nav:     clPagination,
		Button:  clPageBtn.Merge(clPageIdle),
		Current: clPageBtn.Merge(clPageCur),
		Label:   clPageLabel,
	}
}
