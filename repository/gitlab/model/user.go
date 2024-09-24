package model

type User struct {
	ID       int
	Username string
	Name     string
	WebURL   string
	Email    string
	Avatar   UserAvatar
}

type UserAvatar struct {
	AvatarURL         *string
	FallbackAvatarURL *string
}
