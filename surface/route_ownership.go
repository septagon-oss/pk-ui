// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package surface

// route_ownership.go — deterministic concrete-path ownership for surface routes.
//

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

const maxRouteHierarchyDepth = 128

type routeOwnerCandidate struct {
	route      Route
	index      int
	depth      int
	stateAlias bool
}

// RouteOwnersForTarget projects a contribution's navigation graph onto one
// concrete route owner per canonical pathname for target. A contribution may
// legitimately contain multiple links to one page: a clickable parent and its
// default child, or query/fragment variants that select page state. Those links
// remain part of the contribution and its navigation projection, but they must
// not become ambiguous concrete route registrations.
//
// Ownership is deterministic and independent of map iteration:
//   - a query/fragment-free route owns the pathname ahead of a state alias;
//   - the deepest child owns an otherwise identical pathname ahead of an
//     ancestor grouping route;
//
// The returned routes retain declaration order among the selected owners. The
// input contribution is never mutated. Invalid contributions fail closed;
// callers must not register or render a partial projection when err is non-nil.
func RouteOwnersForTarget(contribution Contribution, target Target) ([]Route, error) {
	if err := ValidateContribution(contribution); err != nil {
		return nil, err
	}

	parents := make(map[string]string, len(contribution.Routes))
	for _, route := range contribution.Routes {
		routeID := strings.TrimSpace(route.ID)
		if routeID == "" {
			continue
		}
		parents[routeID] = strings.TrimSpace(route.ParentRouteID)
	}

	owners := make(map[string]routeOwnerCandidate, len(contribution.Routes))
	for index, route := range contribution.Routes {
		if !route.IncludesTarget(target) {
			continue
		}
		path, stateAlias := routeOwnershipPath(route.Path)
		if path == "" || strings.TrimSpace(route.ID) == "" {
			continue
		}

		candidate := routeOwnerCandidate{
			route:      route,
			index:      index,
			depth:      routeHierarchyDepth(route.ID, parents),
			stateAlias: stateAlias,
		}
		current, exists := owners[path]
		if !exists || candidateOwnsRoute(candidate, current) {
			owners[path] = candidate
		}
	}

	selected := make([]routeOwnerCandidate, 0, len(owners))
	for _, owner := range owners {
		selected = append(selected, owner)
	}
	sort.SliceStable(selected, func(i, j int) bool {
		return selected[i].index < selected[j].index
	})

	routes := make([]Route, 0, len(selected))
	for _, owner := range selected {
		routes = append(routes, cloneOwnedRoute(owner.route))
	}
	return routes, nil
}

