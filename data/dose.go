package data

import (
	"time"
)

type Dose struct {
	User      *User
	Substance *Substance
	Route     *Route
	Amount    int
	Created   time.Time
}
