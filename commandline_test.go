package main

import (
	"fmt"
	"io/ioutil"
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

func assertNotNil(t *testing.T, a interface{}, message string) {
	if a != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is nil", a)
	}
	t.Fatal(message)
}

func assertIsNil(t *testing.T, a interface{}, message string) {
	if a == nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not nil", a)
	}
	t.Fatal(message)
}

func assertIsTrue(t *testing.T, a bool, message string) {
	if a {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not true", a)
	}
	t.Fatal(message)
}

func assertIsFalse(t *testing.T, a bool, message string) {
	if a {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not true", a)
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

func TestCustomParser(t *testing.T) {
	o := []string{"file", "command", "", ""}
	cmd := NewCommandline(o)
	assertIsTrue(t, cmd.Custom == "", "custom must be empty")

	o[2] = "-c"
	o[3] = "{\"key\":\"value\"}"
	cmd = NewCommandline(o)
	assertIsTrue(t, cmd.Custom != "", "custom must not be empty")

	d1 := []byte("{\"key\":\"value\"}")
	err := ioutil.WriteFile("../dat1", d1, 0644)
	checkErr(err)

	o[2] = "-c"
	o[3] = "../dat1"
	cmd = NewCommandline(o)
	assertEqual(t, cmd.Custom, string(d1), "issue with file")

}
