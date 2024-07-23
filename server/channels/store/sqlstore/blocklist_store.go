// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"

	sq "github.com/mattermost/squirrel"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
	"github.com/pkg/errors"
)

type SqlBlocklistStore struct {
	*SqlStore
}

func newSqlChannelBlocklistStore(sqlStore *SqlStore) store.BlocklistStore {
	return &SqlBlocklistStore{sqlStore}
}

func NewMapFromChannelBlockUserModel(cm *model.ChannelBlockUser) map[string]any {
	return map[string]any{
		"ChannelId": cm.ChannelId,
		"BlockedId": cm.BlockedId,
		"CreateBy":  cm.CreateBy,
		"CreateAt":  cm.CreateAt,
	}
}

func teamBlockUserSliceColumns() []string {
	return []string{"TeamId", "BlockedId", "CreateBy", "CreateAt"}
}

func teamBlockUserToSlice(teamBlockUser *model.TeamBlockUser) []any {
	resultSlice := []any{}
	resultSlice = append(resultSlice, teamBlockUser.TeamId)
	resultSlice = append(resultSlice, teamBlockUser.BlockedId)
	resultSlice = append(resultSlice, teamBlockUser.CreateBy)
	resultSlice = append(resultSlice, teamBlockUser.CreateAt)
	return resultSlice
}


func channelBlockUserSliceColumns() []string {
	return []string{"ChannelId", "BlockedId", "CreateBy", "CreateAt"}
}

func channelBlockUserToSlice(channelBlockUser *model.ChannelBlockUser) []any {
	resultSlice := []any{}
	resultSlice = append(resultSlice, channelBlockUser.ChannelId)
	resultSlice = append(resultSlice, channelBlockUser.BlockedId)
	resultSlice = append(resultSlice, channelBlockUser.CreateBy)
	resultSlice = append(resultSlice, channelBlockUser.CreateAt)
	return resultSlice
}

func userBlockUserSliceColumns() []string {
	return []string{"UserId", "BlockedId", "CreateAt"}
}

func userBlockUserToSlice(userBlockUser *model.UserBlockUser) []any {
	resultSlice := []any{}
	resultSlice = append(resultSlice, userBlockUser.UserId)
	resultSlice = append(resultSlice, userBlockUser.BlockedId)
	resultSlice = append(resultSlice, userBlockUser.CreateAt)
	return resultSlice
}

func (s *SqlBlocklistStore) GetTeamBlockUser(teamId string, blockedId string) (*model.TeamBlockUser, error) {
	query := s.getQueryBuilder().
		Select(teamBlockUserSliceColumns()...).
		From("TeamBlockUsers tb").Where(
		sq.And{
			sq.Eq{"tb.TeamId": teamId},
			sq.Eq{"tb.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_team_block_user_tosql")
	}

	teamBlockUser := model.TeamBlockUser{}
	err = s.GetReplicaX().Get(&teamBlockUser, queryString, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "failed to get teams blocklist with id %v", teamId)
		}
	}
	return &teamBlockUser, nil
}

func (s *SqlBlocklistStore) GetTeamBlockUserByEmail(teamId string, email string) (*model.TeamBlockUser, error) {
	query := s.getQueryBuilder().
		Select("TeamBlockUsers.TeamId", "TeamBlockUsers.BlockedId", "TeamBlockUsers.CreateBy", "TeamBlockUsers.CreateAt").
		From("TeamBlockUsers").
		InnerJoin("Users ON TeamBlockUsers.BlockedId = Users.Id").
		Where(
		sq.And{
			sq.Eq{"TeamBlockUsers.TeamId": teamId},
			sq.Eq{"Users.Email": email},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_team_block_user_tosql")
	}

	teamBlockUser := model.TeamBlockUser{}
	err = s.GetReplicaX().Get(&teamBlockUser, queryString, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "failed to get teams blocklist with id %v", teamId)
		}
	}
	return &teamBlockUser, nil
}

