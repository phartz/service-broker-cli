package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Slice which stores the registered commands
var sbcommands []SBCommand = make([]SBCommand, 0)

func main() {
	// register commands
	addCommand("marketplace", "m", "List available offerings in the marketplace", marketplace)
	addCommand("services", "s", "List all service instances in the target space", marketplace)
	addCommand("service", "", "Show service instance info", notImplemented)
	addCommand("", "", "", nil)
	addCommand("create-service", "cs", "Create a service instance", notImplemented)
	addCommand("update-service", "us", "Update a service instance", notImplemented)
	addCommand("delete-service", "ds", "Delete a service instance", notImplemented)
	addCommand("rename-service", "rs", "Rename a service instance", notImplemented)
	addCommand("", "", "", nil)
	addCommand("create-service-key", "csk", "Create key for a service instance", notImplemented)
	addCommand("service-keys", "sk", "List keys for a service instance", notImplemented)
	addCommand("service-key", "", "Show service key info", notImplemented)
	addCommand("delete-service-key", "dsk", "Delete a service key", notImplemented)

	if len(os.Args) == 1 {
		usage()
		return
	}

	for _, command := range sbcommands {
		if command.Name == os.Args[1] || command.Shortcut == os.Args[1] {
			command.Function(os.Args)
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
func getServiceBroker() *Servicebroker {
	sb := Servicebroker{Host: "http://localhost:3000", Username: "admin", Password: "admin"}
	return &sb
}

// show a not implemented message
func notImplemented(options []string) {
	fmt.Println("not implemented yet")
}

// retreieves all service instances from the service broker
func services(options []string) {
	sb := getServiceBroker()
	fmt.Printf("Getting services from Servicebroker \x1b[36;1m%s\x1b[0m as \x1b[36;1m%s\x1b[0m\n", sb.Host, sb.Username)

	services, err := sb.instances()
	if err != nil {
		panic(err)
	}

	fmt.Println("\x1b[32;1mOK\x1b[0m\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "\033[1mservice\tplans\tdescription\n")
	for _, service := range services.Resources {
		fmt.Fprintf(w, "\x1b[36;1m%s\x1b[0m\t%s\t%s\n", service.GUIDAtTenant, "", "")
	}
	w.Flush()
	fmt.Println("")
}

// retrieves the available service plans from the service broker
func marketplace(options []string) {
	sb := getServiceBroker()
	fmt.Printf("Getting services from Servicebroker \x1b[36;1m%s\x1b[0m as \x1b[36;1m%s\x1b[0m\n", sb.Host, sb.Username)

	catalog, err := sb.catalog()
	if err != nil {
		panic(err)
	}

	fmt.Println("\x1b[32;1mOK\x1b[0m\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "\033[1mservice\tplans\tdescription\n")
	for _, service := range catalog.Services {
		plans := make([]string, len(service.Plans))

		for i, plan := range service.Plans {
			plans[i] = plan.Name
		}

		fmt.Fprintf(w, "\x1b[36;1m%s\x1b[0m\t%s\t%s\n", service.Name, strings.Join(plans[:], ", "), service.Description)
	}
	w.Flush()
	fmt.Println("")
}
