package webserver

// In order to embed templates in this package, install `jteeuwen/go-bindata`
// and use the following command: `go-bindata -pkg webserver templates/`

import (
	"fmt"
	"net/http"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/davidderus/christopher/config"
	"github.com/davidderus/christopher/dispatcher"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

// WebServer defines a basic HTTP WebServer
type WebServer struct {
	appConfig     *config.Config
	options       *config.WebServerOptions
	authenticator *auth.DigestAuth
	scenario      *dispatcher.Scenario
	router        *mux.Router
	csrf          func(http.Handler) http.Handler
}

// Init initiates the WebServer struct
func (ws *WebServer) Init() {
	// Enables auth if there is users in config
	if len(ws.options.Users) > 0 {
		ws.enableAuthentication()
	}

	// Enable CSRF
	ws.csrf = csrf.Protect([]byte(ws.options.Secret), csrf.Secure(ws.options.SecureCookie))

	// Building router with routes
	ws.buildRouter()
}

// Start starts the webserver
func (ws *WebServer) Start() error {
	ws.Init()

	webServerAddress := fmt.Sprintf("%s:%d", ws.options.Host, ws.options.Port)

	server := &http.Server{
		Handler: ws.csrf(ws.router),
		Addr:    webServerAddress,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting webserver on %s\n", webServerAddress)
	return server.ListenAndServe()
}

// enableAuthentication activates digest authentication for all requests wrapped
// in LoadHandlerWithAuth
func (ws *WebServer) enableAuthentication() {
	ws.authenticator = auth.NewDigestAuthenticator(ws.options.AuthRealm, ws.lookForSecret)
}

// LoadHandlerWithAuth check for any auth infos in config and use it for authentication.
func (ws *WebServer) LoadHandlerWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	// If no auth infos are found, then no auth is set up.
	if ws.authenticator != nil {
		return auth.JustCheck(ws.authenticator, handler)
	}

	return handler
}

// lookForSecret returns a password hash from config for a given existing user
func (ws *WebServer) lookForSecret(user, realm string) string {
	for _, webUser := range ws.options.Users {
		if webUser.Name == user {
			return webUser.Password
		}
	}

	return ""
}

func (ws *WebServer) buildRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/", ws.LoadHandlerWithAuth(ws.HomeHandler))
	router.HandleFunc("/submit", ws.LoadHandlerWithAuth(ws.SubmitHandler)).Methods("POST")

	ws.router = router
}

// NewWebServer instanciates a web server
func NewWebServer(appConfig *config.Config) *WebServer {
	server := &WebServer{}
	server.appConfig = appConfig
	server.options = &appConfig.WebServer

	return server
}