func (s *SqlBlocklistStore) ListTeamBlockUsers(teamId string) (*model.TeamBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(teamBlockUserSliceColumns()...).
		From("TeamBlockUsers tb").Where(
		sq.And{
			sq.Eq{"tb.TeamId": teamId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_team_block_users_tosql")
	}

	teamBlockUserList := model.TeamBlockUserList{}
	err = s.GetReplicaX().Select(&teamBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get teams blocklist with id %v", teamId)
	}
	return &teamBlockUserList, nil
}

func (s *SqlBlocklistStore) ListTeamBlockUsersByBlockedUser(blockedId string) (*model.TeamBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(teamBlockUserSliceColumns()...).
		From("TeamBlockUsers tb").Where(
		sq.And{
			sq.Eq{"tb.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_team_block_users_by_blocked_user_tosql")
	}

	teamBlockUserList := model.TeamBlockUserList{}
	err = s.GetReplicaX().Select(&teamBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get teams blocklist by blocked user_id %v", blockedId)
	}
	return &teamBlockUserList, nil
}

func (s *SqlBlocklistStore) DeleteTeamBlockUser(teamId string, userId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "SetDeleteAt: begin_transaction")
	}
	defer finalizeTransactionX(transaction, &err)

	if _, err := transaction.Exec(`DELETE FROM TeamBlockUsers WHERE TeamId = ? and BlockedId = ?`, teamId, userId); err != nil {
		return errors.Wrapf(err, "failed to delete team block user %s:%s", teamId, userId)
	}
	if err := transaction.Commit(); err != nil {
		return errors.Wrapf(err, "Delete: commit_transaction")
	}
	return nil
}

func (s *SqlBlocklistStore) SaveTeamBlockUser(blockUser *model.TeamBlockUser) (*model.TeamBlockUser, error) {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	blockUser.PreSave()
	defer finalizeTransactionX(transaction, &err)
	query := s.getQueryBuilder().Insert("TeamBlockUsers").Columns(teamBlockUserSliceColumns()...).Values(teamBlockUserToSlice(blockUser)...)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "save_team_block_user_tosql")
	}

	if _, err := transaction.Exec(sql, args...); err != nil {
		if IsUniqueConstraintError(err, []string{"Name", "team_block_users_key"}) {
			dup := model.TeamBlockUser{}
			if serr := s.GetMasterX().Get(&dup, "SELECT * FROM TeamBlockUsers WHERE TeamId = ? AND BlockedId = ?", blockUser.TeamId, blockUser.BlockedId); serr != nil {
				return nil, errors.Wrapf(serr, "error while retrieving existing team block user %s %s", blockUser.TeamId, blockUser.BlockedId)
			}
			return &dup, store.NewErrConflict("TeamBlockUser", err, "id="+blockUser.TeamId)
		}
		return nil, errors.Wrapf(err, "save_team_block_user: team_id=%s blocked_id=%s", blockUser.TeamId, blockUser.BlockedId)
	}
	err = transaction.Commit()
	return blockUser, err
}

// Channel Block User
func (s *SqlBlocklistStore) GetChannelBlockUser(channelId string, blockedId string) (*model.ChannelBlockUser, error) {
	query := s.getQueryBuilder().
		Select(channelBlockUserSliceColumns()...).
		From("ChannelBlockUsers cb").Where(
		sq.And{
			sq.Eq{"cb.ChannelId": channelId},
			sq.Eq{"cb.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_channel_block_user_tosql")
	}

	channelBlockUser := model.ChannelBlockUser{}
	err = s.GetReplicaX().Get(&channelBlockUser, queryString, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "failed to get channels blocklist with id %v", channelId)
		}
	}
	return &channelBlockUser, nil
}

func (s *SqlBlocklistStore) GetChannelBlockUserByEmail(channelId string, email string) (*model.ChannelBlockUser, error) {
	query := s.getQueryBuilder().
		Select("ChannelBlockUsers.ChannelId", "ChannelBlockUsers.BlockedId", "ChannelBlockUsers.CreateBy", "ChannelBlockUsers.CreateAt").
		From("ChannelBlockUsers").
		InnerJoin("Users ON ChannelBlockUsers.BlockedId = Users.Id").
		Where(
		sq.And{
			sq.Eq{"ChannelBlockUsers.ChannelId": channelId},
			sq.Eq{"Users.Email": email},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_channel_block_user_tosql")
	}

	channelBlockUser := model.ChannelBlockUser{}
	err = s.GetReplicaX().Get(&channelBlockUser, queryString, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "failed to get channels blocklist with id %v %s", channelId, queryString)
		}
	}
	return &channelBlockUser, nil
}

