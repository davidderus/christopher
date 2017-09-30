package downloader_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dnaeon/go-vcr/recorder"

	. "github.com/davidderus/christopher/downloader"
)

const (
	validRpcURL = "http://127.0.0.1:6800/jsonrpc"
	RpcToken    = "my-good-token"
)

// note: stop has to be handled manually
func getRecorder(cassette string) *recorder.Recorder {
	recording, recordingError := recorder.New(fmt.Sprintf("../testdata/cassettes/aria_downloader/%s", cassette))
	if recordingError != nil {
		Fail(recordingError.Error())
	}

	return recording
}

func getClientForCassette(cassette string) (*Aria2, *recorder.Recorder) {
	authInfos := make(map[string]interface{})
	authInfos["token"] = RpcToken
	authInfos["rpcURL"] = validRpcURL

	testRecorder := getRecorder(cassette)
	ariaDownloader := &Aria2{CustomTransport: testRecorder}

	authError := ariaDownloader.Auth(authInfos)
	if authError != nil {
		Fail(authError.Error())
	}

	return ariaDownloader, testRecorder
}

var _ = Describe("Aria2", func() {
	Describe(".Auth()", func() {
		Context("With valid auth infos", func() {
			It("Should log the user in", func() {
				authInfos := make(map[string]interface{})
				authInfos["token"] = RpcToken
				authInfos["rpcURL"] = validRpcURL

				ariaDownloader := &Aria2{}
				authError := ariaDownloader.Auth(authInfos)

				Expect(authError).NotTo(HaveOccurred())
			})
		})

		Context("With invalid auth infos", func() {
			It("Should return an error on missing token", func() {
				authInfos := make(map[string]interface{})
				authInfos["token"] = 12121
				authInfos["rpcURL"] = validRpcURL

				ariaDownloader := &Aria2{}
				authError := ariaDownloader.Auth(authInfos)

				Expect(authError.Error()).To(Equal("Invalid token"))
			})

			It("Should not return an error on missing token", func() {
				authInfos := make(map[string]interface{})
				authInfos["token"] = ""
				authInfos["rpcURL"] = validRpcURL

				ariaDownloader := &Aria2{}
				authError := ariaDownloader.Auth(authInfos)

				Expect(authError).NotTo(HaveOccurred())
			})

			It("Should return an error on missing url", func() {
				authInfos := make(map[string]interface{})
				authInfos["token"] = RpcToken
				authInfos["rpcURL"] = ""

				ariaDownloader := &Aria2{}
				authError := ariaDownloader.Auth(authInfos)

				Expect(authError.Error()).To(Equal("Invalid RPC url"))
			})
		})
	})

	Context("Once authenticated", func() {
		Describe(".Download()", func() {
			Context("with an HTTP Link", func() {
				It("Should return a unique ID for the download", func() {
					ariaDownloader, testRecorder := getClientForCassette("download_without_options")

					gid, _ := ariaDownloader.Download("http://google.fr", nil)

					testRecorder.Stop()

					Expect(gid).To(Equal("96676fbc46cbbc04"))
				})

				It("Should accept download options", func() {
					ariaDownloader, testRecorder := getClientForCassette("download_with_options")

					downloadOptions := make(map[string]interface{})
					downloadOptions["max-overall-download-limit"] = "512K"

					gid, _ := ariaDownloader.Download("http://google.fr", downloadOptions)

					testRecorder.Stop()

					Expect(gid).To(Equal("002eda8439d70942"))
				})
			})

			Context("with an invalid link", func() {
				It("Should return an error", func() {
					ariaDownloader, testRecorder := getClientForCassette("download_with_invalid_link")

					_, downloadError := ariaDownloader.Download("not-a-link", nil)

					testRecorder.Stop()

					Expect(downloadError.Error()).To(Equal("No URI to download."))
				})
			})
		})

		Describe(".DownloadStatus()", func() {
			Context("With a valid GID", func() {
				It("Should return the status of a download", func() {
					ariaDownloader01, testRecorder01 := getClientForCassette("download_without_options")
					gid, _ := ariaDownloader01.Download("http://google.fr", nil)
					testRecorder01.Stop()

					ariaDownloader02, testRecorder02 := getClientForCassette("download_status_with_valid_gid")
					status, _ := ariaDownloader02.DownloadStatus(gid)
					testRecorder02.Stop()

					Expect(status["status"]).To(Equal("active"))
				})
			})

			Context("With an invalid GID", func() {
				It("Should return an error", func() {
					ariaDownloader, testRecorder := getClientForCassette("download_status_with_invalid_gid")
					_, statusError := ariaDownloader.DownloadStatus("111")
					testRecorder.Stop()

					Expect(statusError.Error()).To(Equal("GID 111 is not found"))
				})
			})
		})
	})
})
