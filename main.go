package main

import (
	"github.com/Fliegermarzipan/gallipot/template"
	"os"
	"log"
	"strings"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	reqPath := "/login"

	fd := new(template.FrontendData)
	fd.LoggedIn = false
	tp := template.GetTemplates()

	path := strings.Split(reqPath, "/")[1:]
	page := path[0]
	if len(page) == 0 {
		page = "about"
	}

	pagesPublic := []string{"about", "signup", "login"}
	pagesPrivate := []string{"dashboard"}
	pagesAll := append(append([]string{}, pagesPublic...), pagesPrivate...)

	if !contains(pagesAll, page) {
		// TODO: redirect to 404
	}

	if !fd.LoggedIn && contains(pagesPrivate, page) {
		// TODO: redirect to login
	}

	fd.Page = page
	log.Println(fd.Page)

	tp.Execute(os.Stdout, fd)
}