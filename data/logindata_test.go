package data


import(
	"testing"
)


func TestHash(t *testing.T){
	ld := &LoginData{Username:"testuser", Password:"secret"}
	t.Log(ld, ld.Hash([]byte("adslfjksdfajasdjhgljhg")))
}
