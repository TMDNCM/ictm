package sqlite

import (
	"github.com/TMDNCM/ictm/data"
	"testing"
	"os"
	"time"
)

func TestAuth(t *testing.T) {
	os.Remove("/tmp/proto.db")
	p := NewPersistor("/tmp/proto.db")
	p.InitDB()
	ld := data.LoginData{Username: "test", Password: "456"}
	u := p.Register(ld, "user@local")
	ld2 := data.LoginData{Username:"friend", Password:"321"}
	u2 := p.Register(ld2, "email@other.com")
	u2.AddFriend(u)
	u.Log("aspirin", "oral", 100, "mg", time.Now().Add(-2*time.Hour))
	u.Log("paracetamol", "oral", 500, "mg", time.Now().Add(-4*time.Hour))
	u2.Log("ibuprofen", "oral", 400, "mg", time.Now())
	t.Logf("%#+v\n",u.Get())
	t.Logf("%#+v\n",u.History().After(time.Now().Add(-3*time.Hour)))
	t.Logf("%+#v\n",u2.Friends())
	s := p.Authenticate(ld)
	t.Logf("%#+v\n", s.Get())
}
