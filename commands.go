package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davidderus/christopher/config"
	"github.com/davidderus/christopher/dispatcher"
	"github.com/davidderus/christopher/feedwatcher"
	"github.com/davidderus/christopher/webserver"
	"github.com/urfave/cli"
)

var (
	appConfig *config.Config
)

//////////////
// Debrider //
//////////////

// DebriderCli defines the cli args for the debrider
var DebriderCli = cli.Command{
	Name:        "debrid",
	Aliases:     []string{"de"},
	Usage:       "Debrid an URI",
	Description: "Send a URI to the debrider and return a debrided URI.",
	Action:      runDebrider,
	ArgsUsage:   "<URI>",
}

func runDebrider(ctx *cli.Context) error {
	uri := ctx.Args().First()
	if uri == "" {
		return cli.NewExitError("No URI given", 1)
	}

	appConfig, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return cli.NewExitError(configError.Error(), 1)
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDebrider()

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		return cli.NewExitError(runError.Error(), 2)
	}

	fmt.Print(event.Value)

	return nil
}

////////////////
// Downloader //
////////////////

// DownloaderCli defines the cli args for the downloader
var DownloaderCli = cli.Command{
	Name:        "download",
	Aliases:     []string{"do"},
	Usage:       "Send an URI to the downloader",
	Description: "Send a URI to the downloader set in config.",
	Action:      runDownloader,
	ArgsUsage:   "<URI>",
}

func runDownloader(ctx *cli.Context) error {
	uri := ctx.Args().First()
	if uri == "" {
		return cli.NewExitError("No URI given", 1)
	}

	appConfig, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return cli.NewExitError(configError.Error(), 1)
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	// story without debrider
	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDownloader()

	// Setting a custom notifier
	story.SetNotifier(func(event *dispatcher.Event) error {
		log.Printf("%s sent to downloader.\n", event.Value)

		return nil
	})

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		return cli.NewExitError(runError.Error(), 2)
	}

	return nil
}

/////////////////////////
// Download and Debrid //
/////////////////////////

// DownloadAndDebridCli defines the cli args for a combination of the debrider and the downloader
var DownloadAndDebridCli = cli.Command{
	Name:        "debrid-download",
	Action:      downloadAndDebrid,
	Aliases:     []string{"dedo"},
	Usage:       "Debrid and download an URI",
	Description: "If debrider is able to handle URI, it will be debrided and then sent to downloader. Otherwise, it will only be downloaded.",
	ArgsUsage:   "<URI>",
}

func downloadAndDebrid(ctx *cli.Context) error {
	uri := ctx.Args().First()
	if uri == "" {
		return cli.NewExitError("No URI given", 1)
	}

	appConfig, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return cli.NewExitError(configError.Error(), 1)
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDebrider().EnableDownloader()

	// Setting a custom notifier
	story.SetNotifier(func(event *dispatcher.Event) error {
		log.Printf("Download started with ID %s\n", event.Value)

		return nil
	})

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		return cli.NewExitError(runError.Error(), 2)
	}

	return nil
}

/////////////////
// FeedWatcher //
/////////////////

// FeedWatcherCli defines the cli args for the feed watcher
var FeedWatcherCli = cli.Command{
	Name:        "feed-watcher",
	Aliases:     []string{"fw"},
	Usage:       "Start feed watcher",
	Description: "Watch the feeds defined in configuration and send all links to the right service.",
	Action:      runFeedWatcher,
}

func runFeedWatcher(ctx *cli.Context) error {
	// Loading config
	appConfig, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return cli.NewExitError(configError.Error(), 1)
	}

	// Building FeedWatcher from config
	feedWatcherConfig := appConfig.FeedWatcher
	watchInterval := time.Duration(feedWatcherConfig.WatchInterval) * time.Minute
	feedWatcher, feedWatcherError := feedwatcher.NewFeedWatcher(watchInterval)
	if feedWatcherError != nil {
		return cli.NewExitError(feedWatcherError.Error(), 1)
	}

	// Getting feeds
	configFeeds := feedWatcherConfig.Feeds
	feedWatcherFeeds := make([]feedwatcher.RemoteFeed, len(configFeeds))
	for feedIndex, feed := range configFeeds {
		feedProvider := feed.Provider

		providerOptions := appConfig.Providers[feedProvider]

		feedWatcherFeeds[feedIndex] = feedwatcher.RemoteFeed{
			Title:           feed.Title,
			URL:             feed.URL,
			Provider:        feedProvider,
			ProviderOptions: providerOptions,
		}
	}

	feedWatcher.Feeds = feedWatcherFeeds

	// Artificial SinceDate for now
	feedWatcher.SinceDate = time.Now()

	// Using default story to process new links
	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDebrider().EnableDownloader()

	// Setting a custom log notifier
	story.SetNotifier(func(event *dispatcher.Event) error {
		log.Printf("Download started with ID %s\n", event.Value)

		return nil
	})

	scenario := story.Scenario()
	scenario.SetInitialStep("config")

	feedWatcher.Scenario = scenario

	// Running FeedWatcher for eternity
	log.Println("Starting FeedWatcher")
	runSummary, runError := feedWatcher.Run(0)

	// Handling run errors
	if runError != nil {
		return cli.NewExitError(runError.Error(), 1)
	}

	// Logging output if any (will never be reached on Run(0))
	log.Println(runSummary)

	return nil
}

///////////////
// WebServer //
///////////////

// WebServerCli defines the cli args for the feed watcher
var WebServerCli = cli.Command{
	Name:    "webserver",
	Aliases: []string{"ws"},
	Usage:   "Starts the webserver",
	Action:  runWebServer,
}

func runWebServer(ctx *cli.Context) error {
	// Loading config
	appConfig, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return cli.NewExitError(configError.Error(), 1)
	}

	webServer := webserver.NewWebServer(appConfig)

	log.Fatal(webServer.Start())

	return nil
}
