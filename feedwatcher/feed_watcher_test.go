package feedwatcher_test

import (
	"fmt"
	"os"
	"time"

	"github.com/davidderus/christopher/config"
	. "github.com/davidderus/christopher/feedwatcher"

	"github.com/mmcdole/gofeed"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFeedWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FeedWatcher Suite")
}

// feedSinceDateWithItems is the date since when we can get 3 items from both
// basic_feed and directdownload_feed
var feedSinceDateWithItems = time.Date(2017, time.February, 25, 0, 0, 0, 0, time.UTC)
var feedSinceDateWithoutItems = time.Date(2017, time.March, 25, 0, 0, 0, 0, time.UTC)

// customFeedParser is a test parser using gofeed to parse an xml
// from the testdata
func customFeedParser(feedURL string) ([]*RemoteFeedItem, error) {
	feedParser := gofeed.NewParser()

	feedData, feedDataError := os.Open(fmt.Sprintf("../testdata/%s_feed.xml", feedURL))

	if feedDataError != nil {
		panic("Can't read given test feed")
	}

	defer feedData.Close()

	parsedFeed, feedParsingError := feedParser.Parse(feedData)

	if feedParsingError != nil {
		panic("Error while parsing test feed")
	}

	remoteFeedItems := make([]*RemoteFeedItem, len(parsedFeed.Items))

	for index, feedItem := range parsedFeed.Items {
		remoteFeedItems[index] = &RemoteFeedItem{
			Title:       feedItem.Title,
			Link:        feedItem.Link,
			Description: feedItem.Description,
			PublishedAt: *feedItem.PublishedParsed}
	}

	return remoteFeedItems, nil
}

var _ = Describe("FeedWatcher", func() {
	var feedWatcher FeedWatcher

	BeforeEach(func() {
		remoteFeeds := make([]RemoteFeed, 1)
		remoteFeeds[0] = RemoteFeed{Title: "Feed One", URL: "directdownload", Provider: "DirectDownload"}

		feedWatcher = FeedWatcher{Feeds: remoteFeeds, Parser: customFeedParser}
	})

	Describe(".NewLinks()", func() {
		Context("With new items", func() {
			It("should only return feeds with new items", func() {
				downloadLinks, _ := feedWatcher.NewLinks(feedSinceDateWithItems)

				Expect(len(downloadLinks)).To(Equal(3))
			})

			It("should return the correct links", func() {
				downloadLinks, _ := feedWatcher.NewLinks(feedSinceDateWithItems)

				Expect(downloadLinks).To(Equal([]string{"http://www.filefactory.com/file/Zombie-One.mkv", "http://rapidgator.net/file/HTGAWM.mkv", "http://www.filefactory.com/file/Shark-Avocado.mkv"}))
			})

			Context("and some options", func() {
				It("should return the correct links", func() {
					providerOptions := config.ProviderOptions{FavoriteHosts: []string{"uploaded.net", "rapidgator.net"}}
					remoteFeeds := make([]RemoteFeed, 1)
					remoteFeeds[0] = RemoteFeed{Title: "Feed One", URL: "directdownload", Provider: "DirectDownload", ProviderOptions: providerOptions}

					feedWatcher = FeedWatcher{Feeds: remoteFeeds, Parser: customFeedParser}
				})
			})
		})

		Context("With no items", func() {
			It("should return nothing", func() {
				downloadLinks, _ := feedWatcher.NewLinks(feedSinceDateWithoutItems)

				Expect(len(downloadLinks)).To(BeZero())
				Expect(downloadLinks).To(BeNil())
			})
		})

		Context("With invalid items", func() {
			It("should log an error", func() {
				feedWatcher.Parser = failingFeedParser
				_, newItemsError := feedWatcher.NewLinks(feedSinceDateWithItems)

				Expect(newItemsError).To(HaveOccurred())
				Expect(newItemsError.Error()).To(Equal("Feed One: Fatal Feed error"))
			})
		})

		Context("with multiple feeds", func() {
			It("should works too", func() {
				remoteFeeds := make([]RemoteFeed, 3)
				remoteFeeds[0] = RemoteFeed{Title: "Feed One", URL: "directdownload", Provider: "DirectDownload"}
				remoteFeeds[1] = RemoteFeed{Title: "Feed Two", URL: "directdownload", Provider: "DirectDownload"}
				remoteFeeds[2] = RemoteFeed{Title: "Feed Three", URL: "directdownload", Provider: "DirectDownload"}

				feedWatcher.Feeds = remoteFeeds

				downloadLinks, _ := feedWatcher.NewLinks(feedSinceDateWithItems)

				Expect(len(downloadLinks)).To(Equal(9))
			})
		})
	})

	Describe(".Run()", func() {
		It("should work", func() {
			// Creating a new feedWatcher
			feedWatcher, _ := NewFeedWatcher(5 * time.Microsecond)

			// Setting an artificial SinceDate
			feedWatcher.SinceDate = time.Date(2017, time.February, 25, 0, 0, 0, 0, time.UTC)

			// Using customFeedParser instead of default parser for testing purpose
			feedWatcher.Parser = customFeedParser

			// Building a RemoteFeed with 3 expected results
			firstFeed := RemoteFeed{Title: "Run Feed", URL: "directdownload", Provider: "DirectDownload"}

			// Building 1 feed (1*3 results expected)
			feedWatcher.Feeds = []RemoteFeed{firstFeed}

			// Running only twice (~10 microseconds max)
			runSummary, _ := feedWatcher.Run(2)
			Expect(runSummary).To(Equal("2 runs done, 3 items found"))
		})

		It("should exit if there is no feeds", func() {
			feedWatcher, _ := NewFeedWatcher(1 * time.Microsecond)
			_, runError := feedWatcher.Run(1)

			Expect(runError.Error()).To(Equal("No feeds in config"))
		})
	})
})
