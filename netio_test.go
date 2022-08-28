package netio

import (
	"encoding/json"
	"testing"
)

func TestRequest(t *testing.T) {

	ntchan := ConnNtchanStub("test", "testout")
	if ntchan.SrcName != "test" {
		t.Error("setup error")
	}

	go func() { ntchan.REQ_in <- "test" }()

	readout := <-ntchan.REQ_in

	if readout != "test" {
		t.Error("parse error")
	}

}

func TestProcesser(t *testing.T) {

	ntchan := ConnNtchanStub("test", "testout")
	NetConnectorSetupMock(ntchan)

	//go func() { ntchan.Reader_queue <- "REQ PING" }()
	m := MessageJSON{MessageType: "REQ", Command: "PING"}
	jm, _ := json.Marshal(m)

	ntchan.Reader_queue <- string(jm)

	readout := <-ntchan.REQ_in

	if readout != "{\"messagetype\":\"REQ\",\"command\":\"PING\"}" {
		t.Error("process error ", readout)
	}

	// reply := RequestReply(ntchan, readout2)

	// if reply != "{\"messagetype\":\"REP\",\"command\":\"PONG\"}" {
	// 	t.Error("process reply error ", reply)
	// }

}
