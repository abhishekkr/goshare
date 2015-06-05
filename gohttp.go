package goshare

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	golhashmap "github.com/abhishekkr/gol/golhashmap"
	goltime "github.com/abhishekkr/gol/goltime"
	"github.com/abhishekkr/goshare/httpd"
)

/*
DBRest enables HTTP to formulate DBTasks call from HTTP Request.
*/
func DBRest(httpMethod string, w http.ResponseWriter, req *http.Request) {
	var (
		dbAction      string
		responseBytes []byte
		axnStatus     bool
	)

	switch httpMethod {
	case "GET":
		dbAction = "read"

	case "POST", "PUT":
		dbAction = "push"

	case "DELETE":
		dbAction = "delete"

	default:
		// log_this corrupt request
		return
	}

	packet := PacketFromHTTPRequest(dbAction, req)
	if packet.DBAction != "" {
		responseBytes, axnStatus = DBTasksOnPacket(packet)
		DBRestResponse(w, req, responseBytes, axnStatus)
	}
}

/*
DBRestResponse sends proper response back to client based on success, data or error.
*/
func DBRestResponse(w http.ResponseWriter, req *http.Request, responseBytes []byte, axnStatus bool) {
	if !axnStatus {
		errorMsg := fmt.Sprintf("FATAL Error: (DBTasks) %q", req.Form)
		http.Error(w, errorMsg, http.StatusInternalServerError)

	} else if len(responseBytes) == 0 {
		w.Write([]byte("Success"))

	} else {
		w.Write(responseBytes)

	}
}

/*
PacketFromHTTPRequest return Packet identifiable by DBTasksOnAction.
*/
func PacketFromHTTPRequest(dbAction string, req *http.Request) Packet {
	packet := Packet{}
	packet.HashMap = make(golhashmap.HashMap)
	packet.DBAction = dbAction

	req.ParseForm()
	taskType := req.FormValue("type")
	if taskType == "" {
		taskType = "default"
	}
	packet.TaskType = taskType
	taskTypeTokens := strings.Split(taskType, "-")
	packet.KeyType = taskTypeTokens[0]
	if packet.KeyType == "tsds" && packet.DBAction == "push" {
		packet.TimeDot = goltime.TimestampFromHTTPRequest(req)
	}

	if len(taskTypeTokens) > 1 {
		packet.ValType = taskTypeTokens[1]

		if len(taskTypeTokens) == 3 {
			thirdTokenFeatureHTTP(&packet, req)
		}
	}

	dbdata := req.FormValue("dbdata")
	key := req.FormValue("key")
	if key != "" {
		dbdata = fmt.Sprintf("%s\n%s", key, req.FormValue("val"))
	} else if dbdata == "" {
		return Packet{}
	}
	decodeData(&packet, strings.Split(dbdata, "\n"))

	return packet
}

/*
thirdTokenFeatureHTTP  handles third token in taskType.
*/
func thirdTokenFeatureHTTP(packet *Packet, req *http.Request) {
	parentNS := req.FormValue("parentNS")
	if parentNS != "" {
		packet.ParentNS = parentNS
	}
}

/*
DBRestHandler handles DB Call for HTTP Request Method at '/db'.
*/
func DBRestHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	req.ParseForm()

	DBRest(req.Method, w, req)
}

/*
GetReadKey HTTP GET DB-GET call handler at '/get'.
*/
func GetReadKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	req.ParseForm()

	DBRest("GET", w, req)
}

/*
GetPushKey HTTP GET DB-POST call handler at '/put'.
*/
func GetPushKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	req.ParseForm()

	DBRest("POST", w, req)
}

/*
GetDeleteKey HTTP GET DB-POST call handler at 'del'.
*/
func GetDeleteKey(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	req.ParseForm()

	DBRest("DELETE", w, req)
}

/*
GoShareHTTP handles all valid HTTP Requests for DBTasks, documentation
and playground(WIP).
*/
func GoShareHTTP(httpuri string, httpport int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/", abkhttpd.Index)
	http.HandleFunc("/quickstart", abkhttpd.QuickStart)
	http.HandleFunc("/help-http", abkhttpd.HelpHTTP)
	http.HandleFunc("/help-zmq", abkhttpd.HelpZMQ)
	http.HandleFunc("/concept", abkhttpd.Concept)
	http.HandleFunc("/status", abkhttpd.Status)

	http.HandleFunc("/db", DBRestHandler)
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