// ValidateRouteOwnership rejects ambiguous concrete-path declarations that a
// deterministic projection must not silently hide. Exact-path declarations may
// form one ancestor/default-child chain; the deepest child is the owner. Query
// or fragment aliases must point at a query-free owner and carry the same
// capability set. Exact sibling declarations are ambiguous even when their
// permissions happen to match.
func ValidateRouteOwnership(contribution Contribution) error {
	if err := validateRouteHierarchy(contribution, false); err != nil {
		return err
	}
	for _, route := range contribution.Routes {
		if _, err := parseRouteLocation(route.Path); err != nil {
			return fmt.Errorf(
				"surface contribution invalid for module %q route %q: %w",
				strings.TrimSpace(contribution.ModuleID), strings.TrimSpace(route.ID), err,
			)
		}
	}

	targets := make(map[Target]struct{})
	for _, route := range contribution.Routes {
		for _, target := range route.Targets {
			targets[target] = struct{}{}
		}
	}

	orderedTargets := make([]Target, 0, len(targets))
	for target := range targets {
		orderedTargets = append(orderedTargets, target)
	}
	slices.Sort(orderedTargets)

	for _, target := range orderedTargets {
		groups := make(map[string][]routeOwnerCandidate)
		parents := make(map[string]string, len(contribution.Routes))
		for _, route := range contribution.Routes {
			routeID := strings.TrimSpace(route.ID)
			if routeID != "" {
				parents[routeID] = strings.TrimSpace(route.ParentRouteID)
			}
		}

		for index, route := range contribution.Routes {
			if !route.IncludesTarget(target) {
				continue
			}
			path, stateAlias := routeOwnershipPath(route.Path)
			if path == "" || strings.TrimSpace(route.ID) == "" {
				continue
			}
			groups[path] = append(groups[path], routeOwnerCandidate{
				route:      route,
				index:      index,
				depth:      routeHierarchyDepth(route.ID, parents),
				stateAlias: stateAlias,
			})
		}

		orderedPaths := make([]string, 0, len(groups))
		for path := range groups {
			orderedPaths = append(orderedPaths, path)
		}
		sort.Strings(orderedPaths)

		for _, path := range orderedPaths {
			candidates := groups[path]
			if len(candidates) == 1 {
				if candidates[0].stateAlias {
					return fmt.Errorf(
						"surface contribution invalid for module %q: target %q path %q: navigation-state aliases require one query-free concrete owner",
						strings.TrimSpace(contribution.ModuleID), target, path,
					)
				}
				continue
			}

			owner, err := validateRouteOwnerGroup(candidates, parents)
			if err != nil {
				return fmt.Errorf(
					"surface contribution invalid for module %q: target %q path %q: %w",
					strings.TrimSpace(contribution.ModuleID), target, path, err,
				)
			}
			if owner.route.Hidden {
				for _, candidate := range candidates {
					if !candidate.route.Hidden {
						return fmt.Errorf(
							"surface contribution invalid for module %q: target %q path %q: hidden owner %q shadows visible navigation route %q",
							strings.TrimSpace(contribution.ModuleID), target, path, owner.route.ID, candidate.route.ID,
						)
					}
				}
			}
			for _, candidate := range candidates {
				if candidate.route.ID == owner.route.ID {
					continue
				}
				if !candidate.stateAlias && isRouteAncestor(candidate.route.ID, owner.route.ID, parents) {
					// Ancestors are navigation groupings for this pathname; the
					// deepest child is its sole concrete authorization owner. An
					// empty capability set marks a purely structural alias. A
					// non-empty declaration must match the owner so dead security
					// metadata cannot mislead consumers.
					if len(canonicalCapabilities(candidate.route.CapabilityTags)) == 0 ||
						sameCapabilities(candidate.route.CapabilityTags, owner.route.CapabilityTags) {
						continue
					}
				}
				if !sameCapabilities(candidate.route.CapabilityTags, owner.route.CapabilityTags) {
					return fmt.Errorf(
						"surface contribution invalid for module %q: target %q path %q aliases %q and owner %q with different capabilities",
						strings.TrimSpace(contribution.ModuleID), target, path, candidate.route.ID, owner.route.ID,
					)
				}
			}
		}
	}
	return nil
}

