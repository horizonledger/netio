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
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

//TODO move to implementation
func RequestReply(ntchan Ntchan, msgString string) string {

	//TODO separate namespace

	var reply_msg string
	//var reply_msg netio.Message
	//msg, _ := ParseLine(msgString)
	msg, _ := ParseLineJson(msgString)

	fmt.Sprintf("Handle cmd %v", msg.Command)

	switch msg.Command {

	case CMD_PING:
		rmsg := MessageJSON{MessageType: "REP", Command: "PONG"}
		//reply_msg := "REP PONG"
		reply_msg, _ := json.Marshal(rmsg)
		return string(reply_msg)
		//reply := HandlePing(msg)
		//msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
		//reply_msg = netio.ToJSONMessage(reply)

	case CMD_EXIT:
		//TODO close chans?
		err := ntchan.Conn.(*net.TCPConn).SetLinger(0)
		if err != nil {
			log.Printf("Error when setting linger: %s", err)
		} else {
			fmt.Println("connection closed")
			//quite all
			//ntchan.Conn.
			ntchan.quitchan <- true

		}

	case CMD_TIME:
		dt := time.Now()
		reply_msg := dt.String()
		return reply_msg

	case CMD_REGISTERALIAS:
		//TODO only pointer is set
		ntchan.Alias = "123"
		reply_msg := fmt.Sprintf("new alias %v", ntchan.Alias)
		//fmt.Printf("new alias %v", ntchan.Alias)
		return reply_msg

	//handshake
	case CMD_REGISTERPEER:
		reply_msg := "todo"
		return reply_msg

	case CMD_BALANCE:
		//TODO
		//balance := t.Mgr.State.Accounts[a]
		//fmt.Println("balance for ", a, balance, t.Mgr.State.Accounts)

		balance := 100
		balJson, _ := json.Marshal(balance)

		//rmsg := MessageJSON{MessageType: "REP", Command: "BALANCE", Data: &balJson}
		rmsg := Message{MessageType: "REP", Command: "BALANCE", Data: balJson}
		reply_msg, _ := json.Marshal(rmsg)
		return string(reply_msg)

	default:
		errormsg := "Error: not found command"
		fmt.Println(errormsg)
		xjson, _ := json.Marshal("")
		msg := Message{MessageType: REP, Command: CMD_ERROR, Data: []byte(xjson)}
		reply_msg = ToJSONMessage(msg)
	}

	return reply_msg
}
