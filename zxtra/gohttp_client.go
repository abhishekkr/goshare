package main

import (
  "fmt"
  "flag"
  "io/ioutil"
  "net/http"
)

var (
  httphost   = flag.String("host", "127.0.0.1", "what Host to run at")
  httpport   = flag.Int("port", 9999, "what Socket PORT to connect")
)

func GetURL(host string, port int, key string) string{
  return fmt.Sprintf("http://%s:%d/get?key='%s'", host, port, key)
}

func PutURL(host string, port int, key string, val string) string{
  return fmt.Sprintf("http://%s:%d/put?key='%s'&val='%s'",
                     host, port, key, val)
}

func HttpGet(url string) string{
  resp, err := http.Get(url)
  if err != nil {
    return "Error: " + url + " failed for HTTP GET"
  }
  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  return string(body)
}

func main(){
  flag.Parse()

  fmt.Println("Creating a Key: 'myname'")
  resp := HttpGet(PutURL(*httphost, *httpport, "myname", "anon"))
  fmt.Println("Response:", resp)
  fmt.Println("Reading a Key: 'myname'")
  resp = HttpGet(GetURL(*httphost, *httpport, "myname"))
  fmt.Println("Response:", resp)

  fmt.Println("Updating a Key: 'myname'")
  resp = HttpGet(PutURL(*httphost, *httpport, "myname", "anonymous"))
  fmt.Println("Response:", resp)
  fmt.Println("Reading a Key: 'myname'")
  resp = HttpGet(GetURL(*httphost, *httpport, "myname"))
  fmt.Println("Response:", resp)
}
