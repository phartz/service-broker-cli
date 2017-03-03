package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSBErrorTexts(t *testing.T) {
	bytes := []byte("{\"description\":\"Insufficient parameters for un-provisioning instance!\"}")
	var sbError = new(SBError)
	tempErr := json.Unmarshal(bytes, &sbError)

	var err error
	if (sbError != nil && (sbError.Error != "" || sbError.Description != "")) || tempErr != nil {
		err = errors.New(fmt.Sprintf("%s / %s", sbError.Description, sbError.Error))
		return
	}

	assertNotNil(t, err, "")
}

func TestJSONInsertion(t *testing.T) {
	testPayload := BindPayload{ServiceID: "expected"}

	bytes := []byte("{\"parameters\":{\"description\":\"Insufficient parameters for un-provisioning instance!\"}}")
	var data = new(CustomPaylod)
	err := json.Unmarshal(bytes, &data)
	assertIsNil(t, err, "")

	testPayload.Parameters = data.Parameters

	payloadBytes, err := json.Marshal(testPayload)
	assertIsNil(t, err, "")

	assertIsTrue(t, strings.Contains(string(payloadBytes), "expected"), "")
}
