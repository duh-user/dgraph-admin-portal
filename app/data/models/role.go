package models

import "time"

// Role is used for access control
type Role struct {
	UID          string    `json:"uid"`
	Name         string    `json:"role_name,omitempty"`
	DateCreated  time.Time `json:"date_created,omitempty"`
	LastSeen     time.Time `json:"last_seen,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	ReverseEdge  []User    `json:"~role,omitempty"`
}
