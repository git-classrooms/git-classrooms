package model

type User struct {
	ID         int
	Username   string
	Name       string
	WebUrl     string
	Email      string
	Avatar     UserAvatar
	Permission *AccessLevelValue
}

type UserAvatar struct {
	AvatarURL         *string
	FallbackAvatarURL *string
}
