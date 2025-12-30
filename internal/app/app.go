package app

import (
	"enactus/internal/database"
	"enactus/internal/models"
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

	commRepo := repository.CommentRepository{Pool: pool}

	newComm := models.TaskComment{
		Comment: "updated",
		TaskId:  2,
		UserId:  2,
	}

	updated, err := commRepo.UpdateComment(3, newComm)
	if err != nil {
		log.Fatalf("fail: %v", err)
	}

	fmt.Println(updated)
}
