package dbrepo

import (
	"context"
)

var ctx = context.Background()

func (p *DBRepo) RedisGet(taskID string) (string, error) {
	res, err := p.RedisConn.Get(ctx, taskID).Result()
	return res, err
}
