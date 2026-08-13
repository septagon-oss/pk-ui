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
	"bytes"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

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
	no := false
	return []g.Node{
		Avatar(atoms.AvatarProps{Name: "Ada Lovelace", Size: "xs"}),
		Avatar(atoms.AvatarProps{Initials: "GH", Alt: "Grace Hopper", Size: "sm", Shape: "rounded", Tone: "brand", Status: "online"}),
		Avatar(atoms.AvatarProps{Src: "/avatar.png", Alt: "Lin", Size: "lg", Shape: "square", Status: "busy", StatusPosition: "top-left"}),
		Avatar(atoms.AvatarProps{FallbackIcon: "user", Alt: "Account", Size: "2xl", Shape: "pill", Tone: "info", Status: "away", StatusPosition: "bottom-left"}),

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
		ButtonWithSlots(
			atoms.ButtonProps{Label: "Iconed", Variant: "primary", Tone: "neutral", Size: "md"},
			ButtonSlots{IconEnd: []g.Node{Icon(atoms.IconProps{Name: "plus", Size: "md", Tone: "neutral"})}},
		),

		Badge(atoms.BadgeProps{Label: "Default"}),
		Badge(atoms.BadgeProps{Label: "New", Variant: "primary", Tone: "brand", Dot: true}),
		Badge(atoms.BadgeProps{Label: "OK", Tone: "success"}),
		BadgeWithSlots(
			atoms.BadgeProps{Label: "Careful", Tone: "warning"},
			BadgeSlots{IconStart: []g.Node{Icon(atoms.IconProps{Name: "alert", Size: "sm", Tone: "warning"})}},
		),
		Badge(atoms.BadgeProps{Label: "Bad", Tone: "danger"}),
		Badge(atoms.BadgeProps{Label: "FYI", Tone: "info"}),
		Badge(atoms.BadgeProps{Label: "Two", Variant: "secondary"}),
		Badge(atoms.BadgeProps{Label: "Outlined", Variant: "outline"}),
		Badge(atoms.BadgeProps{Label: "Messages", Count: 120, Removable: true, Live: true}),

		Alert(atoms.AlertProps{Title: "Saved", Message: "All good.", Tone: "success"}),
		alertWithSlots(
			atoms.AlertProps{Message: "Careful now.", Tone: "warning"},
			[]g.Node{Icon(atoms.IconProps{Name: "alert", Size: "md", Tone: "warning"})},
			nil,
		),
		Alert(atoms.AlertProps{Message: "That failed.", Tone: "danger"}),
		Alert(atoms.AlertProps{Message: "Broken.", Tone: "danger"}),
		Alert(atoms.AlertProps{Message: "Heads up.", Tone: "info"}),
		Alert(atoms.AlertProps{Message: "Compact.", Tone: "info", Compact: true}),
		Alert(atoms.AlertProps{Message: "Accented.", Tone: "warning", Bordered: true}),
		Alert(atoms.AlertProps{Message: "Dismiss me.", Tone: "success", Dismissible: true}),

		Toast(atoms.ToastProps{Title: "Saved", Message: "All changes are live.", Tone: "success", Duration: 5000}),
		Toast(atoms.ToastProps{Message: "Connection interrupted.", Tone: "warning", Persistent: true}),
		Toast(atoms.ToastProps{Message: "Background sync failed.", Tone: "danger", Position: "bottom-left", Closable: &no}),
		toastWithSlots(
			atoms.ToastProps{Message: "New activity", Tone: "info"},
			[]g.Node{Icon(atoms.IconProps{Name: "bell", Size: "sm", Tone: "info"})},
		),

		Input(atoms.InputProps{Name: "email", Type: "email", Label: "Email",
			Placeholder: "you@example.test", HelpText: "We never share it.", Required: true}),
		Input(atoms.InputProps{Name: "slug", Label: "Slug", Value: "hello",
			Error: "Already taken.", Pattern: "[a-z-]+"}),
		Input(atoms.InputProps{Name: "quiet"}),
		DatePicker(molecules.DatePickerProps{
			Name: "start_date", Label: "Start date", Required: true,
			Min: "2026-01-01", Max: "2026-12-31", HelpText: "Choose a date in 2026.",
		}),
		DatePicker(molecules.DatePickerProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "due_date", Label: "Due date", Error: "Choose a due date.", Format: "YYYY-MM-DD",
			HTMXProps: contracts.HTMXProps{Post: "/dates/preview", Target: "#preview", Swap: "outerHTML"},
		}),
		FileUpload(molecules.FileUploadProps{
			ComponentProps: contracts.ComponentProps{ID: "evidence-upload"},
			HTMXProps:      contracts.HTMXProps{Post: "/evidence", Target: "#evidence-list"},
			Name:           "evidence", Label: "Evidence", Accept: ".pdf,.png", Multiple: true,
			Required: true, MaxSize: 5 * 1024 * 1024, Preview: true,
			HelpText: "Attach supporting evidence.",
		}),
		FileUpload(molecules.FileUploadProps{
			Name: "logo_url", Label: "Company logo", Accept: "image/*", Preview: true,
			UploadURL: "/api/v1/files/upload", UploadCategory: "image",
			Value: "/media/logo.png", CurrentName: "logo.png",
		}),
		FileUpload(molecules.FileUploadProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "attachment", Label: "Attachment", DropZone: &no, ShowList: &no,
		}),
		Autocomplete(molecules.AutocompleteProps{
			Name: "assignee_id", Label: "Assign to", Placeholder: "Search users...", Required: true,
			Options: []molecules.Option{
				{Value: "ada", Label: "Ada Lovelace"},
				{Value: "grace", Label: "Grace Hopper"},
			},
		}),
		Autocomplete(molecules.AutocompleteProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "company_id", Label: "Company", SearchURL: "/companies/search",
			Value: "c1", DisplayValue: "Acme", Error: "Choose a valid company.",
		}),
		Dropdown(molecules.DropdownProps{
			ComponentProps: contracts.ComponentProps{ID: "country-dropdown"},
			Name:           "country", Label: "Country", Placeholder: "Select a country", Value: "pt", Clearable: true,
			Options: []molecules.Option{{Value: "us", Label: "United States"}, {Value: "pt", Label: "Portugal"}},
		}),
		Dropdown(molecules.DropdownProps{
			Name: "topics", AriaLabel: "Topics", Searchable: true, Multiple: true,
			Selected: []string{"alerts", "release"},
			Options: []molecules.Option{
				{Value: "alerts", Label: "Alerts", Icon: "bell", Group: "Operations"},
				{Value: "release", Label: "Releases", Group: "Operations"},
				{Value: "billing", Label: "Billing", Disabled: true, Group: "Finance"},
			},
		}),
		ActionMenu(molecules.ActionMenuProps{
			ComponentProps: contracts.ComponentProps{ID: "row-actions"},
			Items: []molecules.ActionMenuItem{
				{Label: "View", Icon: "eye", Href: "/items/1"},
				{Label: "Retry", Icon: "arrow-path", HxPost: "/items/1/retry", HxTarget: "#items", HxSwap: "innerHTML"},
			},
			Sections: []molecules.ActionMenuSection{{Label: "Danger zone", Items: []molecules.ActionMenuItem{
				{Label: "Delete", Icon: "trash", Danger: true, HxDelete: "/items/1", HxConfirm: "Delete it?"},
			}}},
		}),
		ActionMenu(molecules.ActionMenuProps{TriggerLabel: "Manage", Align: "start", Width: "sm", Items: []molecules.ActionMenuItem{{Label: "Edit"}}}),
		ActionMenu(molecules.ActionMenuProps{ComponentProps: contracts.ComponentProps{Disabled: true}, Width: "lg", Items: []molecules.ActionMenuItem{{Label: "Locked", Disabled: true}}}),
		Drawer(molecules.DrawerProps{ComponentProps: contracts.ComponentProps{ID: "right-drawer"}, Title: "Details", Description: "Review this record.", Body: "Drawer body", Footer: "Drawer actions", Position: "right", Size: "small", Open: true}),
		Drawer(molecules.DrawerProps{Title: "Left navigation", Position: "left", Size: "medium"}),
		Drawer(molecules.DrawerProps{Title: "Bottom sheet", Position: "bottom", Size: "large"}),
		Drawer(molecules.DrawerProps{Title: "Tall bottom sheet", Position: "bottom", Size: "xl"}),
		Drawer(molecules.DrawerProps{Position: "bottom", Size: "full", Closable: &no, CloseOnOverlay: &no, CloseOnEscape: &no, ShowOverlay: &no, OpenOnSwap: true}),
		Modal(molecules.ModalProps{ComponentProps: contracts.ComponentProps{ID: "confirm-modal"}, Title: "Archive", Description: "This action cannot be undone.", Body: "Review the affected records.", Footer: "Confirm or cancel.", Size: "small", Open: true}),
		Modal(molecules.ModalProps{Title: "Edit record", Size: "medium"}),
		Modal(molecules.ModalProps{Title: "Large review", Size: "large"}),
		Modal(molecules.ModalProps{Title: "Wide review", Size: "xl"}),
		Modal(molecules.ModalProps{AriaLabel: "Required decision", Size: "full", Closable: &no, CloseOnOverlay: &no, CloseOnEscape: &no, ShowClose: &no, ShowOverlay: &no, Centered: &no}),
		Modal(molecules.ModalProps{ComponentProps: contracts.ComponentProps{ID: "server-modal"}, AriaLabel: "Server dialog", Deferred: true, OpenOnSwap: true}),

		Textarea(atoms.TextareaProps{Name: "body", Label: "Body", Rows: 6,
			HelperText: "Markdown is fine.", MaxLength: 500}),
		Textarea(atoms.TextareaProps{Name: "bad", Label: "Bad", ErrorMessage: "Too long."}),
		Textarea(atoms.TextareaProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "details", Label: "Details", Value: "Existing", HelperText: "Add context.",
			ErrorMessage: "More detail is required.", AutoResize: true, MinRows: 3, MaxRows: 15, FullWidth: true,
		}),

		Checkbox(atoms.CheckboxProps{Name: "agree", Label: "I agree", Required: true}),
		Checkbox(atoms.CheckboxProps{Name: "done", Label: "Done", Checked: true}),
		Checkbox(atoms.CheckboxProps{Name: "some", Label: "Some", Indeterminate: true}),
		Checkbox(atoms.CheckboxProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "disabled",
			Label:          "Disabled",
		}),
		Radio(atoms.RadioProps{Name: "plan", Value: "standard", Label: "Standard"}),
		Radio(atoms.RadioProps{Name: "plan", Value: "pro", Label: "Pro", Checked: true, Required: true}),
		Radio(atoms.RadioProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "plan", Value: "enterprise", Label: "Enterprise",
		}),
		Radio(atoms.RadioProps{Name: "delivery", Value: "express"}),
		Slider(atoms.SliderProps{Name: "opacity", Label: "Opacity", Min: 0, Max: 100, Step: 1, Value: 75, ShowValue: true}),
		Slider(atoms.SliderProps{Name: "cpu", Label: "CPU", Min: 10, Max: 500, Step: 5, Value: 100, Tone: "warning", AriaValueText: "100 millicores"}),
		Slider(atoms.SliderProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "volume", Min: 0, Max: 100, Value: 40,
		}),
		Toggle(atoms.ToggleProps{Name: "notifications", Label: "Notifications", Size: "md"}),
		Toggle(atoms.ToggleProps{Name: "dark_mode", Label: "Dark mode", Checked: true, Size: "lg"}),
		Toggle(atoms.ToggleProps{Name: "marketing", Label: "Marketing", HideLabel: true, Checked: true, Size: "sm"}),
		Toggle(atoms.ToggleProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Name:           "feature", Label: "Feature", Size: "md",
		}),

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
		Progress(atoms.ProgressProps{Value: 75, Max: 100, AriaLabel: "Upload progress", ShowText: true}),
		Progress(atoms.ProgressProps{Value: 45, Max: 100, Label: "Importing records", ShowText: true, Tone: "warning", Size: "lg"}),
		Progress(atoms.ProgressProps{Value: 100, Max: 100, AriaLabel: "Import completion", ShowText: true, Tone: "success", Size: "sm"}),
		Progress(atoms.ProgressProps{Label: "Preparing export", Indeterminate: true, Tone: "info"}),
		Tooltip(atoms.TooltipProps{Content: "Top help", Position: "top"}, h.Button(g.Text("Top"))),
		Tooltip(atoms.TooltipProps{Content: "Bottom help", Position: "bottom"}, h.Button(g.Text("Bottom"))),
		Tooltip(atoms.TooltipProps{Content: "Left help", Position: "left"}, h.Button(g.Text("Left"))),
		Tooltip(atoms.TooltipProps{Content: "Right help", Position: "right"}, h.Button(g.Text("Right"))),

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
		Accordion(molecules.AccordionProps{
			DefaultOpen: []string{"intro"},
			Items: []molecules.AccordionItem{
				{ID: "intro", Title: "Introduction", Subtitle: "Start here", Icon: "information-circle", Content: "Welcome."},
				{ID: "limits", Title: "Limits", Content: "Usage limits.", Disabled: true},
			},
		}),
		AccordionWithSections(
			molecules.AccordionProps{Multiple: true, Bordered: &no},
			AccordionSection{ID: "rich", Title: "Rich content", Open: true, Content: []g.Node{h.Strong(g.Text("Composed panel"))}},
		),
		Stepper(molecules.StepperProps{
			CurrentStep: 1,
			Steps: []molecules.StepItem{
				{Key: "account", Label: "Account", Icon: "user"},
				{Key: "profile", Label: "Profile", Description: "Tell us about yourself"},
				{Key: "review", Label: "Review"},
			},
		}),
		Stepper(molecules.StepperProps{
			Orientation: "vertical", CurrentStep: 1,
			Steps: []molecules.StepItem{
				{Label: "Create account", Description: "Set up your credentials"},
				{Label: "Profile", Description: "Tell us about yourself"},
				{Label: "Preferences", Description: "Customize your experience", Status: "error"},
			},
		}),
		Stepper(molecules.StepperProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			CurrentStep:    2, Clickable: true, Compact: true, StepAction: "click->checkout#step",
			Steps: []molecules.StepItem{
				{Key: "cart", Label: "Cart"},
				{Key: "shipping", Label: "Shipping"},
				{Key: "payment", Label: "Payment"},
			},
		}),
		SidebarWithSlots(molecules.SidebarProps{
			Current: "/admin/customers/accounts",
			Sections: []molecules.SidebarSection{{
				ID: "operate", Label: "Operate", Glyph: "O", Tone: "brand",
				Items: []molecules.SidebarItem{
					{Label: "Dashboard", Href: "/admin", Icon: "home"},
					{Label: "Customers", Href: "/admin/customers", Icon: "users", Badge: "24", Children: []molecules.SidebarItem{
						{Label: "Accounts", Href: "/admin/customers/accounts"},
					}},
				},
			}},
		}, SidebarSlots{Footer: []g.Node{g.Text("Signed in")}}),
		SidebarWithSlots(molecules.SidebarProps{
			ComponentProps: contracts.ComponentProps{ID: "docs-sidebar"},
			Flavor:         "content", Current: "#figma", NavigationLabel: "Documentation sections",
			Sections: []molecules.SidebarSection{{
				ID: "runtime", Label: "Runtime", Glyph: "R", Tone: "info", SearchText: "runtime docs",
				Items: []molecules.SidebarItem{
					{ID: "overview", Label: "Overview", Href: "#overview", Prefix: "01", SearchText: "overview"},
					{ID: "figma", Label: "Figma handoff", Href: "#figma", Prefix: "02", SearchText: "figma handoff"},
				},
			}},
		}, SidebarSlots{
			Brand:  []g.Node{h.Strong(g.Text("Documentation"))},
			Footer: []g.Node{g.Text("Version 1")},
		}),
		Sidebar(molecules.SidebarProps{
			ComponentProps: contracts.ComponentProps{Disabled: true},
			Collapsible:    true, Collapsed: true,
			Items: []molecules.SidebarItem{
				{Label: "Home", Href: "/admin", Icon: "home"},
				{Label: "Reports", Href: "/admin/reports", Icon: "chart", Disabled: true},
			},
		}),

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
		DashboardWidget(
			organisms.DashboardWidgetProps{
				Title: "Active seats", Subtitle: "Current workspace posture",
				Value: "1,234", PreviousValue: "1,102", Change: "+12%",
				Trend: "up", Icon: "activity", Span: "2",
			},
			DashboardWidgetSlots{Footer: []g.Node{g.Text("Updated now")}},
		),

		Tabs(molecules.TabsProps{ActiveTab: "b", Items: []molecules.TabItem{
			{Key: "a", Label: "First", URL: "/tab/a"}, {Key: "b", Label: "Second"},
		}}),
		TabsWithPanels(
			molecules.TabsProps{ActiveTab: "profile", Orientation: "vertical", Variant: "pills"},
			TabSlot{ID: "profile", Label: "Profile", Icon: "user", Badge: "New", Content: []g.Node{g.Text("Profile panel")}},
			TabSlot{ID: "security", Label: "Security", Disabled: true, Content: []g.Node{g.Text("Security panel")}},
			TabSlot{ID: "activity", Label: "Activity", HxGet: "/activity"},
		),
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

