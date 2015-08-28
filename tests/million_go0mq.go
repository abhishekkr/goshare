package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	golzmq "github.com/abhishekkr/gol/golzmq"
)

var (
	httphost         = flag.String("host", "127.0.0.1", "what Host to run at")
	httpport         = flag.Int("port", 9999, "what Socket PORT to connect")
	request_port01   = flag.Int("req-port01", 9797, "what Socket PORT to run at")
	request_port02   = flag.Int("req-port02", 9898, "what Socket PORT to run at")
	zmqSock          = golzmq.ZmqRequestSocket("127.0.0.1", []int{*request_port01, *request_port02})
	expected, result string
	err              error
	daycount         = 0
	days             = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// return Get URL for task_type, key
func GetURL(host string, port int, key_type, key string) string {
	return fmt.Sprintf("http://%s:%d/get?type=%s&key=%s", host, port, key_type, key)
}

// return Push URL for task_type, key, val
func PutURL(host string, port int, key_type, key, val string) string {
	return fmt.Sprintf("http://%s:%d/put?type=%s&key=%s&val=%s",
		host, port, key_type, key, val)
}

// return Delete URL for task_type, key
func DelURL(host string, port int, key_type, key string) string {
	return fmt.Sprintf("http://%s:%d/del?type=%s&key=%s", host, port, key_type, key)
}

// return Push URL for TSDS type key, val, time-elements
func TSDSPutURL(host string, port int, key, val, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s&type=tsds&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, key, val, year, month, day, hr, min, sec)
}

// makes HTTP call for given URL and returns response body
func HttpGet(url string) (int, string) {
	resp, err := http.Get(url)
	if err != nil {
		return 404, "Error: " + url + " failed for HTTP GET"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Printf("Url: %s; with\n result:\n%s\n", url, string(body))
	return resp.StatusCode, string(body)
}

func everysecond(yy, mm, dd int) {
	for hr := 1; hr <= 24; hr++ {
		for min := 1; min <= 60; min++ {
			for sec := 1; sec <= 60; sec++ {
				//result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", string(yy), string(mm), string(dd), string(hr), string(min), string(sec), "dayState", string(daycount))
				//result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", string(yy), string(mm), string(dd), string(hr), string(min), string(sec), "dayState")
				//result, err = golzmq.ZmqRequest(zmqSock, "delete", "tsds", string(yy), string(mm), string(dd), string(hr), string(min), string(sec), "dayState")

				HttpGet(TSDSPutURL(*httphost, *httpport, "dayState", string(daycount), string(yy), string(mm), string(dd), string(hr), string(min), string(sec)))
				//_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname"))
				//_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "myname"))

				daycount++
				//fmt.Println(result, err)
			}
		}
	}
}

func everyday(fromYear, toYear int) {
	for yy := fromYear; yy <= toYear; yy++ {
		for mm := 1; mm <= 12; mm++ {
			for dd := 1; dd <= days[mm]; dd++ {
				everysecond(yy, mm, dd)
			}
		}
	}
}

func main() {
	flag.Parse()
	fmt.Printf("client ZeroMQ REP/REQ... at %d, %d\n", *request_port01, *request_port02)
	//everyday(2001, 2001)
	everysecond(2015, 6, 21)
	fmt.Println(daycount)
	/*
		//for i := 0; i < 1000000; i++ {
		for i := 0; i < 10000; i++ {
			_i := fmt.Sprintf("%d", i)
			result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", "2014", "2", "10", "9", "18", "37", "people"+_i, "bob"+_i)
			//result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "people"+_i+",bob"+_i)
			//result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default", "people"+_i)
			//result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns-default", "people"+_i)

				//_i := fmt.Sprintf("people:group%d", i)
				//result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv-parent", "2014", "2", "10", "9", "18", "37", _i, "myname,bob")
				//result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default-parent", _i, "myname")
				//result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns-default-parent", _i, "myname")

				//ZmqKeyVal("push", "default", "k_"+_i, _i)
				//ZmqKey("read", "default", "k_"+_i)
				//ZmqKey("delete", "default", "k_"+_i)
		}
	*/
}
