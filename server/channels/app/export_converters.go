// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"strings"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/imports"
)

func ImportLineFromTeam(team *model.TeamForExport) *imports.LineImportData {
	return &imports.LineImportData{
		Type: "team",
		Team: &imports.TeamImportData{
			Name:            &team.Name,
			DisplayName:     &team.DisplayName,
			Type:            &team.Type,
			Description:     &team.Description,
			AllowOpenInvite: &team.AllowOpenInvite,
			Scheme:          team.SchemeName,
		},
	}
}

func ImportLineFromChannel(channel *model.ChannelForExport) *imports.LineImportData {
	return &imports.LineImportData{
		Type: "channel",
		Channel: &imports.ChannelImportData{
			Team:        &channel.TeamName,
			Name:        &channel.Name,
			DisplayName: &channel.DisplayName,
			Type:        &channel.Type,
			Header:      &channel.Header,
			Purpose:     &channel.Purpose,
			Scheme:      channel.SchemeName,
			DeletedAt:   &channel.DeleteAt,
		},
	}
}

func ImportLineFromDirectChannel(channel *model.DirectChannelForExport, favoritedBy []string) *imports.LineImportData {
	channelMembers := *channel.Members
	if len(channelMembers) == 1 {
		channelMembers = []string{channelMembers[0], channelMembers[0]}
	}

	line := &imports.LineImportData{
		Type: "direct_channel",
		DirectChannel: &imports.DirectChannelImportData{
			Header:  &channel.Header,
			Members: &channelMembers,
		},
	}

	if len(favoritedBy) != 0 {
		line.DirectChannel.FavoritedBy = &favoritedBy
	}

	return line
}

func ImportLineFromUser(user *model.User, exportedPrefs map[string]*string) *imports.LineImportData {
	// Bulk Importer doesn't accept "empty string" for AuthService.
	var authService *string
	if user.AuthService != "" {
		authService = &user.AuthService
	}

	return &imports.LineImportData{
		Type: "user",
		User: &imports.UserImportData{
			Username:           &user.Username,
			Email:              &user.Email,
			AuthService:        authService,
			AuthData:           user.AuthData,
			Nickname:           &user.Nickname,
			FirstName:          &user.FirstName,
			LastName:           &user.LastName,
			Position:           &user.Position,
			Roles:              &user.Roles,
			Locale:             &user.Locale,
			Mobilephone: 		user.Mobilephone,
			UseMarkdownPreview: exportedPrefs["UseMarkdownPreview"],
			UseFormatting:      exportedPrefs["UseFormatting"],
			ShowUnreadSection:  exportedPrefs["ShowUnreadSection"],
			Theme:              exportedPrefs["Theme"],
			UseMilitaryTime:    exportedPrefs["UseMilitaryTime"],
			CollapsePreviews:   exportedPrefs["CollapsePreviews"],
			MessageDisplay:     exportedPrefs["MessageDisplay"],
			ColorizeUsernames:  exportedPrefs["ColorizeUsernames"],
			ChannelDisplayMode: exportedPrefs["ChannelDisplayMode"],
			TutorialStep:       exportedPrefs["TutorialStep"],
			EmailInterval:      exportedPrefs["EmailInterval"],
			DeleteAt:           &user.DeleteAt,
		},
	}
}

func ImportUserTeamDataFromTeamMember(member *model.TeamMemberForExport) *imports.UserTeamImportData {
	rolesList := strings.Fields(member.Roles)
	if member.SchemeAdmin {
		rolesList = append(rolesList, model.TeamAdminRoleId)
	}
	if member.SchemeModerator {
		rolesList = append(rolesList, model.TeamModeratorRoleId)
	}
	if member.SchemeVerified {
		rolesList = append(rolesList, model.TeamVerifiedRoleId)
	}
	if member.SchemeUser {
		rolesList = append(rolesList, model.TeamUserRoleId)
	}
	if member.SchemeGuest {
		rolesList = append(rolesList, model.TeamGuestRoleId)
	}
	roles := strings.Join(rolesList, " ")
	return &imports.UserTeamImportData{
		Name:  &member.TeamName,
		Roles: &roles,
	}
}

