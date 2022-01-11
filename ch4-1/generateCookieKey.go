// package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gorilla/securecookie"
)

func main() {
	var key = securecookie.GenerateRandomKey(32)
	err := ioutil.WriteFile("key.txt",key, 0644)
	if err != nil {
		fmt.Println(err)
	}
}