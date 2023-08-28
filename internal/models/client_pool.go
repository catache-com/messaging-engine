package models

import (
	"sync"
)

type ClientPool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	rwMutex    sync.RWMutex
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
		}
	}
}

func (ClientPool *ClientPool) GetTheClient(clientId string) *Client {
	ClientPool.rwMutex.Lock()
	defer ClientPool.rwMutex.Unlock()

	if client, ok := ClientPool.Clients[clientId]; ok {
		return client
	}

	return nil
}

func (ClientPool *ClientPool) SendMsgToClient(message Message) {
	ClientPool.rwMutex.Lock()
	defer ClientPool.rwMutex.Unlock()

	foundClient := ClientPool.Clients[message.SendTo]
	if foundClient != nil {
		foundClient.Write(message)
	}
}
