package template

import (
	//"net/http"
	"fmt"
	"github.com/Fliegermarzipan/gallipot/data"
	"github.com/dustin/go-humanize"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
		"contains": func(haystack []string, needle string) bool {
			for _, elem := range haystack {
				if elem == needle {
					return true
				}
			}
			return false
		},
		"combine": func(slices ...[]string) []string {
			combined := []string{}
			for _, elem := range slices {
				combined = append(combined, elem...)
			}
			return combined
		},
		"list": func(elems ...interface{}) []interface{} {
			// This turns all arguments into a slice,
			//  as those cannot be directly created from within templates
			return elems
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
		"prettyTime": func(t time.Time) string {
			return humanize.RelTime(t, time.Now(), "ago", "in the future")
		},
		"getUser": func(username string) *data.User {
			// TODO: return data.User of username if exists, else nil
			return nil
		},
		"getFriends": func() []*data.User {
			// TODO: return slice of friends as data.User objects
			// TODO: remove mock data below
			friend1 := new(data.User)
			friend1.Username = "burgerman420"
			friend1.Displayname = "Bob"

			friend2 := new(data.User)
			friend2.Username = "wonderland69"
			friend2.Displayname = "Alice"

			friend3 := new(data.User)
			friend3.Username = "eavesdr0pper"
			friend3.Displayname = "Eve"

			return []*data.User{friend1, friend2, friend3}
		},
		"getLog": func() []*data.LogEntry {
			return []*data.LogEntry{}
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
			return template.HTML(fmt.Sprintf("<meta http-equiv=refresh content='0; url = %s'>", to))
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
