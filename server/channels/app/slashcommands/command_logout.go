// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/i18n"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

type LogoutProvider struct {
}

const (
	CmdLogout = "logout"
)

func init() {
	app.RegisterCommandProvider(&LogoutProvider{})
}

func (*LogoutProvider) GetTrigger() string {
	return CmdLogout
}

func (*LogoutProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdLogout,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_logout.desc"),
		AutoCompleteHint: "",
		DisplayName:      T("api.command_logout.name"),
	}
}

func (*LogoutProvider) DoCommand(a *app.App, _ request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	// Actual logout is handled client side.
	return &model.CommandResponse{GotoLocation: "/login"}
}
