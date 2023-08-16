package server

import (
	"github.com/sirupsen/logrus"
	"messaging-engine/internal/config"
	"net/http"
	"strconv"
	"sync"
)

var authMiddleware AuthMiddleware
var messagingEngineConfig config.MessagingEngineConfig

func StartMessagingEngine(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when the goroutine completes

	r := NewRouter(authMiddleware)
	http.Handle("/", r)

	logrus.Infof("Starting messaging engine at %v", messagingEngineConfig.Port)
	logrus.Infof("Engine Id: %v", config.EngineId)

	err := http.ListenAndServe(":"+strconv.Itoa(messagingEngineConfig.Port), nil)

	if err != nil {
		logrus.Errorf("error starting engine service: %v", err)
	}
}
