package redis

//func Connect(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
//	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
//	client := redis.NewClient(&redis.Options{
//		Addr:        addr,
//		Password:    cfg.Password,
//		DB:          cfg.Database,
//		PoolSize:    cfg.PoolSize,
//		PoolTimeout: time.Duration(cfg.Timeout) * time.Second,
//	})
//
//	_, err := client.Ping(ctx).Result()
//
//	if err != nil {
//		return nil, err
//	}
//	return client, nil
//}0
