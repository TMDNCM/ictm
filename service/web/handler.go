package web

import (
	"github.com/TMDNCM/ictm/data"
	"github.com/TMDNCM/ictm/persistence"
	"github.com/TMDNCM/ictm/template"
	"log"
	"net/http"
	"strings"
)

var pageVisibility = map[string]string{
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

type WebHandler struct {
	pageVisibility map[string]string
	persistor      persistence.Persistor
	mux            *http.ServeMux
}

func NewHandler(p persistece.Persistor) *WebHandler {
	h := new(WebHandler)
	h.persistor = p
	return h
}

func baseRenderer(p persistence.Persistor, r *http.Request) template.BaseRenderer {
	var b http.BaseRenderer
	var session persistence.Session
	cookie, err := r.Cookie(token)
	if err == nil { //login cookie
		if session = p.GetSession(cookie.Value); session.Valid() {
			b.User = session.User().Get()
			b.LoggedIn = true
		}
	}
	b.Path = strings.Split(r.URL.Path, "/")[1:]
	return b
}

func makeServeMux(p persistence.Persistor) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Url.Path == "/" {
			b := baseRenderer(p, r)
			template.DashboardHtml{BaseRenderer: b}.Render(w)
		} else {
			http.NotFound(w, r)
		}
	})

	m.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if len(b.Path) == 2 {
			user := p.GetUser(b.Path[1]).Get()
			if user != nil {
				template.UserHtml{BaseRenderer: b, Userpage: user}.Render(w)
			} else {
				http.NotFound(w, r)
			}
		} else if len(b.Path) == 1 {
			template.UserHtml{BaseRenderer: b, Userpage: b.User}.Render(w)
		} else {
			http.NotFound(w, r)
		}
	})

	m.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		entries := p.GetUser(b.User.Username).History()
		if r.FormValue("after") != "" {
			entries = entries.After(time.Unix(r.FormValue("after")))
		}
		if r.FormValue("since") != "" {
			entries = entries.Since(time.Unix(r.FormValue("since")))
		}
		if r.FormValue("substance") != "" {
			entries = entries.Substance(r.FormValue("substance"))
		}
		if count, err := strconv.ParseUint(r.FormValue("count")); err != nil {
			count = 100
		}
		entries = entries.LastX(count)

		template.LogHtml{BaseRenderer: b, Entries: entries.Get()}.Render(w)
	})

	m.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		friends := make([]data.User, 0, len(p.GetUser(b.User.Username).Friends()))
		for _, v := range p.GetUser(b.User.Username).Friends() {
			friends = append(friends, v.Get())
		}
		template.FriendsHtml{BaseRenderer: b, Friends: friends}.Render(w)
	})

	m.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		var ld *data.LoginData
		if r.FormValue("username") != "" && r.FormValue("password") != "" {
			ld = new(data.LoginData)
			ld.Username = r.FormValue("username")
			ld.Password = r.FormValue("password")
			sess := Authenticate(ld)
			if sess.Valid() { //successful login
				sessionData := sess.Get()
				b.User = sessionData.User
				http.SetCookie(w, &http.Cookie{Name: "token", Value: sessionData.Token,
					Expires: sessionData.Expiry, SameSite:http.SameSiteStrictMode})
				template.LoginHtml{BaseRenderer: b, LoginAttempted: true,
					LoginSuccessful: true, LoginData: ld}.Render(w)
			} else { //unsuccessful login
				template.LoginHtml{BaseRenderer: b, LoginAttempted: true, LoginData: ld}.Render(w)
			}
		} else { //no login attempted
			template.LoginHtml{BaseRenderer: b}.Render(w)
		}
	})

	return m
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method)
	template.NewRenderer(r, h.persistor).Render(w)
}
