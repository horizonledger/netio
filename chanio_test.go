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
	NetConnectorSetupMockEcho(ntchan)

	//put message in request queue
	//ntchan.REQ_in <- string(jm)
	ntchan.Reader_queue <- "test"
	fmt.Println("????....")

	//go func() {
	fmt.Println("wait")
	req := <-ntchan.Writer_queue
	fmt.Println("request received ", req)
	ntchan.REP_out <- "ok"

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
