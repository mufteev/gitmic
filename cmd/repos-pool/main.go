package main

import (
	"gitmic/internal/repos"
	"log"
)

func main() {
	if err := repos.RunPool(); err != nil {
		log.Panicf("run pool: %v", err)
	}
}
