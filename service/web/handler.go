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

type WebHandler struct {
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
		if session = p.GetSession(cookie.Value);session!=nil&& session.Valid() {
			b.User = session.User().Get()
			b.LoggedIn = true
		}
	}
	b.Path = strings.Split(r.URL.Path, "/")[1:]
	return b
}

func haveRequiredLogin(b template.BaseRenderer, w http.ResponseWriter) bool {
	if !b.LoggedIn { // 403
		forbidden(b, w)
		return false
	}
	return true
}

func forbidden(b template.BaseRenderer, w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	t := &template.ForbiddenHtml{BaseRenderer: b}
	t.Register(t)
	t.Render(w)
}

func notFound(b template.BaseRenderer, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	t := &template.NotFoundHtml{BaseRenderer: b}
	t.Register(t)
	t.Render(w)
}

func makeServeMux(p persistence.Persistor) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if r.URL.Path == "/" {
			if b.LoggedIn {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther) // redirect to dashboard
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther) // redirect to signup
			}
		} else { // 404
			notFound(b, w)
		}
	})

	m.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if !haveRequiredLogin(b, w) {
			return
		}

		t := &template.DashboardHtml{BaseRenderer: b}
		t.Register(t)
		t.Render(w)
	})

	m.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if !haveRequiredLogin(b, w) {
			return
		}

		friends := make([]data.User, 0, len(p.GetUser(b.User.Username).Friends()))
		for _, v := range p.GetUser(b.User.Username).Friends() {
			friends = append(friends, *(v.Get()))
		}
		t := &template.FriendsHtml{BaseRenderer: b, Friends: friends}
		t.Register(t)
		t.Render(w)
	})

	m.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if !haveRequiredLogin(b, w) {
			return
		}

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
				t := &template.LoginHtml{BaseRenderer: b, LoginAttemptedAs: ld.Username}
				t.Register(t)
				t.Render(w)
			}
		} else { //no login attempted
			t := &template.LoginHtml{BaseRenderer: b}
			t.Register(t)
			t.Render(w)
		}
	})

	m.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if !haveRequiredLogin(b, w) {
			return
		}

		t := &template.ProfileHtml{BaseRenderer: b}
		t.Register(t)
		t.Render(w)
	})

	m.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if r.FormValue("password") != r.FormValue("password-again") {
			t := &template.SignupHtml{BaseRenderer: b, SignupAttempt: "Your passwords do not match"}
			t.Register(t)
			t.Render(w)
		} else if r.FormValue("username") != "" && r.FormValue("email") != "" && r.FormValue("password") != "" {
			if p.GetUser(r.FormValue("username")) != nil { // username already exists
				t := &template.SignupHtml{BaseRenderer: b, SignupAttempt: "Username already taken"}
				t.Register(t)
				t.Render(w)
			} else { // username available
				var ld *data.LoginData
				ld = new(data.LoginData)
				ld.Username = r.FormValue("username")
				ld.Password = r.FormValue("password")
				// persist the registration
				if p.Register(*ld, r.FormValue("email")) == nil {
					t := &template.SignupHtml{BaseRenderer: b, SignupAttempt: "An error occurred while trying to sign up"}
					t.Register(t)
					t.Render(w)
				} else { // registration done, attempt to log in
					sess := p.Authenticate(*ld)
					if sess.Valid() { //successful
						sessionData := sess.Get()
						b.User = sessionData.User
						http.SetCookie(w, &http.Cookie{Name: "token", Value: sessionData.Token,
							Expires: sessionData.Expiry, SameSite: http.SameSiteStrictMode})
						http.Redirect(w, r, "/", http.StatusSeeOther) // redirect to front page
					} else { //unsuccessful
						t := &template.LoginHtml{BaseRenderer: b, LoginAttemptedAs: ld.Username}
						t.Register(t)
						t.Render(w)
					}
				}
			}
		} else { //no signup attempted
			t := &template.SignupHtml{BaseRenderer: b}
			t.Register(t)
			t.Render(w)
		}
	})

	m.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		b := baseRenderer(p, r)
		if !haveRequiredLogin(b, w) {
			return
		}

		if len(b.Path) == 2 {
			user := p.GetUser(b.Path[1]).Get()
			if user != nil {
				t := &template.UserHtml{BaseRenderer: b, Userpage: user}
				t.Register(t)
				t.Render(w)
			} else {
				notFound(b, w)
			}
		} else if len(b.Path) == 1 {
			t := &template.UserHtml{BaseRenderer: b, Userpage: b.User}
			t.Register(t)
			t.Render(w)
		} else {
			notFound(b, w)
		}
	})

	return m
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)

}
