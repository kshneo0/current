package main

import (
	"fmt"
	"time"
)


func factor(c chan int, n int) {
	for i:=2; i<= n/2; i++{
		if n%i == 0 {
			c <- i
		}
	}
	close(c)
}

func main() {
	c := make(chan int)

	go factor(c, 100)

	// automatically block/unblock the channel util it closes
	for val := range c {
		fmt.Println(val)
	}
	time.Sleep(10 * time.Millisecond)
}