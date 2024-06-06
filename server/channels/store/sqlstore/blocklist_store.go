// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
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
		return nil, errors.Wrapf(err, "failed to get channels blocklist with id %v", channelId)
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
		return nil, errors.Wrapf(err, "failed to get user blocklist with id %v %v", userId, blockedId)
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

func (s *SqlBlocklistStore) DeleteUserBlockUser(userId string, blockedId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "SetDeleteAt: begin_transaction")
	}
	defer finalizeTransactionX(transaction, &err)

	if _, err := transaction.Exec(`DELETE FROM UserBlockUsers WHERE UserId = ? and BlockedId = ?`, userId, blockedId); err != nil {
		return errors.Wrapf(err, "failed to delete channel block user %s %s", userId, blockedId)
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
	query := s.getQueryBuilder().Insert("UserBlockUsers").Columns(userBlockUserSliceColumns()...).Values(userBlockUserToSlice(userBlockUser)...)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "save_user_block_user_tosql")
	}

	if _, err := transaction.Exec(sql, args...); err != nil {
		if IsUniqueConstraintError(err, []string{"Name", "user_block_users_key"}) {
			dup := model.UserBlockUser{}
			if serr := s.GetMasterX().Get(&dup, "SELECT * FROM UserBlockUsers WHERE UserId = ? AND BlockedId = ?", userBlockUser.UserId, userBlockUser.BlockedId); serr != nil {
				return nil, errors.Wrapf(serr, "error while retrieving existing user block user %s %s", userBlockUser.UserId, userBlockUser.BlockedId)
			}
			return &dup, store.NewErrConflict("UserBlockUser", err, "id="+userBlockUser.UserId+":"+userBlockUser.BlockedId)
		}
		return nil, errors.Wrapf(err, "save_channel_block_user: user_id=%s blocked_id=%s", userBlockUser.UserId, userBlockUser.BlockedId)
	}
	err = transaction.Commit()
	return userBlockUser, err
}
