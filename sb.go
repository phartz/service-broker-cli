package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Slice which stores the registered commands
var sbcommands []SBCommand = make([]SBCommand, 0)

func main() {
	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Println("Trace is activated...")
		fmt.Println()
	}

	// register commands
	addCommand("status", "st", "Shows the status", Status)
	addCommand("target", "t", "Sets or gets the target", Target)
	addCommand("login", "li", "Login to the target", Login)
	addCommand("logout", "lo", "Logout from the target", Logout)
	addCommand("auth", "a", "Authenticate to the target", Auth)
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

	if len(os.Args) == 1 {
		usage()
		return
	}

	c := new(Commandline)
	c.Parse(os.Args)

	for _, command := range sbcommands {
		if command.Name == c.Command || command.Shortcut == c.Command {
			command.Function(c)
			return
		}
	}

	usage()
}

// fnuction to add the commands into the command slice
func addCommand(name string, shortcut string, helptext string, function SBFunction) {
	sbcommands = append(sbcommands, SBCommand{Name: name, Shortcut: shortcut, Helptext: helptext, Function: function})
}

// prints the usage text
func usage() {
	fmt.Println(UsageText)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "\033[1mcommand\tshortcut\tdescription\n")
	for _, command := range sbcommands {
		fmt.Fprintf(w, "\x1b[36;1m%s\x1b[0m\t%s\t%s\n", command.Name, command.Shortcut, command.Helptext)
	}
	w.Flush()
}

// creates the Servicebroker client, in later version the user credentials should be read out of a file
func getServiceBroker() *SBClient {
	c := Config{}
	c.load()

	return createSBClient(&Credentials{Host: c.Host, Password: c.Password, Username: c.Username})
}

// creates the Servicebroker client, in later version the user credentials should be read out of a file
func createSBClient(c *Credentials) *SBClient {
	var sb SBClient
	sb.Host = c.Host
	sb.Username = c.Username
	sb.Password = c.Password

	return &sb
}

// show a not implemented message
func notImplemented(options []string) {
	fmt.Println("not implemented yet")
}
