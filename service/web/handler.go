package web

import (
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

	fd := new(template.FrontendData)
	fd.LoggedIn = false
	tp := template.GetTemplates()

	path := strings.Split(reqPath, "/")[1:]
	page := path[0]
	if len(page) == 0 {
		page = "about"
	}

	if h.pageVisibility[page] == "" {
		// TODO: redirect to 404
		return
	}

	if !fd.LoggedIn && h.pageVisibility[page] == "private" {
		// TODO: redirect to login
	}

	fd.Page = page
	log.Println(fd.Page)

	tp.Execute(w, fd)
}
