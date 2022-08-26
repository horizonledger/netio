package netio

import (
	"encoding/json"
	"fmt"
	"testing"
)

//go test -run=TestBasicencoding
func TestBasicencoding(t *testing.T) {

	//m := MessageJSON{MessageType: "REQ", Command: "TIME"}
	m := MessageJSON{MessageType: "REQ", Command: "TIME"}

	jm, err := json.Marshal(m)

	if err != nil {
		fmt.Println("error: ", err)
		t.Error("error")
	}

	if string(jm) != "{\"messagetype\":\"REQ\",\"command\":\"TIME\"}" {
		t.Error("encode json error ", string(jm))
	}

}

func TestBasicencodingDataString(t *testing.T) {

	d := []byte(`"test"`)
	m := MessageJSON{MessageType: "REQ", Command: "TIME", Data: (*json.RawMessage)(&d)}

	jm, err := json.Marshal(m)

	if err != nil {
		fmt.Println("error: ", err)
		t.Error("error")
	}

	if string(jm) != "{\"messagetype\":\"REQ\",\"command\":\"TIME\",\"data\":\"test\"}" {
		t.Error("encode json error ", string(jm))
	}

}

type Data struct {
	Name string
	Id   int
}

func TestBasicencodingDataInt(t *testing.T) {

	tmp := Data{"test", 2}
	b, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println("Error ", err.Error())
	}
	raw := json.RawMessage(b)
	m := MessageJSON{MessageType: "REQ", Command: "TIME", Data: &raw}

	jm, err := json.Marshal(m)

	if err != nil {
		fmt.Println("error: ", err)
		t.Error("error")
	}

	if string(jm) != "{\"messagetype\":\"REQ\",\"command\":\"TIME\",\"data\":{\"Name\":\"test\",\"Id\":2}}" {
		t.Error("encode json error ", string(jm))
	}

}
