package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phartz/service-broker-cli/sbcli"
	"github.com/phartz/service-broker-cli/tests"
)

func TestHelpTexts(t *testing.T) {
	registerCommands()

	for _, command := range sbcommands {
		if command.Name == "" {
			continue
		}

		text := sbcli.GetHelpText(command.Name)
		tests.AssertEqual(t, strings.HasPrefix(text, "Sorry"), false, fmt.Sprintf("No help text found for %s", command.Name))

		if command.Shortcut == "" {
			continue
		}

		text = sbcli.GetHelpText(command.Shortcut)
		tests.AssertEqual(t, strings.HasPrefix(text, "Sorry"), false, fmt.Sprintf("No help text found for %s", command.Shortcut))
	}
}
