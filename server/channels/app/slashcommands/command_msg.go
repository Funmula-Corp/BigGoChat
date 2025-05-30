// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"errors"
	"strings"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type msgProvider struct {
}

const (
	CmdMsg = "msg"
)

func init() {
	app.RegisterCommandProvider(&msgProvider{})
}

func (*msgProvider) GetTrigger() string {
	return CmdMsg
}

func (*msgProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdMsg,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_msg.desc"),
		AutoCompleteHint: T("api.command_msg.hint"),
		DisplayName:      T("api.command_msg.name"),
	}
}

func (*msgProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	splitMessage := strings.SplitN(message, " ", 2)

	parsedMessage := ""
	targetUsername := ""

	if len(splitMessage) > 1 {
		parsedMessage = strings.SplitN(message, " ", 2)[1]
	}
	targetUsername = strings.SplitN(message, " ", 2)[0]
	targetUsername = strings.TrimPrefix(targetUsername, "@")

	userProfile, nErr := a.Srv().Store().User().GetByUsername(targetUsername)
	if nErr != nil {
		c.Logger().Error(nErr.Error())
		return &model.CommandResponse{Text: args.T("api.command_msg.missing.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	if userProfile.Id == args.UserId {
		return &model.CommandResponse{Text: args.T("api.command_msg.missing.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	canSee, err := a.UserCanSeeOtherUser(c, args.UserId, userProfile.Id)
	if err != nil {
		c.Logger().Error(err.Error())
		return &model.CommandResponse{Text: args.T("api.command_msg.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}
	if !canSee {
		return &model.CommandResponse{Text: args.T("api.command_msg.missing.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	// Find the channel based on this user
	channelName := model.GetDMNameFromIds(args.UserId, userProfile.Id)

	targetChannelID := ""
	if channel, channelErr := a.Srv().Store().Channel().GetByName(args.TeamId, channelName, true); channelErr != nil {
		var nfErr *store.ErrNotFound
		if errors.As(channelErr, &nfErr) {
			if !a.HasPermissionTo(args.UserId, model.PermissionCreateDirectChannel) {
				return &model.CommandResponse{Text: args.T("api.command_msg.permission.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
			}

			var directChannel *model.Channel
			if directChannel, err = a.GetOrCreateDirectChannel(c, args.UserId, userProfile.Id); err != nil {
				c.Logger().Error(err.Error())
				return &model.CommandResponse{Text: args.T(err.Id), ResponseType: model.CommandResponseTypeEphemeral}
			}
			targetChannelID = directChannel.Id
		} else {
			c.Logger().Error(channelErr.Error())
			return &model.CommandResponse{Text: args.T("api.command_msg.dm_fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
	} else {
		targetChannelID = channel.Id
	}

	if parsedMessage != "" {
		post := &model.Post{}
		post.Message = parsedMessage
		post.ChannelId = targetChannelID
		post.UserId = args.UserId
		if _, err = a.CreatePostMissingChannel(c, post, true, true); err != nil {
			return &model.CommandResponse{Text: args.T("api.command_msg.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
	}

	team, err := a.GetTeam(args.TeamId)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_msg.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + channelName, Text: "", ResponseType: model.CommandResponseTypeEphemeral}
}
