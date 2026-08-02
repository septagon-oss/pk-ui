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

import "github.com/septagon-oss/tw"

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

	clAvatarBase = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).Position(tw.PositionRelative).
			Overflow(tw.OverflowHidden).FlexShrink0()
	clAvatarSize = map[string]tw.ClassList{
		"xs":  tw.New().Width(tw.S6).Height(tw.S6),
		"sm":  tw.New().Width(tw.S8).Height(tw.S8),
		"md":  tw.New().Width(tw.S10).Height(tw.S10),
		"lg":  tw.New().Width(tw.S12).Height(tw.S12),
		"xl":  tw.New().Width(tw.S14).Height(tw.S14),
		"2xl": tw.New().Width(tw.S16).Height(tw.S16),
	}
	clAvatarShape = map[string]tw.ClassList{
		"circle":  tw.New().Rounded(tw.RadiusFull),
		"rounded": tw.New().Rounded(tw.RadiusMD),
		"square":  tw.New().Rounded(tw.RadiusNone),
		"pill":    tw.New().Rounded(tw.RadiusFull),
	}
	clAvatarTone = map[string]tw.ClassList{
		"neutral": tw.New().Bg(tw.SurfaceTertiary).TextColor(tw.FgPrimary),
		"brand":   tw.New().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand),
		"success": tw.New().Bg(tw.SurfaceSuccess).TextColor(tw.FgOnBrand),
		"warning": tw.New().Bg(tw.SurfaceWarning).TextColor(tw.FgOnBrand),
		"danger":  tw.New().Bg(tw.SurfaceDanger).TextColor(tw.FgOnBrand),
		"info":    tw.New().Bg(tw.SurfaceInfo).TextColor(tw.FgOnBrand),
	}
	clAvatarInitials     = tw.New().FontWeight(tw.FontMedium)
	clAvatarInitialsSize = map[string]tw.ClassList{
		"xs":  tw.New().FontSize(tw.TextXS),
		"sm":  tw.New().FontSize(tw.TextSM),
		"md":  tw.New().FontSize(tw.TextBase),
		"lg":  tw.New().FontSize(tw.TextLG),
		"xl":  tw.New().FontSize(tw.TextXL),
		"2xl": tw.New().FontSize(tw.Text2XL),
	}
	clAvatarImage  = tw.New().Width(tw.SFull).Height(tw.SFull).ObjectCover()
	clAvatarStatus = tw.New().Position(tw.PositionAbsolute).Display(tw.DisplayBlock).
			Rounded(tw.RadiusFull).Ring(tw.Ring2).RingColor(tw.SurfacePrimary)
	clAvatarStatusTone = map[string]tw.ClassList{
		"online":  tw.New().Bg(tw.SurfaceSuccess),
		"offline": tw.New().Bg(tw.FgPlaceholder),
		"busy":    tw.New().Bg(tw.SurfaceDanger),
		"away":    tw.New().Bg(tw.SurfaceWarning),
	}
	clAvatarStatusSize = map[string]tw.ClassList{
		"xs":  tw.New().Width(tw.S2).Height(tw.S2),
		"sm":  tw.New().Width(tw.S2_5).Height(tw.S2_5),
		"md":  tw.New().Width(tw.S3).Height(tw.S3),
		"lg":  tw.New().Width(tw.S3).Height(tw.S3),
		"xl":  tw.New().Width(tw.S4).Height(tw.S4),
		"2xl": tw.New().Width(tw.S4).Height(tw.S4),
	}
	clAvatarStatusPosition = map[string]tw.ClassList{
		"top-right": tw.New().Top(tw.S0).Right(tw.S0).
			TranslateX(tw.TranslateHalf).TranslateY(tw.TranslateNegHalf),
		"top-left": tw.New().Top(tw.S0).Left(tw.S0).
			TranslateX(tw.TranslateNegHalf).TranslateY(tw.TranslateNegHalf),
		"bottom-left": tw.New().Bottom(tw.S0).Left(tw.S0).
			TranslateX(tw.TranslateNegHalf).TranslateY(tw.TranslateHalf),
		"bottom-right": tw.New().Bottom(tw.S0).Right(tw.S0).
			TranslateX(tw.TranslateHalf).TranslateY(tw.TranslateHalf),
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

	clButtonFull         = tw.New().Width(tw.SFull)
	clButtonIconOnly     = tw.New().MinWidth(tw.S11).MinHeight(tw.S11)
	clButtonDisabledLink = tw.New().
				Cursor(tw.CursorNotAllowed).
				Opacity(tw.Opacity50).
				PointerEvents(tw.PointerNone)

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

	clBadgeDot     = tw.New().Width(tw.S1_5).Height(tw.S1_5).Rounded(tw.RadiusFull)
	clBadgeDotTone = map[string]tw.ClassList{
		"neutral": tw.New().Bg(tw.FgSecondary),
		"brand":   tw.New().Bg(tw.FgBrand),
		"success": tw.New().Bg(tw.FgSuccess),
		"warning": tw.New().Bg(tw.FgWarning),
		"danger":  tw.New().Bg(tw.FgDanger),
		"info":    tw.New().Bg(tw.FgInfo),
	}
	clBadgeCount  = tw.New().MarginLeft(tw.S1).FontWeight(tw.FontSemibold).TabularNums()
	clBadgeRemove = tw.New().MarginLeft(tw.S1).Display(tw.DisplayInlineFlex).
			Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
			Bg(tw.ColorTransparent).Border(tw.Border0).Padding(tw.S0).
			Rounded(tw.RadiusFull).Cursor(tw.CursorPointer).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Opacity(tw.Opacity75) }).
			Merge(clFocusRing)

	// Alert.
	clAlertBase = tw.New().
			Display(tw.DisplayFlex).Items(tw.ItemsStart).Gap(tw.S3).Rounded(tw.RadiusLG).
			Border(tw.Border1)
	clAlertRegular  = tw.New().Padding(tw.S4)
	clAlertCompact  = tw.New().PaddingX(tw.S3).PaddingY(tw.S2)
	clAlertBordered = tw.New().BorderLeft(tw.Border4)

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
	clAlertIcon    = tw.New().
			MarginTop(tw.S0_5).Display(tw.DisplayFlex).Height(tw.S9).Width(tw.S9).
			FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).Rounded(tw.RadiusFull)
	clAlertActions = tw.New().MarginTop(tw.S3).Display(tw.DisplayFlex).FlexWrap().
			Items(tw.ItemsCenter).Gap(tw.S3).FontSize(tw.TextSM)
	clAlertClose = tw.New().MarginLeft(tw.SAuto).Display(tw.DisplayInlineFlex).
			FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
			Rounded(tw.RadiusMD).Padding(tw.S1_5).Cursor(tw.CursorPointer).
			Transition(tw.TransitionColors).Merge(clFocusRing)

	// Toast.
	clToastBase = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart).Gap(tw.S3).
			MaxW("sm").Width(tw.SFull).Rounded(tw.RadiusLG).Border(tw.Border1).
			Shadow(tw.ShadowLG).PointerEvents(tw.PointerAuto).Overflow(tw.OverflowHidden).
			Padding(tw.S4)
	clToastTone = map[string]tw.ClassList{
		"success": clAlertVariant["success"].Merge(tw.New().BorderLeft(tw.Border4)),
		"warning": clAlertVariant["warning"].Merge(tw.New().BorderLeft(tw.Border4)),
		"danger":  clAlertVariant["danger"].Merge(tw.New().BorderLeft(tw.Border4)),
		"info":    clAlertVariant["info"].Merge(tw.New().BorderLeft(tw.Border4)),
	}
	clToastIcon = tw.New().MarginTop(tw.S0_5).Display(tw.DisplayFlex).FlexShrink0()
	clToastBody = tw.New().MinWidth(tw.S0).Flex1().Display(tw.DisplayFlex).
			FlexDir(tw.FlexCol).Gap(tw.S1)
	clToastTitle   = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clToastMessage = tw.New().FontSize(tw.TextSM).TextColor(tw.FgSecondary)
	clToastLead    = tw.New().FontWeight(tw.FontMedium).TextColor(tw.FgPrimary)
	clToastClose   = tw.New().MarginLeft(tw.SAuto).Display(tw.DisplayInlineFlex).
			FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
			Rounded(tw.RadiusMD).Padding(tw.S1_5).Cursor(tw.CursorPointer).
			Transition(tw.TransitionColors).Merge(clFocusRing)

	// Inputs.
	clFieldWrap     = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1_5)
	clFieldWrapFull = tw.New().Width(tw.SFull)
	clLabel         = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).TextColor(tw.FgPrimary)
	clHelp          = tw.New().FontSize(tw.TextXS).TextColor(tw.FgMuted)
	clFieldErr      = tw.New().FontSize(tw.TextXS).TextColor(tw.FgDanger)
	clRequired      = tw.New().TextColor(tw.FgDanger)

	clInput = tw.New().
		Display(tw.DisplayBlock).Width(tw.SFull).
		Rounded(tw.RadiusMD).Border(tw.Border1).
		Bg(tw.SurfacePrimary).TextColor(tw.FgPrimary).
		On(tw.StatePlaceholder, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPlaceholder) }).
		On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceDisabled).Cursor(tw.CursorNotAllowed)
		}).
		Merge(clFocusRing)

	clInputNormal   = tw.New().BorderColor(tw.BorderPrimary)
	clInputError    = tw.New().BorderColor(tw.BorderDanger)
	clInputReadOnly = tw.New().Bg(tw.SurfaceSecondary).TextColor(tw.FgSecondary).
			Shadow(tw.ShadowNone).Cursor(tw.CursorDefault)
	clInputTone = map[string]tw.ClassList{
		"neutral": tw.New().BorderColor(tw.BorderPrimary),
		"success": tw.New().BorderColor(tw.BorderSuccess),
		"warning": tw.New().BorderColor(tw.BorderWarning),
		"danger":  tw.New().BorderColor(tw.BorderDanger),
	}
	clInputSize = map[string]tw.ClassList{
		"sm": tw.New().PaddingX(tw.S3).PaddingY(tw.S1_5).FontSize(tw.TextSM),
		"md": tw.New().PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM),
		"lg": tw.New().PaddingX(tw.S4).PaddingY(tw.S2_5).FontSize(tw.TextBase),
	}
	clTextareaManual = tw.New().ResizeY()
	clTextareaAuto   = tw.New().ResizeNone().Overflow(tw.OverflowHidden)
	clTextareaFull   = tw.New().Width(tw.SFull)
	clTextareaMeta   = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart).
				Justify(tw.JustifyBetween).Gap(tw.S2)
	clTextareaSupporting = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1)
	clTextareaCounter    = tw.New().FontSize(tw.TextXS).TextColor(tw.FgMuted).TabularNums().FlexShrink0()
	clInputIconWrap      = tw.New().Position(tw.PositionRelative)
	clInputIconStart     = tw.New().
				Position(tw.PositionAbsolute).InsetY(tw.S0).Left(tw.S0).
				Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				PaddingLeft(tw.S3).PointerEvents(tw.PointerNone)
	clInputIconEnd = tw.New().
			Position(tw.PositionAbsolute).InsetY(tw.S0).Right(tw.S0).
			Display(tw.DisplayFlex).Items(tw.ItemsCenter).
			PaddingRight(tw.S3).PointerEvents(tw.PointerNone)
	clInputPadStart = tw.New().PaddingLeft(tw.S10)
	clInputPadEnd   = tw.New().PaddingRight(tw.S10)

	// Autocomplete extends the input family with an APG combobox listbox.
	clAutocompleteControl = tw.New().Position(tw.PositionRelative).Width(tw.SFull)
	clAutocompletePanel   = tw.New().Position(tw.PositionAbsolute).ZIndex(tw.ZDropdown).
				MarginTop(tw.S1).Width(tw.SFull).
				Overflow(tw.OverflowAuto).Rounded(tw.RadiusMD).Border(tw.Border1).
				BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).Shadow(tw.ShadowXL)
	clAutocompleteList   = tw.New().PaddingY(tw.S1)
	clAutocompleteOption = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).
				Cursor(tw.CursorPointer).Items(tw.ItemsCenter).PaddingX(tw.S3).
				PaddingY(tw.S2).TextAlign(tw.TextLeft).FontSize(tw.TextSM).
				TextColor(tw.FgPrimary).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceSecondary) })
	clAutocompleteOptionActive = tw.New().Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand)
	clAutocompleteIndicator    = tw.New().Position(tw.PositionAbsolute).InsetY(tw.S0).
					Right(tw.S0).Display(tw.DisplayFlex).Items(tw.ItemsCenter).
					PaddingRight(tw.S3).PointerEvents(tw.PointerNone)
	clAutocompleteSpinner = tw.New().AnimateSpin().Height(tw.S4).Width(tw.S4).
				Rounded(tw.RadiusFull).Border(tw.Border2).BorderColor(tw.BorderBrand).
				BorderTop(tw.Border2).BorderStyle(tw.BorderSolid)

	// FileUpload owns one progressively enhanced DOM contract for native form
	// submission and provider-backed remote uploads. Dynamic list nodes use the
	// same declared class lists, so stylesheet emission never depends on a JS
	// source scan.
	clFileUploadRoot = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).
				FlexDir(tw.FlexCol).Gap(tw.S2)
	clFileUploadDropZone = tw.New().Position(tw.PositionRelative).Width(tw.SFull).
				Border(tw.Border2).BorderStyle(tw.BorderStyle("dashed")).
				BorderColor(tw.BorderSecondary).Rounded(tw.RadiusLG).
				Cursor(tw.CursorPointer).Transition(tw.TransitionColors).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.BorderColor(tw.BorderPrimary) }).
				On(tw.StateFocusWithin, func(c tw.ClassList) tw.ClassList {
			return c.Ring(tw.Ring2).RingColor(tw.RingFocus).RingOffset(tw.RingOffset1)
		})
	clFileUploadDropZoneDisabled = tw.New().PointerEvents(tw.PointerNone).
					Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity60)
	clFileUploadDropZoneInner = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).
					Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
					PaddingX(tw.S4).PaddingY(tw.S8).TextAlign(tw.TextCenter)
	clFileUploadInputHidden = tw.New().SrOnly().MinHeight(tw.S0).MinWidth(tw.S0)
	clFileUploadInput       = clInput.Merge(clInputNormal).Merge(clInputSize["md"]).Width(tw.SFull)
	clFileUploadIcon        = tw.New().Height(tw.S10).Width(tw.S10).
				TextColor(tw.FgTertiary).MarginBottom(tw.S3)
	clFileUploadPrompt = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).
				Items(tw.ItemsCenter).Justify(tw.JustifyCenter)
	clFileUploadPromptText   = tw.New().FontSize(tw.TextSM).TextColor(tw.FgSecondary)
	clFileUploadPromptAction = tw.New().FontWeight(tw.FontMedium).
					TextColor(tw.FgBrand)
	clFileUploadHint = tw.New().MarginTop(tw.S1).FontSize(tw.TextXS).
				TextColor(tw.FgTertiary)
	clFileUploadLoading = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).
				Items(tw.ItemsCenter).Justify(tw.JustifyCenter)
	clFileUploadLoadingIcon = tw.New().AnimateSpin().Height(tw.S8).Width(tw.S8).
				TextColor(tw.FgMuted).MarginBottom(tw.S2)
	clFileUploadList  = tw.New().MarginTop(tw.S1).SpaceY(tw.S2)
	clFileUploadError = tw.New().FontSize(tw.TextSM).TextColor(tw.FgDanger)
	clFileUploadItem  = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Justify(tw.JustifyBetween).PaddingX(tw.S3).PaddingY(tw.S2).
				Bg(tw.SurfaceSecondary).Rounded(tw.RadiusMD)
	clFileUploadItemMeta = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Gap(tw.S3).MinWidth(tw.S0)
	clFileUploadPreview = tw.New().Height(tw.S10).Width(tw.S10).
				Rounded(tw.RadiusBase).ObjectCover().FlexShrink0()
	clFileUploadPreviewLarge = tw.New().Height(tw.S16).Width(tw.S16).
					Rounded(tw.RadiusBase).ObjectCover().FlexShrink0()
	clFileUploadItemCopy = tw.New().MinWidth(tw.S0)
	clFileUploadItemName = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).
				TextColor(tw.FgPrimary).Truncate()
	clFileUploadItemSize = tw.New().FontSize(tw.TextXS).TextColor(tw.FgTertiary)
	clFileUploadRemove   = tw.New().FlexShrink0().MarginLeft(tw.S2).Padding(tw.S1).
				TextColor(tw.FgTertiary).Rounded(tw.RadiusBase).
				Transition(tw.TransitionColors).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgDanger) })
	clFileUploadDragActive = tw.New().BorderColor(tw.BorderBrand).Bg(tw.SurfaceBrandSoft)
	clFileUploadDragIdle   = tw.New().BorderColor(tw.BorderSecondary)

	// Dropdown is the canonical select-only listbox shell. The controller owns
	// open state, filtering, selection, hidden-input synchronization, and focus.
	clDropdownRoot = tw.New().Position(tw.PositionRelative).Display(tw.DisplayFlex).
			Width(tw.SFull).FlexDir(tw.FlexCol)
	clDropdownTrigger = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).
				Items(tw.ItemsCenter).Overflow(tw.OverflowHidden).Rounded(tw.RadiusMD).
				Border(tw.Border1).BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).
				TextColor(tw.FgPrimary).Shadow(tw.ShadowSM).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceSecondary) }).
				On(tw.StateFocusWithin, func(c tw.ClassList) tw.ClassList {
			return c.Outline(tw.OutlineNone).Ring(tw.Ring2).RingColor(tw.RingFocus).RingOffset(tw.RingOffset1)
		})
	clDropdownTriggerSize = map[string]tw.ClassList{
		"sm": tw.New().MinHeight(tw.S9).FontSize(tw.TextSM),
		"md": tw.New().MinHeight(tw.S11).FontSize(tw.TextSM),
		"lg": tw.New().MinHeight(tw.S12).FontSize(tw.TextBase),
	}
	clDropdownButton = tw.New().Display(tw.DisplayFlex).MinWidth(tw.S0).Flex1().
				Items(tw.ItemsCenter).TextAlign(tw.TextLeft).Outline(tw.OutlineNone).
				On(tw.StateDisabled, func(c tw.ClassList) tw.ClassList {
			return c.Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
		})
	clDropdownButtonSize = map[string]tw.ClassList{
		"sm": tw.New().MinHeight(tw.S9).PaddingX(tw.S2).PaddingY(tw.S1_5),
		"md": tw.New().MinHeight(tw.S11).PaddingX(tw.S3).PaddingY(tw.S2),
		"lg": tw.New().MinHeight(tw.S12).PaddingX(tw.S4).PaddingY(tw.S2_5),
	}
	clDropdownTriggerLabel   = tw.New().MinWidth(tw.S0).Flex1().Truncate()
	clDropdownTriggerActions = tw.New().Display(tw.DisplayFlex).FlexShrink0().
					Items(tw.ItemsCenter).Gap(tw.S1).PaddingRight(tw.S2)
	clDropdownIconButton = tw.New().Display(tw.DisplayInlineFlex).Height(tw.S8).Width(tw.S8).
				FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
				Rounded(tw.RadiusMD).TextColor(tw.FgTertiary).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceSecondary).TextColor(tw.FgPrimary)
		})
	clDropdownChevron = tw.New().Display(tw.DisplayInlineFlex).Height(tw.S8).Width(tw.S8).
				FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).
				TextColor(tw.FgTertiary).Transition(tw.TransitionTransform)
	clDropdownPanel = tw.New().Position(tw.PositionAbsolute).ZIndex(tw.ZDropdown).
			MarginTop(tw.S1).Width(tw.SFull).Overflow(tw.OverflowHidden).
			Rounded(tw.Radius2XL).Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Bg(tw.SurfacePrimary).Shadow(tw.ShadowXL)
	clDropdownSearchWrap = tw.New().BorderBottom(tw.Border1).BorderColor(tw.BorderPrimary).
				Padding(tw.S2)
	clDropdownSearch = tw.New().MinHeight(tw.S10).Width(tw.SFull).Rounded(tw.RadiusLG).
				Border(tw.Border1).BorderColor(tw.BorderSecondary).Bg(tw.SurfacePrimary).
				PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM).TextColor(tw.FgPrimary).
				On(tw.StatePlaceholder, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgTertiary) }).
				On(tw.StateFocus, func(c tw.ClassList) tw.ClassList {
			return c.Outline(tw.OutlineNone).Ring(tw.Ring2).RingColor(tw.RingFocus).RingOffset(tw.RingOffset1)
		})
	clDropdownOptions = tw.New().MaxHeight(tw.S60).OverflowY(tw.OverflowAuto).PaddingY(tw.S1)
	clDropdownOption  = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).Items(tw.ItemsCenter).
				Gap(tw.S2).PaddingX(tw.S3).PaddingY(tw.S2).TextAlign(tw.TextLeft).
				FontSize(tw.TextSM).TextColor(tw.FgPrimary).Cursor(tw.CursorPointer).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceSecondary) })
	clDropdownOptionSelected = tw.New().Bg(tw.SurfaceSecondary)
	clDropdownOptionDisabled = tw.New().Opacity(tw.Opacity50).Cursor(tw.CursorNotAllowed)
	clDropdownOptionMark     = tw.New().Display(tw.DisplayFlex).Height(tw.S4).Width(tw.S4).
					FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter)
	clDropdownOptionSpacer = tw.New().Height(tw.S4).Width(tw.S4).FlexShrink0()
	clDropdownOptionIcon   = tw.New().Display(tw.DisplayFlex).Height(tw.S4).Width(tw.S4).
				FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).TextColor(tw.FgTertiary)
	clDropdownOptionLabel = tw.New().MinWidth(tw.S0).Flex1().Truncate()
	clDropdownGroupLabel  = tw.New().PaddingX(tw.S3).PaddingY(tw.S1_5).
				FontSize(tw.TextXS).FontWeight(tw.FontMedium).TextColor(tw.FgTertiary)

	// ActionMenu implements the APG menu-button visual contract. Alignment
	// and width are governed presets; callers never pass raw utility classes.
	clActionMenuRoot = tw.New().Position(tw.PositionRelative).Display(tw.DisplayInlineBlock).
				TextAlign(tw.TextLeft)
	clActionMenuTrigger = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
				Justify(tw.JustifyCenter).Gap(tw.S1_5).Rounded(tw.RadiusLG).
				Border(tw.Border1).BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).
				PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM).
				FontWeight(tw.FontMedium).TextColor(tw.FgPrimary).Shadow(tw.ShadowSM).
				Transition(tw.TransitionColors).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceSecondary) }).
				Merge(clFocusRing)
	clActionMenuPanel = tw.New().Position(tw.PositionAbsolute).ZIndex(tw.ZDropdown).
				MarginTop(tw.S1).Rounded(tw.RadiusLG).Border(tw.Border1).
				BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).Shadow(tw.ShadowLG).
				Ring(tw.Ring1).RingColor(tw.BorderPrimary).
				On(tw.StateFocus, func(c tw.ClassList) tw.ClassList { return c.Outline(tw.OutlineNone) })
	clActionMenuAlign = map[string]tw.ClassList{
		"start": tw.New().Left(tw.S0),
		"end":   tw.New().Right(tw.S0),
	}
	clActionMenuWidth = map[string]tw.ClassList{
		"sm": tw.New().Width(tw.S48),
		"md": tw.New().Width(tw.S56),
		"lg": tw.New().Width(tw.S64),
	}
	clActionMenuPanelInner = tw.New().PaddingY(tw.S1)
	clActionMenuSeparator  = tw.New().MarginY(tw.S1).BorderTop(tw.Border1).
				BorderColor(tw.BorderPrimary)
	clActionMenuSectionLabel = tw.New().PaddingX(tw.S3).PaddingY(tw.S1_5).
					FontSize(tw.TextXS).FontWeight(tw.FontSemibold).Uppercase().
					Tracking(tw.TrackingWider).TextColor(tw.FgMuted)
	clActionMenuItem = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).
				Items(tw.ItemsCenter).Gap(tw.S2).PaddingX(tw.S3).PaddingY(tw.S2).
				FontSize(tw.TextSM).Transition(tw.TransitionColors).TextAlign(tw.TextLeft)
	clActionMenuItemTone = map[string]tw.ClassList{
		"neutral": tw.New().TextColor(tw.FgPrimary).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceSecondary) }),
		"danger": tw.New().TextColor(tw.FgDanger).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceDangerSoft) }),
	}
	clActionMenuItemDisabled = tw.New().TextColor(tw.FgMuted).
					Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
	clActionMenuItemIcon  = tw.New().Display(tw.DisplayFlex).FlexShrink0()
	clActionMenuItemLabel = tw.New().Flex1().Truncate()

	// Drawer is the governed modal edge panel used by mobile navigation,
	// filters, and server-loaded details.
	clDrawerRoot    = tw.New().Position(tw.PositionFixed).Inset(tw.S0).ZIndex(tw.ZOverlay)
	clDrawerOverlay = tw.New().Position(tw.PositionAbsolute).Inset(tw.S0).
			BgOpacity(tw.SurfaceOverlay, string(tw.Opacity50)).Transition(tw.TransitionOpacity)
	clDrawerPanel = tw.New().Position(tw.PositionAbsolute).Display(tw.DisplayFlex).
			FlexDir(tw.FlexCol).Bg(tw.SurfacePrimary).Shadow(tw.ShadowXL)
	clDrawerPosition = map[string]tw.ClassList{
		"left":   tw.New().InsetY(tw.S0).Left(tw.S0),
		"right":  tw.New().InsetY(tw.S0).Right(tw.S0),
		"bottom": tw.New().InsetX(tw.S0).Bottom(tw.S0),
	}
	clDrawerWidth = map[string]tw.ClassList{
		"small":  tw.New().Width(tw.S80).MaxWScaled(tw.MaxWFull),
		"medium": tw.New().Width(tw.S96).MaxWScaled(tw.MaxWFull),
		"large":  tw.New().Width(tw.SFull).MaxWScaled(tw.MaxWLG),
		"xl":     tw.New().Width(tw.SFull).MaxWScaled(tw.MaxW2XL),
		"full":   tw.New().Width(tw.SFull),
	}
	clDrawerBottomSize = map[string]tw.ClassList{
		"small":  tw.New().HeightViewport(tw.VH25).Width(tw.SFull).Rounded(tw.Radius("t-xl")),
		"medium": tw.New().HeightViewport(tw.VH50).Width(tw.SFull).Rounded(tw.Radius("t-xl")),
		"large":  tw.New().HeightViewport(tw.VH75).Width(tw.SFull).Rounded(tw.Radius("t-xl")),
		"xl":     tw.New().HeightViewport(tw.VH85).Width(tw.SFull).Rounded(tw.Radius("t-xl")),
		"full":   tw.New().Height(tw.SFull).Width(tw.SFull),
	}
	clDrawerHeader = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart).
			Justify(tw.JustifyBetween).Gap(tw.S4).PaddingX(tw.S6).PaddingY(tw.S4).
			BorderBottom(tw.Border1).BorderColor(tw.BorderPrimary).FlexShrink0()
	clDrawerTitleBlock  = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1).MinWidth(tw.S0)
	clDrawerTitle       = tw.New().FontSize(tw.TextLG).FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clDrawerDescription = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted)
	clDrawerBody        = tw.New().Flex1().OverflowY(tw.OverflowAuto).PaddingX(tw.S6).PaddingY(tw.S4)
	clDrawerFooter      = tw.New().FlexShrink0().BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).
				PaddingX(tw.S6).PaddingY(tw.S4)
	clDrawerClose = tw.New().Display(tw.DisplayInlineFlex).FlexShrink0().Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).Rounded(tw.RadiusMD).Padding(tw.S1_5).
			TextColor(tw.FgMuted).Bg(tw.ColorTransparent).Border(tw.Border0).
			Cursor(tw.CursorPointer).Transition(tw.TransitionColors).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPrimary).Bg(tw.SurfaceHover) }).
			Merge(clFocusRing)

	// Modal is the governed centered-dialog / mobile-sheet overlay. The root
	// also doubles as the empty HTMX swap target used by server-loaded forms.
	clModalRoot = tw.New().Position(tw.PositionFixed).Inset(tw.S0).ZIndex(tw.ZModal).
			Display(tw.DisplayFlex).Justify(tw.JustifyCenter).Padding(tw.S4).
			OverflowY(tw.OverflowAuto)
	clModalCentered    = tw.New().Items(tw.ItemsCenter)
	clModalBottomSheet = tw.New().Items(tw.ItemsEnd).
				Breakpoint(tw.BreakpointSM, func(c tw.ClassList) tw.ClassList { return c.Items(tw.ItemsCenter) })
	clModalOverlay = tw.New().Position(tw.PositionAbsolute).Inset(tw.S0).
			BgOpacity(tw.SurfaceOverlay, string(tw.Opacity50)).Transition(tw.TransitionOpacity)
	clModalPanel = tw.New().Position(tw.PositionRelative).Display(tw.DisplayFlex).
			FlexDir(tw.FlexCol).Width(tw.SFull).MaxHeightViewport(tw.VH85).
			Overflow(tw.OverflowHidden).Rounded(tw.Radius2XL).
			Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Bg(tw.SurfacePrimary).Shadow(tw.Shadow2XL).TextAlign(tw.TextLeft)
	clModalPanelSize = map[string]tw.ClassList{
		"small":  tw.New().MaxWScaled(tw.MaxWSM),
		"medium": tw.New().MaxWScaled(tw.MaxWLG),
		"large":  tw.New().MaxWScaled(tw.MaxW3XL),
		"xl":     tw.New().MaxWScaled(tw.MaxW5XL),
		"full":   tw.New().MaxWScaled(tw.MaxWFull),
	}
	clModalHeader = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart).
			Justify(tw.JustifyBetween).Gap(tw.S4).PaddingX(tw.S6).PaddingY(tw.S4).
			BorderBottom(tw.Border1).BorderColor(tw.BorderPrimary).
			Bg(tw.SurfaceSecondary).FlexShrink0()
	clModalTitleBlock  = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S1).MinWidth(tw.S0)
	clModalTitle       = tw.New().FontSize(tw.TextLG).FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clModalDescription = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted)
	clModalBody        = tw.New().Flex1().OverflowY(tw.OverflowAuto).Padding(tw.S6)
	clModalFooter      = tw.New().FlexShrink0().BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).
				Bg(tw.SurfaceSecondary).PaddingX(tw.S6).PaddingY(tw.S4)
	clModalClose = tw.New().Display(tw.DisplayInlineFlex).FlexShrink0().Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).Rounded(tw.RadiusMD).Padding(tw.S1_5).
			TextColor(tw.FgMuted).Bg(tw.ColorTransparent).Border(tw.Border0).
			Cursor(tw.CursorPointer).Transition(tw.TransitionColors).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.TextColor(tw.FgPrimary).Bg(tw.SurfaceHover) }).
			Merge(clFocusRing)
	clModalCancel = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).Rounded(tw.RadiusMD).Border(tw.Border1).
			BorderColor(tw.BorderPrimary).Bg(tw.SurfacePrimary).PaddingX(tw.S4).PaddingY(tw.S2).
			FontSize(tw.TextSM).FontWeight(tw.FontMedium).TextColor(tw.FgSecondary).
			Cursor(tw.CursorPointer).Transition(tw.TransitionColors).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Bg(tw.SurfaceHover) }).
			Merge(clFocusRing)

	clCheckbox = tw.New().
			Width(tw.S4).Height(tw.S4).Rounded(tw.RadiusSM).
			Border(tw.Border1).BorderColor(tw.BorderPrimary).
			Cursor(tw.CursorPointer).Merge(clFocusRing)

	clCheckRow     = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S2)
	clCheckboxRoot = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsStart).
			Gap(tw.S3).Cursor(tw.CursorPointer)
	clCheckboxRootDisabled = tw.New().Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
	clCheckboxInput        = tw.New().Position(tw.PositionAbsolute).
				Height(tw.SPX).Width(tw.SPX).MinHeight(tw.S0).MinWidth(tw.S0).
				AppearanceNone().Opacity(tw.Opacity0).PointerEvents(tw.PointerNone)
	clCheckboxIndicator = tw.New().MarginTop(tw.S0_5).Display(tw.DisplayFlex).
				Height(tw.S5).Width(tw.S5).FlexShrink0().Items(tw.ItemsCenter).
				Justify(tw.JustifyCenter).Rounded(tw.RadiusMD).Border(tw.Border1).
				Transition(tw.TransitionColors)
	clCheckboxIndicatorIdle = tw.New().BorderColor(tw.BorderPrimary).
				Bg(tw.SurfacePrimary).TextColor(tw.ColorTransparent)
	clCheckboxIndicatorActive = tw.New().BorderColor(tw.BorderBrand).
					Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand)
	clCheckboxCheckmark = tw.New().Height(tw.S3).Width(tw.S3)
	clCheckboxBar       = tw.New().Height(tw.S0_5).Width(tw.S2_5).
				Rounded(tw.RadiusFull).Bg(tw.SurfacePrimary)
	clCheckboxLabel = tw.New().Truncate().PaddingTop(tw.S0_5).
			FontSize(tw.TextSM).TextColor(tw.FgPrimary)

	// Radio uses the native control so browser group selection updates its
	// checked projection without a JavaScript controller. The design system
	// still owns size, semantic accent, focus, disabled, and label treatment.
	clRadioRoot = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
			Gap(tw.S3).Cursor(tw.CursorPointer).UserSelect(tw.SelectNone)
	clRadioRootDisabled = tw.New().Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
	clRadioInput        = tw.New().Height(tw.S5).Width(tw.S5).FlexShrink0().
				Rounded(tw.RadiusFull).Border(tw.Border1).BorderColor(tw.BorderPrimary).
				Accent(tw.FgBrand).Cursor(tw.CursorPointer).Merge(clFocusRing)
	clRadioInputDisabled = tw.New().Cursor(tw.CursorNotAllowed)
	clRadioLabel         = tw.New().Truncate().FontSize(tw.TextSM).TextColor(tw.FgPrimary)
	clRadioDot           = tw.New().Height(tw.S2).Width(tw.S2).Rounded(tw.RadiusFull).Bg(tw.SurfaceBrand)

	// Choice groups compose the canonical Checkbox and Radio atoms inside a
	// native fieldset. The group owns shared labelling, validation, layout,
	// and enhancement while each atom retains native control behavior.
	clChoiceGroupRoot        = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S2)
	clChoiceGroupLegend      = clLabel
	clChoiceGroupDescription = clHelp
	clChoiceGroupOptions     = tw.New().Display(tw.DisplayFlex).Gap(tw.S2)
	clChoiceGroupVertical    = tw.New().FlexDir(tw.FlexCol)
	clChoiceGroupHorizontal  = tw.New().FlexDir(tw.FlexRow).FlexWrap().Gap(tw.S4)
	clChoiceGroupError       = clFieldErr

	clSliderRoot = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S3).
			Width(tw.SFull).MaxWScaled(tw.MaxWXS)
	clSliderLabel = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).
			TextColor(tw.FgPrimary)
	clSliderRow   = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S3)
	clSliderInput = tw.New().Flex1().Height(tw.S2).Bg(tw.SurfaceTertiary).
			Rounded(tw.RadiusFull).AppearanceNone().Cursor(tw.CursorPointer).Merge(clFocusRing)
	clSliderInputDisabled = tw.New().Cursor(tw.CursorNotAllowed).Opacity(tw.Opacity50)
	clSliderValue         = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).
				TextColor(tw.FgSecondary).TabularNums().MinWidth(tw.S8).TextAlign(tw.TextRight)
	clSliderTone = map[string]tw.ClassList{
		"brand":   tw.New().Accent(tw.SurfaceBrand),
		"success": tw.New().Accent(tw.SurfaceSuccess),
		"warning": tw.New().Accent(tw.SurfaceWarning),
		"danger":  tw.New().Accent(tw.SurfaceDanger),
		"info":    tw.New().Accent(tw.SurfaceInfo),
	}

	clToggleRoot         = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter)
	clToggleRootDisabled = tw.New().Opacity(tw.Opacity50)
	clToggleControl      = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).Gap(tw.S3).
				Bg(tw.ColorTransparent).Border(tw.Border0).Padding(tw.S0).Cursor(tw.CursorPointer).
				On(tw.StateDisabled, func(cl tw.ClassList) tw.ClassList {
			return cl.Cursor(tw.CursorNotAllowed)
		})
	clToggleInput = tw.New().SrOnly().MinHeight(tw.S0).MinWidth(tw.S0)
	clToggleTrack = tw.New().Relative().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
			FlexShrink0().Overflow(tw.OverflowHidden).Rounded(tw.RadiusFull).
			Border(tw.Border1).PaddingX(tw.S0_5).Shadow(tw.ShadowSM).
			Transition(tw.TransitionColors).Duration(tw.Duration200)
	clToggleTrackSize = map[string]tw.ClassList{
		"sm": tw.New().Height(tw.S5).Width(tw.S9),
		"md": tw.New().Height(tw.S6).Width(tw.S11),
		"lg": tw.New().Height(tw.S8).Width(tw.S14),
	}
	clToggleTrackState = map[string]tw.ClassList{
		"checked":   tw.New().Bg(tw.SurfaceBrand).BorderColor(tw.SurfaceBrand),
		"unchecked": tw.New().Bg(tw.SurfaceTertiary).BorderColor(tw.BorderPrimary),
	}
	clToggleKnob = tw.New().PointerEvents(tw.PointerNone).Display(tw.DisplayInlineBlock).
			Rounded(tw.RadiusFull).Bg(tw.SurfacePrimary).Shadow(tw.ShadowSM).
			Ring(tw.Ring1).RingColor(tw.BorderPrimary).Transition(tw.TransitionAll).
			Duration(tw.Duration200).Transform()
	clToggleKnobSize = map[string]tw.ClassList{
		"sm": tw.New().Height(tw.S4).Width(tw.S4),
		"md": tw.New().Height(tw.S5).Width(tw.S5),
		"lg": tw.New().Height(tw.S7).Width(tw.S7),
	}
	clToggleKnobChecked = map[string]tw.ClassList{
		"sm": tw.New().TranslateXStep(tw.S4),
		"md": tw.New().TranslateXStep(tw.S5),
		"lg": tw.New().TranslateXStep(tw.S6),
	}
	clToggleKnobUnchecked = tw.New().TranslateXStep(tw.S0)
	clToggleLabel         = tw.New().Truncate().FontSize(tw.TextSM).FontWeight(tw.FontMedium).TextColor(tw.FgSecondary)

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
		"primary": tw.FgPrimary, "secondary": tw.FgSecondary, "tertiary": tw.FgTertiary, "muted": tw.FgMuted,
		"brand": tw.FgBrand, "success": tw.FgSuccess, "warning": tw.FgWarning,
		"danger": tw.FgDanger, "error": tw.FgDanger, "info": tw.FgInfo,
	}
	clTextSize = map[string]tw.FontSize{
		"xs": tw.TextXS, "sm": tw.TextSM, "base": tw.TextBase, "md": tw.TextBase,
		"lg": tw.TextLG, "xl": tw.TextXL, "2xl": tw.Text2XL,
		"3xl": tw.Text3XL, "4xl": tw.Text4XL, "5xl": tw.Text5XL,
	}
	clTextWeight = map[string]tw.FontWeight{
		"thin": tw.FontThin, "extralight": tw.FontExtralight, "light": tw.FontLight,
		"normal": tw.FontNormal, "medium": tw.FontMedium, "semibold": tw.FontSemibold,
		"bold": tw.FontBold, "extrabold": tw.FontExtrabold, "black": tw.FontBlack,
	}
	clTextAlign = map[string]tw.TextAlign{
		"left": tw.TextLeft, "center": tw.TextCenter,
		"right": tw.TextRight, "justify": tw.TextJustify,
	}
	clTextTransform = map[string]tw.ClassList{
		"none": {}, "uppercase": tw.New().Uppercase(),
		"lowercase": tw.New().Lowercase(), "capitalize": tw.New().Capitalize(),
	}
	clTextItalic    = tw.New().Italic()
	clTextUnderline = tw.New().Underline()
	clTextNoWrap    = tw.New().WhitespaceNowrap()
	clTruncate      = tw.New().Truncate()

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

	// Progress. Determinate fill widths are owned by the renderer stylesheet's
	// bounded data-progress-percent rules, keeping delivery standalone without
	// inline styles or host-private CSS.
	clProgressRoot = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).
			Width(tw.SFull).Gap(tw.S2)
	clProgressHeader = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Justify(tw.JustifyBetween).Gap(tw.S3)
	clProgressLabel = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).
			Tracking(tw.TrackingTight).TextColor(tw.FgPrimary)
	clProgressValue = tw.New().FontSize(tw.TextSM).TextColor(tw.FgSecondary).
			TabularNums()
	clProgressTrack = tw.New().Width(tw.SFull).Overflow(tw.OverflowHidden).
			Rounded(tw.RadiusFull).Bg(tw.SurfaceTertiary)
	clProgressTrackSize = map[string]tw.ClassList{
		"sm": tw.New().Height(tw.S1_5),
		"md": tw.New().Height(tw.S2_5),
		"lg": tw.New().Height(tw.S3_5),
	}
	clProgressFill = tw.New().Height(tw.SFull).Rounded(tw.RadiusFull).
			Transition(tw.TransitionAll).Duration(tw.Duration300).Easing(tw.EaseOut)
	clProgressTone = map[string]tw.ClassList{
		"brand":   tw.New().Bg(tw.SurfaceBrand),
		"success": tw.New().Bg(tw.SurfaceSuccess),
		"warning": tw.New().Bg(tw.SurfaceWarning),
		"danger":  tw.New().Bg(tw.SurfaceDanger),
		"info":    tw.New().Bg(tw.SurfaceInfo),
	}
	clProgressIndeterminate = tw.New().Width(tw.SFull).AnimatePulse()

	// Tooltip. The controller only toggles the popup's hidden attribute and
	// aria-hidden state; geometry and visual treatment remain renderer-owned.
	clTooltipContainer = tw.New().
				Position(tw.PositionRelative).
				Display(tw.DisplayInlineBlock)
	clTooltipTrigger = tw.New().Display(tw.DisplayContents)
	clTooltipPopup   = tw.New().
				Position(tw.PositionAbsolute).
				ZLayer(tw.ZTooltip).
				PaddingX(tw.S2).
				PaddingY(tw.S1).
				FontSize(tw.TextSM).
				TextColor(tw.FgOnBrand).
				Bg(tw.SurfaceInverse).
				Rounded(tw.RadiusBase).
				Shadow(tw.ShadowLG).
				WhitespaceNowrap().
				PointerEvents(tw.PointerNone)
	clTooltipPosition = map[string]tw.ClassList{
		"top": tw.New().
			Bottom(tw.SFull).
			LeftOffset(tw.PositionHalf).
			TranslateX(tw.TranslateNegHalf).
			MarginBottom(tw.S2),
		"bottom": tw.New().
			Top(tw.SFull).
			LeftOffset(tw.PositionHalf).
			TranslateX(tw.TranslateNegHalf).
			MarginTop(tw.S2),
		"left": tw.New().
			Right(tw.SFull).
			TopOffset(tw.PositionHalf).
			TranslateY(tw.TranslateNegHalf).
			MarginRight(tw.S2),
		"right": tw.New().
			Left(tw.SFull).
			TopOffset(tw.PositionHalf).
			TranslateY(tw.TranslateNegHalf).
			MarginLeft(tw.S2),
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

	// Card. clCard remains the compiled default frame used by delivery and
	// skeleton projections; the smaller fragments let CardWithSlots own
	// section padding and variants without a private downstream style stack.
	clCardFrame = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S2).
			Rounded(tw.RadiusLG).Bg(tw.SurfacePrimary).Overflow(tw.OverflowHidden)
	clCardSectioned    = tw.New().Gap(tw.S0)
	clCardBorder       = tw.New().Border(tw.Border1).BorderColor(tw.BorderPrimary)
	clCardPadNone      = tw.New().Padding(tw.S0)
	clCardPadSmall     = tw.New().Padding(tw.S3)
	clCardPadDefault   = tw.New().Padding(tw.S5)
	clCardPadMedium    = tw.New().Padding(tw.S6)
	clCardPadLarge     = tw.New().Padding(tw.S8)
	clCardShadowSmall  = tw.New().Shadow(tw.ShadowSM)
	clCardShadowMedium = tw.New().Shadow(tw.ShadowBase)
	clCardShadowLarge  = tw.New().Shadow(tw.ShadowLG)
	clCard             = clCardFrame.Merge(clCardBorder).Merge(clCardPadDefault).Merge(clCardShadowSmall)
	clCardClickable    = tw.New().Cursor(tw.CursorPointer).Transition(tw.TransitionShadow).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Shadow(tw.ShadowMD) }).
				Merge(clFocusRing)
	clCardHoverable = tw.New().Transition(tw.TransitionShadow).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Shadow(tw.ShadowLG) })
	clCardTitle = tw.New().FontFamily(tw.FontSerif).FontSize(tw.TextLG).
			FontWeight(tw.FontSemibold).TextColor(tw.FgPrimary)
	clCardDesc            = tw.New().FontSize(tw.TextSM).TextColor(tw.FgMuted)
	clCardHeader          = tw.New().BorderBottom(tw.Border1).BorderColor(tw.BorderPrimary)
	clCardFooter          = tw.New().BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).Bg(tw.SurfaceSecondary)
	clCardImageVertical   = tw.New().Width(tw.SFull).ObjectCover()
	clCardImageHorizontal = tw.New().Width(tw.S48).ObjectCover()
	clCardHorizontal      = tw.New().Display(tw.DisplayFlex)
	clCardVertical        = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Flex1()

	// Breadcrumb.
	clBreadcrumb = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).FontSize(tw.TextSM).
			ListStyle("none").Margin(tw.S0).Padding(tw.S0)
	clBreadcrumbSep = tw.New().TextColor(tw.FgTertiary)
	clBreadcrumbCur = tw.New().TextColor(tw.FgPrimary).FontWeight(tw.FontMedium)

	// Accordion.
	clAccordionRoot     = tw.New().Width(tw.SFull).Overflow(tw.OverflowHidden)
	clAccordionBordered = tw.New().Border(tw.Border1).BorderColor(tw.BorderPrimary).
				Rounded(tw.RadiusLG).Bg(tw.SurfacePrimary)
	clAccordionUnbordered = tw.New().Rounded(tw.RadiusLG).Bg(tw.SurfacePrimary)
	clAccordionSeparator  = tw.New().BorderTop(tw.Border1).BorderColor(tw.BorderPrimary)
	clAccordionTrigger    = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).
				Items(tw.ItemsCenter).Justify(tw.JustifyBetween).Gap(tw.S3).
				PaddingX(tw.S4).PaddingY(tw.S3).TextAlign(tw.TextLeft).
				Transition(tw.TransitionColors).
				On(tw.StateHover, hoverBg(tw.SurfaceSecondary)).Merge(clFocusRing)
	clAccordionTriggerDisabled = tw.New().Opacity(tw.Opacity50).Cursor(tw.CursorNotAllowed)
	clAccordionLead            = tw.New().Display(tw.DisplayFlex).MinWidth(tw.S0).Flex1().Items(tw.ItemsStart).Gap(tw.S3)
	clAccordionItemIcon        = tw.New().MarginTop(tw.S0_5).Display(tw.DisplayFlex).Height(tw.S5).Width(tw.S5).
					FlexShrink0().Items(tw.ItemsCenter).Justify(tw.JustifyCenter).TextColor(tw.FgTertiary)
	clAccordionTitleBlock = tw.New().Display(tw.DisplayFlex).MinWidth(tw.S0).Flex1().FlexDir(tw.FlexCol).Gap(tw.S1)
	clAccordionTitle      = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).TextColor(tw.FgPrimary)
	clAccordionSubtitle   = tw.New().FontSize(tw.TextSM).TextColor(tw.FgTertiary)
	clAccordionChevron    = tw.New().Display(tw.DisplayFlex).FlexShrink0().TextColor(tw.FgTertiary).
				Transition(tw.TransitionTransform)
	clAccordionChevronOpen = tw.New().Rotate("180")
	clAccordionPanel       = tw.New().PaddingX(tw.S4).PaddingBottom(tw.S4).FontSize(tw.TextSM).TextColor(tw.FgSecondary)

	// Stepper.
	clStepperListHorizontal = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter)
	clStepperListVertical   = tw.New().SpaceY(tw.S4)
	clStepperItemHorizontal = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Flex1()
	clStepperItemLast       = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter)
	clStepperRowHorizontal  = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter)
	clStepperRowVertical    = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart)
	clStepperVerticalItem   = tw.New().Position(tw.PositionRelative)
	clStepperVerticalIcons  = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Items(tw.ItemsCenter)
	clStepperVerticalText   = tw.New().MarginLeft(tw.S4).MinWidth(tw.S0)
	clStepperIndicator      = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Justify(tw.JustifyCenter).Rounded(tw.RadiusFull).FlexShrink0().
				Transition(tw.TransitionColors)
	clStepperIndicatorRegular = tw.New().Height(tw.S8).Width(tw.S8)
	clStepperIndicatorCompact = tw.New().Height(tw.S6).Width(tw.S6)
	clStepperClickable        = tw.New().Cursor(tw.CursorPointer).
					On(tw.StateHover, func(c tw.ClassList) tw.ClassList { return c.Opacity(tw.Opacity80) }).
					Merge(clFocusRing)
	clStepperDisabled = tw.New().Opacity(tw.Opacity50).Cursor(tw.CursorNotAllowed)
	clStepperGlyph    = tw.New().FontSize(tw.TextXS).FontWeight(tw.FontMedium)
	clStepperLabel    = tw.New().MarginLeft(tw.S3).FontSize(tw.TextSM).
				FontWeight(tw.FontMedium).WhitespaceNowrap()
	clStepperLabelCompact = tw.New().MarginLeft(tw.S2).FontSize(tw.TextXS).
				FontWeight(tw.FontMedium).WhitespaceNowrap()
	clStepperLabelBlock          = tw.New().Display(tw.DisplayBlock)
	clStepperDescription         = tw.New().FontSize(tw.TextXS).FontWeight(tw.FontNormal)
	clStepperConnectorHorizontal = tw.New().Flex1().MarginX(tw.S4).Height(tw.S0_5).
					Transition(tw.TransitionColors)
	clStepperConnectorVertical = tw.New().Width(tw.S0_5).Height(tw.S8).MarginTop(tw.S2).
					Transition(tw.TransitionColors)
	clStepperIndicatorState = map[string]tw.ClassList{
		"completed": tw.New().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand),
		"active": tw.New().Border(tw.Border2).BorderColor(tw.BorderBrand).
			Bg(tw.SurfacePrimary).TextColor(tw.FgBrand),
		"error": tw.New().Bg(tw.SurfaceDanger).TextColor(tw.FgOnBrand),
		"pending": tw.New().Border(tw.Border2).BorderColor(tw.BorderSecondary).
			Bg(tw.SurfacePrimary).TextColor(tw.FgTertiary),
	}
	clStepperLabelState = map[string]tw.ClassList{
		"completed": tw.New().TextColor(tw.FgBrand),
		"active":    tw.New().TextColor(tw.FgBrand),
		"error":     tw.New().TextColor(tw.FgDanger),
		"pending":   tw.New().TextColor(tw.FgTertiary),
	}
	clStepperDescriptionState = map[string]tw.ClassList{
		"completed": tw.New().TextColor(tw.FgTertiary),
		"active":    tw.New().TextColor(tw.FgSecondary),
		"error":     tw.New().TextColor(tw.FgDanger),
		"pending":   tw.New().TextColor(tw.FgTertiary),
	}
	clStepperConnectorState = map[string]tw.ClassList{
		"completed": tw.New().Bg(tw.SurfaceBrand),
		"pending":   tw.New().Bg(tw.BorderSecondary),
	}

	// Sidebar.
	clSidebarRootAdmin = tw.New().Display(tw.DisplayHidden).
				Breakpoint(tw.BreakpointLG, func(c tw.ClassList) tw.ClassList {
			return c.Display(tw.DisplayFlex).FlexShrink0()
		}).Transition(tw.TransitionAll)
	clSidebarRootContent = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).
				Transition(tw.TransitionAll)
	clSidebarWidthCollapsed = tw.New().Breakpoint(tw.BreakpointLG, func(c tw.ClassList) tw.ClassList {
		return c.Width(tw.S16)
	})
	clSidebarWidthExpanded = tw.New().Breakpoint(tw.BreakpointLG, func(c tw.ClassList) tw.ClassList {
		return c.Width(tw.S64)
	})
	clSidebarDisabled    = tw.New().Opacity(tw.Opacity50)
	clSidebarInner       = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Height(tw.SFull)
	clSidebarColumnAdmin = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Flex1().
				Bg(tw.SurfaceInverse).PaddingTop(tw.S5).PaddingBottom(tw.S4).OverflowY(tw.OverflowAuto)
	clSidebarColumnContent = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Flex1().
				Bg(tw.SurfacePrimary).OverflowY(tw.OverflowVisible)
	clSidebarBrandAdmin = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).FlexShrink0().
				PaddingX(tw.S4).MarginBottom(tw.S8)
	clSidebarBrandContent = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S3).
				FlexShrink0().MarginBottom(tw.S4)
	clSidebarBrandLink      = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Merge(clFocusRing)
	clSidebarBrandText      = tw.New().FontSize(tw.TextXL).FontWeight(tw.FontBold).TextColor(tw.FgOnInverse)
	clSidebarNavWrapAdmin   = tw.New().MarginTop(tw.S5).Flex1().Display(tw.DisplayFlex).FlexDir(tw.FlexCol)
	clSidebarNavWrapContent = tw.New().Flex1().Display(tw.DisplayFlex).FlexDir(tw.FlexCol)
	clSidebarNavAdmin       = tw.New().Flex1().PaddingX(tw.S2).SpaceY(tw.S1)
	clSidebarNavContent     = tw.New().Flex1().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S4)
	clSidebarLinkAdmin      = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S3).
				FontSize(tw.TextSM).FontWeight(tw.FontMedium).Rounded(tw.RadiusMD).
				Transition(tw.TransitionColors).Merge(clFocusRing)
	clSidebarLinkContent = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsStart).Gap(tw.S2).
				FontSize(tw.TextSM).FontWeight(tw.FontMedium).Rounded(tw.RadiusMD).
				Transition(tw.TransitionColors).Merge(clFocusRing)
	clSidebarLinkPadExpanded  = tw.New().PaddingX(tw.S2).PaddingY(tw.S2)
	clSidebarLinkPadCollapsed = tw.New().PaddingX(tw.S2).PaddingY(tw.S2).
					Justify(tw.JustifyCenter)
	clSidebarLinkActiveAdmin = tw.New().Bg(tw.SurfaceOverlay).TextColor(tw.FgOnInverse)
	clSidebarLinkIdleAdmin   = tw.New().TextColor(tw.FgOnInverse).
					On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceOverlay).TextColor(tw.FgOnInverse)
		})
	clSidebarLinkActiveContent = tw.New().Bg(tw.SurfaceBrandSoft).TextColor(tw.FgBrand)
	clSidebarLinkIdleContent   = tw.New().TextColor(tw.FgSecondary).
					On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceSecondary).TextColor(tw.FgPrimary)
		})
	clSidebarItemDisabled = tw.New().Opacity(tw.Opacity50).PointerEvents(tw.PointerNone).
				Cursor(tw.CursorNotAllowed)
	clSidebarSection       = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S2)
	clSidebarSectionHeader = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).Gap(tw.S2).
				PaddingX(tw.S2).PaddingBottom(tw.S2).BorderBottom(tw.Border1).
				BorderColor(tw.BorderPrimary).FontSize(tw.TextXS).FontWeight(tw.FontBold).
				TextColor(tw.FgTertiary)
	clSidebarSectionGlyph = tw.New().TextColor(tw.FgBrand)
	clSidebarSectionList  = tw.New().Display(tw.DisplayFlex).FlexDir(tw.FlexCol).Gap(tw.S0_5)
	clSidebarPrefixAdmin  = tw.New().MinWidth(tw.S10).FontSize(tw.TextXS).
				FontWeight(tw.FontSemibold).TextColor(tw.FgTertiary)
	clSidebarPrefixContent = tw.New().MinWidth(tw.S10).FontSize(tw.TextXS).
				FontWeight(tw.FontSemibold).TextColor(tw.FgSecondary)
	clSidebarLabelVisible  = tw.New().MinWidth(tw.S0).Flex1()
	clSidebarLabelHidden   = tw.New().MinWidth(tw.S0).Flex1().Display(tw.DisplayHidden)
	clSidebarLabelContent  = tw.New().MinWidth(tw.S0).Flex1().BreakWords()
	clSidebarNestedGroup   = tw.New().SpaceY(tw.S1)
	clSidebarNestedIndent  = tw.New().MarginLeft(tw.S4).SpaceY(tw.S1)
	clSidebarFooterAdmin   = tw.New().MarginTop(tw.SAuto).PaddingX(tw.S4).PaddingTop(tw.S4)
	clSidebarFooterContent = tw.New().MarginTop(tw.S4)

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
	clTabsRoot                    = tw.New().Display(tw.DisplayFlex).Width(tw.SFull).Gap(tw.S4)
	clTabsRootHorizontal          = tw.New().FlexDir(tw.FlexCol)
	clTabsRootVertical            = tw.New().FlexDir(tw.FlexRow)
	clTabsListBase                = tw.New().Display(tw.DisplayFlex).Gap(tw.S1)
	clTabsListHorizontal          = tw.New().FlexDir(tw.FlexRow)
	clTabsListVertical            = tw.New().FlexDir(tw.FlexCol)
	clTabsListUnderlineHorizontal = tw.New().BorderBottom(tw.Border1).
					BorderColor(tw.BorderPrimary)
	clTabsListUnderlineVertical = tw.New().BorderRight(tw.Border1).
					BorderColor(tw.BorderPrimary).PaddingRight(tw.S4)
	clTabsButtonBase = tw.New().Display(tw.DisplayInlineFlex).Items(tw.ItemsCenter).
				PaddingX(tw.S3).PaddingY(tw.S2).FontSize(tw.TextSM).
				FontWeight(tw.FontMedium).Transition(tw.TransitionColors).Merge(clFocusRing)
	clTabsButtonPills               = tw.New().Rounded(tw.RadiusMD)
	clTabsButtonUnderlineHorizontal = tw.New().BorderBottom(tw.Border2).NegMargin("b", tw.SPX)
	clTabsButtonUnderlineVertical   = tw.New().BorderRight(tw.Border2).NegMargin("r", tw.SPX)
	clTabsPillsActive               = tw.New().Bg(tw.SurfaceBrand).TextColor(tw.FgOnBrand)
	clTabsPillsIdle                 = tw.New().TextColor(tw.FgSecondary).
					On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.TextColor(tw.FgPrimary).Bg(tw.SurfaceSecondary)
		})
	clTabsUnderlineActive = tw.New().BorderColor(tw.BorderBrand).TextColor(tw.FgBrand)
	clTabsUnderlineIdle   = tw.New().BorderColor(tw.ColorTransparent).TextColor(tw.FgSecondary).
				On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.TextColor(tw.FgPrimary).BorderColor(tw.BorderSecondary)
		})
	clTabsDisabled = tw.New().Opacity(tw.Opacity50).Cursor(tw.CursorNotAllowed)
	clTabsIcon     = tw.New().MarginRight(tw.S2)
	clTabsBadge    = tw.New().MarginLeft(tw.S2).Display(tw.DisplayInlineFlex).
			Items(tw.ItemsCenter).Rounded(tw.RadiusFull).Bg(tw.SurfaceTertiary).
			PaddingX(tw.S2).PaddingY(tw.S0_5).FontSize(tw.TextXS).
			FontWeight(tw.FontMedium).TextColor(tw.FgSecondary)
	clTabsPanels = tw.New().MinWidth(tw.S0).Flex1()
	clTabsPanel  = tw.New().MinWidth(tw.S0)
	clTabsLazy   = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
			Justify(tw.JustifyCenter).PaddingY(tw.S8)
	clTabsLazyLabel = tw.New().TextColor(tw.FgSecondary)

	clDashboardWidget = tw.New().Bg(tw.SurfacePrimary).Rounded(tw.RadiusLG).
				Border(tw.Border1).BorderColor(tw.BorderPrimary).Shadow(tw.ShadowSM).
				Overflow(tw.OverflowHidden)
	clDashboardWidgetSpan = map[string]tw.ClassList{
		"1": tw.New().ColSpan(1),
		"2": tw.New().ColSpan(2),
		"3": tw.New().ColSpan(3),
		"4": tw.New().ColSpan(4),
	}
	clDashboardWidgetHeader = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Gap(tw.S3).PaddingX(tw.S4).PaddingTop(tw.S4).PaddingBottom(tw.S2)
	clDashboardWidgetIcon = tw.New().Display(tw.DisplayFlex).FlexShrink0().Padding(tw.S2).
				Bg(tw.SurfaceBrandSoft).Rounded(tw.RadiusLG)
	clDashboardWidgetCopy  = tw.New().MinWidth(tw.S0).Flex1()
	clDashboardWidgetTitle = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium).
				TextColor(tw.FgPrimary).Truncate()
	clDashboardWidgetSubtitle = tw.New().FontSize(tw.TextXS).
					TextColor(tw.FgSecondary).MarginTop(tw.S0_5)
	clDashboardWidgetBody  = tw.New().PaddingX(tw.S4).PaddingBottom(tw.S4)
	clDashboardWidgetValue = tw.New().FontSize(tw.Text2XL).FontWeight(tw.FontBold).
				TextColor(tw.FgPrimary)
	clDashboardWidgetTrend = tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter).
				Gap(tw.S1).MarginTop(tw.S1)
	clDashboardWidgetChange    = tw.New().FontSize(tw.TextSM).FontWeight(tw.FontMedium)
	clDashboardWidgetTrendTone = map[string]tw.ClassList{
		"up":   tw.New().TextColor(tw.FgSuccess),
		"down": tw.New().TextColor(tw.FgDanger),
		"flat": tw.New().TextColor(tw.FgTertiary),
	}
	clDashboardWidgetPrevious = tw.New().FontSize(tw.TextXS).TextColor(tw.FgSecondary)
	clDashboardWidgetFooter   = tw.New().PaddingX(tw.S4).PaddingY(tw.S3).
					Bg(tw.SurfaceSecondary).BorderTop(tw.Border1).BorderColor(tw.BorderPrimary).
					FontSize(tw.TextSM).TextColor(tw.FgTertiary)
)

