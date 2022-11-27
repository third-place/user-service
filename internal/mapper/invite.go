package mapper

import (
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/model"
)

func MapInviteEntityToModel(invite *entity.Invite) *model.Invite {
	return &model.Invite{
		Code:    invite.Code,
		Claimed: invite.Claimed,
	}
}

func MapInviteEntitiesToModels(invites []*entity.Invite) []*model.Invite {
	inviteModels := make([]*model.Invite, len(invites))
	for i, v := range invites {
		inviteModels[i] = MapInviteEntityToModel(v)
	}
	return inviteModels
}