func (s *SqlBlocklistStore) ListChannelBlockUsers(channelId string) (*model.ChannelBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(channelBlockUserSliceColumns()...).
		From("ChannelBlockUsers cb").Where(
		sq.And{
			sq.Eq{"cb.ChannelId": channelId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_channel_block_users_tosql")
	}

	channelBlockUserList := model.ChannelBlockUserList{}
	err = s.GetReplicaX().Select(&channelBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get channels blocklist with id %v", channelId)
	}
	return &channelBlockUserList, nil
}

func (s *SqlBlocklistStore) ListChannelBlockUsersByBlockedUser(blockedId string) (*model.ChannelBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(channelBlockUserSliceColumns()...).
		From("ChannelBlockUsers cb").Where(
		sq.And{
			sq.Eq{"cb.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_channel_block_users_by_blocked_user_tosql")
	}

	channelBlockUserList := model.ChannelBlockUserList{}
	err = s.GetReplicaX().Select(&channelBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get channels blocklist by blocked user_id %v", blockedId)
	}
	return &channelBlockUserList, nil
}

func (s *SqlBlocklistStore) DeleteChannelBlockUser(channelId string, userId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "SetDeleteAt: begin_transaction")
	}
	defer finalizeTransactionX(transaction, &err)

	if _, err := transaction.Exec(`DELETE FROM ChannelBlockUsers WHERE ChannelId = ? and BlockedId = ?`, channelId, userId); err != nil {
		return errors.Wrapf(err, "failed to delete channel block user %s:%s", channelId, userId)
	}
	if err := transaction.Commit(); err != nil {
		return errors.Wrapf(err, "Delete: commit_transaction")
	}
	return nil
}

func (s *SqlBlocklistStore) SaveChannelBlockUser(blockUser *model.ChannelBlockUser) (*model.ChannelBlockUser, error) {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	blockUser.PreSave()
	defer finalizeTransactionX(transaction, &err)
	query := s.getQueryBuilder().Insert("ChannelBlockUsers").Columns(channelBlockUserSliceColumns()...).Values(channelBlockUserToSlice(blockUser)...)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "save_channel_block_user_tosql")
	}

	if _, err := transaction.Exec(`DELETE FROM ChannelMembers WHERE ChannelId = ? and UserId = ?`, blockUser.ChannelId, blockUser.BlockedId); err != nil {
		return nil, errors.Wrapf(err, "failed to delete blocked user %s from channel %s", blockUser.BlockedId, blockUser.ChannelId)
	}
	if _, err := transaction.Exec(sql, args...); err != nil {
		if IsUniqueConstraintError(err, []string{"Name", "channel_block_users_key"}) {
			dup := model.ChannelBlockUser{}
			if serr := s.GetMasterX().Get(&dup, "SELECT * FROM ChannelBlockUsers WHERE ChannelId = ? AND BlockedId = ?", blockUser.ChannelId, blockUser.BlockedId); serr != nil {
				return nil, errors.Wrapf(serr, "error while retrieving existing channel block user %s %s", blockUser.ChannelId, blockUser.BlockedId)
			}
			return &dup, store.NewErrConflict("ChannelBlockUser", err, "id="+blockUser.ChannelId)
		}
		return nil, errors.Wrapf(err, "save_channel_block_user: channel_id=%s blocked_id=%s", blockUser.ChannelId, blockUser.BlockedId)
	}
	err = transaction.Commit()
	return blockUser, err
}

func (s *SqlBlocklistStore) GetUserBlockUser(userId string, blockedId string) (*model.UserBlockUser, error) {
	query := s.getQueryBuilder().
		Select(userBlockUserSliceColumns()...).
		From("UserBlockUsers ub").Where(
		sq.And{
			sq.Eq{"ub.UserId": userId},
			sq.Eq{"ub.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "get_user_block_user_tosql")
	}

	userBlockUser := model.UserBlockUser{}
	err = s.GetReplicaX().Get(&userBlockUser, queryString, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "failed to get user blocklist with id %v %v", userId, blockedId)
		}
	}
	return &userBlockUser, nil
}

