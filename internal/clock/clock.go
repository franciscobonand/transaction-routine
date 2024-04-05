//go:generate mockgen -destination=./../../tests/mocks/mock_clock.go -package=mocks -source=clock.go
package clock

import "time"

type Clock interface {
	Now() time.Time
}

type clock struct {
	loc *time.Location
}

func New() Clock {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.UTC
	}
	return &clock{
		loc: loc,
	}
}

func (c *clock) Now() time.Time {
	return time.Now().In(c.loc)
}
