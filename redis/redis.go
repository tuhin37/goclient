package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedis(host string, port string, password string, db int) RedisClient {
	r := RedisClient{}
	r.client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return r
}

// defer after creating redis client
func (r *RedisClient) Close() {
	r.client.Close()
}

// check if redis is up and running
func (r *RedisClient) Ping() bool {
	_, err := r.client.Ping().Result()
	return err == nil
}

// set a new key value pair in redis
func (r *RedisClient) Set(key string, value string, minutes ...uint16) bool {
	// check variadic input
	mts := 0
	if len(minutes) > 1 {
		fmt.Println("ERROR: i")
		return false
	} else if len(minutes) == 1 {
		mts = int(minutes[0])
	} else {
		mts = 0
	}

	// set minutes to live
	minutesToLive := time.Duration(mts) * time.Minute
	err := r.client.Set(key, value, minutesToLive).Err()

	// return status
	return err == nil
}

// get value from	 key
func (r *RedisClient) Get(key string) (string, bool) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}

// get value and TTL for key
func (r *RedisClient) GetWithTTL(key string) (string, int, bool) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", 0, false
	}

	// get TTL of key
	ttl, err := r.client.TTL(key).Result()
	if err != nil {
		return "", 0, false
	}

	return val, int(ttl.Seconds()), true
}

// delete key
func (r *RedisClient) Unset(key string) error {
	// Use the client to send the DEL command to Redis.
	err := r.client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

// queue operations-------------------------------------------------------------

// push
func (r *RedisClient) Push(queue string, value string) error {
	err := r.client.LPush(queue, value).Err()
	return err
}

// pop
func (r *RedisClient) Pop(queue string) (string, error) {
	result, err := r.client.RPop(queue).Result()
	return result, err
}

// get length of queue
func (r *RedisClient) Qlength(queue string) (int64, error) {
	length, err := r.client.LLen(queue).Result()
	return length, err
}

// --------------------------------- Example usage --------------------------------------

/*
	redis := client.NewRedis("127.0.01", 6379, "", 0)
	redis.Push("boo", "bar")
	redis.Push("boo", "new")
	val, err := redis.Length("boo")
	redis.Push("boo", "goo")
	redis.Pop("boo")
	val, err := redis.Length("boo")

	// get
	val, _ := redis.Get("noo")
	fmt.Println(val)

	// get with ttl
	val, ttl, _ := redis.GetWithTTL("noo")
	fmt.Println(val)
	fmt.Println(ttl)

*/
