package models

import "time"

type AuthUser struct {
	Id        int
	Email     string
	Password  string
	DeletedAt *time.Time
}
