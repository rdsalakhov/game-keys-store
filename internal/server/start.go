package server

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rdsalakhov/game-keys-store/internal/store/mysqlStore"
	"net/http"
	"os"
	"strconv"
)

func Start(config *Config) error {
	db, err := newDb(config.DbConnection)
	if err != nil {
		return err
	}
	defer db.Close()
	store := mysqlStore.New(db)
	redis, err := newRedis(config.RedisConnection, config.RedisPassword)
	if err != nil {
		return err
	}
	writeToEnv(config)

	server := NewServer(store, redis)
	return http.ListenAndServe(config.Port, server)
}

func newDb(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newRedis(redisConnection string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConnection, //redis port
		Password: password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func writeToEnv(config *Config) {
	os.Setenv("ACCESS_SECRET", config.AccessSecret)
	os.Setenv("REFRESH_SECRET", config.RefreshSecret)
	os.Setenv("PLATFORM_FEE_SHARE", strconv.FormatFloat(config.PlatformFeeShare, 'E', -1, 64))
	os.Setenv("PLATFORM_ACCOUNT", config.PlatformAccount)
	os.Setenv("PLATFORM_EMAIL", config.PlatformEmail)
	os.Setenv("PLATFORM_EMAIL_PASSWORD", config.PlatformEmailPassword)
}
