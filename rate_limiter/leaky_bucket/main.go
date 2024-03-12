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
	RemainWater int
	LastLeakTm  int64
	Lock        *sync.Mutex
}

func (l *Limiter) Allow() bool {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	now := time.Now().UnixMilli()

	tmFlow := now - l.LastLeakTm

	fmt.Println("tmFlow:", tmFlow)

	leak := tmFlow * l.LeakyRate / int64(1000)
	fmt.Println("leak:", leak)

	if leak > 0 {
		l.RemainWater -= int(leak)

		if l.RemainWater < 0 {
			l.RemainWater = 0
		}
	}

	l.RemainWater++

	if l.RemainWater > l.Capacity {
		l.RemainWater--
		return false
	}
	l.LastLeakTm = now
	return true
}

func NewLimiter(rate int64, capacity int) *Limiter {
	return &Limiter{
		LeakyRate:   rate,
		Capacity:    capacity,
		RemainWater: 0,
		LastLeakTm:  time.Now().UnixMilli(),
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
