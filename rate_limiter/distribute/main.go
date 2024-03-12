package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const script = `
	local key = KEYS[1]
	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])

	local tm = redis.call('TIME')
	local now = tonumber(tm[1]*1000 + tm[2]/1000)
	local token, lastTokenTm = unpack(redis.call('HMGET', key, 'token', 'lastTokenTm'))
	token = tonumber(token)
	lastTokenTm = tonumber(lastTokenTm)

	if not token or not lastTokenTm then
		token = capacity
		lastTokenTm = now
	else
		local tmFlow = now - lastTokenTm
		local tokenGen = tmFlow * rate / 1000
		token = math.min(token + tokenGen, capacity)
	end

	if token > 0 then
		token = token - 1
		lastTokenTm = now
		redis.call('HMSET', key, 'token', token, 'lastTokenTm', lastTokenTm)
		return 1
	else 
		return 0
	end
`

type Limiter struct {
	Ctx         context.Context
	RedisClient *redis.Client
	LeakyRate   int64
	Capacity    int
}

func (l *Limiter) Allow() bool {
	result, err := l.RedisClient.Eval(l.Ctx, script, []string{"rate_limiter"}, l.LeakyRate, l.Capacity).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return fmt.Sprintf("%v", result) == "1"
}

func NewLimiter(rate int64, capacity int) *Limiter {
	limiter := &Limiter{
		Ctx:       context.Background(),
		LeakyRate: rate,
		Capacity:  capacity,
	}

	limiter.RedisClient = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)
	return limiter
}

func main() {

	limiter := NewLimiter(5, 5)

	for i := 0; i < 10; i++ {
		tempId := i
		go RandomReq(limiter, tempId)
	}

	time.Sleep(time.Minute)
}

func RandomReq(l *Limiter, id int) {
	for {
		sleep := r.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		if l.Allow() {
			fmt.Printf("id:%d,  req allow | ", id)
		} else {
			fmt.Printf("id:%d, req reject \n", id)
		}

	}
}
