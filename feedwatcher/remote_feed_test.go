package feedwatcher_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/davidderus/christopher/config"
	. "github.com/davidderus/christopher/feedwatcher"
)

func failingFeedParser(feedURL string) ([]*RemoteFeedItem, error) {
	return nil, errors.New("Fatal Feed error")
}

var _ = Describe("RemoteFeed", func() {
	var myRemoteFeed RemoteFeed

	Context("With new items", func() {
		It("should return new items", func() {
			myRemoteFeed = RemoteFeed{Title: "New items feed", URL: "basic", Provider: "BasicProvider"}

			newItems, _ := myRemoteFeed.NewItems(feedSinceDateWithItems, customFeedParser)
			newItemsLength := len(newItems)

			newItemsTitles := make([]string, newItemsLength)

			for index, item := range newItems {
				newItemsTitles[index] = item.Title
			}

			Expect(newItemsLength).To(Equal(3))
			Expect(newItemsTitles).To(Equal([]string{"Zombie One", "How to Get Away With Mummy", "Shark Avocado"}))
		})
	})

	Context("With no new items", func() {
		It("should be nil", func() {
			myRemoteFeed = RemoteFeed{Title: "New items feed", URL: "basic", Provider: "BasicProvider"}

			newItems, newItemsError := myRemoteFeed.NewItems(feedSinceDateWithoutItems, customFeedParser)

			Expect(len(newItems)).To(BeZero())
			Expect(newItemsError).NotTo(HaveOccurred())
		})
	})

	Context("With invalid items", func() {
		It("should log an error", func() {
			_, newItemsError := myRemoteFeed.NewItems(feedSinceDateWithItems, failingFeedParser)

			Expect(newItemsError).To(HaveOccurred())
			Expect(newItemsError.Error()).To(Equal("Fatal Feed error"))
		})
	})

	Context("A specific feed", func() {
		Context("With new items", func() {
			It("should return new items links", func() {
				myRemoteFeed := RemoteFeed{Title: "New items feed", URL: "directdownload", Provider: "DirectDownload"}

				links, linksError := myRemoteFeed.NewItemsLinks(feedSinceDateWithItems, customFeedParser)

				Expect(linksError).NotTo(HaveOccurred())

				Expect(links).To(Equal([]string{"http://www.filefactory.com/file/Zombie-One.mkv", "http://rapidgator.net/file/HTGAWM.mkv", "http://www.filefactory.com/file/Shark-Avocado.mkv"}))
			})
		})

		Context("With new items and some provider's options", func() {
			It("should return new items links", func() {
				// Setting some config options
				providerOptions := config.ProviderOptions{FavoriteHosts: []string{"uploaded.net", "rapidgator.net"}}
				myRemoteFeed := RemoteFeed{Title: "New items feed", URL: "directdownload", Provider: "DirectDownload", ProviderOptions: providerOptions}

				links, linksError := myRemoteFeed.NewItemsLinks(feedSinceDateWithItems, customFeedParser)

				Expect(linksError).NotTo(HaveOccurred())

				Expect(links).To(Equal([]string{"http://uploaded.net/file/Zombie-One.mkv", "http://rapidgator.net/file/HTGAWM.mkv", ""}))
			})
		})
	})
})
