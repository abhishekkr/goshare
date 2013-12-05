package abkzeromq

import (
  "fmt"
  "strings"
  zmq "github.com/alecthomas/gozmq"
)

type READ func(key string) string
type PUSH func(key string, val string) bool
type DELETE func(key string) bool

func ZmqRep(req_port int, rep_port int, read READ, push PUSH, del DELETE) {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%d", req_port))
  socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%d", rep_port))

  fmt.Printf("ZMQ REQ/REP Daemon at port %d and %d\n", req_port, rep_port)
  for {
    msg, _ := socket.Recv(0)
    msg_arr := strings.Fields(string(msg))
    return_value := ""
    if msg_arr[0] == "read" {
      return_value = read(msg_arr[1])
    } else if msg_arr[0] == "push" {
      _value := strings.Join(msg_arr[2:], " ")
      if push(msg_arr[1], _value) {
        return_value = _value
      }
    } else if msg_arr[0] == "delete" {
      if del(msg_arr[1]) {
        return_value = msg_arr[1]
      }
    } else {
      fmt.Printf("unhandled request sent: %s", msg)
    }
    socket.Send([]byte(return_value), 0)
    fmt.Println("Got: [ ", string(msg), " ]; Sent: [ ", return_value, " ]")
  }
}

func ZmqReq(req_port int, rep_port int, dat ...string) []byte{
  fmt.Printf("ZMQ REQ/REP Client at port %d and %d\n", req_port, rep_port)
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", req_port))
  socket.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", rep_port))

  var msg string
  msg = strings.Join(dat, " ")
  socket.Send([]byte(msg), 0)
  response, _ := socket.Recv(0)
  fmt.Printf("msg: %s\nresponse: %s\n\n", msg, response)
  return response
}
