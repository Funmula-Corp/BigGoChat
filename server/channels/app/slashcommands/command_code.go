// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/i18n"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

type CodeProvider struct {
}

const (
	CmdCode = "code"
)

func init() {
	app.RegisterCommandProvider(&CodeProvider{})
}

func (*CodeProvider) GetTrigger() string {
	return CmdCode
}

func (*CodeProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdCode,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_code.desc"),
		AutoCompleteHint: T("api.command_code.hint"),
		DisplayName:      T("api.command_code.name"),
	}
}

func (*CodeProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	if message == "" {
		return &model.CommandResponse{Text: args.T("api.command_code.message.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}
	rmsg := "    " + strings.Join(strings.Split(message, "\n"), "\n    ")
	return &model.CommandResponse{ResponseType: model.CommandResponseTypeInChannel, Text: rmsg, SkipSlackParsing: true}
}
