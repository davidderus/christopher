package debrider_test

import (
	"fmt"

	"github.com/dnaeon/go-vcr/recorder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/davidderus/christopher/debrider"
)

var validInfos = [2]string{"valid-username", "valid-password"}

// note: stop has to be handled manually
func getRecorder(cassette string) *recorder.Recorder {
	recording, recordingError := recorder.New(fmt.Sprintf("../testdata/cassettes/alldebrid/%s", cassette))
	if recordingError != nil {
		Fail(recordingError.Error())
	}

	return recording
}

func getClientForCassette(cassette string) (*AllDebrid, *recorder.Recorder) {
	authInfos := make(map[string]string)
	authInfos["username"] = validInfos[0]
	authInfos["password"] = validInfos[1]

	testRecorder := getRecorder(cassette)
	allDebrid := &AllDebrid{CustomTransport: testRecorder}
	allDebrid.Init()

	authError := allDebrid.Auth(authInfos)
	if authError != nil {
		Fail(authError.Error())
	}

	return allDebrid, testRecorder
}

var _ = Describe("AllDebrid", func() {
	Describe(".Auth()", func() {
		Context("With valid auth infos", func() {
			It("should log the user in", func() {
				authInfos := make(map[string]string)
				authInfos["username"] = validInfos[0]
				authInfos["password"] = validInfos[1]

				testRecorder := getRecorder("auth_success")
				allDebrid := &AllDebrid{CustomTransport: testRecorder}
				allDebrid.Init()

				authError := allDebrid.Auth(authInfos)

				testRecorder.Stop()

				Expect(authError).NotTo(HaveOccurred())
			})

			It("should build the hosts regexp", func() {
				allDebrid, testRecorder := getClientForCassette("auth_success")

				testRecorder.Stop()

				Expect(allDebrid.SupportedHostsRegex).NotTo(BeNil())
			})
		})

		Context("With invalid auth infos", func() {
			It("should return an error", func() {
				authInfos := make(map[string]string)
				authInfos["username"] = "wrong-username"
				authInfos["password"] = "wrong-password"

				testRecorder := getRecorder("auth_failure")
				allDebrid := &AllDebrid{CustomTransport: testRecorder}
				allDebrid.Init()

				authError := allDebrid.Auth(authInfos)

				testRecorder.Stop()

				Expect(authError.Error()).To(Equal("Invalid credentials"))
			})
		})
	})

	Describe(".Debrid()", func() {
		Context("With a valid link", func() {
			It("should debrid the link", func() {
				link := "http://rapidgator.net/file/HTGAWM.mkv"

				allDebrid, testRecorder := getClientForCassette("debrid_valid_link")

				debridedLink, debridError := allDebrid.Debrid(link, nil)

				testRecorder.Stop()

				Expect(debridError).NotTo(HaveOccurred())
				Expect(debridedLink).To(Equal("https://subdomain.alld.io/dl/ABC/HTGAWM.mkv"))
			})
		})

		Context("With an unsupported link or offline host", func() {
			It("should fail", func() {
				link := "http://google.fr"

				allDebrid, testRecorder := getClientForCassette("debrid_unsupported_link")

				_, debridError := allDebrid.Debrid(link, nil)

				testRecorder.Stop()

				Expect(debridError.Error()).To(Equal("This link is not valid or not supported"))
			})
		})
	})

	Describe(".IsDebridable()", func() {
		var debrider *AllDebrid

		Context("with Auth()", func() {
			BeforeEach(func() {
				var testRecorder *recorder.Recorder
				debrider, testRecorder = getClientForCassette("auth_success")

				testRecorder.Stop()
			})

			Context("with a compatible host", func() {
				It("should be true", func() {
					var assertion bool

					// Rapidgator
					assertion = debrider.IsDebridable("http://rapidgator.net/file/08987898765/HTGAWM.mkv")
					Expect(assertion).To(BeTrue())

					// Uploaded.net
					assertion = debrider.IsDebridable("http://uploaded.net/file/Zombie-One.mkv")
					Expect(assertion).To(BeTrue())

					// 4shared
					assertion = debrider.IsDebridable("https://4shared.com/file/HTGAWM/dl.html")
					Expect(assertion).To(BeTrue())

					// Ironfiles
					assertion = debrider.IsDebridable("https://mediafire.com/download/sharkavocado/")
					Expect(assertion).To(BeTrue())

					// K2s
					assertion = debrider.IsDebridable("http://mega.co.nz/#!Zombie-MinusOne!AOUE6")
					Expect(assertion).To(BeTrue())
				})
			})

			Context("with an uncompatible or unknown host", func() {
				It("should be false", func() {
					assertion := debrider.IsDebridable("http://google.fr")
					Expect(assertion).To(BeFalse())
				})
			})
		})

		Context("without Auth()", func() {
			It("should work too", func() {
				debrider := &AllDebrid{}

				By("Initializing the debrider first")
				debrider.Init()

				assertion := debrider.IsDebridable("http://rapidgator.net/file/08987898765/HTGAWM.mkv")
				Expect(assertion).To(BeTrue())
			})
		})
	})
})
