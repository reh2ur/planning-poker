package room

import "testing"

func TestDeckByName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		numeric bool
		ok      bool
	}{
		{"fibonacci", "fibonacci", true, true},
		{"tshirt", "tshirt", false, true},
		{"custom", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, ok := DeckByName(tt.name)
			if ok != tt.ok || d.Name != tt.want || d.Numeric != tt.numeric {
				t.Fatalf("got %+v ok=%v", d, ok)
			}
		})
	}
}
