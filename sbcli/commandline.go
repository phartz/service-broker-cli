package sbcli

import (
	"flag"
	"io/ioutil"
	"strings"
)

type Commandline struct {
	Executable string
	Command    string
	Options    []string
	NoOptions  bool
	Force      bool
	Plan       string
	Tags       string
	Custom     string
	JSON       bool
	NoFilter   bool
}

func (c *Commandline) Parse(options []string) (err error) {
	c.Command = "help"
	c.Executable = options[0]
	c.NoOptions = false

	err = nil
	if len(options) == 1 {
		return
	}

	c.Command = options[1]

	if len(options) == 2 {
		return
	}

	flagPos := -1
	// find first occurence of '-'
	for i, option := range options[2:] {
		if strings.HasPrefix(option, "-") {
			flagPos = i + 2
			break
		}
	}

	if flagPos == -1 {
		c.Options = options[2:]
		c.NoOptions = len(c.Options) > 0
		return
	}

	c.Options = options[2:flagPos]
	c.NoOptions = len(c.Options) > 0

	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	force := flagSet.Bool("f", false, "")
	plan := flagSet.String("p", "", "")
	tags := flagSet.String("t", "", "")
	custom := flagSet.String("c", "", "")
	json := flagSet.Bool("j", false, "")
	noFilter := flagSet.Bool("no-filter", false, "")

	flagSet.Parse(options[flagPos:])
	c.Force = *force
	c.Plan = *plan
	c.Tags = *tags
	c.Custom = c.checkCustom(*custom)
	c.JSON = *json
	c.NoFilter = *noFilter

	return
}

func (c *Commandline) checkCustom(unknown string) string {
	if unknown == "" {
		return ""
	}

	if strings.HasPrefix(unknown, "{") {
		return unknown
	}

	bytes, err := ioutil.ReadFile(unknown)
	if len(bytes) > 0 {
		return string(bytes)
	}

	// check Home Folder
	if strings.HasPrefix(unknown, "~") {
		bytes, err = ioutil.ReadFile(getUserHome() + unknown[1:])
		CheckErr(err)
		return string(bytes)
	}

	return unknown
}

func NewCommandline(options []string) *Commandline {
	c := new(Commandline)
	c.Parse(options)
	return c
}
