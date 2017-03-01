package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestHelpTexts(t *testing.T) {
	registerCommands()

	for _, command := range sbcommands {
		if command.Name == "" {
			continue
		}

		text := GetHelpText(command.Name)
		assertEqual(t, strings.HasPrefix(text, "Sorry"), false, fmt.Sprintf("No help text found for %s", command.Name))

		if command.Shortcut == "" {
			continue
		}

		text = GetHelpText(command.Shortcut)
		assertEqual(t, strings.HasPrefix(text, "Sorry"), false, fmt.Sprintf("No help text found for %s", command.Shortcut))
	}
}