func TestButtonOwnsNativeLinkSlotsAndStateContract(t *testing.T) {
	t.Parallel()

	var native strings.Builder
	if err := ButtonWithSlots(atoms.ButtonProps{
		ComponentProps: contracts.ComponentProps{
			ID:       "save",
			Disabled: true,
			Attrs:    map[string]string{"data-action": "save"},
		},
		HTMXProps: contracts.HTMXProps{Post: "/save", Target: "#result", Include: "#save-form"},
		Label:     "Saving",
		Variant:   "primary",
		Tone:      "success",
		Size:      "lg",
		Type:      "submit",
		Loading:   true,
	}, ButtonSlots{}).Render(&native); err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{
		`<button`, `id="save"`, `disabled`, `type="submit"`,
		`data-component="button"`, `data-variant="primary"`, `data-tone="success"`,
		`data-loading="true"`, `aria-busy="true"`, `hx-post="/save"`,
		`hx-target="#result"`, `hx-include="#save-form"`, `data-action="save"`, `Saving`,
	} {
		if !strings.Contains(native.String(), fragment) {
			t.Errorf("native Button output is missing %q: %s", fragment, native.String())
		}
	}

	var link strings.Builder
	if err := ButtonWithSlots(atoms.ButtonProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Label:          "Continue",
		Href:           "/continue",
		AriaLabel:      "Continue to checkout",
	}, ButtonSlots{Content: []g.Node{h.Span(g.Text("Compound action"))}}).Render(&link); err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{
		`<a`, `href="/continue"`, `data-button-as-link="true"`,
		`aria-disabled="true"`, `tabindex="-1"`,
		`aria-label="Continue to checkout"`, `Compound action`,
	} {
		if !strings.Contains(link.String(), fragment) {
			t.Errorf("link Button output is missing %q: %s", fragment, link.String())
		}
	}
	for _, forbidden := range []string{` type=`, ` disabled="`, `>Continue<`} {
		if strings.Contains(link.String(), forbidden) {
			t.Errorf("link Button output unexpectedly contains %q: %s", forbidden, link.String())
		}
	}
}

func TestTextOwnsSemanticElementTypographyAndClampContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := Text(atoms.TextProps{
		ComponentProps: contracts.ComponentProps{
			ID: "summary", Class: "product-copy", Attrs: map[string]string{"data-owner": "profile"},
		},
		Content: "Important account summary", Element: "STRONG", Size: "5xl",
		Align: "justify", Weight: "extra bold", Color: "tertiary",
		Transform: "uppercase", Italic: true, Underline: true, NoWrap: true,
		Truncate: true, Lines: 20,
	}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<strong`, `id="summary"`, `product-copy`, `data-owner="profile"`,
		`data-component="text"`, `data-element="strong"`, `data-size="5xl"`,
		`data-align="justify"`, `data-weight="extrabold"`, `data-color="tertiary"`,
		`text-5xl`, `text-justify`, `font-extrabold`, `text-fg-tertiary`,
		`uppercase`, `italic`, `underline`, `whitespace-nowrap`, `line-clamp-6`,
		`data-lines="6"`, `Important account summary`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Text missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, ` truncate`) {
		t.Fatalf("line-clamped Text must not also use the single-line truncate contract: %s", html)
	}
}

func TestTextRejectsHeadingTagsAndNormalizesDefaults(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := Text(atoms.TextProps{Content: "Body copy", Element: "h1", Color: "danger"}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<p`, `data-element="p"`, `data-size="base"`, `data-align="left"`,
		`data-weight="normal"`, `data-color="danger"`, `text-fg-danger`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("normalized Text missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `<h1`) {
		t.Fatalf("Text must not bypass the canonical Heading component: %s", html)
	}
}

func TestLinkWithTrailingAdornmentUsesCanonicalAnchor(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := LinkWithTrailingAdornment(
		atoms.LinkProps{
			ComponentProps: contracts.ComponentProps{
				Class: "error-action",
				Attrs: map[string]string{"data-action": "retry"},
			},
			Label: "Try again",
			Href:  "/retry",
		},
		Icon(atoms.IconProps{Name: "arrow-right", Size: "sm"}),
	).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<a`,
		`href="/retry"`,
		`data-action="retry"`,
		`Try again`,
		`data-pk-icon="arrow-right"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("slotted Link output is missing %q: %s", fragment, html)
		}
	}
}

func TestBreadcrumbOwnsTruncationAndIconRendering(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Breadcrumb(molecules.BreadcrumbProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "account-trail",
			Class: "client-breadcrumb",
			Attrs: map[string]string{"data-client": "collect"},
		},
		HTMXProps: contracts.HTMXProps{Boost: true},
		Separator: "›",
		MaxItems:  3,
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
		`id="account-trail"`,
		`class="client-breadcrumb"`,
		`data-client="collect"`,
		`data-component="breadcrumb"`,
		`hx-boost="true"`,
		`data-pk-icon="home"`,
		`href="/"`,
		`>…</li>`,
		`href="/settings"`,
		`aria-current="page">Profile`,
		`>›</li>`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Breadcrumb output is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "Workspace") {
		t.Errorf("Breadcrumb did not collapse its middle items: %s", html)
	}
}

