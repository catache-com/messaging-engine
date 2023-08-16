package server

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"messaging-engine/internal/db/mongo"
	"messaging-engine/internal/models"
	"messaging-engine/internal/util"
	"net/http"
	"time"
)

var ClientPool *models.ClientPool

func AcceptConnection(w http.ResponseWriter, r *http.Request) {
	clientId, ok := r.URL.Query()["client_id"]

	if !ok || len(clientId[0]) < 1 {
		util.WriteJSONResponse(
			w,
			http.StatusBadRequest,
			[]byte("client_id query param is missing"),
		)
		return
	}

	wsConnection, err := models.UpgradeHTTPToWS(w, r)
	if err != nil {
		logrus.Errorf("error accepting connection, %v", err)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error connecting"))
		return
	}

	// create a new client
	client := models.NewClient(clientId[0], wsConnection, ClientPool, HandleMessage)

	// register client into pool
	ClientPool.Register <- client

	// make client listening for new messages
	go client.Read()
}

func HandleSendMessageToClient(w http.ResponseWriter, r *http.Request) {
	type got struct {
		ClientId string         `json:"client_id"`
		Message  models.Message `json:"message"`
	}

	var g got
	err := util.DecodeJSONBody(w, r, &g)
	if err != nil {
		logrus.Errorf("error decoding JSON body when HandleSendMessageToClient, %v", err)
		util.WriteJSONResponse(w, http.StatusBadRequest, []byte("error"))
		return
	}

	// check the client we are sending to is existing
	clientId := g.ClientId
	foundClient := ClientPool.GetTheClient(clientId)
	if foundClient != nil {
		// server handles the message
		err := HandleMessage(g.Message)
		if err != nil {
			util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
			return
		} else {
			// ready to send to another end client
			ClientPool.SendMsgToClients([]string{clientId}, g.Message)
			util.WriteJSONResponse(w, http.StatusOK, []byte("OK"))
			return
		}
	}

	util.WriteJSONResponse(w, http.StatusAccepted, []byte("Client not online."))
}

func HandleGetChannelMessages(w http.ResponseWriter, r *http.Request) {
	type got struct {
		ChannelId      string    `json:"channel_id"`
		DatetimeAnchor time.Time `json:"datetime_anchor"`
		Pagination     int64     `json:"pagination"`
	}

	var g got
	err := util.DecodeJSONBody(w, r, &g)
	if err != nil {
		logrus.Errorf("error decoding JSON body when HandleGetChannelMessages, %v", err)
		util.WriteJSONResponse(w, http.StatusBadRequest, []byte("error"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foundMessages, err := mongo.FindChannelMessagesByChannelId(
		ctx,
		g.ChannelId,
		g.DatetimeAnchor,
		g.Pagination,
	)
	if err != nil {
		logrus.Errorf(
			"error db.FindChannelMessagesByChannelId for ChannelId: %s, DatetimeAnchor: %s, Pagination: %d : %v",
			g.ChannelId,
			g.DatetimeAnchor,
			g.Pagination,
			err,
		)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
	}

	byteFoundMessages, err := json.Marshal(foundMessages)
	if err != nil {
		logrus.Errorf("error json.Marshal foundMessages, %v", err)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, byteFoundMessages)
}

func HandleGetThreadMessages(w http.ResponseWriter, r *http.Request) {
	type got struct {
		ThreadId       string    `json:"thread_id"`
		DatetimeAnchor time.Time `json:"datetime_anchor"`
		Pagination     int64     `json:"pagination"`
	}

	var g got
	err := util.DecodeJSONBody(w, r, &g)
	if err != nil {
		logrus.Errorf("error decoding JSON body when HandleGetThreadMessages, %v", err)
		util.WriteJSONResponse(w, http.StatusBadRequest, []byte("error"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foundMessages, err := mongo.FindThreadMessagesByThreadId(
		ctx,
		g.ThreadId,
		g.DatetimeAnchor,
		g.Pagination,
	)
	if err != nil {
		logrus.Errorf(
			"error db.FindThreadMessagesByThreadId for ThreadId: %s, DatetimeAnchor: %s, Pagination: %d : %v",
			g.ThreadId,
			g.DatetimeAnchor,
			g.Pagination,
			err,
		)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
	}

	byteFoundMessages, err := json.Marshal(foundMessages)
	if err != nil {
		logrus.Errorf("error json.Marshal foundMessages, %v", err)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, byteFoundMessages)
}

func HandleMakeNewChannel(w http.ResponseWriter, r *http.Request) {
	type got struct {
		ChannelClients []string `json:"channel_clients"`
	}

	var g got
	err := util.DecodeJSONBody(w, r, &g)
	if err != nil {
		logrus.Errorf("error decoding JSON body when HandleMakeNewChannel, %v", err)
		util.WriteJSONResponse(w, http.StatusBadRequest, []byte("error"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ChannelId := uuid.New()

	newChannel := models.Channel{
		Id:      ChannelId.String(),
		Clients: g.ChannelClients,
	}

	err = mongo.NewChannel(
		ctx,
		newChannel,
	)
	if err != nil {
		logrus.Errorf(
			"error db.NewChannel: %v", err,
		)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
	}
}

func HandleMakeNewThread(w http.ResponseWriter, r *http.Request) {
	type got struct {
		ChannelId     string `json:"channel_id"`
		RootMessageId string `json:"root_message_id"`
	}

	var g got
	err := util.DecodeJSONBody(w, r, &g)
	if err != nil {
		logrus.Errorf("error decoding JSON body when HandleMakeNewThread, %v", err)
		util.WriteJSONResponse(w, http.StatusBadRequest, []byte("error"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	threadId := uuid.New()

	newThread := models.Thread{
		Id:            threadId.String(),
		ChannelId:     g.ChannelId,
		RootMessageId: g.RootMessageId,
	}

	err = mongo.NewThread(
		ctx,
		newThread,
	)
	if err != nil {
		logrus.Errorf(
			"error db.NewThread: %v", err,
		)
		util.WriteJSONResponse(w, http.StatusInternalServerError, []byte("error"))
	}
}
