package netio

import "time"

// only output the message read
func NetConnectorSetupMock(ntchan Ntchan) {

	vlog(ntchan, "NetConnectorSetup "+ntchan.SrcName+" "+ntchan.DestName)

	//process of reads
	go ReadProcessor(ntchan)
	//processor of REQ_out REP_out
	go WriteProcessor(ntchan)

	//go RequestLoop(ntchan)

}

// read from reader queue and echo all messages back
// TODO in mock space
func ReadProcessorEcho(ntchan Ntchan) {

	for {
		logmsgd(ntchan, "ReadProcessor", "loop")
		msgString := <-ntchan.Reader_queue
		logmsgd(ntchan, "ReadProcessor", msgString)

		if len(msgString) > 0 {
			logmsgc(ntchan, ntchan.SrcName, "ReadProcessor", msgString) //, ntchan.Reader_processed)

			reply := "echo: " + msgString
			ntchan.Writer_queue <- reply

		}
	}

}

// echo pipeline
// TODO separate namespace
func NetConnectorSetupEcho2(ntchan Ntchan, RequestReplyF func(Ntchan, string) string) {

	vlog(ntchan, "NetConnectorSetup "+ntchan.SrcName+" "+ntchan.DestName)

	// read_loop_time := 800 * time.Millisecond
	// read_time_chan := 300 * time.Millisecond
	// write_loop_time := 300 * time.Millisecond

	//reads from the actual "physical" network
	go ReadLoop(ntchan)
	//process of reads
	go ReadProcessorEcho(ntchan)
	//write to network whatever is in writer queue
	go WriteLoop(ntchan, 300*time.Millisecond)

	go func() {
		for {
			msg := <-ntchan.BROAD_in
			ntchan.BROAD_signal <- msg
		}
	}()

	//TODO
	//go WriteProducer(ntchan)
}

// simple echo net
func MockNetConnectorSetupEcho(ntchan Ntchan) {

	//read loop
	go func() {
		d := 300 * time.Millisecond
		for {
			//vlog(ntchan, "iter ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
			msg, _ := NetMsgRead(ntchan)
			//currently can be empty or len, shoudl fix one style
			if len(msg) > 0 { //&& msg != EMPTY_MSG {
				vlog(ntchan, "ntwk read => "+msg)
				logmsgc(ntchan, ntchan.SrcName, "ReadLoop", msg)
				vlog(ntchan, "put on Reader queue "+msg)
				reply := "echo >>> " + msg
				ntchan.Reader_queue <- reply
			}

			time.Sleep(d)
		}
	}()

	//echo back
	go func() {
		for {
			msgString := <-ntchan.Reader_queue
			NetWrite(ntchan, msgString)
		}
	}()
}

func NetConnectorSetupMockEcho(ntchan Ntchan) {

	for {
		msgString := <-ntchan.Reader_queue

		if len(msgString) > 0 {
			reply := "echo: " + msgString
			ntchan.Writer_queue <- reply
		}
	}

}
