package goshare

import (
  "fmt"
  "net/http"
  "runtime"
  "time"

  "github.com/abhishekkr/goshare/httpd"
)

func GetReadKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  task_type := req.Form["type"]

  if len(task_type) > 0 {
    _get_val := GetValTask(task_type[0])
  }

  if len(keys) > 0 {
    val := _get_val(req.Form["key"][0], db)
  }
  w.Write([]byte(val))
}

func GetPushKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  vals := req.Form["val"]
  task_type := req.Form["type"]

  if len(task_type) > 0 {
    _push_keyval := PushKeyValTask(task_type[0])
  }

  if len(keys) > 0 && len(vals) > 0 {
    status := _push_keyval(keys[0], vals[0])
    if status != true {
      http.Error(w, "FATAL Error", http.StatusInternalServerError)
    }
  }
  w.Write([]byte("Success"))
}

func GetDeleteKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  task_type := req.Form["type"]

  if len(task_type) > 0 {
    _del_key := DelKeyTask(task_type[0])
  }

  if len(keys) > 0 {
    status := _del_key(keys[0])
    if status != true {
      http.Error(w, "FATAL Error", http.StatusInternalServerError)
    }
  }
  w.Write([]byte("Success"))
}

func GoShareHTTP(httpuri string, httpport int) {
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
