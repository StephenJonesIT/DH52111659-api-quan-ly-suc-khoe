package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)
type RedisStore interface {
	StoreOTP(ctx context.Context, email, otp string) error
	VerifyOTP(ctx context.Context, email, otp string) (bool, error)
}


type RedisStoreImpl struct {
	client *redis.Client
}

func NewRedisStore() (*RedisStoreImpl, error) {
	if config.AppConfig.RedisHost == ""{
		return nil, fmt.Errorf("chưa cấu hình Redis host") // or return an error if you prefer
	}

	// Initialize Redis client with configuration from AppConfig
	redisOptions := &redis.Options{
		Addr: config.AppConfig.RedisHost,
		Password: config.AppConfig.RedisPass, // No password set
		DB: 0, // use default DB
		DialTimeout:  5 * time.Second,  // Thêm timeout
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
	}
	

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	client := redis.NewClient(redisOptions)
	if err := client.Ping(ctx).Err(); err != nil {
		 if closeErr := client.Close(); closeErr != nil {
            log.Printf("Lỗi khi đóng kết nối Redis: %v", closeErr)
        }
        return nil, fmt.Errorf("không thể kết nối tới Redis tại %s: %v", config.AppConfig.RedisHost, err)
	}

	return &RedisStoreImpl{client: client}, nil
}

func (r *RedisStoreImpl) StoreOTP(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.client.SetEX(ctx, key, otp, 10*time.Minute).Err()
}

func (r *RedisStoreImpl) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	key := fmt.Sprintf("otp:%s", email)
	storeOTP, err := r.client.Get(ctx, key).Result()
	
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return storeOTP == otp, nil
}
