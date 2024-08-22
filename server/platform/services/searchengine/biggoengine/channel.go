package biggoengine

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func (be *BiggoEngine) DeleteChannel(channel *model.Channel) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexChannel(rctx request.CTX, channel *model.Channel, userIDs, teamMemberIDs []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexChannelsBulk(channels []*model.Channel) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchChannels(teamId, userID, term string, isGuest bool) (result []string, aErr *model.AppError) {
	return
}
