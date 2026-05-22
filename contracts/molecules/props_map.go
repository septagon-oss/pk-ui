package molecules

import "encoding/json"

// propsToMap converts any Props struct to map[string]any via JSON round-trip.
// Used by ToMap() methods on Props types to enable unified Component construction.
func propsToMap(v any) map[string]any {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return nil
	}
	return m
}
