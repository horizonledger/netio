package netio

import (
	"bytes"
	"testing"
)

func TestBasic(t *testing.T) {

	m := Message{MessageType: "REQ", Command: "TIME"}

	if m.MessageType != "REQ" {
		t.Error("wrong type")
	}

	if m.Command != "TIME" {
		t.Error("wrong command")
	}

	m2 := ConstructMsg("REQ", "TIME", "test")

	if m2.MessageType != "REQ" {
		t.Error("wrong type")
	}

	if m2.Command != "TIME" {
		t.Error("wrong command")
	}

	res := bytes.Compare(m2.Data, []byte("test"))

	if res == 1 {
		t.Error("wrong data")
	}

	m3 := ConstructMsgBytes("REQ", "TIME", []byte("test"))

	if m3.MessageType != "REQ" {
		t.Error("wrong type")
	}

	if m3.Command != "TIME" {
		t.Error("wrong command")
	}

	res2 := bytes.Compare(m3.Data, []byte("test"))

	if res2 == 1 {
		t.Error("wrong data")
	}

	if !validCMD("TIME") {
		t.Error("invalid command")
	}
}

func TestParse(t *testing.T) {

	m, _ := ParseLine("REQ TIME")

	if m.MessageType != "REQ" {
		t.Error("parse error")
	}

	if m.Command != "TIME" {
		t.Error("parse error")
	}

}