func TestAccordionOwnsDisclosureStateRichContentAndNormalizedIDs(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := AccordionWithSections(
		molecules.AccordionProps{
			ComponentProps: contracts.ComponentProps{ID: "faq", Class: "client-faq"},
			Multiple:       true,
			DefaultOpen:    []string{"shipping"},
		},
		AccordionSection{
			ID:       "shipping",
			Title:    "Shipping",
			Subtitle: "Delivery details",
			Icon:     "information-circle",
			Content:  []g.Node{h.Strong(g.Text("Ships worldwide"))},
		},
		AccordionSection{
			ID:      "shipping",
			Title:   "Duplicate ID",
			Open:    true,
			Content: []g.Node{g.Text("Still independently addressable")},
		},
		AccordionSection{
			Title:    "Unavailable",
			Icon:     "bad<script",
			Disabled: true,
			Content:  []g.Node{g.Text("Disabled panel")},
		},
	).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="faq"`, `class="`, `client-faq`,
		`data-component="accordion"`, `data-controller="accordion"`,
		`data-accordion-multiple-value="true"`,
		`data-accordion-item="shipping"`, `data-accordion-item="shipping-2"`,
		`data-accordion-item="item-3"`, `id="faq-header-shipping"`,
		`id="faq-panel-shipping"`, `role="region"`,
		`aria-labelledby="faq-header-shipping"`, `data-pk-icon="information-circle"`,
		`<strong>Ships worldwide</strong>`, `disabled`, `aria-disabled="true"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Accordion is missing %q: %s", fragment, html)
		}
	}
	if got := strings.Count(html, `aria-expanded="true"`); got != 2 {
		t.Errorf("open sections = %d, want 2: %s", got, html)
	}
	if strings.Contains(html, "bad<script") {
		t.Fatalf("invalid icon names must not become markup: %s", html)
	}
}

func TestAccordionPortableItemsEscapeContentAndSupportKeyAlias(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Accordion(molecules.AccordionProps{Items: []molecules.AccordionItem{{
		Key:     "portable",
		Title:   "Portable",
		Content: `<script>alert("no")</script>`,
		Open:    true,
	}}}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-accordion-item="portable"`,
		`aria-expanded="true"`,
		`&lt;script&gt;alert(&#34;no&#34;)&lt;/script&gt;`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("portable Accordion is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "<script>") {
		t.Fatalf("portable content must be escaped: %s", html)
	}
	if node := Accordion(molecules.AccordionProps{}); node != nil {
		t.Fatalf("empty Accordion = %#v, want nil", node)
	}
}

func TestStepperDerivesStatusesAndPreservesExplicitError(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Stepper(molecules.StepperProps{
		ComponentProps: contracts.ComponentProps{ID: "setup", Class: "client-stepper"},
		CurrentStep:    2,
		Steps: []molecules.StepItem{
			{Key: "account", Label: "Account", Icon: "user"},
			{Key: "profile", Label: "Profile", Icon: "bad<script"},
			{Key: "billing", Label: "Billing", Description: "Choose a plan"},
			{Key: "review", Label: "Review", Status: "error"},
		},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<nav id="setup"`, `client-stepper`, `aria-label="Progress"`,
		`data-step-key="account"`, `data-step-status="completed"`,
		`data-step-key="billing"`, `data-step-status="active"`,
		`aria-current="step"`, `data-step-key="review"`,
		`data-step-status="error"`, `data-pk-icon="x-mark"`,
		`data-stepper-connector=""`, `data-connector-status="completed"`,
		`data-stepper-label-description=""`, `Choose a plan`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Stepper is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "bad<script") {
		t.Fatalf("invalid icon names must not become markup: %s", html)
	}
	if got := strings.Count(html, `data-step-status="error"`); got != 1 {
		t.Errorf("explicit error status count = %d, want 1: %s", got, html)
	}
}

func TestStepperVerticalWholeStepActionsAndDisabledState(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Stepper(molecules.StepperProps{
		ComponentProps:  contracts.ComponentProps{Disabled: true},
		Orientation:     "vertical",
		StepAction:      "click->wizard#go",
		NavigationLabel: "Setup progress",
		Steps: []molecules.StepItem{
			{Key: "one", Label: "One", Description: "First"},
			{Key: "two", Label: "Two", Description: "Second"},
		},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-orientation="vertical"`, `aria-label="Setup progress"`,
		`aria-disabled="true"`, `<button type="button"`,
		`data-action="click-&gt;wizard#go"`, `data-step-key="one"`,
		`data-stepper-indicator=""`, `data-stepper-label=""`,
		`data-stepper-label-description=""`, `disabled`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("vertical Stepper is missing %q: %s", fragment, html)
		}
	}
	if got := strings.Count(html, `<button`); got != 2 {
		t.Errorf("whole-step button count = %d, want 2: %s", got, html)
	}
}

func TestStepperClickableIndicatorsAndEmptyInput(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Stepper(molecules.StepperProps{
		Clickable: true,
		Steps:     []molecules.StepItem{{Label: "One"}, {Label: "Two"}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(rendered.String(), `<button`); got != 2 {
		t.Errorf("indicator button count = %d, want 2: %s", got, rendered.String())
	}
	if node := Stepper(molecules.StepperProps{}); node != nil {
		t.Fatalf("empty Stepper = %#v, want nil", node)
	}
}

func TestSidebarOwnsNestedActiveNavigationAndPortableBrand(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Sidebar(molecules.SidebarProps{
		ComponentProps: contracts.ComponentProps{ID: "workspace", Class: "client-sidebar"},
		Current:        "/admin/customers/accounts",
		BrandLabel:     "Acme Control",
		BrandHref:      "/control",
		Sections: []molecules.SidebarSection{{
			ID: "operate", Label: "Operate", Glyph: "O", Tone: "brand",
			Items: []molecules.SidebarItem{
				{ID: "dashboard", Label: "Dashboard", Href: "/admin", Icon: "home"},
				{ID: "customers", Label: "Customers", Href: "/admin/customers", Icon: "users", Badge: "24", Children: []molecules.SidebarItem{
					{ID: "accounts", Label: "Accounts", Href: "/admin/customers/accounts"},
				}},
			},
		}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<aside id="workspace"`, `client-sidebar`, `data-component="sidebar"`,
		`data-sidebar-flavor="admin"`, `aria-label="Admin navigation"`,
		`data-sidebar-brand=""`, `href="/control"`, `Acme Control`,
		`data-sidebar-section="operate"`, `data-sidebar-tone="brand"`,
		`data-sidebar-item="customers"`, `data-sidebar-item="accounts"`,
		`data-active="true"`, `aria-current="page"`, `data-has-badge="true"`,
		`data-component="badge"`, `data-sidebar-submenu="customers"`,
		`data-pk-icon="users"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Sidebar is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `data-sidebar-item="dashboard" data-sidebar-depth="0" data-active="true"`) {
		t.Fatalf("the /admin root must not prefix-match every admin route: %s", html)
	}
}

func TestDatePickerOwnsAccessibleConstraintsAndHTMXDefaults(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := DatePicker(molecules.DatePickerProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "booking-date",
			Class: "client-date",
			Attrs: map[string]string{"data-owner": "booking"},
		},
		HTMXProps: contracts.HTMXProps{
			Get: "/availability", Target: "#availability", Swap: "outerHTML",
		},
		Name: "booking_date", Label: "Booking date", Required: true,
		Value: "2026-04-06", Min: "2026-01-01", Max: "2026-12-31",
		HelpText: "Choose a date in 2026.", Format: "YYYY-MM-DD",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-component="datepicker"`, `type="date"`, `id="booking-date"`,
		`name="booking_date"`, `value="2026-04-06"`, `min="2026-01-01"`,
		`max="2026-12-31"`, `required`, `client-date`,
		`data-owner="booking"`, `data-date-format="YYYY-MM-DD"`,
		`hx-get="/availability"`, `hx-target="#availability"`,
		`hx-swap="outerHTML"`, `hx-trigger="change"`,
		`aria-describedby="booking-date-help"`, `data-pk-icon="calendar"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical DatePicker is missing %q: %s", fragment, html)
		}
	}
}

func TestDatePickerErrorAndDisabledStatesRemainAccessible(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := DatePicker(molecules.DatePickerProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Name:           "due_date", Label: "Due date", Error: "Choose a due date.",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="pk-datepicker-due_date"`, `disabled`, `aria-invalid="true"`,
		`aria-describedby="pk-datepicker-due_date-error"`, `role="alert"`,
		`Choose a due date.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical DatePicker error state is missing %q: %s", fragment, html)
		}
	}
}

func TestFileUploadOwnsNativeMultipartControllerAndHTMXContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := FileUpload(molecules.FileUploadProps{
		ComponentProps: contracts.ComponentProps{
			ID: "evidence-upload", Class: "client-upload",
			Attrs: map[string]string{"data-owner": "compliance"},
		},
		HTMXProps: contracts.HTMXProps{
			Post: "/evidence", Target: "#evidence-list", Swap: "outerHTML",
		},
		Name: "evidence", Label: "Evidence", Accept: ".pdf,.png", Multiple: true,
		Required: true, MaxSize: 5 * 1024 * 1024, Preview: true,
		HelpText: "Attach supporting evidence.",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="evidence-upload"`, `client-upload`, `data-owner="compliance"`,
		`data-component="fileupload"`, `data-controller="file-upload"`, `data-state="idle"`,
		`data-file-upload-multiple-value="true"`, `data-file-upload-disabled-value="false"`,
		`data-file-upload-preview-value="true"`, `data-file-upload-max-size-value="5242880"`,
		`data-file-upload-accept-value=".pdf,.png"`, `data-file-upload-remove-label-value="Remove file"`,
		`id="evidence-upload-input"`, `type="file"`, `name="evidence"`, `multiple`, `required`,
		`accept=".pdf,.png"`, `aria-label="Choose Evidence"`,
		`aria-describedby="evidence-upload-input-help"`, `data-file-upload-input="true"`,
		`data-file-upload-dropzone="true"`, `for="evidence-upload-input"`, `aria-label="Upload Evidence"`,
		`change-&gt;file-upload#handleChange`, `drop-&gt;file-upload#drop`,
		`data-file-upload-loading="true"`, `data-file-upload-prompt="true"`,
		`Click to upload`, `or drag and drop`, `Max size: 5.0 MB`,
		`data-file-upload-list="true"`, `data-file-upload-error="true"`,
		`id="evidence-upload-input-help"`, `Attach supporting evidence.`,
		`hx-post="/evidence"`, `hx-target="#evidence-list"`, `hx-swap="outerHTML"`,
		`hx-trigger="change from:[data-file-upload-input]"`,
		`hx-include="find [data-file-upload-input]"`, `hx-encoding="multipart/form-data"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical FileUpload is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `role="button"`) || strings.Contains(html, `openFileDialog`) {
		t.Fatalf("FileUpload must use native label/input activation instead of a simulated button: %s", html)
	}
}

func TestFileUploadRemoteAndSimpleModesPreserveSubmissionSemantics(t *testing.T) {
	t.Parallel()

	var remote strings.Builder
	err := FileUpload(molecules.FileUploadProps{
		Name: "logo_url", Label: "Company logo", Accept: "image/*", Multiple: true,
		Required: true, Preview: true, UploadURL: "/api/v1/files/upload",
		UploadCategory: "image", Value: "/media/logo.png", CurrentName: "logo.png",
	}).Render(&remote)
	if err != nil {
		t.Fatal(err)
	}
	html := remote.String()
	for _, fragment := range []string{
		`type="hidden"`, `name="logo_url"`, `value="/media/logo.png"`,
		`data-file-upload-value-input="true"`, `data-file-upload-upload-url-value="/api/v1/files/upload"`,
		`data-file-upload-upload-category-value="image"`,
		`data-file-upload-current-url-value="/media/logo.png"`,
		`data-file-upload-current-name-value="logo.png"`, `data-file-upload-multiple-value="false"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("remote FileUpload is missing %q: %s", fragment, html)
		}
	}

	no := false
	var simple strings.Builder
	err = FileUpload(molecules.FileUploadProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Name:           "attachment", Label: "Attachment", Required: true,
		DropZone: &no, ShowList: &no,
	}).Render(&simple)
	if err != nil {
		t.Fatal(err)
	}
	simpleHTML := simple.String()
	for _, fragment := range []string{
		`data-file-upload-disabled-value="true"`, `aria-disabled="true"`,
		`type="file"`, `name="attachment"`, `required`, `disabled`,
		`for="attachment"`,
	} {
		if !strings.Contains(simpleHTML, fragment) {
			t.Errorf("simple FileUpload is missing %q: %s", fragment, simpleHTML)
		}
	}
	for _, forbidden := range []string{`data-file-upload-dropzone`, `data-file-upload-list`} {
		if strings.Contains(simpleHTML, forbidden) {
			t.Errorf("simple FileUpload unexpectedly emitted %q: %s", forbidden, simpleHTML)
		}
	}
}

