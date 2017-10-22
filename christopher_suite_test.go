package main

import (
	"bytes"
	"fmt"

	"github.com/davidderus/christopher/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"

	"testing"
)

func TestChristopher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Christopher Suite")
}

var _ = Describe("Christopher", func() {
	var cliApp *cli.App

	BeforeEach(func() {
		cliApp = christopherApp()
	})

	Context("version", func() {
		It("should return the correct version number", func() {
			Expect(cliApp.Version).To(Equal(christopherVersion))
		})
	})

	Context("feed-watcher --help", func() {
		It("should show the feed watcher help", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"feed-watcher", "--help"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("Starts the feed watcher"))
		})
	})

	Context("download --help", func() {
		It("should show the downloader help", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"download", "--help"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("Sends an URI to the downloader"))
		})
	})

	Context("debrid --help", func() {
		It("should show the debrider help", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"debrid", "--help"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("Debrids an URI"))
		})
	})

	Context("debrid-download --help", func() {
		It("should show the debrid-downloader help", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"debrid-download", "--help"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("Debrids and downloads an URI"))
		})
	})

	Context("--config", func() {
		It("should show the config flag", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"--help"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("--config"))

			// Also contains a default value for config
			defaultConfigString := fmt.Sprintf("Load configuration from FILE (default: \"%s\")", config.DefaultConfigPath())
			Expect(fwOutput).To(ContainSubstring(defaultConfigString))
		})
	})

	Context("webserver", func() {
		It("should show the webserver help", func() {
			cliBuffer := new(bytes.Buffer)
			cliApp.Writer = cliBuffer

			fwErr := cliApp.Run([]string{"webserver"})
			fwOutput := cliBuffer.String()

			Expect(fwErr).To(BeNil())
			Expect(fwOutput).To(ContainSubstring("Starts the webserver"))
		})
	})
})
