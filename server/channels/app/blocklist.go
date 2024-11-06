package app

import (
	"errors"
	"net/http"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
)

func (a *App) AddTeamBlockUser(rctx request.CTX, teamId string, blockedId string) (*model.TeamBlockUser, *model.AppError) {
	userId := rctx.Session().UserId
	newTBU := model.TeamBlockUser{
		TeamId: teamId,
		BlockedId: blockedId,
		CreateBy:  userId,
	}
	if userId == blockedId {
		return nil, model.NewAppError("AddTeamBlockUser", "app.team.add_blocklist.add_self.app_err", nil, "", http.StatusBadRequest)
	}
	if tm, err := a.GetTeamMember(rctx, teamId, blockedId); err != nil{
		var nfErr *store.ErrNotFound
		if !errors.As(err, &nfErr) {
			return nil, err
		}
	}else if tm.SchemeAdmin || tm.SchemeModerator {
		return nil, model.NewAppError("AddTeamBlockUser", "app.team.add_blocklist.team_moderator.app_err", nil, "", http.StatusBadRequest)
	}
	if err := a.RemoveUserFromTeam(rctx, teamId, blockedId, userId); err != nil {
		// block a user who is not a team member.
		if err.Id != "api.team.remove_user_from_team.missing.app_error" {
			return nil, err
		}
	}
	if saved, err := a.Srv().Store().Blocklist().SaveTeamBlockUser(&newTBU); err != nil {
		return nil, model.NewAppError("AddTeamBlockUser", "app.team.add_blocklist.add.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		return saved, nil
	}
}

func (a *App) DeleteTeamBlockUser(rctx request.CTX, teamId string, blockedId string) *model.AppError {
	if err := a.Srv().Store().Blocklist().DeleteTeamBlockUser(teamId, blockedId); err != nil {
		return model.NewAppError("DeleteTeamBlockUser", "app.delete_blocklist.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return nil
	}
}

func (a *App) ListTeamBlockUsers(rctx request.CTX, teamId string) (*model.TeamBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListTeamBlockUsers(teamId); err != nil {
		return nil, model.NewAppError("ListTeamBlockUsers", "app.team.get_blocklist.list.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbul, nil
	}
}

func (a *App) GetTeamBlockUser(rctx request.CTX, teamId string, blockedId string) (*model.TeamBlockUser, *model.AppError) {
	if cbu, err := a.Srv().Store().Blocklist().GetTeamBlockUser(teamId, blockedId); err != nil {
		return nil, model.NewAppError("GetTeamBlockUser", "app.team.get_blocklist.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbu, nil
	}
}

func (a *App) ListTeamsByBlockedUser(rctx request.CTX, blockedId string) (*model.TeamBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListTeamBlockUsersByBlockedUser(blockedId); err != nil {
		return nil, model.NewAppError("ListTeamByBlockedUser", "app.user.get_team_blocklist.by_blocked_user.app_err", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbul, nil
	}
}

//
func (a *App) AddChannelBlockUser(rctx request.CTX, channelId string, blockedId string) (*model.ChannelBlockUser, *model.AppError) {
	userId := rctx.Session().UserId
	newCBU := model.ChannelBlockUser{
		ChannelId: channelId,
		BlockedId: blockedId,
		CreateBy:  userId,
	}
	if userId == blockedId {
		return nil, model.NewAppError("AddChannelBlockUser", "app.channel.add_blocklist.add_self.app_err", nil, "", http.StatusBadRequest)
	}
	if cm, err := a.GetChannelMember(rctx, channelId, blockedId); err != nil {
		var nfErr *store.ErrNotFound
		if !errors.As(err, &nfErr) {
			return nil, err
		}
	} else if cm.SchemeAdmin {
		return nil, model.NewAppError("AddChanelBlockUser", "app.channel.add_blocklist.channel_admin.app_err", nil, "", http.StatusBadRequest)
	}
	if saved, err := a.Srv().Store().Blocklist().SaveChannelBlockUser(&newCBU); err != nil {
		return nil, model.NewAppError("AddChannelBlockUser", "app.channel.add_blocklist.add.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		return saved, nil
	}
}

func (a *App) DeleteChannelBlockUser(rctx request.CTX, channelId string, blockedId string) *model.AppError {
	if err := a.Srv().Store().Blocklist().DeleteChannelBlockUser(channelId, blockedId); err != nil {
		return model.NewAppError("DeleteChannelBlockUser", "app.delete_blocklist.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return nil
	}
}

func (a *App) ListChannelBlockUsers(rctx request.CTX, channelId string) (*model.ChannelBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListChannelBlockUsers(channelId); err != nil {
		return nil, model.NewAppError("ListChannelBlockUsers", "app.channel.get_blocklist.list.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbul, nil
	}
}

func (a *App) GetChannelBlockUser(rctx request.CTX, channelId string, blockedId string) (*model.ChannelBlockUser, *model.AppError) {
	if cbu, err := a.Srv().Store().Blocklist().GetChannelBlockUser(channelId, blockedId); err != nil {
		return nil, model.NewAppError("GetChannelBlockUser", "app.channel.get_blocklist.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbu, nil
	}
}

func (a *App) ListChannelsByBlockedUser(rctx request.CTX, blockedId string) (*model.ChannelBlockUserList, *model.AppError) {
	if cbul, err := a.Srv().Store().Blocklist().ListChannelBlockUsersByBlockedUser(blockedId); err != nil {
		return nil, model.NewAppError("ListChannelByBlockedUser", "app.user.get_channel_blocklist.by_blocked_user.app_err", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cbul, nil
	}
}

func (a *App) AddUserBlockUser(rctx request.CTX, userId string, blockedId string) (*model.UserBlockUser, *model.AppError) {
	newUBU := model.UserBlockUser{
		UserId:    userId,
		BlockedId: blockedId,
	}
	if userId == blockedId {
		return nil, model.NewAppError("AddUserBlockUser", "app.user.add_blocklist.add_self.app_err", nil, "", http.StatusBadRequest)
	}
	if _, err := a.Srv().Store().User().Get(rctx.Context(), blockedId); err != nil {
		return nil, model.NewAppError("AddUserBlockUser", MissingAccountError, nil, "", http.StatusBadRequest).Wrap(err)
	}
	if saved, err := a.Srv().Store().Blocklist().SaveUserBlockUser(&newUBU); err != nil {
		return nil, model.NewAppError("AddUserBlockUser", "app.user.add_blocklist.save.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		a.InvalidateCacheForUser(userId)
		return saved, nil
	}
}

func (a *App) DeleteUserBlockUser(rctx request.CTX, userId string, blockedId string) *model.AppError {
	var user, blockedUser *model.User
	var err error
	if user, err = a.Srv().Store().User().Get(rctx.Context(), userId); err != nil {
		return model.NewAppError("DeleteUserBlockUser", MissingAccountError, nil, "", http.StatusInternalServerError).Wrap(err)
	}
	if blockedUser, err = a.Srv().Store().User().Get(rctx.Context(), blockedId); err != nil {
		return model.NewAppError("DeleteUserBlockUser", MissingAccountError, nil, "", http.StatusInternalServerError).Wrap(err)
	}
	if err = a.Srv().Store().Blocklist().DeleteUserBlockUser(userId, blockedId, user.IsVerified(), blockedUser.IsVerified()); err != nil {
		return model.NewAppError("DeleteUserBlockUser", "app.user.delete_blocklist.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		a.InvalidateCacheForUser(blockedId)
		a.InvalidateCacheForUser(userId)
		return nil
	}
}

func (a *App) ListUserBlockUsers(rctx request.CTX, userId string) (*model.UserBlockUserList, *model.AppError) {
	if cub, err := a.Srv().Store().Blocklist().ListUserBlockUsers(userId); err != nil {
		return nil, model.NewAppError("ListUserBlockUsers", "app.user.get_blocklist.list.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return cub, nil
	}
}

func (a *App) GetUserBlockUser(rctx request.CTX, userId string, blockedId string) (*model.UserBlockUser, *model.AppError) {
	if ubu, err := a.Srv().Store().Blocklist().GetUserBlockUser(userId, blockedId); err != nil {
		return nil, model.NewAppError("GetUserBlockUser", "app.user.get_blocklist.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return ubu, nil
	}
}

func (a *App) ListUsersByBlockedUser(rctx request.CTX, blockedId string) (*model.UserBlockUserList, *model.AppError) {
	if ubul, err := a.Srv().Store().Blocklist().ListUserBlockUsersByBlockedUser(blockedId); err != nil {
		return nil, model.NewAppError("ListUserByBlockedUser", "app.user.get_blocklist.by_blocked_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	} else {
		return ubul, nil
	}
}
