package data


import(
	"crypto/sha512"
	"io"
)

type LoginData struct{
	Username string
	Password string
}


func (ld *LoginData) Hash(salt []byte)(hash []byte){
	h := sha512.New()
	io.WriteString(h, ld.Password)
	h.Write(salt)
	hash = h.Sum(hash)
	return
}
