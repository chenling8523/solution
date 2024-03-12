package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Limiter struct {
	LeakyRate   int64
	Capacity    int
	Token       int
	LastTokenTm int64
	Lock        *sync.Mutex
}

func (l *Limiter) Allow() bool {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	now := time.Now().UnixMilli()

	tmFlow := now - l.LastTokenTm

	tokenGen := tmFlow * l.LeakyRate / int64(1000)
	l.Token += int(tokenGen)
	if l.Token > l.Capacity {
		l.Token = l.Capacity
	}

	if l.Token > 0 {
		l.Token--
		l.LastTokenTm = now
		return true
	}
	return false
}

func NewLimiter(rate int64, capacity int) *Limiter {
	return &Limiter{
		LeakyRate:   rate,
		Capacity:    capacity,
		Token:       0,
		LastTokenTm: time.Now().UnixMilli(),
		Lock:        &sync.Mutex{},
	}
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
