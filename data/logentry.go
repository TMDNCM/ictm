package data

import(
	"time"
)

type LogEntry struct{
	User *User
	Substance *Substance
	Route *Route
	Amount int
	Created time.Time
}
