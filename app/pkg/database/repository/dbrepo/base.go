package dbrepo

import (
	"time"

	getJob "github.com/MinhNHHH/get-job/pkg/app"
)

type PostgresDBRepo struct {
	DB *getJob.Database
}

const dbTimeout = time.Second * 3
