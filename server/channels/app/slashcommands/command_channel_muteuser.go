package slashcommands

import (
	"encoding/json"
	"slices"
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/i18n"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app"
)

type MuteuserProvider struct {
}

type UnMuteuserProvider struct {
}

const (
	CmdMuteuser   = "muteuser"
	CmdUnMuteuser = "unmuteuser"
)

func init() {
	app.RegisterCommandProvider(&MuteuserProvider{})
	app.RegisterCommandProvider(&UnMuteuserProvider{})
}

func (*MuteuserProvider) GetTrigger() string {
	return CmdMuteuser
}

func (*MuteuserProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdMuteuser,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_channel_muteuser.desc"),
		AutoCompleteHint: T("api.command_channel_muteuser.hint"),
		DisplayName:      T("api.command_channel_muteuser.name"),
	}
}

func (*MuteuserProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	return muteUserCommandHandler(a, c, args, message, true)
}

func (*UnMuteuserProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdUnMuteuser,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_channel_unmuteuser.desc"),
		AutoCompleteHint: T("api.command_channel_unmuteuser.hint"),
		DisplayName:      T("api.command_channel_unmuteuser.name"),
	}
}

func (*UnMuteuserProvider) GetTrigger() string {
	return CmdUnMuteuser
}

func (*UnMuteuserProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	return muteUserCommandHandler(a, c, args, message, false)
}

func muteUserCommandHandler(a *app.App, c request.CTX, args *model.CommandArgs, message string, mute bool) *model.CommandResponse {
	channel, err := a.GetChannel(c, args.ChannelId)
	if err != nil {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_muteuser.channel.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	switch channel.Type {
	case model.ChannelTypeOpen:
		fallthrough
	case model.ChannelTypePrivate:
		if !a.HasPermissionToChannel(c, args.UserId, args.ChannelId, model.PermissionManageChannelRoles) {
			return &model.CommandResponse{
				Text:         args.T("api.command_channel_muteuser.permission.app_error"),
				ResponseType: model.CommandResponseTypeEphemeral,
			}
		}
	default:
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_muteuser.direct_group.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	splitMessage := strings.Fields(message)
	channelMembers := make(model.ChannelMembers, 0, 1)

	for _, msg := range splitMessage {
		if len(msg) != 0 && msg[0] == '@' {
			mentionName := strings.TrimPrefix(msg, "@")
			if member, err := getChannelMemberByMentionName(a, c, args.ChannelId, mentionName); err == nil {
				channelMembers = append(channelMembers, *member)
			}
		}
	}

	if len(channelMembers) == 0 {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_muteuser.param.length_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	for _, member := range channelMembers {
		muteChannelMember(a, c, &member, mute)
	}

	return &model.CommandResponse{}
}

func getChannelMemberByMentionName(a *app.App, c request.CTX, channelId, mentionName string) (*model.ChannelMember, *model.AppError) {
	user, err := a.GetUserByUsername(mentionName)
	if err != nil {
		return nil, err
	}

	return a.GetChannelMember(c, channelId, user.Id)
}

func muteChannelMember(a *app.App, c request.CTX, member *model.ChannelMember, mute bool) *model.AppError {
	excludePermissions := member.GetExcludePermissions()
	muted := slices.Contains(excludePermissions, model.PermissionCreatePost.Id)
	if mute && !muted {
		excludePermissions = append(excludePermissions, model.PermissionCreatePost.Id)
	} else if !mute && muted {
		excludePermissions = slices.DeleteFunc(excludePermissions, func(s string) bool { return s == model.PermissionCreatePost.Id })
	} else {
		// no need to update
		return nil
	}

	newMember, err := a.UpdateChannelMemberExcludePermissions(c, member.ChannelId, member.UserId, strings.Join(excludePermissions, " "))
	if err != nil {
		return err
	}

	a.Srv().Store().Channel().InvalidateAllChannelMembersForUser(newMember.UserId)

	evt := model.NewWebSocketEvent(model.WebsocketEventChannelMemberUpdated, "", "", newMember.UserId, nil, "")
	memberJSON, jsonErr := json.Marshal(newMember)
	if jsonErr == nil {
		evt.Add("channelMember", string(memberJSON))
		a.Publish(evt)
	}

	return nil
}
