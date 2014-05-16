package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	golhashmap "github.com/abhishekkr/gol/golhashmap"
)

var (
	httphost = flag.String("host", "127.0.0.1", "what Host to run at")
	httpport = flag.Int("port", 9999, "what Socket PORT to connect")
)

/* KeyVal Support */

func GetURL(host string, port int, key_type, key string) string {
	return fmt.Sprintf("http://%s:%d/get?type=%s&key=%s", host, port, key_type, key)
}

func PutURL(host string, port int, key_type, key, val string) string {
	return fmt.Sprintf("http://%s:%d/put?type=%s&key=%s&val=%s",
		host, port, key_type, key, val)
}

func DelURL(host string, port int, key_type, key string) string {
	return fmt.Sprintf("http://%s:%d/del?type=%s&key=%s", host, port, key_type, key)
}

func TSDSPutURL(host string, port int, key, val, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?key=%s&val=%s&type=tsds&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, key, val, year, month, day, hr, min, sec)
}

func MultiValPutURL(host string, port int, key_type, val_type, multi_value string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=%s-%s", host, port, multi_value, key_type, val_type)
}

func MultiTSDSPutURL(host string, port int, val_type, multi_value, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=tsds-%s&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, multi_value, val_type, year, month, day, hr, min, sec)
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

func assertEqual(body interface{}, expected_body interface{}) {
	if body != expected_body {
		panic(fmt.Sprintf("[FAIL]\nExpected: '%s'\nRecieved: '%s'\n\n", expected_body, body))
	}
}

func main() {
	flag.Parse()

	//jsonmap := golhashmap.GetHashMapEngine("json")
	csvmap := golhashmap.GetHashMapEngine("csv")

	assertEqual(HttpGet(PutURL(*httphost, *httpport, "default", "myname", "anon")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "default", "myname")), "anon")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "default", "myname", "anonymous")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "default", "myname")), "anonymous")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "default", "myname")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "default", "myname")), "FATAL Error: (DBTasks) map[\"type\":[\"default\"] \"key\":[\"myname\"]]\n")

	assertEqual(HttpGet(PutURL(*httphost, *httpport, "ns", "myname:last:first", "anon")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last:first")), "myname:last:first,anon\n")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "ns", "myname:last", "ymous")), "Success")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "ns", "myname", "anonymous")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last")), "myname:last,ymous\nmyname:last:first,anon\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "ns", "myname")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last")), "FATAL Error: (DBTasks) map[\"type\":[\"ns\"] \"key\":[\"myname:last\"]]\n")

	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname:last", "ymous", "2014", "2", "10", "1", "1", "1")), "Success")
	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "anonymous", "2014", "2", "10", "9", "8", "7")), "Success")
	assertEqual(HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "untitled", "2014", "2", "10", "0", "1", "7")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:last")), "myname:last:2014:February:10:1:1:1,ymous\n")

	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "myname")), "myname:last:2014:February:10:1:1:1,ymous\nmyname:2014:February:10:9:8:7,anonymous\nmyname:2014:February:10:0:1:7,untitled\n")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:2014:February:10")), "myname:2014:February:10:9:8:7,anonymous\nmyname:2014:February:10:0:1:7,untitled\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "myname")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "myname")), "FATAL Error: (DBTasks) map[\"type\":[\"tsds\"] \"key\":[\"myname\"]]\n")

	assertEqual(HttpGet(MultiTSDSPutURL(*httphost, *httpport, "csv", "yourname:last:first,trudy", "2014", "2", "10", "1", "1", "1")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "yourname:last:first:2014:February")), "yourname:last:first:2014:February:10:1:1:1,trudy\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "yourname")), "Success")

	assertEqual(HttpGet(MultiValPutURL(*httphost, *httpport, "ns", "csv", "yname:frend:first,monica")), "Success")
	assertEqual(HttpGet(MultiValPutURL(*httphost, *httpport, "ns", "csv", "yname:frend:second,lolita")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:frend")), "yname:frend:first,monica\nyname:frend:second,lolita\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "yname")), "Success")

	assertEqual(HttpGet(MultiValPutURL(*httphost, *httpport, "ns", "csv", "yname:frend:first,monica%0D%0Ayname:frend:second,lolita%0D%0Auname:frend:second,juno%0D%0A")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:frend")), "yname:frend:first,monica\nyname:frend:second,lolita\n")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "uname:frend")), "uname:frend:second,juno\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "yname")), "Success")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "uname")), "Success")

	assertEqual(HttpGet(MultiValPutURL(*httphost, *httpport, "ns", "json", "{\"power:first\":\"yay\",\"power:second\":\"way\",\"rower:second\":\"kay\"}")), "Success")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "power")), "power:first,yay\npower:second,way\n")
	assertEqual(HttpGet(GetURL(*httphost, *httpport, "tsds", "rower")), "rower:second,kay\n")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "power")), "Success")
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "rower")), "Success")

	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "yname")), "Success")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "now", "yname:last:first", "zodiac")), "Success")
	assertEqual(len(csvmap.ToHashMap(HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:last:first")))), 1)
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "yname")), "Success")

	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "myname")), "Success")
	assertEqual(HttpGet(PutURL(*httphost, *httpport, "now", "myname:last:first", "ripper")), "Success")
	assertEqual(len(csvmap.ToHashMap(HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:last:first")))), 1)
	assertEqual(HttpGet(DelURL(*httphost, *httpport, "tsds", "myname")), "Success")
}
