package feedwatcher

import (
	"time"

	"github.com/davidderus/christopher/config"
)

// RemoteFeed is the access informations of a feed on the Internet
type RemoteFeed struct {
	Title           string                 // Remote feed title
	URL             string                 // URL to the feed
	Provider        string                 // The feed provider
	ProviderOptions config.ProviderOptions // The feed provider options
	remoteFeedItems []*RemoteFeedItem      // Storing last parsed feed for functions to consume
}

// RemoteFeedItem represents a simplified RSS item
// Here we only includes relevant infos for future download
type RemoteFeedItem struct {
	Title       string
	Link        string
	Description string
	PublishedAt time.Time
}

func (rf *RemoteFeed) itemsSince(date time.Time) []*RemoteFeedItem {
	newItems := []*RemoteFeedItem{}

	for _, feedItem := range rf.remoteFeedItems {
		if feedItem.PublishedAt.After(date) {
			newItems = append(newItems, feedItem)
		}
	}

	return newItems
}

// DownloadLink extracts the link from a RemoteFeedItem
// thanks to a given FeedExtractor
func (rfi *RemoteFeedItem) DownloadLink(extractor FeedExtractor) string {
	return extractor.Extract(rfi)
}

// NewItems returns the feed new items since the given date
func (rf *RemoteFeed) NewItems(sinceDate time.Time, feedParserFunction FeedParser) ([]*RemoteFeedItem, error) {
	parsedFeedItems, parsingError := feedParserFunction(rf.URL)

	if parsingError != nil {
		return nil, parsingError
	}

	rf.remoteFeedItems = parsedFeedItems

	return rf.itemsSince(sinceDate), nil
}

// NewItemsLinks returns the feed new items links since the given date
func (rf *RemoteFeed) NewItemsLinks(sinceDate time.Time, feedParserFunction FeedParser) ([]string, error) {
	newItems, newItemsError := rf.NewItems(sinceDate, feedParserFunction)

	if newItemsError != nil {
		return nil, newItemsError
	}

	links := make([]string, len(newItems))

	extractor, extractorError := NewFeedExtractor(rf.Provider, rf.ProviderOptions)
	if extractorError != nil {
		return nil, extractorError
	}

	for index, item := range newItems {
		links[index] = item.DownloadLink(extractor)
	}

	return links, nil
}
