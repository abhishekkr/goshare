package main

import (
  "fmt"
  "flag"
  "log"
  "net/http"
  "html/template"
  "os"
  "runtime"
  "runtime/pprof"
  "time"

  "github.com/jmhodges/levigo"
  "./leveldb"
)

var (
  db *levigo.DB
  dbpath     = flag.String("dbpath", "/tmp/GO.DB", "the path to DB")
  httpport   = flag.Int("port", 9797, "what Socket PORT to run at")
  cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func HomePage(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/index.html")
  t.Execute(w, nil)
}

func ReadKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  val := abkleveldb.GetValues(req.Form["key"][0], db)
  w.Write([]byte(val))
}

func PushKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  status := abkleveldb.PushKeyVal(req.Form["key"][0], req.Form["val"][0], db)
  if status != true {
    http.Error(w, "FATAL Error", http.StatusInternalServerError)
  }
  w.Write([]byte("Success"))
}

func main() {
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

  http.HandleFunc("/", HomePage)
  http.HandleFunc("/get", ReadKey)
  http.HandleFunc("/put", PushKey)

  srv := &http.Server{
    Addr:        fmt.Sprintf(":%d", *httpport),
    Handler:     http.DefaultServeMux,
    ReadTimeout: time.Duration(5) * time.Second,
  }
  fmt.Printf("access your goshare at http://<IP>:%d\n", *httpport)
  srv.ListenAndServe()

}
