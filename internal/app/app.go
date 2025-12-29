package app

import (
	"enactus/internal/database"
	"fmt"
	"log"
)

func Run() {
	pool, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	defer pool.Close()

	fmt.Println("Success", pool)
}
