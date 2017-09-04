package tests

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/phartz/service-broker-cli/sbcli"
)

func TestSBErrorTexts(t *testing.T) {
	bytes := []byte("{\"description\":\"Insufficient parameters for un-provisioning instance!\"}")
	var sbError = new(sbcli.SBError)
	tempErr := json.Unmarshal(bytes, &sbError)

	var err error
	if (sbError != nil && (sbError.Error != "" || sbError.Description != "")) || tempErr != nil {
		err = fmt.Errorf("%s / %s", sbError.Description, sbError.Error)
		return
	}

	AssertNotNil(t, err, "")
}

func TestJSONInsertion(t *testing.T) {
	testPayload := sbcli.BindPayload{ServiceID: "expected"}

	bytes := []byte("{\"parameters\":{\"description\":\"Insufficient parameters for un-provisioning instance!\"}}")
	var data = new(sbcli.CustomPaylod)
	err := json.Unmarshal(bytes, &data)
	AssertIsNil(t, err, "")

	testPayload.Parameters = data.Parameters

	payloadBytes, err := json.Marshal(testPayload)
	AssertIsNil(t, err, "")
	AssertIsTrue(t, strings.Contains(string(payloadBytes), "expected"), "")
}
