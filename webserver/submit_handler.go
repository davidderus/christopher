package webserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"

	"github.com/davidderus/christopher/dispatcher"
)

type submitRequest struct {
	Urls string
}

type submitResponse struct {
	Count  int      `json:"count"`
	Errors []string `json:"errors"`
}

var uriMatcher = regexp.MustCompile(`(https?:\/\/[\da-z\.-]+\.[a-z\.]{2,6}[\/\w \.-]*\/?)`)

func (ws *WebServer) dispatchURI(uri string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	scenario := ws.loadScenario()
	scenario.Play(event)
}

// SubmitHandler handles submitted links
func (ws *WebServer) SubmitHandler(w http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var submittedRequest submitRequest

	unmarshallError := json.Unmarshal(body, &submittedRequest)
	if unmarshallError != nil {
		http.Error(w, unmarshallError.Error(), http.StatusInternalServerError)
		return
	}

	uris := uriMatcher.FindAllString(submittedRequest.Urls, -1)
	urisCount := len(uris)

	var waitGroup sync.WaitGroup
	waitGroup.Add(urisCount)

	for _, uri := range uris {
		go ws.dispatchURI(string(uri), &waitGroup)
	}

	waitGroup.Wait()

	marshaledJSON, jsonError := json.Marshal(submitResponse{Count: urisCount, Errors: nil})
	if jsonError != nil {
		http.Error(w, jsonError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledJSON)
}
