package repository

import (
	"errors"
	"github.com/third-place/user-service/internal/entity"
	"github.com/jinzhu/gorm"
)

type InviteRepository struct {
	conn *gorm.DB
}

func CreateInviteRepository(conn *gorm.DB) *InviteRepository {
	return &InviteRepository{conn}
}

func (r *InviteRepository) FindInvites(offset int) []*entity.Invite {
	var invites []*entity.Invite
	r.conn.Table("invites").
		Where("invites.deleted_at IS NULL").
		Order("id desc").
		Limit(25).
		Offset(offset).
		Find(&invites)
	return invites
}

func (r *InviteRepository) FindOneByCode(code string) (*entity.Invite, error) {
	invite := &entity.Invite{}
	r.conn.Where("code = ?", code).Find(invite)
	if invite.ID == 0 {
		return nil, errors.New("no invite found")
	}
	return invite, nil
}

func (r *InviteRepository) Create(invite *entity.Invite) *gorm.DB {
	return r.conn.Create(invite)
}

func (r *InviteRepository) Save(invite *entity.Invite) *gorm.DB {
	return r.conn.Save(invite)
}
