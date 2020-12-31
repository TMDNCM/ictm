package template

import(
	"net/http"
	"html/template"
	"github.com/Fliegermarzipan/gallipot/data"
)

type FrontendData struct{
	request *http.Request
	loggedIn bool
	user *data.User
}
