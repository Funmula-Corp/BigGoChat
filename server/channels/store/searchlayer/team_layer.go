// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchlayer

import (
	store "git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
	model "git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

type SearchTeamStore struct {
	store.TeamStore
	rootStore *SearchStore
}

func (s SearchTeamStore) SaveMember(rctx request.CTX, teamMember *model.TeamMember, maxUsersPerTeam int) (*model.TeamMember, error) {
	member, err := s.TeamStore.SaveMember(rctx, teamMember, maxUsersPerTeam)
	if err == nil {
		s.rootStore.indexUserFromID(rctx, member.UserId)
	}
	return member, err
}

func (s SearchTeamStore) UpdateMember(rctx request.CTX, teamMember *model.TeamMember) (*model.TeamMember, error) {
	member, err := s.TeamStore.UpdateMember(rctx, teamMember)
	if err == nil {
		s.rootStore.indexUserFromID(rctx, member.UserId)
	}
	return member, err
}

func (s SearchTeamStore) RemoveMember(rctx request.CTX, teamId string, userId string) error {
	err := s.TeamStore.RemoveMember(rctx, teamId, userId)
	if err == nil {
		s.rootStore.indexUserFromID(rctx, userId)
	}
	return err
}

func (s SearchTeamStore) RemoveAllMembersByUser(rctx request.CTX, userId string) error {
	err := s.TeamStore.RemoveAllMembersByUser(rctx, userId)
	if err == nil {
		s.rootStore.indexUserFromID(rctx, userId)
	}
	return err
}
