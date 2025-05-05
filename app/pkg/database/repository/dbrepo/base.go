package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

const dbTimeout = time.Second * 3

var ctx = context.Background()

type DBRepo struct {
	SqlConn   *sql.DB
	RedisConn *redis.Client
}

func (p *DBRepo) SQLConnection() *sql.DB {
	return p.SqlConn
}

func (r *DBRepo) RedisConnection() *redis.Client {
	return r.RedisConn
}
