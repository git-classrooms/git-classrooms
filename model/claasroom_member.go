package model

type MemberRoll int

const (
	Student MemberRoll = iota
	Prof
)

type ClassroomMember struct {
	Roll MemberRoll
	User User
}
