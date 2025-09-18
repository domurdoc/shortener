package db

import "fmt"

func NewPGArger() arger {
	return &pgArger{}
}

type pgArger struct {
	pos int
}

func (a *pgArger) next() string {
	a.pos += 1
	return fmt.Sprintf("$%d", a.pos)
}
