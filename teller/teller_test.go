package teller_test

import (
	"bytes"

	. "github.com/davidderus/christopher/teller"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTeller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Teller Suite")
}

var _ = Describe("Teller", func() {
	var consoleBuffer *bytes.Buffer
	var teller *Teller

	BeforeEach(func() {
		consoleBuffer = new(bytes.Buffer)
		teller = NewTeller("debug", "text")
		teller.SetLogOutput(consoleBuffer)
	})

	Describe("Log", func() {
		Context("with debug level", func() {
			It("should log everything", func() {
				teller.Log().Infoln("This is an info log")
				teller.Log().Debugln("This is a debug log")
				teller.Log().Errorf("This is an %s log", "error")
				teller.Log().WithField("url", "https://google.fr").Warnln("This is a warning log")

				consoleOutput := consoleBuffer.String()

				Expect(consoleOutput).To(ContainSubstring("level=info msg=\"This is an info log"))
				Expect(consoleOutput).To(ContainSubstring("level=debug msg=\"This is a debug log"))
				Expect(consoleOutput).To(ContainSubstring("level=error msg=\"This is an error log"))
				Expect(consoleOutput).To(ContainSubstring("level=warning msg=\"This is a warning log\" url=\"https://google.fr\""))
			})
		})

		Context("with a specific level", func() {
			It("should log this level and everything above", func() {
				By("Setting log level to Info")
				teller.SetLogLevel("error")

				teller.Log().Debugln("This is a debug log")
				teller.Log().Infoln("This is an info log")
				teller.Log().Warnln("This is an warning log")
				teller.Log().Errorln("This is an error log")

				consoleOutput := consoleBuffer.String()

				Expect(consoleOutput).To(ContainSubstring("level=error msg=\"This is an error log"))
				Expect(consoleOutput).NotTo(ContainSubstring("level=info msg=\"This is an info log"))
				Expect(consoleOutput).NotTo(ContainSubstring("level=warning msg=\"This is an warning log"))
				Expect(consoleOutput).NotTo(ContainSubstring("level=debug msg=\"This is an debug log"))
			})
		})

		It("should handle JSON", func() {
			teller.SetLogFormatter("json")

			teller.Log().Debugln("This is a debug log")

			consoleOutput := consoleBuffer.String()
			expectedJSON := `{"level":"debug","msg":"This is a debug log","time":".*"}`

			Expect(consoleOutput).To(MatchRegexp(expectedJSON))
		})
	})

	Describe("LogWithFields", func() {
		It("should log all fields", func() {
			teller.SetLogFormatter("json")

			teller.LogWithFields(map[string]interface{}{
				"url":  "https://google.fr",
				"step": "Step1",
			}).Debugln("This is a debug log")

			consoleOutput := consoleBuffer.String()
			expectedJSON := `{"level":"debug","msg":"This is a debug log","step":"Step1","time":".*","url":"https://google.fr"}`

			Expect(consoleOutput).To(MatchRegexp(expectedJSON))
		})
	})
})
