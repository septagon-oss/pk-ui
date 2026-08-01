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
	"github.com/septagon-oss/pk-ui/contracts"
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
			deliveryExample(componentType, "sortable", map[string]any{
				"columns": []map[string]any{
					{"key": "name", "label": "Name", "sortable": true, "primary": true},
					{"key": "status", "label": "Status", "sortable": true},
					{"key": "updated", "label": "Updated", "sortable": true},
				},
				"rows":     tableExampleProps()["rows"],
				"sortable": true,
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
	return newSlottedDeliveryDefinition(
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
			Attrs:        deliveryRepeatingSlotAttrs(),
		}},
		deliveryDesign("summary", "03 Molecules", "Surfaces", root),
		[]DeliveryExample{
			cardDeliveryExample(
				canonicalDeliveryExample(componentType, "summary", map[string]any{
					"title": "Quarterly summary", "description": "Performance is tracking above plan.",
				}),
				"Review the current operational details.",
			),
			cardDeliveryExample(
				mobileDeliveryExample(componentType, "mobile-summary", map[string]any{
					"title": "Open reviews", "description": "Three items need attention.",
				}),
				"Open the review queue.",
			),
			cardDeliveryExample(
				deliveryExample(componentType, "basic", map[string]any{
					"title": "Quarterly summary", "description": "Performance is tracking above plan.",
				}),
				"Review the current operational details.",
			),
		},
		func(props molecules.CardProps, slots DeliverySlotChildren) g.Node {
			return Card(props, slots.Nodes("children")...)
		},
	)
}

func cardDeliveryExample(
	example DeliveryExample,
	content string,
) DeliveryExample {
	return withDeliveryExampleSlots(
		example,
		deliveryExampleSlot(
			"children",
			deliveryExampleComponent(
				"body",
				"Text",
				map[string]any{
					"content": content,
					"size":    "sm",
					"color":   "muted",
				},
			),
		),
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
			deliveryExample(componentType, "with-chevrons", map[string]any{
				"separator": "›",
				"items": []map[string]any{
					{"label": "Home", "href": "/"},
					{"label": "Products", "href": "/products"},
					{"label": "Widget Pro", "current": true},
				},
			}),
		},
		Breadcrumb,
	)
}

