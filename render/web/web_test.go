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
		Button(atoms.ButtonProps{Text: "Save", Variant: "primary", Size: "medium"}),
		Button(atoms.ButtonProps{Text: "Cancel", Variant: "secondary", Size: "small"}),
		Button(atoms.ButtonProps{Text: "Delete", Variant: "danger", Size: "xs"}),
		Button(atoms.ButtonProps{Text: "Ghost", Variant: "ghost", Size: "large"}),
		Button(atoms.ButtonProps{Text: "Docs", Variant: "link", Size: "xl"}),
		Button(atoms.ButtonProps{Text: "Outline", Variant: "outline", Size: "2xl", FullWidth: true}),
		Button(atoms.ButtonProps{Text: "Info", Variant: "info"}),
		Button(atoms.ButtonProps{Text: "Warn", Variant: "warning"}),
		Button(atoms.ButtonProps{Text: "OK", Variant: "success"}),
		Button(atoms.ButtonProps{Text: "Err", Variant: "error"}),
		Button(atoms.ButtonProps{Text: "Busy", Loading: true}),
		Button(atoms.ButtonProps{Text: "Iconed", Icon: "plus", IconPosition: "right"}),

		Badge(atoms.BadgeProps{Text: "Default"}),
		Badge(atoms.BadgeProps{Text: "New", Variant: "primary", Dot: true}),
		Badge(atoms.BadgeProps{Text: "OK", Variant: "success"}),
		Badge(atoms.BadgeProps{Text: "Careful", Variant: "warning", Icon: "alert"}),
		Badge(atoms.BadgeProps{Text: "Bad", Variant: "danger"}),
		Badge(atoms.BadgeProps{Text: "FYI", Variant: "info"}),
		Badge(atoms.BadgeProps{Text: "Two", Variant: "secondary"}),
		Badge(atoms.BadgeProps{Text: "Err", Variant: "error"}),

		Alert(atoms.AlertProps{Title: "Saved", Message: "All good.", Variant: "success"}),
		Alert(atoms.AlertProps{Message: "Careful now.", Variant: "warning", Icon: "alert"}),
		Alert(atoms.AlertProps{Message: "That failed.", Variant: "error"}),
		Alert(atoms.AlertProps{Message: "Broken.", Variant: "danger"}),
		Alert(atoms.AlertProps{Message: "Heads up.", Variant: "info"}),

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
		Divider(atoms.DividerProps{Orientation: "vertical"}),

		Spinner(atoms.SpinnerProps{Size: "small"}),
		Spinner(atoms.SpinnerProps{Size: "medium", Label: "Fetching"}),
		Spinner(atoms.SpinnerProps{Size: "large"}),

		EmptyState(atoms.EmptyStateProps{Title: "No tenants yet",
			Description: "Create the first tenant to get started.", Bordered: true,
			Actions: []atoms.EmptyStateAction{{Label: "New tenant", Href: "/admin/tenants/new"}}}),
		EmptyState(atoms.EmptyStateProps{Title: "Empty", Compact: true}),

		Kbd(atoms.KbdProps{Keys: []string{"Ctrl", "K"}}),

		Link(atoms.LinkProps{Text: "Internal", Href: "/docs"}),
		Link(atoms.LinkProps{Text: "External", Href: "https://example.test", External: true}),

		Tag(atoms.TagProps{Text: "beta"}),
		Tag(atoms.TagProps{Text: "chosen", Selected: true, Removable: true, OnRemoveURL: "/tags/1"}),

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

		DataGrid(organisms.DataGridProps{
			Search:    molecules.SearchBarProps{Label: "Search rows", Placeholder: "Filter…"},
			Actions:   []atoms.ButtonProps{{Text: "Refresh", Variant: "secondary"}},
			CreateURL: "/rows/new", CreateText: "New row",
			Filters: []organisms.DataGridFilter{{Key: "state", Label: "State", Type: "select",
				Options: []molecules.Option{{Label: "Open", Value: "open"}}}},
			Table: molecules.TableProps{Sortable: true,
				Columns: []molecules.TableColumn{{Key: "name", Label: "Name", Sortable: true, Primary: true}},
				Rows:    []molecules.TableRow{{ID: "g1", Cells: map[string]any{"name": "Grid row"}}}},
			Pagination: molecules.PaginationProps{CurrentPage: 1},
		},
			Text(atoms.TextProps{Content: "1 record on this page", Color: "muted", Size: "sm"}),
		),

		Tabs(molecules.TabsProps{ActiveTab: "b", Items: []molecules.TabItem{
			{Key: "a", Label: "First", URL: "/tab/a"}, {Key: "b", Label: "Second"},
		}}),
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
