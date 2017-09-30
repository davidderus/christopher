package webserver

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// HomeHandler prints the christopher homepage
func (ws *WebServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	ws.writeWithTemplate(w, "Home", "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
