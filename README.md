# Christopher v1.0.0-alpha.5

[![Build Status](https://travis-ci.org/davidderus/christopher.svg?branch=master)](https://travis-ci.org/davidderus/christopher)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidderus/christopher)](https://goreportcard.com/report/github.com/davidderus/christopher)
[![GoDoc](https://godoc.org/github.com/davidderus/christopher?status.svg)](https://godoc.org/github.com/davidderus/christopher)

## Description

Christopher is your everyday direct-download companion.

It automatically grabs new episodes from RSS feeds, debrids their URL if needed,
and send them to a downloader service.

It also offers an integrated webserver for remote URL submission and accepts
debrid and download instruction from the command line.

## Usage

```shell
# Starts a FeedWatcher
christopher feed-watcher
# shorter version: christopher fw

# Runs a webserver with a simple interface to debrid and download URIs
christopher webserver
# shorter version: christopher ws

# Debrids an URI
christopher debrid "http://rapidgator.net/file/HTGAWM.mkv"
# shorter version: christopher de "http://rapidgator.net/file/HTGAWM.mkv"

# Downloads an URI
christopher download "https://google.fr"
# shorter version: christopher do "https://google.fr"

# Debrids and downloads an URI
christopher debrid-download "http://rapidgator.net/file/HTGAWM.mkv"
# shorter version: christopher dedo "http://rapidgator.net/file/HTGAWM.mkv"

# Use a custom config file (default in ~/.config/christopher/config.toml)
christopher -c ~/my/custom/config.toml […]
```

## Configuration

Christopher looks for a toml configuration file at
`$HOME/.config/christopher/config.toml`.

```toml
# Download configuration (required)
# The downloader is an external service Christopher pushes links to.
[downloader]
  # Name of the service downloading URIs
  name = "aria2"

  [downloader.auth_infos]
    token = "my-good-token"
    rpc_url = "http://127.0.0.1:6800/jsonrpc"

# Debrider configuration (optional)
# The debrider converts links from specific services to a downloadable link.
# Each link sent to Christopher is first tested against each debriders
# to see if it can be debrided.
[debrider]
  name = "AllDebrid"

  [debrider.auth_infos]
    username = "valid-username"
    password = "valid-password"
    base_url = "https://alldebrid.com"

# FeedWatcher configuration (optional)
# The feedwatcher watch some feeds and send every new links
# to the debriders/downloader
[feedwatcher]
  # Defining a custom watch interval in minutes (default to 30 min.)
  watch_interval = 5

  # Adding one feed for the feedwatcher to look for
  [[feedwatcher.feeds]]
    title = "DirectDownload Feed"

    # URL to the provider feed
    url = "https://directdownload.tv"

    provider = "DirectDownload"

# Providers configuration (optional)
# For each feed provider, you can setup some specific config like a list
# of host to send to the debrider.
[providers]
  # Provider config key and feed provider key must be the same
  [providers.DirectDownload]
    # If specified, favorite_hosts will be looked for in the provider links.
    # If none of them are found, nothing will be downloaded.
    # If not specified, the first link available is downloaded.
    favorite_hosts = ["uploaded.net", "rapidgator.net"]

# WebServer configuration (optional)
# The webserver accepts valid URLs and sent them to Christopher
# debriders/downloader
[webserver]
  # Setting a custom host (default to "127.0.0.1")
  host = "0.0.0.0"

  # Setting a custom webserver port (default to 8000)
  port = 8080

  # Defining a secret for the CSRF token generation (required)
  secret = "my-strong-secret"

  # Constraining CSRF cookie to be HTTPS only if true (default to false)
  secure_cookie = true

  # Setting a realm for the HTTP digest auth (default to "christopher.local")
  auth_realm = "download-helper.local"

  # Users are a list of allowed users.
	# If no users are given, no Digest auth is setup.
  [[webserver.users]]
  name = "johndoe"

  # Setting a password via a MD5 string using:
  # `echo -n "$username:download-helper.local:$password" | md5`
  # Also works with all the algorithms supported by HTTP Digest
  password = "36f0730e47562fbdf9b91434130f91f2"

# Teller configuration (optional)
# The Teller handles logging across the whole application.
# It supports text and JSON logging with variables.
[teller]
  # A level from when the Teller must log things
  # Default is `info`, meaning info logs and above will be shown.
  log_level = "debug"

  # The log items format
  # Default is `text`
  log_formatter = "json"
```

## Supported services

Here is a complete list of all supported services.

Do not hesitate to write a PR with some tests to add more.

### Providers

- DirectDownload (`provider = "DirectDownload" # or dd, directdownload, directdownload.tv`)

### Debriders

- AllDebrid (`name = "AllDebrid" # or alldebrid, Alldebrid, ad`)

### Downloaders

- Aria2 (`name = "aria2" # or Aria2, aria`)

## Upcoming features

- [x] A complete logger with log levels handling
- [ ] A download history to avoid duplicates and show status in webserver
- [ ] A successful download notifier (_push or email_)
- [ ] A lighter Docker Image
- [ ] A better communication between the webserver and Christopher core
- [ ] A documentation about the dispatcher and its stories/scenarios

## Docker

You can run Christopher with Docker with the following command:

```shell
# First build the image
docker build -t christopher .

# As a command line debrider
docker run -v $HOME/.config/christopher/config.toml:/christopher/config.toml:ro christopher debrid "http://rapidgator.net/file/HTGAWM.mkv"
```

In order to download files, you must give the christopher container access
to a Downloader.

Looking at the `docker-compose.yml` file, you will see a basic example of a
working christopher + aria2 setup. Try it with `docker-compose up`.

Remember that you need to update your christopher config file with the right
aria2 RPC token and use `http://aria2:6800/jsonrpc` as the `rpc_url`.

## Development

To build bindata file, use `go-bindata -pkg webserver -o webserver/bindata.go -prefix webserver webserver/templates/`.

## Licence

MIT Licence. Click [here](LICENCE) to see the full text.
