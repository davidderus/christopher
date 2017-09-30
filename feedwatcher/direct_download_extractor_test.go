package feedwatcher_test

import (
	"github.com/davidderus/christopher/config"
	. "github.com/davidderus/christopher/feedwatcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FeedExtractor", func() {
	Context("for DirectDownload", func() {
		Context("With a known feed", func() {
			var myRemoteFeed RemoteFeed
			var newItems []*RemoteFeedItem
			var extractor FeedExtractor
			var extractorError error

			BeforeEach(func() {
				myRemoteFeed = RemoteFeed{Title: "New items feed", URL: "directdownload", Provider: "DirectDownload"}
				newItems, _ = myRemoteFeed.NewItems(feedSinceDateWithItems, customFeedParser)
				extractor, extractorError = NewFeedExtractor(myRemoteFeed.Provider, nil)
			})

			It("Should return the right extractor", func() {
				feedExtractor := &DirectDownload{}
				feedExtractor.Init()
				Expect(extractor).To(BeEquivalentTo(feedExtractor))
			})

			Context("Without any config", func() {
				It("should return first download links", func() {
					downloadLinks := make([]string, 3)

					for index, item := range newItems {
						downloadLinks[index] = item.DownloadLink(extractor)
					}

					Expect(downloadLinks).To(Equal([]string{"http://www.filefactory.com/file/Zombie-One.mkv", "http://rapidgator.net/file/HTGAWM.mkv", "http://www.filefactory.com/file/Shark-Avocado.mkv"}))
				})
			})

			Context("With a specified host order", func() {
				It("should return the more accurate download link", func() {
					options := config.ProviderOptions{FavoriteHosts: []string{"uploaded.net", "rapidgator.net"}}

					customExtractor, _ := NewFeedExtractor(myRemoteFeed.Provider, options)
					downloadLinks := make([]string, 3)

					for index, item := range newItems {
						downloadLinks[index] = item.DownloadLink(customExtractor)
					}

					Expect(len(downloadLinks)).To(Equal(3))

					Expect(downloadLinks).To(Equal([]string{"http://uploaded.net/file/Zombie-One.mkv", "http://rapidgator.net/file/HTGAWM.mkv", ""}))
				})
			})
		})

		Context("With an unknown feed", func() {
			It("should return nil", func() {
				remoteFeed := RemoteFeed{Title: "New items feed", URL: "basic", Provider: "BBDown"}
				extractor, _ := NewFeedExtractor(remoteFeed.Provider, nil)

				Expect(extractor).To(BeNil())
			})
		})
	})
})
