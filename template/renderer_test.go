package template

import (
	"github.com/Fliegermarzipan/gallipot/data"
	"os"
	"testing"
)

func TestRenderer(t *testing.T) {
	testData := FrontendData{
		true,
		new(data.User),
	}
	testData.User.Username = "test"
	testData.User.Displayname = "test2"
	testData.User.Email = "test@test.com"

	tpl := GetTemplates()
	for _, v := range tpl.Templates() {
		t.Logf("%#v - %#v", *v, v.Name())
	}

	execErr := tpl.Execute(os.Stderr, testData)
	if execErr != nil {
		t.Fatal(execErr)
	}
}
