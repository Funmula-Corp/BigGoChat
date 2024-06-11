package api4

import (
	"encoding/json"
	"net/http"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

func (api *API) InitBlocklist() {
	api.BaseRoutes.User.Handle("/blockuser", api.APISessionRequired(listUserBlockUsers)).Methods("GET")
	api.BaseRoutes.User.Handle("/blockuser/{blocked_user_id:[A-za-z0-9]+}", api.APISessionRequired(addUserBlockUser)).Methods("PUT")
	api.BaseRoutes.User.Handle("/blockuser/{blocked_user_id:[A-za-z0-9]+}", api.APISessionRequired(deleteUserBlockUser)).Methods("DELETE")

	api.BaseRoutes.Channel.Handle("/blockuser", api.APISessionRequired(getChannelBlockUsers)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/blockuser/{blocked_user_id:[A-za-z0-9]+}", api.APISessionRequired(addChannelBlockUser)).Methods("PUT")
	api.BaseRoutes.Channel.Handle("/blockuser/{blocked_user_id:[A-za-z0-9]+}", api.APISessionRequired(deleteChannelBlockUser)).Methods("DELETE")
}

func addUserBlockUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	c.RequireBlockedId()
	if c.Err != nil {
		return
	}

	blockedId := c.Params.BlockedUserId
	// userId := c.AppContext.Session().UserId
	if !c.App.SessionHasPermissionToUserOrBot(c.AppContext, *c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}
	uBU, err := c.App.AddUserBlockUser(c.AppContext, c.Params.UserId, blockedId)
	if err != nil {
		c.Err = err
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(uBU); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func deleteUserBlockUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	c.RequireBlockedId()
	if c.Err != nil {
		return
	}
	blockedId := c.Params.BlockedUserId
	if !c.App.SessionHasPermissionToUserOrBot(c.AppContext, *c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}
	if err := c.App.DeleteUserBlockUser(c.AppContext, c.Params.UserId, blockedId); err != nil {
		c.Err = err
		return
	}
}

func listUserBlockUsers(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}
	if !c.App.SessionHasPermissionToUserOrBot(c.AppContext, *c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}
	blockList, err := c.App.ListUserBlockUsers(c.AppContext, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}
	if err := json.NewEncoder(w).Encode(blockList); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func addChannelBlockUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	c.RequireBlockedId()
	if c.Err != nil {
		return
	}
	channel, err := c.App.GetChannel(c.AppContext, c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}
	userId := c.AppContext.Session().UserId
	if _, err := c.App.GetUser(userId); err != nil {
		c.Err = err
		return
	}

	switch channel.Type {
	case model.ChannelTypePrivate:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
			c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
			return
		}
	case model.ChannelTypeOpen:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePublicChannelMembers) {
			c.SetPermissionError(model.PermissionManagePublicChannelMembers)
			return
		}
	case model.ChannelTypeDirect:
		c.Err = model.NewAppError("addChannelBlockUser", "api.channel.add_blocklist.channel_type_direct.app_error", nil, "", http.StatusBadRequest)
		return
	case model.ChannelTypeGroup:
		if _, errGet := c.App.GetChannelMember(c.AppContext, channel.Id, userId); errGet != nil {
			c.Err = model.NewAppError("addChannelBlockUser", "api.channel.add_blocklist.forbidden.app_error", nil, "", http.StatusForbidden)
			return
		}
	default:
		c.Err = model.NewAppError("addChannelBlockUser", "api.channel.add_blocklist.unkown_channel_type.app_error", nil, "", http.StatusBadRequest)
		return
	}

	blockedId := c.Params.BlockedUserId
	var cbu *model.ChannelBlockUser
	var errApp *model.AppError

	if cbu, errApp = c.App.AddChannelBlockUser(c.AppContext, c.Params.ChannelId, blockedId); errApp != nil {
		c.Err = errApp
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(cbu); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}

func deleteChannelBlockUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	c.RequireBlockedId()
	if c.Err != nil {
		return
	}
	channel, err := c.App.GetChannel(c.AppContext, c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	blockedId := c.Params.BlockedUserId
	userId := c.AppContext.Session().UserId
	switch channel.Type {
	case model.ChannelTypePrivate:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
			c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
			return
		}
	case model.ChannelTypeOpen:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePublicChannelMembers) {
			c.SetPermissionError(model.PermissionManagePublicChannelMembers)
			return
		}
	case model.ChannelTypeDirect:
		c.Err = model.NewAppError("deleteChannelBlockUser", "api.channel.delete_blocklist.channel_type_direct.app_error", nil, "", http.StatusBadRequest)
		return
	case model.ChannelTypeGroup:
		if _, errGet := c.App.GetChannelMember(c.AppContext, channel.Id, userId); errGet != nil {
			c.Err = model.NewAppError("deleteChannelBlockUser", "api.channel.delete_blocklist.forbidden.app_error", nil, "", http.StatusForbidden)
			return
		}
	default:
		c.Err = model.NewAppError("deleteChannelBlockUser", "api.channel.delete_blocklist.unkown_channel_type.app_error", nil, "", http.StatusBadRequest)
		return
	}
	if errApp := c.App.DeleteChannelBlockUser(c.AppContext, c.Params.ChannelId, blockedId); errApp != nil {
		c.Err = errApp
		return
	}
}

func getChannelBlockUsers(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}
	userId := c.AppContext.Session().UserId
	channel, err := c.App.GetChannel(c.AppContext, c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}
	switch channel.Type {
	case model.ChannelTypePrivate:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
			c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
			return
		}
	case model.ChannelTypeOpen:
		if !c.App.SessionHasPermissionToChannel(c.AppContext, *c.AppContext.Session(), channel.Id, model.PermissionManagePublicChannelMembers) {
			c.SetPermissionError(model.PermissionManagePublicChannelMembers)
			return
		}
	case model.ChannelTypeDirect:
		c.Err = model.NewAppError("getChannelBlockUsers", "api.channel.get_blocklist.channel_type_direct.app_error", nil, "", http.StatusBadRequest)
		return
	case model.ChannelTypeGroup:
		if _, errGet := c.App.GetChannelMember(c.AppContext, channel.Id, userId); errGet != nil {
			c.Err = model.NewAppError("getChannelBlockUsers", "api.channel.get_blocklist.forbidden.app_error", nil, "", http.StatusForbidden)
			return
		}
	default:
		c.Err = model.NewAppError("getChannelBlockUsers", "api.channel.get_blocklist.unkown_channel_type.app_error", nil, "", http.StatusBadRequest)
		return
	}

	var cBUL *model.ChannelBlockUserList
	var errApp *model.AppError

	if cBUL, errApp = c.App.ListChannelBlockUsers(c.AppContext, c.Params.ChannelId); errApp != nil {
		c.Err = errApp
		return
	}
	if err := json.NewEncoder(w).Encode(cBUL); err != nil {
		c.Logger.Warn("Error while writing response", mlog.Err(err))
	}
}
