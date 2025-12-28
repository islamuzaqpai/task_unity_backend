package app

import (
	"enactus/internal/database"
	"enactus/internal/repository"
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

	attRepo := repository.AttendanceRepository{Pool: pool}

	atts, err := attRepo.GetAllAttendance()
	if err != nil {
		log.Fatalf("fail: %v", err)
	}

	fmt.Println(atts)
}
