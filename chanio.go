package netio

// network communication layer (netio)

// netio -> semantics of channels
// TCP/IP -> golang net

// TODO
// create a channel wrapper struct
// which has a flag if its in or out flow
// see whitepaper for details
// type Nchain {
// c chan string
// inflow }

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

type Ntchan struct {
	//TODO is only single connection
	Conn net.Conn
	//flag for if the peer has done a handshake
	registered bool
	//Name     string
	SrcName      string //TODO doc
	DestName     string
	Alias        string
	Reader_queue chan string
	Writer_queue chan string
	//inflow
	REQ_in       chan string
	REP_in       chan string
	BROAD_in     chan string
	BROAD_signal chan string
	SEND_in      chan string
	HEART_in     chan string
	//outflow
	REP_out chan string
	REQ_out chan string
	//BROAD_out chan string

	SEND_out  chan string
	HEART_out chan string
	//
	//
	PUB_out chan string
	SUB_out chan string
	//PUB_time_quit chan int
	verbose       bool
	lastheartbeat time.Time
	quitchan      chan bool
	// SUB_request_out   chan string
	// SUB_request_in    chan string
	// UNSUB_request_out chan string
	// UNSUB_request_in  chan string

	// Reader_processed int
	// Writer_processed int
}

func vlog(ntchan Ntchan, s string) {
	//fmt.Println(s)
	if ntchan.verbose {
		log.Println("vlog ", s)
	}
	//log.Println("vlog ", s)
}

func logmsgd(ntchan Ntchan, src string, msg string) {
	s := fmt.Sprintf("[%s] ### %v", src, msg)
	vlog(ntchan, s)
}

func logmsgc(ntchan Ntchan, name string, src string, msg string) {
	s := fmt.Sprintf("%s [%s] ### %v", name, src, msg)
	vlog(ntchan, s)
}

func logmsge(ntchan Ntchan, name string, src string, dest string, msg string) {
	s := fmt.Sprintf("%s [%s] [%s] ### %v", name, src, dest, msg)
	vlog(ntchan, s)
}

// func ConnNtchan(conn net.Conn, SrcName string, DestName string, verbose bool, BROAD_signal chan string) Ntchan {
func ConnNtchan(conn net.Conn, SrcName string, DestName string, verbose bool) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.BROAD_in = make(chan string)
	//ntchan.BROAD_signal = BROAD_signal

	ntchan.HEART_in = make(chan string)
	ntchan.HEART_out = make(chan string)

	ntchan.REQ_out = make(chan string)
	ntchan.PUB_out = make(chan string)
	//ntchan.PUB_time_quit = make(chan int)
	// ntchan.Reader_processed = 0
	// ntchan.Writer_processed = 0
	ntchan.Conn = conn
	ntchan.SrcName = SrcName
	ntchan.DestName = DestName

	ntchan.Alias = fmt.Sprintf("peer-%v", rand.Intn(1000))

	return ntchan
}

// for testing
func ConnNtchanStub(SrcName string, DestName string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.PUB_out = make(chan string)
	//ntchan.BROAD_in = make(chan string)
	//ntchan.PUB_time_quit = make(chan int)
	//ntchan.Reader_processed = 0
	//ntchan.Writer_processed = 0
	ntchan.SrcName = SrcName
	ntchan.DestName = DestName

	return ntchan
}

// all major processes to operate
// RequestReplyF func(Ntchan, string) string
func NetConnectorSetup(ntchan Ntchan) {

	vlog(ntchan, "NetConnectorSetup "+ntchan.SrcName+" "+ntchan.DestName)

	// read_loop_time := 800 * time.Millisecond
	// read_time_chan := 300 * time.Millisecond
	// write_loop_time := 300 * time.Millisecond

	//quite all channel
	ntchan.quitchan = make(chan bool)

	//reads from the actual "physical" network
	go ReadLoop(ntchan)
	//process of reads in X_in chans
	go ReadProcessor(ntchan)
	//processor of X_out chans
	go WriteProcessor(ntchan)
	//write to network whatever is in writer queue
	go WriteLoop(ntchan, 300*time.Millisecond)

	//go HeartbeatPub(ntchan)

	//node logic

	//go RequestReplyLoop(ntchan, RequestReplyF)

	go HeartBeatProcess(ntchan)

	//TODO
	//go WriteProducer(ntchan)
}

