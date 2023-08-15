package util

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Connection", "close")
	w.WriteHeader(status)

	_, err := w.Write(data)
	if err != nil {
		logrus.Errorf("error writing json response %v", err)
		return
	}
}
