package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

const (
	// defaultWatchInterval defines the time in minutes between two FeedWatcher scans
	defaultWatchInterval = 30

	// defaultDBName is the default filename for the database
	defaultDBName = "database.db"

	// defaultHost sets up the webserver host default as a local host
	defaultHost = "127.0.0.1"

	// defaultPort sets up the webserver port to a non-standard one
	defaultPort = 8000

	// defaultAuthRealm is the default realm used for the web server digest authentication
	defaultAuthRealm = "christopher.local"

	// defaultLogLevel sets the minimum level for logging infos
	defaultLogLevel = "info"

	// defaultLogFormatter sets the log items formatter
	defaultLogFormatter = "text"
)

// Feed is a Feed Representation
type feed struct {
	Title    string // Remote feed title
	URL      string // URL to the feed
	Provider string // The feed provider
}

// FeedWatcherOptions defines some options for the FeedWatcher
type feedWatcherOptions struct {
	WatchInterval int `toml:"watch_interval"` // In Minutes
	Feeds         []*feed
}

// DownloaderOptions defines options for the downloader
type DownloaderOptions struct {
	Name            string
	AuthInfos       map[string]interface{} `toml:"auth_infos"`
	DownloadOptions map[string]interface{} `toml:"download_options"`
}

// DebriderOptions defines name and auth info for the debrider
type DebriderOptions struct {
	Name      string
	AuthInfos map[string]string `toml:"auth_infos"`
}

// ProviderOptions specify options for a given provider
type ProviderOptions struct {
	FavoriteHosts []string `toml:"favorite_hosts"`
}

type webUser struct {
	Name     string
	Password string
}

// WebServerOptions defines some of the webserver options
type WebServerOptions struct {
	Port int
	Host string

	// Secret is the secret for the CSRF token generation
	Secret string

	// SecureCookie constrains CSRF cookie to be HTTPS only if true
	SecureCookie bool `toml:"secure_cookie"`

	// AuthRealm is the realm for the HTTP digest auth
	AuthRealm string `toml:"auth_realm"`

	// Users are a list of allowed users.
	// If no users are given, no Digest auth is setup.
	Users []webUser
}

// TellerOptions defines logging options for the Teller
type TellerOptions struct {
	// LogLevel is a level from when the Teller must log things
	LogLevel string `toml:"log_level"`

	// LogFormatter is the log items format
	LogFormatter string `toml:"log_formatter"`
}

// Config defines the Christopher configuration
type Config struct {
	configPath string `toml:"config_path"`

	DBPath      string `toml:"db_path"`
	FeedWatcher feedWatcherOptions

	Downloader DownloaderOptions

	Debrider DebriderOptions

	Providers map[string]ProviderOptions

	WebServer WebServerOptions

	Teller TellerOptions
}

// Load loads config from the default config path
func Load() (*Config, error) {
	return LoadFromFile(DefaultConfigPath())
}

// LoadFromFile loads the configuration from a file
func LoadFromFile(configPath string) (*Config, error) {
	// Initializing configuration
	var config Config

	// Setting defaults for required elements
	config.setDefaults()

	// Parsing config file, thus overriding defaults if needed
	_, configError := toml.DecodeFile(configPath, &config)
	if configError != nil {
		if os.IsNotExist(configError) {
			return nil, fmt.Errorf("Missing config file in %s", configPath)
		}

		return nil, fmt.Errorf("Invalid config file: %v", configError)
	}

	// Validating resulting configuration
	validationError := config.validate()
	if validationError != nil {
		return nil, validationError
	}

	return &config, nil
}

// UserDir returns the user home directory
func UserDir() (string, error) {
	user, userError := user.Current()
	if userError != nil {
		return "", userError
	}

	return user.HomeDir, nil
}

// DefaultConfigPath returns the default path to the config file
func DefaultConfigPath() string {
	userDir, _ := UserDir()
	return path.Join(userDir, ".config", "christopher", "config.toml")
}

func (c *Config) validate() error {
	// Validating DBPath
	if c.DBPath == "" {
		return errors.New("DBPath can't be blank")
	}

	// Must have 32 bytes secret for CSRF protection
	if c.WebServer.Secret == "" {
		return errors.New("A 32 bytes secret token must be set")
	}

	return nil
}

func (c *Config) setDefaults() {
	userDir, _ := UserDir()

	// Setting default DBPath
	c.DBPath = path.Join(userDir, ".config", "christopher", defaultDBName)

	// Setting default WatchInterval
	c.FeedWatcher.WatchInterval = defaultWatchInterval

	// WebServer defaults
	c.WebServer.Host = defaultHost
	c.WebServer.Port = defaultPort
	c.WebServer.AuthRealm = defaultAuthRealm

	c.Teller.LogLevel = defaultLogLevel
	c.Teller.LogFormatter = defaultLogFormatter
}
