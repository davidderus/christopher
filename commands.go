package main

import (
	"fmt"
	"time"

	"github.com/davidderus/christopher/config"
	"github.com/davidderus/christopher/dispatcher"
	"github.com/davidderus/christopher/feedwatcher"
	"github.com/davidderus/christopher/teller"
	"github.com/davidderus/christopher/webserver"
	"github.com/urfave/cli"
)

var (
	appConfig *config.Config
	appTeller *teller.Teller
)

func loadRequirements() error {
	configError := loadConfig()

	if configError != nil {
		return configError
	}

	loadTeller()

	return nil
}

func loadConfig() error {
	config, configError := config.LoadFromFile(appConfigPath)
	if configError != nil {
		return configError
	}

	appConfig = config
	return nil
}

func loadTeller() {
	appTeller = teller.NewTeller(appConfig.Teller.LogLevel, appConfig.Teller.LogFormatter)
}

//////////////
// Debrider //
//////////////

// DebriderCli defines the cli args for the debrider
var DebriderCli = cli.Command{
	Name:        "debrid",
	Aliases:     []string{"de"},
	Usage:       "Debrids an URI",
	Description: "Sends an URI to the debrider and return a debrided URI.",
	Action:      runDebrider,
	ArgsUsage:   "<URI>",
}

func runDebrider(ctx *cli.Context) error {
	// Loading command requirements
	loadError := loadRequirements()
	if loadError != nil {
		return cli.NewExitError(loadError.Error(), 1)
	}

	uri := ctx.Args().First()
	if uri == "" {
		appTeller.Log().Fatalln("No URI given")
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDebrider()
	story.SetTeller(appTeller)

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		appTeller.Log().Fatalln(runError)
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
	Usage:       "Sends an URI to the downloader",
	Description: "Sends a URI to the downloader set in the config file.",
	Action:      runDownloader,
	ArgsUsage:   "<URI>",
}

func runDownloader(ctx *cli.Context) error {
	// Loading command requirements
	loadError := loadRequirements()
	if loadError != nil {
		return cli.NewExitError(loadError.Error(), 1)
	}

	uri := ctx.Args().First()
	if uri == "" {
		appTeller.Log().Fatalln("No URI given")
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	// story without debrider
	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDownloader()
	story.SetTeller(appTeller)

	// Setting a custom notifier
	story.SetNotifier(func(event *dispatcher.Event) error {
		appTeller.LogWithFields(map[string]interface{}{
			"downloadURI":     event.Value,
			"downloadHandler": appConfig.Downloader.Name,
		}).Infoln("URI sent to downloader")

		return nil
	})

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		appTeller.Log().Fatalln(runError)
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
	Usage:       "Debrids and downloads an URI",
	Description: "If debrider is able to handle the URI, it will be debrided and then sent to downloader. Otherwise, it will only be downloaded.",
	ArgsUsage:   "<URI>",
}

func downloadAndDebrid(ctx *cli.Context) error {
	// Loading command requirements
	loadError := loadRequirements()
	if loadError != nil {
		return cli.NewExitError(loadError.Error(), 1)
	}

	uri := ctx.Args().First()
	if uri == "" {
		appTeller.Log().Fatalln("No URI given")
	}

	event := &dispatcher.Event{Origin: "cli", Value: uri}

	story := &dispatcher.ChristopherStory{}
	story.SetConfig(appConfig).EnableDebrider().EnableDownloader()
	story.SetTeller(appTeller)

	scenario := story.Scenario()

	scenario.SetInitialStep("config")

	scenario.Play(event)

	runError := scenario.RunError()
	if runError != nil {
		appTeller.Log().Fatalln(runError)
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
	Usage:       "Starts the feed watcher",
	Description: "Watch the feeds defined in configuration and send all links to the right service.",
	Action:      runFeedWatcher,
}

func runFeedWatcher(ctx *cli.Context) error {
	// Loading command requirements
	loadError := loadRequirements()
	if loadError != nil {
		return cli.NewExitError(loadError.Error(), 1)
	}

	// Building FeedWatcher from config
	feedWatcherConfig := appConfig.FeedWatcher
	watchInterval := time.Duration(feedWatcherConfig.WatchInterval) * time.Minute
	feedWatcher, feedWatcherError := feedwatcher.NewFeedWatcher(watchInterval)
	if feedWatcherError != nil {
		appTeller.Log().Fatalln(feedWatcherError)
	}

	// Adding Teller
	feedWatcher.SetTeller(appTeller)

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
	story.SetTeller(appTeller)

	scenario := story.Scenario()
	scenario.SetInitialStep("config")

	feedWatcher.Scenario = scenario

	// Running FeedWatcher for eternity
	runSummary, runError := feedWatcher.Run(0)

	// Handling run errors
	if runError != nil {
		appTeller.Log().Fatalln(runError)
	}

	// Logging output if any (will never be reached on Run(0))
	appTeller.Log().Infoln(runSummary)

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
	// Loading command requirements
	loadError := loadRequirements()
	if loadError != nil {
		return cli.NewExitError(loadError.Error(), 1)
	}

	webServer := webserver.NewWebServer(appConfig, appTeller)

	appTeller.Log().Fatalln(webServer.Start())

	return nil
}