func TestAutocompleteOwnsComboboxControllerAndStaticOptionContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Autocomplete(molecules.AutocompleteProps{
		ComponentProps: contracts.ComponentProps{ID: "owner-field", Class: "client-owner"},
		Name:           "owner_id", Label: "Owner", Placeholder: "Search users...", Required: true,
		Value: "user-1", DisplayValue: "Ada Lovelace", MinChars: 1,
		Options: []molecules.Option{
			{Value: "user-1", Label: "Ada Lovelace"},
			{Value: "user-2", Label: "Grace Hopper", Disabled: true},
		},
		HelpText: "Choose one teammate.",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="owner-field"`, `class="flex flex-col gap-1.5 client-owner"`,
		`data-component="autocomplete"`, `data-controller="autocomplete"`,
		`data-autocomplete-min-chars-value="1"`, `data-autocomplete-static-value="true"`,
		`data-autocomplete-initial-query-value="Ada Lovelace"`,
		`data-autocomplete-initial-value-value="user-1"`,
		`type="hidden"`, `name="owner_id"`, `value="user-1"`,
		`id="owner_id-input"`, `name="q"`, `role="combobox"`,
		`aria-autocomplete="list"`, `aria-controls="owner_id-listbox"`,
		`aria-expanded="false"`, `autocomplete="off"`, `required`,
		`input-&gt;autocomplete#handleInput`, `keydown-&gt;autocomplete#handleKeydown`,
		`data-autocomplete-panel="true"`, `data-autocomplete-listbox="true"`,
		`htmx:afterSwap-&gt;autocomplete#resultsLoaded`,
		`data-autocomplete-value="user-1"`, `data-autocomplete-label="Ada Lovelace"`,
		`aria-disabled="true"`, `aria-describedby="owner_id-input-help"`,
		`Choose one teammate.`, `data-pk-icon="search"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Autocomplete is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `hx-get=`) {
		t.Fatalf("static autocomplete unexpectedly emits an HTMX request: %s", html)
	}
}

func TestAutocompleteServerSearchHasSafeDefaultsAndAccessibleError(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Autocomplete(molecules.AutocompleteProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Name:           "company_id", Label: "Company", SearchURL: "/companies/search",
		QueryName: "term", Error: "Choose a valid company.",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`name="term"`, `hx-get="/companies/search"`,
		`hx-trigger="keyup changed delay:300ms"`, `hx-target="#company_id-listbox"`,
		`hx-swap="innerHTML"`, `hx-indicator="#company_id-indicator"`,
		`id="company_id-indicator"`, `htmx-indicator`, `disabled`,
		`aria-invalid="true"`, `aria-describedby="company_id-input-error"`,
		`role="alert"`, `Choose a valid company.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("server-backed Autocomplete is missing %q: %s", fragment, html)
		}
	}
}

func TestDropdownOwnsAccessibleListboxAndSubmissionContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Dropdown(molecules.DropdownProps{
		ComponentProps: contracts.ComponentProps{ID: "topic-picker", Class: "client-picker"},
		Name:           "topics", Label: "Topics", Searchable: true, Clearable: true, Multiple: true,
		Selected: []string{"alerts", "release", "alerts"}, Size: "large",
		Options: []molecules.Option{
			{Value: "alerts", Label: "Alerts", Icon: "bell", Group: "Operations"},
			{Value: "release", Label: "Releases", Group: "Operations"},
			{Value: "billing", Label: "Billing", Group: "Finance", Disabled: true},
		},
		HTMXProps: contracts.HTMXProps{Get: "/issues", Target: "#issue-grid", Swap: "outerHTML"},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="topic-picker"`, `client-picker`, `data-component="dropdown"`,
		`data-controller="dropdown"`, `data-dropdown-contract="1"`,
		`data-dropdown-kind-value="select"`, `data-dropdown-multiple-value="true"`,
		`data-dropdown-searchable-value="true"`,
		`data-dropdown-selected-value="[&#34;alerts&#34;,&#34;release&#34;]"`,
		`data-dropdown-input-name-value="topics"`, `data-state="closed"`,
		`aria-haspopup="listbox"`, `aria-expanded="false"`,
		`aria-controls="topic-picker-panel"`, `aria-label="Topics"`,
		`data-dropdown-trigger-label="true"`, `>Alerts, Releases</span>`,
		`data-dropdown-clear-button="true"`, `data-dropdown-chevron-button="true"`,
		`data-dropdown-panel="true"`, `role="listbox"`, `aria-multiselectable="true"`,
		`data-dropdown-search-input="true"`, `input-&gt;dropdown#filterOptions`,
		`role="group"`, `aria-label="Operations"`, `data-dropdown-option="alerts"`,
		`data-dropdown-option-icon="bell"`, `aria-selected="true"`,
		`data-dropdown-option="billing"`, `aria-disabled="true"`,
		`data-dropdown-hidden-inputs="true"`, `data-dropdown-hidden-value="true"`,
		`name="topics"`, `value="alerts"`, `value="release"`,
		`hx-get="/issues"`, `hx-target="#issue-grid"`, `hx-swap="outerHTML"`,
		`hx-trigger="dropdown:changed"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Dropdown is missing %q: %s", fragment, html)
		}
	}
	triggerIndex := strings.Index(html, `data-dropdown-button="true"`)
	triggerEnd := strings.Index(html[triggerIndex:], `</button>`)
	clearIndex := strings.Index(html, `data-dropdown-clear-button="true"`)
	if triggerIndex < 0 || triggerEnd < 0 || clearIndex < triggerIndex+triggerEnd {
		t.Fatalf("Dropdown action buttons must remain outside the focusable trigger button: %s", html)
	}
}

func TestDropdownSingleAndDisabledDefaultsRemainDeterministic(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Dropdown(molecules.DropdownProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Name:           "country", AriaLabel: "Country", Value: "pt",
		Options: []molecules.Option{{Value: "pt", Label: "Portugal"}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`data-dropdown-placeholder-value="Select an option"`,
		`data-dropdown-multiple-value="false"`, `data-dropdown-searchable-value="false"`,
		`data-dropdown-selected-value="[&#34;pt&#34;]"`, `>Portugal</span>`,
		`data-dropdown-hidden-input="true"`, `name="country"`, `value="pt"`,
		`aria-disabled="true"`, `disabled`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical disabled Dropdown is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `data-dropdown-clear-button`) {
		t.Fatalf("non-clearable Dropdown emitted a clear action: %s", html)
	}
}

func TestActionMenuOwnsAccessibleMenuButtonAndHTMXContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := ActionMenu(molecules.ActionMenuProps{
		ComponentProps: contracts.ComponentProps{ID: "job-42-actions", Class: "client-actions"},
		Align:          "start", Width: "lg",
		Items: []molecules.ActionMenuItem{{Label: "Details", Icon: "eye", Href: "/jobs/42"}},
		Sections: []molecules.ActionMenuSection{{
			Label: "Job actions",
			Items: []molecules.ActionMenuItem{{
				Label: "Retry", Icon: "arrow-path", Action: "retry-job",
				HxPost: "/jobs/42/retry", HxTarget: "#jobs", HxSwap: "innerHTML", HxConfirm: "Retry now?",
				Attrs: map[string]string{"data-job-id": "42"},
			}},
		}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="job-42-actions"`, `class="relative inline-block text-left client-actions"`,
		`data-component="action-menu"`, `data-controller="action-menu"`, `data-state="closed"`,
		`data-action-menu-align-value="start"`, `data-action-menu-trigger`,
		`aria-haspopup="menu"`, `aria-expanded="false"`, `aria-label="Actions"`,
		`aria-controls="job-42-actions-panel"`, `data-pk-icon="ellipsis-vertical"`,
		`id="job-42-actions-panel"`, `data-action-menu-panel`, `role="menu"`,
		`aria-orientation="vertical"`, `hidden`, `left-0`, `w-64`,
		`role="group"`, `aria-label="Job actions"`, `role="separator"`,
		`role="menuitem"`, `href="/jobs/42"`, `data-item-action="retry-job"`,
		`hx-post="/jobs/42/retry"`, `hx-target="#jobs"`, `hx-swap="innerHTML"`,
		`hx-confirm="Retry now?"`, `data-job-id="42"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical ActionMenu is missing %q: %s", fragment, html)
		}
	}
}

func TestActionMenuDisabledStateCannotNavigateOrMutate(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := ActionMenu(molecules.ActionMenuProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Items:          []molecules.ActionMenuItem{{Label: "Delete", Href: "/danger", Danger: true, HxDelete: "/danger"}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{`aria-disabled="true"`, `disabled`, `cursor-not-allowed`, `role="menuitem"`} {
		if !strings.Contains(html, fragment) {
			t.Errorf("disabled ActionMenu is missing %q: %s", fragment, html)
		}
	}
	for _, fragment := range []string{`href="/danger"`, `hx-delete="/danger"`, `text-fg-danger`} {
		if strings.Contains(html, fragment) {
			t.Errorf("disabled ActionMenu must not expose %q: %s", fragment, html)
		}
	}
}

func TestDrawerOwnsAccessibleOverlayAndFocusContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := DrawerWithSlots(molecules.DrawerProps{
		ComponentProps: contracts.ComponentProps{ID: "account-drawer", Class: "client-drawer"},
		Title:          "Account", Description: "Manage your profile.", Position: "left", Size: "large",
		CloseLabel: "Close account panel", Open: true, OpenOnSwap: true,
	}, DrawerSlots{
		Body:   []g.Node{h.P(g.Text("Rich body"))},
		Footer: []g.Node{h.Button(g.Text("Save"))},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="account-drawer"`, `client-drawer`, `data-component="drawer"`,
		`data-controller="drawer"`, `data-drawer-id-value="account-drawer"`,
		`data-drawer-position-value="left"`, `data-drawer-open-value="true"`,
		`data-drawer-close-on-escape-value="true"`, `data-state="open"`,
		`role="dialog"`, `aria-modal="true"`, `aria-hidden="false"`,
		`aria-labelledby="account-drawer-title"`, `id="account-drawer-title"`,
		`data-action="htmx:afterSwap-&gt;drawer#openDrawer"`,
		`data-drawer-backdrop`, `data-action="click-&gt;drawer#close"`,
		`data-drawer-panel`, `tabindex="-1"`, `left-0`, `max-w-lg`,
		`aria-label="Close account panel"`, `data-pk-icon="x"`,
		`data-drawer-body`, `Rich body`, `data-drawer-footer`, `Save`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Drawer is missing %q: %s", fragment, html)
		}
	}
}

