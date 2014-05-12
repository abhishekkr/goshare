package main

import (
	"flag"
	"fmt"

	abkzeromq "github.com/abhishekkr/goshare/zeromq"
)

var (
	req_port = flag.Int("req-port", 9797, "what Socket PORT to run at")
	rep_port = flag.Int("rep-port", 9898, "what Socket PORT to run at")
)

func main() {
	flag.Parse()
	fmt.Printf("client ZeroMQ REP/REQ... at %d, %d", req_port, rep_port)

	fmt.Println("Checking out levigo based storage...")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "default", "myname", "anon")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "default", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "default", "myname", "anonymous")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "default", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "delete", "default", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "delete", "default", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "default", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "default", "myname", "anon")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "default", "myname")

	fmt.Println("Checking out levigoNS based storage...")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "ns", "myname:last:first", "anon")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "ns", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "ns", "myname:last", "ymous")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "ns", "myname", "anonymous")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "ns", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "delete", "ns", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "ns", "myname")

	fmt.Println("Checking out levigoTSDS based storage...")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "now", "myname:last:first", "anon")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "tsds", "2014", "2", "10", "9", "8", "7", "myname:last:first", "anon")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "now", "myname:last", "ymous")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "tsds", "2014", "2", "10", "9", "18", "37", "myname", "anonymous")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "tsds", "2014", "2", "10", "5", "28", "57", "myname", "untitles")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,bob\nmyemail,bob@b.com")
	abkzeromq.ZmqReq(*req_port, *rep_port, "push", "tsds-csv", "2014", "2", "10", "9", "18", "37", "myname,alice\nmytxt,\"my email, bob@b.com\"")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "myname:2014:February:10")
	abkzeromq.ZmqReq(*req_port, *rep_port, "delete", "tsds", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "myname")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "myemail")
	abkzeromq.ZmqReq(*req_port, *rep_port, "read", "tsds", "mytxt")
}
