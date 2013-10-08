package goshare

import (
  "fmt"
  "runtime"

  "github.com/jmhodges/levigo"

  "github.com/abhishekkr/goshare/leveldb"
  "github.com/abhishekkr/goshare/zeromq"
)

func ReadKey(key string) string{
  return abkleveldb.GetValues(key, db)
}

func PushKey(key string, val string) bool{
  return abkleveldb.PushKeyVal(key, val, db)
}

func GoShareZMQ(leveldb *levigo.DB, req_port int, rep_port int){
  db = leveldb
  fmt.Printf("starting ZeroMQ REP/REQ at %d/%d\n", req_port, rep_port)
  runtime.GOMAXPROCS(runtime.NumCPU())

  abkzeromq.ZmqRep(req_port, rep_port, ReadKey, PushKey)
}
