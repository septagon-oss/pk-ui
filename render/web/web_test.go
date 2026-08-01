// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// web_test.go closes the loop this package promises: every class a renderer
// can emit is declared in classlists.go, every declared class resolves to a
// CSS rule through tw/emission, and the rendered HTML carries the
// accessibility structure the contracts imply. The golden file makes markup
// changes reviewable diffs.

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/layouts"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	"github.com/septagon-oss/pk-ui/contracts/organisms"
	"github.com/septagon-oss/tw/emission"
)

// gallery renders one representative instance of every implemented component,
// exercising variants, sizes, states, errors, and edge branches. It is the
// single source for the golden file AND the class-closure test, so a renderer
// cannot emit a class the tests never see.
func gallery() []g.Node {
	yes := true
	_ = yes
	return []g.Node{
		Button(atoms.ButtonProps{Label: "Save", Variant: "primary", Tone: "neutral", Size: "md"}),
		Button(atoms.ButtonProps{Label: "Cancel", Variant: "secondary", Tone: "neutral", Size: "sm"}),
		Button(atoms.ButtonProps{Label: "Delete", Variant: "primary", Tone: "danger", Size: "xs"}),
		Button(atoms.ButtonProps{Label: "Ghost", Variant: "ghost", Tone: "neutral", Size: "lg"}),
		Button(atoms.ButtonProps{Label: "Docs", Variant: "link", Tone: "neutral", Size: "xl"}),
		Button(atoms.ButtonProps{Label: "Outline", Variant: "outline", Tone: "neutral", Size: "2xl", FullWidth: true}),
		Button(atoms.ButtonProps{Label: "Info", Variant: "primary", Tone: "info"}),
		Button(atoms.ButtonProps{Label: "Warn", Variant: "primary", Tone: "warning"}),
		Button(atoms.ButtonProps{Label: "OK", Variant: "primary", Tone: "success"}),
		Button(atoms.ButtonProps{Label: "Err", Variant: "primary", Tone: "danger"}),
		Button(atoms.ButtonProps{Label: "Busy", Loading: true}),
		buttonWithSlots(
			atoms.ButtonProps{Label: "Iconed", Variant: "primary", Tone: "neutral", Size: "md"},
			nil,
			[]g.Node{Icon(atoms.IconProps{Name: "plus", Size: "md", Tone: "neutral"})},
		),

		Badge(atoms.BadgeProps{Label: "Default"}),
		Badge(atoms.BadgeProps{Label: "New", Variant: "primary", Tone: "brand", Dot: true}),
		Badge(atoms.BadgeProps{Label: "OK", Tone: "success"}),
		badgeWithSlots(
			atoms.BadgeProps{Label: "Careful", Tone: "warning"},
			[]g.Node{Icon(atoms.IconProps{Name: "alert", Size: "sm", Tone: "warning"})},
			nil,
		),
		Badge(atoms.BadgeProps{Label: "Bad", Tone: "danger"}),
		Badge(atoms.BadgeProps{Label: "FYI", Tone: "info"}),
		Badge(atoms.BadgeProps{Label: "Two", Variant: "secondary"}),
		Badge(atoms.BadgeProps{Label: "Outlined", Variant: "outline"}),

		Alert(atoms.AlertProps{Title: "Saved", Message: "All good.", Tone: "success"}),
		alertWithSlots(
			atoms.AlertProps{Message: "Careful now.", Tone: "warning"},
			[]g.Node{Icon(atoms.IconProps{Name: "alert", Size: "md", Tone: "warning"})},
			nil,
		),
		Alert(atoms.AlertProps{Message: "That failed.", Tone: "danger"}),
		Alert(atoms.AlertProps{Message: "Broken.", Tone: "danger"}),
		Alert(atoms.AlertProps{Message: "Heads up.", Tone: "info"}),

		Input(atoms.InputProps{Name: "email", Type: "email", Label: "Email",
			Placeholder: "you@example.test", HelpText: "We never share it.", Required: true}),
		Input(atoms.InputProps{Name: "slug", Label: "Slug", Value: "hello",
			Error: "Already taken.", Pattern: "[a-z-]+"}),
		Input(atoms.InputProps{Name: "quiet"}),

		Textarea(atoms.TextareaProps{Name: "body", Label: "Body", Rows: 6,
			HelperText: "Markdown is fine.", MaxLength: 500}),
		Textarea(atoms.TextareaProps{Name: "bad", Label: "Bad", ErrorMessage: "Too long."}),

		Checkbox(atoms.CheckboxProps{Name: "agree", Label: "I agree", Required: true}),
		Checkbox(atoms.CheckboxProps{Name: "done", Label: "Done", Checked: true}),

		Label(atoms.LabelProps{Text: "Standalone", For: "x", Required: true}),

		Text(atoms.TextProps{Content: "Plain body copy.", Size: "sm", Color: "muted"}),
		Text(atoms.TextProps{Content: "Loud.", Size: "xl", Weight: "bold", Color: "brand"}),
		Text(atoms.TextProps{Content: "Cut off eventually.", Truncate: true}),

		Heading(atoms.HeadingProps{Text: "Page title", Level: 1}),
		Heading(atoms.HeadingProps{Text: "Section", Level: 2, Anchor: "section"}),
		Heading(atoms.HeadingProps{Text: "Sub", Level: 3}),
		Heading(atoms.HeadingProps{Text: "Minor", Level: 4}),
		Heading(atoms.HeadingProps{Text: "Small", Level: 5}),
		Heading(atoms.HeadingProps{Text: "Eyebrow", Level: 6}),

		Divider(atoms.DividerProps{}),
		Divider(atoms.DividerProps{Text: "Or continue with"}),
		Divider(atoms.DividerProps{Orientation: "vertical"}),

		Spinner(atoms.SpinnerProps{Size: "sm", Tone: "brand"}),
		Spinner(atoms.SpinnerProps{Size: "md", Tone: "info", Label: "Fetching"}),
		Spinner(atoms.SpinnerProps{Size: "lg", Tone: "success"}),

		Skeleton(atoms.SkeletonProps{}),
		Skeleton(atoms.SkeletonProps{Shape: "block", Size: "sm"}),
		Skeleton(atoms.SkeletonProps{Shape: "block", Size: "lg"}),
		Skeleton(atoms.SkeletonProps{Shape: "text"}),
		Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3, Size: "sm"}),
		Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 2, Size: "lg"}),
		Skeleton(atoms.SkeletonProps{Shape: "circle", Size: "sm"}),
		Skeleton(atoms.SkeletonProps{Shape: "circle"}),
		Skeleton(atoms.SkeletonProps{Shape: "circle", Size: "lg"}),
		DeferredSlot(atoms.DeferredSlotProps{
			HTMXProps: contracts.HTMXProps{Get: "/fragments/activity"},
		}, Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3})),
		TableSkeleton(molecules.TableSkeletonProps{}),
		TableSkeleton(molecules.TableSkeletonProps{Columns: 2, Rows: 5, Compact: true}),
		CardSkeleton(molecules.CardSkeletonProps{}),
		CardSkeleton(molecules.CardSkeletonProps{Lines: 5}),

		emptyStateWithSlots(
			atoms.EmptyStateProps{
				Title: "No tenants yet", Description: "Create the first tenant to get started.", Bordered: true,
			},
			nil,
			[]g.Node{Link(atoms.LinkProps{Label: "New tenant", Href: "/admin/tenants/new"})},
		),
		EmptyState(atoms.EmptyStateProps{Title: "Empty", Compact: true}),

		Kbd(atoms.KbdProps{Keys: []string{"Ctrl", "K"}}),

		Link(atoms.LinkProps{Label: "Internal", Href: "/docs"}),
		Link(atoms.LinkProps{Label: "External", Href: "https://example.test", External: true}),

		Tag(atoms.TagProps{Label: "beta", Tone: "neutral"}),
		Tag(atoms.TagProps{Label: "chosen", Tone: "neutral", Selected: true, Removable: true, OnRemoveURL: "/tags/1"}),

		Stack(layouts.StackProps{Gap: "2", Align: "start"}, g.Text("a"), g.Text("b")),
		Flex(layouts.FlexProps{Direction: "row", Gap: "4", Align: "center",
			Justify: "between", Wrap: true}, g.Text("l"), g.Text("r")),
		Grid(layouts.GridProps{Columns: "3", Gap: "6"}, g.Text("1"), g.Text("2"), g.Text("3")),
		Container(layouts.ContainerProps{MaxWidth: "4xl"}, g.Text("content")),

		Table(molecules.TableProps{
			Columns: []molecules.TableColumn{{Key: "name", Label: "Name"}, {Key: "role", Label: "Role"}},
			Rows: []molecules.TableRow{
				{ID: "u1", Cells: map[string]any{"name": "Ada", "role": "admin"}},
				{ID: "u2", Cells: map[string]any{"name": "Lin", "role": 7}},
			},
		}),
		Table(molecules.TableProps{
			Columns:   []molecules.TableColumn{{Key: "a", Label: "A"}},
			EmptyText: "No rows.", Compact: true,
		}),
		Table(molecules.TableProps{
			Sortable: true, Striped: true, Selectable: true,
			Columns: []molecules.TableColumn{
				{Key: "name", Label: "Name", Sortable: true, Primary: true},
				{Key: "count", Label: "Count", Sortable: true, Align: "right"},
				{Key: "note", Label: "Note"},
			},
			Rows: []molecules.TableRow{
				{ID: "r1", Cells: map[string]any{"name": "First", "count": 3, "note": "odd"}},
				{ID: "r2", Cells: map[string]any{"name": "Second", "count": 1, "note": "even"}},
				{ID: "r3", Cells: map[string]any{"name": "Third", "count": 2}},
			},
		}),

		Card(molecules.CardProps{Title: "Plain card", Description: "With copy."}),
		Card(molecules.CardProps{Title: "Go somewhere", Clickable: true, Href: "/detail"}),

		Breadcrumb(molecules.BreadcrumbProps{Items: []molecules.BreadcrumbItem{
			{Label: "Home", Href: "/"}, {Label: "Tenants", Href: "/tenants"}, {Label: "Acme"},
		}}),

		Pagination(molecules.PaginationProps{CurrentPage: 5, TotalPages: 12, BaseURL: "/rows"}),
		Pagination(molecules.PaginationProps{CurrentPage: 1, TotalPages: 2, BaseURL: "/few"}),
		Pagination(molecules.PaginationProps{CurrentPage: 1}),
		Pagination(molecules.PaginationProps{CurrentPage: 4}),

		Select(atoms.SelectProps{Name: "kind", Label: "Kind", Required: true,
			Options: []atoms.SelectOption{{Label: "Post", Value: "post"}, {Label: "Page", Value: "page"}},
			Value:   "post", HelpText: "What the entry renders as."}),
		Select(atoms.SelectProps{Name: "state", Label: "State", Placeholder: "Any state",
			Options: []atoms.SelectOption{{Label: "Draft", Value: "draft"}},
			Error:   "Pick a state."}),
		SearchBar(molecules.SearchBarProps{Label: "Search tenants", Placeholder: "Filter this page…"}),

		DataGrid(
			organisms.DataGridProps{},
			DataGridSlots{
				Search: SearchBar(molecules.SearchBarProps{
					Label: "Search rows", Placeholder: "Filter…",
				}),
				Filters: []g.Node{
					Select(atoms.SelectProps{
						Name:  "state",
						Label: "State",
						Options: []atoms.SelectOption{{
							Label: "Open",
							Value: "open",
						}},
					}),
				},
				Actions: []g.Node{
					Button(atoms.ButtonProps{
						Label: "Refresh", Variant: "secondary", Tone: "neutral",
					}),
					Link(atoms.LinkProps{Label: "New row", Href: "/rows/new"}),
				},
				Table: Table(molecules.TableProps{Sortable: true,
					Columns: []molecules.TableColumn{{Key: "name", Label: "Name", Sortable: true, Primary: true}},
					Rows:    []molecules.TableRow{{ID: "g1", Cells: map[string]any{"name": "Grid row"}}}},
				),
				Status: []g.Node{
					Text(atoms.TextProps{Content: "1 record on this page", Color: "muted", Size: "sm"}),
				},
				Pagination: Pagination(molecules.PaginationProps{CurrentPage: 1}),
			},
		),

		Tabs(molecules.TabsProps{ActiveTab: "b", Items: []molecules.TabItem{
			{Key: "a", Label: "First", URL: "/tab/a"}, {Key: "b", Label: "Second"},
		}}),
	}
}

