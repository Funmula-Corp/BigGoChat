// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type LeaveProvider struct {
}

const (
	CmdLeave = "leave"
)

func init() {
	app.RegisterCommandProvider(&LeaveProvider{})
}

func (*LeaveProvider) GetTrigger() string {
	return CmdLeave
}

func (*LeaveProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdLeave,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_leave.desc"),
		DisplayName:      T("api.command_leave.name"),
	}
}

func (*LeaveProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	var channel *model.Channel
	var noChannelErr *model.AppError
	if channel, noChannelErr = a.GetChannel(c, args.ChannelId); noChannelErr != nil {
		return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	team, err := a.GetTeam(args.TeamId)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	err = a.LeaveChannel(c, args.ChannelId, args.UserId)
	if err != nil {
		if channel.Name == model.DefaultChannelName {
			return &model.CommandResponse{Text: args.T("api.channel.leave.default.app_error", map[string]any{"Channel": model.DefaultChannelName}), ResponseType: model.CommandResponseTypeEphemeral}
		}
		return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	member, err := a.GetTeamMember(c, team.Id, args.UserId)
	if err != nil || member.DeleteAt != 0 {
		return &model.CommandResponse{GotoLocation: args.SiteURL + "/"}
	}

	user, err := a.GetUser(args.UserId)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	if user.IsGuest() {
		members, err := a.GetChannelMembersForUser(c, team.Id, args.UserId)
		if err != nil || len(members) == 0 {
			return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
		channel, err := a.GetChannel(c, members[0].ChannelId)
		if err != nil {
			return &model.CommandResponse{Text: args.T("api.command_leave.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
		return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + channel.Name}
	}

	return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + model.DefaultChannelName}
}
