// Package database contains reference types for representing data with gorm
package database

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User is the representation of the user in database
type User struct {
	ID              int               `gorm:"primary_key;autoIncrement:false" json:"id"`
	GitlabUsername  string            `gorm:"unique;not null" json:"gitlabUsername"`
	GitlabEmail     string            `gorm:"unique;not null" json:"gitlabEmail"`
	GitLabAvatar    UserAvatar        `json:"gitlabAvatar"`
	Name            string            `gorm:"not null" json:"name"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	OwnedClassrooms []*Classroom      `gorm:"foreignKey:OwnerID" json:"-"`
	Classrooms      []*UserClassrooms `gorm:"foreignKey:UserID" json:"-"`
} //@Name User

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.ID).Delete(&UserClassrooms{})
	tx.Clauses(clause.Returning{}).Where("owner_id = ?", u.ID).Delete(&Classroom{})
	return
}
