package main

import (
	"flag"
	"fmt"
	"time"

	golhashmap "github.com/abhishekkr/gol/golhashmap"
	golzmq "github.com/abhishekkr/gol/golzmq"
)

var (
	request_port01   = flag.Int("req-port01", 9797, "what Socket PORT to run at")
	request_port02   = flag.Int("req-port02", 9898, "what Socket PORT to run at")
	zmqSock          = golzmq.ZmqRequestSocket("127.0.0.1", []int{*request_port01, *request_port02})
	result, expected string
	err              error
)

// just a helper assert
func assertEqual(result interface{}, expected interface{}) {
	if result != expected {
		panic(fmt.Sprintf("[FAIL]\nExpected: '%q'\nRecieved: '%q'\n\n", expected, result))
	}
}

// for key-type default
func TestDefaultKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "default", "myname", "anon")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "anon"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "default", "myname", "anonymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "anonymous"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "default", "myname")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "default", "myname")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "default", "myname")
	expected = "Error for request sent: read default myname"
	assertEqual(result, expected)
	assertEqual(err, nil)
}

// for key-type ns
func TestNSKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname:last:first", "anon")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "myname:last:first,anon"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname:last", "ymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns", "myname", "anonymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "myname,anonymous\nmyname:last,ymous\nmyname:last:first,anon"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myname")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "myname")
	expected = "Error for request sent: read ns myname"
	assertEqual(result, expected)
	assertEqual(err, nil)
}

// for key-type tsds
func TestTSDSKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", "2014", "2", "10", "9", "8", "7", "myname:last:first", "anon")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	expected = "myname:last:first:2014:February:10:9:8:7,anon"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname:last:first")
	expected = "myname:last:first:2014:February:10:9:8:7,anon"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds", "2014", "2", "10", "9", "18", "37", "myname", "anonymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	expected = "myname:last:first:2014:February:10:9:8:7,anon\nmyname:2014:February:10:9:18:37,anonymous"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,bob\nmyemail,bob@b.com")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	expected = "myname:last:first:2014:February:10:9:8:7,anon\nmyname:2014:February:10:9:18:37,bob"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,alice\nmytxt,\"my email, bob@b.com\"")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myemail")
	expected = "myemail:2014:February:10:9:18:37,bob@b.com"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "mytxt")
	expected = "mytxt:2014:February:10:9:18:37,\"my email, bob@b.com\""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname:2014:February:10")
	expected = "myname:2014:February:10:9:18:37,alice"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myname")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "myemail")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns", "mytxt")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)
}

// for key-type now
func TestNowKeyType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "now", "myname:last:first", "anon")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_length := len(golhashmap.CSVToHashMap(result))
	expected_length := 1
	assertEqual(result_length, expected_length)
	assertEqual(err, nil)

	time.Sleep(1)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "now", "myname:last", "ymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds", "myname")
	result_length = len(golhashmap.CSVToHashMap(result))
	expected_length = 2
	assertEqual(result_length, expected_length)
	assertEqual(err, nil)
}

/* for parentNS for key-type */
func TestParentNSValType() {
	result, err = golzmq.ZmqRequest(zmqSock, "push", "ns-default-parent", "people", "myname", "anonymous")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns", "people:myname")
	expected = "people:myname,anonymous"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "ns-default-parent", "people", "myname")
	expected = "people:myname,anonymous"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "push", "tsds-csv-parent", "2014", "2", "10", "9", "18", "37", "people", "myname,bob")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default-parent", "people", "myname")
	expected = "people:myname,anonymous\npeople:myname:2014:February:10:9:18:37,bob"
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "delete", "ns-default-parent", "people", "myname")
	expected = ""
	assertEqual(result, expected)
	assertEqual(err, nil)

	result, err = golzmq.ZmqRequest(zmqSock, "read", "tsds-default-parent", "people", "myname")
	expected = "Error for request sent: read tsds-default-parent people myname"
	assertEqual(result, expected)
	assertEqual(err, nil)
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