func NetConnectorSetupEcho(ntchan Ntchan) {

	vlog(ntchan, "NetConnectorSetup "+ntchan.SrcName+" "+ntchan.DestName)

	// read_loop_time := 800 * time.Millisecond
	// read_time_chan := 300 * time.Millisecond
	// write_loop_time := 300 * time.Millisecond

	//quite all channel
	ntchan.quitchan = make(chan bool)

	//reads from the actual "physical" network
	go ReadLoop(ntchan)
	//process of reads in X_in chans
	go ReadProcessor(ntchan)
	//processor of X_out chans
	go WriteProcessor(ntchan)
	//write to network whatever is in writer queue
	go WriteLoop(ntchan, 300*time.Millisecond)

	//go HeartbeatPub(ntchan)

	//node logic
	//go RequestReplyLoop()

	go HeartBeatProcess(ntchan)

	//TODO
	//go WriteProducer(ntchan)
}

func BroadSignalSetup() {

	//TODO unused
	// go func() {
	// 	for {
	// 		msg := <-ntchan.BROAD_in
	// 		fmt.Printf("received broadcast %s %s\n", msg, ntchan.Alias)
	// 		ntchan.BROAD_signal <- msg
	// 		//signal back to main
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		ntchan.BROAD_signal <- fmt.Sprintf("test %v", ntchan.SrcName)
	// 		time.Sleep(5000 * time.Millisecond)
	// 	}
	// }()

}

// Net
func WriteLoop(ntchan Ntchan, d time.Duration) {
	//msg_writer_total := 0
	for {
		//log.Println("loop writer")
		//TODO!
		//

		//take from channel and write
		msg := <-ntchan.Writer_queue
		vlog(ntchan, "writeloop "+msg)
		NetWrite(ntchan, msg)
		//logmsg(ntchan.Name, "WriteLoop", msg, msg_writer_total)
		//NetworkWrite(ntchan, msg)

		time.Sleep(d)
		//msg_writer_total++
	}
}

