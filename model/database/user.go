package database

import (
	"gorm.io/gorm"
	"time"
)

// User is the representation of the user in database
type User struct {
	ID                     int `gorm:"primary_key;autoIncrement:false"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt   `gorm:"index"`
	Classrooms             []UserClassrooms `gorm:"foreignKey:UserID"`
	AssignmentRepositories []AssignmentProjects
}

type UserDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	WebUrl    string `json:"webUrl"`
	AvatarUrl string `json:"avatarUrl"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
}

// TODO: Implement this function
// func ToDto(user User, user GitlabUser) UserDTO {
//	 return UserDTO{
//	 	ID:   user.ID,
//	 	Name: user.Name,
//	 }
// }
