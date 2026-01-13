package models

import "time"

type Comment struct {
	Id          int        `json:"id" db:"id"`
	Description string     `json:"description" db:"description"`
	TaskId      int        `json:"task_id" db:"task_id"`
	UserId      int        `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
