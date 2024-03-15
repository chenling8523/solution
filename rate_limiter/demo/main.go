package main

import (
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	limiter := rate.NewLimiter(1, 50)

	for i := 20; i > 0; i-- {
		go func(i int) {
			time.Sleep(time.Second * time.Duration(rand.Intn(10)))
			if limiter.AllowN(time.Now(), 7) {
				println("allow ", i)
			} else {
				println("limit ", i)
			}
		}(i)
	}
	time.Sleep(time.Second * 15)
}
