package template

import (
	//"net/http"
	"github.com/Fliegermarzipan/gallipot/data"
	"html/template"
	"strings"
	"log"
	"path"
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

func templateFromFile(funcMap template.FuncMap, s string, d interface{}) string {
		var buf strings.Builder
		err := template.Must(template.New(path.Base(s)).Funcs(funcMap).ParseFiles(s)).Execute(&buf, d)
		if err != nil {
			log.Print(err)
		}
		return buf.String()
}

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
		return template.HTML(templateFromFile(funcMap, s, d))
	}
	funcMap["includeCSS"] = func(s string, d interface{}) template.CSS {
		return template.CSS(templateFromFile(funcMap, s, d))
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
