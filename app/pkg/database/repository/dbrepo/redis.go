package dbrepo

func (p *DBRepo) RedisGet(key string) (string, error) {
	res, err := p.RedisConn.Get(ctx, key).Result()
	return res, err
}
