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
	MaxReq         int
	ReqRecords     []int64
	Lock           *sync.Mutex
}

func (l *Limiter) Allow() bool {
	l.Lock.Lock()
	defer l.Lock.Unlock()

	now := time.Now().UnixMilli()

	for {
		if len(l.ReqRecords) == 0 {
			break
		}
		if (now - l.ReqRecords[0]) <= int64(l.WindowDuration/time.Millisecond) {
			break
		}
		l.ReqRecords = l.ReqRecords[1:]
	}

	if len(l.ReqRecords) < l.MaxReq {
		l.ReqRecords = append(l.ReqRecords, now)
		return true
	}
	return false
}

func NewLimiter(window time.Duration, maxReq int) *Limiter {
	return &Limiter{
		WindowDuration: window,
		MaxReq:         maxReq,
		ReqRecords:     make([]int64, 0, maxReq),
		Lock:           &sync.Mutex{},
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
			fmt.Printf("id:%d,  req allow | ", id)
		} else {
			fmt.Printf("id:%d, req reject \n", id)
		}

		sleep := r.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}
