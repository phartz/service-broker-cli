package sbcli

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

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

func getUserHome() string {
	usr, err := user.Current()
	CheckErr(err)
	return usr.HomeDir
}

func CheckErr(err error, helpTexts ...string) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "error: %v\n", err)

	for _, text := range helpTexts {
		fmt.Println()
		fmt.Println(text)
	}

	os.Exit(1)
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func prettyPrintJson(jsonString string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonString), "", "  ")
	return string(out.Bytes()), err
}
