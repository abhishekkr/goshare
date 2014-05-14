package main

import (
	"flag"
	"fmt"

	"github.com/abhishekkr/goshare"

	"github.com/ChaosCloud/godaemon"
)

var (
	axn = flag.String("axn", "status", "status|start|stop")
)

func main() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	flag.Parse()
	switch *axn {
	case "stop":
		fmt.Println("Stopping...")
		if godaemon.KillPID("") {
			fmt.Println("Status: Stopped.")
		} else {
			fmt.Printf("Failed to stop. Status: ")
			godaemon.StatusPID("")
		}
	case "status":
		fmt.Printf("Status: ")
		godaemon.StatusPID("")
	default:
		if godaemon.PersistPID("") {
			goshare.GoShare()
		}
	}
}
