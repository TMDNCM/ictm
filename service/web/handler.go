package web

import (
	"github.com/TMDNCM/ictm/data"
	"github.com/TMDNCM/ictm/persistence"
	"github.com/TMDNCM/ictm/template"
	//"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func NewHandler(p persistence.Persistor) *WebHandler {
	h := new(WebHandler)
	h.persistor = p
	h.mux = makeServeMux(p)
	return h
}

func baseRenderer(p persistence.Persistor, r *http.Request) template.BaseRenderer {
	var b template.BaseRenderer
	var session persistence.Session
	cookie, err := r.Cookie("token")
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
		if r.URL.Path == "/" {
			b := baseRenderer(p, r)
			t := &template.DashboardHtml{BaseRenderer: b}
			t.Register(t)
			t.Render(w)
		} else {
			http.NotFound(w, r)
		}
	})

	m.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if len(b.Path) == 2 {
			user := p.GetUser(b.Path[1]).Get()
			if user != nil {
				t := &template.UserHtml{BaseRenderer: b, Userpage: user}
				t.Register(t)
				t.Render(w)
			} else {
				http.NotFound(w, r)
			}
		} else if len(b.Path) == 1 {
			t := &template.UserHtml{BaseRenderer: b, Userpage: b.User}
			t.Register(t)
			t.Render(w)
		} else {
			http.NotFound(w, r)
		}
	})

	m.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		entries := p.GetUser(b.User.Username).History()
		if r.FormValue("after") != "" {
			if timestamp, err := strconv.ParseInt(r.FormValue("after"), 10, 64); err == nil {
				entries = entries.After(time.Unix(timestamp, 0))
			}
		}
		if r.FormValue("since") != "" {
			if timestamp, err := strconv.ParseInt(r.FormValue("before"), 0, 64); err == nil {
				entries = entries.Before(time.Unix(timestamp, 0))
			}
		}
		if r.FormValue("substance") != "" {
			entries = entries.OfSubstance(r.FormValue("substance"))
		}
		if count, err := strconv.ParseUint(r.FormValue("count"), 10, 64); err == nil {
			entries = entries.LastX(count)
		} else {
			entries = entries.LastX(100)
		}
		entrylist := entries.Get()
		doses := make([]data.Dose, 0, len(entrylist))
		for _, v := range entrylist {
			doses = append(doses, *(v.Get()))
		}
		t := &template.LogHtml{BaseRenderer: b, Entries: doses}
		t.Register(t)
		t.Render(w)
	})

	m.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		friends := make([]data.User, 0, len(p.GetUser(b.User.Username).Friends()))
		for _, v := range p.GetUser(b.User.Username).Friends() {
			friends = append(friends, *(v.Get()))
		}
		t := &template.FriendsHtml{BaseRenderer: b, Friends: friends}
		t.Register(t)
		t.Render(w)
	})

	m.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		var ld *data.LoginData
		if r.FormValue("username") != "" && r.FormValue("password") != "" {
			ld = new(data.LoginData)
			ld.Username = r.FormValue("username")
			ld.Password = r.FormValue("password")
			sess := p.Authenticate(*ld)
			if sess.Valid() { //successful login
				sessionData := sess.Get()
				b.User = sessionData.User
				http.SetCookie(w, &http.Cookie{Name: "token", Value: sessionData.Token,
					Expires: sessionData.Expiry, SameSite: http.SameSiteStrictMode})
				http.Redirect(w, r, "/", http.StatusSeeOther) // redirect to front page
			} else { //unsuccessful login
				t := &template.LoginHtml{BaseRenderer: b, LoginAttempted: true}
				t.Register(t)
				t.Render(w)
			}
		} else { //no login attempted
			t := &template.LoginHtml{BaseRenderer: b}
			t.Register(t)
			t.Render(w)
		}
	})

	return m
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)

}
