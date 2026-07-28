// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

// delivery_composites.go defines the molecule, organism, template, and page
// half of the OSS delivery catalog. The page is a real executable composition
// of catalog components and a native instance graph—not a flattened drawing.

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	designcomponent "github.com/septagon-oss/pk-design/pkg/components"
	uicomponent "github.com/septagon-oss/pk-ui/component"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/layouts"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	"github.com/septagon-oss/pk-ui/contracts/organisms"
	"github.com/septagon-oss/tw"
)

func tableDeliveryDefinition() DeliveryDefinition {
	componentType := "Table"
	root := deliveryFrame(
		componentType,
		clTableWrap.Compile(),
		deliveryFrame(
			"Table",
			clTable.Compile(),
			deliveryFrame(
				"Header",
				clTableHead.Compile(),
				deliveryText("Column", clTableTh.Compile(), "Name"),
				deliveryText("Column", clTableTh.Compile(), "Status"),
				deliveryText("Column", clTableTh.Compile(), "Updated"),
			),
			deliveryFrame(
				"Row",
				clTableRow.Compile(),
				deliveryText("PrimaryCell", clTableTd.Merge(clTableTdStrong).Compile(), "Alpha workspace"),
				deliveryText("Cell", clTableTd.Compile(), "Active"),
				deliveryText("Cell", clTableTd.Compile(), "Today"),
			),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible tabular data presentation with sorting and selection hooks.",
		map[string]PropertyContract{
			"columns":    contentProperty("Ordered column definitions."),
			"rows":       contentProperty("Ordered data rows."),
			"sortable":   stateProperty("Whether sortable columns expose controls."),
			"selectable": stateProperty("Whether row selection is available."),
			"striped":    modifierProperty("Whether alternate rows are tinted."),
			"compact":    modifierProperty("Whether row density is compact."),
			"emptyText":  contentProperty("Empty-state copy."),
		},
		nil,
		deliveryDesign("data", "03 Molecules", "Data", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "data", tableExampleProps()),
			deliveryExample(componentType, "empty", map[string]any{
				"columns": []map[string]any{
					{"key": "name", "label": "Name", "primary": true},
					{"key": "status", "label": "Status"},
				},
				"rows":      []map[string]any{},
				"emptyText": "No workspaces found.",
			}),
		},
		Table,
	)
}

func cardDeliveryDefinition() DeliveryDefinition {
	componentType := "Card"
	root := deliveryClassBound(
		deliveryFrame(
			componentType,
			clCard.Compile(),
			deliveryText("Title", clCardTitle.Compile(), "Quarterly summary", "title"),
			deliveryText("Description", clCardDesc.Compile(), "Performance is tracking above plan.", "description"),
			deliverySlot(
				"Content",
				"children",
				"",
				deliveryInstance("Body", "Text", "body", ""),
			),
		),
		"clickable",
		map[string]string{"false": "", "true": clCardClickable.Compile()},
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Composable content surface with optional navigation behavior.",
		map[string]PropertyContract{
			"title":       contentProperty("Optional card heading."),
			"description": contentProperty("Optional supporting copy."),
			"image":       contentProperty("Optional media source."),
			"variant":     variantProperty([]string{"default", "elevated", "outlined"}, "default", "Surface treatment."),
			"clickable":   stateProperty("Whether the whole card is interactive."),
			"href":        contentProperty("Navigation target for clickable cards."),
		},
		[]designcomponent.Slot{{
			Name:         "children",
			Description:  "Ordered card content.",
			AllowedTypes: deliveryAtomicContentTypes(),
			Cardinality:  designcomponent.SlotMany,
		}},
		deliveryDesign("summary", "03 Molecules", "Surfaces", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "summary", map[string]any{
				"title": "Quarterly summary", "description": "Performance is tracking above plan.",
			}),
			mobileDeliveryExample(componentType, "mobile-summary", map[string]any{
				"title": "Open reviews", "description": "Three items need attention.",
			}),
		},
		func(props molecules.CardProps) g.Node {
			return Card(props, Text(textProps("Review the current operational details.", "muted")))
		},
	)
}