func TestDrawerExplicitFalseOptionsRemoveDismissalAffordances(t *testing.T) {
	t.Parallel()

	no := false
	var rendered strings.Builder
	err := Drawer(molecules.DrawerProps{
		ComponentProps: contracts.ComponentProps{ID: "locked", Disabled: true},
		AriaLabel:      "Required decision", Position: "bottom", Size: "full",
		Closable: &no, CloseOnOverlay: &no, CloseOnEscape: &no, ShowOverlay: &no,
		Body: "Choose one option.",
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`aria-label="Required decision"`, `aria-disabled="true"`,
		`data-drawer-open-value="false"`, `data-drawer-close-on-escape-value="false"`,
		`data-state="closed"`, `hidden`, `bottom-0`, `h-full`, `w-full`,
		`tabindex="-1"`, `Choose one option.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("non-dismissible Drawer is missing %q: %s", fragment, html)
		}
	}
	for _, fragment := range []string{`data-drawer-backdrop`, `data-drawer-close=""`, `click-&gt;drawer#close`} {
		if strings.Contains(html, fragment) {
			t.Errorf("non-dismissible Drawer must not expose %q: %s", fragment, html)
		}
	}
}

func TestModalOwnsAccessibleOverlayAndServerBehaviorContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := ModalWithSlots(molecules.ModalProps{
		ComponentProps: contracts.ComponentProps{ID: "archive-modal", Class: "client-modal"},
		Title:          "Archive project", Description: "This cannot be undone.", Size: "large",
		CloseLabel: "Close archive dialog", Open: true, OpenOnSwap: true,
	}, ModalSlots{
		Body:   []g.Node{h.P(g.Text("Rich body"))},
		Footer: []g.Node{h.Button(g.Text("Archive"))},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="archive-modal"`, `client-modal`, `data-component="modal"`,
		`data-controller="htmx-modal"`, `data-htmx-modal-open-value="true"`,
		`data-htmx-modal-close-on-escape-value="true"`, `data-htmx-modal-clear-on-close-value="false"`,
		`data-state="open"`, `role="dialog"`, `aria-modal="true"`, `aria-hidden="false"`,
		`aria-labelledby="archive-modal-title"`, `id="archive-modal-title"`,
		`data-action="htmx:afterSwap-&gt;htmx-modal#show"`, `data-modal-backdrop`,
		`data-action="click-&gt;htmx-modal#close"`, `data-modal-panel`, `tabindex="-1"`,
		`max-w-3xl`, `aria-label="Close archive dialog"`, `data-pk-icon="x"`,
		`data-modal-body`, `Rich body`, `data-modal-footer`, `Archive`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Modal is missing %q: %s", fragment, html)
		}
	}
}

func TestModalExplicitFalseOptionsRemoveDismissalAffordances(t *testing.T) {
	t.Parallel()

	no := false
	var rendered strings.Builder
	err := Modal(molecules.ModalProps{
		ComponentProps: contracts.ComponentProps{ID: "locked", Disabled: true},
		AriaLabel:      "Required decision", Size: "full", Body: "Choose one option.",
		Closable: &no, CloseOnOverlay: &no, CloseOnEscape: &no,
		ShowClose: &no, ShowOverlay: &no, Centered: &no,
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`aria-label="Required decision"`, `aria-disabled="true"`,
		`data-htmx-modal-open-value="false"`, `data-htmx-modal-close-on-escape-value="false"`,
		`data-state="closed"`, `hidden`, `style="display:none"`, `items-end`,
		`sm:items-center`, `max-w-full`, `tabindex="-1"`, `Choose one option.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("non-dismissible Modal is missing %q: %s", fragment, html)
		}
	}
	for _, fragment := range []string{`data-modal-backdrop`, `data-modal-close`, `click-&gt;htmx-modal#close`} {
		if strings.Contains(html, fragment) {
			t.Errorf("non-dismissible Modal must not expose %q: %s", fragment, html)
		}
	}
}

func TestModalDeferredRootAndPanelFragmentsShareOneControllerContract(t *testing.T) {
	t.Parallel()

	var rootRendered strings.Builder
	err := Modal(molecules.ModalProps{
		ComponentProps: contracts.ComponentProps{ID: "entity-form-modal"},
		AriaLabel:      "Entity form", Deferred: true, OpenOnSwap: true,
	}).Render(&rootRendered)
	if err != nil {
		t.Fatal(err)
	}
	rootHTML := rootRendered.String()
	for _, fragment := range []string{
		`id="entity-form-modal"`, `data-controller="htmx-modal"`,
		`data-htmx-modal-clear-on-close-value="true"`,
		`data-action="htmx:afterSwap-&gt;htmx-modal#show click-&gt;htmx-modal#backdropClick"`,
		`hidden`, `style="display:none"`,
	} {
		if !strings.Contains(rootHTML, fragment) {
			t.Errorf("deferred Modal root is missing %q: %s", fragment, rootHTML)
		}
	}
	for _, fragment := range []string{`data-modal-panel`, `role="dialog"`} {
		if strings.Contains(rootHTML, fragment) {
			t.Errorf("deferred Modal root must initially be empty of %q: %s", fragment, rootHTML)
		}
	}

	var panelRendered strings.Builder
	err = ModalPanel(molecules.ModalProps{Title: "Edit entity"}, h.P(g.Text("Form body"))).Render(&panelRendered)
	if err != nil {
		t.Fatal(err)
	}
	panelHTML := panelRendered.String()
	for _, fragment := range []string{
		`data-modal-panel`, `role="dialog"`, `aria-modal="true"`,
		`aria-label="Edit entity"`, `click-&gt;htmx-modal#stopPropagation`, `Form body`,
	} {
		if !strings.Contains(panelHTML, fragment) {
			t.Errorf("server Modal panel is missing %q: %s", fragment, panelHTML)
		}
	}
}

func TestModalActionHelpersUseControllerWithoutInlineJavaScript(t *testing.T) {
	t.Parallel()

	for name, node := range map[string]g.Node{
		"close":  ModalCloseButton("Dismiss", "client-close"),
		"cancel": ModalCancelButton("Go back", "client-cancel"),
		"form":   ModalForm(h.ID("edit-form")),
	} {
		var rendered strings.Builder
		if err := node.Render(&rendered); err != nil {
			t.Fatalf("render %s helper: %v", name, err)
		}
		html := rendered.String()
		if strings.Contains(html, `onclick=`) {
			t.Errorf("%s helper must avoid inline JavaScript: %s", name, html)
		}
		if !strings.Contains(html, `htmx-modal#`) {
			t.Errorf("%s helper must use the canonical controller: %s", name, html)
		}
	}
}

func TestSidebarContentSlotsSearchHooksAndSafeItemData(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := SidebarWithSlots(molecules.SidebarProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "docs",
			Attrs: map[string]string{"data-docs-search-root": "sidebar"},
		},
		Flavor:          "content",
		Current:         "#figma",
		NavigationLabel: "Docs navigation",
		Sections: []molecules.SidebarSection{{
			Label: "Runtime", Glyph: "R", Tone: "info", SearchText: "runtime architecture",
			Items: []molecules.SidebarItem{
				{Label: "Overview", Href: "#overview", Prefix: "01"},
				{Label: "Figma handoff", Href: "#figma", Prefix: "02", SearchText: "figma", Icon: `bad<script`},
			},
		}},
	}, SidebarSlots{
		Brand:  []g.Node{h.Strong(g.Text("Architecture decisions"))},
		Footer: []g.Node{h.Small(g.Text("Updated today"))},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="docs"`, `data-sidebar-flavor="content"`, `aria-label="Docs navigation"`,
		`data-docs-search-root="sidebar"`, `<strong>Architecture decisions</strong>`,
		`data-sidebar-section="runtime"`, `data-sidebar-search-section="true"`,
		`data-sidebar-search-text="runtime architecture"`,
		`data-sidebar-search-item="true"`, `data-sidebar-item-prefix="true"`,
		`data-sidebar-item-label="true"`, `data-active="true"`,
		`data-sidebar-footer=""`, `<small>Updated today</small>`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("content Sidebar is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, "bad<script") {
		t.Fatalf("invalid icon names must not become markup: %s", html)
	}
}

func TestSidebarCollapsedAndDisabledStatesRemainAccessible(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Sidebar(molecules.SidebarProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Collapsible:    true,
		Collapsed:      true,
		Items: []molecules.SidebarItem{{
			Label: "Reports", Href: "/admin/reports", Icon: "chart",
		}},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`id="admin-sidebar-collapsible"`, `data-sidebar-collapsible="true"`,
		`data-sidebar-collapsed="true"`, `data-state="collapsed"`,
		`aria-expanded="false"`, `aria-disabled="true"`,
		`aria-label="Reports"`, `data-disabled="true"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("collapsed Sidebar is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `href="/admin/reports"`) {
		t.Fatalf("disabled Sidebar must not emit operable links: %s", html)
	}
}

func TestBreadcrumbOmitsEmptyNavigationLandmark(t *testing.T) {
	t.Parallel()

	if node := Breadcrumb(molecules.BreadcrumbProps{}); node != nil {
		t.Fatalf("empty breadcrumb = %#v, want nil", node)
	}
}

