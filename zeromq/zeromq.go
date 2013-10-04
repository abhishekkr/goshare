package abkzeromq

import (
  "fmt"
  "strings"
  zmq "github.com/alecthomas/gozmq"
)

type READ func(key string) string
type PUSH func(key string, val string) bool

func ZmqRep(req_port int, rep_port int, read READ, push PUSH) {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%d", req_port))
  socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%d", rep_port))

  fmt.Printf("ZMQ REQ/REP Daemon at port %d and %d\n", req_port, rep_port)
  for {
    msg, _ := socket.Recv(0)
    fmt.Println("Got: %s", string(msg))
    msg_arr := strings.Fields(string(msg))
    if len(msg_arr) == 1 {
      read(msg_arr[0])
    } else if len(msg_arr) == 2 {
      push(msg_arr[0], msg_arr[1])
    }
    socket.Send(msg, 0)
  }
}

func ZmqReq(req_port int, rep_port int, dat ...string) {
  fmt.Printf("ZMQ REQ/REP Client at port %d and %d\n", req_port, rep_port)
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", req_port))
  socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", rep_port))

  var msg string
  msg = strings.Join(dat, " ")
  socket.Send([]byte(msg), 0)
  fmt.Printf("msg: %s\n", msg)
  socket.Recv(0)
}