func breadcrumbDeliveryDefinition() DeliveryDefinition {
	componentType := "Breadcrumb"
	root := deliveryFrame(
		componentType,
		clBreadcrumb.Compile(),
		deliveryText("Parent", clLink.Compile(), "Workspaces"),
		deliveryText("Separator", clBreadcrumbSep.Compile(), "/"),
		deliveryText("Current", clBreadcrumbCur.Compile(), "Alpha"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible hierarchical navigation trail.",
		map[string]PropertyContract{
			"items":     contentProperty("Ordered navigation segments."),
			"separator": contentProperty("Visual separator."),
			"maxItems":  modifierProperty("Maximum visible segments before collapsing."),
		},
		nil,
		deliveryDesign("default", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"items": []map[string]any{
					{"label": "Workspaces", "href": "/workspaces"},
					{"label": "Alpha", "current": true},
				},
			}),
		},
		Breadcrumb,
	)
}

func paginationDeliveryDefinition() DeliveryDefinition {
	componentType := "Pagination"
	root := deliveryFrame(
		componentType,
		clPagination.Compile(),
		deliveryText("Previous", clPageBtn.Merge(clPageIdle).Compile(), "‹"),
		deliveryText("PageOne", clPageBtn.Merge(clPageIdle).Compile(), "1"),
		deliveryText("Current", clPageBtn.Merge(clPageCur).Compile(), "2"),
		deliveryText("PageThree", clPageBtn.Merge(clPageIdle).Compile(), "3"),
		deliveryText("Next", clPageBtn.Merge(clPageIdle).Compile(), "›"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible numbered or cursor pagination.",
		map[string]PropertyContract{
			"currentPage": contentProperty("Current one-based page."),
			"totalPages":  contentProperty("Known total; zero selects cursor mode."),
			"perPage":     modifierProperty("Rows per page."),
			"siblings":    modifierProperty("Visible siblings around the current page."),
			"baseURL":     contentProperty("Navigation base URL."),
		},
		nil,
		deliveryDesign("numbered", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "numbered", map[string]any{
				"currentPage": 2, "totalPages": 7, "siblings": 1, "baseURL": "/workspaces",
			}),
			deliveryExample(componentType, "cursor", map[string]any{
				"currentPage": 1, "totalPages": 0,
			}),
		},
		Pagination,
	)
}

func searchBarDeliveryDefinition() DeliveryDefinition {
	componentType := "SearchBar"
	root := deliveryFrame(
		componentType,
		clSearchWrap.Compile(),
		deliveryText("Icon", clIcon.Compile(), "⌕"),
		deliveryText("Input", clSearchInput.Compile(), "Search workspaces", "placeholder", "label"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Search field shell with progressive-enhancement hooks.",
		map[string]PropertyContract{
			"searchURL":    contentProperty("Search endpoint."),
			"label":        contentProperty("Accessible control name."),
			"placeholder":  contentProperty("Empty-state prompt."),
			"instant":      stateProperty("Whether typing submits progressively."),
			"showClear":    modifierProperty("Whether a clear affordance is shown."),
			"showShortcut": modifierProperty("Whether a keyboard hint is shown."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"searchURL": "/workspaces/search", "label": "Search workspaces",
				"placeholder": "Search workspaces", "instant": true,
			}),
		},
		SearchBar,
	)
}

func tabsDeliveryDefinition() DeliveryDefinition {
	componentType := "Tabs"
	root := deliveryFrame(
		componentType,
		clTabList.Compile(),
		deliveryText("Overview", clTab.Merge(clTabActive).Compile(), "Overview"),
		deliveryText("Activity", clTab.Merge(clTabIdle).Compile(), "Activity"),
		deliveryText("Settings", clTab.Merge(clTabIdle).Compile(), "Settings"),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible tab navigation whose panels remain page-owned.",
		map[string]PropertyContract{
			"items":       contentProperty("Ordered tabs."),
			"activeTab":   stateProperty("Active tab key."),
			"orientation": variantProperty([]string{"horizontal", "vertical"}, "horizontal", "Tab direction."),
			"variant":     variantProperty([]string{"underline", "pills"}, "underline", "Visual treatment."),
		},
		nil,
		deliveryDesign("underline", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "underline", map[string]any{
				"activeTab": "overview", "orientation": "horizontal", "variant": "underline",
				"items": []map[string]any{
					{"key": "overview", "label": "Overview", "url": "/overview"},
					{"key": "activity", "label": "Activity", "url": "/activity"},
					{"key": "settings", "label": "Settings", "url": "/settings"},
				},
			}),
		},
		Tabs,
	)
}

