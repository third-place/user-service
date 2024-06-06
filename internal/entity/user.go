package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/model"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CognitoId     uuid.UUID
	Name          string
	Username      string `gorm:"unique;not null"`
	ProfilePic    string
	BioMessage    string
	Role          string `gorm:"default:'user'"`
	IsBanned      bool   `gorm:"default:false"`
	AddressStreet string
	AddressCity   string
	AddressZip    string
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"null"`
	Birthday      string
	Verified      bool `gorm:"not null"`
	InviteID      uint
	OTP           string
	Emails        []*Email
	Passwords     []*Password
}

func (u *User) UpdateUserFromModel(user *model.User) {
	u.Name = user.Name
	u.BioMessage = user.BioMessage
	u.Birthday = user.Birthday
	if user.ProfilePic != "" {
		u.ProfilePic = user.ProfilePic
	}
}

func (u *User) ToJson() []byte {
	data, _ := json.Marshal(u)
	return data
}
