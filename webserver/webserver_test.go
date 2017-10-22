package webserver_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/davidderus/christopher/config"
	"github.com/davidderus/christopher/teller"
	. "github.com/davidderus/christopher/webserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebServer Suite")
}

const validConfigSampleFile = "../testdata/config_valid_sample.toml"

// TODO Test basic auth
var _ = Describe("WebServer", func() {
	var webServer *WebServer

	BeforeEach(func() {
		appConfig, _ := config.LoadFromFile(validConfigSampleFile)
		appTeller := teller.NewTeller(appConfig.Teller.LogLevel, appConfig.Teller.LogFormatter)
		appTeller.SetLogOutput(ioutil.Discard)

		webServer = NewWebServer(appConfig, appTeller)
		webServer.Init()
	})

	Describe("/", func() {
		Context("With no auth", func() {
			It("should return the homepage", func() {
				request, requestError := http.NewRequest("GET", "/", nil)
				Expect(requestError).NotTo(HaveOccurred())

				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(webServer.HomeHandler)

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				bodyString := recorder.Body.String()
				Expect(bodyString).To(ContainSubstring("Submit links"))
				Expect(bodyString).To(ContainSubstring("Christopher"))
			})
		})

		Context("With invalid auth", func() {
			It("should fail", func() {
				request, requestError := http.NewRequest("GET", "/", nil)
				Expect(requestError).NotTo(HaveOccurred())

				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(webServer.LoadHandlerWithAuth(webServer.HomeHandler))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("/submit", func() {
		var jsonString []byte

		BeforeEach(func() {
			jsonString = []byte(`{"urls":"http://google.fr\nhttp://google.com\nhttp://google.de"}`)
		})

		It("should process the given json", func() {
			request, requestError := http.NewRequest("POST", "/submit", bytes.NewBuffer(jsonString))
			Expect(requestError).NotTo(HaveOccurred())

			request.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(webServer.SubmitHandler)

			handler.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))

			validResponse := `{"count":3,"errors":null}`
			Expect(recorder.Body.String()).To(Equal(validResponse))
		})
	})
})
