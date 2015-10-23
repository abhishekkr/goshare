package goshare

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	golconfig "github.com/abhishekkr/gol/golconfig"
	golerror "github.com/abhishekkr/gol/golerror"
	golkeyvalTSDS "github.com/abhishekkr/gol/golkeyvalTSDS"
	gollist "github.com/abhishekkr/gol/gollist"
)

var (
	tsds golkeyvalTSDS.TSDSDBEngine
)

/* banner just brand print */
func banner() {
	fmt.Println("**************************************************")
	fmt.Println("  ___  ____      ___        __   _   __")
	fmt.Println("  |    |  |      |    |  | /  \\ | ) |")
	fmt.Println("  | =| |  |  =~  |==| |==| |==| |=  |=")
	fmt.Println("  |__| |__|      ___| |  | |  | | \\ |__")
	fmt.Println("")
	fmt.Println("**************************************************")
}

/*
DoYouWannaContinue checking if you still wanna keep the goshare up.
*/
func DoYouWannaContinue() {
	var input string
	for {
		fmt.Println("Do you wanna exit. (yes|no):\n\n")

		fmt.Scanf("%s", &input)

		if input == "yes" || input == "y" {
			break
		}
	}
}

/*
goshareDB returns golkeyval DBEngine for it.
*/
func goshareDB(config golconfig.FlatConfig) {
	if config["TSEngine"] == "namespace" {
		tsds = golkeyvalTSDS.GetNamespaceEngine(config)
	} else {
		panic("Unhandled TimeSeries Engine required.")
	}
}

/*
GoShareEngine putting together base engine for GoShare as per config.
dbpath, server_uri, httpport, rep_port, *string
*/
func GoShareEngine(config golconfig.FlatConfig) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// remember it will be same DB instance shared across goshare package
	goshareDB(config)

	if config["cpuprofile"] != "" {
		f, err := os.Create(config["cpuprofile"])
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		go func() {
			time.Sleep(100 * time.Second)
			pprof.StopCPUProfile()
		}()
	}

	_httpPort, errHTTPPort := strconv.Atoi(config["http-port"])
	_replyPorts, errReplyPorts := gollist.CSVToNumbers(config["rep-ports"])
	if errHTTPPort == nil && errReplyPorts == nil {
		go GoShareHTTP(config["server-uri"], _httpPort)
		go GoShareZMQ(config["server-uri"], _replyPorts)
	} else {
		golerror.Boohoo("Port parameters to bind, error-ed while conversion to number.", true)
	}
}

/*
GoShare is daddy-o of goshare instance.
*/
func GoShare() {
	banner()
	runtime.GOMAXPROCS(runtime.NumCPU())
	GoShareEngine(ConfigFromFlags())
	DoYouWannaContinue()
}