// ClassLists returns every ClassList the renderers compose, base and variant
// alike. Applications derive their stylesheet from it:
//
//	sheet, err := emission.For(web.ClassLists()...)
func ClassLists() []tw.ClassList {
	out := []tw.ClassList{
		clIcon, clAvatarBase, clAvatarInitials, clAvatarImage, clAvatarStatus,
		clFocusRing, clButtonBase, clButtonFull, clButtonIconOnly, clButtonDisabledLink,
		clBadgeBase, clBadgeDot, clBadgeCount, clBadgeRemove,
		clAlertBase, clAlertRegular, clAlertCompact, clAlertBordered,
		clAlertTitle, clAlertMessage, clAlertBody, clAlertIcon, clAlertActions, clAlertClose,
		clToastBase, clToastIcon, clToastBody, clToastTitle, clToastMessage, clToastLead, clToastClose,
		clFieldWrap, clFieldWrapFull, clLabel, clHelp, clFieldErr, clRequired,
		clInput, clInputNormal, clInputError, clInputReadOnly,
		clInputIconWrap, clInputIconStart, clInputIconEnd, clInputPadStart, clInputPadEnd,
		clAutocompleteControl, clAutocompletePanel, clAutocompleteList,
		clAutocompleteOption, clAutocompleteOptionActive, clAutocompleteIndicator, clAutocompleteSpinner,
		clFileUploadRoot, clFileUploadDropZone, clFileUploadDropZoneDisabled,
		clFileUploadDropZoneInner, clFileUploadInputHidden, clFileUploadInput,
		clFileUploadIcon, clFileUploadPrompt, clFileUploadPromptText,
		clFileUploadPromptAction, clFileUploadHint, clFileUploadLoading,
		clFileUploadLoadingIcon, clFileUploadList, clFileUploadError,
		clFileUploadItem, clFileUploadItemMeta, clFileUploadPreview,
		clFileUploadPreviewLarge, clFileUploadItemCopy, clFileUploadItemName,
		clFileUploadItemSize, clFileUploadRemove, clFileUploadDragActive,
		clFileUploadDragIdle,
		clDropdownRoot, clDropdownTrigger, clDropdownButton, clDropdownTriggerLabel,
		clDropdownTriggerActions, clDropdownIconButton, clDropdownChevron,
		clDropdownPanel, clDropdownSearchWrap, clDropdownSearch, clDropdownOptions,
		clDropdownOption, clDropdownOptionSelected, clDropdownOptionDisabled,
		clDropdownOptionMark, clDropdownOptionSpacer, clDropdownOptionIcon,
		clDropdownOptionLabel, clDropdownGroupLabel,
		clActionMenuRoot, clActionMenuTrigger, clActionMenuPanel, clActionMenuPanelInner,
		clActionMenuSeparator, clActionMenuSectionLabel, clActionMenuItem,
		clActionMenuItemDisabled, clActionMenuItemIcon, clActionMenuItemLabel,
		clDrawerRoot, clDrawerOverlay, clDrawerPanel, clDrawerHeader,
		clDrawerTitleBlock, clDrawerTitle, clDrawerDescription, clDrawerBody,
		clDrawerFooter, clDrawerClose,
		clModalRoot, clModalCentered, clModalBottomSheet, clModalOverlay,
		clModalPanel, clModalHeader, clModalTitleBlock, clModalTitle,
		clModalDescription, clModalBody, clModalFooter, clModalClose, clModalCancel,
		clTextareaManual, clTextareaAuto, clTextareaFull, clTextareaMeta,
		clTextareaSupporting, clTextareaCounter,
		clCheckbox, clCheckRow, clCheckboxRoot, clCheckboxRootDisabled,
		clCheckboxInput, clCheckboxIndicator, clCheckboxIndicatorIdle,
		clCheckboxIndicatorActive, clCheckboxCheckmark, clCheckboxBar, clCheckboxLabel,
		clRadioRoot, clRadioRootDisabled, clRadioInput, clRadioInputDisabled, clRadioLabel, clRadioDot,
		clChoiceGroupRoot, clChoiceGroupLegend, clChoiceGroupDescription,
		clChoiceGroupOptions, clChoiceGroupVertical, clChoiceGroupHorizontal, clChoiceGroupError,
		clSliderRoot, clSliderLabel, clSliderRow, clSliderInput, clSliderInputDisabled, clSliderValue,
		clToggleRoot, clToggleRootDisabled, clToggleControl, clToggleInput, clToggleTrack,
		clToggleKnob, clToggleKnobUnchecked, clToggleLabel,
		clSearchWrap, clSearchInput, clSrOnly,
		clTruncate, clHeadingBase,
		clDividerH, clDividerV, clDividerText, clDividerTextLine, clDividerTextLabel,
		clSpinner, clProgressRoot, clProgressHeader, clProgressLabel, clProgressValue,
		clProgressTrack, clProgressFill, clProgressIndeterminate,
		clTooltipContainer, clTooltipTrigger, clTooltipPopup,
		clEmpty, clEmptyPad, clEmptyBordered, clEmptyCompact, clEmptyTitle, clEmptyDesc,
		clSkeleton, clSkeletonText, clSkeletonLine, clSkeletonLineLast,
		clKbd, clLink, clTagBase, clTagIdle, clTagSelected,
		clTextItalic, clTextUnderline, clTextNoWrap, clTruncate,
		clStack, clFlex, clGrid, clContainer, clDataManagementPage,
		clTableWrap, clTable, clTableHead, clTableThBase, clTableTh, clTableTd, clTableRow, clTableTdC,
		clTableThSort, clTableSortBtn, clTableRowAlt, clTableTdStrong, clTableActions, clTableCellNote,
		clCard, clCardFrame, clCardSectioned, clCardBorder,
		clCardPadNone, clCardPadSmall, clCardPadDefault, clCardPadMedium, clCardPadLarge,
		clCardShadowSmall, clCardShadowMedium, clCardShadowLarge,
		clCardClickable, clCardHoverable, clCardTitle, clCardDesc,
		clCardHeader, clCardFooter, clCardImageVertical, clCardImageHorizontal,
		clCardHorizontal, clCardVertical,
		clBreadcrumb, clBreadcrumbSep, clBreadcrumbCur,
		clAccordionRoot, clAccordionBordered, clAccordionUnbordered,
		clAccordionSeparator, clAccordionTrigger,
		clAccordionTriggerDisabled, clAccordionLead, clAccordionItemIcon,
		clAccordionTitleBlock, clAccordionTitle, clAccordionSubtitle,
		clAccordionChevron, clAccordionChevronOpen, clAccordionPanel,
		clStepperListHorizontal, clStepperListVertical,
		clStepperItemHorizontal, clStepperItemLast,
		clStepperRowHorizontal, clStepperRowVertical,
		clStepperVerticalItem, clStepperVerticalIcons, clStepperVerticalText,
		clStepperIndicator, clStepperIndicatorRegular, clStepperIndicatorCompact,
		clStepperClickable, clStepperDisabled, clStepperGlyph,
		clStepperLabel, clStepperLabelCompact, clStepperLabelBlock,
		clStepperDescription, clStepperConnectorHorizontal, clStepperConnectorVertical,
		clSidebarRootAdmin, clSidebarRootContent, clSidebarWidthCollapsed,
		clSidebarWidthExpanded, clSidebarDisabled, clSidebarInner,
		clSidebarColumnAdmin, clSidebarColumnContent,
		clSidebarBrandAdmin, clSidebarBrandContent, clSidebarBrandLink, clSidebarBrandText,
		clSidebarNavWrapAdmin, clSidebarNavWrapContent, clSidebarNavAdmin, clSidebarNavContent,
		clSidebarLinkAdmin, clSidebarLinkContent,
		clSidebarLinkPadExpanded, clSidebarLinkPadCollapsed,
		clSidebarLinkActiveAdmin, clSidebarLinkIdleAdmin,
		clSidebarLinkActiveContent, clSidebarLinkIdleContent,
		clSidebarItemDisabled, clSidebarSection, clSidebarSectionHeader,
		clSidebarSectionGlyph, clSidebarSectionList,
		clSidebarPrefixAdmin, clSidebarPrefixContent,
		clSidebarLabelVisible, clSidebarLabelHidden, clSidebarLabelContent,
		clSidebarNestedGroup, clSidebarNestedIndent,
		clSidebarFooterAdmin, clSidebarFooterContent,
		clGridSection, clGridToolbar, clGridActions,
		clWindow, clWindowItems, clWindowState, clWindowLoading, clWindowEmpty,
		clWindowError, clWindowTitle, clWindowDescription, clWindowRetry,
		clWindowFooter, clWindowNavigationError,
		clPagination, clPageBtn, clPageIdle, clPageCur, clPageLabel,
		clTabsRoot, clTabsRootHorizontal, clTabsRootVertical,
		clTabsListBase, clTabsListHorizontal, clTabsListVertical,
		clTabsListUnderlineHorizontal, clTabsListUnderlineVertical,
		clTabsButtonBase, clTabsButtonPills, clTabsButtonUnderlineHorizontal,
		clTabsButtonUnderlineVertical, clTabsPillsActive, clTabsPillsIdle,
		clTabsUnderlineActive, clTabsUnderlineIdle, clTabsDisabled, clTabsIcon,
		clTabsBadge, clTabsPanels, clTabsPanel, clTabsLazy, clTabsLazyLabel,
		clDashboardWidget, clDashboardWidgetHeader, clDashboardWidgetIcon,
		clDashboardWidgetCopy, clDashboardWidgetTitle, clDashboardWidgetSubtitle,
		clDashboardWidgetBody, clDashboardWidgetValue, clDashboardWidgetTrend,
		clDashboardWidgetChange, clDashboardWidgetPrevious, clDashboardWidgetFooter,
	}
	for _, m := range []map[string]tw.ClassList{
		clButtonVariant, clButtonTone, clButtonSize,
		clBadgeVariant, clBadgeTone, clBadgeSize, clBadgeDotTone, clAlertVariant, clToastTone,
		clIconSize, clIconTone, clAvatarSize, clAvatarShape, clAvatarTone,
		clAvatarInitialsSize, clAvatarStatusTone, clAvatarStatusSize, clAvatarStatusPosition,
		clInputSize, clInputTone, clDropdownTriggerSize, clDropdownButtonSize,
		clSliderTone, clToggleTrackSize, clToggleTrackState, clToggleKnobSize, clToggleKnobChecked,
		clSpinnerSize, clSpinnerTone, clProgressTrackSize, clProgressTone,
		clTooltipPosition, clKbdSize, clTagTone,
		clSkeletonBlockSize, clSkeletonLineSize, clSkeletonCircleSize,
		clTextTransform, clStepperIndicatorState, clStepperLabelState,
		clStepperDescriptionState, clStepperConnectorState,
		clActionMenuAlign, clActionMenuWidth, clActionMenuItemTone,
		clDrawerPosition, clDrawerWidth, clDrawerBottomSize,
		clModalPanelSize,
		clDashboardWidgetSpan, clDashboardWidgetTrendTone,
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
	for _, align := range clTextAlign {
		out = append(out, tw.New().TextAlign(align))
	}
	for lines := 1; lines <= 6; lines++ {
		out = append(out, tw.New().LineClamp(lines))
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
