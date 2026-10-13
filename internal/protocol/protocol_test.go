package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewIsVersionedAndRoundTrips(t *testing.T) {
	m, err := New("vote", map[string]string{"card": "5"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var got Message
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got.V != Version || got.Type != "vote" || string(got.Payload) != `{"card":"5"}` {
		t.Fatalf("unexpected message: %+v", got)
	}
}
