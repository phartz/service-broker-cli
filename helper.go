package main

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func getPassword(prompt string) (password string, err error) {
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, state)
	term := terminal.NewTerminal(os.Stdout, "")
	password, err = term.ReadPassword(prompt)
	if err != nil {
		log.Fatal(err)
	}
	return
}
