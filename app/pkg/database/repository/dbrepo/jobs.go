package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/MinhNHHH/get-job/pkg/database/data"
)

func (p *DBRepo) AllJobs() ([]*data.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, location, description from jobs order by created_at`

	rows, err := p.SqlConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*data.Job
	for rows.Next() {
		var job data.Job
		err := rows.Scan(
			&job.Id,
			&job.Title,
			&job.Location,
			&job.Description,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (p *DBRepo) InsertJob(j *data.Job) (int, error) {
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
