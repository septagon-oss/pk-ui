// Validates: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package surface

import (
	"encoding/json"
	"testing"
)

func TestErrorDocumentJSONContract(t *testing.T) {
	content, err := json.Marshal(ErrorDocument{
		StatusCode:  404,
		Title:       "Not found",
		Description: "The requested resource does not exist.",
		HomeURL:     "/",
		HomeLabel:   "Home",
	})
	if err != nil {
		t.Fatalf("marshal error document: %v", err)
	}

	const want = `{"statusCode":404,"title":"Not found","description":"The requested resource does not exist.","homeUrl":"/","homeLabel":"Home"}`
	if string(content) != want {
		t.Fatalf("ErrorDocument JSON = %s, want %s", content, want)
	}
}
