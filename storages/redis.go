package storages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	Client *redis.Client
}

type RedisStoreInterface interface {
	Set(key string, value interface{}) error
	Get(key string, typeInfo interface{}) (interface{}, error)
	Del(key string) error
	FlushAll() error
}

func NewRedisDB() (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %s", err)
		return nil, err
	}

	redisStore := &RedisStore{Client: rdb}

	log.Println("Connected to Redis")
	return redisStore, nil
}

func (r *RedisStore) Set(key string, value interface{}) error {
	data, _ := json.Marshal(value)
	err := r.Client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) Get(key string, typeInfo interface{}) (interface{}, error) {
	data, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(data), &typeInfo)
	return typeInfo, nil
}

func (r *RedisStore) Del(key string) error {
	data, err := r.Client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if data == 0 {
		return fmt.Errorf("key not found")
	}

	return nil
}

func (r *RedisStore) FlushAll() error {
	err := r.Client.FlushAll(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil
}
