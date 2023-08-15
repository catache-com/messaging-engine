package models

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type ClientPool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (ClientPool *ClientPool) Start() error {
	for {
		select {
		case client := <-ClientPool.Register:
			ClientPool.Clients[client.ID] = client
			break
		case client := <-ClientPool.Unregister:
			delete(ClientPool.Clients, client.ID)
		}
	}
}

func (ClientPool *ClientPool) GetTheClient(clientId string) *Client {
	if client, ok := ClientPool.Clients[clientId]; ok {
		return client
	}

	return nil
}

func (ClientPool *ClientPool) SendMsgToClients(clientIds []string, message Message) {
	var wg sync.WaitGroup

	for _, id := range clientIds {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			foundClient := ClientPool.Clients[id]
			if foundClient != nil {
				foundClient.Write(message)
			}
		}(id)
	}

	wg.Wait()
}

func (ClientPool *ClientPool) ClientExitFromPool(clientId string) {
	foundClient := ClientPool.GetTheClient(clientId)
	if foundClient != nil {
		ClientPool.Unregister <- foundClient
		err := foundClient.Conn.Close()
		if err != nil {
			logrus.Infof("failed to close client connection for %s: %v", foundClient.ID, err)
		}
	}
}
