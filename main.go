package main

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/notifications"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/utils"
)

var (
	_ = flag.String("f", "config.json", "config file path")
)

func main() {

	if !rrutils.FlagIsSet("f") {
		rruitls.FlagHelp()
		return
	}

	path := rrutils.FlagGetString("f")
	jc, err := rrconfig.LoadJsonConfigFromFile(path)
	if err != nil {
		logs.Error(err)
		return
	}

	http.HandleFunc("/events", eventHandler)
	err = http.ListenAndServeTLS(httpConnectionString, ctx.Config.Server.Ssl.Cert, ctx.Config.Server.Ssl.CertKey, nil)
	if err != nil {
		glog.Exit(err)
	}

	glog.Info("Exiting.")
}

func eventHandler(w http.ResponseWriter, req *http.Request) {

	// The docker registry sends events to HTTP endpoints and queues them up if
	// the endpoint refuses to accept those events. We are only interested in
	// manifest updates, therefore we ignore all others by answering with an HTTP
	// 200 OK. This should prevent the docker registry from getting too full.

	// A request needs to be made via POST
	if req.Method != "POST" {
		http.Error(w, fmt.Sprintf("Ignoring request. Required method is \"POST\" but got \"%s\".\n", req.Method), http.StatusOK)
		return
	}

	// A request must have a body.
	if req.Body == nil {
		http.Error(w, "Ignoring request. Required non-empty request body.\n", http.StatusOK)
		return
	}

	// Test for correct mimetype and reject all others
	// Even the documentation on docker notfications says that we shouldn't be to
	// picky about the mimetype. But we are and let the caller know this.
	contentType := req.Header.Get("Content-Type")
	if contentType != notifications.EventsMediaType {
		http.Error(w, fmt.Sprintf("Ignoring request. Required mimetype is \"%s\" but got \"%s\"\n", notifications.EventsMediaType, contentType), http.StatusOK)
		return
	}

	// Try to decode HTTP body as Docker notification envelope
	decoder := json.NewDecoder(req.Body)
	var envelope notifications.Envelope
	err := decoder.Decode(&envelope)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode envelope: %s\n", err), http.StatusBadRequest)
		return
	}

	for index, event := range envelope.Events {

		// Handle all three cases: push, pull, and delete
		if event.Action == notifications.EventActionPull || event.Action == notifications.EventActionPush {
			logs.Info(event.Action, "event")

		} else if event.Action == notifications.EventActionDelete {

			logs.Info(event.Action, "event")

		} else {

			http.Error(w, fmt.Sprintf("Invalid event action: %s\n", event.Action), http.StatusBadRequest)
			return

		}

	}

	http.Error(w, fmt.Sprintf("Done\n"), http.StatusOK)
}
