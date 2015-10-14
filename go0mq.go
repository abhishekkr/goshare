package goshare

import (
	"fmt"
	"runtime"
	"strings"

	zmq "github.com/alecthomas/gozmq"

	golzmq "github.com/abhishekkr/gol/golzmq"
)

/*
goShareZmqRep handles Read/Push/Delete tasks diversion based on task-type.
*/
func goShareZmqRep(socket *zmq.Socket) {
	var errResponse string
	for {
		msg, _ := socket.Recv(0)
		messageArray := strings.Fields(string(msg))
		responseBytes, axnStatus := DBTasks(messageArray)

		if axnStatus {
			socket.Send([]byte(responseBytes), 0)
		} else {
			errResponse = fmt.Sprintf("Error for request sent: %s", msg)
			socket.Send([]byte(errResponse), 0)
		}
	}
}

/*
GoShareZMQ starts a Daemon communicating of provided array ports over ZMQ Reply.
*/
func GoShareZMQ(ip string, replyPorts []int) {
	fmt.Printf("starting ZeroMQ REP/REQ at %v\n", replyPorts)
	runtime.GOMAXPROCS(runtime.NumCPU())

	socket := golzmq.ZmqReplySocket(ip, replyPorts)
	goShareZmqRep(socket)
}
