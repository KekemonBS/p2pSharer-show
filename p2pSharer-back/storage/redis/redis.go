package redis

import (
	"context"
	"encoding/json"

	"github.com/KekemonBS/p2pSharer-back/tracker/models"
	"github.com/redis/go-redis/v9"
)

type CacheImpl struct {
	Client *redis.Client
}

func New(c *redis.Client) *CacheImpl {
	return &CacheImpl{
		Client: c,
	}
}

// SaveUser stores a user in Redis with the username as the key and password hash as the value
func (c CacheImpl) SaveUser(ctx context.Context, u models.User) error {
	key := "user:" + u.Name
	value, err := json.Marshal(u.PassHash)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, value, 0).Err() // No expiration
}

// ReadUser retrieves a user from Redis by username
func (c CacheImpl) ReadUser(ctx context.Context, username string) (models.User, error) {
	key := "user:" + username
	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return models.User{}, err
		}
		return models.User{}, err
	}

	var user models.User
	user.Name = username // Set the username explicitly
	err = json.Unmarshal([]byte(value), &user.PassHash)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// DeleteUser removes a user from Redis
func (c CacheImpl) DeleteUser(ctx context.Context, u models.User) error {
	key := "user:" + u.Name
	return c.Client.Del(ctx, key).Err()
}

// Save stores a folder for a user in Redis with a unique key
func (c CacheImpl) Save(ctx context.Context, u models.User, f models.Folder) error {
	uniqueKey := "folder:" + u.Name + ":" + f.Name
	value, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, uniqueKey, value, 0).Err() // No expiration
}

// Read retrieves all folders associated with a given user
func (c CacheImpl) Read(ctx context.Context, u models.User) ([]models.Folder, error) {
	keyPattern := "folder:" + u.Name + ":*" // Pattern to match keys for the user's folders

	var folders []models.Folder
	iter := c.Client.Scan(ctx, 0, keyPattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		value, err := c.Client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var folder models.Folder
		err = json.Unmarshal([]byte(value), &folder)
		if err != nil {
			return nil, err
		}

		folders = append(folders, folder)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

// Delete removes a specific folder for a user from Redis
func (c CacheImpl) Delete(ctx context.Context, u models.User, folderName string) error {
	key := "folder:" + u.Name + ":" + folderName
	return c.Client.Del(ctx, key).Err()
}
