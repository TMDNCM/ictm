package sqlite

import (
	"github.com/Fliegermarzipan/gallipot/data"
	"testing"
)

func TestAuth(t *testing.T) {

	p := NewSQLitePersist("proto.db")
	ld := data.LoginData{Username: "test", Password: "456"}
	u := p.Register(ld, "user@local")
	t.Log(u)
	s := p.Authenticate(ld)
	t.Log(s)
}
