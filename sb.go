package main

import (
	"bufio"
	"errors"
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
	addCommand("logout", "", "Logout from the target", logout)
	addCommand("auth", "", "Authenticate to the target", auth)
	addCommand("", "", "", nil)
	addCommand("marketplace", "m", "List available offerings in the marketplace", marketplace)
	addCommand("services", "s", "List all service instances in the target space", services)
	addCommand("service", "", "Show service instance info", service)
	addCommand("", "", "", nil)
	addCommand("create-service", "cs", "Create a service instance", createService)
	addCommand("update-service", "", "Update a service instance", updateService)
	addCommand("delete-service", "ds", "Delete a service instance", deleteService)
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
func getServiceBroker() *SBClient {
	c := Config{}
	c.load()

	return createServiceBroker(&Credentials{Host: c.Host, Password: c.Password, Username: c.Username})
}

// creates the Servicebroker client, in later version the user credentials should be read out of a file
func createServiceBroker(c *Credentials) *SBClient {
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

// retreieves all service instances from the service broker
func logout(options []string) {
	fmt.Println("Loggin out...")

	c := Config{}
	c.load()
	c.Password = ""
	c.Username = ""
	c.save()

	fmt.Printf("\x1b[92mOK\x1b[0m\n\n")
}

// retreieves all service instances from the service broker
func services(options []string) {
	sb := getServiceBroker()
	fmt.Printf("Getting services from Servicebroker \x1b[96m%s\x1b[0m as \x1b[96m%s\x1b[0m\n", sb.Host, sb.Username)

	catalog, err := sb.Catalog()
	plans := make(map[string]string)
	for _, plan := range catalog.Services[0].Plans {
		plans[plan.ID] = plan.Name
	}

	services, err := sb.Instances()
	if err != nil {
		panic(err)
	}

	// update states
	for _, service := range services.Resources {
		if service.State == "deleted" {
			continue
		}

		_, _ = sb.LastState(service.GUIDAtTenant)
	}
	services, err = sb.Instances()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\x1b[92mOK\x1b[0m\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.FilterHTML)
	bold := color.New(color.Bold)
	bold.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "name", "service", "plan", "bound apps", "last operation")
	for _, service := range services.Resources {
		if service.State == "deleted" {
			continue
		}

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
		err := sb.TestConnection()
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
		color.Set(color.FgRed)
		fmt.Printf("Missing arguments!\n\n")
		color.Unset()
		fmt.Printf("sb auth <username> <password>\n")
		fmt.Printf("No authentication given.\n")
		return
	}
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		color.Set(color.FgRed)
		fmt.Printf("No target set!\n")
		color.Unset()
		return
	}
	fmt.Printf("Target: \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createServiceBroker(&Credentials{Host: conf.Host})
	err := sb.TestConnection()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Printf("Failed!\n")
		color.Unset()
		fmt.Println(err.Error())
		return
	}
	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()

	fmt.Printf("\nAuthenticating...")
	sb = createServiceBroker(&Credentials{Host: conf.Host, Username: options[0], Password: options[1]})
	_, err = sb.Catalog()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Printf("Failed!\n\n")
		color.Unset()
		return
	}
	conf.Username = options[0]
	conf.Password = options[1]
	conf.save()

	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()
}

func login(options []string) {
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		fmt.Printf("\033[1mNo target set!\nSet target first!\n")
		return
	}
	fmt.Printf("Target: \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createServiceBroker(&Credentials{Host: conf.Host})
	err := sb.TestConnection()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Printf("Failed!\n")
		color.Unset()
		fmt.Println(err.Error())
		return
	}
	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()

	c := Credentials{Host: conf.Host}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username> ")
	c.Username, _ = reader.ReadString('\n')
	c.Username = strings.TrimSpace(c.Username)

	if c.Username == "" {
		fmt.Printf("\x1b[31;1mNo username given, break!\x1b[0m\n")
		return
	}

	fmt.Println()

	ok := false
	for i := 0; i < 3; i++ {
		c.Password, _ = getPassword("Password> ")

		fmt.Printf("\nAuthenticating...")
		sb := createServiceBroker(&c)
		_, err = sb.Catalog()
		if err != nil {
			fmt.Printf("\033[31mFailed!\x1b[0m\n\n")
			continue
		}
		color.Set(color.FgGreen)
		fmt.Printf("OK\n\n")
		color.Unset()
		ok = true
		break
	}

	if ok {
		conf.Username = c.Username
		conf.Password = c.Password
		conf.save()
		color.Set(color.FgCyan)
		fmt.Printf("Target:   %s\n", conf.Host)
		fmt.Printf("Username: %s\n", conf.Host)
		color.Unset()
	}
}

