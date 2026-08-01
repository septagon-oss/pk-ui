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
	clIcon     = tw.New().Display(tw.DisplayInlineBlock).FlexShrink0()
	clIconSize = map[string]tw.ClassList{
		"xs":  tw.New().Width(tw.S3).Height(tw.S3),
		"sm":  tw.New().Width(tw.S4).Height(tw.S4),
		"md":  tw.New().Width(tw.S5).Height(tw.S5),
		"lg":  tw.New().Width(tw.S6).Height(tw.S6),
		"xl":  tw.New().Width(tw.S8).Height(tw.S8),
		"2xl": tw.New().Width(tw.S12).Height(tw.S12),
	}
	clIconTone = map[string]tw.ClassList{
		"neutral": tw.New(),
		"brand":   tw.New().TextColor(tw.FgBrand),
		"success": tw.New().TextColor(tw.FgSuccess),
		"warning": tw.New().TextColor(tw.FgWarning),
		"danger":  tw.New().TextColor(tw.FgDanger),
		"info":    tw.New().TextColor(tw.FgInfo),
	}

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
			Cursor(tw.CursorPointer).
			Transition(tw.TransitionColors).
			On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
		}).
		Merge(clFocusRing)

	clButtonVariant = map[string]tw.ClassList{
		"primary": tw.New().NoUnderline().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, hoverBg(tw.SurfaceBrandHover)),
		"secondary": tw.New().NoUnderline().Bg(tw.SurfacePrimary).TextColor(tw.FgPrimary).
			BorderColor(tw.BorderPrimary).
			On(tw.StateHover, hoverBg(tw.SurfaceHover)),
		"outline": tw.New().NoUnderline().Bg(tw.ColorTransparent).TextColor(tw.FgBrand).
			BorderColor(tw.BorderBrand).
			On(tw.StateHover, hoverBg(tw.SurfaceBrandSoft)),
		"ghost": tw.New().NoUnderline().Bg(tw.ColorTransparent).TextColor(tw.FgPrimary).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, hoverBg(tw.SurfaceHover)),
		"link": tw.New().Bg(tw.ColorTransparent).TextColor(tw.FgLink).BorderColor(tw.ColorTransparent).Underline().
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgLinkHover) }),
	}
	clButtonTone = map[string]tw.ClassList{
		"neutral": tw.New(),
		"brand": tw.New().NoUnderline().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent).
			On(tw.StateHover, hoverBg(tw.SurfaceBrandHover)),
		"success": tw.New().NoUnderline().Bg(tw.SurfaceSuccess).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"warning": tw.New().NoUnderline().Bg(tw.SurfaceWarning).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"danger":  tw.New().Bg(tw.SurfaceDanger).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
		"info":    tw.New().Bg(tw.SurfaceInfo).TextColor(tw.FgOnBrand).BorderColor(tw.ColorTransparent),
	}

	clButtonSize = map[string]tw.ClassList{
		"xs":  tw.New().FontSize(tw.TextXS).PaddingX(tw.S2).PaddingY(tw.S1),
		"sm":  tw.New().FontSize(tw.TextSM).PaddingX(tw.S3).PaddingY(tw.S1_5),
		"md":  tw.New().FontSize(tw.TextSM).PaddingX(tw.S4).PaddingY(tw.S2),
		"lg":  tw.New().FontSize(tw.TextBase).PaddingX(tw.S5).PaddingY(tw.S2_5),
		"xl":  tw.New().FontSize(tw.TextLG).PaddingX(tw.S6).PaddingY(tw.S3),
		"2xl": tw.New().FontSize(tw.TextXL).PaddingX(tw.S8).PaddingY(tw.S4),
	}

	clButtonFull     = tw.New().Width(tw.SFull)
	clButtonIconOnly = tw.New().MinWidth(tw.S11).MinHeight(tw.S11)

	// Badge / Tag.
	clBadgeBase = tw.New().
			Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S1).
			Rounded(tw.RadiusFull).FontWeight(tw.FontMedium)

	clBadgeVariant = map[string]tw.ClassList{
		"primary":   tw.New().Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand),
		"secondary": tw.New().NoUnderline().Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary),
		"outline":   tw.New().Border(tw.Border1).BorderColor(tw.BorderPrimary).TextColor(tw.FgPrimary),
	}
	clBadgeTone = map[string]tw.ClassList{
		"neutral": tw.New().Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary),
		"brand":   tw.New().Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand),
		"success": tw.New().Bg(tw.SurfaceSuccessSoft).TextColor(tw.FgSuccess),
		"warning": tw.New().Bg(tw.SurfaceWarningSoft).TextColor(tw.FgWarning),
		"danger":  tw.New().Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger),
		"info":    tw.New().Bg(tw.SurfaceInfoSoft).TextColor(tw.FgInfo),
	}
	clBadgeSize = map[string]tw.ClassList{
		"xs":  tw.New().PaddingX(tw.S2).PaddingY(tw.S0_5).FontSize(tw.TextXS),
		"sm":  tw.New().PaddingX(tw.S2_5).PaddingY(tw.S0_5).FontSize(tw.TextXS),
		"md":  tw.New().PaddingX(tw.S2_5).PaddingY(tw.S0_5).FontSize(tw.TextSM),
		"lg":  tw.New().PaddingX(tw.S3).PaddingY(tw.S1).FontSize(tw.TextSM),
		"xl":  tw.New().PaddingX(tw.S3).PaddingY(tw.S1).FontSize(tw.TextBase),
		"2xl": tw.New().PaddingX(tw.S4).PaddingY(tw.S1).FontSize(tw.TextBase),
	}

	clBadgeDot = tw.New().Width(tw.S1_5).Height(tw.S1_5).Rounded(tw.RadiusFull).Bg(tw.FgBrand)

	// Alert.
	clAlertBase = tw.New().
			Display(tw.DisplayFlex).Gap(tw.S3).Rounded(tw.RadiusLG).Padding(tw.S4).
			Border(tw.Border1)

	clAlertVariant = map[string]tw.ClassList{
		"success": tw.New().NoUnderline().Bg(tw.SurfaceSuccessSoft).TextColor(tw.FgSuccess).BorderColor(tw.BorderSuccess),
		"warning": tw.New().NoUnderline().Bg(tw.SurfaceWarningSoft).TextColor(tw.FgWarning).BorderColor(tw.BorderWarning),
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
		On(tw.StatePlaceholder, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPlaceholder) }).
		On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceDisabled).Cursor(tw.CursorNotAllowed)
		}).
		Merge(clFocusRing)

	clInputNormal = tw.New().BorderColor(tw.BorderPrimary)
	clInputError  = tw.New().BorderColor(tw.BorderDanger)
	clInputSize   = map[string]tw.ClassList{
		"sm": tw.New().PaddingX(tw.S3).PaddingY(tw.S1_5).FontSize(tw.TextSM),
		"md": tw.New().PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM),
		"lg": tw.New().PaddingX(tw.S4).PaddingY(tw.S2_5).FontSize(tw.TextBase),
	}
	clInputIconWrap  = tw.New().Position(tw.PositionRelative)
	clInputIconStart = tw.New().
				Position(tw.PositionAbsolute).InsetY(tw.S0).Left(tw.S0).
				Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				PaddingLeft(tw.S3).PointerEvents(tw.PointerNone)
	clInputIconEnd = tw.New().
			Position(tw.PositionAbsolute).InsetY(tw.S0).Right(tw.S0).
			Display(tw.DisplayFlex).Items(tw.ItemsCenter).
			PaddingRight(tw.S3).PointerEvents(tw.PointerNone)
	clInputPadStart = tw.New().PaddingLeft(tw.S10)
	clInputPadEnd   = tw.New().PaddingRight(tw.S10)

	clCheckbox = tw.New().
			Width(tw.S4).Height(tw.S4).Rounded(tw.RadiusSM).
			Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Cursor(tw.CursorPointer).Merge(clFocusRing)

	clCheckRow = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S2)

	// Search bar: a bordered wrapper that reads as one field, hosting a
	// borderless input so the glyph and the text share the control.
	clSearchWrap = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S2).
			Rounded(tw.RadiusMD).Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Bg(tw.SurfacePrimary).PaddingX(tw.S3).Width(tw.SFull).Flex1().MaxWScaled(tw.MaxWMD).
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
	clDividerH         = tw.New().Border(tw.Border0).BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).MarginY(tw.S4)
	clDividerV         = tw.New().Border(tw.Border0).BorderLeft(tw.Border1).BorderColor(tw.BorderPrimary).MarginX(tw.S4)
	clDividerText      = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S4)
	clDividerTextLine  = tw.New().Flex1().BorderTop(tw.Border1).BorderColor(tw.BorderPrimary)
	clDividerTextLabel = tw.New().FontSize(tw.TextSM).TextColor(tw.FgTertiary)

	clSpinner = tw.New().
			Display(tw.DisplayInlineBlock).Rounded(tw.RadiusFull).
			Border(tw.Border2).BorderColor(tw.BorderSecondary).
			AnimateSpin()
	clSpinnerSize = map[string]tw.ClassList{
		"xs":  tw.New().Width(tw.S3).Height(tw.S3),
		"sm":  tw.New().Width(tw.S4).Height(tw.S4),
		"md":  tw.New().Width(tw.S6).Height(tw.S6),
		"lg":  tw.New().Width(tw.S8).Height(tw.S8),
		"xl":  tw.New().Width(tw.S10).Height(tw.S10),
		"2xl": tw.New().Width(tw.S12).Height(tw.S12),
	}
	clSpinnerTone = map[string]tw.ClassList{
		"brand":   tw.New().BorderTopColor(tw.FgBrand),
		"success": tw.New().BorderTopColor(tw.FgSuccess),
		"warning": tw.New().BorderTopColor(tw.FgWarning),
		"danger":  tw.New().BorderTopColor(tw.FgDanger),
		"info":    tw.New().BorderTopColor(tw.FgInfo),
	}

	// Skeleton — pulsing placeholder that holds the geometry of content that
	// has not arrived yet. Same feedback family as Spinner. The base carries
	// only surface and motion; each shape owns its dimensions and radius, so
	// merges never stack two rounded-* utilities (Merge is a plain append).
	clSkeleton     = tw.New().Bg(tw.SurfaceTertiary).AnimatePulse()
	clSkeletonText = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S2)
	// Text lines fill the row; the last line is capped at the named xs width
	// (percent widths are outside tw's enumerable universe) so a paragraph
	// placeholder reads as prose, not a solid slab.
	clSkeletonLine      = tw.New().Width(tw.SFull).Rounded(tw.RadiusMD)
	clSkeletonLineLast  = tw.New().Width(tw.SFull).MaxWScaled(tw.MaxWXS).Rounded(tw.RadiusMD)
	clSkeletonBlockSize = map[string]tw.ClassList{
		"sm": tw.New().Width(tw.SFull).Height(tw.S16).Rounded(tw.RadiusMD),
		"md": tw.New().Width(tw.SFull).Height(tw.S24).Rounded(tw.RadiusMD),
		"lg": tw.New().Width(tw.SFull).Height(tw.S40).Rounded(tw.RadiusMD),
	}
	clSkeletonLineSize = map[string]tw.ClassList{
		"sm": tw.New().Height(tw.S3),
		"md": tw.New().Height(tw.S4),
		"lg": tw.New().Height(tw.S5),
	}
	clSkeletonCircleSize = map[string]tw.ClassList{
		"sm": tw.New().Width(tw.S8).Height(tw.S8).Rounded(tw.RadiusFull),
		"md": tw.New().Width(tw.S12).Height(tw.S12).Rounded(tw.RadiusFull),
		"lg": tw.New().Width(tw.S16).Height(tw.S16).Rounded(tw.RadiusFull),
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
		Display(tw.DisplayInlineBlock).FontFamily(tw.FontMono).
		Rounded(tw.RadiusSM).Border(tw.Border1).BorderColor(tw.BorderPrimary).
		BorderBottom(tw.Border2).Bg(tw.SurfaceTertiary).TextColor(tw.FgSecondary).
		PaddingY(tw.S0_5)
	clKbdSize = map[string]tw.ClassList{
		"xs": tw.New().FontSize(tw.TextXS).PaddingX(tw.S1),
		"sm": tw.New().FontSize(tw.TextXS).PaddingX(tw.S1_5),
		"md": tw.New().FontSize(tw.TextSM).PaddingX(tw.S2),
		"lg": tw.New().FontSize(tw.TextBase).PaddingX(tw.S2_5),
	}

	clLink = tw.New().TextColor(tw.FgLink).Underline().UnderlineOffset(tw.S2).
		On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgLinkHover) }).
		Merge(clFocusRing)

	clTagBase = tw.New().
			Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S1).
			Rounded(tw.RadiusMD).Border(tw.Border1).
			PaddingX(tw.S2).PaddingY(tw.S0_5).FontSize(tw.TextXS)
	clTagIdle     = tw.New().BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).TextColor(tw.FgSecondary)
	clTagSelected = tw.New().BorderColor(tw.BorderBrand).Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand)
	clTagTone     = map[string]tw.ClassList{
		"neutral": clTagIdle,
		"brand":   clTagSelected,
		"success": tw.New().BorderColor(tw.BorderSuccess).Bg(tw.SurfaceSuccessSoft).TextColor(tw.FgSuccess),
		"warning": tw.New().BorderColor(tw.BorderWarning).Bg(tw.SurfaceWarningSoft).TextColor(tw.FgWarning),
		"danger":  tw.New().BorderColor(tw.BorderDanger).Bg(tw.SurfaceDangerSoft).TextColor(tw.FgDanger),
		"info":    tw.New().BorderColor(tw.BorderInfo).Bg(tw.SurfaceInfoSoft).TextColor(tw.FgInfo),
	}

	// Layouts.
	clStack     = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol)
	clFlex      = tw.New().Display(tw.DisplayFlex)
	clGrid      = tw.New().Display(tw.DisplayGrid)
	clContainer = tw.New().MarginX(tw.SAuto).Width(tw.SFull).PaddingX(tw.S4)

	// Canonical OSS solution page shell. It is declared alongside the
	// component primitives so browser CSS and native design compilation share
	// the exact same utility vocabulary.
	clDataManagementPage = tw.New().MinHeight(tw.S96).Bg(tw.SurfaceSecondary).PaddingY(tw.S8)

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

	// DataGrid organism: arrangement only — every part inside is a
	// molecule or atom carrying its own component classes.
	clGridSection = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S4)
	clGridToolbar = tw.New().
			Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Items(tw.ItemsStretch).Gap(tw.S3).FlexWrap().
			Breakpoint(tw.BreakpointSM, func(cl tw.ClassList) tw.ClassList {
			return cl.FlexDir(tw.FlexRow).Items(tw.ItemsCenter)
		})
	clGridActions = tw.New().
			Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).
			Breakpoint(tw.BreakpointSM, func(cl tw.ClassList) tw.ClassList {
			return cl.MarginLeft(tw.SAuto)
		})

	// WindowedCollection.
	clWindow      = tw.New().MinWidth(tw.S0)
	clWindowItems = tw.New().MinWidth(tw.S0)
	clWindowState = tw.New().MarginTop(tw.S4).Rounded(tw.RadiusLG).
			Border(tw.Border1).Padding(tw.S4).BreakWords()
	clWindowLoading = clWindowState.Bg(tw.SurfaceInfoSoft).
			BorderColor(tw.BorderInfo).TextColor(tw.FgInfo)
	clWindowEmpty = clWindowState.Bg(tw.SurfaceSecondary).
			BorderColor(tw.BorderPrimary).TextColor(tw.FgSecondary)
	clWindowError = clWindowState.Bg(tw.SurfaceDangerSoft).
			BorderColor(tw.BorderDanger).TextColor(tw.FgDanger)
	clWindowTitle           = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontSemibold).BreakWords()
	clWindowDescription     = tw.New().MarginTop(tw.S1).FontSize(tw.TextSM).BreakWords()
	clWindowRetry           = tw.New().MarginTop(tw.S3).Merge(clPageBtn).Merge(clPageIdle).BreakWords()
	clWindowFooter          = tw.New().MarginTop(tw.S4)
	clWindowNavigationError = tw.New().MarginTop(tw.S4).FontSize(tw.TextSM).
				TextColor(tw.FgDanger).BreakWords()

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
		clIcon, clFocusRing, clButtonBase, clButtonFull, clButtonIconOnly,
		clBadgeBase, clBadgeDot,
		clAlertBase, clAlertTitle, clAlertMessage, clAlertBody,
		clFieldWrap, clLabel, clHelp, clFieldErr, clRequired,
		clInput, clInputNormal, clInputError,
		clInputIconWrap, clInputIconStart, clInputIconEnd, clInputPadStart, clInputPadEnd,
		clCheckbox, clCheckRow,
		clSearchWrap, clSearchInput, clSrOnly,
		clTruncate, clHeadingBase,
		clDividerH, clDividerV, clDividerText, clDividerTextLine, clDividerTextLabel,
		clSpinner, clEmpty, clEmptyPad, clEmptyBordered, clEmptyCompact, clEmptyTitle, clEmptyDesc,
		clSkeleton, clSkeletonText, clSkeletonLine, clSkeletonLineLast,
		clKbd, clLink, clTagBase, clTagIdle, clTagSelected,
		clStack, clFlex, clGrid, clContainer, clDataManagementPage,
		clTableWrap, clTable, clTableHead, clTableThBase, clTableTh, clTableTd, clTableRow, clTableTdC,
		clTableThSort, clTableSortBtn, clTableRowAlt, clTableTdStrong, clTableActions, clTableCellNote,
		clCard, clCardClickable, clCardTitle, clCardDesc,
		clBreadcrumb, clBreadcrumbSep, clBreadcrumbCur,
		clGridSection, clGridToolbar, clGridActions,
		clWindow, clWindowItems, clWindowState, clWindowLoading, clWindowEmpty,
		clWindowError, clWindowTitle, clWindowDescription, clWindowRetry,
		clWindowFooter, clWindowNavigationError,
		clPagination, clPageBtn, clPageIdle, clPageCur, clPageLabel,
		clTabList, clTab, clTabIdle, clTabActive,
	}
	for _, m := range []map[string]tw.ClassList{
		clButtonVariant, clButtonTone, clButtonSize,
		clBadgeVariant, clBadgeTone, clBadgeSize, clAlertVariant,
		clIconSize, clIconTone,
		clInputSize, clSpinnerSize, clSpinnerTone, clKbdSize, clTagTone,
		clSkeletonBlockSize, clSkeletonLineSize, clSkeletonCircleSize,
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
