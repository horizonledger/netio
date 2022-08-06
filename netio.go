package netio

// functions relating to network stack
// TCP implementation currently works like this
// we read and write from the network whenever we can
// the types of read/write operations depends on the type of message flow
// any heartbeating and higher level protocols can be added on top of this
// disconnects and throttling can be added

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

func NetWrite(ntchan Ntchan, content string) (int, error) {
	//fmt.Println("NetWrite ", content)

	//since Netread READLINE uses \n add it here
	NEWLINE := '\n'
	respContent := fmt.Sprintf("%s%c", content, NEWLINE)
	//log.Println("write > ", content, respContent)
	writer := bufio.NewWriter(ntchan.Conn)
	n, err := writer.WriteString(respContent)
	if err == nil {
		err = writer.Flush()
	} else {
		fmt.Println("error writing to ", ntchan)
	}
	s := fmt.Sprintf("bytes written %d %s %s", n, ntchan.SrcName, ntchan.DestName)
	vlog(ntchan, s)
	return n, err
}

// func NetReadOpt() {
// 	//TOOD user buffer and scanner
// 	//would work like this. we expect delimiter to denote end of the message
// 	//we dont worry about large message length for now. if we do we need some expection of the
// 	//size of the message, could add this to a header of the message
// 	//conn.Read(buf[0:])
// }

func NetRead(ntchan Ntchan, delim byte) (string, error) {
	//log.Println("NetRead ", ntchan.SrcName, ntchan.DestName)
	reader := bufio.NewReader(ntchan.Conn)

	//TODO handle invalid json!

	var buffer bytes.Buffer
	for {
		//READLINE uses \n
		fmt.Println("NetRead")
		ba, isPrefix, err := reader.ReadLine()
		fmt.Println(">> ", string(ba))
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		buffer.Write(ba)
		if !isPrefix {
			break
		}

		fmt.Println("buffer ", buffer)
	}
	return buffer.String(), nil
}

func NetMsgRead(ntchan Ntchan) (string, error) {
	DELIM := byte('}')
	msg_string, err := NetRead(ntchan, DELIM)
	//msg_string = strings.Trim(msg_string, string(DELIM))
	return msg_string, err
}

func OpenConn(addr string) net.Conn {
	// Dial the remote process
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	// 	//timeoutDuration := 5 * time.Second
	// 	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	return conn
}

func OpenNtchanOut(src string, ip string, Port int, verbose bool) Ntchan {
	fulladdr := ip + ":" + strconv.Itoa(Port)
	return OpenNtchan(src, fulladdr, verbose)
}

func OpenNtchan(src string, dest string, verbose bool) Ntchan {
	conn := OpenConn(dest)
	//name := addr
	BROAD_signal := make(chan string)
	return ConnNtchan(conn, src, dest, verbose, BROAD_signal)
}
