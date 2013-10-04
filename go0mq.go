package main

import (
  "fmt"
  "flag"
  "log"
  "os"
  "runtime"
  "runtime/pprof"
  "time"

  "github.com/jmhodges/levigo"
  "./leveldb"
  "./zeromq"
)

var (
  db *levigo.DB
  dbpath      = flag.String("dbpath", "/tmp/GO.DB", "the path to DB")
  req_port    = flag.Int("req-port", 9797, "what Socket PORT to run at")
  rep_port    = flag.Int("rep-port", 9898, "what Socket PORT to run at")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func ReadKey(key string) string{
  return abkleveldb.GetValues(key, db)
}

func PushKey(key string, val string) bool{
  return abkleveldb.PushKeyVal(key, val, db)
}

func main(){
  fmt.Println("starting ZeroMQ REP/REQ...")
  runtime.GOMAXPROCS(runtime.NumCPU())

  flag.Parse()
  db = abkleveldb.CreateDB(*dbpath)
  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
      log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    go func() {
      time.Sleep(100 * time.Second)
      pprof.StopCPUProfile()
    }()
  }

  abkzeromq.ZmqRep(*req_port, *rep_port, ReadKey, PushKey)
}
