package goshare

import (
	"flag"
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

type Config map[string]*string

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
func do_you_wanna_continue() {
	var input string
	for {
		fmt.Println("Do you wanna exit. (yes|no):\n\n")

		fmt.Scanf("%s", &input)

		if input == "yes" || input == "y" {
			break
		}
	}
}

/* config from flags */
func ConfigFromFlags() Config {
	var config Config
	config = make(Config)
	config["dbpath"] = flag.String("dbpath", "/tmp/GO.DB", "the path to DB")
	config["httpuri"] = flag.String("uri", "0.0.0.0", "what Port to Run HTTP Server at")
	config["httpport"] = flag.String("port", "9999", "what Port to Run HTTP Server at")
	config["req_port"] = flag.String("req-port", "9797", "what PORT to run ZMQ REQ at")
	config["rep_port"] = flag.String("rep-port", "9898", "what PORT to run ZMQ REP at")
	config["cpuprofile"] = flag.String("cpuprofile", "", "write cpu profile to file")

	//flag.Parse()
	return config
}

/*
putting together base engine for GoShare
dbpath, httpuri, httpport, rep_port, req_port *string
*/
func GoShareEngine(config Config) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// remember it will be same DB instance shared across goshare package
	db = abkleveldb.CreateDB(*config["dbpath"])
	if *config["cpuprofile"] != "" {
		f, err := os.Create(*config["cpuprofile"])
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		go func() {
			time.Sleep(100 * time.Second)
			pprof.StopCPUProfile()
		}()
	}

	_httpport, err_httpport := strconv.Atoi(*config["httpport"])
	_req_port, err_req_port := strconv.Atoi(*config["req_port"])
	_rep_port, err_rep_port := strconv.Atoi(*config["rep_port"])
	if err_httpport == nil && err_rep_port == nil && err_req_port == nil {
		go GoShareHTTP(*config["httpuri"], _httpport)
		go GoShareZMQ(_req_port, _rep_port)
	} else {
		golerror.Boohoo("Port parameters to bind, error-ed while conversion to number.", true)
	}
}

/* GoShare DB */
func GoShare() {
	banner()
	GoShareEngine(ConfigFromFlags())
	do_you_wanna_continue()
}
