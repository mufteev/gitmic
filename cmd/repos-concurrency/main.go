package main

import (
	"fmt"
	"gitmic/internal/repos"
	"log"
	"time"
)

func main() {
	t := time.Now()

	if err := repos.RunConcurrency(true); err != nil {
		log.Printf("repos run: %v", err)
	}

	fmt.Println(time.Since(t))
}
