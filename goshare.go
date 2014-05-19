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
	abkleveldb "github.com/abhishekkr/levigoNS/leveldb"

	"github.com/jmhodges/levigo"
)

var (
	db *levigo.DB
)

/* just a banner print */
func banner() {
	fmt.Println("**************************************************")
	fmt.Println("  ___  ____      ___        __   _   __")
	fmt.Println("  |    |  |      |    |  | /  \\ | ) |")
	fmt.Println("  | =| |  |  =~  |==| |==| |==| |=  |=")
	fmt.Println("  |__| |__|      ___| |  | |  | | \\ |__")
	fmt.Println("")
	fmt.Println("**************************************************")
}

/* checking if you still wanna keep the goshare up */
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
putting together base engine for GoShare
dbpath, httpuri, httpport, rep_port, req_port *string
*/
func GoShareEngine(config Config) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// remember it will be same DB instance shared across goshare package
	db = abkleveldb.CreateDB(config["dbpath"])
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

	_httpport, err_httpport := strconv.Atoi(config["http-port"])
	fmt.Println(_httpport)
	_req_port, err_req_port := strconv.Atoi(config["req-port"])
	_rep_port, err_rep_port := strconv.Atoi(config["rep-port"])
	if err_httpport == nil && err_rep_port == nil && err_req_port == nil {
		go GoShareHTTP(config["http-uri"], _httpport)
		go GoShareZMQ(_req_port, _rep_port)
	} else {
		golerror.Boohoo("Port parameters to bind, error-ed while conversion to number.", true)
	}
}

/* GoShare DB */
func GoShare() {
	banner()
	GoShareEngine(ConfigFromFlags())
	DoYouWannaContinue()
}
