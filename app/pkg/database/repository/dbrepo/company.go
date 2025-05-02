package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MinhNHHH/get-job/pkg/database/data"
)

func (p *DBRepo) AllCompanies() ([]*data.Companies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, url, created_at, updated_at from companies order by name`

	rows, err := p.SqlConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*data.Companies
	for rows.Next() {
		var company data.Companies
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

func (p *DBRepo) InsertCompany(c *data.Companies) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := p.SqlConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `INSERT INTO companies (name, url, created_at, updated_at)
			 VALUES ($1, $2, $3, $4) RETURNING id`

	err = tx.QueryRowContext(ctx, stmt,
		c.Name,
		c.Url,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newID, nil
}

func (p *DBRepo) IsExisted(companyName string) (bool, int) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Use parameterized query to avoid SQL injection
	query := "SELECT id FROM companies WHERE name = $1"

	// Execute the query
	rows, err := p.SqlConn.QueryContext(ctx, query, companyName)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false, 0
	}
	defer rows.Close()

	// If a row exists, that means the company exists in the database
	if rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println("Error scanning id:", err)
			return false, 0
		}
		return true, id
	}

	return false, 0
}
