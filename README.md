# Christopher v0.0.0

## Usage

```shell
# Starts a FeedWatcher
christopher feed-watcher
# shorter version: christopher fw

# Debrids an URI
christopher debrid "http://rapidgator.net/file/HTGAWM.mkv"
# shorter version: christopher de "http://rapidgator.net/file/HTGAWM.mkv"

# Downloads an URI
christopher download "https://google.fr"
# shorter version: christopher do "https://google.fr"

# Debrids and downloads an URI
christopher debrid-download "http://rapidgator.net/file/HTGAWM.mkv"
# shorter version: christopher dedo "http://rapidgator.net/file/HTGAWM.mkv"

# Runs a webserver with a simple interface to debrid and download URIs
christopher webserver
# shorter version: christopher ws

# Use a custom config file (default in ~/.config/christopher/config.toml)
christopher -c ~/my/custom/config.toml […]
```

## Configuration

See the example below for a list of the available options:

```toml
# Setting a custom database path (default to ~/.config/christopher/config.toml)
config_path = "/Users/tom/christopher/config.toml"

# Setting a custom db path (default to ~/.config/christopher/database.db)
db_path = "/Users/tom/christopher/mybase.db"

[feedwatcher]
  # Defining a custom watch interval in minutes (default to 30 min.)
  watch_interval = 5

  # Adding one feed for the feedwatcher to look for
  [[feedwatcher.feeds]]
    title = "DirectDownload Feed"

    # URL to the provider feed
    url = "https://directdownload.tv"

    provider = "DirectDownload"

[downloader]
  # Name of the service downloading URIs
  name = "aria2"

  [downloader.auth_infos]
    token = "my-good-token"
    rpcURL = "http://127.0.0.1:6800/jsonrpc"

[debrider]
  name = "AllDebrid"

  [debrider.auth_infos]
    username = "valid-username"
    password = "valid-password"
    base_url = "https://alldebrid.com"

[providers]
  # Provider config key and feed provider key must be the same
  [providers.DirectDownload]
    # If specified, favorite_hosts will be looked for in the provider links.
    # If none of them are found, nothing will be downloaded.
    # If not specified, the first link available is downloaded.
    favorite_hosts = ["uploaded.net", "rapidgator.net"]

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
```

## Supported services

Here is a complete list of all supported services.

Do not hesitate to write a PR to add more.

### Providers

- DirectDownload (`provider = "DirectDownload" # or dd, directdownload, directdownload.tv`)

### Debriders

- AllDebrid (`name = "AllDebrid" # or alldebrid, Alldebrid, ad`)

### Downloaders

- Aria2 (`name = "aria2" # or Aria2, aria`)

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
aria2 RPC token and use `http://aria2:6800/jsonrpc` as the ` rpcURL`.

## Licence

MIT Licence. Click [here](LICENCE) to see the full text.
