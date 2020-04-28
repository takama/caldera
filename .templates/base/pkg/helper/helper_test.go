package helper_test

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"{{[ .Project ]}}/pkg/helper"
)

func TestLogE(t *testing.T) {
	var data bytes.Buffer

	log.SetOutput(&data)
	log.SetFlags(0)
	helper.LogE("new message", nil)

	expected := ""

	if data.String() != expected {
		t.Errorf("Expected message %s, got %s", expected, data.String())
	}

	helper.LogE("message", errors.New("Error"))

	expected = "message: Error\n"

	if data.String() != expected {
		t.Errorf("Expected message %s, got %s", expected, data.String())
	}

	helper.LogE("new message", nil)
}

func TestLogF(t *testing.T) {
	helper.LogF("message", nil)

	if t.Failed() {
		t.Error("Expected success, got failed")
	}
}
