package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisConn struct {
	Cluster bool
	*redis.Client
	*redis.ClusterClient
}

type RedisConfig struct {
	Host          string        `json:"host"`
	Port          int           `json:"port"`
	IsCluster     bool          `json:"isCluster"`
	PoolSize      int           `json:"poolSize"`
	ClusterConfig ClusterConfig `json:"clusterConfig"`
	MaxIdles      int           `json:"maxIdles"`
	Database      int           `json:"database"`
}
type ClusterConfig struct {
	Nodes        []string `json:"nodes"`
	MaxRedirects int      `json:"maxRedirects"`
	MaxRetries   int      `json:"maxRetries"`
	IsReadOnly   bool     `json:"isReadOnly"`
}

func NewConn(addr string, db int, poolsize int) (*RedisConn, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		PoolSize: poolsize,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisConn{false, client, nil}, nil
}

func NewClusterConn(addrs []string, poolsize, maxRedirects, maxRetries int, readOnly bool) (*RedisConn, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		PoolSize:     poolsize,
		MaxRedirects: maxRedirects,
		MaxRetries:   3,
		ReadOnly:     readOnly,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisConn{true, nil, client}, nil
}

func (rc *RedisConn) HSet(key, field string, value interface{}) (bool, error) {

	if rc.Cluster {
		return rc.ClusterClient.HSet(key, field, value).Result()
	}
	return rc.Client.HSet(key, field, value).Result()
}

func (rc *RedisConn) HGetAll(key string) (map[string]string, error) {
	var res *redis.StringStringMapCmd

	if rc.Cluster {
		res = rc.ClusterClient.HGetAll(key)
	} else {
		res = rc.Client.HGetAll(key)
	}

	return res.Result()
}

func (rc *RedisConn) HMSet(key string, fields map[string]interface{}) error {

	if rc.Cluster {
		return rc.ClusterClient.HMSet(key, fields).Err()
	}
	return rc.Client.HMSet(key, fields).Err()

}

func (rc *RedisConn) HGet(key string, field string) (string, error) {

	if rc.Cluster {
		return rc.ClusterClient.HGet(key, field).Result()
	}
	return rc.Client.HGet(key, field).Result()

}

func (rc *RedisConn) HDel(key string, fields ...string) error {

	if rc.Cluster {
		return rc.ClusterClient.HDel(key, fields...).Err()
	}
	return rc.Client.HDel(key, fields...).Err()

}

func (rc *RedisConn) Incr(key string) error {

	if rc.Cluster {
		return rc.ClusterClient.Incr(key).Err()
	}
	return rc.Client.Incr(key).Err()
}

func (rc *RedisConn) Set(key string, value interface{}, expiration time.Duration) error {
	if rc.Cluster {
		return rc.ClusterClient.Set(key, value, expiration).Err()
	}
	return rc.Client.Set(key, value, expiration).Err()
}

func (rc *RedisConn) Keys(pattern string) ([]string, error) {
	if rc.Cluster {
		return rc.ClusterClient.Keys(pattern).Result()
	}
	return rc.Client.Keys(pattern).Result()
}

func (rc *RedisConn) Get(key string) (string, error) {
	if rc.Cluster {
		return rc.ClusterClient.Get(key).Result()
	}
	return rc.Client.Get(key).Result()
}

func (rc *RedisConn) SAdd(key string, values ...interface{}) error {
	if rc.Cluster {
		return rc.ClusterClient.SAdd(key, values...).Err()
	}
	return rc.Client.SAdd(key, values...).Err()

}

func (rc *RedisConn) SRem(key string, values ...interface{}) error {
	if rc.Cluster {
		return rc.ClusterClient.SRem(key, values...).Err()
	}
	return rc.Client.SRem(key, values...).Err()

}

func (rc *RedisConn) SMembers(key string) ([]string, error) {
	if rc.Cluster {
		return rc.ClusterClient.SMembers(key).Result()
	}
	return rc.Client.SMembers(key).Result()

}

func (rc *RedisConn) SMembersMap(key string) (map[string]struct{}, error) {
	if rc.Cluster {
		return rc.ClusterClient.SMembersMap(key).Result()
	}
	return rc.Client.SMembersMap(key).Result()

}

func (rc *RedisConn) Del(key ...string) error {
	if rc.Cluster {
		return rc.ClusterClient.Del(key...).Err()
	}
	return rc.Client.Del(key...).Err()

}

func (rc *RedisConn) Stop() error {
	if rc.Client != nil {
		return rc.Client.Close()
	}
	if rc.ClusterClient != nil {
		return rc.ClusterClient.Close()
	}
	return nil
}

func (rc *RedisConn) Rename(keyOls, keyNew string) error {
	if rc.Cluster {
		return rc.ClusterClient.Rename(keyOls, keyNew).Err()
	}
	return rc.Client.Rename(keyOls, keyNew).Err()

}

func (rc *RedisConn) HIncrBy(key, field string, incr int64) error {
	if rc.Cluster {
		return rc.ClusterClient.HIncrBy(key, field, incr).Err()
	}
	return rc.Client.HIncrBy(key, field, incr).Err()

}
