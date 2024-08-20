package biggoengine

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func (be *BiggoEngine) DeleteUser(user *model.User) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexUser(rctx request.CTX, user *model.User, teamsIds, channelsIds []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (result []string, aErr *model.AppError) {
	return
}
