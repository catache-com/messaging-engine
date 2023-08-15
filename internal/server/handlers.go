package server

import (
	"github.com/sirupsen/logrus"
	"messaging-engine/internal/models"
	"messaging-engine/internal/util"
	"net/http"
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