func dataGridDeliveryDefinition() DeliveryDefinition {
	componentType := "DataGrid"
	root := deliveryFrame(
		componentType,
		clGridSection.Compile(),
		deliveryFrame(
			"Toolbar",
			clGridToolbar.Compile(),
			deliveryInstance("Search", "SearchBar", "basic", "flex-1"),
			deliveryInstance("Create", "Button", "primary", clGridCreate.Compile()),
		),
		deliveryInstance("Table", "Table", "data", ""),
		deliverySlot(
			"Status",
			"status",
			"",
			deliveryInstance("Empty", "EmptyState", "compact", ""),
		),
		deliveryInstance("Pagination", "Pagination", "numbered", ""),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Search, filter, action, table, status, and pagination composition.",
		map[string]PropertyContract{
			"table":      contentProperty("Table contract."),
			"pagination": contentProperty("Pagination contract."),
			"search":     contentProperty("Search contract."),
			"actions":    contentProperty("Toolbar actions."),
			"searchURL":  contentProperty("Search endpoint shorthand."),
			"createURL":  contentProperty("Optional create target."),
			"createText": contentProperty("Create-action label."),
			"filters":    contentProperty("Filter definitions."),
		},
		[]designcomponent.Slot{{
			Name:         "status",
			Description:  "Page-owned loading, error, or empty status between data and pagination.",
			AllowedTypes: []string{"Alert", "EmptyState", "Spinner", "Text"},
			Cardinality:  designcomponent.SlotOne,
		}},
		deliveryDesign("workspace-list", "04 Organisms", "Data", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "workspace-list", dataGridExampleProps()),
		},
		func(props organisms.DataGridProps) g.Node {
			return DataGrid(props)
		},
	)
}

func stackDeliveryDefinition() DeliveryDefinition {
	componentType := "Stack"
	root := deliverySlot(
		componentType,
		"children",
		clStack.Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Card", "summary", ""),
		deliveryInstance("Secondary", "Card", "mobile-summary", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Vertical auto-layout primitive.",
		map[string]PropertyContract{
			"gap":   sizeProperty(gaps, "4", "Space between children."),
			"align": variantProperty([]string{"start", "center", "end", "stretch"}, "stretch", "Cross-axis alignment."),
		},
		deliveryDesign("default", "05 Templates", "Layout", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"gap": "4", "align": "stretch",
			}),
		},
		Stack,
	)
}

func flexDeliveryDefinition() DeliveryDefinition {
	componentType := "Flex"
	root := deliverySlot(
		componentType,
		"children",
		clFlex.Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Button", "primary", ""),
		deliveryInstance("Secondary", "Button", "secondary", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"One-dimensional flexible layout primitive.",
		map[string]PropertyContract{
			"direction": variantProperty([]string{"row", "column"}, "row", "Main-axis direction."),
			"wrap":      modifierProperty("Whether children may wrap."),
			"gap":       sizeProperty(gaps, "4", "Space between children."),
			"align":     variantProperty([]string{"start", "center", "end", "stretch"}, "center", "Cross-axis alignment."),
			"justify":   variantProperty([]string{"start", "center", "end", "between"}, "start", "Main-axis distribution."),
		},
		deliveryDesign("row", "05 Templates", "Layout", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "row", map[string]any{
				"direction": "row", "gap": "4", "align": "center", "justify": "start",
			}),
		},
		Flex,
	)
}

