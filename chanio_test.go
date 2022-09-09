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

	//fmt.Println("request received ", req)
	//ntchan.REP_out <- "ok"

	//}()

	//read from reply
	// readout := <-ntchan.REP_out
	// fmt.Println(readout)

	// m = MessageJSON{MessageType: "REP", Command: "PONG"}
	// jm, _ = json.Marshal(m)

	// if readout != string(jm) {
	// 	t.Error("process reply error ", readout)
	// }

}
