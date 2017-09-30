package feedwatcher_test

import (
	"github.com/davidderus/christopher/feedwatcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewFeedExtractor", func() {
	Context("With a valid feed extractor", func() {
		It("Should return an instantiated feed extractor", func() {
			feedExtractor, feedExtractorError := feedwatcher.NewFeedExtractor("DirectDownload", nil)

			Expect(feedExtractorError).NotTo(HaveOccurred())

			newExtractor := &feedwatcher.DirectDownload{}
			newExtractor.Init()
			Expect(feedExtractor).To(BeEquivalentTo(newExtractor))
		})
	})

	Context("With an invalid feed extractor", func() {
		It("Should return an error", func() {
			_, feedExtractorError := feedwatcher.NewFeedExtractor("Fake", nil)

			Expect(feedExtractorError.Error()).To(Equal("Invalid Feed Extractor"))
		})
	})
})
