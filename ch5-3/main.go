package main

import (
	"fmt"
	"sync"
)


func squares(c chan int, wg *sync.WaitGroup){
	num := <-c

	c <- num * num

	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	

	c := make(chan int, 3)

	for i:=0; i<3;i ++ {
		wg.Add(3)
		go squares(c, &wg)
		go squares(c, &wg)
		go squares(c, &wg)

		c <- 2
		c <- 3
		c <- 5

		wg.Wait()

		fmt.Println(<- c)
		fmt.Println(<- c)
		fmt.Println(<- c)
	}
}