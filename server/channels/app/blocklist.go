package app

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func (a *App) AddChannelBlockUser(rctx request.CTX, channelId string, blockedId string) (*model.ChannelBlockUser, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "", 501)
}

func (a *App) RemoveChannelBlockUser(rctx request.CTX, channelId string, blockedId string) *model.AppError {
	return model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

func (a *App) ListChannelBlockUsers(rctx request.CTX, channelId string)(*model.ChannelBlockUsers, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

func (a *App) GetChannelBlockUser(rctx request.CTX, channelId string, blockedId string)(*model.ChannelBlockUser, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "", 501)
}

func (a *App) ListChannelsByBlockedUser(rctx request.CTX, blockedId string) (*model.ChannelBlockUsers, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}


func (a *App) AddUserBlockUser(rctx request.CTX, userId string, blockedId string) (*model.UserBlockUser, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "", 501)
}

func (a *App) RemoveUserBlockUser(rctx request.CTX, userId string, blockedId string) *model.AppError {
	return model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

func (a *App) ListUserBlockUsers(rctx request.CTX, userId string)(*model.UserBlockUsers, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

func (a *App) GetUserBlockUser(rctx request.CTX, userId string, blockedId string)(*model.UserBlockUser , *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

func (a *App) ListUsersBlockedUser(rctx request.CTX, blockedId string) (*model.UserBlockUsers, *model.AppError) {
	return nil, model.NewAppError("", "", map[string]any{}, "not implemented", 501)
}

