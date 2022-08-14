package netio

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Optional[T any] struct {
    Defined bool
    Value   *T
}

type MessageJSON struct {
	//type of message i.e. the communications protocol
	MessageType string `json:"messagetype"`
	//Specific message command
	Command string `json:"command"`
	//any data, can be empty. gets interpreted downstream to other structs
	//https://stackoverflow.com/questions/36601367/json-field-set-to-null-vs-field-not-there
	Data json.RawMessage `json:"data,omitempty"`
	//timestamp
	//Layer *string `json:"layer,omitempty"`
}

//marshal to json, check command
func NewJSONMessage(m Message) (MessageJSON, error) {
	//fmt.Println("NewJSONMessage")
	valid := validCMD(m.Command)
	if valid {
		//		if m.Data != nil {
		return MessageJSON{
			m.MessageType,
			m.Command,
			m.Data,
		}, nil

	} else {
		fmt.Println("not valid cmd")
	}
	return MessageJSON{}, errors.New("not valid cmd")
}

func ToJSONMessage(m Message) string {
	jsonmsgtype, _ := NewJSONMessage(m)
	jsonmsg, _ := json.Marshal(jsonmsgtype)
	return string(jsonmsg)
}

func ParseLineJson(msg_string string) (Message, error) {
	//TODO parse Data
	var msgu Message
	err := json.Unmarshal([]byte(msg_string), &msgu)
	fmt.Println("error decoding json ", err, msg_string)
	var msgTypes = StrSlice{"REP", "REQ", "HEART", "BROAD", "PUB"}

	if !msgTypes.Has(msgu.MessageType) {
		fmt.Println("unsupported msg type ", msg_string)
		return msgu, errors.New("unsupported msg type")
	}

	return msgu, nil

}
