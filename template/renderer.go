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
	"fmt"
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
		"exists": func(filename string) bool {
			if _, err := os.Stat(filepath.Join(getTemplateDir(), filename)); err == nil {
				return true
			} else if os.IsNotExist(err) {
				return false
			}
			log.Print("file exists fucked up")
			return false
		},
		"getUser": func(username string) *data.User {
			// TODO: return data.User of username if exists, else nil
			return nil
		},
		"getFriends": func() []*data.User {
			// TODO: return slice of friends as data.User objects
			return []*data.User{}
		},
		"getNotificationCount": func() int {
			// TODO: return amount of unread notifications
			return 420
		},
		"getFriendCount": func() int {
			// TODO: return amount of friends
			return 3
		},
		"getLogUnreadCount": func() int {
			// TODO: return amount of unread log entries
			return 69
		},
		"redirect": func(to string) template.HTML {
			return template.HTML(fmt.Sprintf("<meta http-equiv=refresh content='0; url = %s'", to))
		},
	}
	funcMap["include"] = func(filename string, fd FrontendData) interface{} {
		ext := filepath.Ext(filename)

		var buf strings.Builder
		err := template.Must(
			template.New(path.Base(filename)).
				Funcs(funcMap).
				ParseFiles(filepath.Join(getTemplateDir(), filename))).
			Execute(&buf, fd)
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
