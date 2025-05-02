package role

import "time"

// Role is used for access control
type Role struct {
	UID          string    `json:"uid"`
	Name         string    `json:"role_name"`
	DateCreated  time.Time `json:"date_created"`
	LastSeen     time.Time `json:"last_seen"`
	LastModified time.Time `json:"last_modified"`
}
