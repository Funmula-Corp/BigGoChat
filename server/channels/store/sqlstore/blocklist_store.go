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

func (s *SqlBlocklistStore) GetChannelBlockUser(channelId string, blockId string) (*model.ChannelBlockUser, error) {
	query := s.getQueryBuilder().
		Select(channelBlockUserSliceColumns()...).
		From("ChannelBlocklist cb").Where(
		sq.And{
			sq.Eq{"cb.ChannelId": channelId},
			sq.Eq{"cb.BlockedId": blockId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "channel_blacklist_get_tosql")
	}

	blackUser := model.ChannelBlockUser{}
	err = s.GetReplicaX().Get(&blackUser, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get channels blacklist with id %v", channelId)
	}
	return &blackUser, nil
}

func (s *SqlBlocklistStore) ListChannelBlockUsers(channelId string) (*model.ChannelBlockUsers, error) {
	query := s.getQueryBuilder().
		Select(channelBlockUserSliceColumns()...).
		From("ChannelBlocklist cb").Where(
		sq.And{
			sq.Eq{"cb.ChannelId": channelId},
		},
	)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "channel_blackusers_get_tosql")
	}

	blackUsers := model.ChannelBlockUsers{}
	err = s.GetReplicaX().Select(&blackUsers, queryString, args...)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get channels blacklist with id %v", channelId)
	}
	return &blackUsers, nil

}

func (s *SqlBlocklistStore) ListChannelBlockUsersByBlockedUser(blockedId string) (*model.ChannelBlockUsers, error) {
	return nil, errors.Errorf("Not Implemented")
}

func (s *SqlBlocklistStore) DeleteChannelBlockUser(channelId string, userId string) error {
	transaction, err := s.GetMasterX().Beginx()
	if err != nil {
		return errors.Wrap(err, "SetDeleteAt: begin_transaction")
	}
	defer finalizeTransactionX(transaction, &err)

	if _, err := transaction.Exec(`DELETE ChannelBlockUsers WHERE ChannelId = ? and BlockedId = ?`, channelId, userId); err != nil {
		return errors.Wrapf(err, "failed to delete channel block user %s %s", channelId, userId)
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
	query := s.getQueryBuilder().Insert("ChannelBlackUsers").Columns(channelBlockUserSliceColumns()...).Values(channelBlockUserToSlice(blockUser)...)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "channel_black_user_tosql")
	}

	if _, err := transaction.Exec(sql, args...); err != nil {
		if IsUniqueConstraintError(err, []string{"Name", "channel_black_users_key"}) {
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
	return nil, errors.Errorf("Not Implemented")
}

func (s *SqlBlocklistStore) ListUserBlockUsers(userId string) (*model.UserBlockUsers, error) {
	return nil, errors.Errorf("Not Implemented")
}
func (s *SqlBlocklistStore) ListUserBlockUsersByBlockedUser(blockedId string) (*model.UserBlockUsers, error) {
	return nil, errors.Errorf("Not Implemented")
}
func (s *SqlBlocklistStore) DeleteUserBlockUser(userId string, blockedId string) error {
	return errors.Errorf("Not Implemented")
}
func (s *SqlBlocklistStore) SaveUserBlockUser(userBlockUser *model.UserBlockUser) (*model.UserBlockUser, error) {
	return nil, errors.Errorf("Not Implemented")
}
