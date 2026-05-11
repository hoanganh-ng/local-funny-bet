package user

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	AvatarURL *string
	CreatedAt time.Time
}
