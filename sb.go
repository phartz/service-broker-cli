package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Slice which stores the registered commands
var sbcommands []SBCommand = make([]SBCommand, 0)

// Variables to identiy the build
var (
	Version   string
	Build     string
	BuildTime string
)

// main function
func main() {
	// Check trace
	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Println("Trace is activated...")
		fmt.Println()
	}

	// register all commands
	registerCommands()

	if len(os.Args) == 1 {
		usage()
		return
	}

	c := NewCommandline(os.Args)

	for _, command := range sbcommands {
		if command.Name == c.Command || command.Shortcut == c.Command {
			command.Function(c)
			return
		}
	}

	usage()
}

func registerCommands() {
	// register commands
	addCommand("help", "h", "Show help", Help)
	addCommand("", "", "", nil)
	addCommand("target", "t", "Sets or gets the target", Target)
	addCommand("login", "l", "Login to the target", Login)
	addCommand("logout", "lo", "Logout from the target", Logout)
	addCommand("auth", "", "Authenticate to the target", Auth)
	addCommand("version", "-v", "Print the version", ShowVersion)
	addCommand("", "", "", nil)
	addCommand("marketplace", "m", "List available offerings in the marketplace", Marketplace)
	addCommand("services", "s", "List all service instances in the target space", Services)
	addCommand("service", "", "Show service instance info", Service)
	addCommand("", "", "", nil)
	addCommand("create-service", "cs", "Create a service instance", CreateService)
	addCommand("update-service", "", "Update a service instance", UpdateService)
	addCommand("delete-service", "ds", "Delete a service instance", DeleteService)
	addCommand("", "", "", nil)
	addCommand("create-service-key", "csk", "Create key for a service instance", CreateServiceKey)
	addCommand("service-keys", "sk", "List keys for a service instance", ServiceKeys)
	addCommand("delete-service-key", "dsk", "Delete a service key", DeleteServiceKey)
}

// fnuction to add the commands into the command slice
func addCommand(name string, shortcut string, helptext string, function SBFunction) {
	sbcommands = append(sbcommands, SBCommand{Name: name, Shortcut: shortcut, Helptext: helptext, Function: function})
}

func Help(cmd *Commandline) {
	if len(cmd.Options) == 0 {
		usage()
		return
	}
	fmt.Println(GetHelpText(cmd.Options[0]))
}

// prints the usage text
func usage() {
	fmt.Println(UsageText)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "command\tshortcut\tdescription\n")
	for _, command := range sbcommands {
		fmt.Fprintf(w, "%s\t%s\t%s\n", command.Name, command.Shortcut, command.Helptext)
	}
	w.Flush()
}

// creates the Servicebroker client, in later version the user credentials should be read out of a file
func NewSBClient(cred ...*Credentials) *SBClient {
	var sb SBClient

	if len(cred) == 0 {
		conf := LoadConfig()
		sb.Host = conf.Host
		sb.Password = conf.Password
		sb.Username = conf.Username
	} else {
		sb.Host = cred[0].Host
		sb.Username = cred[0].Username
		sb.Password = cred[0].Password
	}

	return &sb
}

func ShowVersion(cmd *Commandline) {
	fmt.Printf("sb version %s (%s)\n", getVersionString(), BuildTime)
}
