package downloader_test

import (
	"github.com/davidderus/christopher/downloader"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDownloader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downloader Suite")
}

var _ = Describe("NewDownloader", func() {
	Context("With a valid downloader", func() {
		It("Should return an instantiated downloader", func() {
			dlInstance, downloaderError := downloader.NewDownloader("Aria2", nil)

			Expect(downloaderError).NotTo(HaveOccurred())
			Expect(dlInstance).To(BeEquivalentTo(&downloader.Aria2{}))
		})
	})

	Context("With an invalid downloader", func() {
		It("Should return an error", func() {
			_, downloaderError := downloader.NewDownloader("Fake", nil)

			Expect(downloaderError.Error()).To(Equal("Invalid downloader given"))
		})
	})
})
