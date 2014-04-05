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
    _axn, _type, _key := message_array[0], message_array[1], message_array[2]
    return_value := ""
    _axn_result := false

    switch _axn {
      case "read":
        return_value = GetValTask(_type, _key)
        if return_value != "" { _axn_result = true }

      case "push":
        if PushKeyValByType(_type, message_array){ _axn_result = true }

      case "delete":
        if DelKeyTask(_type, _key){ _axn_result = true }
    }

    socket.Send([]byte(return_value), 0)
    if ! _axn_result {
      fmt.Printf("Error for request sent: %s\n", msg)
    }
  }
}


/* start a Daemon communicating over 2 ports over ZMQ Rep/Req */
func GoShareZMQ(req_port int, rep_port int){
  fmt.Printf("starting ZeroMQ REP/REQ at %d/%d\n", req_port, rep_port)
  runtime.GOMAXPROCS(runtime.NumCPU())

  goShareZmqRep(req_port, rep_port)
}
