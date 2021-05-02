package persistence

import (
	"github.com/TMDNCM/ictm/data"
	"time"
)


type Persistor interface {
	Authenticate(ld data.LoginData) Session
	Register(ld data.LoginData, email string) User
	GetSession(token string) Session
	GetUser(username string) User
}

type Session interface {
	Get() *data.Session
	User() User
	Valid() bool
	Invalidate()
}

type User interface {
	Get() *data.User
	SetUsername(username string) User
	SetEmail(email string) User
	SetDisplayname(displayname string) User
	Friends() []User
	History() Doses
	Log(substance string, route string, dose int, unit string, time time.Time)
}

type Doses interface {
	Since(t time.Time) Doses
	After(t time.Time) Doses
	OfSubstance(substance string) Doses
	LastX(x uint64) Doses
	Get() []Dose
}

type Dose interface {
	Get() *data.Dose
	SetWhen(t time.Time) Dose
	SetAmount(amount int) Dose
	SetSubstance(substance string) Dose
	SetRoute(route string) Dose
}
