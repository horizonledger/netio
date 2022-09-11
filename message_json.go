package netio

import (
	"encoding/json"
	"errors"
	"fmt"
)

// type Optional[T any] struct {
//     Defined bool
//     Value   *T
// }

// func (o *Optional[T]) UnmarshalJSON(data []byte) error {
// 	o.Defined = true
// 	return json.Unmarshal(data, &o.Value)
// }

type MessageJSON struct {
	//type of message i.e. the communications protocol
	MessageType string `json:"messagetype"`
	//Specific message command
	Command string `json:"command"`
	//any data, can be empty. gets interpreted downstream to other structs
	Data json.RawMessage `json:"data,omitempty"`
	//Data Optional[json.RawMessage] `json:"data,omitempty"`
	//timestamp
}

// marshal to json, check command
func NewJSONMessage(m Message) (MessageJSON, error) {
	//fmt.Println("NewJSONMessage")
	valid := validCMD(m.Command)
	if valid {
		return MessageJSON{
			m.MessageType,
			m.Command,
			m.Data,
		}, nil

		// if m.Data != nil {
		// 	return MessageJSON{
		// 		m.MessageType,
		// 		m.Command,
		// 		m.Data,
		// 	}, nil
		// } else {
		// 	return MessageJSON{
		// 		m.MessageType,
		// 		m.Command,
		// 		nil,
		// 	}, nil
		// }

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

func ParseLineJson(msg_string string) (MessageJSON, error) {
	//TODO parse Data
	var msgu MessageJSON
	err := json.Unmarshal([]byte(msg_string), &msgu)
	if err != nil {
		fmt.Println("error decoding json ", err, "msg: ", msg_string)
		return msgu, errors.New("error decoding json")
	}
	var msgTypes = StrSlice{"REP", "REQ", "HEART", "BROAD", "PUB"}

	if !msgTypes.Has(msgu.MessageType) {
		fmt.Println("unsupported msg type ", msg_string)
		return msgu, errors.New("unsupported msg type")
	}

	return msgu, nil

}
