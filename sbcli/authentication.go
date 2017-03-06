package sbcli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// retreieves all service instances from the service broker
func Logout(cmd *Commandline) {
	fmt.Println("Loggin out...")

	c := Config{}
	c.load()
	c.Password = ""
	c.Username = ""
	c.save()

	fmt.Println("OK")
}

func Target(cmd *Commandline) {
	c := Config{}
	c.load()

	if len(cmd.Options) == 0 {
		if c.Host == "" {
			fmt.Printf("No target set!\n")
		} else {
			fmt.Println("API endpoint: %s", c.Host)
			fmt.Println("User:         %s", c.Username)
		}
	} else {
		sb := NewSBClient(&Credentials{Host: cmd.Options[0]})
		err := sb.TestConnection()
		CheckErr(err)

		c.Host = cmd.Options[0]
		c.Password = ""
		c.Username = ""
		c.save()

		fmt.Printf("Target set to %s\n\n", c.Host)
		fmt.Printf("You have to login now.\n")
		fmt.Printf("\tsb login\n")
	}
}

func Auth(cmd *Commandline) {
	if len(cmd.Options) != 2 {
		CheckErr(errors.New("Missing arguments!"), GetHelpText("Auth"))
	}
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		CheckErr(errors.New("No target set."))
	}
	fmt.Printf("Target: %s...", conf.Host)

	// check if host is reachable
	sb := NewSBClient(&Credentials{Host: conf.Host})
	err := sb.TestConnection()
	CheckErr(err)

	fmt.Printf("OK\n\n")

	fmt.Printf("\nAuthenticating...")
	sb = NewSBClient(&Credentials{Host: conf.Host, Username: cmd.Options[0], Password: cmd.Options[1]})
	_, err = sb.Catalog()
	CheckErr(err)

	conf.Username = cmd.Options[0]
	conf.Password = cmd.Options[1]
	conf.save()

	fmt.Printf("OK\n\n")
}

func Login(cmd *Commandline) {
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		CheckErr(errors.New("No target set!"))
	}
	fmt.Printf("Target: %s...", conf.Host)

	// check if host is reachable
	sb := NewSBClient(&Credentials{Host: conf.Host})
	err := sb.TestConnection()
	CheckErr(err)

	fmt.Printf("OK\n\n")

	c := Credentials{Host: conf.Host}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username> ")
	c.Username, _ = reader.ReadString('\n')
	c.Username = strings.TrimSpace(c.Username)

	if c.Username == "" {
		fmt.Printf("No username given, break!\n")
		os.Exit(1)
	}

	fmt.Println()

	ok := false
	for i := 0; i < 3; i++ {
		c.Password, _ = getPassword("Password> ")

		fmt.Printf("\nAuthenticating...")
		sb := NewSBClient(&c)
		_, err = sb.Catalog()
		if err != nil {
			fmt.Printf("Failed!\n\n")
			continue
		}
		fmt.Printf("OK\n\n")
		ok = true
		break
	}

	if ok {
		conf.Username = c.Username
		conf.Password = c.Password
		conf.save()
		fmt.Printf("Target:   %s\n", conf.Host)
		fmt.Printf("Username: %s\n", conf.Host)
	}
}
