package main

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
)

type Config struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
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
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		config, err = findConfig(usr.HomeDir, false)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(config, ConfigFile), nil
}

func getCon() (*Config, error) {
	file, err := getConfig()
	if err != nil {
		return nil, err
	}

	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func testFolder(folder string) {
	fmt.Print("Test: ", folder, " ... ")

	file := filepath.Join(folder, ConfigFile)
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	found, err := getConfig()
	if err != nil {
		fmt.Println("nope")
	} else {
		fmt.Println("found", found)
	}

	os.Remove(file)
}

/*
func main() {
	c, err := getCon()
	if err != nil {
		panic(err)
	}
	fmt.Println(c.Password)
}
*/