func validateRouteHierarchy(contribution Contribution, requireAllParents bool) error {
	moduleID := strings.TrimSpace(contribution.ModuleID)
	routesByID := make(map[string]Route, len(contribution.Routes))
	for _, route := range contribution.Routes {
		routeID := strings.TrimSpace(route.ID)
		if routeID != "" {
			routesByID[routeID] = route
		}
	}

	for _, route := range contribution.Routes {
		routeID := strings.TrimSpace(route.ID)
		parentID := strings.TrimSpace(route.ParentRouteID)
		if routeID == "" || parentID == "" {
			continue
		}
		if parentID == routeID {
			return fmt.Errorf("surface contribution invalid for module %q: route %q cannot parent itself", moduleID, routeID)
		}
		parent, exists := routesByID[parentID]
		if !exists {
			if requireAllParents {
				return fmt.Errorf("surface contribution invalid for module %q: route %q references missing parent %q", moduleID, routeID, parentID)
			}
			// ParentRouteID may deliberately reference a route published by a
			// different contribution. The complete graph validator resolves it.
			continue
		}
		for _, target := range route.Targets {
			if !parent.IncludesTarget(target) {
				return fmt.Errorf(
					"surface contribution invalid for module %q: route %q target %q is not declared by parent %q",
					moduleID, routeID, target, parentID,
				)
			}
		}
	}

	for _, route := range contribution.Routes {
		routeID := strings.TrimSpace(route.ID)
		if routeID == "" {
			continue
		}
		seen := make(map[string]struct{}, len(routesByID))
		currentID := routeID
		depth := 0
		for currentID != "" {
			if _, duplicate := seen[currentID]; duplicate {
				return fmt.Errorf("surface contribution invalid for module %q: parent route cycle includes %q", moduleID, currentID)
			}
			seen[currentID] = struct{}{}

			current, exists := routesByID[currentID]
			if !exists {
				break
			}
			parentID := strings.TrimSpace(current.ParentRouteID)
			if parentID == "" {
				break
			}
			depth++
			if depth > maxRouteHierarchyDepth {
				return fmt.Errorf(
					"surface contribution invalid for module %q: route hierarchy exceeds maximum depth %d at route %q",
					moduleID, maxRouteHierarchyDepth, routeID,
				)
			}
			currentID = parentID
		}
	}
	return nil
}

func cloneOwnedRoute(route Route) Route {
	clone := route
	clone.CapabilityTags = append([]string(nil), route.CapabilityTags...)
	clone.Targets = append([]Target(nil), route.Targets...)
	if route.Section != nil {
		section := *route.Section
		clone.Section = &section
	}
	if route.RendersEntity != nil {
		entity := *route.RendersEntity
		clone.RendersEntity = &entity
	}
	if route.Page != nil {
		clone.Page = ClonePageContract(route.Page)
	}
	return clone
}

func validateRouteOwnerGroup(candidates []routeOwnerCandidate, parents map[string]string) (routeOwnerCandidate, error) {
	queryFree := make([]routeOwnerCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if !candidate.stateAlias {
			queryFree = append(queryFree, candidate)
		}
	}
	if len(queryFree) == 0 {
		return routeOwnerCandidate{}, fmt.Errorf("navigation-state aliases require one query-free concrete owner")
	}

	owner := queryFree[0]
	for _, candidate := range queryFree[1:] {
		switch {
		case isRouteAncestor(owner.route.ID, candidate.route.ID, parents):
			owner = candidate
		case isRouteAncestor(candidate.route.ID, owner.route.ID, parents):
			// The current owner is already the deeper route.
		default:
			return routeOwnerCandidate{}, fmt.Errorf(
				"query-free routes %q and %q are siblings; declare one concrete owner",
				owner.route.ID, candidate.route.ID,
			)
		}
	}
	return owner, nil
}

// SameNavigationLocation reports whether two navigation URLs identify the
// same path and state. Unlike concrete ownership identity, query and fragment
// suffixes remain significant so filtered links stay visible in navigation.
func SameNavigationLocation(left, right string) bool {
	leftLocation, leftErr := parseRouteLocation(left)
	rightLocation, rightErr := parseRouteLocation(right)
	return leftErr == nil && rightErr == nil &&
		leftLocation.path == rightLocation.path && leftLocation.suffix == rightLocation.suffix
}

// SameConcreteRouteLocation reports whether two URLs resolve to the same
// concrete pathname, ignoring query and fragment navigation state.
func SameConcreteRouteLocation(left, right string) bool {
	leftPath, _ := routeOwnershipPath(left)
	rightPath, _ := routeOwnershipPath(right)
	return leftPath != "" && leftPath == rightPath
}

// CanonicalNavigationLocation normalizes a navigation URL's pathname while
// preserving its query and fragment state.
func CanonicalNavigationLocation(raw string) string {
	return canonicalNavigationLocation(raw)
}

func candidateOwnsRoute(candidate, current routeOwnerCandidate) bool {
	if candidate.stateAlias != current.stateAlias {
		return !candidate.stateAlias
	}
	if candidate.depth != current.depth {
		return candidate.depth > current.depth
	}
	return candidate.index < current.index
}

