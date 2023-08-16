package server

import (
	"net/http"
)

type Route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var MessagingEngineProtectedRoutes = Routes{}

var MessagingEngineOpenRoutes = Routes{
	// ---------- probing ----------
	Route{
		Name:    "HealthCheck",
		Method:  "GET",
		Pattern: "/healthz",
		HandlerFunc: func(writer http.ResponseWriter, r *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte("OK"))
			if err != nil {
				return
			}
		},
	},

	// ---------- messaging engine ----------

	Route{
		Name:        "connect",
		Method:      "GET",
		Pattern:     "/connect",
		HandlerFunc: AcceptConnection,
	},

	Route{
		Name:        "send message to a client",
		Method:      "POST",
		Pattern:     "/client/send",
		HandlerFunc: HandleSendMessageToClient,
	},

	Route{
		Name:        "send message to a client",
		Method:      "POST",
		Pattern:     "/clients/send",
		HandlerFunc: HandleSendMessageToClients,
	},

	Route{
		Name:        "find messages in a channel",
		Method:      "GET",
		Pattern:     "/messages/channel",
		HandlerFunc: HandleGetChannelMessages,
	},

	Route{
		Name:        "find messages in a channel",
		Method:      "GET",
		Pattern:     "/messages/thread",
		HandlerFunc: HandleGetThreadMessages,
	},

	Route{
		Name:        "initialise a new channel",
		Method:      "POST",
		Pattern:     "/channel/new",
		HandlerFunc: HandleMakeNewChannel,
	},

	Route{
		Name:        "initialise a new thread",
		Method:      "POST",
		Pattern:     "/thread/new",
		HandlerFunc: HandleMakeNewThread,
	},
}
