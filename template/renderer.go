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
	Page string
	User     *data.User
}

var (
	indexTemplate *template.Template
)

func LoadTemplates() {
	funcMap := template.FuncMap{
		"split": strings.Split,
		"contains": func(s []string, e string) bool {
			for _, a := range s {
				if a == e {
					return true
				}
			}
			return false
		},
		"combine": func(s ...[]string) []string {
			ret := []string{}
			for i := range s {
				ret = append(ret, s[i]...)
			}
			return ret
		},
		"list": func(s ...interface{}) []interface{} {
			return s
		},
		"title": func(s string) string {
			return strings.Title(strings.ToLower(s))
		},
	}
	funcMap["include"] = func(s string, d interface{}) template.HTML {
		var buf strings.Builder
		template.Must(template.New(s).Funcs(funcMap).ParseFiles(s)).Execute(&buf, d)
		return template.HTML(buf.String())
	}

	//indexTemplate = template.Must(template.ParseFiles("index.html", "sidebar.html",
	//	"head.html", "login.html")).Funcs(funcMap)

	indexTemplate = template.Must(template.New("").Funcs(funcMap).ParseGlob("*.html")).Lookup("index.html")
}

func GetTemplates() *template.Template {
	if indexTemplate == nil {
		LoadTemplates()
	}
	return indexTemplate
}
