package entities

import "time"

type User struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Role         Role
	Username     string
	Email        string
	PasswordHash string
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)
