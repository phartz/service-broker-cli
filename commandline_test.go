package main

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestCmd(t *testing.T) {
	options := []string{"sb", "command", "option1", "option2", "-f", "-t", "tags"}
	c := NewCommandline(options)
	assertEqual(t, c.Command, "command", "Expect command")
	assertEqual(t, len(c.Options), 2, "Expect len options == 2")
	assertEqual(t, c.Force, true, "Expect flag force = true")
	assertEqual(t, c.Tags, "tags", "Expect flag Tags = tags")
}