func TestCardOwnsSectionsMediaAndSurfaceVariants(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := CardWithSlots(molecules.CardProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "account-card",
			Class: "client-card",
			Attrs: map[string]string{"data-client": "collect"},
		},
		Image:         "/avatar.webp",
		ImageAlt:      "Account owner",
		ImagePosition: "left",
		Variant:       "plain",
		Padding:       "small",
		Shadow:        "none",
		Hoverable:     true,
	}, CardSlots{
		Header:  []g.Node{h.H2(g.Text("Account"))},
		Content: []g.Node{h.P(g.Text("Current plan"))},
		Footer:  []g.Node{h.Button(g.Text("Manage"))},
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`id="account-card"`,
		`class="`,
		`client-card`,
		`data-client="collect"`,
		`data-component="card"`,
		`src="/avatar.webp"`,
		`alt="Account owner"`,
		`>Account<`,
		`>Current plan<`,
		`>Manage<`,
		clCardPadSmall.Compile(),
		clCardHeader.Compile(),
		clCardFooter.Compile(),
		clCardHoverable.Compile(),
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("CardWithSlots output is missing %q: %s", fragment, html)
		}
	}
	for _, fragment := range []string{clCardBorder.Compile(), clCardShadowSmall.Compile()} {
		if strings.Contains(html, fragment) {
			t.Errorf("plain shadowless CardWithSlots unexpectedly contains %q: %s", fragment, html)
		}
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

func TestTooltipOwnsTriggerBehaviorAndAccessibilityContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Tooltip(atoms.TooltipProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "billing-help",
			Class: "client-tooltip",
			Attrs: map[string]string{"data-owner": "collect"},
		},
		Content:  "Invoices are issued monthly.",
		Position: "right",
		Delay:    150,
	}, h.Button(g.Text("Billing help"))).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`id="billing-help"`,
		`client-tooltip`,
		`data-owner="collect"`,
		`data-component="tooltip"`,
		`data-controller="tooltip"`,
		`data-tooltip-delay-value="150"`,
		`data-tooltip-position="right"`,
		`mouseenter-&gt;tooltip#scheduleShow`,
		`focusin-&gt;tooltip#scheduleShow`,
		`data-tooltip-trigger="true"`,
		`aria-describedby="billing-help-tooltip"`,
		`id="billing-help-tooltip"`,
		`data-tooltip-popup="true"`,
		`role="tooltip"`,
		`aria-hidden="true"`,
		`hidden=""`,
		`<button>Billing help</button>`,
		`Invoices are issued monthly.`,
		clTooltipPosition["right"].Compile(),
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Tooltip output is missing %q: %s", fragment, html)
		}
	}
	for _, forbidden := range []string{`x-data=`, `x-show=`, `@mouseenter`} {
		if strings.Contains(html, forbidden) {
			t.Errorf("Tooltip output still contains legacy behavior %q: %s", forbidden, html)
		}
	}
}

func TestTooltipFallsBackWithoutOrphanPopup(t *testing.T) {
	t.Parallel()

	trigger := h.Button(g.Text("Plain trigger"))
	var rendered strings.Builder
	if err := Tooltip(atoms.TooltipProps{}, trigger).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	if html := rendered.String(); html != `<button>Plain trigger</button>` {
		t.Fatalf("empty Tooltip rendered an orphan popup: %s", html)
	}
	if node := Tooltip(atoms.TooltipProps{Content: "No owner"}); node != nil {
		t.Fatalf("triggerless Tooltip = %#v, want nil", node)
	}
}

func TestRadioOwnsNativeGroupAndAccessibleNamingContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Radio(atoms.RadioProps{
		ComponentProps: contracts.ComponentProps{
			ID:    "plan-pro",
			Class: "client-radio",
			Attrs: map[string]string{"data-owner": "collect"},
		},
		Name:     "plan",
		Label:    "Pro plan",
		Value:    "pro",
		Checked:  true,
		Required: true,
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`<label`,
		`client-radio`,
		`data-component="radio"`,
		`data-state="checked"`,
		`type="radio"`,
		`id="plan-pro"`,
		`name="plan"`,
		`value="pro"`,
		`data-owner="collect"`,
		`checked`,
		`required`,
		`Pro plan`,
		`accent-fg-brand`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("Radio output is missing %q: %s", fragment, html)
		}
	}
}

func TestRadioWithoutVisibleLabelUsesItsValueAsAccessibleName(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := Radio(atoms.RadioProps{
		ComponentProps: contracts.ComponentProps{Class: "compact-radio"},
		Name:           "delivery",
		Value:          "express",
	}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{
		`<input`, `type="radio"`, `class="`, `compact-radio`, `aria-label="express"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("unlabelled Radio output is missing %q: %s", fragment, html)
		}
	}
	if strings.Contains(html, `<label`) {
		t.Fatalf("unlabelled Radio unexpectedly rendered a label wrapper: %s", html)
	}
}

func TestSliderOwnsNativeRangeAndLiveReadoutContract(t *testing.T) {
	t.Parallel()
	node := Slider(atoms.SliderProps{
		ComponentProps: contracts.ComponentProps{ID: "cpu", Attrs: map[string]string{"data-op-kind": "limit"}},
		Name:           "cpu_limit", Label: "CPU limit", Min: 10, Max: 500, Step: 5,
		Value: 100, ShowValue: true, Tone: "warning", AriaValueText: "100 millicores",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`data-component="slider"`,
		`data-controller="slider"`,
		`data-slider-input="true"`,
		`data-slider-value="true"`,
		`input-&gt;slider#sync change-&gt;slider#sync`,
		`id="cpu-input"`,
		`for="cpu-input"`,
		`name="cpu_limit"`,
		`min="10"`,
		`max="500"`,
		`step="5"`,
		`value="100"`,
		`aria-valuetext="100 millicores"`,
		`accent-surface-warning`,
		`data-op-kind="limit"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Slider missing %s: %s", fragment, html)
		}
	}
}

func TestSliderWithoutVisibleLabelUsesNameAsAccessibleName(t *testing.T) {
	t.Parallel()
	node := Slider(atoms.SliderProps{Name: "volume", Min: 0, Max: 100, Value: 40})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `aria-label="volume"`) {
		t.Fatalf("unlabelled Slider has no name-derived accessible label: %s", output.String())
	}
}

func TestToggleOwnsOneAccessibleControlAndSynchronizedStateContract(t *testing.T) {
	t.Parallel()
	node := Toggle(atoms.ToggleProps{
		ComponentProps: contracts.ComponentProps{ID: "dark-mode", Class: "client-toggle", Attrs: map[string]string{"data-owner": "settings"}},
		Name:           "dark_mode", Label: "Dark mode", Checked: true, Size: "lg",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`id="dark-mode"`,
		`client-toggle`,
		`data-owner="settings"`,
		`data-component="toggle"`,
		`data-controller="switch"`,
		`data-state="checked"`,
		`data-size="lg"`,
		`data-switch-checked-value="true"`,
		`data-switch-active-class-value="bg-surface-brand border-surface-brand"`,
		`data-switch-checked-knob-class-value="translate-x-6"`,
		`type="checkbox"`,
		`name="dark_mode"`,
		`value="true"`,
		`tabindex="-1"`,
		`aria-hidden="true"`,
		`<button`,
		`role="switch"`,
		`aria-checked="true"`,
		`aria-label="Dark mode"`,
		`data-switch-button="true"`,
		`data-switch-track="true"`,
		`data-switch-knob="true"`,
		`click-&gt;switch#toggle`,
		`Dark mode`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Toggle missing %s: %s", fragment, html)
		}
	}
	if strings.Contains(html, `<label`) {
		t.Fatalf("Toggle must not nest a button inside a label: %s", html)
	}
	if count := strings.Count(html, `role="switch"`); count != 1 {
		t.Fatalf("Toggle exposes %d switch roles, want one: %s", count, html)
	}
}

func TestToggleHiddenLabelKeepsAccessibleName(t *testing.T) {
	t.Parallel()
	node := Toggle(atoms.ToggleProps{Name: "marketing", Label: "Marketing cookies", HideLabel: true})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if !strings.Contains(html, `aria-label="Marketing cookies"`) {
		t.Fatalf("hidden-label Toggle lost its accessible name: %s", html)
	}
	if strings.Contains(html, `>Marketing cookies</span>`) {
		t.Fatalf("hidden-label Toggle rendered duplicate visible text: %s", html)
	}
}

func TestSpinnerOwnsAccessibleStatusAndLabelDefault(t *testing.T) {
	t.Parallel()

	var labelled strings.Builder
	if err := Spinner(atoms.SpinnerProps{Label: "Fetching results", Size: "lg", Tone: "info"}).Render(&labelled); err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{`role="status"`, `aria-label="Fetching results"`} {
		if !strings.Contains(labelled.String(), fragment) {
			t.Errorf("canonical Spinner missing %s: %s", fragment, labelled.String())
		}
	}

	var defaultLabel strings.Builder
	if err := Spinner(atoms.SpinnerProps{}).Render(&defaultLabel); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(defaultLabel.String(), `aria-label="Loading"`) {
		t.Fatalf("canonical Spinner lost its accessible default label: %s", defaultLabel.String())
	}
}

func TestProgressOwnsNormalizedAccessibleDeterminateContract(t *testing.T) {
	t.Parallel()
	node := Progress(atoms.ProgressProps{
		ComponentProps: contracts.ComponentProps{
			ID: "upload", Class: "client-progress", Attrs: map[string]string{"data-owner": "importer"},
		},
		Value: 25, Max: 50, Label: "Upload progress", ShowText: true,
		Tone: "warning", Size: "lg",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`id="upload"`,
		`client-progress`,
		`data-owner="importer"`,
		`data-component="progress"`,
		`role="progressbar"`,
		`aria-labelledby="upload-label"`,
		`aria-valuemin="0"`,
		`aria-valuemax="50"`,
		`aria-valuenow="25"`,
		`data-progress-percent="50"`,
		`data-state="determinate"`,
		`data-size="lg"`,
		`data-tone="warning"`,
		`id="upload-label"`,
		`data-progress-label="true"`,
		`data-progress-value="true"`,
		`50%`,
		`bg-surface-warning`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Progress missing %s: %s", fragment, html)
		}
	}
	if strings.Contains(html, `style=`) || strings.Contains(html, `pk-progress-fill`) {
		t.Fatalf("Progress must not depend on inline width or host-private CSS: %s", html)
	}
}

func TestProgressIndeterminateOmitsValueAndKeepsAccessibleName(t *testing.T) {
	t.Parallel()
	node := Progress(atoms.ProgressProps{
		AriaLabel: "Preparing export", ShowText: true, Indeterminate: true,
		Tone: "info",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`aria-label="Preparing export"`,
		`aria-busy="true"`,
		`data-state="indeterminate"`,
		`data-tone="info"`,
		`animate-pulse`,
		`w-full`,
		`bg-surface-info`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("indeterminate Progress missing %s: %s", fragment, html)
		}
	}
	for _, forbidden := range []string{`aria-valuenow=`, `data-progress-value=`, `%</span>`} {
		if strings.Contains(html, forbidden) {
			t.Errorf("indeterminate Progress unexpectedly contains %s: %s", forbidden, html)
		}
	}
}