func TestIconRendersAccessibleEditableSVG(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Icon(atoms.IconProps{
		Name:      "check",
		Size:      "lg",
		Tone:      "success",
		AriaLabel: "Complete",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<svg`,
		`viewBox="0 0 256 256"`,
		`fill="currentColor"`,
		`focusable="false"`,
		`data-pk-icon="check"`,
		`role="img"`,
		`aria-label="Complete"`,
		`<path`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Icon output is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "<span") || strings.Contains(html, "aria-hidden") {
		t.Errorf("labelled Icon has invalid wrapper or a11y state: %s", html)
	}
}

func TestIconUnknownNameUsesVisibleDecorativeFallback(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := Icon(atoms.IconProps{Name: "not-a-real-icon"}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-pk-icon="not-a-real-icon"`,
		`data-pk-icon-fallback="true"`,
		`aria-hidden="true"`,
		`<path`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("fallback Icon output is missing %q: %s", fragment, html)
		}
	}
}

func TestBreadcrumbOwnsTruncationAndIconRendering(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Breadcrumb(molecules.BreadcrumbProps{
		MaxItems: 3,
		Items: []molecules.BreadcrumbItem{
			{Label: "Home", Href: "/", Icon: "home"},
			{Label: "Workspace", Href: "/workspace"},
			{Label: "Settings", Href: "/settings"},
			{Label: "Profile", Current: true},
		},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`data-pk-icon="home"`,
		`href="/"`,
		`>…</li>`,
		`href="/settings"`,
		`aria-current="page">Profile`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Breadcrumb output is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "Workspace") {
		t.Errorf("Breadcrumb did not collapse its middle items: %s", html)
	}
}

func TestDividerPreservesLabelAndTrustedBaseAttributes(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Divider(atoms.DividerProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "auth-divider",
			Class: "my-divider",
			Attrs: map[string]string{"data-auth-divider": "true"},
		},
		Text: "Or continue with",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="auth-divider"`,
		`class="flex items-center gap-4 my-divider"`,
		`data-auth-divider="true"`,
		`role="presentation"`,
		`aria-hidden="true"`,
		`Or continue with`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("labelled Divider output is missing %q: %s", fragment, html)
		}
	}
}

func TestSearchBarOwnsProgressiveEnhancementContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := SearchBar(molecules.SearchBarProps{
		ComponentProps: contracts.ComponentProps{
			ID:       "workspace-search",
			Disabled: true,
			Attrs:    map[string]string{"data-owner": "catalog"},
		},
		HTMXProps: contracts.HTMXProps{
			Target: "#results",
			Swap:   "innerHTML",
		},
		SearchURL:    "/search",
		Label:        "Search workspaces",
		Placeholder:  "Find a workspace",
		Name:         "q",
		Value:        "existing",
		Instant:      true,
		ShowClear:    true,
		ShowShortcut: true,
		DebounceMS:   450,
		MinChars:     2,
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`id="workspace-search"`,
		`data-owner="catalog"`,
		`data-component="searchbar"`,
		`data-controller="searchbar"`,
		`name="q"`,
		`value="existing"`,
		`disabled`,
		`hx-get="/search"`,
		`hx-trigger="keyup changed delay:450ms, search"`,
		`hx-target="#results"`,
		`hx-swap="innerHTML"`,
		`hx-indicator="#q-indicator"`,
		`data-search-min-chars="2"`,
		`data-searchbar-clear-button="true"`,
		`data-searchbar-shortcut="true"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("SearchBar output is missing %q: %s", fragment, html)
		}
	}
}

