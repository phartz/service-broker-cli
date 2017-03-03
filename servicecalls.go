package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// retreieves all service instances from the service broker
func Services(cmd *Commandline) {
	sb := NewSBClient()
	fmt.Printf("Getting services from Servicebroker %s as %s\n", sb.Host, sb.Username)

	catalog, err := sb.Catalog()
	checkErr(err)

	plans := make(map[string]string)
	for _, plan := range catalog.Services[0].Plans {
		plans[plan.ID] = plan.Name
	}

	services, err := sb.Instances()
	checkErr(err)

	// update states
	for _, service := range services.Resources {
		if service.State == "deleted" {
			continue
		}

		_, _ = sb.LastState(service.GUIDAtTenant)
	}
	services, err = sb.Instances()
	checkErr(err)

	fmt.Printf("OK\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "name", "service", "plan", "bound apps", "last operation")
	for _, service := range services.Resources {
		if service.State == "deleted" {
			continue
		}

		planName := "unknown"
		if name, found := plans[service.PlanGUID]; found {
			planName = name
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", service.GUIDAtTenant, catalog.Services[0].Name, planName, "unknown", service.State)
	}
	w.Flush()
	fmt.Println("")
}

// retrieves the available service plans from the service broker
func Marketplace(cmd *Commandline) {
	sb := NewSBClient()

	fmt.Printf("Getting services from Servicebroker %s as %s\n", sb.Host, sb.Username)

	catalog, err := sb.Catalog()
	checkErr(err)

	fmt.Println("OK\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, "service\tplans\tdescription\n")
	for _, service := range catalog.Services {
		plans := make([]string, len(service.Plans))

		for i, plan := range service.Plans {
			plans[i] = plan.Name
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n", service.Name, strings.Join(plans[:], ", "), service.Description)
	}
	w.Flush()
	fmt.Println("")
}

func getServiceIDPlanID(servicename string) (*BindPayload, error) {
	sb := NewSBClient()

	services, err := sb.Instances()
	checkErr(err)

	for _, service := range services.Resources {
		if service.GUIDAtTenant == servicename {
			return &BindPayload{ServiceID: service.ServiceGUID, PlanID: service.PlanGUID}, nil
		}
	}

	return nil, errors.New("Service not found!")
}

func Service(cmd *Commandline) {
	if len(cmd.Options) != 1 {
		checkErr(errors.New("Missing arguments!"), GetHelpText("Service"))
	}

	sb := NewSBClient()

	catalog, err := sb.Catalog()
	checkErr(err)

	plans := make(map[string]string)
	for _, plan := range catalog.Services[0].Plans {
		plans[plan.ID] = plan.Name
	}

	services, err := sb.Instances()
	checkErr(err)

	for _, service := range services.Resources {
		if service.GUIDAtTenant == cmd.Options[0] {
			_, _ = sb.LastState(service.GUIDAtTenant)
			break
		}
	}

	services, err = sb.Instances()
	checkErr(err)

	for _, service := range services.Resources {
		if service.GUIDAtTenant == cmd.Options[0] {
			fmt.Println("")
			planName := "unknown"
			if name, found := plans[service.PlanGUID]; found {
				planName = name
			}
			lastState, err := sb.LastState(cmd.Options[0])
			checkErr(err)

			fmt.Printf("Service instance: %s\n", service.GUIDAtTenant)
			fmt.Printf("Service: %s\n", catalog.Services[0].Name)
			fmt.Printf("Bound apps: %s\n", "unknown")
			fmt.Printf("Tags:%s\n", strings.Join(catalog.Services[0].Tags, ", "))
			fmt.Printf("Plan: %s\n", planName)
			fmt.Printf("Description: %s\n", catalog.Services[0].Description)
			fmt.Printf("Documentation url: \n")
			fmt.Printf("Dashboard: \n")
			fmt.Printf("\n")
			fmt.Printf("Last Operation: %s\n", lastState.State)
			fmt.Printf("Status: %s\n", service.State)
			fmt.Printf("Message: %s\n", lastState.Description)
			fmt.Printf("Started: %s\n", service.CreatedAt)
			fmt.Printf("Updated: %s\n", service.UpdatedAt)
			return
		}
	}

	checkErr(errors.New("Service instance not found."))
}

func CreateService(cmd *Commandline) {
	sb := NewSBClient()
	fmt.Printf("Creating service at Servicebroker %s as %s\n", sb.Host, sb.Username)
	if len(cmd.Options) != 3 {
		checkErr(errors.New("Missing arguments!"), GetHelpText("CreateService"))
	}

	catalog, err := sb.Catalog()
	checkErr(err)

	if catalog.Services[0].Name != cmd.Options[0] {
		checkErr(errors.New("Service offering not found. Check the marketplace."))
	}

	var planID string
	for _, plan := range catalog.Services[0].Plans {
		if cmd.Options[1] == plan.Name {
			planID = plan.ID
		}
	}

	if planID == "" {
		checkErr(errors.New("Service plan not found. Check the marketplace."))
	}

	orgID, _ := newUUID()
	spaceID, _ := newUUID()

	data := ProvisonPayload{
		PlanID:           planID,
		OrganizationGUID: orgID,
		SpaceGUID:        spaceID,
		ServiceID:        catalog.Services[0].ID,
	}

	err = sb.Provision(&data, cmd.Options[2])
	checkErr(err)

	fmt.Printf("OK\n\n")

	fmt.Printf("Create in progress. Use 'sb services' or 'sb service %s' to check operation status.\n", cmd.Options[2])
}

func DeleteService(cmd *Commandline) {
	if len(cmd.Options) != 1 {
		checkErr(errors.New("Missing arguments!"), GetHelpText("DeleteService"))
	}

	sb := NewSBClient()

	instances, err := sb.Instances()
	checkErr(err)

	found := false
	for _, instance := range instances.Resources {
		if instance.GUIDAtTenant == cmd.Options[0] {
			found = true
			break
		}
	}

	if !found {
		checkErr(errors.New("Service not found!"))
	}

	if !cmd.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nReally delete service %s> ", cmd.Options[0])
		asked, _ := reader.ReadString('\n')
		asked = strings.TrimSpace(asked)
		if asked != "yes" {
			fmt.Println("Delete cancelled!")
			return
		}
	}

	data, err := getServiceIDPlanID(cmd.Options[0])
	checkErr(err)

	err = sb.Deprovision(data, cmd.Options[0])
	checkErr(err)

	fmt.Printf("Deleting service %s at %s as %s...\n", cmd.Options[0], sb.Host, sb.Username)
	fmt.Printf("OK\n\n")

	fmt.Printf("Delete in progress. Use 'sb services' or 'sb service %s' to check operation status.\n", cmd.Options[0])
}

func UpdateService(cmd *Commandline) {
	if len(cmd.Options) != 1 {
		checkErr(errors.New("Missing arguments!"), GetHelpText("UpdateService"))
	}

	sb := NewSBClient()

	instances, err := sb.Instances()
	checkErr(err)

	found := false
	for _, instance := range instances.Resources {
		if instance.GUIDAtTenant == cmd.Options[0] {
			found = true
			break
		}
	}

	if !found {
		checkErr(errors.New("Service not found!"))
	}

	err = sb.UpdateService(cmd.Options[0])
	checkErr(err)

	fmt.Printf("Updating service %s at %s as %s...\n", cmd.Options[0], sb.Host, sb.Username)
	fmt.Printf("OK\n\n")

	fmt.Printf("Update in progress. Use 'sb services' or 'sb service %s' to check operation status.\n", cmd.Options[0])
}
