package feedwatcher

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/davidderus/christopher/dispatcher"
	"github.com/davidderus/christopher/teller"
)

// FeedWatcher is the main process struct
type FeedWatcher struct {
	Feeds     []RemoteFeed
	Parser    FeedParser
	SinceDate time.Time
	Scenario  *dispatcher.Scenario

	interval time.Duration
	teller   *teller.Teller
}

// NewFeedWatcher returns a FeedWatcher for a given time interval
func NewFeedWatcher(interval time.Duration) (*FeedWatcher, error) {
	if interval == 0 {
		return nil, errors.New("Invalid interval")
	}

	feedWatcher := &FeedWatcher{}
	feedWatcher.interval = interval
	feedWatcher.Parser = GofeedParser

	return feedWatcher, nil
}

// feedNewItems get all new items for a given feed
func (fw *FeedWatcher) feedNewItems(feed *RemoteFeed, sinceDate time.Time, linksChan chan []string, errorsChan chan string) {
	newItems, newItemsError := feed.NewItemsLinks(sinceDate, fw.Parser)

	if newItemsError != nil {
		errorsChan <- fmt.Sprintf("%s: %s", feed.Title, newItemsError)
		return
	}

	linksChan <- newItems
}

// NewLinks returns new links across all feeds
func (fw *FeedWatcher) NewLinks(sinceDate time.Time) ([]string, error) {
	feedsCount := len(fw.Feeds)

	var newLinks []string
	newLinksChan := make(chan []string, feedsCount)
	defer close(newLinksChan)

	var errorMessages []string
	errorsMessagesChan := make(chan string, feedsCount)
	defer close(errorsMessagesChan)

	// Parsing feeds concurrently
	for _, feed := range fw.Feeds {
		go fw.feedNewItems(&feed, sinceDate, newLinksChan, errorsMessagesChan)
	}

	// Waiting for answers
	for feedIndex := 0; feedIndex < feedsCount; feedIndex++ {
		select {
		case newItemsLinks := <-newLinksChan:
			newLinks = append(newLinks, newItemsLinks...)
		case newError := <-errorsMessagesChan:
			errorMessages = append(errorMessages, newError)
		}
	}

	var finalError error

	if len(errorMessages) > 0 {
		finalError = errors.New(strings.Join(errorMessages, "\n"))
	}

	return newLinks, finalError
}

// processNewLinks send new links to others (download, debrid…)
// TODO Handle errors
// TODO Allow concurrent dispatch
func (fw *FeedWatcher) processNewLinks(sinceDate time.Time) (int, error) {
	newLinks, linkErrors := fw.NewLinks(sinceDate)

	var currentEvent *dispatcher.Event

	scenario := fw.Scenario
	if scenario != nil {
		for _, newLink := range newLinks {
			currentEvent = &dispatcher.Event{Origin: "feed-watcher", Value: newLink}
			scenario.SetInitialStep("config")
			scenario.Play(currentEvent)
		}
	}

	return len(newLinks), linkErrors
}

// Run starts the FeedWatcher cycle.
//
// It gets new links every tick based on FeedWatcher's interval.
//
// If you want it to stop after a certain number of iteration, use maxRunCount,
// otherwise, it will run forever.
//
// If run is stopped (maxRunCount > 0), then it return a summary of the run
func (fw *FeedWatcher) Run(maxRunCount int) (string, error) {
	feedsCount := len(fw.Feeds)
	if feedsCount == 0 {
		return "", errors.New("No feeds in config")
	}

	if fw.SinceDate.IsZero() {
		return "", errors.New("Invalid SinceDate")
	}

	hasLimit := maxRunCount > 0
	runCount := 0
	newItemsTotal := 0

	tick := time.Tick(fw.interval)

	// TODO replace by previous launch time
	sinceDate := fw.SinceDate

	fw.teller.LogWithFields(map[string]interface{}{
		"feedsCount":  feedsCount,
		"interval":    fw.interval,
		"maxRunCount": maxRunCount,
		"sinceDate":   fw.SinceDate,
	}).Infoln("Starting FeedWatcher")

	// Starts a new go routine every tick to get new links
	for _ = range tick {
		newItemsCount, newItemErrors := fw.processNewLinks(sinceDate)
		if newItemErrors != nil {
			fw.teller.Log().WithField("errors", newItemErrors).Errorln("Error with new items")
		}

		// But also some success
		newItemsTotal += newItemsCount

		// Anyway we keep going
		sinceDate = time.Now()
		runCount++

		fw.teller.Log().WithField("runCount", runCount).Infof("%d new items found", newItemsCount)

		// Except if a given limit is reached
		if hasLimit && runCount == maxRunCount {
			fw.teller.Log().Debugln("Reaching maxRunCount, breaking!")
			break
		}
	}

	return fmt.Sprintf("%d runs done, %d items found", runCount, newItemsTotal), nil
}

// SetTeller defines the teller used to report the feed watcher adventures
func (fw *FeedWatcher) SetTeller(teller *teller.Teller) *FeedWatcher {
	fw.teller = teller
	return fw
}
