package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	httphost = flag.String("host", "127.0.0.1", "what Host to run at")
	httpport = flag.Int("port", 9999, "what Socket PORT to connect")
)

/* KeyVal Support */

func GetURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/get?key=%s", host, port, key)
}

func PutURL(host string, port int, key string, val string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s",
		host, port, key, val)
}

func DelURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/del?key=%s", host, port, key)
}

/* KeyVal NS Support */

func NSGetURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/get?key=%s&type=ns", host, port, key)
}

func NSPutURL(host string, port int, key string, val string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s&type=ns",
		host, port, key, val)
}

func NSDelURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/del?key=%s&type=ns", host, port, key)
}

/* KeyVal TSDS Support */

func TSDSGetURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/get?key=%s&type=tsds", host, port, key)
}

func TSDSPutURL(host string, port int, key, val, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s&type=tsds&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, key, val, year, month, day, hr, min, sec)
}

func NowTSDSPutURL(host string, port int, key string, val string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s&type=now",
		host, port, key, val)
}

func CSVNSPutURL(host string, port int, csv_value string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=ns-csv",
		host, port, csv_value)
}

func CSVTSDSPutURL(host string, port int, csv_value, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=tsds-csv&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, csv_value, year, month, day, hr, min, sec)
}

func TSDSDelURL(host string, port int, key string) string {
	return fmt.Sprintf("http://%s:%d/del?key=%s&type=tsds", host, port, key)
}

func HttpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "Error: " + url + " failed for HTTP GET"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("Url: %s; with\n result:\n%s\n", url, string(body))
	return string(body)
}

func assertEqual(body string, expected_body string) {
	if body != expected_body {
		panic(fmt.Sprintf("[FAIL]\nExpected: '%s'\nRecieved: '%s'\n\n", expected_body, body))
	}
}

func main() {
	flag.Parse()

	assertEqual(HttpGet(PutURL(*httphost, *httpport, "myname", "anon")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "myname")), "anon")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "myname", "anonymous")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "myname")), "anonymous")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "myname")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "myname")), "FATAL Error: (DBTasks) map[\"key\":[\"myname\"]]\n")

	assertEqual(HttpGet(NSPutURL(*httphost, *httpport, "myname:last:first", "anon")), "Success")
	assertEqual(HttpGet(NSGetURL(*httphost, *httpport, "myname:last:first")), "myname:last:first,anon\n")
	assertEqual(HttpGet(NSPutURL(*httphost, *httpport, "myname:last", "ymous")), "Success")
	assertEqual(HttpGet(NSPutURL(*httphost, *httpport, "myname", "anonymous")), "Success")
	assertEqual(HttpGet(NSGetURL(*httphost, *httpport, "myname:last")), "myname:last,ymous\nmyname:last:first,anon\n")
	assertEqual(HttpGet(NSDelURL(*httphost, *httpport, "myname")), "Success")
	assertEqual(HttpGet(NSGetURL(*httphost, *httpport, "myname:last")), "FATAL Error: (DBTasks) map[\"key\":[\"myname:last\"] \"type\":[\"ns\"]]\n")

	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname:last", "ymous", "2014", "2", "10", "1", "1", "1")), "Success")
	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "anonymous", "2014", "2", "10", "9", "8", "7")), "Success")
	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "untitled", "2014", "2", "10", "0", "1", "7")), "Success")
	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "myname:last")), "myname:last:2014:February:10:1:1:1,ymous\n")

	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "myname")), "myname:last:2014:February:10:1:1:1,ymous\nmyname:2014:February:10:9:8:7,anonymous\nmyname:2014:February:10:0:1:7,untitled\n")
	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "myname:2014:February:10")), "myname:2014:February:10:9:8:7,anonymous\nmyname:2014:February:10:0:1:7,untitled\n")
	assertEqual(HttpGet(TSDSDelURL(*httphost, *httpport, "myname")), "Success")
	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "myname")), "FATAL Error: (DBTasks) map[\"key\":[\"myname\"] \"type\":[\"tsds\"]]\n")

	assertEqual(HttpGet(CSVTSDSPutURL(*httphost, *httpport, "yourname:last:first,trudy", "2014", "2", "10", "1", "1", "1")), "Success")
	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "yourname:last:first:2014:February")), "yourname:last:first:2014:February:10:1:1:1,trudy\n")
	assertEqual(HttpGet(TSDSDelURL(*httphost, *httpport, "yourname")), "Success")

	assertEqual(HttpGet(CSVNSPutURL(*httphost, *httpport, "yname:frend:first,monica")), "Success")
	assertEqual(HttpGet(CSVNSPutURL(*httphost, *httpport, "yname:frend:second,lolita")), "Success")
	assertEqual(HttpGet(TSDSGetURL(*httphost, *httpport, "yname:frend")), "yname:frend:first,monica\nyname:frend:second,lolita\n")
	assertEqual(HttpGet(TSDSDelURL(*httphost, *httpport, "yname")), "Success")

	HttpGet(NowTSDSPutURL(*httphost, *httpport, "yname:last:first", "anon"))
	HttpGet(TSDSGetURL(*httphost, *httpport, "yname:last:first"))
	HttpGet(NowTSDSPutURL(*httphost, *httpport, "myname:last:first", "anon"))
	HttpGet(TSDSGetURL(*httphost, *httpport, "myname:last:first:2014:February"))
	HttpGet(TSDSDelURL(*httphost, *httpport, "myname"))
}
