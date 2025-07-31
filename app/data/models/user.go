package models

import (
	"time"
)

// User is a generac type used to reperesent users in all roles
type User struct {
	UID          string    `json:"uid"`
	Name         string    `json:"name"`
	UserName     string    `json:"user_name"`
	PassHash     string    `json:"pass_hash"`
	Email        string    `json:"email"`
	Role         []Role    `json:"role"`
	DateCreated  time.Time `json:"date_created"`
	LastSeen     time.Time `json:"last_seen"`
	LastModified time.Time `json:"last_modified"`
}

// NewUser is used to hold details during user creation
type NewUser struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Pass     string `json:"pass"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
