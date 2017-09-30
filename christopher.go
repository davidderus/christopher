package main

import (
	"os"

	"github.com/davidderus/christopher/config"
	"github.com/urfave/cli"
)

var appConfigPath string

func christopherApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Christopher"
	app.Usage = "Your everyday direct-download companion"
	app.Version = christopherVersion

	// @see commands.go
	app.Commands = []cli.Command{
		FeedWatcherCli,
		DownloaderCli,
		DebriderCli,
		DownloadAndDebridCli,
		WebServerCli,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       config.DefaultConfigPath(),
			Usage:       "Load configuration from `FILE`",
			Destination: &appConfigPath,
		},
	}

	return app
}

func main() {
	app := christopherApp()
	app.Run(os.Args)
}
