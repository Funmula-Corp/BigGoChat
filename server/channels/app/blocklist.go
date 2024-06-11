package app

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func (a *App) AddChannelBlockUser(rctx request.CTX, channelId string, blockedId string) (*model.ChannelBlockUser, *model.AppError) {
	userId := rctx.Session().UserId
	newCBU := model.ChannelBlockUser{
		ChannelId: channelId,
		BlockedId: blockedId,
		CreateBy:  userId,
	}
	if saved, err := a.Srv().Store().Blocklist().SaveChannelBlockUser(&newCBU); err != nil {
		return nil, model.NewAppError("AddChannelBlockUser", "app.channel.add_blocklist.add.app_error", nil, "", 500).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		return saved, nil
	}
}

func (a *App) DeleteChannelBlockUser(rctx request.CTX, channelId string, blockedId string) *model.AppError {
	if err := a.Srv().Store().Blocklist().DeleteChannelBlockUser(channelId, blockedId); err != nil {
		return model.NewAppError("DeleteChannelBlockUser", "app.delete_blocklist.delete.app_error", nil, "", 500).Wrap(err)
	} else {
		return nil
	}
}

func (a *App) ListChannelBlockUsers(rctx request.CTX, channelId string) (*model.ChannelBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListChannelBlockUsers(channelId); err != nil {
		return nil, model.NewAppError("ListChannelBlockUsers", "app.channel.get_blocklist.list.app_error", nil, "", 500).Wrap(err)
	} else {
		return cbul, nil
	}
}

func (a *App) GetChannelBlockUser(rctx request.CTX, channelId string, blockedId string) (*model.ChannelBlockUser, *model.AppError) {
	if cbu, err := a.Srv().Store().Blocklist().GetChannelBlockUser(channelId, blockedId); err != nil {
		return nil, model.NewAppError("GetChannelBlockUser", "app.channel.get_blocklist.get.app_error", nil, "", 500).Wrap(err)
	} else {
		return cbu, nil
	}
}

func (a *App) ListChannelsByBlockedUser(rctx request.CTX, blockedId string) (*model.ChannelBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListChannelBlockUsersByBlockedUser(blockedId); err != nil {
		return nil, model.NewAppError("ListChannelByBlockedUser", "app.user.get_channel_blocklist.by_blocked_user.app_err", nil, "", 500).Wrap(err)
	} else {
		return cbul, nil
	}
}

func (a *App) AddUserBlockUser(rctx request.CTX, userId string, blockedId string) (*model.UserBlockUser, *model.AppError) {
	newUBU := model.UserBlockUser{
		UserId:    userId,
		BlockedId: blockedId,
	}
	if saved, err := a.Srv().Store().Blocklist().SaveUserBlockUser(&newUBU); err != nil {
		return nil, model.NewAppError("AddUserBlockUser", "app.user.add_blocklist.add.app_error", nil, "", 500).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		a.InvalidateCacheForUser(userId)
		return saved, nil
	}
}

func (a *App) DeleteUserBlockUser(rctx request.CTX, userId string, blockedId string) *model.AppError {
	if err := a.Srv().Store().Blocklist().DeleteUserBlockUser(userId, blockedId); err != nil {
		return model.NewAppError("DeleteUserBlockUser", "app.user.delete_blocklist.delete.app_error", nil, "", 500).Wrap(err)
	} else {
		return nil
	}
}

func (a *App) ListUserBlockUsers(rctx request.CTX, userId string) (*model.UserBlockUserList, *model.AppError) {
	if cub, err := a.Srv().Store().Blocklist().ListUserBlockUsers(userId); err != nil {
		return nil, model.NewAppError("ListUserBlockUsers", "app.user.get_blocklist.list.app_error", nil, "", 500).Wrap(err)
	} else {
		return cub, nil
	}
}

func (a *App) GetUserBlockUser(rctx request.CTX, userId string, blockedId string) (*model.UserBlockUser, *model.AppError) {
	if ubu, err := a.Srv().Store().Blocklist().GetUserBlockUser(userId, blockedId); err != nil {
		return nil, model.NewAppError("GetUserBlockUser", "app.user.get_blocklist.get.app_error", nil, "", 500).Wrap(err)
	} else {
		return ubu, nil
	}
}

func (a *App) ListUsersByBlockedUser(rctx request.CTX, blockedId string) (*model.UserBlockUserList, *model.AppError) {
	if ubul, err := a.Srv().Store().Blocklist().ListUserBlockUsersByBlockedUser(blockedId); err != nil {
		return nil, model.NewAppError("ListUserByBlockedUser", "app.user.get_blocklist.by_blocked_user.app_error", nil, "", 500).Wrap(err)
	} else {
		return ubul, nil
	}
}
