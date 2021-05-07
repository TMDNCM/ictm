package template

import (
	"embed"
	"fmt"
	"github.com/TMDNCM/ictm/data"
	"html/template"
	"io"
	"log"
)

var (
	//go:embed *
	tplFiles embed.FS

	templates *template.Template
	pages     = make(map[string]*template.Template)
)

type Renderer interface {
	TemplateName() string
	Render(w io.Writer) error
}

func Render(r Renderer, w io.Writer) error {
	t := pages[r.TemplateName()]
	if t == nil {
		t = template.Must(template.Must(template.Must(GetTemplates().Clone()).
			Parse("{{define \"content\"}} {{template \""+r.TemplateName()+"\" .}} {{end}}")).ParseFS(tplFiles, r.TemplateName()))
	}
	if err := t.ExecuteTemplate(w, "index.html", r); err != nil {
		return err
	}

	return nil
}

type CommonFields struct {
	Renderer
	//Request *http.Request
	LoggedIn bool
	Path     []string
	User     *data.User
}

type BaseMethods struct {
	templateName func(Renderer) string
	render       func(Renderer, io.Writer) error
}

func (b BaseMethods) Dict(elems ...interface{}) map[interface{}]interface{} {
	// This turns all arguments into a map,
	//  where key & val are all arguments as pairs of two,
	//  as those cannot be directly created from within templates
	m := make(map[interface{}]interface{})
	for i := range elems[:len(elems)/2] {
		m[elems[i*2]] = elems[i*2+1]
	}
	return m
}

func (b BaseMethods) List(elems ...interface{}) []interface{} {
	// This turns all arguments into a slice,
	//  as those cannot be directly created from within templates
	return elems
}

var defaultMethods = BaseMethods{
	TemplateName,
	Render,
}

type BaseRenderer struct {
	CommonFields
	BaseMethods
}

func (r *BaseRenderer) TemplateName() string {
	if r.templateName == nil {
		r.templateName = defaultMethods.templateName
	}
	return r.templateName(r.CommonFields.Renderer)
}

func (r *BaseRenderer) Render(w io.Writer) error {
	if r.render == nil {
		r.render = defaultMethods.render
	}
	return r.render(r.CommonFields.Renderer, w)
}

func (r *BaseRenderer) Register(self Renderer) Renderer {
	r.CommonFields.Renderer = self
	return self
}

func TemplateName(r Renderer) string {
	switch (r).(type) {
	case *NotFoundHtml:
		return "404.html"
	case *ForbiddenHtml:
		return "403.html"
	case *LogHtml:
		return "log.html"
	case *FriendsHtml:
		return "friends.html"
	case *LoginHtml:
		return "login.html"
	case *DashboardHtml:
		return "dashboard.html"
	case *ProfileHtml:
		return "profile.html"
	case *SignupHtml:
		return "signup.html"
	case *UserHtml:
		return "user.html"
	case *AboutHtml:
		return "about.html"
	case *NotificationsHtml:
		return "notifications.html"
	case *StockHtml:
		return "stock.html"
	default:
		return fmt.Sprintf("type is %T", r)
	}
}

type NotFoundHtml struct {
	BaseRenderer
}

type ForbiddenHtml struct {
	BaseRenderer
}

type LogHtml struct {
	BaseRenderer
	Entries []data.Dose
}

type FriendsHtml struct {
	BaseRenderer
	Friends []data.User
}

type LoginHtml struct {
	BaseRenderer
	LoginAttemptedAs string
}

type AboutHtml struct {
	BaseRenderer
}

type SignupHtml struct {
	BaseRenderer
	LoginData *data.LoginData
	Email     string
}

type DashboardHtml struct {
	BaseRenderer
}

type ProfileHtml struct {
	BaseRenderer
}

type NotificationsHtml struct {
	BaseRenderer
}

type StockHtml struct {
	BaseRenderer
}

type UserHtml struct {
	BaseRenderer
	Userpage *data.User
}

func LoadTemplates() {

	templates = template.Must(template.New("").ParseFS(tplFiles, "layout/*"))
	templates = template.Must(templates.ParseFS(tplFiles, "assets/*"))

	log.Println(GetTemplates())

}

func GetTemplates() *template.Template {
	if templates == nil {
		LoadTemplates()
	}
	return templates
}
