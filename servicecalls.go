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

// retreieves all service instances from the service broker
func Services(options []string) {
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

// retrieves the available service plans from the service broker
func Marketplace(options []string) {
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

func Service(options []string) {
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

func CreateService(options []string) {
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

func DeleteService(options []string) {
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

func UpdateService(options []string) {
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