func gridDeliveryDefinition() DeliveryDefinition {
	componentType := "Grid"
	root := deliverySlot(
		componentType,
		"children",
		clGrid.GridCols(2).Gap(tw.S4).Compile(),
		deliveryInstance("Primary", "Card", "summary", ""),
		deliveryInstance("Secondary", "Card", "mobile-summary", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Two-dimensional responsive grid primitive.",
		map[string]PropertyContract{
			"columns": variantProperty([]string{"1", "2", "3", "4", "6", "12"}, "1", "Column count."),
			"gap":     sizeProperty(gaps, "4", "Space between cells."),
		},
		deliveryDesign("two-columns", "05 Templates", "Layout", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "two-columns", map[string]any{
				"columns": "2", "gap": "4",
			}),
			mobileDeliveryExample(componentType, "single-column", map[string]any{
				"columns": "1", "gap": "3",
			}),
		},
		Grid,
	)
}

func containerDeliveryDefinition() DeliveryDefinition {
	componentType := "Container"
	root := deliverySlot(
		componentType,
		"children",
		clContainer.MaxWScaled(tw.MaxW7XL).Compile(),
		deliveryInstance("Content", "Stack", "default", ""),
	)
	return newLayoutDeliveryDefinition(
		componentType,
		stableDeliveryID(componentType),
		"Centered maximum-width page container.",
		map[string]PropertyContract{
			"maxWidth": sizeProperty([]string{"sm", "md", "lg", "xl", "2xl", "4xl", "7xl", "full"}, "7xl", "Maximum content width."),
			"padding":  sizeProperty(gaps, "4", "Horizontal padding."),
		},
		deliveryDesign("default", "05 Templates", "Layout", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "default", map[string]any{
				"maxWidth": "7xl", "padding": "4",
			}),
		},
		Container,
	)
}

// DataManagementPageProps is the typed contract for the canonical OSS
// solution page. It deliberately models product-neutral data management so
// proprietary modules and clients can extend the composition without forking
// its atoms, molecules, organism, or layout primitives.
type DataManagementPageProps struct {
	Title       string               `json:"title"`
	Description string               `json:"description,omitempty"`
	Rows        []molecules.TableRow `json:"rows,omitempty"`
}

func dataManagementPageDeliveryDefinition() DeliveryDefinition {
	componentType := "DataManagementPage"
	root := deliveryFrame(
		componentType,
		clDataManagementPage.Compile(),
		deliveryFrame(
			"Container",
			clContainer.MaxWScaled(tw.MaxW7XL).Compile(),
			deliveryFrame(
				"Stack",
				clStack.Gap(tw.S6).Compile(),
				deliveryInstance("Breadcrumb", "Breadcrumb", "default", ""),
				deliveryFrame(
					"Header",
					clFlex.Gap(tw.S4).Items(tw.ItemsStart).Justify(tw.JustifyBetween).Compile(),
					deliveryFrame(
						"Copy",
						clStack.Gap(tw.S2).Compile(),
						deliveryInstance("Title", "Heading", "level-one", ""),
						deliveryInstance("Description", "Text", "muted", ""),
					),
					deliveryInstance("Create", "Button", "primary", ""),
				),
				deliveryInstance("Tabs", "Tabs", "underline", ""),
				deliveryInstance("Results", "DataGrid", "workspace-list", ""),
			),
		),
	)
	design := deliveryDesign(
		"workspace-management",
		"06 OSS Solutions",
		"Data Management",
		root,
	)
	design.Taxonomy.FlowGroup = "data-management"
	design.Taxonomy.FlowOrder = 10
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierPage,
		stableDeliveryID(componentType),
		"Complete OSS data-management solution composed exclusively from shared catalog building blocks.",
		map[string]PropertyContract{
			"title":       contentProperty("Page heading."),
			"description": contentProperty("Page purpose."),
			"rows":        contentProperty("Rows displayed by the shared data grid."),
		},
		nil,
		design,
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "workspace-management", map[string]any{
				"title":       "Workspaces",
				"description": "Manage access, lifecycle, and operational state.",
				"rows": []map[string]any{
					{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
					{"id": "bravo", "cells": map[string]any{"name": "Bravo workspace", "status": "Review", "updated": "Yesterday"}},
					{"id": "charlie", "cells": map[string]any{"name": "Charlie workspace", "status": "Paused", "updated": "3 days ago"}},
				},
			}),
			mobileDeliveryExample(componentType, "workspace-management-mobile", map[string]any{
				"title":       "Workspaces",
				"description": "Manage access and lifecycle.",
				"rows": []map[string]any{
					{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
				},
			}),
		},
		renderDataManagementPage,
	)
}

