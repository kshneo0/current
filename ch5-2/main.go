package main

import (
	"fmt"
	"time"
)

func factor_f (c chan<-int, done <- chan bool ){
	//transmitter

	i := 2

	for{
		select {
		case <- done:
			fmt.Println("done")
			close(c)
			return
		default:
			c <- i
			time.Sleep(1 * time.Millisecond)
		}
		i++
	}
}

func main(){
	//receiver

	c := make(chan int)
	done := make(chan bool)

	f := 91

	go factor_f(c, done)

	for {
		select {
		case i := <- c:
			if f%i == 0 {
				fmt.Println(i)
				done <- true
				return
			}
		}
	}
}