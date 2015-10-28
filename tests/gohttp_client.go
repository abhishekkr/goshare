package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	golassert "github.com/abhishekkr/gol/golassert"
	golhashmap "github.com/abhishekkr/gol/golhashmap"
)

var (
	httphost = flag.String("host", "127.0.0.1", "what Host to run at")
	httpport = flag.Int("port", 9999, "what Socket PORT to connect")
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

// return Push URL for multi-val-type on task-type and multi-val
func MultiValPutURL(host string, port int, task_type, multi_value string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=%s", host, port, multi_value, task_type)
}

// return Push URL for TSDS multi-val-type on task-type, val-type and multi-val
func MultiTSDSPutURL(host string, port int, val_type, multi_value, year, month, day, hr, min, sec string) string {
	return fmt.Sprintf("http://%s:%d/put?dbdata=%s&type=tsds-%s&year=%s&month=%s&day=%s&hour=%s&min=%s&sec=%s",
		host, port, multi_value, val_type, year, month, day, hr, min, sec)
}

// append url values for parentNS and return URL
func URLAppendParentNS(url string, parentNS string) string {
	return fmt.Sprintf("%s&parentNS=%s", url, parentNS)
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

// for default key-type
func testDefaultKeyType() {
	_, body := HttpGet(PutURL(*httphost, *httpport, "default", "myname", "anon"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "default", "myname"))
	golassert.AssertEqual(body, "myname,anon")

	_, body = HttpGet(PutURL(*httphost, *httpport, "default", "myname", "anonymous"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "default", "myname"))
	golassert.AssertEqual(body, "myname,anonymous")

	_, body = HttpGet(DelURL(*httphost, *httpport, "default", "myname"))
	golassert.AssertEqual(body, "Success")

	code, _ := HttpGet(GetURL(*httphost, *httpport, "default", "myname"))
	golassert.AssertEqual(code, 500)

	fmt.Println("No panic for 'default' key  type.")
}

// for ns key-type
func testNamespaceKeyType() {
	_, body := HttpGet(PutURL(*httphost, *httpport, "ns", "myname:last:first", "anon"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last:first"))
	golassert.AssertEqual(body, "myname:last:first,anon")

	_, body = HttpGet(PutURL(*httphost, *httpport, "ns", "myname:last", "ymous"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(PutURL(*httphost, *httpport, "ns", "myname", "anonymous"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last"))
	strings2 := strings.Split(body, "\n")
	strings2_expected := []string{"myname:last,ymous",
		"myname:last:first,anon"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	_, body = HttpGet(DelURL(*httphost, *httpport, "ns", "myname"))
	golassert.AssertEqual(body, "Success")

	code, _ := HttpGet(GetURL(*httphost, *httpport, "ns", "myname:last"))
	golassert.AssertEqual(code, 500)

	fmt.Println("No panic for 'namespace' key  type.")
}

// for tsds key-type
func testTSDSKeyType() {
	_, body := HttpGet(TSDSPutURL(*httphost, *httpport, "myname:last", "ymous", "2014", "2", "10", "1", "1", "1"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "anonymous", "2014", "2", "10", "9", "8", "7"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(TSDSPutURL(*httphost, *httpport, "myname", "untitled", "2014", "2", "10", "0", "1", "7"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:last"))
	golassert.AssertEqual(body, "myname:last:2014:2:10:1:1:1:0,ymous")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname"))
	strings3 := strings.Split(body, "\n")
	strings3_expected := []string{"myname:last:2014:2:10:1:1:1:0,ymous",
		"myname:2014:2:10:9:8:7:0,anonymous",
		"myname:2014:2:10:0:1:7:0,untitled"}
	golassert.AssertEqualStringArray(strings3, strings3_expected)

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:2014:2:10"))
	strings2 := strings.Split(body, "\n")
	strings2_expected := []string{"myname:2014:2:10:9:8:7:0,anonymous",
		"myname:2014:2:10:0:1:7:0,untitled"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "myname"))
	golassert.AssertEqual(body, "Success")

	code, _ := HttpGet(GetURL(*httphost, *httpport, "tsds", "myname"))
	golassert.AssertEqual(code, 500)

	fmt.Println("No panic for 'timeseries' key  type.")
}

// for now key-type
func testNowKeyType() {
	csvmap := golhashmap.GetHashMapEngine("csv")

	_, body := HttpGet(DelURL(*httphost, *httpport, "tsds", "yname"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(PutURL(*httphost, *httpport, "now", "yname:last:first", "zodiac"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:last:first"))
	golassert.AssertEqual(len(csvmap.ToHashMap(body)), 1)

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "yname"))
	golassert.AssertEqual(body, "Success")

	code, body := HttpGet(GetURL(*httphost, *httpport, "tsds", "yname"))
	golassert.AssertEqual(code, 500)

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "myname"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(PutURL(*httphost, *httpport, "now", "myname:last:first", "ripper"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname:last:first"))
	golassert.AssertEqual(len(csvmap.ToHashMap(body)), 1)

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "myname"))
	golassert.AssertEqual(body, "Success")

	code, _ = HttpGet(GetURL(*httphost, *httpport, "tsds", "myname"))
	golassert.AssertEqual(code, 500)

	fmt.Println("No panic for 'timeseries now' key  type.")
}

// for csv val-type
func testForCSV() {
	_, body := HttpGet(MultiTSDSPutURL(*httphost, *httpport, "csv", "yourname:last:first,trudy", "2014", "2", "10", "1", "1", "1"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "yourname:last:first:2014:2"))
	golassert.AssertEqual(body, "yourname:last:first:2014:2:10:1:1:1:0,trudy")

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "yourname"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(MultiValPutURL(*httphost, *httpport, "ns-csv", "yname:frend:first,monica"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(MultiValPutURL(*httphost, *httpport, "ns-csv", "yname:frend:second,lolita"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:frend"))
	strings2 := strings.Split(body, "\n")
	strings2_expected := []string{"yname:frend:first,monica",
		"yname:frend:second,lolita"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "yname"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(MultiValPutURL(*httphost, *httpport, "ns-csv", "yname:frend:first,monica%0D%0Ayname:frend:second,lolita%0D%0Auname:frend:second,juno%0D%0A"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "yname:frend"))
	strings2 = strings.Split(body, "\n")
	strings2_expected = []string{"yname:frend:first,monica",
		"yname:frend:second,lolita"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "uname:frend"))
	golassert.AssertEqual(body, "uname:frend:second,juno")

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "yname"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "uname"))
	golassert.AssertEqual(body, "Success")

	fmt.Println("No panic for 'csv' key  type.")
}

// for JSON val-type
func testForJSON() {
	_, body := HttpGet(MultiValPutURL(*httphost, *httpport, "ns-json", "{\"power:first\":\"yay\",\"power:second\":\"way\",\"rower:second\":\"kay\"}"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "power"))
	strings2 := strings.Split(body, "\n")
	strings2_expected := []string{"power:first,yay", "power:second,way"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	_, body = HttpGet(GetURL(*httphost, *httpport, "tsds", "rower"))
	golassert.AssertEqual(body, "rower:second,kay")

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "power"))
	golassert.AssertEqual(body, "Success")

	_, body = HttpGet(DelURL(*httphost, *httpport, "tsds", "rower"))
	golassert.AssertEqual(body, "Success")

	fmt.Println("No panic for 'json' key  type.")
}

// for &parentNS=parent:namespace
func testWithParentNS() {
	var url string

	url = MultiValPutURL(*httphost, *httpport, "ns-csv-parent", "yname:frend:first,monica")
	_, body := HttpGet(URLAppendParentNS(url, "animal:people"))
	golassert.AssertEqual(body, "Success")

	url = MultiValPutURL(*httphost, *httpport, "ns-csv-parent", "yname:frend:second,lolita")
	_, body = HttpGet(URLAppendParentNS(url, "animal:people"))
	golassert.AssertEqual(body, "Success")

	url = URLAppendParentNS(GetURL(*httphost, *httpport, "tsds-csv-parent", "yname:frend"), "animal:people")
	_, body = HttpGet(url)
	strings2 := strings.Split(body, "\n")
	strings2_expected := []string{"animal:people:yname:frend:first,monica",
		"animal:people:yname:frend:second,lolita"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	url = URLAppendParentNS(GetURL(*httphost, *httpport, "tsds-csv-parent", "people:yname:frend"), "animal")
	_, body = HttpGet(url)
	strings2 = strings.Split(body, "\n")
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	url = GetURL(*httphost, *httpport, "tsds-csv-parent", "animal:people:yname:frend")
	_, body = HttpGet(url)
	strings2 = strings.Split(body, "\n")
	strings2_expected = []string{"animal:people:yname:frend:first,monica",
		"animal:people:yname:frend:second,lolita"}
	golassert.AssertEqualStringArray(strings2, strings2_expected)

	url = URLAppendParentNS(DelURL(*httphost, *httpport, "tsds-csv-parent", "yname"), "animal:people")
	_, body = HttpGet(url)
	golassert.AssertEqual(body, "Success")

	fmt.Println("No panic for 'Parent NS' key  type.")
}

func main() {
	flag.Parse()

	testDefaultKeyType()
	testNamespaceKeyType()
	testTSDSKeyType()
	testNowKeyType()
	testForCSV()
	testForJSON()
	testWithParentNS()
	fmt.Println("passed not panic")
}
