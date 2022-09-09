package netio

import (
	//"encoding/json"
	"fmt"
	"testing"
)

func TestBasicNetio(t *testing.T) {
	fmt.Println("run")

	//establish network
	ntchan := ConnNtchanStub("test", "testout")
	go NetConnectorSetupMockEcho(ntchan)

	//put message in request queue
	ntchan.Reader_queue <- "test"

	fmt.Println("wait")
	req := <-ntchan.Writer_queue

	if req != "echo: test" {
		t.Error("wrong echo message")
	}

}

// func RequestReplyLoop(ntchan Ntchan) {
// 	for {
// 		msg := <-ntchan.REQ_in
// 		m := MessageJSON{MessageType: "REP", Command: "PONG"}
// 		jm, _ := json.Marshal(m)
// 		reply := jm
// 		ntchan.REP_out <- reply
// 	}
// }

// // TODO
// func TestNetioPingPong(t *testing.T) {

// 	m := MessageJSON{MessageType: "REQ", Command: "PING"}
// 	jm, _ := json.Marshal(m)

// 	ntchan.Reader_queue <- string(jm)

// 	go RequestReplyLoop(ntchan)
// }
