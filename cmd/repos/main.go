package main

import (
	"gitmic/internal/repos"
	"log"
)

// Исходная точка приложения, где оно начнёт выполняться

func main() {
	if err := repos.Run(); err != nil {
		log.Printf("repos run: %v", err)
	}
}