func paginationDeliveryDefinition() DeliveryDefinition {
	componentType := "Pagination"
	numberedTotals := []any{4, 7}
	pageOne := deliveryClassBound(
		deliveryVisibleWhen(
			deliveryText("PageOne", clPageBtn.Merge(clPageIdle).Compile(), "1"),
			map[string]any{"totalPages": numberedTotals},
		),
		"currentPage",
		map[string]string{
			"1": clPageBtn.Merge(clPageCur).Compile(),
			"2": clPageBtn.Merge(clPageIdle).Compile(),
		},
	)
	pageTwo := deliveryClassBound(
		deliveryVisibleWhen(
			deliveryText("PageTwo", clPageBtn.Merge(clPageIdle).Compile(), "2"),
			map[string]any{"totalPages": numberedTotals},
		),
		"currentPage",
		map[string]string{
			"1": clPageBtn.Merge(clPageIdle).Compile(),
			"2": clPageBtn.Merge(clPageCur).Compile(),
		},
	)
	root := deliveryFrame(
		componentType,
		clPagination.Compile(),
		deliveryVisibleWhen(
			deliveryText("Previous", clPageBtn.Merge(clPageIdle).Compile(), "‹"),
			map[string]any{"currentPage": 2},
		),
		pageOne,
		pageTwo,
		deliveryVisibleWhen(
			deliveryText("PageThree", clPageBtn.Merge(clPageIdle).Compile(), "3"),
			map[string]any{"totalPages": 7},
		),
		deliveryVisibleWhen(
			deliveryText("Ellipsis", clBreadcrumbSep.Compile(), "…"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("LastPage", clPageBtn.Merge(clPageIdle).Compile(), "7", "totalPages"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("Next", clPageBtn.Merge(clPageIdle).Compile(), "›"),
			map[string]any{"totalPages": numberedTotals},
		),
		deliveryVisibleWhen(
			deliveryText("CursorPrevious", clPageBtn.Merge(clPageIdle).Compile(), "← Previous"),
			map[string]any{"totalPages": 0},
		),
		deliveryVisibleWhen(
			deliveryText("CursorLabel", clPageLabel.Compile(), "Page 1"),
			map[string]any{"totalPages": 0},
		),
		deliveryVisibleWhen(
			deliveryText("CursorNext", clPageBtn.Merge(clPageIdle).Compile(), "Next →"),
			map[string]any{"totalPages": 0},
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Accessible numbered or cursor pagination.",
		map[string]PropertyContract{
			"currentPage":     contentProperty("Current one-based page."),
			"totalPages":      contentProperty("Known total; zero selects cursor mode."),
			"perPage":         modifierProperty("Rows per page."),
			"siblings":        modifierProperty("Visible siblings around the current page."),
			"baseURL":         contentProperty("Numbered navigation base URL."),
			"cursorMode":      variantProperty([]string{"previous-next", "load-more"}, "previous-next", "Cursor navigation presentation."),
			"previousURL":     contentProperty("Opaque backward continuation URL."),
			"nextURL":         contentProperty("Opaque forward continuation URL."),
			"previousLabel":   contentProperty("Localized previous-window label."),
			"nextLabel":       contentProperty("Localized next-window label."),
			"loadMoreLabel":   contentProperty("Localized forward-only label."),
			"navigationLabel": contentProperty("Accessible navigation label."),
		},
		nil,
		deliveryDesign("numbered", "03 Molecules", "Navigation", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "numbered", map[string]any{
				"currentPage": 2, "totalPages": 7, "siblings": 1, "baseURL": "/workspaces",
			}),
			deliveryExample(componentType, "workspace-page-one", map[string]any{
				"currentPage": 1, "totalPages": 4, "baseURL": "/workspaces",
			}),
			deliveryExample(componentType, "cursor", map[string]any{
				"currentPage": 1, "totalPages": 0,
				"nextURL": "/workspaces?after=opaque",
			}),
		},
		Pagination,
	)
}

func searchBarDeliveryDefinition() DeliveryDefinition {
	componentType := "SearchBar"
	root := deliverySemantics(
		deliveryFrame(
			componentType,
			clSearchWrap.Compile(),
			deliveryInstance("Icon", "Icon", "Search / Fg Tertiary / 20", ""),
			deliveryText("Input", clSearchInput.Compile(), "Search workspaces", "placeholder", "label"),
			deliveryHiddenWhen(
				deliveryInstance("Clear", "Button", "button-icon-only-ghost", ""),
				map[string]any{"showClear": false},
			),
			deliveryHiddenWhen(
				deliveryInstance("Shortcut", "Kbd", "shortcut", ""),
				map[string]any{"showShortcut": false},
			),
		),
		"field",
		"field",
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
			"name":         contentProperty("Submitted query field name."),
			"value":        contentProperty("Initial query value."),
			"instant":      stateProperty("Whether typing submits progressively."),
			"showClear":    modifierProperty("Whether a clear affordance is shown."),
			"showShortcut": modifierProperty("Whether a keyboard hint is shown."),
			"debounceMs":   modifierProperty("Typing debounce in milliseconds."),
			"minChars":     modifierProperty("Minimum query length exposed to enhancement controllers."),
			"clearLabel":   contentProperty("Localized clear-action label."),
			"shortcutKey":  contentProperty("Displayed keyboard shortcut."),
		},
		nil,
		deliveryDesign("basic", "03 Molecules", "Forms", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "basic", map[string]any{
				"searchURL": "/workspaces/search", "label": "Search workspaces",
				"placeholder": "Search workspaces", "instant": true,
			}),
			deliveryExample(componentType, "with-shortcut", map[string]any{
				"searchURL": "/search", "label": "Search everything",
				"placeholder": "Search everything...", "instant": true, "showShortcut": true,
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
			deliverySlot(
				"Search",
				"search",
				"flex-1",
				deliveryInstance("SearchFixture", "SearchBar", "basic", ""),
			),
			deliverySlot("Filters", "filters", ""),
			deliverySlot(
				"Actions",
				"actions",
				clGridActions.Compile(),
				deliveryInstance("CreateFixture", "Link", "default", ""),
			),
		),
		deliverySlot(
			"Table",
			"table",
			"",
			deliveryInstance("TableFixture", "Table", "data", ""),
		),
		deliverySlot(
			"Status",
			"status",
			"",
			deliveryInstance("Empty", "EmptyState", "compact", ""),
		),
		deliverySlot(
			"Pagination",
			"pagination",
			"",
			deliveryInstance("PaginationFixture", "Pagination", "numbered", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Search, filter, action, table, status, and pagination composition.",
		map[string]PropertyContract{},
		[]designcomponent.Slot{
			{
				Name:         "search",
				Description:  "Optional search control.",
				AllowedTypes: []string{"SearchBar"},
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "filters",
				Description:  "Ordered filter controls.",
				AllowedTypes: []string{"Input", "Select"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "actions",
				Description:  "Toolbar actions.",
				AllowedTypes: []string{"Button", "Link"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "table",
				Description:  "Primary tabular result view.",
				AllowedTypes: []string{"Table"},
				Required:     true,
				Cardinality:  designcomponent.SlotOne,
			},
			{
				Name:         "status",
				Description:  "Loading, error, or empty status between data and pagination.",
				AllowedTypes: []string{"Alert", "EmptyState", "Spinner", "Text"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "pagination",
				Description:  "Optional pagination control.",
				AllowedTypes: []string{"Pagination"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("workspace-list", "04 Organisms", "Data", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "workspace-list", nil),
				dataGridDeliveryExampleSlots()...,
			),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "server-side-grid", map[string]any{
					"hx-get":     "/workspaces",
					"hx-target":  "#workspace-grid",
					"hx-trigger": "workspace:changed from:body",
				}),
				dataGridDeliveryExampleSlots()...,
			),
		},
		func(props organisms.DataGridProps, slots DeliverySlotChildren) g.Node {
			return DataGrid(props, DataGridSlots{
				Search:     firstDeliverySlotNode(slots, "search"),
				Filters:    slots.Nodes("filters"),
				Actions:    slots.Nodes("actions"),
				Table:      firstDeliverySlotNode(slots, "table"),
				Status:     slots.Nodes("status"),
				Pagination: firstDeliverySlotNode(slots, "pagination"),
			})
		},
	)
}

func dataGridDeliveryExampleSlots() []DeliveryExampleSlot {
	return []DeliveryExampleSlot{
		deliveryExampleSlot(
			"search",
			deliveryExampleComponent(
				"workspace-search",
				"SearchBar",
				map[string]any{
					"searchURL":   "/workspaces/search",
					"label":       "Search workspaces",
					"placeholder": "Search workspaces",
					"instant":     true,
				},
			),
		),
		deliveryExampleSlot(
			"actions",
			deliveryExampleComponent(
				"create-workspace",
				"Link",
				map[string]any{
					"label": "New workspace",
					"href":  "/workspaces/new",
				},
			),
		),
		deliveryExampleSlot(
			"table",
			deliveryExampleComponent(
				"workspace-table",
				"Table",
				tableExampleProps(),
			),
		),
		deliveryExampleSlot(
			"pagination",
			deliveryExampleComponent(
				"workspace-pagination",
				"Pagination",
				map[string]any{
					"currentPage": 1,
					"totalPages":  4,
					"baseURL":     "/workspaces",
				},
			),
		),
	}
}

func windowedCollectionDeliveryDefinition() DeliveryDefinition {
	componentType := "WindowedCollection"
	root := deliveryFrame(
		componentType,
		clWindow.Compile(),
		deliverySlot(
			"Items",
			"items",
			clWindowItems.Compile(),
			deliveryInstance("FirstItem", "Card", "summary", ""),
			deliveryInstance("SecondItem", "Card", "summary", ""),
		),
		deliveryVisibleWhen(
			deliveryText(
				"Loading",
				clWindowLoading.Compile(),
				"Loading results",
				"loadingLabel",
			),
			map[string]any{"state": "loading"},
		),
		deliveryVisibleWhen(
			deliveryText(
				"Error",
				clWindowError.Compile(),
				"Unable to load results",
				"errorTitle",
			),
			map[string]any{"state": "error"},
		),
		deliveryVisibleWhen(
			deliveryText(
				"Empty",
				clWindowEmpty.Compile(),
				"No results",
				"emptyTitle",
			),
			map[string]any{"itemCount": 0},
		),
		deliverySlot(
			"Controls",
			"controls",
			clWindowFooter.Compile(),
			deliveryInstance("Pagination", "Pagination", "cursor", ""),
		),
	)
	return newSlottedDeliveryDefinition(
		componentType,
		uicomponent.TierOrganism,
		stableDeliveryID(componentType),
		"Bounded collection shell with explicit item, state, and cursor-control composition.",
		map[string]PropertyContract{
			"state": variantProperty(
				[]string{"ready", "loading", "error"},
				"ready",
				"Current collection state.",
			),
			"itemCount":             contentProperty("Number of bounded item roots."),
			"maxItems":              modifierProperty("Hard render-node budget."),
			"collectionLabel":       contentProperty("Accessible collection label."),
			"emptyTitle":            contentProperty("Localized empty-state title."),
			"emptyDescription":      contentProperty("Localized empty-state description."),
			"loadingLabel":          contentProperty("Localized loading status."),
			"errorTitle":            contentProperty("Localized error title."),
			"errorDescription":      contentProperty("Localized error description."),
			"retryLabel":            contentProperty("Localized retry label."),
			"retryURL":              contentProperty("Progressive retry URL."),
			"navigationUnavailable": contentProperty("Cursor-navigation failure message."),
		},
		[]designcomponent.Slot{
			{
				Name:         "items",
				Description:  "Bounded item roots in stable source order.",
				AllowedTypes: []string{"Card", "DataGrid", "Table", "Text"},
				Cardinality:  designcomponent.SlotMany,
				Attrs:        deliveryRepeatingSlotAttrs(),
			},
			{
				Name:         "controls",
				Description:  "Optional bounded cursor navigation.",
				AllowedTypes: []string{"Pagination"},
				Cardinality:  designcomponent.SlotOne,
			},
		},
		deliveryDesign("cards", "04 Organisms", "Data", root),
		[]DeliveryExample{
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "cards", map[string]any{
					"state": "ready", "itemCount": 2, "maxItems": 100,
				}),
				deliveryExampleSlot(
					"items",
					deliveryExampleComponent(
						"window-card-one",
						"Card",
						map[string]any{"title": "First result"},
					),
					deliveryExampleComponent(
						"window-card-two",
						"Card",
						map[string]any{"title": "Second result"},
					),
				),
				deliveryExampleSlot(
					"controls",
					deliveryExampleComponent(
						"window-pagination",
						"Pagination",
						map[string]any{
							"currentPage": 1,
							"totalPages":  0,
							"nextURL":     "/workspaces?after=opaque",
						},
					),
				),
			),
			deliveryExample(componentType, "empty", map[string]any{
				"state": "ready", "itemCount": 0, "maxItems": 100,
			}),
		},
		func(
			props organisms.WindowedCollectionProps,
			slots DeliverySlotChildren,
		) g.Node {
			return WindowedCollection(props, WindowedCollectionSlots{
				Items:    slots.Nodes("items"),
				Controls: firstDeliverySlotNode(slots, "controls"),
			})
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
	root = deliveryClassBound(root, "gap", deliveryGapClassMap())
	root = deliveryClassBound(root, "align", deliveryAlignClassMap())
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
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"gap": "4", "align": "stretch",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"primary-card",
						"Quarterly summary",
						"Performance is tracking above plan.",
					),
					deliveryExampleCard(
						"secondary-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
			withDeliveryExampleSlots(
				deliveryExample(componentType, "compact", map[string]any{
					"gap": "2", "align": "stretch",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"compact-heading",
						"Heading",
						map[string]any{"level": 2, "text": "Compact stack"},
					),
					deliveryExampleComponent(
						"compact-copy",
						"Text",
						map[string]any{
							"color":   "muted",
							"content": "Closely related supporting content.",
							"size":    "base",
						},
					),
				),
			),
		},
		Stack,
	)
}

func deliveryGapClassMap() map[string]string {
	out := make(map[string]string, len(clGapScale))
	for name, gap := range clGapScale {
		out[name] = tw.New().Gap(gap).Compile()
	}
	return out
}

func deliveryAlignClassMap() map[string]string {
	out := make(map[string]string, len(alignItems))
	for name, align := range alignItems {
		out[name] = tw.New().Items(align).Compile()
	}
	return out
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
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "row", map[string]any{
					"direction": "row", "gap": "4", "align": "center", "justify": "start",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"primary-action",
						"Button",
						map[string]any{
							"label": "Continue", "variant": "primary",
							"tone": "neutral", "size": "md",
						},
					),
					deliveryExampleComponent(
						"secondary-action",
						"Button",
						map[string]any{
							"label": "Cancel", "variant": "secondary",
							"tone": "neutral", "size": "md",
						},
					),
				),
			),
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
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "two-columns", map[string]any{
					"columns": "2", "gap": "4",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"primary-card",
						"Quarterly summary",
						"Performance is tracking above plan.",
					),
					deliveryExampleCard(
						"secondary-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
			withDeliveryExampleSlots(
				mobileDeliveryExample(componentType, "single-column", map[string]any{
					"columns": "1", "gap": "3",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleCard(
						"mobile-card",
						"Open reviews",
						"Three items need attention.",
					),
				),
			),
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
			withDeliveryExampleSlots(
				canonicalDeliveryExample(componentType, "default", map[string]any{
					"maxWidth": "7xl", "padding": "4",
				}),
				deliveryExampleSlot(
					"children",
					deliveryExampleComponent(
						"content-stack",
						"Stack",
						map[string]any{"gap": "2", "align": "stretch"},
						deliveryExampleSlot(
							"children",
							deliveryExampleComponent(
								"content-heading",
								"Heading",
								map[string]any{
									"text":  "Workspace overview",
									"level": 2,
								},
							),
							deliveryExampleComponent(
								"content-copy",
								"Text",
								map[string]any{
									"content": "Review current operational details.",
									"size":    "base", "color": "muted",
								},
							),
						),
					),
				),
			),
		},
		Container,
	)
}

func deliveryExampleCard(
	id string,
	title string,
	description string,
) DeliveryExampleComponent {
	return deliveryExampleComponent(
		id,
		"Card",
		map[string]any{
			"title":       title,
			"description": description,
		},
		deliveryExampleSlot(
			"children",
			deliveryExampleComponent(
				id+"-body",
				"Text",
				map[string]any{
					"content": description,
					"size":    "sm",
					"color":   "muted",
				},
			),
		),
	)
}

// DataManagementPageProps is the typed contract for the canonical OSS
// solution page. It deliberately models product-neutral data management so
// proprietary modules and clients can extend the composition without forking
// its atoms, molecules, organism, or layout primitives.
type DataManagementPageProps struct {
	contracts.ComponentProps

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
					clStack.Gap(tw.S2).Compile(),
					deliveryInstance("Title", "Heading", "level-one", ""),
					deliveryInstance("Description", "Text", "muted", ""),
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
	main := baseAttrs(props.ComponentProps)
	main = append(main,
		classes(clDataManagementPage.Compile(), props.Class),
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
				Stack(
					layouts.StackProps{Gap: "2"},
					Heading(atoms.HeadingProps{Text: title, Level: 1}),
					Text(atoms.TextProps{Content: description, Size: "sm", Color: "muted"}),
				),
				Tabs(molecules.TabsProps{
					ActiveTab: "all",
					Items: []molecules.TabItem{
						{Key: "all", Label: "All", URL: "/workspaces"},
						{Key: "active", Label: "Active", URL: "/workspaces?status=active"},
						{Key: "archived", Label: "Archived", URL: "/workspaces?status=archived"},
					},
				}),
				DataGrid(
					organisms.DataGridProps{},
					dataGridSlots(rows),
				),
			),
		),
	)
	return h.Main(main...)
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

func dataGridSlots(rows []molecules.TableRow) DataGridSlots {
	return DataGridSlots{
		Search: SearchBar(molecules.SearchBarProps{
			SearchURL:   "/workspaces/search",
			Label:       "Search workspaces",
			Placeholder: "Search workspaces",
			Instant:     true,
		}),
		Actions: []g.Node{
			Link(atoms.LinkProps{
				Label: "New workspace",
				Href:  "/workspaces/new",
			}),
		},
		Table: Table(molecules.TableProps{
			Columns: []molecules.TableColumn{
				{Key: "name", Label: "Name", Sortable: true, Primary: true},
				{Key: "status", Label: "Status", Sortable: true},
				{Key: "updated", Label: "Updated"},
			},
			Rows:     rows,
			Sortable: true,
			Striped:  true,
		}),
		Pagination: Pagination(molecules.PaginationProps{
			CurrentPage: 1,
			TotalPages:  4,
			BaseURL:     "/workspaces",
		}),
	}
}

func tableSkeletonDeliveryDefinition() DeliveryDefinition {
	componentType := "TableSkeleton"
	line := skeletonLine("sm", false).Compile()
	root := deliveryFrame(
		componentType,
		clTableWrap.Compile(),
		deliveryFrame(
			"Table",
			clTable.Compile(),
			deliveryFrame(
				"Header",
				clTableHead.Compile(),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Column", clTableTh.Compile(), deliveryFrame("Line", line)),
			),
			deliveryFrame(
				"Row",
				clTableRow.Compile(),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
				deliveryFrame("Cell", clTableTd.Compile(), deliveryFrame("Line", line)),
			),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Loading rendering of Table: identical wrap, header, and cell classes with pulsing placeholders.",
		map[string]PropertyContract{
			"columns": contentProperty("Placeholder column count."),
			"rows":    contentProperty("Placeholder row count."),
			"compact": modifierProperty("Whether row density is compact."),
		},
		nil,
		deliveryDesign("loading", "03 Molecules", "Data", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "loading", map[string]any{
				"columns": 3, "rows": 3,
			}),
		},
		TableSkeleton,
	)
}

func cardSkeletonDeliveryDefinition() DeliveryDefinition {
	componentType := "CardSkeleton"
	root := deliveryFrame(
		componentType,
		clCard.Compile(),
		deliveryFrame("Title", skeletonLine("lg", true).Compile()),
		deliveryFrame(
			"Body",
			clSkeletonText.Compile(),
			deliveryFrame("Line", skeletonLine("md", false).Compile()),
			deliveryFrame("Line", skeletonLine("md", false).Compile()),
			deliveryFrame("Line", skeletonLine("md", true).Compile()),
		),
	)
	return newDeliveryDefinition(
		componentType,
		uicomponent.TierMolecule,
		stableDeliveryID(componentType),
		"Loading rendering of Card: title and body placeholders inside the real card frame.",
		map[string]PropertyContract{
			"lines": contentProperty("Body placeholder line count."),
		},
		nil,
		deliveryDesign("loading", "03 Molecules", "Surfaces", root),
		[]DeliveryExample{
			canonicalDeliveryExample(componentType, "loading", map[string]any{
				"lines": 3,
			}),
		},
		CardSkeleton,
	)
}
