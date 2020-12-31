package data

import(
	"time"
)
type Session struct{
	User *User
	lastActive time.Time
	expiry time.Time
}


func (s *Session) expired() bool{
	if time.Now().After( s.expiry ){
		return true;
	}else{
		return false;
	}
}
