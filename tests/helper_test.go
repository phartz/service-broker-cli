package tests

import (
	"testing"

	"github.com/phartz/service-broker-cli/sbcli"
)

func TestCleanTargetUIR(t *testing.T) {
	AssertEqual(t, sbcli.CleanTargetURI("10.0.0.1"), "http://10.0.0.1:3000", "")
	AssertEqual(t, sbcli.CleanTargetURI("10.0.0.1:3000"), "http://10.0.0.1:3000", "")
	AssertEqual(t, sbcli.CleanTargetURI("https://10.0.0.1"), "https://10.0.0.1:3000", "")
	AssertEqual(t, sbcli.CleanTargetURI("http://10.0.0.1"), "http://10.0.0.1:3000", "")
	AssertEqual(t, sbcli.CleanTargetURI("https://10.0.0.1:3001"), "https://10.0.0.1:3001", "")
	AssertEqual(t, sbcli.CleanTargetURI("HTTP://10.0.0.1:3001"), "HTTP://10.0.0.1:3001", "")
}
