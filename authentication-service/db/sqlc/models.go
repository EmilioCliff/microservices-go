// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"
)

type User struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int32
	CreatedAt time.Time
	UpdateAt  time.Time
}
