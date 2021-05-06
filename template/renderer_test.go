package template

import (
	"bytes"
	"github.com/TMDNCM/ictm/data"
	_ "github.com/TMDNCM/ictm/persistence/sqlite"
	"testing"
	_ "time"
)

func TestRenderer(t *testing.T) {
	testUser := new(data.User)
	testUser.Username = "test"
	testUser.Displayname = "test2"

	testUser.Email = "test@test.com"

	renderer := new(UserHtml)
	renderer.Register(renderer)
	renderer.Userpage = testUser
	renderer.User = testUser

	buf := new(bytes.Buffer)
	err := renderer.Render(buf)
	if err != nil {
		t.Log("render borken")
		t.Log(err)
	}
	t.Log(buf)
}
