// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type DndProvider struct {
}

const (
	CmdDND = "dnd"
)

func init() {
	app.RegisterCommandProvider(&DndProvider{})
}

func (*DndProvider) GetTrigger() string {
	return CmdDND
}

func (*DndProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdDND,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_dnd.desc"),
		DisplayName:      T("api.command_dnd.name"),
	}
}

func (*DndProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	a.SetStatusDoNotDisturb(args.UserId)

	return &model.CommandResponse{ResponseType: model.CommandResponseTypeEphemeral, Text: args.T("api.command_dnd.success")}
}
