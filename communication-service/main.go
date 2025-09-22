package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Start")

	log.Println("ok")

	for {
		time.Sleep(time.Hour)
	}
}
