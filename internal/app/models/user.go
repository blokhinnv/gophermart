package models

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Salt           string
}
