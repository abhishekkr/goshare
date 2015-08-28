package goshare

import "fmt"

type Log struct {
	Level  string
	Thread chan string
}

func (l Log) LogIt() {
	for {
		msg := <-(l.Thread)
		fmt.Printf("[%s] %s", l.Level, msg)
	}
}
