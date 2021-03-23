package persistence

import (
	"github.com/TMDNCM/ictm/data"
	"time"
)

type Backend interface {
	Authenticate(ld data.LoginData) *data.Session
	Register(ld data.LoginData, email string) data.User
}

type Persistor interface {
	Authenticate(ld data.LoginData) Session
	Register(ld data.LoginData, email string) User
	getSession(token string) Session
}

type Session interface {
	User() User
	Valid() bool
	Invalidate() Session
	Save() error
}

type User interface {
	Username() string
	SetUsername(username string) User
	Email() string
	SetEmail(email string) User
	Displayname() string
	SetDisplayname(displayname string) User
	Save() error
	History() Doses
	Log(substance string, route string, dose int, unit string, time time.Time)
}

type Doses interface {
	Since(t time.Time) Doses
	After(t time.Time) Doses
	Between(x, y time.Time) Doses
	OfSubstance(substance string) Doses
	LastX(x int) Doses
	Get() []Dose
}

type Dose interface {
	When() time.Time
	SetWhen(t time.Time) Dose
	Amount() int
	SetAmount(amount int) Dose
	Substance() string
	setSubstance(substance string) Dose
	Route() string
	SetRoute(route string) Dose
	Save() error
}
