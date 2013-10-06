package main

import (
  "fmt"
  "flag"

  "../zeromq"
)

var (
  req_port    = flag.Int("req-port", 9797, "what Socket PORT to run at")
  rep_port    = flag.Int("rep-port", 9898, "what Socket PORT to run at")
)

func main(){
  flag.Parse()
  fmt.Printf("client ZeroMQ REP/REQ... at %d, %d", req_port, rep_port)

  abkzeromq.ZmqReq(*req_port, *rep_port, "myname", "anon")
  abkzeromq.ZmqReq(*req_port, *rep_port, "myname")
  abkzeromq.ZmqReq(*req_port, *rep_port, "myname", "anonymous")
  abkzeromq.ZmqReq(*req_port, *rep_port, "myname")
}