func TestProgressClampsValuesAndRendersCanonicalTone(t *testing.T) {
	t.Parallel()
	node := Progress(atoms.ProgressProps{Value: 150, Max: 100, AriaLabel: "Completion", Tone: "danger"})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`aria-valuenow="100"`, `data-progress-percent="100"`,
		`data-tone="danger"`, `bg-surface-danger`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("normalized Progress missing %s: %s", fragment, html)
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

func TestTextareaOwnsValidationCounterAndAutoresizeContract(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	err := Textarea(atoms.TextareaProps{
		ComponentProps: contracts.ComponentProps{
			Disabled: true,
			Attrs:    map[string]string{"data-field-owner": "profile"},
		},
		HTMXProps: contracts.HTMXProps{Post: "/notes", Trigger: "blur"},
		Name:      "notes", Label: "Notes", Value: "hello", Placeholder: "Add notes",
		HelperText: "Visible to administrators.", ErrorMessage: "Add more detail.",
		Required: true, ReadOnly: true, MinLength: 2, MaxLength: 20,
		AutoResize: true, MinRows: 3, MaxRows: 15, FullWidth: true,
	}).Render(&rendered)
	if err != nil {
		t.Fatal(err)
	}

	html := rendered.String()
	for _, fragment := range []string{
		`data-component="textarea"`,
		`data-controller="textarea-counter"`,
		`data-controller="autoresize"`,
		`data-autoresize-min-rows-value="3"`,
		`data-autoresize-max-rows-value="15"`,
		`data-textarea-counter-target="input"`,
		`data-textarea-counter-target="display"`,
		`data-action="input-&gt;autoresize#resize input-&gt;textarea-counter#update"`,
		`aria-describedby="pk-textarea-notes-error pk-textarea-notes-helper"`,
		`id="pk-textarea-notes-error"`,
		`id="pk-textarea-notes-helper"`,
		`role="alert"`,
		`aria-live="polite"`,
		`5 / 20`,
		`rows="3"`,
		`minlength="2"`,
		`maxlength="20"`,
		`required`,
		`readonly`,
		`disabled`,
		`hx-post="/notes"`,
		`hx-trigger="blur"`,
		`data-field-owner="profile"`,
		`resize-none`,
		`w-full`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Textarea missing %s: %s", fragment, html)
		}
	}
	if !strings.Contains(html, "Add more detail.") || !strings.Contains(html, "Visible to administrators.") {
		t.Errorf("canonical Textarea must render error and helper together: %s", html)
	}
}

func TestTextareaDefaultsToManualResizeAndAccessibleName(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := Textarea(atoms.TextareaProps{Name: "summary"}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	html := rendered.String()
	for _, fragment := range []string{`rows="4"`, `resize-y`, `aria-label="summary"`} {
		if !strings.Contains(html, fragment) {
			t.Errorf("default Textarea missing %s: %s", fragment, html)
		}
	}
	if strings.Contains(html, `data-controller="autoresize"`) || strings.Contains(html, `textarea-counter`) {
		t.Errorf("optional Textarea controllers must be state-gated: %s", html)
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
		`data-component="table"`,
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

func TestDataGridOwnsStableSurfaceIdentity(t *testing.T) {
	t.Parallel()

	var rendered strings.Builder
	if err := DataGrid(organisms.DataGridProps{}, DataGridSlots{
		Table: Table(molecules.TableProps{
			Columns: []molecules.TableColumn{{Key: "name", Label: "Name"}},
			Rows:    []molecules.TableRow{{ID: "row-1", Cells: map[string]any{"name": "Ada"}}},
		}),
	}).Render(&rendered); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered.String(), `data-component="data-grid"`) {
		t.Fatalf("DataGrid has no stable component identity: %s", rendered.String())
	}
}

func TestTableWithSlotsCarriesTrustedRowAndCellProjection(t *testing.T) {
	var buf bytes.Buffer
	err := TableWithSlots(molecules.TableProps{
		Selectable: true,
		Sortable:   true,
		Columns:    []molecules.TableColumn{{Key: "name", Label: "Name", Sortable: true}},
		Rows:       []molecules.TableRow{{ID: "row-1", Cells: map[string]any{"name": "Ada"}}},
	}, TableSlots{
		RowAttrs: func(molecules.TableRow) []g.Node { return []g.Node{g.Attr("data-state", "selected")} },
		CellAttrs: func(_ molecules.TableRow, c molecules.TableColumn) []g.Node {
			return []g.Node{g.Attr("data-label", c.Label)}
		},
		SortButtonAttrs:  func(molecules.TableColumn) []g.Node { return []g.Node{g.Attr("hx-include", "[name='search']")} },
		SelectRowChecked: func(molecules.TableRow) bool { return true },
	}).Render(&buf)
	if err != nil {
		t.Fatalf("render table slots: %v", err)
	}
	for _, fragment := range []string{
		`data-state="selected"`, `data-label="Name"`, `hx-include="[name=&#39;search&#39;]"`, `checked`,
	} {
		if !strings.Contains(buf.String(), fragment) {
			t.Errorf("TableWithSlots output is missing %q: %s", fragment, buf.String())
		}
	}
}

func TestSelectOwnsGroupedMultipleAndSupportingTextContract(t *testing.T) {
	var output strings.Builder
	err := Select(atoms.SelectProps{
		ComponentProps: contracts.ComponentProps{ID: "team-filter"},
		Name:           "teams",
		Values:         []string{"platform", "design"},
		Multiple:       true,
		FullWidth:      true,
		Error:          "Choose at least one team.",
		HelpText:       "Use Shift to select a range.",
		Options: []atoms.SelectOption{
			{Value: "platform", Label: "Platform", Group: "Engineering", Description: "Core platform team"},
			{Value: "design", Label: "Design", Group: "Engineering"},
			{Value: "operations", Label: "Operations", Group: "Business", Disabled: true},
		},
	}).Render(&output)
	if err != nil {
		t.Fatalf("render Select: %v", err)
	}

	html := output.String()
	for _, fragment := range []string{
		`aria-label="teams"`,
		`aria-describedby="team-filter-error team-filter-help"`,
		`multiple`,
		`size="4"`,
		`optgroup label="Engineering"`,
		`value="platform" selected`,
		`title="Core platform team"`,
		`value="operations" disabled`,
		`role="alert"`,
		`Choose at least one team.`,
		`Use Shift to select a range.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("Select output missing %q: %s", fragment, html)
		}
	}
}

func TestInputOwnsValidationToneConstraintsAndSupportingText(t *testing.T) {
	var output strings.Builder
	err := Input(atoms.InputProps{
		ComponentProps: contracts.ComponentProps{ID: "email-field"},
		Name:           "email",
		Type:           "EMAIL",
		Value:          "name@example.com",
		ReadOnly:       true,
		FullWidth:      true,
		Tone:           "warning",
		Error:          "Review this address.",
		HelpText:       "Use your work email.",
		MinLength:      5,
		MaxLength:      120,
		Pattern:        `.+@.+`,
		Autocomplete:   "email",
	}).Render(&output)
	if err != nil {
		t.Fatalf("render Input: %v", err)
	}

	html := output.String()
	for _, fragment := range []string{
		`id="email-field"`,
		`name="email"`,
		`type="email"`,
		`data-tone="warning"`,
		`data-size="md"`,
		`readonly`,
		`minlength="5"`,
		`maxlength="120"`,
		`pattern=".+@.+"`,
		`autocomplete="email"`,
		`aria-label="email"`,
		`aria-invalid="true"`,
		`aria-describedby="email-field-error email-field-help"`,
		`role="alert"`,
		`Review this address.`,
		`Use your work email.`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("Input output missing %q: %s", fragment, html)
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
		// HTMX owns this runtime sentinel and its visibility rule; it is not a
		// visual utility and therefore deliberately does not enter tw emission.
		if class == "htmx-indicator" {
			continue
		}
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

func TestAlertOwnsCanonicalInteractiveAndVisualStates(t *testing.T) {
	t.Parallel()

	node := Alert(atoms.AlertProps{
		Title:       "Check required",
		Message:     "Review this change.",
		Tone:        "warning",
		Dismissible: true,
		Bordered:    true,
		Compact:     true,
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`role="alert"`,
		`aria-live="assertive"`,
		`data-controller="alert"`,
		`data-alert-dismissible-value="true"`,
		`data-alert-icon=""`,
		`data-pk-icon="exclamation-triangle"`,
		`data-action="click-&gt;alert#dismiss"`,
		`data-alert-close=""`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Alert missing %s: %s", fragment, html)
		}
	}
	for _, classList := range []string{
		clAlertCompact.Compile(),
		clAlertBordered.Compile(),
		clAlertIcon.Compile(),
		clAlertClose.Compile(),
	} {
		for _, class := range strings.Fields(classList) {
			if !strings.Contains(" "+html+" ", class) {
				t.Errorf("canonical Alert missing state class %q: %s", class, html)
			}
		}
	}
	for _, class := range strings.Fields(clAlertRegular.Compile()) {
		if strings.Contains(" "+html+" ", class) {
			t.Errorf("compact Alert retained regular-spacing class %q: %s", class, html)
		}
	}
}

func TestToastOwnsCanonicalUrgencyDismissalAndDefaults(t *testing.T) {
	t.Parallel()

	node := Toast(atoms.ToastProps{
		ComponentProps: contracts.ComponentProps{ID: "save-toast", Class: "custom-toast"},
		Title:          "Check required",
		Message:        "Review this change.",
		Tone:           "warning",
		Duration:       2500,
		Position:       "bottom-left",
		CloseLabel:     "Dismiss update",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`id="save-toast"`,
		`data-component="toast"`,
		`data-controller="toast"`,
		`data-toast-tone="warning"`,
		`data-toast-duration-ms-value="2500"`,
		`data-toast-persistent-value="false"`,
		`data-toast-position-value="bottom-left"`,
		`role="alert"`,
		`aria-live="assertive"`,
		`data-pk-icon="exclamation-triangle"`,
		`data-action="click-&gt;toast#dismiss"`,
		`aria-label="Dismiss update"`,
		`custom-toast`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Toast missing %s: %s", fragment, html)
		}
	}

	defaultHTML := renderNode(t, Toast(atoms.ToastProps{Message: "Done."}))
	for _, fragment := range []string{
		`data-toast-tone="info"`,
		`data-toast-duration-ms-value="5000"`,
		`data-toast-position-value="top-right"`,
		`role="status"`,
		`aria-live="polite"`,
	} {
		if !strings.Contains(defaultHTML, fragment) {
			t.Errorf("default Toast missing %s: %s", fragment, defaultHTML)
		}
	}
}

func TestAvatarOwnsAccessibleFallbackAndPresenceContract(t *testing.T) {
	t.Parallel()

	node := Avatar(atoms.AvatarProps{
		ComponentProps: contracts.ComponentProps{ID: "ada", Class: "custom-avatar"},
		Name:           "Ada Lovelace",
		Size:           "lg",
		Shape:          "rounded",
		Tone:           "brand",
		Status:         "ONLINE",
		StatusPosition: "top-right",
		StatusLabel:    "Ada is online",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`id="ada"`, `custom-avatar`, `data-component="avatar"`,
		`data-avatar-size="lg"`, `data-avatar-shape="rounded"`,
		`data-avatar-tone="brand"`, `role="img"`, `aria-label="Ada Lovelace"`,
		`>AL</span>`, `data-avatar-status="online"`, `aria-label="Ada is online"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Avatar missing %s: %s", fragment, html)
		}
	}
}

func TestAvatarImageAndIconFallbacksKeepNamesOnSemanticContent(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		props atoms.AvatarProps
		want  []string
	}{
		"image": {
			props: atoms.AvatarProps{Src: "/portrait.png", Alt: "Profile portrait"},
			want:  []string{`<img`, `src="/portrait.png"`, `alt="Profile portrait"`},
		},
		"icon": {
			props: atoms.AvatarProps{FallbackIcon: "user", Alt: "Account"},
			want:  []string{`data-pk-icon="user"`, `role="img"`, `aria-label="Account"`},
		},
		"decorative default": {
			props: atoms.AvatarProps{},
			want:  []string{`data-pk-icon="user"`, `aria-hidden="true"`},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var output strings.Builder
			if err := Avatar(test.props).Render(&output); err != nil {
				t.Fatal(err)
			}
			for _, fragment := range test.want {
				if !strings.Contains(output.String(), fragment) {
					t.Errorf("Avatar missing %s: %s", fragment, output.String())
				}
			}
		})
	}
}

