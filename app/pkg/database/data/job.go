package data

import "time"

type Jobs struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	Location    string    `db:"location"`
	CompanyName string    `db:"company_name"`
	CompanyId   int       `db:"company_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
