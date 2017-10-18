package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/phartz/service-broker-cli/sbcli"
)

// Slice which stores the registered commands
var sbcommands []sbcli.SBCommand = make([]sbcli.SBCommand, 0)

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

	c := sbcli.NewCommandline(os.Args)

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
	addCommand("api", "", "Sets or gets the api", sbcli.Api)
	addCommand("login", "l", "Login to the target", sbcli.Login)
	addCommand("logout", "lo", "Logout from the target", sbcli.Logout)
	addCommand("auth", "", "Authenticate to the target", sbcli.Auth)
	addCommand("version", "-v", "Print the version", ShowVersion)
	addCommand("", "", "", nil)
	addCommand("marketplace", "m", "List available offerings in the marketplace", sbcli.Marketplace)
	addCommand("services", "s", "List all service instances in the target space", sbcli.Services)
	addCommand("service", "", "Show service instance info", sbcli.Service)
	addCommand("", "", "", nil)
	addCommand("create-service", "cs", "Create a service instance", sbcli.CreateService)
	addCommand("update-service", "", "Update a service instance", sbcli.UpdateService)
	addCommand("delete-service", "ds", "Delete a service instance", sbcli.DeleteService)
	addCommand("", "", "", nil)
	addCommand("create-service-key", "csk", "Create key for a service instance", sbcli.CreateServiceKey)
	addCommand("service-keys", "sk", "List keys for a service instance", sbcli.ServiceKeys)
	addCommand("delete-service-key", "dsk", "Delete a service key", sbcli.DeleteServiceKey)
}

// fnuction to add the commands into the command slice
func addCommand(name string, shortcut string, helptext string, function sbcli.SBFunction) {
	sbcommands = append(sbcommands, sbcli.SBCommand{Name: name, Shortcut: shortcut, Helptext: helptext, Function: function})
}

func Help(cmd *sbcli.Commandline) {
	if len(cmd.Options) == 0 {
		usage()
		return
	}
	fmt.Println(sbcli.GetHelpText(cmd.Options[0]))
}

// prints the usage text
func usage() {
	fmt.Println(sbcli.UsageText)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "command\tshortcut\tdescription\n")
	for _, command := range sbcommands {
		fmt.Fprintf(w, "%s\t%s\t%s\n", command.Name, command.Shortcut, command.Helptext)
	}
	w.Flush()
}

func ShowVersion(cmd *sbcli.Commandline) {
	fmt.Printf("sb version %s (%s)\n", getVersionString(), BuildTime)
}

func getVersionString() string {
	return fmt.Sprintf("%s+%s", Version, Build)
}
