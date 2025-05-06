package repository

import (
	"database/sql"

	"github.com/MinhNHHH/get-job/pkg/database/data"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

type DatabaseRepo interface {
	SQLConnection() *sql.DB
	RedisConnection() *redis.Client
	InsertCompany(c *data.Companies) (int, error)
	IsExisted(companyName string) (bool, int)
	AllCompanies() ([]*data.Companies, error)
	InsertJob(j *data.Jobs) (int, error)
	AllJobs(title, companyName, location string, page, pageSize int) ([]*data.Jobs, int, error)
	IsJobExisted(jobTitle, location string, companyId int) bool
	RedisGet(taskID string) (string, error)
}
