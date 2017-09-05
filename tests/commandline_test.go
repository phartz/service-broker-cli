package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/phartz/service-broker-cli/sbcli"
)

func TestCmd(t *testing.T) {
	options := []string{"sb", "command", "option1", "option2", "-f", "--no-filter", "-t", "tags"}
	c := sbcli.NewCommandline(options)
	AssertEqual(t, c.Command, "command", "Expect command")
	AssertEqual(t, len(c.Options), 2, "Expect len options == 2")
	AssertEqual(t, c.Force, true, "Expect flag force = true")
	AssertEqual(t, c.Tags, "tags", "Expect flag Tags = tags")
	AssertIsTrue(t, c.NoFilter, "Expect flag no-filter = true")
}

func TestCustomParser(t *testing.T) {
	o := []string{"file", "command", "", ""}
	cmd := sbcli.NewCommandline(o)
	AssertIsTrue(t, cmd.Custom == "", "custom must be empty")

	o[2] = "-c"
	o[3] = "{\"key\":\"value\"}"
	cmd = sbcli.NewCommandline(o)
	AssertIsTrue(t, cmd.Custom != "", "custom must not be empty")

	customFileName := "customFile.json"
	d1 := []byte("{\"key\":\"value\"}")
	err := ioutil.WriteFile(customFileName, d1, 0644)
	sbcli.CheckErr(err)

	o[2] = "-c"
	o[3] = customFileName
	cmd = sbcli.NewCommandline(o)
	AssertEqual(t, cmd.Custom, string(d1), "issue with custom parameter file")

	os.Remove(customFileName)
}
