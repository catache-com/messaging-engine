package models

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type ClientPool struct {
	Register   chan *Client
	Unregister chan *Client
	Messages   chan Message
	Clients    map[string]*Client
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (ClientPool *ClientPool) Start(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when the function returns

	for {
		select {
		case client := <-ClientPool.Register:
			ClientPool.Clients[client.ID] = client
			break
		case client := <-ClientPool.Unregister:
			delete(ClientPool.Clients, client.ID)
			break
		case message := <-ClientPool.Messages:
			ClientPool.SendMsgToClient(message)
		}
	}
}

func (ClientPool *ClientPool) GetTheClient(clientId string) *Client {
	if client, ok := ClientPool.Clients[clientId]; ok {
		return client
	}

	return nil
}

func (ClientPool *ClientPool) SendMsgToClient(message Message) {
	foundClient := ClientPool.Clients[message.SendTo]
	if foundClient != nil {
		foundClient.Write(message)
	}
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
