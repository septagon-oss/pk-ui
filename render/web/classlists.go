// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// classlists.go declares every tw ClassList the renderers compose. This is
// the load-bearing inversion of the usual CSS workflow: the stylesheet is
// derived FROM these declarations via tw/emission.For(ClassLists()...), so a
// class that no component declares is never emitted, and a declared class is
// always backed by a rule. TestRenderedClassesAreDeclared closes the loop
// from the other side.

import (
	"github.com/septagon-oss/tw"
)

// hoverBg returns a hover-state background modifier.
func hoverBg(c tw.Color) func(tw.ClassList) tw.ClassList {
	return func(cl tw.ClassList) tw.ClassList { return cl.Bg(c) }
}

var (
	// Shared fragments.
	clIcon = tw.New().Display(tw.DisplayInlineBlock).Width(tw.S4).Height(tw.S4).FlexShrink0()

	clFocusRing = tw.New().
			On(tw.StateFocusVisible, func(c tw.ClassList) tw.ClassList {
			return c.Ring(tw.Ring2).RingColor(tw.RingFocus).RingOffset(tw.RingOffset2)
		})

	// Button: base plus one list per variant.
	// Variant discipline, applied to every base+variant pair in this file: a
	// base fragment never declares a property any of its variants declares.
	// Two single-class utilities on one element have equal specificity, so
	// stylesheet order — not merge order — would pick the winner.
	// TestComposedListsHaveNoPropertyCollisions enforces this structurally.
	clButtonBase = tw.New().
			Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
			Gap(tw.S2).FontWeight(tw.FontSemibold).
			Rounded(tw.RadiusMD).Border(tw.Border1).
			// Buttons render as <button> or as <a> (EmptyState actions,
			// link-shaped calls to action); an anchor must not inherit the
			// page's link underline.
			NoUnderline().
			Cursor(tw.CursorPointer).
			Transition(tw.TransitionColors).
			On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
		}).
		Merge(clFocusRing)

	clButtonVariant = map[string]tw.ClassList{
		"primary": tw.New().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, hoverBg(tw.SurfaceBrandHover)),
		"secondary": tw.New().Bg(tw.SurfacePrimary).TextColor(tw.FgPrimary).
			BorderColor(tw.BorderPrimary).
			On(tw.StateHover, hoverBg(tw.SurfaceHover)),
		"success": tw.New().Bg(tw.SurfaceSuccess).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"warning": tw.New().Bg(tw.SurfaceWarning).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"error":   tw.New().Bg(tw.SurfaceDanger).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"danger":  tw.New().Bg(tw.SurfaceDanger).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"info":    tw.New().Bg(tw.SurfaceInfo).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"outline": tw.New().Bg(tw.ColorTransparent).TextColor(tw.FgBrand).
			BorderColor(tw.BorderBrand).
			On(tw.StateHover, hoverBg(tw.SurfaceBrandSoft)),
		"ghost": tw.New().Bg(tw.ColorTransparent).TextColor(tw.FgPrimary).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, hoverBg(tw.SurfaceHover)),
		"link": tw.New().Bg(tw.ColorTransparent).TextColor(tw.FgLink).BorderColor(tw.ColorTransparent).Underline().
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgLinkHover) }),
	}

	clButtonSize = map[string]tw.ClassList{
		"xs":     tw.New().FontSize(tw.TextXS).PaddingX(tw.S2).PaddingY(tw.S1),
		"small":  tw.New().FontSize(tw.TextSM).PaddingX(tw.S3).PaddingY(tw.S1_5),
		"medium": tw.New().FontSize(tw.TextSM).PaddingX(tw.S4).PaddingY(tw.S2),
		"large":  tw.New().FontSize(tw.TextBase).PaddingX(tw.S5).PaddingY(tw.S2_5),
		"xl":     tw.New().FontSize(tw.TextLG).PaddingX(tw.S6).PaddingY(tw.S3),
		"2xl":    tw.New().FontSize(tw.TextXL).PaddingX(tw.S8).PaddingY(tw.S4),
	}

	clButtonFull = tw.New().Width(tw.SFull)

	// Badge / Tag.
	clBadgeBase = tw.New().
			Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S1).
			Rounded(tw.RadiusFull).FontWeight(tw.FontMedium).
			PaddingX(tw.S2_5).PaddingY(tw.S0_5).FontSize(tw.TextXS)

	clBadgeVariant = map[string]tw.ClassList{
		"default":   tw.New().Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary),
		"primary":   tw.New().Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand),
		"secondary": tw.New().Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary),
		"success":   tw.New().Bg(tw.SurfaceSuccessSoft).TextColor(tw.FgSuccess),
		"warning":   tw.New().Bg(tw.SurfaceWarningSoft).TextColor(tw.FgWarning),
		"error":     tw.New().Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger),
		"danger":    tw.New().Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger),
		"info":      tw.New().Bg(tw.SurfaceInfoSoft).TextColor(tw.FgInfo),
	}

	clBadgeDot = tw.New().Width(tw.S1_5).Height(tw.S1_5).Rounded(tw.RadiusFull).Bg(tw.FgBrand)

	// Alert.
	clAlertBase = tw.New().
			Display(tw.DisplayFlex).Gap(tw.S3).Rounded(tw.RadiusLG).Padding(tw.S4).
			Border(tw.Border1)

	clAlertVariant = map[string]tw.ClassList{
		"success": tw.New().Bg(tw.SurfaceSuccessSoft).TextColor(tw.FgSuccess).BorderColor(tw.BorderSuccess),
		"warning": tw.New().Bg(tw.SurfaceWarningSoft).TextColor(tw.FgWarning).BorderColor(tw.BorderWarning),
		"error":   tw.New().Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger).BorderColor(tw.BorderDanger),
		"danger":  tw.New().Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger).BorderColor(tw.BorderDanger),
		"info":    tw.New().Bg(tw.SurfaceInfoSoft).TextColor(tw.FgInfo).BorderColor(tw.BorderInfo),
	}

	clAlertTitle   = tw.New().FontWeight(tw.FontSemibold).FontSize(tw.TextSM)
	clAlertMessage = tw.New().FontSize(tw.TextSM)
	clAlertBody    = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1).Flex1()

	// Inputs.
	clFieldWrap = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1_5)
	clLabel     = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).TextColor(tw.FgPrimary)
	clHelp      = tw.New().FontSize(tw.TextXS).TextColor(tw.FgMuted)
	clFieldErr  = tw.New().FontSize(tw.TextXS).TextColor(tw.FgDanger)
	clRequired  = tw.New().TextColor(tw.FgDanger)

	clInput = tw.New().
		Display(tw.DisplayBlock).Width(tw.SFull).
		Rounded(tw.RadiusMD).Border(tw.Border1).
		Bg(tw.SurfacePrimary).TextColor(tw.FgPrimary).
		PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM).
		On(tw.StatePlaceholder, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPlaceholder) }).
		On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceDisabled).Cursor(tw.CursorNotAllowed)
		}).
		Merge(clFocusRing)

	clInputNormal = tw.New().BorderColor(tw.BorderPrimary)
	clInputError  = tw.New().BorderColor(tw.BorderDanger)

	clCheckbox = tw.New().
			Width(tw.S4).Height(tw.S4).Rounded(tw.RadiusSM).
			Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Cursor(tw.CursorPointer).Merge(clFocusRing)

	clCheckRow = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S2)

	// Search bar: a bordered wrapper that reads as one field, hosting a
	// borderless input so the glyph and the text share the control.
	clSearchWrap = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S2).
			Rounded(tw.RadiusMD).Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Bg(tw.SurfacePrimary).PaddingX(tw.S3).Flex1().MaxWScaled(tw.MaxWMD).
			TextColor(tw.FgMuted)
	clSearchInput = tw.New().Width(tw.SFull).Border(tw.Border0).Bg(tw.ColorTransparent).
			TextColor(tw.FgPrimary).PaddingY(tw.S2).FontSize(tw.TextSM).
			On(tw.StatePlaceholder, func(cl tw.ClassList) tw.ClassList { return cl.TextColor(tw.FgPlaceholder) })
	clSrOnly = tw.New().SrOnly()

	// Text and headings.
	clTextColor = map[string]tw.Color{
		"primary": tw.FgPrimary, "secondary": tw.FgSecondary, "muted": tw.FgMuted,
		"brand": tw.FgBrand, "success": tw.FgSuccess, "warning": tw.FgWarning,
		"danger": tw.FgDanger, "error": tw.FgDanger, "info": tw.FgInfo,
	}
	clTextSize = map[string]tw.FontSize{
		"xs": tw.TextXS, "sm": tw.TextSM, "base": tw.TextBase, "md": tw.TextBase,
		"lg": tw.TextLG, "xl": tw.TextXL, "2xl": tw.Text2XL,
	}
	clTextWeight = map[string]tw.FontWeight{
		"normal": tw.FontNormal, "medium": tw.FontMedium,
		"semibold": tw.FontSemibold, "bold": tw.FontBold,
	}
	clTruncate = tw.New().Truncate()

	clHeadingBase  = tw.New().FontFamily(tw.FontSerif).TextColor(tw.FgPrimary).FontWeight(tw.FontSemibold)
	clHeadingLevel = map[int]tw.ClassList{
		1: tw.New().FontSize(tw.Text3XL),
		2: tw.New().FontSize(tw.Text2XL),
		3: tw.New().FontSize(tw.TextXL),
		4: tw.New().FontSize(tw.TextLG),
		5: tw.New().FontSize(tw.TextBase),
		6: tw.New().FontSize(tw.TextSM).Uppercase().Tracking(tw.TrackingWider),
	}

	// Structure.
	clDividerH = tw.New().Border(tw.Border0).BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).MarginY(tw.S4)
	clDividerV = tw.New().Border(tw.Border0).BorderLeft(tw.Border1).BorderColor(tw.BorderPrimary).MarginX(tw.S4)

	clSpinner = tw.New().
			Display(tw.DisplayInlineBlock).Rounded(tw.RadiusFull).
			Border(tw.Border2).BorderColor(tw.BorderSecondary).
			BorderTopColor(tw.FgBrand).AnimateSpin()
	clSpinnerSize = map[string]tw.ClassList{
		"small":  tw.New().Width(tw.S4).Height(tw.S4),
		"medium": tw.New().Width(tw.S6).Height(tw.S6),
		"large":  tw.New().Width(tw.S8).Height(tw.S8),
	}

	clEmpty = tw.New().
		Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Items(tw.ItemsCenter).
		Justify(tw.JustifyCenter).Gap(tw.S3).TextAlign(tw.TextCenter)
	clEmptyPad      = tw.New().Padding(tw.S12)
	clEmptyBordered = tw.New().Border(tw.Border1).BorderColor(tw.BorderPrimary).
			BorderStyle(tw.BorderStyle("dashed")).Rounded(tw.RadiusLG)
	clEmptyCompact = tw.New().Padding(tw.S6)
	clEmptyTitle   = tw.New().FontFamily(tw.FontSerif).FontSize(tw.TextLG).
			FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clEmptyDesc = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted).MaxWScaled(tw.MaxWMD)

	clKbd = tw.New().
		Display(tw.DisplayInlineBlock).FontFamily(tw.FontMono).FontSize(tw.TextXS).
		Rounded(tw.RadiusSM).Border(tw.Border1).BorderColor(tw.BorderPrimary).
		BorderBottom(tw.Border2).Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary).
		PaddingX(tw.S1_5).PaddingY(tw.S0_5)

	clLink = tw.New().TextColor(tw.FgLink).Underline().UnderlineOffset(tw.S2).
		On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgLinkHover) }).
		Merge(clFocusRing)

	clTagBase = tw.New().
			Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S1).
			Rounded(tw.RadiusMD).Border(tw.Border1).
			PaddingX(tw.S2).PaddingY(tw.S0_5).FontSize(tw.TextXS)
	clTagIdle     = tw.New().BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).TextColor(tw.FgSecondary)
	clTagSelected = tw.New().BorderColor(tw.BorderBrand).Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand)

	// Layouts.
	clStack     = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol)
	clFlex      = tw.New().Display(tw.DisplayFlex)
	clGrid      = tw.New().Display(tw.DisplayGrid)
	clContainer = tw.New().MarginX(tw.SAuto).Width(tw.SFull).PaddingX(tw.S4)

	clGapScale = map[string]tw.Spacing{
		"0": tw.S0, "1": tw.S1, "2": tw.S2, "3": tw.S3, "4": tw.S4,
		"5": tw.S5, "6": tw.S6, "8": tw.S8,
	}

	// Table.
	clTableWrap = tw.New().Width(tw.SFull).Overflow(tw.OverflowAuto).
			Rounded(tw.RadiusLG).Border(tw.Border1).BorderColor(tw.BorderPrimary)
	clTable       = tw.New().Width(tw.SFull).FontSize(tw.TextSM).TextColor(tw.FgPrimary)
	clTableHead   = tw.New().Bg(tw.SurfaceSecondary).TextAlign(tw.TextLeft)
	clTableThBase = tw.New().FontWeight(tw.FontSemibold).
			FontSize(tw.TextXS).Uppercase().Tracking(tw.TrackingWider).TextColor(tw.FgMuted)
	clTableTh  = clTableThBase.Merge(tw.New().PaddingX(tw.S4).PaddingY(tw.S3))
	clTableTd  = tw.New().PaddingX(tw.S4).PaddingY(tw.S3).BorderTop(tw.Border1).BorderColor(tw.BorderPrimary)
	clTableRow = tw.New().On(tw.StateHover, hoverBg(tw.SurfaceHover))
	clTableTdC = tw.New().PaddingX(tw.S4).PaddingY(tw.S2)
	// Sortable headers render a real button so keyboards and readers get the
	// affordance; aria-sort lives on the th, per WAI-ARIA sortable-table
	// practice. The cell's padding moves onto the button to keep the whole
	// header surface clickable.
	clTableThSort  = tw.New().PaddingX(tw.S0).PaddingY(tw.S0)
	clTableSortBtn = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S1).
			Width(tw.SFull).PaddingX(tw.S4).PaddingY(tw.S3).
			FontWeight(tw.FontSemibold).FontSize(tw.TextXS).Uppercase().
			Tracking(tw.TrackingWider).TextColor(tw.FgMuted).
			Bg(tw.ColorTransparent).Border(tw.Border0).Cursor(tw.CursorPointer).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPrimary) }).
			Merge(clFocusRing)
	clTableRowAlt   = tw.New().Bg(tw.SurfaceSecondary)
	clTableTdStrong = tw.New().FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clTableActions  = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).Justify(tw.JustifyEnd)
	clTableCellNote = tw.New().TextColor(tw.FgMuted).FontSize(tw.TextXS)

	// Card.
	clCard = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S2).
		Rounded(tw.RadiusLG).Border(tw.Border1).BorderColor(tw.BorderPrimary).
		Bg(tw.SurfacePrimary).Padding(tw.S5).Shadow(tw.ShadowSM)
	clCardClickable = tw.New().Cursor(tw.CursorPointer).Transition(tw.TransitionShadow).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Shadow(tw.ShadowMD) }).
			Merge(clFocusRing)
	clCardTitle = tw.New().FontFamily(tw.FontSerif).FontSize(tw.TextLG).
			FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clCardDesc = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted)

	// Breadcrumb.
	clBreadcrumb = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).FontSize(tw.TextSM).
			ListStyle("none").Margin(tw.S0).Padding(tw.S0)
	clBreadcrumbSep = tw.New().TextColor(tw.FgTertiary)
	clBreadcrumbCur = tw.New().TextColor(tw.FgPrimary).FontWeight(tw.FontMedium)

	// Pagination.
	clPagination = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S1)
	clPageBtn    = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).MinWidth(tw.S8).Height(tw.S8).
			Rounded(tw.RadiusMD).FontSize(tw.TextSM).Merge(clFocusRing)
	clPageIdle = tw.New().TextColor(tw.FgSecondary).
			On(tw.StateHover, hoverBg(tw.SurfaceHover))
	clPageCur   = tw.New().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand).FontWeight(tw.FontSemibold)
	clPageLabel = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted).PaddingX(tw.S2)

	// Tabs.
	clTabList = tw.New().Display(tw.DisplayFlex).Gap(tw.S1).
			BorderBottom(tw.Border1).BorderColor(tw.BorderPrimary)
	clTab = tw.New().PaddingX(tw.S4).PaddingY(tw.S2).FontSize(tw.TextSM).
		FontWeight(tw.FontMedium).
		BorderBottom(tw.Border2).
		Merge(clFocusRing)
	clTabIdle = tw.New().TextColor(tw.FgMuted).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPrimary) })
	clTabActive = tw.New().TextColor(tw.FgBrand).BorderColor(tw.BorderBrand)
)

