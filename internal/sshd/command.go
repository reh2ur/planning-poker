package sshd

import "strings"

// RoomFromCommand extracts the room ID from the SSH command, like Git passes a repo path.
func RoomFromCommand(cmd string) string {
	return strings.TrimSpace(cmd)
}
