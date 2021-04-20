package log

import (
	"github.com/go-redis/redis"
)

func NewRedisWriter(key string, cli *redis.Client) *redisWriter {
	return &redisWriter{
		cli: cli, listKey: key,
	}
}

// 为 logger 提供写入 redis 队列的 io 接口
type redisWriter struct {
	cli     *redis.Client
	listKey string
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}
