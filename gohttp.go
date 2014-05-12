package goshare

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/abhishekkr/goshare/httpd"
)

/*
DB Call HTTP Handler
*/
func DBRest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	req.ParseForm()

	var (
		response_bytes []byte
		axn_status     bool
	)

	key_type, message_array := MessageArrayRest(req)

	if key_type != "" {
		switch req.Method {
		case "GET":
			response_bytes, axn_status = DBTasks("read", key_type, message_array)

		case "POST", "PUT":
			response_bytes, axn_status = DBTasks("push", key_type, message_array)

		case "DELETE":
			response_bytes, axn_status = DBTasks("delete", key_type, message_array)

		default:
			// log_this corrupt request
		}
	} // else log_this corrupt request

	if !axn_status {
		error_msg := fmt.Sprintf("FATAL Error: (DBTasks) %q \n", req.Form)
		http.Error(w, error_msg, http.StatusInternalServerError)

	} else if len(response_bytes) == 0 {
		w.Write([]byte("Success"))

	} else {
		w.Write(response_bytes)

	}
}

/*
return key_type and data as message_array identifiable by DBTasks
*/
func MessageArrayRest(req *http.Request) (string, []string) {
	req.ParseForm()
	key_type := req.FormValue("type")

	key := req.FormValue("key")
	val := req.FormValue("val")
	dbdata := req.FormValue("dbdata")

	if key != "" {
		dbdata = fmt.Sprintf("%s %s", key, val)
	}

	if key_type == "tsds" {
		timedot := fmt.Sprintf("%s %s %s %s %s %s",
			req.FormValue("year"), req.FormValue("month"), req.FormValue("day"),
			req.FormValue("hour"), req.FormValue("min"), req.FormValue("sec"))

		dbdata = fmt.Sprintf("%s %s", timedot, dbdata)
	}

	return key_type, strings.Fields(dbdata)
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

	http.HandleFunc("/db", DBRest)

	srv := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", httpuri, httpport),
		Handler:     http.DefaultServeMux,
		ReadTimeout: time.Duration(5) * time.Second,
	}

	fmt.Printf("access your goshare at http://%s:%d\n", httpuri, httpport)
	err := srv.ListenAndServe()
	fmt.Println("Game Over:", err)
}
