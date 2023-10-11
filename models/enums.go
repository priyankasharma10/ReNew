package models

type Context string

type UserRoles string

const (
	UserContext Context   = "userContext"
	Admin       UserRoles = "admin"
	Users       UserRoles = "User"
)
