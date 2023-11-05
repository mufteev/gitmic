package main

import (
	"gitmic/internal/repos"
	"log"
)

// Исходная точка приложения, где оно начнёт выполняться

func main() {
	// if err := repos.RunSimple(); err != nil {
	// 	log.Printf("repos run: %v", err)
	// }
	if err := repos.RunConcurrency(true); err != nil {
		log.Printf("repos run: %v", err)
	}
}
