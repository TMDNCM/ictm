package web

import (
	"github.com/TMDNCM/ictm/persistence/sqlite"
	"net/http"
	"os"
	"testing"
)

func TestWeb(t *testing.T) {

	os.Remove("/tmp/proto.db")

	p := sqlite.NewPersistor("/tmp/proto.db")
	p.InitDB()
	h := NewHandler(p)
	http.Handle("/", h)
	t.Log(http.ListenAndServe(":8080", nil))

}
