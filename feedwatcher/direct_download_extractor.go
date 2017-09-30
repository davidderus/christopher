package feedwatcher

import (
	"fmt"
	"regexp"

	"github.com/davidderus/christopher/config"
)

var urlMatcher = regexp.MustCompile(`(https?:\/\/[\da-z\.-]+\.[a-z\.]{2,6}[\/\w \.-]*\/?)`)

// DirectDownload is an extractor for directdownload.tv
type DirectDownload struct {
	options               config.ProviderOptions
	favoriteHostsMatchers []*regexp.Regexp
}

// Init sets up various extractor's internals
func (dd *DirectDownload) Init() error {
	dd.buildFavoriteHostsMatchers()
	return nil
}

// SourceURL returns DirectDownload host
func (dd *DirectDownload) SourceURL() string {
	return "directdownload.tv"
}

// SetOptions validates and defines the options for DirectDownload
func (dd *DirectDownload) SetOptions(options interface{}) error {
	dd.options = options.(config.ProviderOptions)
	return nil
}

func (dd *DirectDownload) extractByFavoriteHosts(description string) string {
	for _, matcher := range dd.favoriteHostsMatchers {
		match := matcher.FindString(description)

		if match != "" {
			return match
		}
	}

	return ""
}

func (dd *DirectDownload) extractFirst(description string) string {
	return urlMatcher.FindString(description)
}

// buildFavoriteHostsMatchers compiles the FavoriteHosts regexp once and for all
func (dd *DirectDownload) buildFavoriteHostsMatchers() {
	dd.favoriteHostsMatchers = make([]*regexp.Regexp, len(dd.options.FavoriteHosts))

	for hostIndex, host := range dd.options.FavoriteHosts {
		hostURL := fmt.Sprintf("https?://*.?%s.*", host)
		dd.favoriteHostsMatchers[hostIndex] = regexp.MustCompile(hostURL)
	}
}

// Extract extracts links from DirectDownload feed items
func (dd *DirectDownload) Extract(feedItem *RemoteFeedItem) string {
	description := feedItem.Description

	if len(dd.options.FavoriteHosts) > 0 {
		return dd.extractByFavoriteHosts(description)
	}

	return dd.extractFirst(description)
}
