package template

import (
	//"net/http"
	"github.com/Fliegermarzipan/gallipot/data"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type UserAlert struct {
	Title   string
	Message string
}

type FrontendData struct {
	//Request *http.Request
	LoggedIn bool
	Path     []string
	Page     string
	User     *data.User
	Alert    *UserAlert
}

var (
	indexTemplate *template.Template
)

func getTemplateDir() string {
	exec := os.Args[0]
	execDir := filepath.Dir(exec)
	templateDir := filepath.Join(execDir, "template")
	log.Println(templateDir)
	return templateDir
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
		"lower": strings.ToLower,
		"title": strings.Title,
		"exists": func(s string) bool {
			if _, err := os.Stat(filepath.Join(getTemplateDir(), s)); err == nil {
				return true
			} else if os.IsNotExist(err) {
				return false
			}
			log.Print("file exists fucked up")
			return false
		},
	}
	funcMap["include"] = func(s string, d interface{}) interface{} {
		ext := filepath.Ext(s)

		var buf strings.Builder
		err := template.Must(
			template.New(path.Base(s)).
				Funcs(funcMap).
				ParseFiles(filepath.Join(getTemplateDir(), s))).
			Execute(&buf, d)
		if err != nil {
			log.Print(err)
		}
		rendered := buf.String()

		if ext == ".html" {
			return template.HTML(rendered)
		} else if ext == ".css" {
			return template.CSS(rendered)
		}
		return rendered
	}

	indexTemplate = template.Must(template.New("").Funcs(funcMap).
		ParseGlob(filepath.Join(getTemplateDir(), "*.html"))).Lookup("index.html")
	log.Println(indexTemplate)
}

func GetTemplates() *template.Template {
	if indexTemplate == nil {
		LoadTemplates()
	}
	return indexTemplate
}
