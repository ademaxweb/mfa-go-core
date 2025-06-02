package data

import "time"

type User struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}
