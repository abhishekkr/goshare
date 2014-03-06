package goshare


import (
  "fmt"
  "net/http"
  "runtime"
  "time"

  "github.com/abhishekkr/goshare/httpd"
)


/*
Get Records
Multiple Keys can be passed in HTTP Request,
but one request shall be just for one TaskType (default, ns, tsds, tsds-now)
*/
func GetReadKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  task_type := req.Form["type"]

  if len(task_type) == 0 {
    task_type = make([]string, 1)
    task_type[0] = "default"
  }

  w.Write( []byte(GetValTask(task_type[0], keys[0])) )
}


/*
Push Records
Only one Key,Val can passed in HTTP Request,
and one request shall be just for one TaskType (default, ns, tsds, tsds-now)
*/
func GetPushKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  vals := req.Form["val"]
  task_type := req.Form["type"]
  status := false

  if len(task_type) == 0 {
    task_type = make([]string, 1)
    task_type[0] = "default"
  }
  if task_type[0] == "tsds" {
    message_array := []string{
      req.Form["year"][0],
      req.Form["month"][0],
      req.Form["day"][0],
      req.Form["hour"][0],
      req.Form["min"][0],
      req.Form["sec"][0],
      vals[0],
    }
    status = PushKeyMsgArrayTSDS(keys[0], message_array)
  } else {
    status = PushKeyValTask(task_type[0], keys[0], vals[0])
  }

  if status { http.Error(w, "FATAL Error", http.StatusInternalServerError) }
  w.Write([]byte("Success"))
}


/*
Delete Records
Multiple Keys can passed in HTTP Request,
but one request shall be just for one TaskType (default, ns, tsds)
*/
func GetDeleteKey(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  req.ParseForm()
  keys := req.Form["key"]
  task_type := req.Form["type"]

  if len(task_type) == 0 {
    task_type = make([]string, 1)
    task_type[0] = "default"
  }

  for _, key := range keys {
    if ! DelKeyTask(task_type[0], key) {
      http.Error(w, "FATAL Error", http.StatusInternalServerError)
    }
  }
  w.Write([]byte("Success"))
}


/*
GoShare Handler for HTTP Requests
*/
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