func routeOwnershipPath(raw string) (string, bool) {
	location, err := parseRouteLocation(raw)
	if err != nil {
		return "", false
	}
	return location.path, location.stateAlias
}

func canonicalNavigationLocation(raw string) string {
	location, err := parseRouteLocation(raw)
	if err != nil {
		return ""
	}
	return location.path + location.suffix
}

type routeLocation struct {
	path       string
	suffix     string
	stateAlias bool
}

// parseRouteLocation accepts only canonical, same-origin navigation
// locations. Surface paths are emitted into hrefs and route registries, so
// interpreting schemes, network-path references, encoded separators, dot
// segments, or router-dependent slash forms would make validation disagree
// with the browser or HTTP router.
func parseRouteLocation(raw string) (routeLocation, error) {
	if raw == "" {
		return routeLocation{}, fmt.Errorf("path is required")
	}
	if raw != strings.TrimSpace(raw) {
		return routeLocation{}, fmt.Errorf("path must not contain surrounding whitespace")
	}
	for _, character := range raw {
		if character == '\\' || unicode.IsControl(character) || unicode.IsSpace(character) {
			return routeLocation{}, fmt.Errorf("path contains unsafe whitespace, control, or backslash characters")
		}
	}
	if strings.Contains(raw, "%") {
		return routeLocation{}, fmt.Errorf("path must not contain percent-encoded or ambiguous bytes")
	}
	if !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") {
		return routeLocation{}, fmt.Errorf("path must be an absolute local pathname beginning with exactly one slash")
	}

	path := raw
	suffix := ""
	if index := strings.IndexAny(raw, "?#"); index >= 0 {
		path = raw[:index]
		suffix = raw[index:]
	}
	if path == "" || strings.Contains(path, "//") {
		return routeLocation{}, fmt.Errorf("path must use a non-empty canonical pathname without repeated slashes")
	}
	for segment := range strings.SplitSeq(path, "/") {
		if segment == "." || segment == ".." {
			return routeLocation{}, fmt.Errorf("path must not contain dot segments")
		}
	}
	if suffix != "" {
		if suffix == "?" || suffix == "#" || strings.Count(suffix, "#") > 1 {
			return routeLocation{}, fmt.Errorf("navigation state suffix is malformed")
		}
		if fragment := strings.IndexByte(suffix, '#'); fragment >= 0 && strings.Contains(suffix[fragment+1:], "?") {
			return routeLocation{}, fmt.Errorf("query state must precede fragment state")
		}
	}
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	return routeLocation{path: path, suffix: suffix, stateAlias: suffix != ""}, nil
}

func routeHierarchyDepth(routeID string, parents map[string]string) int {
	depth := 0
	current := strings.TrimSpace(routeID)
	seen := make(map[string]struct{}, len(parents))
	for current != "" {
		if _, exists := seen[current]; exists {
			break
		}
		seen[current] = struct{}{}

		parent := strings.TrimSpace(parents[current])
		if parent == "" {
			break
		}
		depth++
		current = parent
	}
	return depth
}

func isRouteAncestor(candidateID, routeID string, parents map[string]string) bool {
	candidateID = strings.TrimSpace(candidateID)
	current := strings.TrimSpace(routeID)
	if candidateID == "" || current == "" || candidateID == current {
		return false
	}
	seen := make(map[string]struct{}, len(parents))
	for current != "" {
		if _, exists := seen[current]; exists {
			return false
		}
		seen[current] = struct{}{}
		current = strings.TrimSpace(parents[current])
		if current == candidateID {
			return true
		}
	}
	return false
}

func canonicalCapabilities(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	capabilities := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		capabilities = append(capabilities, value)
	}
	sort.Strings(capabilities)
	return capabilities
}

func sameCapabilities(left, right []string) bool {
	left = canonicalCapabilities(left)
	right = canonicalCapabilities(right)
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