func TestTableWithSlotsOwnsRichCellsAndServerSorting(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := TableWithSlots(molecules.TableProps{
		HTMXProps: contracts.HTMXProps{
			Target: "#users-grid",
			Swap:   "outerHTML",
		},
		Sortable:   true,
		Selectable: true,
		Columns: []molecules.TableColumn{
			{Key: "name", Label: "Name", Sortable: true, Width: "12rem"},
			{Key: "actions", Label: "Actions"},
		},
		Rows: []molecules.TableRow{{
			ID: "user-1",
			Cells: map[string]any{
				"name":    "Ada",
				"actions": "plain fallback",
			},
		}},
	}, TableSlots{
		Cell: func(row molecules.TableRow, column molecules.TableColumn) g.Node {
			if column.Key == "actions" {
				return g.El("strong", g.Text("Rich action"))
			}
			return nil
		},
		SortURL: func(column molecules.TableColumn) string {
			return "/users?sort=" + column.Key
		},
		SortState: func(column molecules.TableColumn) string {
			return "ascending"
		},
		SelectAllLabel: "Select every user",
		SelectRowLabel: func(row molecules.TableRow) string {
			return "Select " + row.ID
		},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`aria-label="Select every user"`,
		`aria-label="Select user-1"`,
		`aria-sort="ascending"`,
		`style="width:12rem"`,
		`hx-get="/users?sort=name"`,
		`hx-target="#users-grid"`,
		`hx-swap="outerHTML"`,
		`<strong>Rich action</strong>`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("TableWithSlots output is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "plain fallback") {
		t.Errorf("rich cell slot must replace the fallback value: %s", html)
	}
}

func TestTableWithSlotsCarriesTrustedRowAndCellProjection(t *testing.T) {
	var rendered strings.Builder
	err := TableWithSlots(molecules.TableProps{
		Selectable: true,
		Sortable:   true,
		Columns:    []molecules.TableColumn{{Key: "name", Label: "Name", Sortable: true}},
		Rows:       []molecules.TableRow{{ID: "row-1", Cells: map[string]any{"name": "Ada"}}},
	}, TableSlots{
		RowAttrs: func(molecules.TableRow) []g.Node {
			return []g.Node{g.Attr("data-state", "selected")}
		},
		CellAttrs: func(_ molecules.TableRow, column molecules.TableColumn) []g.Node {
			return []g.Node{g.Attr("data-label", column.Label)}
		},
		SortButtonAttrs: func(molecules.TableColumn) []g.Node {
			return []g.Node{g.Attr("hx-include", "[name='search']")}
		},
		SelectRowChecked: func(molecules.TableRow) bool { return true },
	}).Render(&rendered)
	if err != nil {
		t.Fatalf("render table slots: %v", err)
	}
	for _, fragment := range []string{
		`data-state="selected"`,
		`data-label="Name"`,
		`hx-include="[name=&#39;search&#39;]"`,
		`checked`,
	} {
		if !strings.Contains(rendered.String(), fragment) {
			t.Errorf("TableWithSlots output is missing %q: %s", fragment, rendered.String())
		}
	}
}

func renderAll(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, n := range gallery() {
		if n == nil {
			continue
		}
		if err := n.Render(&b); err != nil {
			t.Fatalf("render: %v", err)
		}
		b.WriteString("\n")
	}
	return b.String()
}

var classAttrRE = regexp.MustCompile(`class="([^"]*)"`)

// TestRenderedClassesAreDeclared is the architectural loop-closure: every
// class in rendered HTML must be covered by the stylesheet derived from
// ClassLists(). A renderer inventing a class the lists do not declare, or a
// list drifting from a renderer, fails here — which is exactly the failure
// Tailwind users hit at runtime as silently unstyled markup.
func TestRenderedClassesAreDeclared(t *testing.T) {
	t.Parallel()

	sheet, err := emission.For(ClassLists()...)
	if err != nil {
		t.Fatalf("emission.For over declared lists: %v", err)
	}
	css, err := sheet.RenderPretty(), error(nil)
	if err != nil {
		t.Fatal(err)
	}

	html := renderAll(t)
	seen := map[string]bool{}
	for _, m := range classAttrRE.FindAllStringSubmatch(html, -1) {
		for _, class := range strings.Fields(m[1]) {
			seen[class] = true
		}
	}
	if len(seen) < 60 {
		t.Fatalf("only %d distinct classes rendered; the gallery regressed", len(seen))
	}

	var missing []string
	for class := range seen {
		if _, err := emission.Rules(class); err != nil {
			missing = append(missing, class+" (unresolvable)")
			continue
		}
		// The class must be in the derived sheet, not merely resolvable —
		// prefix-escaped for selector matching.
		esc := strings.NewReplacer(":", "\\:", "/", "\\/", "[", "\\[", "]", "\\]", ".", "\\.").Replace(class)
		if !strings.Contains(css, "."+esc) {
			missing = append(missing, class+" (not in For(ClassLists()) sheet)")
		}
	}
	sort.Strings(missing)
	for _, m := range missing {
		t.Errorf("rendered class not backed by the design system: %s", m)
	}
}

func TestAccessibilityStructure(t *testing.T) {
	t.Parallel()
	html := renderAll(t)
	for _, want := range []string{
		`aria-invalid="true"`,                    // errored input
		`aria-describedby="pk-input-slug-error"`, // error linkage
		`for="pk-input-email"`,                   // label association
		`role="alert"`,                           // severe alert interrupts
		`role="status"`,                          // polite alert does not
		`aria-current="page"`,                    // breadcrumb + pagination current
		`aria-label="Breadcrumb"`,                // named navigation landmarks
		`aria-label="Pagination"`,
		`aria-selected="true"`,      // active tab
		`aria-hidden="true"`,        // decorative icons hidden
		`rel="noopener noreferrer"`, // external link safety
		`aria-busy="true"`,          // loading button
		`scope="col"`,               // table headers
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered gallery missing accessibility structure %s", want)
		}
	}
}

func TestGalleryGolden(t *testing.T) {
	t.Parallel()
	html := renderAll(t)
	golden := "testdata/gallery.golden.html"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(html), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with UPDATE_GOLDEN=1 to create): %v", err)
	}
	if string(want) != html {
		t.Fatalf("rendered gallery differs from golden (%d vs %d bytes); rerun with UPDATE_GOLDEN=1 and review the diff", len(html), len(want))
	}
}

