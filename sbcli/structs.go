package sbcli

type SBFunction func(*Commandline)

type SBCommand struct {
	Name     string
	Shortcut string
	Helptext string
	Function SBFunction
}
