package tests

import (
	"fmt"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func AssertNotNil(t *testing.T, a interface{}, message string) {
	if a != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is nil", a)
	}
	t.Fatal(message)
}

func AssertIsNil(t *testing.T, a interface{}, message string) {
	if a == nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not nil", a)
	}
	t.Fatal(message)
}

func AssertIsTrue(t *testing.T, a bool, message string) {
	if a {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not true", a)
	}
	t.Fatal(message)
}

func AssertIsFalse(t *testing.T, a bool, message string) {
	if a {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v is not true", a)
	}
	t.Fatal(message)
}
