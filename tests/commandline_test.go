package tests

import (
	"io/ioutil"
	"testing"

	"phartz.dedyn.io/gogs/phartz/service-broker-cli/sbcli"
)

func TestCmd(t *testing.T) {
	options := []string{"sb", "command", "option1", "option2", "-f", "-t", "tags"}
	c := sbcli.NewCommandline(options)
	AssertEqual(t, c.Command, "command", "Expect command")
	AssertEqual(t, len(c.Options), 2, "Expect len options == 2")
	AssertEqual(t, c.Force, true, "Expect flag force = true")
	AssertEqual(t, c.Tags, "tags", "Expect flag Tags = tags")
}

func TestCustomParser(t *testing.T) {
	o := []string{"file", "command", "", ""}
	cmd := sbcli.NewCommandline(o)
	AssertIsTrue(t, cmd.Custom == "", "custom must be empty")

	o[2] = "-c"
	o[3] = "{\"key\":\"value\"}"
	cmd = sbcli.NewCommandline(o)
	AssertIsTrue(t, cmd.Custom != "", "custom must not be empty")

	d1 := []byte("{\"key\":\"value\"}")
	err := ioutil.WriteFile("../dat1", d1, 0644)
	sbcli.CheckErr(err)

	o[2] = "-c"
	o[3] = "../dat1"
	cmd = sbcli.NewCommandline(o)
	AssertEqual(t, cmd.Custom, string(d1), "issue with file")
}
