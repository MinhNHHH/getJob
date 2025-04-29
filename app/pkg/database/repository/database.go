package repository

import (
	getJob "github.com/MinhNHHH/get-job/pkg/app"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseRepository interface {
	NewDB() getJob.Database
}
