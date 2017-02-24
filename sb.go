package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type SBFunction func([]string)

func main() {
	m := map[string]SBFunction{
		"m":           marketplace,
		"marketplace": marketplace,
		"s":           services,
		"services":    services,
	}

	if len(os.Args) == 1 {
		usage()
		return
	}

	if f, ok := m[os.Args[1]]; ok {
		f(os.Args)
		return
	}

	usage()
}

func usage() {
	fmt.Println("nope..... you're wrong")
}

func getServiceBroker() *Servicebroker {
	sb := Servicebroker{Host: "http://localhost:3000", Username: "admin", Password: "admin"}
	return &sb
}

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
