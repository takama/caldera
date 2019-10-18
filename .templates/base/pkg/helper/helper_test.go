package helper

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

func TestLogE(t *testing.T) {
	var data bytes.Buffer

	log.SetOutput(&data)
	log.SetFlags(0)
	LogE("new message", nil)

	expected := ""

	if data.String() != expected {
		t.Errorf("Expected message %s, got %s", expected, data.String())
	}

	LogE("message", errors.New("Error"))

	expected = "message: Error\n"

	if data.String() != expected {
		t.Errorf("Expected message %s, got %s", expected, data.String())
	}

	LogE("new message", nil)
}

func TestLogF(t *testing.T) {
	LogF("message", nil)

	if t.Failed() {
		t.Error("Expected success, got failed")
	}
}
