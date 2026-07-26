// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package component

import "slices"

// FoundationOwner is the public owner of the renderer-neutral vocabulary.
const FoundationOwner = "septagon-oss/pk-ui"

// Foundation component wire IDs. This is the portable vocabulary every
// PlatformKit renderer may implement. Extensions are not required to modify
// this list; they contribute canonical namespaced IDs through Contribution.
const (
	TypeText    = "Text"
	TypeButton  = "Button"
	TypeInput   = "Input"
	TypeBadge   = "Badge"
	TypeIcon    = "Icon"
	TypeImage   = "Image"
	TypeLink    = "Link"
	TypeAvatar  = "Avatar"
	TypeDivider = "Divider"
	TypeSpinner = "Spinner"
	TypeToggle  = "Toggle"

	TypeFormField  = "FormField"
	TypeCard       = "Card"
	TypeAlert      = "Alert"
	TypeTabs       = "Tabs"
	TypeModal      = "Modal"
	TypeTooltip    = "Tooltip"
	TypeDropdown   = "Dropdown"
	TypeBreadcrumb = "Breadcrumb"
	TypeSection    = "Section"
	TypeMetric     = "Metric"
	TypeActionCard = "ActionCard"
	TypeChipRow    = "ChipRow"
	TypeDetailList = "DetailList"
	TypeListItem   = "ListItem"

	TypeDataTable        = "DataTable"
	TypeSearchInput      = "SearchInput"
	TypePagination       = "Pagination"
	TypeActionBar        = "ActionBar"
	TypeFormActions      = "FormActions"
	TypeDetailField      = "DetailField"
	TypeDetailSection    = "DetailSection"
	TypeDetailActions    = "DetailActions"
	TypeNavigation       = "Navigation"
	TypeMap              = "Map"
	TypeMovementTimeline = "MovementTimeline"
	TypeHero             = "Hero"
	TypeEmptyState       = "EmptyState"

	TypeContainer = "Container"
	TypeGrid      = "Grid"
	TypeStack     = "Stack"
	TypeFlex      = "Flex"
)

var foundationCatalog = []Descriptor{
	foundation(TypeText, TierAtom),
	foundation(TypeButton, TierAtom),
	foundation(TypeInput, TierAtom),
	foundation(TypeBadge, TierAtom),
	foundation(TypeIcon, TierAtom),
	foundation(TypeImage, TierAtom),
	foundation(TypeLink, TierAtom),
	foundation(TypeAvatar, TierAtom),
	foundation(TypeDivider, TierAtom),
	foundation(TypeSpinner, TierAtom),
	foundation(TypeToggle, TierAtom),

	foundation(TypeFormField, TierMolecule),
	foundation(TypeCard, TierMolecule),
	foundation(TypeAlert, TierMolecule),
	foundation(TypeTabs, TierMolecule),
	foundation(TypeModal, TierMolecule),
	foundation(TypeTooltip, TierMolecule),
	foundation(TypeDropdown, TierMolecule),
	foundation(TypeBreadcrumb, TierMolecule),
	foundation(TypeSection, TierMolecule),
	foundation(TypeMetric, TierMolecule),
	foundation(TypeActionCard, TierMolecule),
	foundation(TypeChipRow, TierMolecule),
	foundation(TypeDetailList, TierMolecule),
	foundation(TypeListItem, TierMolecule),

	foundation(TypeDataTable, TierOrganism),
	foundation(TypeSearchInput, TierOrganism),
	foundation(TypePagination, TierOrganism),
	foundation(TypeActionBar, TierOrganism),
	foundation(TypeFormActions, TierOrganism),
	foundation(TypeDetailField, TierOrganism),
	foundation(TypeDetailSection, TierOrganism),
	foundation(TypeDetailActions, TierOrganism),
	foundation(TypeNavigation, TierOrganism),
	foundation(TypeMap, TierOrganism),
	foundation(TypeMovementTimeline, TierOrganism),
	foundation(TypeHero, TierOrganism),
	foundation(TypeEmptyState, TierOrganism),

	foundation(TypeContainer, TierTemplate),
	foundation(TypeGrid, TierTemplate),
	foundation(TypeStack, TierTemplate),
	foundation(TypeFlex, TierTemplate),
}

// FoundationCatalog returns a defensive copy of the portable component
// vocabulary in stable declaration order.
func FoundationCatalog() []Descriptor {
	return slices.Clone(foundationCatalog)
}

// FoundationDescriptor returns the descriptor for a portable component ID.
// False means the ID is an extension, not that it is invalid.
func FoundationDescriptor(id string) (Descriptor, bool) {
	for _, descriptor := range foundationCatalog {
		if string(descriptor.ID) == id {
			return descriptor, true
		}
	}
	return Descriptor{}, false
}

func foundation(id string, tier Tier) Descriptor {
	return Descriptor{
		ID:      ID(id),
		Tier:    tier,
		Owner:   FoundationOwner,
		Version: "1.0.0",
	}
}
