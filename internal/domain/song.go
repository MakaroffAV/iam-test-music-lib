package domain

import "time"

type Song struct {
	ID       string
	Youtube  string
	Released time.Time
	Album    string
	Title    string
	Artist   string
}
