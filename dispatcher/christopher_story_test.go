package dispatcher_test

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/davidderus/christopher/config"
	. "github.com/davidderus/christopher/dispatcher"
	"github.com/davidderus/christopher/teller"

	"github.com/dnaeon/go-vcr/recorder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const validConfigSampleFile = "../testdata/config_valid_sample.toml"

// note: stop has to be handled manually
func getRecorder(cassette string) *recorder.Recorder {
	recording, recordingError := recorder.New(fmt.Sprintf("../testdata/cassettes/christopher_story/%s", cassette))
	if recordingError != nil {
		Fail(recordingError.Error())
	}

	return recording
}

var _ = Describe("ChristopherStory", func() {
	var story *ChristopherStory
	var appConfig *config.Config
	var defaultHTTPTransport http.RoundTripper
	var logBuffer *bytes.Buffer
	var tellerInstance *teller.Teller

	BeforeEach(func() {
		defaultHTTPTransport = http.DefaultTransport
		appConfig, _ = config.LoadFromFile(validConfigSampleFile)

		logBuffer = &bytes.Buffer{}

		// Getting logger
		tellerInstance = teller.NewTeller("debug", "text")
		tellerInstance.SetLogOutput(logBuffer)
	})

	// Resetting http.DefaultTransport to initial value to avoid messing with
	// the test suite
	AfterEach(func() {
		http.DefaultTransport = defaultHTTPTransport
	})

	Context("with Downloader", func() {
		Context("without Debrider", func() {
			It("should download the link", func() {
				Expect(appConfig.Downloader.Name).To(Equal("aria2"))

				testRecorder := getRecorder("downloader")
				// NOTE May not be a good idea
				http.DefaultTransport = testRecorder

				event := &Event{Origin: "test", Value: "http://google.fr"}

				story = &ChristopherStory{}
				story.SetConfig(appConfig).EnableDownloader()
				story.SetTeller(tellerInstance)

				scenario := story.Scenario()
				scenario.SetInitialStep("config")
				scenario.Play(event)

				testRecorder.Stop()

				Expect(scenario.RunError()).To(BeNil())
				Expect(event.Value).To(Equal("96676fbc46cbbc04"))
				Expect(event.Origin).To(Equal("downloader"))
			})
		})

		Context("with Debrider", func() {
			It("should debrid and download the link", func() {
				Expect(appConfig.Debrider.Name).To(Equal("AllDebrid"))

				testRecorder := getRecorder("debrider_and_downloader")
				// NOTE May not be a good idea
				http.DefaultTransport = testRecorder

				event := &Event{Origin: "test", Value: "http://rapidgator.net/file/08987898765/HTGAWM.mkv"}

				story = &ChristopherStory{}
				story.SetConfig(appConfig).EnableDebrider().EnableDownloader()
				story.SetTeller(tellerInstance)

				scenario := story.Scenario()
				scenario.SetInitialStep("config")
				scenario.Play(event)

				testRecorder.Stop()

				Expect(scenario.RunError()).To(BeNil())
				Expect(event.Value).To(Equal("96676fbc46cbbaaz"))
				Expect(event.Origin).To(Equal("downloader"))
			})
		})

		Context("with an undebridable link", func() {
			It("should still try", func() {
				Expect(appConfig.Debrider.Name).To(Equal("AllDebrid"))

				testRecorder := getRecorder("failing_debrider_and_downloader")
				// NOTE May not be a good idea
				http.DefaultTransport = testRecorder

				// Here Value is an invalid URI format for AllDebrid
				event := &Event{Origin: "test", Value: "http://rapidgator.net/HTGAWM.mkv"}

				story = &ChristopherStory{}
				story.SetConfig(appConfig).EnableDebrider().EnableDownloader()
				story.SetTeller(tellerInstance)

				scenario := story.Scenario()
				scenario.SetInitialStep("config")
				scenario.Play(event)

				testRecorder.Stop()

				Expect(scenario.RunError()).To(BeNil())
				Expect(event.Value).To(Equal("98676zbc46c00c31"))
				Expect(event.Origin).To(Equal("downloader"))
			})
		})

		Context("reusing the same story for multiple events", func() {
			It("should work", func() {
				story = &ChristopherStory{}
				story.SetConfig(appConfig).EnableDebrider().EnableDownloader()
				story.SetTeller(tellerInstance)

				scenario := story.Scenario()

				By("Using a valid link")

				testRecorder := getRecorder("debrider_and_downloader")
				http.DefaultTransport = testRecorder

				event := &Event{Origin: "test", Value: "http://rapidgator.net/file/08987898765/HTGAWM.mkv"}

				scenario.SetInitialStep("config")
				scenario.Play(event)

				testRecorder.Stop()

				Expect(scenario.RunError()).To(BeNil())
				Expect(event.Value).To(Equal("96676fbc46cbbaaz"))
				Expect(event.Origin).To(Equal("downloader"))

				By("Using an invalid link")

				testRecorder = getRecorder("failing_debrider_and_downloader")
				http.DefaultTransport = testRecorder

				// Here Value is an invalid URI format for AllDebrid
				event = &Event{Origin: "test", Value: "http://rapidgator.net/HTGAWM.mkv"}

				// Initial step must be reset
				scenario.SetInitialStep("config")
				scenario.Play(event)

				testRecorder.Stop()

				Expect(scenario.RunError()).To(BeNil())
				Expect(event.Value).To(Equal("98676zbc46c00c31"))
				Expect(event.Origin).To(Equal("downloader"))
			})
		})
	})

	Context("without Downloader and Debrider", func() {
		It("should do nothing", func() {
			event := &Event{Origin: "test", Value: "http://google.fr"}
			story = &ChristopherStory{}
			story.SetConfig(appConfig)
			story.SetTeller(tellerInstance)

			scenario := story.Scenario()
			scenario.SetInitialStep("config")
			scenario.Play(event)

			Expect(scenario.RunError()).To(BeNil())
			Expect(event.Value).To(Equal("http://google.fr"))
			Expect(event.Origin).To(Equal("test"))
		})
	})

	Context("with a logger", func() {
		It("should log some things", func() {
			By("Configuring the story")
			testRecorder := getRecorder("debrider_and_downloader")
			// NOTE May not be a good idea
			http.DefaultTransport = testRecorder

			event := &Event{Origin: "test", Value: "http://rapidgator.net/file/08987898765/HTGAWM.mkv"}

			story = &ChristopherStory{}
			story.SetConfig(appConfig).EnableDebrider().EnableDownloader()

			By("Setting up the teller")
			story.SetTeller(tellerInstance)

			By("By playing the scenario")
			scenario := story.Scenario()
			scenario.SetInitialStep("config")
			scenario.Play(event)

			logString := logBuffer.String()

			Expect(scenario.RunError()).To(BeNil())
			Expect(logString).To(ContainSubstring(`level=debug msg="Enabling downloader"`))
			Expect(logString).To(ContainSubstring(`level=debug msg="Enabling debrider"`))
			Expect(logString).To(ContainSubstring(`level=debug msg="Loading config"`))

			Expect(logString).To(ContainSubstring(`level=info msg="URI is debridable" debridHandler=AllDebrid initialURI="http://rapidgator.net/file/08987898765/HTGAWM.mkv"`))

			Expect(logString).To(ContainSubstring(`level=debug msg="URI is debrided" debridHandler=AllDebrid debridURI="https://subdomain.alld.io/dl/ABC/HTGAWM.mkv" initialURI="http://rapidgator.net/file/08987898765/HTGAWM.mkv"`))

			Expect(logString).To(ContainSubstring(`level=info msg="Download started" downloadHandler=aria2 downloadID=96676fbc46cbbaaz downloadOptions=map[] downloadURI="https://subdomain.alld.io/dl/ABC/HTGAWM.mkv"`))
		})
	})
})
