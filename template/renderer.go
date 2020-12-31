package template

import (
	//"net/http"
	"github.com/Fliegermarzipan/gallipot/data"
	"html/template"
	"strings"
)

type FrontendData struct {
	//Request *http.Request
	LoggedIn bool
	User     *data.User
}

var (
	indexTemplate *template.Template
)

func LoadTemplates() {
	funcMap := template.FuncMap{
		"Split": strings.Split,
		"Contains": func(s []string, e string) bool {
			for _, a := range s {
				if a == e {
					return true
				}
			}
			return false
		},
		"Combine": func(s ...[]string) []string {
			ret := []string{}
			for i := range s {
				ret = append(ret, s[i])
			}
			return ret
		},
	}

	//indexTemplate = template.Must(template.ParseFiles("index.html", "sidebar.html",
	//	"head.html", "login.html")).Funcs(funcMap)

	indexTemplate = template.Must(template.ParseGlob("*.html")).Funcs(funcMap).Lookup("index.html")
}

func GetTemplates() *template.Template {
	if indexTemplate == nil {
		LoadTemplates()
	}
	return indexTemplate
}
