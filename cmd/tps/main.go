package main

import (
	"flag"
	"time"
	"github.com/jncornett/tps"
	"os/exec"
	"log"
	"os"
)

func main() {
	var (
		transactionsPerSecond = flag.Float64("tps", 1, "transactions per second")
		duration = flag.Duration("duration", time.Minute, "duration")
	)
	flag.Parse()

	t := tps.New(*transactionsPerSecond)
	cancelTps := make(chan struct{})
	go t.Run(cancelTps)

	cancelHandler := make(chan struct{})
	doneHandler := make(chan struct{})
	go func(cancel <- chan struct{}) {
		for {
			select {
			case <-cancel:
				doneHandler <- struct{}{}
				return
			case <-t.Events():
				err := runCmd(flag.Args()...)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}(cancelHandler)

	<-time.After(*duration)
	cancelTps <- struct{}{}
	cancelHandler <- struct{}{}
	<-t.Done()
	<-doneHandler
}

func runCmd(args...string) error {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
