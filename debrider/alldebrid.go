package debrider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

// AllDebrid is an interface for alldebrid.com
type AllDebrid struct {
	// baseURL is the base URL for alldebrid
	baseURL string

	// client is an HTTP Client for the current session
	client *http.Client

	// Allow to set a custom HTTP transport (for test purposes)
	CustomTransport http.RoundTripper

	SupportedHostsRegex []*regexp.Regexp
}

// allDebridResponse represents parts of a debrid response
type allDebridResponse struct {
	Error string
	Link  string
}

const (
	defaultBaseURL     = "https://alldebrid.com"
	authPath           = "register/"
	returnPath         = "account/"
	debridPath         = "service.php"
	supportedHostsPath = "extension/getSupportedHosts.php"
	defaultTimeOut     = 10 * time.Second
)

// Init initializes some Alldebrid things
func (ad *AllDebrid) Init() error {
	// Building hosts regexp
	ad.buildSupportedHosts()

	return nil
}

// Auth initializes AllDebrid
func (ad *AllDebrid) Auth(infos map[string]string) error {
	var username, password, baseURL string

	username = infos["username"]
	if username == "" {
		return errors.New("Invalid username")
	}

	password = infos["password"]
	if password == "" {
		return errors.New("Invalid password")
	}

	baseURL = infos["base_url"]
	if baseURL != "" {
		ad.baseURL = baseURL
	} else {
		ad.baseURL = defaultBaseURL
	}

	cookieJar, _ := cookiejar.New(nil)

	ad.client = &http.Client{
		Timeout:   defaultTimeOut,
		Transport: ad.CustomTransport,
		Jar:       cookieJar,
	}

	query := url.Values{}
	query.Add("action", "login")
	query.Add("returnpage", returnPath)
	query.Add("login_login", username)
	query.Add("login_password", password)

	finalURL := fmt.Sprintf("%s/%s?%s", ad.baseURL, authPath, query.Encode())

	response, responseError := ad.client.Get(finalURL)
	if responseError != nil {
		return responseError
	}

	defer response.Body.Close()

	newLocation := response.Request.URL.String()
	if newLocation != fmt.Sprintf("%s/%s", ad.baseURL, returnPath) {
		return errors.New("Invalid credentials")
	}

	return nil
}

// Debrid debrid a given uri
func (ad *AllDebrid) Debrid(uri string, options map[string]interface{}) (string, error) {
	query := url.Values{}
	query.Add("link", uri)
	query.Add("json", "true")

	getURL := fmt.Sprintf("%s/%s?%s", ad.baseURL, debridPath, query.Encode())

	// Hum, only GET seems to be supported…
	response, responseError := ad.client.Get(getURL)
	if responseError != nil {
		return "", responseError
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	var debridResponse allDebridResponse

	unmarshallError := json.Unmarshal(body, &debridResponse)
	if unmarshallError != nil {
		return "", unmarshallError
	}

	if debridResponse.Error != "" {
		return "", errors.New(debridResponse.Error)
	}

	return debridResponse.Link, nil
}

func (ad *AllDebrid) buildSupportedHosts() {
	hostsCount := len(allDebridSupportedHosts)
	hostRegexps := make([]*regexp.Regexp, hostsCount)

	for hostIndex, hostRegexp := range allDebridSupportedHosts {
		hostRegexps[hostIndex] = regexp.MustCompile(hostRegexp)
	}

	ad.SupportedHostsRegex = hostRegexps
}

// IsDebridable indicates if an uri is supported by AllDebrid
func (ad *AllDebrid) IsDebridable(uri string) bool {
	var hasMatch bool

	if len(ad.SupportedHostsRegex) < 1 {
		return false
	}

	for _, hostRegexp := range ad.SupportedHostsRegex {
		hasMatch = hostRegexp.MatchString(uri)

		if hasMatch {
			return true
		}
	}

	return false
}
