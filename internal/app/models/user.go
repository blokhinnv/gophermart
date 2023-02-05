package models

type User struct {
	Username       string
	HashedPassword string
	Salt           string
}
