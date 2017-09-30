package feedwatcher

import "errors"

// FeedExtractor takes feed items and return urls to download
type FeedExtractor interface {
	Init() error
	SetOptions(options interface{}) error
	SourceURL() string
	Extract(feedItem *RemoteFeedItem) string
}

// NewFeedExtractor returns an extractor with some options
func NewFeedExtractor(name string, options interface{}) (FeedExtractor, error) {
	var extractor FeedExtractor

	switch name {
	case "DirectDownload", "directdownload", "dd", "directdownload.tv":
		extractor = &DirectDownload{}
	default:
		return nil, errors.New("Invalid Feed Extractor")
	}

	if options != nil {
		optionsError := extractor.SetOptions(options)
		if optionsError != nil {
			return nil, optionsError
		}
	}

	extractor.Init()

	return extractor, nil
}
