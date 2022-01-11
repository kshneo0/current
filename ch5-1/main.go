package main

import (
	"fmt"
)

func hello(c chan string) {
	fmt.Println("Hello, " + <- c + "!")
	close(c)
}

func squares(c chan int) {
	for i :=0; i< 10; i++ {
		c <-i * i
	}

	close(c)
}

func main() {
	// c := make(chan string)

	// go hello(c)

	// c <- "William"

	c := make(chan int)

	go squares(c)

	for {
		val, ok := <- c

		if ok == false {
			fmt.Println(val, ok, "not ok")
			break
		} else {
			fmt.Println(val, ok)
		}
	}

	// time.Sleep(100 * time.Millisecond)
}