func (s *SqlBlocklistStore) ListUserBlockUsers(userId string) (*model.UserBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(userBlockUserSliceColumns()...).
		From("UserBlockUsers ub").Where(
		sq.And{
			sq.Eq{"ub.UserId": userId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_user_blockusers_tosql")
	}

	userBlockUserList := model.UserBlockUserList{}
	err = s.GetReplicaX().Select(&userBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users blocklist with id %v", userId)
	}
	return &userBlockUserList, nil
}

func (s *SqlBlocklistStore) ListUserBlockUsersByBlockedUser(blockedId string) (*model.UserBlockUserList, error) {
	query := s.getQueryBuilder().
		Select(userBlockUserSliceColumns()...).
		From("UserBlockUsers ub").Where(
		sq.And{
			sq.Eq{"ub.BlockedId": blockedId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "list_user_block_users_by_blocked_user_tosql")
	}

	userBlockUserList := model.UserBlockUserList{}
	err = s.GetReplicaX().Select(&userBlockUserList, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list user block users by blocked user id %v", blockedId)
	}
	return &userBlockUserList, nil
}

func (s *SqlBlocklistStore) DeleteUserBlockUser(userId, blockedId string, userIsVerified, blockedIsVerified bool) error {
	var blockedByPeer model.UserBlockUser

	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "SetDeleteAt: begin_transaction")
	}
	defer finalizeTransactionX(transaction, &err)

	deleteResult, err := transaction.Exec(`DELETE FROM UserBlockUsers WHERE UserId = ? and BlockedId = ?`, userId, blockedId)
	if err != nil {
		return errors.Wrapf(err, "failed to delete user block user %s %s", userId, blockedId)
	}
	if r, err := deleteResult.RowsAffected(); err != nil {
		return errors.Wrapf(err, "failed to count delete user block user rows %s %s", userId, blockedId)
	} else if r == 0 {
		return nil
	}

	err = transaction.Get(&blockedByPeer, `SELECT UserId, BlockedId, CreateAt FROM UserBlockUsers WHERE UserId = ? and BlockedId = ? FOR SHARE`, blockedId, userId)

	if err == sql.ErrNoRows {
		if _, err = transaction.Exec("UPDATE ChannelMembers SET SchemeVerified = ?, SchemeUser=true WHERE ChannelId IN (SELECT Id FROM Channels WHERE Name= ?) AND UserId = ?", userIsVerified, model.GetDMNameFromIds(userId, blockedId), userId); err != nil {
			return errors.Wrapf(err, "unmark_dm_readonly: user_id=%s blocked_id %s", userId, blockedId)
		}
		if _, err = transaction.Exec("UPDATE ChannelMembers SET SchemeVerified = ?, SchemeUser=true WHERE ChannelId IN (SELECT Id FROM Channels WHERE Name= ?) AND UserId = ?", blockedIsVerified, model.GetDMNameFromIds(userId, blockedId), blockedId); err != nil {
			return errors.Wrapf(err, "unmark_dm_readonly: user_id=%s blocked_id %s", userId, blockedId)
		}
	} else if err != nil {
		return errors.Wrapf(err, "list: user_id=%s blocked_id %s", userId, blockedId)
	}

	if err := transaction.Commit(); err != nil {
		return errors.Wrapf(err, "Delete: commit_transaction")
	}
	return nil
}

func (s *SqlBlocklistStore) SaveUserBlockUser(userBlockUser *model.UserBlockUser) (*model.UserBlockUser, error) {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "begin_transaction")
	}
	userBlockUser.PreSave()
	defer finalizeTransactionX(transaction, &err)

	// this hold "DELETE FROM UserBlockUsers" query from the other user.
	_, err = transaction.Exec(
		"SELECT * FROM UserBlockUsers WHERE UserId = ? and BlockedId = ? FOR SHARE",
		userBlockUser.BlockedId, userBlockUser.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "error while select user block user: %s %s", userBlockUser.BlockedId, userBlockUser.UserId)
	}

	query := s.getQueryBuilder().Insert("UserBlockUsers").Columns(userBlockUserSliceColumns()...).Values(userBlockUserToSlice(userBlockUser)...)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "save_user_block_user_tosql")
	}
	if _, err = transaction.Exec(sql, args...); err != nil {
		if IsUniqueConstraintError(err, []string{"Name", "user_block_users_key"}) {
			dup := model.UserBlockUser{}
			if serr := s.GetMasterX().Get(&dup, "SELECT * FROM UserBlockUsers WHERE UserId = ? AND BlockedId = ?", userBlockUser.UserId, userBlockUser.BlockedId); serr != nil {
				return nil, errors.Wrapf(serr, "error while retrieving existing user block user %s %s", userBlockUser.UserId, userBlockUser.BlockedId)
			}
			return &dup, store.NewErrConflict("UserBlockUser", err, "id="+userBlockUser.UserId+":"+userBlockUser.BlockedId)
		}
		return nil, errors.Wrapf(err, "save_user_block_user: user_id=%s blocked_id=%s", userBlockUser.UserId, userBlockUser.BlockedId)
	}

	if _, err = transaction.Exec("UPDATE ChannelMembers SET SchemeVerified=false, SchemeUser=false, Roles='channel_readonly' WHERE ChannelId IN (SELECT Id FROM Channels WHERE Name= ?)", userBlockUser.GetDMName()); err != nil {
		return nil, errors.Wrapf(err, "mark_dm_readonly: user_id=%s blocked_id %s", userBlockUser.UserId, userBlockUser.BlockedId)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, errors.Wrapf(err, "SaveUserBlockUser: commit_transaction")
	}
	return userBlockUser, nil
}
