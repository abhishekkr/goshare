package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	golassert "github.com/abhishekkr/gol/golassert"
	golhashmap "github.com/abhishekkr/gol/golhashmap"
	golzmq "github.com/abhishekkr/gol/golzmq"

	goshare "github.com/abhishekkr/goshare"
	goshare_requestor "github.com/abhishekkr/goshare/requestor"
)

var (
	request_port01   = flag.Int("req-port01", 9797, "what Socket PORT to run at")
	request_port02   = flag.Int("req-port02", 9898, "what Socket PORT to run at")
	zmqSock          = golzmq.ZmqRequestSocket("127.0.0.1", []int{*request_port01, *request_port02})
	expected, result string
	err              error
)

// for key-type default
func TestDefaultKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "default", "myname", "anon")
	fmt.Println(*request_port01)
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "myname,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default-csv", "myname")
	expected = "myname,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default-json", "[\"myname\"]")
	expected = "{\"myname\":\"anon\"}"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "default", "myname", "anonymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "myname,anonymous"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "default", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "default", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "Error for request sent: read default myname"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

// for key-type ns
func TestNSKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname:last:first", "anon")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "myname:last:first,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname:last", "ymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname", "anonymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected_arr := []string{"myname,anonymous", "myname:last,ymous", "myname:last:first,anon"}
	result_arr := strings.Split(result, "\n")
	golassert.AssertEqualStringArray(expected_arr, result_arr)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "Error for request sent: read ns myname"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

// for key-type tsds
func TestTSDSKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", "2014", "2", "10", "9", "8", "7", "myname:last:first", "anon")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "myname:last:first:2014:2:10:9:8:7:0,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	expected = "myname:last:first:2014:2:10:9:8:7:0,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname:last:first")
	expected = "myname:last:first:2014:2:10:9:8:7:0,anon"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", "2014", "2", "10", "9", "18", "37", "myname", "anonymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_arr := strings.Split(result, "\n")
	expected_arr := []string{"myname:last:first:2014:2:10:9:8:7:0,anon",
		"myname:2014:2:10:9:18:37:0,anonymous"}
	golassert.AssertEqualStringArray(expected_arr, result_arr)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,bob\nmyemail,bob@b.com")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37.0", "myname,bob\nmyemail,bob@b.com")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_arr = strings.Split(result, "\n")
	expected_arr = []string{"myname:last:first:2014:2:10:9:8:7:0,anon",
		"myname:2014:2:10:9:18:37:0,bob"}
	golassert.AssertEqualStringArray(expected_arr, result_arr)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,alice\nmytxt,\"my email, bob@b.com\"")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myemail")
	expected = "myemail:2014:2:10:9:18:37:0,bob@b.com"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "mytxt")
	expected = "mytxt:2014:2:10:9:18:37:0,\"my email, bob@b.com\""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname:2014:2:10")
	expected = "myname:2014:2:10:9:18:37:0,alice"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myemail")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "mytxt")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	expected = "Error for request sent: read tsds myname"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

// for key-type now
func TestNowKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "now", "myname:last:first", "anon")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_length := len(golhashmap.CSVToHashMap(result))
	expected_length := 1
	golassert.AssertEqual(result_length, expected_length)
	golassert.AssertEqual(err, nil)

	time.Sleep(1)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "now", "myname:last", "ymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_length = len(golhashmap.CSVToHashMap(result))
	expected_length = 2
	golassert.AssertEqual(result_length, expected_length)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "Error for request sent: read ns myname"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

/* for parentNS for key-type */
func TestParentNSValType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns-default-parent", "people", "myname", "anonymous")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "people:myname")
	expected = "people:myname,anonymous"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns-default-parent", "people", "myname")
	expected = "people:myname,anonymous"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv-parent", "2014", "2", "10", "9", "18", "37", "people", "myname,bob")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default-parent", "people", "myname")
	expected = "people:myname,anonymous\npeople:myname:2014:2:10:9:18:37:0,bob"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns-default-parent", "people", "myname")
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default-parent", "people", "myname")
	expected = "Error for request sent: read tsds-default-parent people myname"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

func TestRequestorZeromq() {
	_packet := goshare.Packet{}
	_packet.KeyType = "default"

	_packet.DBAction = "push"
	_packet.HashMap = make(golhashmap.HashMap)
	_packet.HashMap["an0n"] = "ymous"
	strPacket := string(goshare_requestor.RequestPacketBytes(&_packet))
	result, err = golzmq.ZmqRequest(zmqSock, strPacket)
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	_packet.DBAction = "read"
	_packet.KeyList = []string{"an0n"}
	strPacket = string(goshare_requestor.RequestPacketBytes(&_packet))
	result, err = golzmq.ZmqRequest(zmqSock, strPacket)
	expected = "an0n,ymous"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	_packet.DBAction = "delete"
	strPacket = string(goshare_requestor.RequestPacketBytes(&_packet))
	result, err = golzmq.ZmqRequest(zmqSock, strPacket)
	expected = ""
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)

	_packet.DBAction = "read"
	strPacket = string(goshare_requestor.RequestPacketBytes(&_packet))
	result, err = golzmq.ZmqRequest(zmqSock, strPacket)
	expected = "Error for request sent: read default-default an0n"
	golassert.AssertEqual(expected, result)
	golassert.AssertEqual(err, nil)
}

func main() {
	flag.Parse()
	fmt.Printf("client ZeroMQ REP/REQ... at %d, %d\n", *request_port01, *request_port02)

	fmt.Println("Checking out levigo based storage...")
	TestDefaultKeyType()

	fmt.Println("Checking out levigoNS based storage...")
	TestNSKeyType()

	fmt.Println("Checking out levigoTSDS based storage...")
	TestTSDSKeyType()
	TestNowKeyType()

	TestParentNSValType()
}