func ImportUserChannelDataFromChannelMemberAndPreferences(member *model.ChannelMemberForExport, preferences *model.Preferences) *imports.UserChannelImportData {
	rolesList := strings.Fields(member.Roles)
	if member.SchemeAdmin {
		rolesList = append(rolesList, model.ChannelAdminRoleId)
	}
	if member.SchemeVerified {
		rolesList = append(rolesList, model.ChannelVerifiedRoleId)
	}
	if member.SchemeUser {
		rolesList = append(rolesList, model.ChannelUserRoleId)
	}
	if member.SchemeGuest {
		rolesList = append(rolesList, model.ChannelGuestRoleId)
	}
	props := member.NotifyProps
	notifyProps := imports.UserChannelNotifyPropsImportData{}

	desktop, exist := props[model.DesktopNotifyProp]
	if exist {
		notifyProps.Desktop = &desktop
	}
	mobile, exist := props[model.PushNotifyProp]
	if exist {
		notifyProps.Mobile = &mobile
	}
	markUnread, exist := props[model.MarkUnreadNotifyProp]
	if exist {
		notifyProps.MarkUnread = &markUnread
	}

	favorite := false
	for _, preference := range *preferences {
		if member.ChannelId == preference.Name {
			favorite = true
		}
	}

	roles := strings.Join(rolesList, " ")
	return &imports.UserChannelImportData{
		Name:               &member.ChannelName,
		Roles:              &roles,
		NotifyProps:        &notifyProps,
		Favorite:           &favorite,
		MentionCount:       &member.MentionCount,
		MentionCountRoot:   &member.MentionCountRoot,
		UrgentMentionCount: &member.UrgentMentionCount,
		MsgCount:           &member.MsgCount,
		MsgCountRoot:       &member.MsgCountRoot,
		LastViewedAt:       &member.LastViewedAt,
	}
}

func ImportLineForPost(post *model.PostForExport) *imports.LineImportData {
	return &imports.LineImportData{
		Type: "post",
		Post: &imports.PostImportData{
			Team:     &post.TeamName,
			Channel:  &post.ChannelName,
			User:     &post.Username,
			Type:     &post.Type,
			Message:  &post.Message,
			Props:    &post.Props,
			CreateAt: &post.CreateAt,
			EditAt:   &post.EditAt,
		},
	}
}

func ImportLineForDirectPost(post *model.DirectPostForExport) *imports.LineImportData {
	channelMembers := *post.ChannelMembers
	if len(channelMembers) == 1 {
		channelMembers = []string{channelMembers[0], channelMembers[0]}
	}
	return &imports.LineImportData{
		Type: "direct_post",
		DirectPost: &imports.DirectPostImportData{
			ChannelMembers: &channelMembers,
			User:           &post.User,
			Type:           &post.Type,
			Message:        &post.Message,
			Props:          &post.Props,
			CreateAt:       &post.CreateAt,
			EditAt:         &post.EditAt,
		},
	}
}

func ImportReplyFromPost(post *model.ReplyForExport) *imports.ReplyImportData {
	return &imports.ReplyImportData{
		User:     &post.Username,
		Type:     &post.Type,
		Message:  &post.Message,
		CreateAt: &post.CreateAt,
		EditAt:   &post.EditAt,
	}
}

func ImportReactionFromPost(user *model.User, reaction *model.Reaction) *imports.ReactionImportData {
	return &imports.ReactionImportData{
		User:      &user.Username,
		EmojiName: &reaction.EmojiName,
		CreateAt:  &reaction.CreateAt,
	}
}

func ImportLineFromEmoji(emoji *model.Emoji, filePath string) *imports.LineImportData {
	return &imports.LineImportData{
		Type: "emoji",
		Emoji: &imports.EmojiImportData{
			Name:  &emoji.Name,
			Image: &filePath,
		},
	}
}

func ImportRoleDataFromRole(role *model.Role) *imports.RoleImportData {
	return &imports.RoleImportData{
		Name:          &role.Name,
		DisplayName:   &role.DisplayName,
		Description:   &role.Description,
		Permissions:   &role.Permissions,
		SchemeManaged: &role.SchemeManaged,
	}
}

func ImportLineFromRole(role *model.Role) *imports.LineImportData {
	return &imports.LineImportData{
		Type: "role",
		Role: ImportRoleDataFromRole(role),
	}
}

func ImportLineFromScheme(scheme *model.Scheme, rolesMap map[string]*model.Role) *imports.LineImportData {
	data := &imports.SchemeImportData{
		Name:        &scheme.Name,
		DisplayName: &scheme.DisplayName,
		Description: &scheme.Description,
		Scope:       &scheme.Scope,
	}

	if scheme.Scope == model.SchemeScopeTeam {
		data.DefaultTeamAdminRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultTeamAdminRole])
		data.DefaultTeamModeratorRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultTeamModeratorRole])
		data.DefaultTeamVerifiedRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultTeamVerifiedRole])
		data.DefaultTeamUserRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultTeamUserRole])
		data.DefaultTeamGuestRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultTeamGuestRole])
	}

	if scheme.Scope == model.SchemeScopeTeam || scheme.Scope == model.SchemeScopeChannel {
		data.DefaultChannelAdminRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultChannelAdminRole])
		data.DefaultChannelVerifiedRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultChannelVerifiedRole])
		data.DefaultChannelUserRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultChannelUserRole])
		data.DefaultChannelGuestRole = ImportRoleDataFromRole(rolesMap[scheme.DefaultChannelGuestRole])
	}

	return &imports.LineImportData{
		Type:   "scheme",
		Scheme: data,
	}
}
