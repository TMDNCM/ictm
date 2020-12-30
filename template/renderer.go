package template


import(
	"html/template"
	"github.com//Fliegermarzipan/gallipot/data"
)



type FrontendData struct{
	loggedIn bool
	user *data.User
}