// ClassLists returns every ClassList the renderers compose, base and variant
// alike. Applications derive their stylesheet from it:
//
//	sheet, err := emission.For(web.ClassLists()...)
func ClassLists() []tw.ClassList {
	out := []tw.ClassList{
		clIcon, clFocusRing, clButtonBase, clButtonFull,
		clBadgeBase, clBadgeDot,
		clAlertBase, clAlertTitle, clAlertMessage, clAlertBody,
		clFieldWrap, clLabel, clHelp, clFieldErr, clRequired,
		clInput, clInputNormal, clInputError, clCheckbox, clCheckRow,
		clSearchWrap, clSearchInput, clSrOnly,
		clTruncate, clHeadingBase,
		clDividerH, clDividerV,
		clSpinner, clEmpty, clEmptyPad, clEmptyBordered, clEmptyCompact, clEmptyTitle, clEmptyDesc,
		clKbd, clLink, clTagBase, clTagIdle, clTagSelected,
		clStack, clFlex, clGrid, clContainer,
		clTableWrap, clTable, clTableHead, clTableThBase, clTableTh, clTableTd, clTableRow, clTableTdC,
		clTableThSort, clTableSortBtn, clTableRowAlt, clTableTdStrong, clTableActions, clTableCellNote,
		clCard, clCardClickable, clCardTitle, clCardDesc,
		clBreadcrumb, clBreadcrumbSep, clBreadcrumbCur,
		clPagination, clPageBtn, clPageIdle, clPageCur, clPageLabel,
		clTabList, clTab, clTabIdle, clTabActive,
	}
	for _, m := range []map[string]tw.ClassList{
		clButtonVariant, clButtonSize, clBadgeVariant, clAlertVariant,
		clSpinnerSize,
	} {
		for _, cl := range m {
			out = append(out, cl)
		}
	}
	for _, cl := range clHeadingLevel {
		out = append(out, cl)
	}
	// Enumerated one-offs composed inline by renderers.
	for _, sz := range clTextSize {
		out = append(out, tw.New().FontSize(sz))
	}
	for _, w := range clTextWeight {
		out = append(out, tw.New().FontWeight(w))
	}
	for _, c := range clTextColor {
		out = append(out, tw.New().TextColor(c))
	}
	for _, s := range clGapScale {
		out = append(out, tw.New().Gap(s))
	}
	for _, n := range []int{1, 2, 3, 4, 6, 12} {
		out = append(out, tw.New().GridCols(n))
	}
	out = append(out,
		tw.New().Items(tw.ItemsStart), tw.New().Items(tw.ItemsCenter),
		tw.New().Items(tw.ItemsEnd), tw.New().Items(tw.ItemsStretch),
		tw.New().Justify(tw.JustifyStart), tw.New().Justify(tw.JustifyCenter),
		tw.New().Justify(tw.JustifyEnd), tw.New().Justify(tw.JustifyBetween),
		tw.New().FlexDir(tw.FlexRow), tw.New().FlexDir(tw.FlexCol), tw.New().FlexWrap(),
		tw.New().MaxWScaled(tw.MaxWSM), tw.New().MaxWScaled(tw.MaxWMD),
		tw.New().MaxWScaled(tw.MaxWLG), tw.New().MaxWScaled(tw.MaxWXL),
		tw.New().MaxWScaled(tw.MaxW2XL), tw.New().MaxWScaled(tw.MaxW4XL),
		tw.New().MaxWScaled(tw.MaxW7XL), tw.New().MaxWScaled(tw.MaxWFull),
	)
	return out
}
