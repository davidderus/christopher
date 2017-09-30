package feedwatcher

import "github.com/mmcdole/gofeed"

// FeedParser abstracts a basic parser function
type FeedParser func(feedURL string) ([]*RemoteFeedItem, error)

// GofeedParser is a an abstraction of the gofeed library returning feed items
func GofeedParser(feedURL string) ([]*RemoteFeedItem, error) {
	feedParser := gofeed.NewParser()

	parsedFeed, parsingError := feedParser.ParseURL(feedURL)
	if parsingError != nil {
		return nil, parsingError
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
