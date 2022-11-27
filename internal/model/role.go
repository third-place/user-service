package model

type Role string

// List of Role
const (
	USER      Role = "user"
	MODERATOR Role = "moderator"
	ADMIN     Role = "admin"
)
