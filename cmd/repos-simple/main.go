package main

import (
	"fmt"
	"gitmic/internal/repos"
	"log"
	"time"
)

// Исходная точка приложения, где оно начнёт выполняться

func main() {
	t := time.Now()

	if err := repos.RunSimple(true); err != nil {
		log.Printf("repos run: %v", err)
	}

	fmt.Println(time.Since(t))
}
