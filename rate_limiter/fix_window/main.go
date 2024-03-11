package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Limiter struct {
	WindowDuration time.Duration
	CurrentReq     int
	MaxReq         int
	LastReSetTm    int64
	Lock           *sync.Mutex
}

func (l *Limiter) Allow() bool {
	l.Lock.Lock()
	defer l.Lock.Unlock()

	now := time.Now().UnixMilli()
	if now-l.LastReSetTm > int64(l.WindowDuration/time.Millisecond) {
		l.CurrentReq = 0
		l.LastReSetTm = now
	}
	if l.CurrentReq < l.MaxReq {
		l.CurrentReq++
		return true
	}
	return false

}

func NewLimiter(window time.Duration, maxReq int) *Limiter {
	return &Limiter{
		WindowDuration: window,
		MaxReq:         maxReq,
		Lock:           &sync.Mutex{},
		LastReSetTm:    time.Now().UnixMilli(),
	}
}

func main() {
	limiter := NewLimiter(time.Second, 5)

	for i := 0; i < 10; i++ {
		tempId := i
		go RandomReq(limiter, tempId)
	}

	time.Sleep(time.Minute)
}

func RandomReq(l *Limiter, id int) {
	for {
		if l.Allow() {
			fmt.Printf("id:%d, req allow\n", id)
		} else {
			fmt.Printf("id:%d, req reject\n", id)
		}

		sleep := r.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}
