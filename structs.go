package main

type SBFunction func([]string)

type SBCommand struct {
	Name     string
	Shortcut string
	Helptext string
	Function SBFunction
}
