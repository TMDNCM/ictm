package web

import (
	"github.com/Fliegermarzipan/gallipot/template"
	"github.com/Fliegermarzipan/gallipot/data"
	htemplate "html/template"
	"log"
	"net/http"
	"strings"
)

type WebHandler struct {
	t              *htemplate.Template
	pageVisibility map[string]string
}

func NewHandler() *WebHandler {
	h := new(WebHandler)
	h.t = template.GetTemplates()

	h.pageVisibility = map[string]string{
		"about":     "public",
		"signup":    "public",
		"login":     "public",
		"dashboard": "private",
		"profile":   "private",
	}

	return h
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	reqPath := r.URL.Path
	tp := template.GetTemplates()

	// TODO: stop using fake login
	fd := template.FrontendData{}
	fd.LoggedIn = true
	fd.User = new(data.User)
	fd.User.Username = "someonespecial"
	fd.User.Displayname = "Someone Special"
	fd.User.Email = "foo@example.com"


	path := strings.Split(reqPath, "/")[1:]
	fd.Page = path[0]
	if len(fd.Page) == 0 {
		fd.Page = "about"
	}

	if h.pageVisibility[fd.Page] == "" {
		// TODO: redirect to 404
		return
	}

	if !fd.LoggedIn && h.pageVisibility[fd.Page] == "private" {
		// TODO: redirect to login
	}

	log.Println(fd.Page)

	tp.Execute(w, fd)
}
