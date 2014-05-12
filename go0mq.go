package goshare

import (
	"fmt"
	"runtime"
	"strings"

	abkzeromq "github.com/abhishekkr/goshare/zeromq"
)

/* handling Read/Push/Delete tasks diversion based on task-type */
func goShareZmqRep(req_port int, rep_port int) {
	socket := abkzeromq.ZmqRep(req_port, rep_port)
	for {
		msg, _ := socket.Recv(0)
		message_array := strings.Fields(string(msg))

		axn, key_type := message_array[0], message_array[1]

		response_bytes, axn_status := DBTasks(axn, key_type, message_array[2:])

		socket.Send([]byte(response_bytes), 0)
		if !axn_status {
			fmt.Printf("Error for request sent: %s\n", msg)
		}
	}
}

/* start a Daemon communicating over 2 ports over ZMQ Rep/Req */
func GoShareZMQ(req_port int, rep_port int) {
	fmt.Printf("starting ZeroMQ REP/REQ at %d/%d\n", req_port, rep_port)
	runtime.GOMAXPROCS(runtime.NumCPU())

	goShareZmqRep(req_port, rep_port)
}
