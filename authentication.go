package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// retreieves all service instances from the service broker
func Logout(options []string) {
	fmt.Println("Loggin out...")

	c := Config{}
	c.load()
	c.Password = ""
	c.Username = ""
	c.save()

	fmt.Printf("\x1b[92mOK\x1b[0m\n\n")
}

func Target(options []string) {
	c := Config{}
	c.load()

	if len(options) == 0 {
		if c.Host == "" {
			fmt.Printf("\033[1mNo target set!\n")
		} else {
			fmt.Printf("Actual target \033[1m%s\n", c.Host)
		}
	} else {
		sb := createSBClient(&Credentials{Host: options[0]})
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

func Auth(options []string) {
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
	sb := createSBClient(&Credentials{Host: conf.Host})
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
	sb = createSBClient(&Credentials{Host: conf.Host, Username: options[0], Password: options[1]})
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

func Status(options []string) {
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		fmt.Printf("\033[1mNo target set!\nSet target first!\n")
		return
	}
	fmt.Printf("Target: \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createSBClient(&Credentials{Host: conf.Host})
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

	if conf.Password != "" {
		color.Set(color.FgCyan)
		fmt.Printf("Username: %s\n", conf.Host)
		color.Unset()
		return
	}

	color.Set(color.FgRed)
	fmt.Println("Your're not logged in!")
	color.Unset()
}

func Login(options []string) {
	conf := Config{}
	conf.load()

	// check host
	if conf.Host == "" {
		fmt.Printf("\033[1mNo target set!\nSet target first!\n")
		return
	}
	fmt.Printf("Target: \033[1m%s\x1b[0m...", conf.Host)

	// check if host is reachable
	sb := createSBClient(&Credentials{Host: conf.Host})
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
		sb := createSBClient(&c)
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
