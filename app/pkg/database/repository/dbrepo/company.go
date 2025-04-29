package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MinhNHHH/get-job/pkg/database/data"
)

func (p *PostgresDBRepo) AllCompanies() ([]*data.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, url from company order by name`

	rows, err := p.DB.SQLConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*data.Company
	for rows.Next() {
		var company data.Company
		err := rows.Scan(
			&company.Id,
			&company.Name,
			&company.Url,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		companies = append(companies, &company)
	}
	return companies, nil
}

func (p *PostgresDBRepo) InsertCompany(c *data.Company) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into company (name, url, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := p.DB.SQLConn.QueryRowContext(ctx, stmt,
		c.Name,
		c.Url,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (p *PostgresDBRepo) IsExisted(companyName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Use parameterized query to avoid SQL injection
	query := "SELECT id FROM company WHERE name = $1"

	// Execute the query
	rows, err := p.DB.SQLConn.QueryContext(ctx, query, companyName)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false
	}
	defer rows.Close()

	// If a row exists, that means the company exists in the database
	return rows.Next()
}
