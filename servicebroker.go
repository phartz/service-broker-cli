package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Servicebroker struct {
	Credentials
}

func (s *Servicebroker) SetCredentials(c Credentials) {
	s.Host = c.Host
	s.Username = c.Username
	s.Password = c.Password
}

func (s *Servicebroker) catalog() (*Catalog, error) {
	result, err := s.getResultFromBroker("v2/catalog", "GET", "{}")
	if err != nil {
		return nil, err
	}

	var c = new(Catalog)
	err = json.Unmarshal(result, &c)
	if err != nil {
		return nil, err
	}
	return c, err
}
func (s *Servicebroker) testConnection() error {
	resp, err := http.Get(s.Host)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *Servicebroker) instances() (*Instances, error) {
	result, err := s.getResultFromBroker("instances", "GET", "{}")
	if err != nil {
		return nil, err
	}

	var i = new(Instances)
	err = json.Unmarshal(result, &i)
	if err != nil {
		return nil, err
	}
	return i, err
}

func (s *Servicebroker) getResultFromBroker(url string, method string, json string) ([]byte, error) {
	body := strings.NewReader(json)
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.Host, url), body)
	if err != nil {
		return nil, err
	}
	if s.Username != "" {
		req.SetBasicAuth(s.Username, s.Password)
	}
	req.Header.Set("Content-Type", "application/json") //"application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	return result, err
}

func (s *Servicebroker) deleteService() {
	body := strings.NewReader(`{ "service_id":$SERVICE_ID, "plan_id":$PLAN_ID, "organization_id":$ORGANIZATION_ID }`)
	req, err := http.NewRequest("DELETE", os.ExpandEnv("$1/v2/service_instances/$2"), body)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth("password", "user")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}

func (s *Servicebroker) provision() {
	type Payload struct {
		OrganizationGUID string `json:"organization_guid"`
		PlanID           string `json:"plan_id"`
		ServiceID        string `json:"service_id"`
		SpaceGUID        string `json:"space_guid"`
		Parameters       struct {
			Parameter1 int    `json:"parameter1"`
			Parameter2 string `json:"parameter2"`
		} `json:"parameters"`
	}

	data := Payload{
	// fill struct
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", "http://username:password@broker-url/v2/service_instances/:instance_id", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("X-Broker-Api-Version", "2.10")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}
