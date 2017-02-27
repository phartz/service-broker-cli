package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Slice which stores the registered commands
var sbcommands []SBCommand = make([]SBCommand, 0)

func main() {
	// register commands
	addCommand("target", "", "Sets or gets the target", target)
	addCommand("login", "", "Login to the target", login)
	addCommand("auth", "", "Authenticate to the target", auth)
	addCommand("", "", "", nil)
	addCommand("marketplace", "m", "List available offerings in the marketplace", marketplace)
	addCommand("services", "s", "List all service instances in the target space", services)
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
			command.Function(os.Args[2:])
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
	c := Config{}
	c.load()

	return createServiceBroker(&Credentials{Host: c.Host, Password: c.Password, Username: c.Username})
}

// creates the Servicebroker client, in later version the user credentials should be read out of a file
func createServiceBroker(c *Credentials) *Servicebroker {
	var sb Servicebroker
	sb.Host = c.Host
	sb.Username = c.Username
	sb.Password = c.Password

	return &sb
}

// show a not implemented message
func notImplemented(options []string) {
	fmt.Println("not implemented yet")
}

// retreieves all service instances from the service broker
func services(options []string) {
	sb := getServiceBroker()
	fmt.Printf("Getting services from Servicebroker \x1b[96m%s\x1b[0m as \x1b[96m%s\x1b[0m\n", sb.Host, sb.Username)

	catalog, err := sb.catalog()
	plans := make(map[string]string)
	for _, plan := range catalog.Services[0].Plans {
		plans[plan.ID] = plan.Name
	}

	services, err := sb.instances()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\x1b[92mOK\x1b[0m\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.FilterHTML)
	bold := color.New(color.Bold)
	bold.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "name", "service", "plan", "bound apps", "last operation")
	for _, service := range services.Resources {
		planName := "unknown"
		if name, found := plans[service.PlanGUID]; found {
			planName = name
		}

		fmt.Fprintf(w, "\x1b[96m%s\x1b[0m\t%s\t%s\t%s\t%s\n", service.GUIDAtTenant, catalog.Services[0].Name, planName, "unknown", service.State)
	}
	w.Flush()
	fmt.Println("")
}

func target(options []string) {
	c := Config{}
	c.load()

	if len(options) == 0 {
		if c.Host == "" {
			fmt.Printf("\033[1mNo target set!\n")
		} else {
			fmt.Printf("Actual target \033[1m%s\n", c.Host)
		}
	} else {
		sb := createServiceBroker(&Credentials{Host: options[0]})
		err := sb.testConnection()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		c.Host = options[0]
		c.Password = ""
		c.Username = ""
		c.save()

		fmt.Printf("Target set to \033[1m%s\n\n", c.Host)
		fmt.Printf("\x1b[0mYou have to login now.\n")
		fmt.Printf("\033[1m\tsb login\n")
	}
}

func auth(options []string) {
	if len(options) != 2 {
		fmt.Printf("\033[31m\tsb auth <username> <password>\n")
		fmt.Printf("\033[31mNo authentication given.\n")
		return
	}
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		fmt.Printf("\033[1mNo target set!\nSet target first!\n")
		return
	}
	fmt.Printf("Actual target \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createServiceBroker(&Credentials{Host: conf.Host})
	err := sb.testConnection()
	if err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\x1b[32;1mOK\x1b[0m\n")

	sb = createServiceBroker(&Credentials{Host: conf.Host, Username: options[0], Password: options[1]})
	_, err = sb.catalog()
	if err != nil {
		fmt.Printf("\033[31mWrong password!\x1b[0m\n\n")
		return
	}

	conf.Username = options[0]
	conf.Password = options[1]
	conf.save()
	fmt.Printf("\033[1mLogin successfull!\x1b[0m\n")
}

func login(options []string) {
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		fmt.Printf("\033[1mNo target set!\nSet target first!\n")
		return
	}
	fmt.Printf("Actual target \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createServiceBroker(&Credentials{Host: conf.Host})
	err := sb.testConnection()
	if err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n")
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\x1b[32;1mOK\x1b[0m\n")

	c := Credentials{Host: conf.Host}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username> ")
	c.Username, _ = reader.ReadString('\n')
	c.Username = strings.TrimSpace(c.Username)

	if c.Username == "" {
		fmt.Printf("\x1b[31;1mNo username given, break!\x1b[0m\n")
		return
	}

	ok := false
	for i := 0; i < 3; i++ {
		c.Password, _ = getPassword("Password> ")

		sb := createServiceBroker(&c)
		_, err = sb.catalog()
		if err != nil {
			fmt.Printf("\033[31mWrong password!\x1b[0m\n\n")
			continue
		}

		ok = true
		break
	}

	if ok {
		conf.Username = c.Username
		conf.Password = c.Password
		conf.save()
		fmt.Printf("\033[1mLogin successfull!\x1b[0m\n")
	}
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
