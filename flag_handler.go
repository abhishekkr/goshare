package goshare

import (
	"flag"
	"fmt"

	"github.com/abhishekkr/gol/golconfig"
)

// flags
var (
	flagConfig = flag.String("config", "", "the path to overriding config file")

	flagDBEngine = flag.String("DBEngine", "leveldb", "the type of KeyVal DB backend to be used")
	flagNSEngine = flag.String("NSEngine", "delimited", "the type of NameSpace DB backend to be used")
	flagTSEngine = flag.String("TSEngine", "namespace", "the type of TimeSeries backend to be used")

	flagDBPath     = flag.String("DBPath", "/tmp/GO.DB", "the path to DB")
	flagServerUri  = flag.String("server-uri", "0.0.0.0", "what Port to Run HTTP Server at")
	flagHTTPPort   = flag.String("http-port", "9999", "what Port to Run HTTP Server at")
	flagRepPorts   = flag.String("rep-ports", "9898,9797", "what PORT to run ZMQ REP at")
	flagCPUProfile = flag.String("cpuprofile", "", "write cpu profile to file")
)

/* assignIfEmpty assigns val to *key only if it's empty */
func assignIfEmpty(mapper golconfig.FlatConfig, key string, val string) {
	if mapper[key] == "" {
		mapper[key] = val
	}
}

/*
ConfigFromFlags configs from values provided to flags.
*/
func ConfigFromFlags() golconfig.FlatConfig {
	flag.Parse()

	var config golconfig.FlatConfig
	config = make(golconfig.FlatConfig)
	if *flagConfig != "" {
		configFile := golconfig.GetConfigurator("json")
		configFile.ConfigFromFile(*flagConfig, &config)
	}

	assignIfEmpty(config, "DBEngine", *flagDBEngine)
	assignIfEmpty(config, "NSEngine", *flagNSEngine)
	assignIfEmpty(config, "TSEngine", *flagTSEngine)
	assignIfEmpty(config, "DBPath", *flagDBPath)
	assignIfEmpty(config, "server-uri", *flagServerUri)
	assignIfEmpty(config, "http-port", *flagHTTPPort)
	assignIfEmpty(config, "rep-ports", *flagRepPorts)
	assignIfEmpty(config, "cpuprofile", *flagCPUProfile)

	fmt.Println("GoShare config:")
	for cfg, val := range config {
		fmt.Printf("[ %v : %v ]\n", cfg, val)
	}
	return config
}