// retrieves the available service plans from the service broker
func marketplace(options []string) {
	sb := getServiceBroker()

	fmt.Printf("Getting services from Servicebroker \x1b[36;1m%s\x1b[0m as \x1b[36;1m%s\x1b[0m\n", sb.Host, sb.Username)

	catalog, err := sb.Catalog()
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

func service(options []string) {
	sb := getServiceBroker()

	catalog, err := sb.Catalog()
	if err != nil {
		printErr(err)
		return
	}

	plans := make(map[string]string)
	for _, plan := range catalog.Services[0].Plans {
		plans[plan.ID] = plan.Name
	}

	services, err := sb.Instances()
	if err != nil {
		printErr(err)
		return
	}

	for _, service := range services.Resources {
		if service.GUIDAtTenant == options[0] {
			_, _ = sb.LastState(service.GUIDAtTenant)
			break
		}
	}

	services, err = sb.Instances()
	if err != nil {
		printErr(err)
		return
	}

	for _, service := range services.Resources {
		if service.GUIDAtTenant == options[0] {
			fmt.Println("")
			planName := "unknown"
			if name, found := plans[service.PlanGUID]; found {
				planName = name
			}
			col := color.New(color.FgCyan)

			lastState, err := sb.LastState(options[0])
			if err != nil {
				fmt.Printf("Failed!\n")
				printErr(err)
				return
			}

			fmt.Printf("Service instance: %s\n", col.Sprint(service.GUIDAtTenant))
			fmt.Printf("Service: %s\n", col.Sprint(catalog.Services[0].Name))
			fmt.Printf("Bound apps: %s\n", col.Sprint("unknown"))
			fmt.Printf("Tags:%s\n", col.Sprint(strings.Join(catalog.Services[0].Tags, ", ")))
			fmt.Printf("Plan: %s\n", col.Sprint(planName))
			fmt.Printf("Description: %s\n", col.Sprint(catalog.Services[0].Description))
			fmt.Printf("Documentation url: \n")
			fmt.Printf("Dashboard: \n")
			fmt.Printf("\n")
			fmt.Printf("Last Operation: %s\n", col.Sprint(lastState.State))
			fmt.Printf("Status: %s\n", col.Sprint(service.State))
			fmt.Printf("Message: %s\n", col.Sprint(lastState.Description))
			fmt.Printf("Started: %s\n", col.Sprint(service.CreatedAt))
			fmt.Printf("Updated: %s\n", col.Sprint(service.UpdatedAt))
			return
		}
	}

	color.Set(color.FgRed)
	fmt.Printf("Failed!\n")
	color.Unset()
	fmt.Println("Service instance not found.")
}

func createService(options []string) {
	sb := getServiceBroker()
	fmt.Printf("Creating service at Servicebroker \x1b[36;1m%s\x1b[0m as \x1b[36;1m%s\x1b[0m\n", sb.Host, sb.Username)
	if len(options) != 3 {
		color.Set(color.FgRed)
		fmt.Printf("Failed!\n")
		color.Unset()
		fmt.Println("sb create-service <servicetyp> <planname> <ervicename>")
		return
	}

	catalog, err := sb.Catalog()
	if err != nil {
		printErr(err)
		return
	}

	if catalog.Services[0].Name != options[0] {
		fmt.Printf("Failed!\n")
		color.Unset()
		fmt.Println("Service offering not found. Check the marketplace.")
		fmt.Println("\tsb marketplace")
		return
	}

	var planID string
	for _, plan := range catalog.Services[0].Plans {
		if options[1] == plan.Name {
			planID = plan.ID
		}
	}

	if planID == "" {
		fmt.Printf("Failed!\n")
		color.Unset()
		fmt.Println("Service plan not found. Check the marketplace.")
		fmt.Println("\tsb marketplace")
		return
	}

	orgID, _ := newUUID()
	spaceID, _ := newUUID()

	data := ProvisonPayload{
		PlanID:           planID,
		OrganizationGUID: orgID,
		SpaceGUID:        spaceID,
		ServiceID:        catalog.Services[0].ID,
	}

	err = sb.Provision(&data, options[2])

	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()

	fmt.Printf("Create in progress. Use '")
	color.Set(color.FgGreen)
	fmt.Printf("cf services")
	color.Unset()
	fmt.Printf("' or '")
	color.Set(color.FgGreen)
	fmt.Printf("cf service %s", options[2])
	color.Unset()
	fmt.Printf("' to check operation status.\n")
}

func deleteService(options []string) {
	sb := getServiceBroker()

	instances, err := sb.Instances()
	if err != nil {
		printErr(err)
		return
	}

	found := false
	for _, instance := range instances.Resources {
		if instance.GUIDAtTenant == options[0] {
			found = true
			break
		}
	}

	if !found {
		printErr(errors.New("Service not found!"))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nReally delete service %s> ", options[0])
	asked, _ := reader.ReadString('\n')
	asked = strings.TrimSpace(asked)
	if asked != "yes" {
		color.Set(color.FgMagenta)
		fmt.Println("Delete cancelled!")
		color.Unset()
		return
	}

	err = sb.Deprovision(options[0])

	fmt.Printf("Deleting service %s at %s as %s...\n", options[0], sb.Host, sb.Username)
	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()

	fmt.Printf("Delete in progress. Use '")
	color.Set(color.FgGreen)
	fmt.Printf("cf services")
	color.Unset()
	fmt.Printf("' or '")
	color.Set(color.FgGreen)
	fmt.Printf("cf service %s", options[0])
	color.Unset()
	fmt.Printf("' to check operation status.\n")
}

func updateService(options []string) {
	sb := getServiceBroker()

	instances, err := sb.Instances()
	if err != nil {
		printErr(err)
		return
	}

	found := false
	for _, instance := range instances.Resources {
		if instance.GUIDAtTenant == options[0] {
			found = true
			break
		}
	}

	if !found {
		printErr(errors.New("Service not found!"))
		return
	}

	err = sb.UpdateService(options[0])

	fmt.Printf("Updating service %s at %s as %s...\n", options[0], sb.Host, sb.Username)
	color.Set(color.FgGreen)
	fmt.Printf("OK\n\n")
	color.Unset()

	fmt.Printf("Update in progress. Use '")
	color.Set(color.FgGreen)
	fmt.Printf("cf services")
	color.Unset()
	fmt.Printf("' or '")
	color.Set(color.FgGreen)
	fmt.Printf("cf service %s", options[0])
	color.Unset()
	fmt.Printf("' to check operation status.\n")
}
