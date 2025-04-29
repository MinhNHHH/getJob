package data

import "time"

type Company struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Url       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
