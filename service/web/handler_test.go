package web

import (
	"github.com/TMDNCM/ictm/persistence/sqlite"
	"github.com/TMDNCM/ictm/data"
	"net/http"
	"testing"
)

func TestWeb(t *testing.T) {

	testdir := t.TempDir()

	p := sqlite.NewPersistor(testdir+"/proto.db")
	t.Log(testdir)
	p.InitDB()
	p.Register(data.LoginData{Username:"test", Password:"123"}, "test@example.com")
	h := NewHandler(p)
	http.Handle("/", h)
	t.Log(http.ListenAndServe(":8080", nil))

}
