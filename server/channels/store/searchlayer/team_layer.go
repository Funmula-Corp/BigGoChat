// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchlayer

import (
	model "git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	store "git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
)

type SearchTeamStore struct {
	store.TeamStore
	rootStore *SearchStore
}

func (s SearchTeamStore) AnalyticsTeamCount(opts *model.TeamSearch) (int64, error) {
	var term string
	if opts != nil {
		term = opts.Term
	}
	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			_, total, aErr := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, 0, 25)
			if aErr != nil {
				return 0, aErr
			}
			return total, nil
		}
	}
	return s.TeamStore.AnalyticsTeamCount(opts)
}

func (s SearchTeamStore) GetAllPage(offset int, limit int, opts *model.TeamSearch) ([]*model.Team, error) {
	var (
		page    int = offset / limit
		perPage int = limit
		term    string
	)
	if opts != nil {
		term = opts.Term
		if opts.Page != nil {
			page = *opts.Page
		}
		if opts.PerPage != nil {
			perPage = *opts.PerPage
		}
	}

	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			ids, _, err := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, page, perPage)
			if err != nil {
				return nil, err
			}
			return s.TeamStore.GetMany(ids)
		}
	}
	return s.TeamStore.GetAllPage(offset, limit, opts)
}

func (s SearchTeamStore) SearchAll(opts *model.TeamSearch) ([]*model.Team, error) {
	var (
		page    int = 0
		perPage int = 25
		term    string
	)
	if opts != nil {
		term = opts.Term
		if opts.Page != nil {
			page = *opts.Page
		}
		if opts.PerPage != nil {
			perPage = *opts.PerPage
		}
	}

	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			ids, _, err := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, page, perPage)
			if err != nil {
				return nil, err
			}
			return s.TeamStore.GetMany(ids)
		}
	}
	return s.TeamStore.SearchAll(opts)
}

func (s SearchTeamStore) SearchAllPaged(opts *model.TeamSearch) ([]*model.Team, int64, error) {
	var (
		page    int = 0
		perPage int = 25
		term    string
	)
	if opts != nil {
		term = opts.Term
		if opts.Page != nil {
			page = *opts.Page
		}
		if opts.PerPage != nil {
			perPage = *opts.PerPage
		}
	}

	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			ids, total, aErr := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, page, perPage)
			if aErr != nil {
				return nil, 0, aErr
			}
			teams, err := s.TeamStore.GetMany(ids)
			return teams, total, err
		}
	}
	return s.TeamStore.SearchAllPaged(opts)
}

func (s SearchTeamStore) SearchOpen(opts *model.TeamSearch) ([]*model.Team, error) {
	var (
		page    int = 0
		perPage int = 25
		term    string
	)
	if opts != nil {
		term = opts.Term
		if opts.Page != nil {
			page = *opts.Page
		}
		if opts.PerPage != nil {
			perPage = *opts.PerPage
		}
	}

	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			ids, _, aErr := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, page, perPage)
			if aErr != nil {
				return nil, aErr
			}
			teams, err := s.TeamStore.GetMany(ids)
			return teams, err
		}
	}
	return s.TeamStore.SearchOpen(opts)
}

func (s SearchTeamStore) SearchPrivate(opts *model.TeamSearch) ([]*model.Team, error) {
	var (
		page    int = 0
		perPage int = 25
		term    string
	)
	if opts != nil {
		term = opts.Term
		if opts.Page != nil {
			page = *opts.Page
		}
		if opts.PerPage != nil {
			perPage = *opts.PerPage
		}
	}

	for _, engine := range s.rootStore.searchEngine.GetActiveEngines() {
		if engine.IsSearchEnabled() {
			ids, _, aErr := engine.SearchTeams("", []*model.SearchParams{{
				Terms: term,
			}}, page, perPage)
			if aErr != nil {
				return nil, aErr
			}
			teams, err := s.TeamStore.GetMany(ids)
			return teams, err
		}
	}
	return s.TeamStore.SearchPrivate(opts)
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
