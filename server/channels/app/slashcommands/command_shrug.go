// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type ShrugProvider struct {
}

const (
	CmdShrug = "shrug"
)

func init() {
	app.RegisterCommandProvider(&ShrugProvider{})
}

func (*ShrugProvider) GetTrigger() string {
	return CmdShrug
}

func (*ShrugProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdShrug,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_shrug.desc"),
		AutoCompleteHint: T("api.command_shrug.hint"),
		DisplayName:      T("api.command_shrug.name"),
	}
}

func (*ShrugProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	rmsg := `¯\\\_(ツ)\_/¯`
	if message != "" {
		rmsg = message + " " + rmsg
	}

	return &model.CommandResponse{ResponseType: model.CommandResponseTypeInChannel, Text: rmsg}
}