func TestBadgeOwnsCountRemovalLiveAndSlotContract(t *testing.T) {
	t.Parallel()

	node := BadgeWithSlots(
		atoms.BadgeProps{
			ComponentProps: contracts.ComponentProps{ID: "messages", Class: "custom-badge"},
			Label:          "Messages",
			Variant:        "secondary",
			Tone:           "success",
			Size:           "sm",
			Dot:            true,
			Count:          125,
			Removable:      true,
			RemoveLabel:    "Clear message filter",
			Live:           true,
		},
		BadgeSlots{IconStart: []g.Node{Icon(atoms.IconProps{Name: "envelope", Size: "xs"})}},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`id="messages"`, `custom-badge`, `data-component="badge"`,
		`data-variant="secondary"`, `data-tone="success"`, `data-size="sm"`,
		`role="status"`, `aria-live="polite"`, `data-badge-dot="true"`,
		`data-pk-icon="envelope"`, `data-badge-count="true">99+`,
		`data-badge-remove="true"`, `aria-label="Clear message filter"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Badge missing %s: %s", fragment, html)
		}
	}
}

func TestTabsOwnCanonicalPanelsAndRejectIconMarkup(t *testing.T) {
	t.Parallel()

	node := TabsWithPanels(
		molecules.TabsProps{ActiveTab: "disabled", Variant: "underline"},
		TabSlot{ID: "disabled", Label: "Disabled", Disabled: true},
		TabSlot{ID: "safe", Label: "Safe", Icon: `<img src=x onerror=alert(1)>`, Content: []g.Node{g.Text("Panel")}},
		TabSlot{ID: "lazy", Label: "Lazy", HxGet: "/lazy"},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`data-controller="tabs"`,
		`data-tabs-contract="1"`,
		`data-tabs-active-tab-value="safe"`,
		`data-action="click-&gt;tabs#activate"`,
		`role="tabpanel"`,
		`hx-trigger="tabs:activate from:this once"`,
		`aria-disabled="true"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Tabs missing %s: %s", fragment, html)
		}
	}
	for _, unsafe := range []string{"<img", "onerror=", "<script"} {
		if strings.Contains(html, unsafe) {
			t.Errorf("tab icon name rendered unsafe markup %q: %s", unsafe, html)
		}
	}
}

func TestCheckboxOwnsIndeterminateControllerProjection(t *testing.T) {
	t.Parallel()

	node := Checkbox(atoms.CheckboxProps{
		ComponentProps: contracts.ComponentProps{ID: "select-some", Class: "w-fit"},
		Name:           "selection",
		Label:          "Select some rows",
		Indeterminate:  true,
		HelpText:       "Some rows are selected.",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`class="flex flex-col gap-1.5"`,
		`class="inline-flex items-start gap-3 cursor-pointer w-fit"`,
		`data-controller="checkbox"`,
		`data-checkbox-indeterminate-value="true"`,
		`data-state="indeterminate"`,
		`id="select-some"`,
		`name="selection"`,
		`aria-checked="mixed"`,
		`data-checkbox-input="true"`,
		`data-checkbox-box="true"`,
		`data-checkbox-checkmark="true"`,
		`data-checkbox-bar="true"`,
		`aria-describedby="select-some-help"`,
		`id="select-some-help"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical Checkbox missing %s: %s", fragment, html)
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
		"atoms.Skeleton",
	}
	if len(unimplemented) == 0 {
		t.Skip("everything implemented")
	}
	// The list is prose-verified: it exists so the package doc and release
	// notes can point at one authoritative statement of scope.
}

func TestPaginationCursorUsesOpaqueContinuationLinks(t *testing.T) {
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

func TestPaginationCursorMaterializesOpaqueContinuationURLs(t *testing.T) {
	t.Parallel()
	node := Pagination(molecules.PaginationProps{
		TotalPages:      0,
		BaseURL:         "/albums?after=stale&filter=caf%C3%A9&tag=a&tag=b#results",
		PreviousCursor:  "previous-token",
		NextCursor:      "pkps:v1:a+b/&?",
		BeforeParameter: "shelfBefore",
		AfterParameter:  "shelfAfter",
		NavigationLabel: "Collection navigation",
	})
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`aria-label="Collection navigation"`,
		`href="/albums?after=stale&amp;filter=caf%C3%A9&amp;shelfBefore=previous-token&amp;tag=a&amp;tag=b#results"`,
		`href="/albums?after=stale&amp;filter=caf%C3%A9&amp;shelfAfter=pkps%3Av1%3Aa%2Bb%2F%26%3F&amp;tag=a&amp;tag=b#results"`,
		`data-cursor-direction="backward"`,
		`data-cursor-direction="forward"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("materialized cursor pagination missing %q: %s", fragment, html)
		}
	}
}

func TestPaginationCursorIntentFailsClosed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		props molecules.PaginationProps
	}{
		{name: "unsafe base URL", props: molecules.PaginationProps{BaseURL: "javascript:alert(1)", NextCursor: "next"}},
		{name: "unsafe direct URL", props: molecules.PaginationProps{NextURL: "javascript:alert(1)"}},
		{name: "blank cursor", props: molecules.PaginationProps{BaseURL: "/albums", NextCursor: " \t "}},
		{name: "oversized cursor", props: molecules.PaginationProps{BaseURL: "/albums", NextCursor: strings.Repeat("x", maxPaginationCursorBytes+1)}},
		{name: "invalid parameter", props: molecules.PaginationProps{BaseURL: "/albums", NextCursor: "next", AfterParameter: "after[]"}},
		{name: "duplicate parameters", props: molecules.PaginationProps{BaseURL: "/albums", NextCursor: "next", AfterParameter: "cursor", BeforeParameter: "cursor"}},
		{name: "parameters without continuation", props: molecules.PaginationProps{BaseURL: "/albums", AfterParameter: "after"}},
		{name: "load more without forward continuation", props: molecules.PaginationProps{BaseURL: "/albums", PreviousCursor: "previous", CursorMode: "load-more"}},
		{name: "invalid mode", props: molecules.PaginationProps{BaseURL: "/albums", NextCursor: "next", CursorMode: "unknown"}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var output strings.Builder
			if err := Pagination(test.props).Render(&output); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(output.String(), "<nav") {
				t.Fatalf("unsafe cursor intent rendered navigation: %s", output.String())
			}
		})
	}
}

func TestPaginationOffsetPreservesQueryAndUsesPerPageHTMXURLs(t *testing.T) {
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

func TestWindowedCollectionProjectsSemanticGridLayout(t *testing.T) {
	t.Parallel()

	node := WindowedCollection(
		organisms.WindowedCollectionProps{
			Layout:      "grid",
			ItemCount:   2,
			MaxItems:    100,
			Title:       "Your stickers",
			Description: "Every sticker you own.",
		},
		WindowedCollectionSlots{Items: []g.Node{
			Text(atoms.TextProps{Content: "First"}),
			Text(atoms.TextProps{Content: "Second"}),
		}},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`data-windowed-layout="grid"`,
		`data-windowed-items`,
		`grid`,
		`sm:grid-cols-2`,
		`lg:grid-cols-3`,
		`Your stickers`,
		`Every sticker you own.`,
		`bg-surface-brand`,
	} {
		if !strings.Contains(html, fragment) {
			t.Fatalf("grid window missing %q: %s", fragment, html)
		}
	}
}

func TestWindowedCollectionDefaultsUnknownLayoutToList(t *testing.T) {
	t.Parallel()

	node := WindowedCollection(
		organisms.WindowedCollectionProps{Layout: "masonry", ItemCount: 1, MaxItems: 100},
		WindowedCollectionSlots{Items: []g.Node{Text(atoms.TextProps{Content: "Only"})}},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if !strings.Contains(html, `data-windowed-layout="list"`) ||
		strings.Contains(html, `sm:grid-cols-2`) {
		t.Fatalf("unknown layout did not fail closed to list: %s", html)
	}
}

func TestDashboardWidgetOwnsTypedStatAndRefreshContract(t *testing.T) {
	t.Parallel()
	node := DashboardWidget(
		organisms.DashboardWidgetProps{
			ComponentProps: contracts.ComponentProps{ID: "seat-widget", Class: "client-widget"},
			Title:          "Active seats", Subtitle: "Current workspace posture",
			Value: "1,234", PreviousValue: "1,102", Change: "+12%",
			Trend: "up", Icon: "activity", Span: "2",
			RefreshURL: "/dashboard/seats", RefreshOn: "seat:changed", RefreshSeconds: 60,
			DetailURL: "/dashboard/seats/detail",
		},
		DashboardWidgetSlots{Footer: []g.Node{g.Text("Updated now")}},
	)
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, fragment := range []string{
		`data-component="dashboard-widget"`, `id="seat-widget"`, `client-widget`,
		`data-widget-type="stat"`, `data-pk-icon="activity"`, `Active seats`,
		`href="/dashboard/seats/detail"`, `data-dashboard-widget-value`, `1,234`,
		`data-dashboard-widget-trend="up"`, `↑ +12%`, `from 1,102`, `Updated now`,
		`hx-get="/dashboard/seats"`, `hx-swap="innerHTML"`,
		`hx-trigger="every 60s, seat:changed from:body"`,
	} {
		if !strings.Contains(html, fragment) {
			t.Errorf("canonical DashboardWidget is missing %q: %s", fragment, html)
		}
	}
}

func renderNodeToString(t *testing.T, node g.Node) string {
	t.Helper()
	var output strings.Builder
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	return output.String()
}

func TestNeutralToneRendersAcrossFeedbackAtoms(t *testing.T) {
	t.Parallel()

	spinner := renderNodeToString(t, Spinner(atoms.SpinnerProps{Tone: "neutral"}))
	if !strings.Contains(spinner, "border-t-fg-secondary") {
		t.Errorf("neutral Spinner missing neutral border tone: %s", spinner)
	}

	progress := renderNodeToString(t, Progress(atoms.ProgressProps{Value: 40, Tone: "neutral"}))
	for _, fragment := range []string{`data-tone="neutral"`, "bg-fg-secondary"} {
		if !strings.Contains(progress, fragment) {
			t.Errorf("neutral Progress missing %s: %s", fragment, progress)
		}
	}

	alert := renderNodeToString(t, Alert(atoms.AlertProps{Message: "Note", Tone: "neutral"}))
	for _, fragment := range []string{`data-alert-tone="neutral"`, `role="status"`, "bg-surface-tertiary"} {
		if !strings.Contains(alert, fragment) {
			t.Errorf("neutral Alert missing %s: %s", fragment, alert)
		}
	}

	toast := renderNodeToString(t, Toast(atoms.ToastProps{Message: "Saved", Tone: "neutral"}))
	if !strings.Contains(toast, "bg-surface-tertiary") {
		t.Errorf("neutral Toast missing neutral surface: %s", toast)
	}
}
