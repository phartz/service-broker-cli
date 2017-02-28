package main

import (
	"flag"
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
	Output     string
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

	c.Options = options[2:flagPos]
	c.NoOptions = len(c.Options) > 0

	if flagPos == -1 {
		return
	}

	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	force := flagSet.Bool("f", false, "")
	plan := flagSet.String("p", "", "")
	tags := flagSet.String("t", "", "")
	custom := flagSet.String("c", "", "")
	output := flagSet.String("o", "", "")
	flagSet.Parse(options[flagPos:])
	c.Force = *force
	c.Plan = *plan
	c.Tags = *tags
	c.Custom = *custom
	c.Output = *output

	return
}
