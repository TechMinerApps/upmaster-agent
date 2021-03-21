package main

import (
	"os"
	"os/signal"

	"github.com/TechMinerApps/upmaster-agent/app"
)

func main() {
	agent := app.NewAgent()
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan)

	// Graceful Shutdown
	go func() {
		<-sigchan
		agent.Stop()
	}()
	agent.Start()
}