func renderDataManagementPage(props DataManagementPageProps) g.Node {
	rows := props.Rows
	if len(rows) == 0 {
		rows = []molecules.TableRow{{
			ID: "alpha",
			Cells: map[string]any{
				"name": "Alpha workspace", "status": "Active", "updated": "Today",
			},
		}}
	}
	title := props.Title
	if title == "" {
		title = "Workspaces"
	}
	description := props.Description
	if description == "" {
		description = "Manage access, lifecycle, and operational state."
	}
	gridProps := dataGridProps(rows)
	return h.Main(
		h.Class(clDataManagementPage.Compile()),
		Container(
			layouts.ContainerProps{MaxWidth: "7xl"},
			Stack(
				layouts.StackProps{Gap: "6"},
				Breadcrumb(molecules.BreadcrumbProps{
					Items: []molecules.BreadcrumbItem{
						{Label: "Operations", Href: "/operations"},
						{Label: title, Current: true},
					},
				}),
				Flex(
					layouts.FlexProps{Gap: "4", Align: "start", Justify: "between"},
					Stack(
						layouts.StackProps{Gap: "2"},
						Heading(atoms.HeadingProps{Text: title, Level: 1}),
						Text(atoms.TextProps{Content: description, Size: "sm", Color: "muted"}),
					),
					Button(atoms.ButtonProps{Text: "New workspace", Variant: "primary", Size: "medium"}),
				),
				Tabs(molecules.TabsProps{
					ActiveTab: "all",
					Items: []molecules.TabItem{
						{Key: "all", Label: "All", URL: "/workspaces"},
						{Key: "active", Label: "Active", URL: "/workspaces?status=active"},
						{Key: "archived", Label: "Archived", URL: "/workspaces?status=archived"},
					},
				}),
				DataGrid(gridProps),
			),
		),
	)
}

func tableExampleProps() map[string]any {
	return map[string]any{
		"columns": []map[string]any{
			{"key": "name", "label": "Name", "sortable": true, "primary": true},
			{"key": "status", "label": "Status", "sortable": true},
			{"key": "updated", "label": "Updated"},
		},
		"rows": []map[string]any{
			{"id": "alpha", "cells": map[string]any{"name": "Alpha workspace", "status": "Active", "updated": "Today"}},
			{"id": "bravo", "cells": map[string]any{"name": "Bravo workspace", "status": "Review", "updated": "Yesterday"}},
		},
		"sortable": true,
		"striped":  true,
	}
}

func dataGridExampleProps() map[string]any {
	return map[string]any{
		"table": tableExampleProps(),
		"pagination": map[string]any{
			"currentPage": 1, "totalPages": 4, "baseURL": "/workspaces",
		},
		"search": map[string]any{
			"searchURL": "/workspaces/search", "label": "Search workspaces",
			"placeholder": "Search workspaces", "instant": true,
		},
		"createURL":  "/workspaces/new",
		"createText": "New workspace",
	}
}

func dataGridProps(rows []molecules.TableRow) organisms.DataGridProps {
	return organisms.DataGridProps{
		Table: molecules.TableProps{
			Columns: []molecules.TableColumn{
				{Key: "name", Label: "Name", Sortable: true, Primary: true},
				{Key: "status", Label: "Status", Sortable: true},
				{Key: "updated", Label: "Updated"},
			},
			Rows:     rows,
			Sortable: true,
			Striped:  true,
		},
		Pagination: molecules.PaginationProps{
			CurrentPage: 1,
			TotalPages:  4,
			BaseURL:     "/workspaces",
		},
		Search: molecules.SearchBarProps{
			SearchURL:   "/workspaces/search",
			Label:       "Search workspaces",
			Placeholder: "Search workspaces",
			Instant:     true,
		},
		CreateURL:  "/workspaces/new",
		CreateText: "New workspace",
	}
}
