// Package config defines and parses configuration for controllers, cameras,
// notifiers and watchers
package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

// CameraOptions lists the options allowed for a camera
type CameraOptions struct {
	// Device address like /dev/video0
	Device string

	// Custom input (default is 8)
	Input int

	// Basic Motion options
	Width     int
	Height    int
	Framerate int

	// MotionThreshold is `threshold` in Motion config
	MotionThreshold int `toml:"motion_threshold"`

	// EventGap is `gap` in Motion config
	EventGap int `toml:"event_gap"`

	// Role is one of []string{"stream", "watch"}
	Role string

	// Autostart defines if the camera should be started at boot
	Autostart bool `toml:"auto_start"`
}

// NotifierOptions includes the option for a notifier
type NotifierOptions struct {
	// Notifying service name
	Service string

	// Notifications recipients
	Recipients []string

	// Service options
	ServiceOptions map[string]string `toml:"options"`
}

type webUser struct {
	Name     string
	Password string
}

// WebServerOptions defines some of the webserver options
type WebServerOptions struct {
	Port      int
	Host      string
	AuthRealm string `toml:"auth_realm"`

	User []webUser
}

// Config is the default config object
type Config struct {
	Port int
	Host string

	// Countdown before a notification is sent
	Countdown int

	// Path to motion binary
	MotionPath string `toml:"motion_path"`

	// Directory where logs and generated config files are stored
	WorkingDir string `toml:"working_dir"`

	// Listing of Camera with their options
	Cameras map[string]*CameraOptions

	// All cameras with a watch role will use the given Notifiers
	Notifiers map[string]*NotifierOptions

	// Defines the webserver options
	WebServer WebServerOptions
}

// TemplatesDirectory is where the main and thread config are stored
const TemplatesDirectory = "templates"

// ConfigDirectoryName is the name for the thread configs directory
const ConfigDirectoryName = "configs"

// LogsDirectoryName is the name for the directory where the motion logs are stored
const LogsDirectoryName = "logs"

// CapturesDirectoryName is the main folder where all the pictures and videos are saved
const CapturesDirectoryName = "captures"

// MainConfigFileTemplate is the default motion config
const MainConfigFileTemplate = "motion.conf.tpl"

// ThreadBaseName is the model name for a thread configuration file
const ThreadBaseName = "dicam-thread-%s"

// DefaultConfigMode is the file mode for a config file
const DefaultConfigMode = 0700

// DefaultHost sets up the controller and webserver host as Internet-open hosts
const DefaultHost = "0.0.0.0"

// DefaultAuthRealm is the default realm used for the web server digest authentication
const DefaultAuthRealm = "dicam.local"

// DefaultWaitTime is the time before firing an event
//
// This is set in order not to immediately alert when detecting a motion and
// letting some time for the user to deactivate the notifier (ie: when entering
// his property)
const DefaultWaitTime = 10

// Read reads config for dicam
func Read() (*Config, error) {
	user, userError := user.Current()
	if userError != nil {
		return nil, userError
	}

	userHomeDir := user.HomeDir

	configFullPath := path.Join(userHomeDir, ".config/dicam/config.toml")

	// Initializing configuration
	var config Config

	// Setting defaults for required elements
	config.setDefaults(userHomeDir)

	// Parsing config file, thus overriding defaults if needed
	_, configError := toml.DecodeFile(configFullPath, &config)
	if configError != nil {
		return nil, configError
	}

	// Validating resulting configuration
	validationError := config.validate()
	if validationError != nil {
		return nil, validationError
	}

	// Setting up working directory for later use
	populateError := config.populateWorkingDir()
	if populateError != nil {
		return nil, populateError
	}

	return &config, nil
}

// setDefaults defines default options in configuration such as motion path,
// controller port and host…
func (c *Config) setDefaults(userDir string) {
	defaultMotionPath, _ := exec.LookPath("motion")

	c.Countdown = DefaultWaitTime
	c.Port = 8888
	c.Host = DefaultHost
	c.MotionPath = defaultMotionPath
	c.WorkingDir = path.Join(userDir, ".dicam")

	c.WebServer.Host = DefaultHost
	c.WebServer.Port = 8000
	c.WebServer.AuthRealm = DefaultAuthRealm
}

// validate validates a few config options to prevent further errors
func (c *Config) validate() error {
	if c.Port == 0 {
		return errors.New("App port is invalid")
	}

	if c.MotionPath == "" {
		return errors.New("Path to motion is invalid or motion is not available")
	}

	return nil
}

// populateWorkingDir creates the configs and logs directories based on the
// WorkingDir
func (c *Config) populateWorkingDir() error {
	userDirError := os.MkdirAll(c.WorkingDir, DefaultConfigMode)
	if userDirError != nil {
		return userDirError
	}

	mkdirConfigError := os.MkdirAll(path.Join(c.WorkingDir, ConfigDirectoryName), DefaultConfigMode)
	if mkdirConfigError != nil {
		return mkdirConfigError
	}

	mkdirLogsError := os.MkdirAll(path.Join(c.WorkingDir, LogsDirectoryName), DefaultConfigMode)
	if mkdirLogsError != nil {
		return mkdirLogsError
	}

	mkdirCapturesError := os.MkdirAll(path.Join(c.WorkingDir, CapturesDirectoryName), DefaultConfigMode)
	if mkdirCapturesError != nil {
		return mkdirCapturesError
	}

	return nil
}

// ListCamsToStart returns ids of cameras to start at boot time
func (c *Config) ListCamsToStart() []string {
	availableCams := c.Cameras
	toStart := []string{}

	for name, config := range availableCams {
		if config.Autostart == true {
			toStart = append(toStart, name)
		}
	}

	return toStart
}

// GetCameraOptions returns the CameraOptions for a given cameraID
func (c *Config) GetCameraOptions(cameraID string) (*CameraOptions, error) {
	availableCams := c.Cameras

	for id, options := range availableCams {
		if id == cameraID {
			return options, nil
		}
	}

	return nil, fmt.Errorf("No options available for camera %s", cameraID)
}
