package sshd

import "testing"

func TestRoomFromCommand(t *testing.T) {
	tests := map[string]string{
		"abc123":          "abc123",
		"  team-checkout": "team-checkout",
		"":                "",
	}
	for in, want := range tests {
		if got := RoomFromCommand(in); got != want {
			t.Errorf("RoomFromCommand(%q) = %q, want %q", in, got, want)
		}
	}
}
