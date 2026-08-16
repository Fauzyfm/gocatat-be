package model

import (
	"time"

)

type RoleType string

const (
	RoleUser RoleType = "user"
	RoleAdmin RoleType = "admin"
	RoleSuperAdmin RoleType = "super_admin"
)

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	ConfirmPassword string `json:"-"`
	Role RoleType `json:"role"`

	CreatedAt time.Time	`json:"createdAt"`
	UpdateAt	time.Time	`json:"updateAt"`
	DeletedAt	*time.Time	`json:"-"`

	Balance []Balance	`json:"balances,omitempty"`
}

