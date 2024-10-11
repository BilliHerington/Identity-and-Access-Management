package redisUsersHandlers

func (repo *RedisUsersRepository) SavePassCode(passCode, userID string) error {
	err := repo.RDB.HSet(ctx, "user:"+userID, passCode, 0).Err()
	return err
}