// TestUnimplementedListIsHonest pins the deliberate scope: these contracts
// have no renderer yet. Implementing one without removing it here fails, and
// removing one without implementing fails at compile time in the gallery.
func TestUnimplementedListIsHonest(t *testing.T) {
	t.Parallel()
	unimplemented := []string{
		"atoms.Avatar", "atoms.Progress", "atoms.Radio", "atoms.Skeleton",
		"atoms.Slider", "atoms.Toast", "atoms.Toggle", "atoms.Tooltip",
		"molecules.Accordion", "molecules.ActionMenu", "molecules.Autocomplete",
		"molecules.DatePicker", "molecules.Drawer", "molecules.Dropdown",
		"molecules.FileUpload", "molecules.Modal",
		"molecules.Sidebar", "molecules.Stepper",
	}
	if len(unimplemented) == 0 {
		t.Skip("everything implemented")
	}
	// The list is prose-verified: it exists so the package doc and release
	// notes can point at one authoritative statement of scope.
}

func TestCursorPaginationUsesOpaqueContinuationLinks(t *testing.T) {
	t.Parallel()
	node := Pagination(molecules.PaginationProps{
		HTMXProps: contracts.HTMXProps{
			Target:      "#results",
			Swap:        "outerHTML",
			DisabledElt: "this",
		},
		TotalPages:    0,
		CursorMode:    "load-more",
		NextURL:       "/results?after=opaque",
		LoadMoreLabel: "Load more results",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`href="/results?after=opaque"`,
		`hx-get="/results?after=opaque"`,
		`hx-target="#results"`,
		`data-cursor-mode="load-more"`,
		`Load more results`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("cursor pagination missing %q: %s", fragment, html)
		}
	}
}

