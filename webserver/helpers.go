package webserver

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/davidderus/christopher/dispatcher"
)

func (ws *WebServer) writeWithTemplate(response http.ResponseWriter, templateName string, templatePath string, data interface{}) error {
	templates := template.New("")

	var (
		templateDir  string
		templateData []byte
		assetError   error
	)

	templatesToRender := []string{"layout.html", "navbar.html", "scripts.js", templatePath}

	for _, templateName := range templatesToRender {
		templateDir = filepath.Join("templates", templateName)
		templateData, assetError = Asset(templateDir)

		// Skipping template on asset error
		if assetError == nil {
			templates.New(templateDir).Parse(string(templateData))
		}
	}

	templates.ExecuteTemplate(response, "layout", data)

	return nil
}

func (ws *WebServer) loadScenario() *dispatcher.Scenario {
	story := &dispatcher.ChristopherStory{}
	story.SetConfig(ws.appConfig).EnableDebrider().EnableDownloader()
	story.SetTeller(ws.appTeller)

	scenario := story.Scenario()
	scenario.SetInitialStep("config")

	return scenario
}
