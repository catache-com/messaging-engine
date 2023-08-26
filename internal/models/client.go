package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
)

type Client struct {
	ID                string
	Conn              *websocket.Conn
	ClientPool        *ClientPool
	HandleMessageFunc func(message Message) error
	readMu            sync.Mutex
	writeMu           sync.Mutex
}

func (c *Client) safeRead() (int, []byte, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()

	return c.Conn.ReadMessage()
}

func (c *Client) SafeWriteJson(message Message) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.Conn.WriteJSON(message)
}

func (c *Client) Read() {
	defer func() {
		c.ClientPool.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			logrus.Infof("failed to close client connection for %s: %v", c.ID, err)
			return
		}
	}()

	for {
		_, m, err := c.safeRead()
		if err != nil {
			logrus.Debugf("client connection error reading message for %s: %v", c.ID, err)
			return
		}

		var message Message
		err = json.Unmarshal(m, &message)
		if err != nil {
			logrus.Errorf("error unmarshal client received message: %v", err)
		} else {
			// client handles the message
			err := c.HandleMessageFunc(message)

			// we only proceed sending to other clients once it's processed
			if err == nil {
				// ready to send to another end client
				c.ClientPool.Messages <- message
			}
		}
	}
}

func (c *Client) Write(message Message) {
	if err := c.SafeWriteJson(message); err != nil {
		logrus.Errorf("failed to write to client %s: %v", c.ID, err)
	}
}

func (c *Client) Leave() {
	c.ClientPool.Unregister <- c
}

func NewClient(
	clientId string,
	conn *websocket.Conn,
	ClientPool *ClientPool,
	HandleMessageFunc func(message Message) error,
) *Client {
	logrus.Infof("creating client %s", clientId)
	return &Client{
		ID:                clientId,
		Conn:              conn,
		ClientPool:        ClientPool,
		HandleMessageFunc: HandleMessageFunc,
	}
}
