package goshare

import (
  "fmt"
  "net/http"
  "runtime"
  "time"

  "github.com/jmhodges/levigo"
  abkleveldb "github.com/abhishekkr/levigoNS/leveldb"

  "github.com/abhishekkr/goshare/httpd"
)

func GetReadKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  val := abkleveldb.GetVal(req.Form["key"][0], db)
  w.Write([]byte(val))
}

func GetPushKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  status := abkleveldb.PushKeyVal(req.Form["key"][0], req.Form["val"][0], db)
  if status != true {
    http.Error(w, "FATAL Error", http.StatusInternalServerError)
  }
  w.Write([]byte("Success"))
}

func GetDeleteKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  status := abkleveldb.DelKey(req.Form["key"][0], db)
  if status != true {
    http.Error(w, "FATAL Error", http.StatusInternalServerError)
  }
  w.Write([]byte("Success"))
}

func GoShareHTTP(leveldb *levigo.DB, httpuri string, httpport int) {
  db = leveldb
  runtime.GOMAXPROCS(runtime.NumCPU())

  http.HandleFunc("/", abkhttpd.F1)
  http.HandleFunc("/help-http", abkhttpd.HelpHTTP)
  http.HandleFunc("/help-zmq", abkhttpd.HelpZMQ)
  http.HandleFunc("/status", abkhttpd.Status)

  http.HandleFunc("/get", GetReadKey)
  http.HandleFunc("/put", GetPushKey)
  http.HandleFunc("/del", GetDeleteKey)

  srv := &http.Server{
    Addr:        fmt.Sprintf("%s:%d", httpuri, httpport),
    Handler:     http.DefaultServeMux,
    ReadTimeout: time.Duration(5) * time.Second,
  }

  fmt.Printf("access your goshare at http://%s:%d\n", httpuri, httpport)
  err := srv.ListenAndServe()
  fmt.Println("Game Over:", err) 
}
