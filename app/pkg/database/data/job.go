package data

import "time"

type Job struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	Location    string    `db:"location"`
	Description string    `db:"description"`
	CompanyId   int       `db:"company_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
