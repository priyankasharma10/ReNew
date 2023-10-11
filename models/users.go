package models

import (
	"errors"
	"time"

	"github.com/unidoc/timestamp"
	"github.com/volatiletech/null"
)

var ErrorPasswordNotMatched = errors.New("password not matched")

type User struct {
	ID         int64               `json:"id" db:"id"`
	Name       string              `json:"name" db:"name"`
	Password   string              `json:"password" db:"password"`
	Role       string              `json:"role" db:"role"`
	Address    string              `json:"address" db:"address"`
	Email      string              `json:"email" db:"email"`
	Phone      string              `json:"phone" db:"phone"`
	Pincode    int                 `json:"pincode" db:"pincode"`
	Aadharcard int                 `json:"aadharcard" db:"aadharcard"`
	City       string              `json:"city" db:"city"`
	Country    string              `json:"country" db:"country"`
	CreatedAt  time.Time           `db:"created_at"`
	UpdatedAt  timestamp.Timestamp `db:"updated_at"`
}

type Session struct {
	ID        int64      `db:"id" json:"id"`
	Token     string     `db:"token" json:"token" sql:"not null"`
	UserID    int        `db:"user_id" json:"user_id" sql:"references users(id)"`
	StartTime time.Time  `db:"start_time" json:"start_time" sql:"not null"`
	EndTime   *time.Time `db:"end_time" json:"end_time"`
	DeviceID  string     `db:"device_id" json:"device_id"`
	Platform  string     `db:"platform" json:"platform"`
	ModelName string     `db:"model_name" json:"model_name"`
	OSVersion string     `db:"os_version" json:"os_version"`
}

type UserContextData struct {
	UserId     int    `json:"userId" db:"id"`
	SessionID  string `json:"sessionID" db:"token"`
	Name       string `json:"name" db:"name"`
	Role       string `json:"role" db:"role"`
	Address    string `json:"address" db:"address"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phone" db:"phone"`
	Pincode    int    `json:"pincode" db:"pincode"`
	Aadharcard int    `json:"aadharcard" db:"aadharcard"`
	City       string `json:"city" db:"city"`
	Country    string `json:"country" db:"city"`
}

type FetchUserData struct {
	UserId     int       `json:"userId" db:"id"`
	Name       string    `json:"name" db:"name"`
	Role       string    `json:"role" db:"role"`
	Address    string    `json:"address" db:"address"`
	Aadharcard int       `json:"aadharcard" db:"aadharcard"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	Pincode    int       `json:"pincode" db:"pincode"`
	City       string    `json:"city" db:"city"`
	Country    string    `json:"country" db:"country"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type FetchUserSessionsData struct {
	ID        int       `json:"id" db:"id"`
	User      int       `json:"userId" db:"user_id"`
	UUIDToken string    `json:"UUIDToken" db:"token"`
	EndTime   time.Time `json:"endTime" db:"end_time"`
}

type CreateNewUserRequest struct {
	Name       string              `json:"name"  db:"name"`
	Email      null.String         `json:"email"  db:"email"`
	Role       UserRoles           `json:"role"  db:"role"`
	Password   string              `json:"password"  db:"password"`
	Address    string              `json:"address" db:"address"`
	Phone      null.String         `json:"phone" db:"phone"`
	Pincode    int                 `json:"pincode" db:"pincode"`
	CreatedAt  time.Time           `db:"created_at"`
	UpdatedAt  timestamp.Timestamp `db:"updated_at"`
	City       string              `json:"city" db:"city"`
	Country    string              `json:"country" db:"country"`
	Aadharcard int                 `json:"aadharcard" db:"aadharcard"`
}

type UserData struct {
	UserID int `json:"userId" db:"id"`
}

type AuthLoginRequest struct {
	Platform  string      `json:"platform"`
	ModelName null.String `json:"modelName"`
	OSVersion null.String `json:"osVersion"`
	DeviceID  null.String `json:"deviceId"`
	Role      string      `json:"role"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
}

type GetUserDataByEmail struct {
	UserId     int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Role       UserRoles `json:"role" db:"role"`
	Address    string    `json:"address" db:"address"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	Pincode    int       `json:"pincode" db:"pincode"`
	City       string    `json:"city" db:"city"`
	Country    string    `json:"country" db:"country"`
	Aadharcard int       `json:"aadharcard" db:"aadharcard"`
	ArchivedAt time.Time `json:"archived_at" db:"archived_at"`
}

type EmailAndPassword struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
