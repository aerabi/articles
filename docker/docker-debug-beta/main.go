package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello, Wowlrd!")
		time.Sleep(5 * time.Second)
	}
}