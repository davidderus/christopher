package config_test

import (
	. "github.com/davidderus/christopher/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	var config *Config
	var loadError error
	validConfigSampleFile := "../testdata/config_valid_sample.toml"

	Describe("Load()", func() {
		Context("With no user config file", func() {
			It("should return an error", func() {
				_, loadError := Load()

				Expect(loadError).NotTo(BeNil())
				Expect(loadError.Error()).To(ContainSubstring("Missing config file"))
			})
		})
	})

	Describe(".LoadFromFile()", func() {
		Context("with a valid file", func() {
			BeforeEach(func() {
				config, loadError = LoadFromFile(validConfigSampleFile)
			})

			It("should load the config", func() {
				Expect(loadError).NotTo(HaveOccurred())
				Expect(config).NotTo(BeNil())

				// FeedWatcher
				By("Parsing FeedWatcher config")
				feedWatcherConfig := config.FeedWatcher
				Expect(len(feedWatcherConfig.Feeds)).To(Equal(1))
				Expect(feedWatcherConfig.Feeds[0].Provider).To(Equal("DirectDownload"))

				// Downloader
				By("Parsing Downloader config")
				downloaderConfig := config.Downloader
				Expect(downloaderConfig.Name).To(Equal("aria2"))
				Expect(downloaderConfig.AuthInfos["token"]).To(Equal("my-good-token"))
				Expect(downloaderConfig.AuthInfos["rpcURL"]).To(Equal("http://127.0.0.1:6800/jsonrpc"))

				// Debrider
				By("Parsing Debrider config")
				debriderConfig := config.Debrider
				Expect(debriderConfig.Name).To(Equal("AllDebrid"))
				Expect(debriderConfig.AuthInfos["username"]).To(Equal("valid-username"))
				Expect(debriderConfig.AuthInfos["password"]).To(Equal("valid-password"))
				Expect(debriderConfig.AuthInfos["base_url"]).To(Equal("https://alldebrid.com"))

				By("Parsing Providers config")
				providersConfig := config.Providers
				Expect(len(providersConfig)).To(Equal(1))

				ddProviderConfig, configExists := providersConfig["DirectDownload"]
				Expect(configExists).To(BeTrue())
				Expect(ddProviderConfig.FavoriteHosts).To(Equal([]string{"uploaded.net", "rapidgator.net"}))

				By("Parsing WebServer config")
				webserverConfig := config.WebServer
				Expect(webserverConfig.AuthRealm).To(Equal("christopher.local")) // Default
				Expect(webserverConfig.Host).To(Equal("127.0.0.1"))
				Expect(webserverConfig.Port).To(Equal(8080))
				Expect(webserverConfig.Secret).To(Equal("Ahgho7aKetho4aiceiVoa3eiKu0chouY"))
				Expect(webserverConfig.SecureCookie).To(BeFalse())

				webserverUsers := webserverConfig.Users
				Expect(len(webserverUsers)).To(Equal(1))
				Expect(webserverUsers[0].Name).To(Equal("johndoe"))
			})

			It("should set some defaults for the missing values", func() {
				Expect(config.FeedWatcher.WatchInterval).To(Equal(30))
			})

			It("should not set defaults for existing values", func() {
				Expect(config.DBPath).NotTo(ContainSubstring("database.db"))
			})

			It("should load some Feeds", func() {
				firstFeed := config.FeedWatcher.Feeds[0]

				// Checking that firstFeed item is valid
				Expect(firstFeed.Title).To(Equal("DirectDownload Feed"))
				Expect(firstFeed.URL).To(Equal("https://directdownload.tv"))
				Expect(firstFeed.Provider).To(Equal("DirectDownload"))
			})
		})

		Context("with an erroneous configuration", func() {
			It("return an error", func() {
				_, loadError := LoadFromFile("../testdata/config_invalid_sample.toml")
				Expect(loadError.Error()).To(Equal("DBPath can't be blank"))
			})
		})

		Context("with an invalid file", func() {
			It("return an error", func() {
				_, loadError := LoadFromFile("../testdata/basic_feed.xml")
				Expect(loadError.Error()).To(ContainSubstring("Invalid config file"))
			})
		})

		Context("with a missing file", func() {
			It("return an error", func() {
				_, loadError := LoadFromFile("../testdata/missing_config.toml")
				Expect(loadError.Error()).To(ContainSubstring("Missing config file"))
			})
		})
	})

	Describe(".UserDir()", func() {
		It("should return the current user home directory", func() {
			userDir, dirError := UserDir()

			Expect(userDir).To(BeADirectory())
			Expect(dirError).NotTo(HaveOccurred())
		})
	})
})