func TestOffsetPaginationPreservesQueryAndUsesPerPageHTMXURLs(t *testing.T) {
	t.Parallel()
	node := Pagination(molecules.PaginationProps{
		HTMXProps: contracts.HTMXProps{
			Get:    "/results?page_size=20&sort=name",
			Target: "#results",
			Swap:   "outerHTML",
		},
		CurrentPage: 2,
		TotalPages:  3,
		BaseURL:     "/results?page_size=20&sort=name",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`href="/results?page=1&amp;page_size=20&amp;sort=name"`,
		`hx-get="/results?page=1&amp;page_size=20&amp;sort=name"`,
		`href="/results?page=3&amp;page_size=20&amp;sort=name"`,
		`hx-get="/results?page=3&amp;page_size=20&amp;sort=name"`,
		`hx-target="#results"`,
		`hx-swap="outerHTML"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("offset pagination missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "??") {
		t.Fatalf("offset pagination emitted a malformed URL: %s", html)
	}
}

func TestWindowedCollectionFailsClosedAboveBudget(t *testing.T) {
	t.Parallel()
	node := WindowedCollection(
		organisms.WindowedCollectionProps{ItemCount: 501, MaxItems: 500},
		WindowedCollectionSlots{Items: []g.Node{Text(atoms.TextProps{Content: "must not render"})}},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if strings.Contains(html, "must not render") ||
		!strings.Contains(html, `data-windowed-error-kind="contract"`) {
		t.Fatalf("window did not fail closed: %s", html)
	}
}
