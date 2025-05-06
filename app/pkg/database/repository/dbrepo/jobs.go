package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MinhNHHH/get-job/pkg/database/data"
)

func (p *DBRepo) AllJobs(title, companyName, location string, page, pageSize int) ([]*data.Jobs, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query for counting total records
	countQuery := `select count(*) from jobs as j 
		inner join companies as c on j.company_id = c.id 
		where 1=1`

	// Base query for fetching records
	query := `select j.id, j.title, j.location, j.description, j.company_id, c.name
		from jobs as j 
		inner join companies as c on j.company_id = c.id 
		where 1=1`

	args := []interface{}{}
	argCount := 1

	if title != "" {
		query += fmt.Sprintf(" AND j.title ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND j.title ILIKE $%d", argCount)
		args = append(args, "%"+title+"%")
		argCount++
	}

	if companyName != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		args = append(args, "%"+companyName+"%")
		argCount++
	}

	if location != "" {
		query += fmt.Sprintf(" AND j.location ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND j.location ILIKE $%d", argCount)
		args = append(args, "%"+location+"%")
		argCount++
	}

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" order by j.created_at desc LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int
	err := p.SqlConn.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := p.SqlConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*data.Jobs
	for rows.Next() {
		var job data.Jobs
		err := rows.Scan(
			&job.Id,
			&job.Title,
			&job.Location,
			&job.Description,
			&job.CompanyId,
			&job.CompanyName,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, 0, err
		}

		jobs = append(jobs, &job)
	}
	return jobs, total, nil
}

func (p *DBRepo) InsertJob(j *data.Jobs) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into jobs (title, location, description, company_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := p.SqlConn.QueryRowContext(ctx, stmt,
		j.Title,
		j.Location,
		j.Description,
		j.CompanyId,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (p *DBRepo) IsJobExisted(jobTitle, location string, companyId int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id FROM jobs WHERE title = $1 and location = $2 and company_id = $3"

	rows, err := p.SqlConn.QueryContext(ctx, query, jobTitle, location, companyId)
	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}

// func (p *DBRepo) UpdateJob(j data.Job) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	stmt := `update jobs set
// 		title = $1,
// 		location = $2,
// 		description = $3,
// 		updated_at = $5
// 		where id = $6
// 	`

// 	_, err := p.DB.SQLConn.ExecContext(ctx, stmt,
// 		j.Title,
// 		j.Location,
// 		j.Description,
// 		time.Now(),
// 		j.Id,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
