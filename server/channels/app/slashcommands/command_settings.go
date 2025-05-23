// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type SettingsProvider struct {
}

const (
	CmdSettings = "settings"
)

func init() {
	app.RegisterCommandProvider(&SettingsProvider{})
}

func (settings *SettingsProvider) GetTrigger() string {
	return CmdSettings
}

func (settings *SettingsProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdSettings,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_settings.desc"),
		AutoCompleteHint: "",
		DisplayName:      T("api.command_settings.name"),
	}
}

func (settings *SettingsProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	// This command is handled client-side and shouldn't hit the server.
	return &model.CommandResponse{
		Text:         args.T("api.command_settings.unsupported.app_error"),
		ResponseType: model.CommandResponseTypeEphemeral,
	}
}
