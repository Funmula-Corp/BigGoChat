package slashcommands

import (
	"encoding/json"
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/i18n"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app"
)

type AnnouncProvider struct {
}

const (
	CmdAnnounc = "announc"
)

func init() {
	app.RegisterCommandProvider(&AnnouncProvider{})
}

func (*AnnouncProvider) GetTrigger() string {
	return CmdAnnounc
}

func (*AnnouncProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdAnnounc,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_channel_announc.desc"),
		AutoCompleteHint: T("api.command_channel_announc.hint"),
		DisplayName:      T("api.command_channel_announc.name"),
	}
}

func (*AnnouncProvider) DoCommand(a *app.App, c request.CTX, args *model.CommandArgs, message string) *model.CommandResponse {
	channel, err := a.GetChannel(c, args.ChannelId)
	if err != nil {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_announc.channel.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	switch channel.Type {
	case model.ChannelTypeOpen:
		fallthrough
	case model.ChannelTypePrivate:
		if !a.HasPermissionToChannel(c, args.UserId, args.ChannelId, model.PermissionManageChannelRoles) {
			return &model.CommandResponse{
				Text:         args.T("api.command_channel_announc.permission.app_error"),
				ResponseType: model.CommandResponseTypeEphemeral,
			}
		}
	default:
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_announc.direct_group.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	params := strings.Fields(args.Command)[1:]
	active := false

	if len(params) == 0 {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_announc.param.length_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	} else if params[0] == "yes" {
		active = true
	} else if params[0] == "no" {
		active = false
	} else {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_announc.param.param_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	if active && (channel.SchemeId == nil || *channel.SchemeId == "") {
		channel.SchemeId = model.NewString(model.ChannelReadOnlySchemeId)
		_, err = a.UpdateChannel(c, channel)
	} else if !active && channel.SchemeId != nil && *channel.SchemeId == model.ChannelReadOnlySchemeId {
		channel.SchemeId = model.NewString("")
		_, err = a.UpdateChannel(c, channel)
	}

	// plan B, use channel moderation
	// createPost := model.PermissionCreatePost.Id
	// _, err = a.PatchChannelModerationsForChannel(c, channel, []*model.ChannelModerationPatch{
	// 	{
	// 		Name: &createPost,
	// 		Roles: &model.ChannelModeratedRolesPatch{Members: model.NewBool(false), Guests: model.NewBool(false), Verified: model.NewBool(false)},
	// 	},
	// })

	if err != nil {
		text := args.T("api.command_channel_announc.update_channel.app_error")

		return &model.CommandResponse{
			Text:         text,
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	forEachChannelMember(a, channel.Id, func(channelMember model.ChannelMember) error {
		a.Srv().Store().Channel().InvalidateAllChannelMembersForUser(channelMember.UserId)

		evt := model.NewWebSocketEvent(model.WebsocketEventChannelMemberUpdated, "", "", channelMember.UserId, nil, "")
		memberJSON, jsonErr := json.Marshal(channelMember)
		if jsonErr != nil {
			return jsonErr
		}
		evt.Add("channelMember", string(memberJSON))
		a.Publish(evt)

		return nil
	})

	return &model.CommandResponse{}
}

func forEachChannelMember(a *app.App, channelID string, f func(model.ChannelMember) error) error {
	perPage := 100
	page := 0

	for {
		channelMembers, err := a.Srv().Store().Channel().GetMembers(channelID, page*perPage, perPage)
		if err != nil {
			return err
		}

		for _, channelMember := range channelMembers {
			if err = f(channelMember); err != nil {
				return err
			}
		}

		length := len(channelMembers)
		if length < perPage {
			break
		}

		page++
	}

	return nil
}