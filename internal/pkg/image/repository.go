package image

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v7"
)

type StatusRepository struct {
	RedisClient *redis.Client
}

func (r *StatusRepository) generateKey(uuid string) string {
	return fmt.Sprintf("image:uuid:%s", uuid)
}

func (r *StatusRepository) Set(uuid string, status Status) error {
	value, err := json.Marshal(&status)
	if err != nil {
		return err
	}

	key := r.generateKey(uuid)
	if err := r.RedisClient.Set(key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r *StatusRepository) Get(uuid string) (Status, error) {
	key := r.generateKey(uuid)
	value, err := r.RedisClient.Get(key).Result()
	if err != nil {
		return Status{}, err
	}

	var status Status
	if err := json.Unmarshal([]byte(value), &status); err != nil {
		return Status{}, err
	}

	return status, nil
}
