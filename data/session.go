package data

import(
	"time"
)
type Session struct{
	User *User
	Token string
	LastActive time.Time
	Expiry time.Time
}


func (s *Session) Expired() bool{
	if time.Now().After( s.Expiry ){
		return true;
	}else{
		return false;
	}
}
