// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type OnlineProvider struct {
}

const (
	CmdOnline = "online"
)

func init() {
	app.RegisterCommandProvider(&OnlineProvider{})
}

func (*OnlineProvider) GetTrigger() string {
	return CmdOnline
}

func (*OnlineProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdOnline,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_online.desc"),
		DisplayName:      T("api.command_online.name"),
	}
}

func (*OnlineProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	a.SetStatusOnline(args.UserId, true)

	return &model.CommandResponse{ResponseType: model.CommandResponseTypeEphemeral, Text: args.T("api.command_online.success")}
}