// read from network and put in reader channel queue
func ReadLoop(ntchan Ntchan) {
	vlog(ntchan, "init ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
	d := 300 * time.Millisecond
	//msg_reader_total := 0
	for {
		select {
		case <-ntchan.quitchan:
			fmt.Println("quit readloop")
			return

		default:
			//read from network and put in channel
			vlog(ntchan, "iter ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
			msg, err := NetMsgRead(ntchan)

			//handle close connection
			if err != nil {
				return
			}

			//handle cases
			//currently can be empty or len, shoudl fix one style
			if len(msg) > 0 { //&& msg != EMPTY_MSG {
				vlog(ntchan, "ntwk read => "+msg)
				logmsgc(ntchan, ntchan.SrcName, "ReadLoop", msg)
				vlog(ntchan, "put on Reader queue "+msg)
				//put in the queue to process
				ntchan.Reader_queue <- msg
			}

			time.Sleep(d)
			//fix: need ntchan to be a pointer
			//msg_reader_total++
		}
	}

}

func HeartBeatProcess(ntchan Ntchan) {
	//TODO set status if changed
	// need to check how long we didnt receive a message from a node
	//separate loop
	//if now - ntchan.lastheartbeat > 2 seconds
	//terminate connection

	for {
		msg := <-ntchan.HEART_in
		vlog(ntchan, "heart  "+msg)
		//TODO set last received heartbeat
		// layout := "2006-01-02T15:04:05.000Z"
		// str := "2014-11-12T11:45:26.371Z"
		// t, err := time.Parse(layout, str)

		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}

//read from reader queue and process by forwarding to right channel

// TODO remove
func ReadProcessorJson(ntchan Ntchan) {

	for {
		logmsgd(ntchan, "ReadProcessor", "loop")
		msgString := <-ntchan.Reader_queue
		logmsgd(ntchan, "ReadProcessor", msgString)

		if len(msgString) > 0 {
			logmsgc(ntchan, ntchan.SrcName, "ReadProcessor", msgString) //, ntchan.Reader_processed)
			//TODO try catch

			//msg := EdnParseMessageMap(msgString)
			msg := FromJSON(msgString)

			if msg.MessageType == REQ {

				//msg_string := MsgString(msg)
				//msg_string := EdnConstructMsgMapS(msg)
				logmsgd(ntchan, "ReadProcessor", "REQ_in")

				//TODO!
				ntchan.REQ_in <- msgString

			} else if msg.MessageType == REP {
				//TODO!
				//msg_string := MsgString(msg)
				//msg_string := EdnConstructMsgMapS(msg)

				logmsgc(ntchan, "ReadProcessor", "REP_in", msgString)
				ntchan.REP_in <- msgString

				// x := <-ntchan.REP_in
				// vlog(ntchan, "x "+x)
			}
			//  else if msg.MessageType == SUB {

			// }

			//else unkown message type

			//ntchan.Reader_processed++
			//log.Println(" ", ntchan.Reader_processed, ntchan)
		}
	}

}

// read from reader queue and echo all messages back
// TODO move to peer&client
// remove
func ReadProcessor(ntchan Ntchan) {

	for {
		select {
		case <-ntchan.quitchan:
			fmt.Println("quit ReadProcessor")
			return

		default:
			logmsgd(ntchan, "ReadProcessor", "loop")
			msgString := <-ntchan.Reader_queue
			logmsgd(ntchan, "ReadProcessor", msgString)

			if len(msgString) > 0 {
				//logmsge(ntchan, ntchan.SrcName, ntchan.DestName, "ReadProcessor", msgString) //, ntchan.Reader_processed)
				logmsge(ntchan, ntchan.Alias, ntchan.DestName, "ReadProcessor", msgString) //, ntchan.Reader_processed)

				//msg := FromJSON(msgString)
				//msg, err := ParseLine(msgString)
				msg, err := ParseLineJson(msgString)

				if err == nil {
					logmsgc(ntchan, ntchan.SrcName, "ReadProcessor Msg", msg.MessageType)

					switch msg.MessageType {
					case "REQ":
						ntchan.REQ_in <- msgString

					case "REP":
						ntchan.REP_in <- msgString

					case "BROAD":
						ntchan.BROAD_in <- msgString

					case "HEART":
						ntchan.HEART_in <- msgString

						//TODO
						//case PUB
						//case SUB
					default:
						//reply := "echo >>> " + msgString
						//handle unknown command
						ntchan.Writer_queue <- "unknown command"
					}

				} else {

					CloseOut(ntchan)
				}

			}
		}
	}

}

// process from higher level chans into write queue
func WriteProcessor(ntchan Ntchan) {
	//TODO have a list and iterate over it
	for {

		select {
		case msg := <-ntchan.HEART_out:
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.REP_out:
			//read from REQ_out
			//log.Println("[Writeprocessor]  REP_out", msg)
			logmsgc(ntchan, "WriteProcessor", "REP_out", msg)
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.REQ_out:
			logmsgc(ntchan, "Writeprocessor", "REQ_out ", msg)
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.PUB_out:
			ntchan.Writer_queue <- msg
			//PUB?
		}
	}
}

func CloseOut(ntchan Ntchan) {
	ntchan.Conn.Close()

	//TODO currently no reply on fault messages
	//reply := "error parsing >>> " + msgString + " (" + ntchan.DestName + ")"
	//ntchan.Writer_queue <- reply
	// reply := "error parsing >>> " + msgString + " (" + ntchan.DestName + ")"
	// ntchan.Writer_queue <- reply
}

func HeartbeatPub(ntchan Ntchan) {
	fmt.Print("HeartbeatPub")
	for range time.Tick(time.Second * 1) {
		msg := "heartbeat"
		vlog(ntchan, "put on Writer_queue "+msg)
		ntchan.Writer_queue <- msg
	}
}
