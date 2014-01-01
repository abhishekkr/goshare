package goshare

import (
  "fmt"
  "runtime"

  "github.com/jmhodges/levigo"

  "github.com/abhishekkr/goshare/zeromq"
)

func ReadKey(key string) string{
  return GetVal(key)
}

func PushKey(key string, val string) bool{
  return PushKeyVal(key, val)
}

func DeleteKey(key string) bool{
  return DelKey(key)
}

func GoShareZMQ(req_port int, rep_port int){
  fmt.Printf("starting ZeroMQ REP/REQ at %d/%d\n", req_port, rep_port)
  runtime.GOMAXPROCS(runtime.NumCPU())

  abkzeromq.ZmqRep(req_port, rep_port, ReadKey, PushKey, DeleteKey)
}
