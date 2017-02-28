package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func CreateServiceKey(cmd *Commandline) {
	if len(cmd.Options) != 2 {
		color.Set(color.FgRed)
		fmt.Printf("Missing arguments!\n\n")
		color.Unset()
		fmt.Printf("sb create-service-key <service_name> <key_name>\n")
		return
	}

	sb := getServiceBroker()

	col := color.New(color.FgCyan)
	fmt.Printf("Creating service key %s for service instance %s as %s...\n", col.Sprint(cmd.Options[1]), col.Sprint(cmd.Options[0]), col.Sprint(sb.Username))

	serviceID, planID, err := getServiceIDPlanID(cmd.Options[0])
	if err != nil {
		printErr(err)
		return
	}

	data := BindPayload{ServiceID: serviceID, PlanID: planID}
	result, err := sb.Bind(&data, cmd.Options[0], cmd.Options[1])
	if err != nil {
		printErr(err)
		return
	}

	fmt.Println("\x1b[32;1mOK\x1b[0m\n\n")

	fmt.Println(prettyPrintJson(result))
}

func ServiceKeys(cmd *Commandline) {
	sb := getServiceBroker()

	if len(cmd.Options) != 1 {
		color.Set(color.FgRed)
		fmt.Printf("Missing arguments!\n\n")
		color.Unset()
		fmt.Printf("sb service-keys <service_name>\n")
		return
	}

	col := color.New(color.FgCyan)
	fmt.Printf("Getting service keys for service instance %s as %s...\n\n", col.Sprint(cmd.Options[0]), col.Sprint(sb.Username))

	services, err := sb.Instances()
	if err != nil {
		printErr(err)
		return
	}

	for _, service := range services.Resources {
		if service.GUIDAtTenant == cmd.Options[0] {
			if len(service.Credentials) == 0 {
				color.Set(color.FgCyan)
				fmt.Println("none")
				color.Unset()
				return
			}

			fmt.Println("name")
			color.Set(color.FgCyan)
			for _, c := range service.Credentials {
				fmt.Println(c.GUIDAtTenant)
			}
			color.Unset()
			return
		}
	}
}

func DeleteServiceKey(cmd *Commandline) {
	sb := getServiceBroker()

	if len(cmd.Options) != 2 {
		color.Set(color.FgRed)
		fmt.Printf("Missing arguments!\n\n")
		color.Unset()
		fmt.Printf("sb create-service-key <service_name> <key_name>\n")
		return
	}

	if !cmd.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nReally delete service key %s> ", cmd.Options[1])
		asked, _ := reader.ReadString('\n')
		asked = strings.TrimSpace(asked)
		if asked != "yes" {
			color.Set(color.FgMagenta)
			fmt.Println("Delete cancelled!")
			color.Unset()
			return
		}
	}

	col := color.New(color.FgCyan)
	fmt.Printf("Delete key %s for service instance %s as %s...\n\n", col.Sprint(cmd.Options[1]), col.Sprint(cmd.Options[0]), col.Sprint(sb.Username))

	serviceID, planID, err := getServiceIDPlanID(cmd.Options[0])
	if err != nil {
		printErr(err)
		return
	}
	data := BindPayload{ServiceID: serviceID, PlanID: planID}

	err = sb.UnBind(&data, cmd.Options[0], cmd.Options[1])
	if err != nil {
		printErr(err)
		return
	}

	fmt.Println("\x1b[32;1mOK\x1b[0m\n")
}
