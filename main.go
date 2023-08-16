package main

import (
	"github.com/sirupsen/logrus"
	"messaging-engine/internal/models"
	"messaging-engine/internal/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	handleSigterm(func() {
		logrus.Info("Captured Ctrl+C")
	})

	var wg sync.WaitGroup
	wg.Add(2)

	// initialise ClientPool
	ClientPool := models.NewClientPool()
	// start websocket clients ClientPool to receive messages
	go ClientPool.Start(&wg)

	server.ClientPool = ClientPool
	// start messaging-engine as a service
	go server.StartMessagingEngine(&wg)

	wg.Wait() // Wait for all the goroutines to finish
}

func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
