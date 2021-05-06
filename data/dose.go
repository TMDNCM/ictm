package data

import (
	"time"
)

type Dose struct {
	User      *User
	Substance string
	Route     string
	Amount    int
	Taken     time.Time
}
