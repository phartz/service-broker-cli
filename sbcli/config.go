package sbcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
)

type Credentials struct {
	Host              string `json:"host"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	SkipSslValidation bool   `json:"skip_ssl_validation"`
}

type Config struct {
	Credentials
}

const (
	ConfigFile = ".sb"
)

func findConfig(dir string, recursive ...bool) (string, error) {
	// check file
	_, err := os.Stat(filepath.Join(dir, ConfigFile))
	if err == nil {
		return dir, nil
	}

	// if not found and root, return error
	if dir == "/" || (len(recursive) > 0 && !recursive[0]) {
		return "", errors.New("config: config not found")
	}

	// remove trailing slash
	if len(dir) > 0 && dir[len(dir)-1] == '/' {
		dir = dir[0 : len(dir)-1]
	}

	// split path and call function again
	parent, _ := path.Split(dir)
	return findConfig(parent)
}

func getConfig() (string, error) {
	dir, _ := filepath.Abs(".")
	config, err := findConfig(dir)

	if err != nil {
		config, err = findConfig(getUserHome(), false)
		return "", err
	}

	return filepath.Join(config, ConfigFile), nil
}

func (c *Config) load() error {
	if os.Getenv("SB_HOST") != "" &&
		os.Getenv("SB_USERNAME") != "" &&
		os.Getenv("SB_PASSWORD") != "" {
		c.Credentials.Host = CleanTargetURI(os.Getenv("SB_HOST"))
		c.Credentials.Username = os.Getenv("SB_USERNAME")
		c.Credentials.Password = os.Getenv("SB_PASSWORD")
		if skipVerify, ok := os.LookupEnv("SB_SKIP_SSL_VERIFY"); ok {
			value, err := strconv.ParseBool(skipVerify)
			if err != nil {
				log.Fatal(err)
			}
			c.SkipSslValidation = value
		}
		return nil
	}

	file, err := getConfig()
	if err != nil {
		return err
	}

	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonFile, c)
	if err != nil {
		return err
	}

	return nil
}

func testFolder(folder string) {
	fmt.Print("Test: ", folder, " ... ")

	file := filepath.Join(folder, ConfigFile)
	f, err := os.Create(file)
	defer f.Close()
	CheckErr(err)

	found, err := getConfig()
	if err != nil {
		fmt.Println("nope")
	} else {
		fmt.Println("found", found)
	}

	os.Remove(file)
}

func (c *Config) save() error {
	configJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}

	dir, _ := filepath.Abs(".")
	err = ioutil.WriteFile(filepath.Join(dir, ConfigFile), configJSON, 0600)
	if err != nil {
		// try to save in users home path
		usr, err := user.Current()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(usr.HomeDir, ConfigFile), configJSON, 0600)
	}
	return nil
}

func LoadConfig() *Config {
	c := Config{}
	c.load()
	return &c
}
