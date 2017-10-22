package dispatcher

import (
	"github.com/davidderus/christopher/config"
	"github.com/davidderus/christopher/debrider"
	"github.com/davidderus/christopher/downloader"
	"github.com/davidderus/christopher/teller"
)

// ChristopherStory is a story to handle a new URI
//
// Currenty, it does the following:
//
// - Checking if URI is debridable
// - If URI is debridable:
//  - Debriding the URI using default debrider
//  - Sending debrided URI to downloader
// - If URI is not debridable:
//  - Sending URI to downloader
//- Notifying the user about the download event
type ChristopherStory struct {
	notifierFunc func(event *Event) error

	withDebrider   bool
	withDownloader bool

	config *config.Config

	teller *teller.Teller
}

const (
	debridableStep = "debridable"
	debriderStep   = "debrider"
	downloaderStep = "downloader"
)

// Scenario is the main scenario of ChristopherStory
func (cs *ChristopherStory) Scenario() *Scenario {
	var (
		afterConfigStepName string
		afterDebridStepName string
		debriderConfig      *config.DebriderOptions
		debriderInstance    debrider.Debrider
		downloaderConfig    *config.DownloaderOptions
		dlInstance          downloader.Downloader
		err                 error
		isDebridable        bool
	)

	// By default we explicitly do nothing
	afterConfigStepName = "doNothing"
	isDebridableFunc := func() bool { return isDebridable }

	scenario := &Scenario{}

	// Defining a config step for current scenario
	if cs.withDownloader {
		afterConfigStepName = downloaderStep
		afterDebridStepName = downloaderStep

		cs.teller.Log().Debugln("Enabling downloader")
	}

	if cs.withDebrider {
		afterConfigStepName = debridableStep

		cs.teller.Log().Debugln("Enabling debrider")
	}

	scenario.From("config").To(afterConfigStepName).Do(func(_ *Event) error {
		cs.teller.Log().Debugln("Loading config")

		debriderConfig = &cs.config.Debrider
		downloaderConfig = &cs.config.Downloader

		return nil
	})

	scenario.From(debridableStep).To(debriderStep).Do(func(event *Event) error {
		var tempDebriderInstance debrider.Debrider

		tempDebriderInstance, err = debrider.NewDebrider(debriderConfig.Name, nil)
		if err != nil {
			cs.teller.Log().Errorln(err)
			return err
		}

		isDebridable = tempDebriderInstance.IsDebridable(event.Value)

		if isDebridable {
			cs.teller.LogWithFields(map[string]interface{}{
				"debridHandler": debriderConfig.Name,
				"initialURI":    event.Value,
			}).Infoln("URI is debridable")
		}

		return nil
	})

	scenario.From(debriderStep).To("debrided").Do(func(_ *Event) error {
		debriderInstance, err = debrider.NewDebrider(debriderConfig.Name, debriderConfig.AuthInfos)
		if err != nil {
			cs.teller.Log().Errorln(err)
			return err
		}

		return nil
	}).If(isDebridableFunc)

	// Skipping if not debridable to go to downloader
	// afterDebridStepName may be "" if we want to step just after debrid
	scenario.From("debrided").To(afterDebridStepName).Do(func(event *Event) error {
		var debridedURI string

		debridedURI, err = debriderInstance.Debrid(event.Value, nil)
		if err != nil {
			cs.teller.Log().Errorln(err)
			return err
		}

		cs.teller.LogWithFields(map[string]interface{}{
			"debridHandler": debriderConfig.Name,
			"debridURI":     debridedURI,
			"initialURI":    event.Value,
		}).Debugln("URI is debrided")

		event.Origin = debriderStep
		event.Value = debridedURI

		return nil
	}).If(isDebridableFunc)

	scenario.From(downloaderStep).To("downloading").Do(func(_ *Event) error {
		dlInstance, err = downloader.NewDownloader(downloaderConfig.Name, downloaderConfig.AuthInfos)
		if err != nil {
			cs.teller.Log().Errorln(err)
			return err
		}

		return nil
	})

	scenario.From("downloading").To("notified").Do(func(event *Event) error {
		var downloadID string

		downloadID, err = dlInstance.Download(event.Value, downloaderConfig.DownloadOptions)
		if err != nil {
			cs.teller.Log().Errorln(err)
			return err
		}

		cs.teller.LogWithFields(map[string]interface{}{
			"downloadHandler": downloaderConfig.Name,
			"downloadID":      downloadID,
			"downloadOptions": downloaderConfig.DownloadOptions,
			"downloadURI":     event.Value,
		}).Infoln("Download started")

		event.Origin = downloaderStep
		event.Value = downloadID

		return nil
	})

	// Ending current story with a notification
	// NOTE notifierFunc must be set before scenario's play in order for the step
	// to be run
	if cs.notifierFunc != nil {
		scenario.From("notified").Do(cs.notifierFunc)
	}

	// Or ending with a print if no step are used
	scenario.From("doNothing").Do(func(_ *Event) error {
		return nil
	})

	return scenario
}

// SetNotifier defines a nofier for the story
func (cs *ChristopherStory) SetNotifier(notifierFunc func(event *Event) error) *ChristopherStory {
	cs.notifierFunc = notifierFunc
	return cs
}

// EnableDebrider enables the debrider for the story play
func (cs *ChristopherStory) EnableDebrider() *ChristopherStory {
	cs.withDebrider = true
	return cs
}

// EnableDownloader enables the downloader for the story play
func (cs *ChristopherStory) EnableDownloader() *ChristopherStory {
	cs.withDownloader = true
	return cs
}

// SetConfig sets a given config instead of the default one
func (cs *ChristopherStory) SetConfig(config *config.Config) *ChristopherStory {
	cs.config = config
	return cs
}

// SetTeller boots the story teller
func (cs *ChristopherStory) SetTeller(teller *teller.Teller) *ChristopherStory {
	cs.teller = teller
	return cs
}
