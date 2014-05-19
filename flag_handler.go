package goshare

import (
	"flag"
	"fmt"

	"github.com/abhishekkr/gol/golconfig"
)

type Config map[string]string

// flags
var (
	flag_config     = flag.String("config", "", "the path to overriding config file")
	flag_dbpath     = flag.String("dbpath", "/tmp/GO.DB", "the path to DB")
	flag_httpuri    = flag.String("http-uri", "0.0.0.0", "what Port to Run HTTP Server at")
	flag_httpport   = flag.String("http-port", "9999", "what Port to Run HTTP Server at")
	flag_req_port   = flag.String("req-port", "9797", "what PORT to run ZMQ REQ at")
	flag_rep_port   = flag.String("rep-port", "9898", "what PORT to run ZMQ REP at")
	flag_cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

/* assign val to *key only if it's empty */
func assignIfEmpty(mapper Config, key string, val string) {
	if mapper[key] == "" {
		mapper[key] = val
	}
}

/* config from flags */
func ConfigFromFlags() Config {
	flag.Parse()

	var config Config
	config = make(Config)
	if *flag_config != "" {
		config_file := golconfig.GetConfig("json")
		config_file.ConfigFromFile(*flag_config, &config)
	}

	assignIfEmpty(config, "dbpath", *flag_dbpath)
	assignIfEmpty(config, "http-uri", *flag_httpuri)
	assignIfEmpty(config, "http-port", *flag_httpport)
	assignIfEmpty(config, "req-port", *flag_req_port)
	assignIfEmpty(config, "rep-port", *flag_rep_port)
	assignIfEmpty(config, "cpuprofile", *flag_cpuprofile)

	fmt.Println("Starting for:", config)
	return config
}
