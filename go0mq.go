package goshare

import (
  "fmt"
  "runtime"

  "github.com/jmhodges/levigo"
  abkleveldb "github.com/abhishekkr/levigoNS/leveldb"

  "github.com/abhishekkr/goshare/zeromq"
)

func ReadKey(key string) string{
  return abkleveldb.GetVal(key, db)
}

func PushKey(key string, val string) bool{
  return abkleveldb.PushKeyVal(key, val, db)
}

func DeleteKey(key string) bool{
  return abkleveldb.DelKey(key, db)
}

func GoShareZMQ(leveldb *levigo.DB, req_port int, rep_port int){
  db = leveldb
  fmt.Printf("starting ZeroMQ REP/REQ at %d/%d\n", req_port, rep_port)
  runtime.GOMAXPROCS(runtime.NumCPU())

  abkzeromq.ZmqRep(req_port, rep_port, ReadKey, PushKey, DeleteKey)
}
