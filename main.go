package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/notifications"
	//"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/connector/redis"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/utils"
)

var (
	//_ = flag.String("f", "config.json", "config file path")
	_ = flag.String("listen", "0.0.0.0:8080", "docker-rec http server listen address")
	_ = flag.String("redis", "0.0.0.0:6379", "redis connection string")
	_ = flag.String("crt", "certs/domain.crt", "key pair crt file")
	_ = flag.String("key", "certs/domain.key", "key pair key file")
)

var (
	RC  *rrredis.RedisClient
	err error
)

func main() {

	//if !rrutils.FlagIsSet("f") || !rrutils.FlagIsSet("redis") {
	//	rrutils.FlagHelp()
	//	return
	//}

	//crt, _ := rrutils.FlagGetString("crt")
	//key, _ := rrutils.FlagGetString("key")

	//path, _ := rrutils.FlagGetString("f")
	//jc, err := rrconfig.LoadJsonConfigFromFile(path)
	//if err != nil {
	//	logs.Error(err)
	//	return
	//}

	// connect redis
	connStr, _ := rrutils.FlagGetString("redis")
	err, RC = rrredis.GetRedisClient(connStr)
	if err != nil {
		logs.Error(err)
		return
	}

	// listen
	listen, _ := rrutils.FlagGetString("listen")
	http.HandleFunc("/events", eventHandler)
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Info("Exiting")
}

func eventHandler(w http.ResponseWriter, req *http.Request) {

	// The docker registry sends events to HTTP endpoints and queues them up if
	// the endpoint refuses to accept those events. We are only interested in
	// manifest updates, therefore we ignore all others by answering with an HTTP
	// 200 OK. This should prevent the docker registry from getting too full.

	// A request needs to be made via POST
	if req.Method != "POST" {
		logs.Error("not post")
		http.Error(w, fmt.Sprintf("Ignoring request. Required method is \"POST\" but got \"%s\".\n", req.Method), http.StatusOK)
		return
	}

	// A request must have a body.
	if req.Body == nil {
		logs.Error("body is nil")
		http.Error(w, "Ignoring request. Required non-empty request body.\n", http.StatusOK)
		return
	}

	// Test for correct mimetype and reject all others
	// Even the documentation on docker notfications says that we shouldn't be to
	// picky about the mimetype. But we are and let the caller know this.
	contentType := req.Header.Get("Content-Type")
	if contentType != notifications.EventsMediaType {
		logs.Error("Content-Type invalid")
		http.Error(w, fmt.Sprintf("Ignoring request. Required mimetype is \"%s\" but got \"%s\"\n", notifications.EventsMediaType, contentType), http.StatusOK)
		return
	}

	// Try to decode HTTP body as Docker notification envelope
	decoder := json.NewDecoder(req.Body)
	var envelope notifications.Envelope
	err := decoder.Decode(&envelope)
	if err != nil {
		logs.Error("Failed to decode envelope")
		http.Error(w, fmt.Sprintf("Failed to decode envelope: %s\n", err), http.StatusBadRequest)
		return
	}

	for _, event := range envelope.Events {
		logs.Debug(event.Target.MediaType, schema1.MediaTypeManifest)
		if event.Target.MediaType != schema2.MediaTypeManifest &&
			event.Target.MediaType != schema1.MediaTypeManifest {
			http.Error(w, fmt.Sprintf("Wrong event.Target.MediaType: \"%s\". Expected: \"%s\"", event.Target.MediaType, schema2.MediaTypeManifest), http.StatusBadRequest)
			return
		}

		// Handle all three cases: push, pull, and delete
		if event.Action == notifications.EventActionPull || event.Action == notifications.EventActionPush {
			logs.Info(event.Action, "event", event)

		} else if event.Action == notifications.EventActionDelete {

			logs.Info(event.Action, "event", event)

		} else {

			http.Error(w, fmt.Sprintf("Invalid event action: %s\n", event.Action), http.StatusBadRequest)
			return

		}

	}

	http.Error(w, fmt.Sprintf("Done\n"), http.StatusOK)
}
