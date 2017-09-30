package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	// Instead of writing our own json2 support, we're using this one
	"github.com/gorilla/rpc/v2/json2"
)

// Aria2 is a downloader interface for aria2
// see https://aria2.github.io/manual/en/html/aria2c.html#methods for more infos
type Aria2 struct {
	client *http.Client

	// token is the Aria2 RPC token if any
	token string

	// rpcURL is the full RPC URL to the aria2 JSON-RPC endpoint
	rpcURL string

	// timeOut is the maximum time for a request to take in seconds
	timeOut time.Duration

	// Allow to set a custom HTTP transport (for test purposes)
	CustomTransport http.RoundTripper
}

const (
	// ariaDownloaderDefaulttimeOut is the default timeout for http requests in seconds
	ariaDownloaderDefaulttimeOut = 10
)

// Auth initializes the Aria2
func (ad *Aria2) Auth(infos map[string]interface{}) error {
	var rpcURLOkay, tokenOkay bool

	ad.rpcURL, rpcURLOkay = infos["rpcURL"].(string)
	if !rpcURLOkay || ad.rpcURL == "" {
		return errors.New("Invalid RPC url")
	}

	ad.token, tokenOkay = infos["token"].(string)
	if !tokenOkay {
		return errors.New("Invalid token")
	}

	timeOut, timeOutOkay := infos["timeout"].(int)
	if !timeOutOkay || timeOut == 0 {
		timeOut = ariaDownloaderDefaulttimeOut
	}
	ad.timeOut = time.Duration(time.Duration(timeOut) * time.Second)

	ad.client = &http.Client{
		Timeout:   ad.timeOut,
		Transport: ad.CustomTransport,
	}

	return nil
}

// call sends a request to Aria2
func (ad *Aria2) call(method string, params interface{}, result interface{}) error {
	message, encodeError := json2.EncodeClientRequest(method, &params)
	if encodeError != nil {
		return encodeError
	}

	response, responseError := ad.client.Post(ad.rpcURL, "application/json", bytes.NewBuffer(message))
	if responseError != nil {
		return responseError
	}

	defer response.Body.Close()

	decodeError := json2.DecodeClientResponse(response.Body, &result)
	if decodeError != nil {
		return decodeError
	}

	return nil
}

// Download starts the download of a given uri
func (ad *Aria2) Download(uri string, options map[string]interface{}) (string, error) {
	var gid string
	var paramsArray []interface{}

	uris := make([]string, 1)
	uris[0] = uri

	if options != nil {
		paramsArray = ad.appendParams(uris, options)
	} else {
		paramsArray = ad.appendParams(uris)
	}

	callError := ad.call("aria2.addUri", paramsArray, &gid)
	if callError != nil {
		return "", callError
	}

	return gid, nil
}

// DownloadStatus returns some status infos about the download
func (ad *Aria2) DownloadStatus(downloadID string) (map[string]interface{}, error) {
	var status map[string]interface{}

	callError := ad.call("aria2.tellStatus", ad.appendParams(downloadID), &status)
	if callError != nil {
		return nil, callError
	}

	return status, nil
}

// appendParams append all given params and wrap them with a token if any
func (ad *Aria2) appendParams(params ...interface{}) []interface{} {
	paramsArray := make([]interface{}, 0)

	if ad.token != "" {
		paramsArray = append(paramsArray, fmt.Sprintf("token:%s", ad.token))
	}

	for _, param := range params {
		paramsArray = append(paramsArray, param)
	}

	return paramsArray
}
