package web

import (
	"github.com/Fliegermarzipan/gallipot/data"
	"github.com/Fliegermarzipan/gallipot/template"
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
		"about":         "public",
		"signup":        "public",
		"login":         "public",
		"dashboard":     "private",
		"profile":       "private",
		"notifications": "private",
		"log":           "private",
		"friends":       "private",
		"stock":         "private",
		"user":          "private",
	}

	return h
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	log.Println(reqPath)

	tp := template.GetTemplates()

	// TODO: stop using fake login
	fd := template.FrontendData{}
	fd.LoggedIn = true
	fd.User = new(data.User)
	fd.User.Username = "someonespecial"
	fd.User.Displayname = "Someone Special"
	fd.User.Email = "foo@example.com"

	fd.Path = strings.Split(reqPath, "/")[1:]
	fd.Page = fd.Path[0]
	if len(fd.Page) == 0 {
		fd.Page = "about"
	}

	if h.pageVisibility[fd.Page] == "" {
		fd.Page = "404"
		w.WriteHeader(http.StatusNotFound)
	}

	if !fd.LoggedIn && h.pageVisibility[fd.Page] == "private" {
		// TODO: redirect to login
	}

	// TODO: remove example alert usage
	if fd.Page == "login" {
		fd.Alert = new(template.UserAlert)
		fd.Alert.Title = "Login incorrect"
		fd.Alert.Message = "Wrong username or password."
	}

	tp.Execute(w, fd)
}
