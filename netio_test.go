package netio

import (
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
	go func() { ntchan.Reader_queue <- "REQ PING" }()

	readout := <-ntchan.REQ_in

	if readout != "REQ PING" {
		t.Error("process error")
	}

	go func() { ntchan.Reader_queue <- "REQ PING" }()

	readout2 := <-ntchan.REQ_in

	if readout2 == "REQ PING" {
		//ntchan.REP_out <- "REP PONG"
	}

	reply := RequestReply(ntchan, readout2)

	if reply != "REP PONG" {
		t.Error("process reply error ", reply)
	}

}
