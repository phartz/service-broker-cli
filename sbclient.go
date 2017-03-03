package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type SBClient struct {
	Credentials
}

func (s *SBClient) SetCredentials(c Credentials) {
	s.Host = c.Host
	s.Username = c.Username
	s.Password = c.Password
}

func (s *SBClient) Catalog() (*Catalog, error) {
	result, _, _, err := s.getResultFromBroker("v2/catalog", "GET", "{}")
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

func (s *SBClient) TestConnection() error {
	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Println("")
		fmt.Printf("\tTest host %s\n", s.Host)
	}

	resp, err := http.Get(s.Host)

	if err != nil {
		if os.Getenv("SB_TRACE") == "ON" {
			fmt.Println("")
			fmt.Printf("\tError: %s\n", err.Error())
		}
		return err
	}

	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Println("")
		fmt.Println("\tStatus: OK\n")
	}
	defer resp.Body.Close()
	return nil
}

func (s *SBClient) LastState(instanceId string) (*LastState, error) {
	result, _, _, err := s.getResultFromBroker(fmt.Sprintf("v2/service_instances/%s/last_operation", instanceId), "GET", "{}")
	if err != nil {
		return nil, err
	}

	var l = new(LastState)
	err = json.Unmarshal(result, &l)
	if err != nil {
		return nil, err
	}
	return l, err
}

func (s *SBClient) Instances() (*Instances, error) {
	result, _, _, err := s.getResultFromBroker("instances", "GET", "{}")
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

func (s *SBClient) getResultFromBroker(url string, method string, jsonStr string) (bytes []byte, statusCode int, status string, err error) {
	statusCode = 0
	status = ""
	bytes = nil

	body := strings.NewReader(jsonStr)
	target := fmt.Sprintf("%s/%s", s.Host, url)

	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Println("")
		fmt.Printf("\tRequest to: %s\n", target)
		fmt.Printf("\tMethod:     %s\n", method)
		fmt.Printf("\tBody:\n\t%s\n", jsonStr)
	}

	req, err := http.NewRequest(method, target, body)
	if err != nil {
		return
	}
	if s.Username != "" {
		req.SetBasicAuth(s.Username, s.Password)
	}
	req.Header.Set("Content-Type", "application/json") //"application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if os.Getenv("SB_TRACE") == "ON" {
			fmt.Printf("\tError:\n\t%s\n", err.Error())
		}
		return
	}

	defer resp.Body.Close()

	status = resp.Status
	statusCode = resp.StatusCode

	bytes, err = ioutil.ReadAll(resp.Body)
	if os.Getenv("SB_TRACE") == "ON" {
		fmt.Printf("\tStatus: %d/%s\n", resp.StatusCode, resp.Status)
		fmt.Printf("\tBody:\n\t%s\n", string(bytes))
	}

	var sbError = new(SBError)
	tempErr := json.Unmarshal(bytes, &sbError)
	if (sbError != nil && (sbError.Error != "" || sbError.Description != "")) || tempErr != nil {
		err = errors.New(fmt.Sprintf("%s / %s", sbError.Description, sbError.Error))
		return
	}

	return
}

func (s *SBClient) Deprovision(data *BindPayload, instanceID string) error {
	_, statusCode, status, err := s.getResultFromBroker(fmt.Sprintf("v2/service_instances/%s?_id=%s&plan_id=%s", instanceID, data.ServiceID, data.PlanID), "DELETE", "{}")
	if err != nil {
		return err
	}

	if statusCode >= 200 && statusCode <= 202 {
		return nil
	}

	return errors.New(fmt.Sprintf("Deprovision failure code: %d/%s", statusCode, status))
}

func (s *SBClient) UpdateService(instanceID string) error {
	_, statusCode, status, err := s.getResultFromBroker(fmt.Sprintf("v2/service_instances/%s", instanceID), "PATCH", "{}")
	if err != nil {
		return err
	}

	if statusCode >= 200 && statusCode <= 202 {
		return nil
	}

	return errors.New(fmt.Sprintf("Deprovision failure code: %d/%s", statusCode, status))
}

func (s *SBClient) Provision(data *ProvisonPayload, instanceID string) error {
	payloadBytes, err := json.Marshal(data)

	_, statusCode, status, err := s.getResultFromBroker(fmt.Sprintf("v2/service_instances/%s", instanceID), "PUT", string(payloadBytes))
	if err != nil {
		return err
	}

	if statusCode >= 200 && statusCode <= 202 {
		return nil
	}

	return errors.New(fmt.Sprintf("Provision failure code: %d/%s", statusCode, status))
}

func (s *SBClient) Bind(data *BindPayload, instanceID string, bindID string) (string, error) {
	payloadBytes, err := json.Marshal(data)

	bytes, _, _, err := s.getResultFromBroker(fmt.Sprintf("/v2/service_instances/%s/service_bindings/%s", instanceID, bindID), "PUT", string(payloadBytes))
	if err != nil {
		return "", err
	}

	var sbError = new(SBError)
	err = json.Unmarshal(bytes, &sbError)

	if err != nil {
		return "", err
	}

	if sbError != nil && sbError.Error != "" {
		return "", errors.New(sbError.Error)
	}

	return string(bytes), nil
}

func (s *SBClient) UnBind(data *BindPayload, instanceID string, bindID string) error {
	_, statusCode, status, err := s.getResultFromBroker(fmt.Sprintf("v2/service_instances/%s/service_bindings/%s?service_id=%s&plan_id=%s", instanceID, bindID, data.ServiceID, data.PlanID), "DELETE", "{}")
	if err != nil {
		return err
	}

	if statusCode == 200 {
		return nil
	}

	return errors.New(fmt.Sprintf("Unbind failure code: %d/%s", statusCode, status))
}
