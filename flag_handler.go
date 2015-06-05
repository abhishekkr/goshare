package goshare

import (
	"flag"
	"fmt"

	"github.com/abhishekkr/gol/golconfig"
)

/*
Config is a hashmap used here to carry around param=paramValue for GoShare.
*/
type Config map[string]string

// flags
var (
	flagConfig     = flag.String("config", "", "the path to overriding config file")
	flagDBEngine   = flag.String("DBEngine", "leveldb", "the type of KeyVal DB backend to be used")
	flagDBPath     = flag.String("DBPath", "/tmp/GO.DB", "the path to DB")
	flagServerUri  = flag.String("server-uri", "0.0.0.0", "what Port to Run HTTP Server at")
	flagHTTPPort   = flag.String("http-port", "9999", "what Port to Run HTTP Server at")
	flagRepPorts   = flag.String("rep-ports", "9898,9797", "what PORT to run ZMQ REP at")
	flagCPUProfile = flag.String("cpuprofile", "", "write cpu profile to file")
)

/* assignIfEmpty assigns val to *key only if it's empty */
func assignIfEmpty(mapper Config, key string, val string) {
	if mapper[key] == "" {
		mapper[key] = val
	}
}

/*
ConfigFromFlags configs from values provided to flags.
*/
func ConfigFromFlags() Config {
	flag.Parse()

	var config Config
	config = make(Config)
	if *flagConfig != "" {
		configFile := golconfig.GetConfig("json")
		configFile.ConfigFromFile(*flagConfig, &config)
	}

	assignIfEmpty(config, "DBEngine", *flagDBEngine)
	assignIfEmpty(config, "DBPath", *flagDBPath)
	assignIfEmpty(config, "server-uri", *flagServerUri)
	assignIfEmpty(config, "http-port", *flagHTTPPort)
	assignIfEmpty(config, "rep-ports", *flagRepPorts)
	assignIfEmpty(config, "cpuprofile", *flagCPUProfile)

	fmt.Println("Starting for:", config)
	return config
}
