package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func CreateServiceKey(cmd *Commandline) {
	if len(cmd.Options) != 2 {
		if cmd.JSON {
			fmt.Println("{\"error\":\"Missing arguments!\"}")
			return
		} else {
			printErr(errors.New("Missing arguments!"), GetHelpText("CreateServiceKey"))
		}
	}

	sb := NewSBClient()

	if !cmd.JSON {
		fmt.Printf("Creating service key %s for service instance %s as %s...\n", cmd.Options[1], cmd.Options[0], sb.Username)
	}

	serviceID, planID, err := getServiceIDPlanID(cmd.Options[0])
	if err != nil {
		printErr(err)
	}

	data := BindPayload{ServiceID: serviceID, PlanID: planID}
	result, err := sb.Bind(&data, cmd.Options[0], cmd.Options[1])
	if err != nil {
		printErr(err)
	}

	if !cmd.JSON {
		fmt.Println("OK\n\n")
	}

	fmt.Println(prettyPrintJson(result))
}

func ServiceKeys(cmd *Commandline) {
	sb := NewSBClient()

	if len(cmd.Options) != 1 {
		printErr(errors.New("Missing arguments!"), GetHelpText("ServiceKeys"))
	}

	fmt.Printf("Getting service keys for service instance %s as %s...\n\n", cmd.Options[0], sb.Username)

	services, err := sb.Instances()
	if err != nil {
		printErr(err)
	}

	for _, service := range services.Resources {
		if service.GUIDAtTenant == cmd.Options[0] {
			if len(service.Credentials) == 0 {
				fmt.Println("none")
				return
			}

			fmt.Println("name")
			for _, c := range service.Credentials {
				fmt.Println(c.GUIDAtTenant)
			}
			return
		}
	}
}

func DeleteServiceKey(cmd *Commandline) {
	sb := NewSBClient()

	if len(cmd.Options) != 2 {
		printErr(errors.New("Missing arguments!"), GetHelpText("DeleteServiceKey"))
	}

	if !cmd.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nReally delete service key %s> ", cmd.Options[1])
		asked, _ := reader.ReadString('\n')
		asked = strings.TrimSpace(asked)
		if asked != "yes" {
			fmt.Println("Delete cancelled!")
			return
		}
	}

	fmt.Printf("Delete key %s for service instance %s as %s...\n\n", cmd.Options[1], cmd.Options[0], sb.Username)

	serviceID, planID, err := getServiceIDPlanID(cmd.Options[0])
	if err != nil {
		printErr(err)
	}
	data := BindPayload{ServiceID: serviceID, PlanID: planID}

	err = sb.UnBind(&data, cmd.Options[0], cmd.Options[1])
	if err != nil {
		printErr(err)
	}

	fmt.Println("OK\n")
}
