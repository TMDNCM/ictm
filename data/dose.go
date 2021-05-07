package data

import (
	"time"
)

type Dose struct {
	User      *User
	Substance string
	Route     string
	Amount    int
	Unit string
	Taken     time.Time
}
