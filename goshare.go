package goshare

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	golerror "github.com/abhishekkr/gol/golerror"
	golkeyval "github.com/abhishekkr/gol/golkeyval"
	gollist "github.com/abhishekkr/gol/gollist"
	gollog "github.com/abhishekkr/gol/gollog"
)

var (
	db      golkeyval.DBEngine
	logInfo Log
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
func goshareDB(config Config) golkeyval.DBEngine {
	db := golkeyval.GetDBEngine(config["DBEngine"])
	db.Configure(config)
	db.CreateDB()
	return db
}

/*
GoShareEngine putting together base engine for GoShare as per config.
dbpath, server_uri, httpport, rep_port, *string
*/
func GoShareEngine(config Config) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// remember it will be same DB instance shared across goshare package
	db = goshareDB(config)
	logInfo = gollog.Log{Level: "info", Thread: make(chan string)}
	go logInfo.LogIt()

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
