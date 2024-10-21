package storetest

import (
	"testing"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocklistStore(t *testing.T, rctx request.CTX, ss store.Store) {
	t.Run("SaveTeamBlockUser", func(t *testing.T) { testSaveTeamBlockUser(t, rctx, ss) })
	t.Run("SaveChannelBlockUser", func(t *testing.T) { testSaveChannelBlockUser(t, rctx, ss) })
	t.Run("SaveUserBlockUser", func(t *testing.T) { testSaveUserBlockUser(t, rctx, ss) })
	t.Run("SaveUserBlockUserDM", func(t *testing.T) { testSaveUserBlockUserDM(t, rctx, ss) })
	t.Run("ListTeamBlockUser", func(t *testing.T) { testListTeamBlockUser(t, rctx, ss) })
	t.Run("ListChannelBlockUser", func(t *testing.T) { testListChannelBlockUser(t, rctx, ss) })
	t.Run("ListUserBlockUser", func(t *testing.T) { testListUserBlockUser(t, rctx, ss) })
	t.Run("GetTeamBlockUserByEmail", func(t *testing.T) { testGetTeamBlockUserByEmail(t, rctx, ss) })
	t.Run("GetChannelBlockUserByEmail", func(t *testing.T) { testGetChannelBlockUserByEmail(t, rctx, ss) })
}

func testSaveTeamBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	teamBlockUser := model.TeamBlockUser{
		TeamId: "000000",
		BlockedId: "abcdefg",
		CreateBy:  "abcdfg",
	}

	var err error
	var newBlockUser, getBlockUser *model.TeamBlockUser
	newBlockUser, err = ss.Blocklist().SaveTeamBlockUser(&teamBlockUser)
	require.NoError(t, err)
	require.Greater(t, newBlockUser.CreateAt, model.GetMillis()-200)
	assert.Equal(t, teamBlockUser.BlockedId, newBlockUser.BlockedId)
	assert.Equal(t, teamBlockUser.TeamId, newBlockUser.TeamId)
	assert.Equal(t, teamBlockUser.CreateBy, newBlockUser.CreateBy)

	getBlockUser, err = ss.Blocklist().GetTeamBlockUser(teamBlockUser.TeamId, teamBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Equal(t, teamBlockUser.BlockedId, getBlockUser.BlockedId)
	assert.Equal(t, teamBlockUser.TeamId, getBlockUser.TeamId)
	assert.Equal(t, newBlockUser.CreateBy, getBlockUser.CreateBy)

	err = ss.Blocklist().DeleteTeamBlockUser(teamBlockUser.TeamId, teamBlockUser.BlockedId)
	require.NoError(t, err)

	getBlockUser, err = ss.Blocklist().GetTeamBlockUser(teamBlockUser.TeamId, teamBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Nil(t, getBlockUser)

	err = ss.Blocklist().DeleteTeamBlockUser(teamBlockUser.TeamId, teamBlockUser.BlockedId)
	require.NoError(t, err)
}

func testSaveChannelBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	channelBlockUser := model.ChannelBlockUser{
		ChannelId: "000000",
		BlockedId: "abcdefg",
		CreateBy:  "abcdfg",
	}

	var err error
	var newBlockUser, getBlockUser *model.ChannelBlockUser
	newBlockUser, err = ss.Blocklist().SaveChannelBlockUser(&channelBlockUser)
	require.NoError(t, err)
	require.Greater(t, newBlockUser.CreateAt, model.GetMillis()-200)
	assert.Equal(t, channelBlockUser.BlockedId, newBlockUser.BlockedId)
	assert.Equal(t, channelBlockUser.ChannelId, newBlockUser.ChannelId)
	assert.Equal(t, channelBlockUser.CreateBy, newBlockUser.CreateBy)

	getBlockUser, err = ss.Blocklist().GetChannelBlockUser(channelBlockUser.ChannelId, channelBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Equal(t, channelBlockUser.BlockedId, getBlockUser.BlockedId)
	assert.Equal(t, channelBlockUser.ChannelId, getBlockUser.ChannelId)
	assert.Equal(t, newBlockUser.CreateBy, getBlockUser.CreateBy)

	err = ss.Blocklist().DeleteChannelBlockUser(channelBlockUser.ChannelId, channelBlockUser.BlockedId)
	require.NoError(t, err)

	getBlockUser, err = ss.Blocklist().GetChannelBlockUser(channelBlockUser.ChannelId, channelBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Nil(t, getBlockUser)

	err = ss.Blocklist().DeleteChannelBlockUser(channelBlockUser.ChannelId, channelBlockUser.BlockedId)
	require.NoError(t, err)
}

func testSaveUserBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	userBlockUser := model.UserBlockUser{
		UserId:    "000000",
		BlockedId: "abcdefg",
	}

	var err error
	var newBlockUser, getBlockUser *model.UserBlockUser
	newBlockUser, err = ss.Blocklist().SaveUserBlockUser(&userBlockUser)
	require.NoError(t, err)
	require.Greater(t, newBlockUser.CreateAt, model.GetMillis()-200)
	assert.Equal(t, userBlockUser.BlockedId, newBlockUser.BlockedId)
	assert.Equal(t, userBlockUser.UserId, newBlockUser.UserId)

	getBlockUser, err = ss.Blocklist().GetUserBlockUser(userBlockUser.UserId, userBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Equal(t, userBlockUser.BlockedId, getBlockUser.BlockedId)
	assert.Equal(t, userBlockUser.UserId, getBlockUser.UserId)
	assert.Equal(t, newBlockUser.CreateAt, getBlockUser.CreateAt)

	err = ss.Blocklist().DeleteUserBlockUser(userBlockUser.UserId, userBlockUser.BlockedId, true, true)
	require.NoError(t, err)

	getBlockUser, err = ss.Blocklist().GetUserBlockUser(userBlockUser.UserId, userBlockUser.BlockedId)
	require.NoError(t, err)
	assert.Nil(t, getBlockUser)

	err = ss.Blocklist().DeleteUserBlockUser(userBlockUser.UserId, userBlockUser.BlockedId, true, true)
	require.NoError(t, err)
}

func testSaveUserBlockUserDM(t *testing.T, ctx request.CTX, ss store.Store) {
	userA := "biggouyyyyyyyyyyyyyyyyyyyy"
	userB := "biggouyyyyyyyyyyyyyyyyyyyb"
	channel := model.Channel{
		TeamId: "teamid",
		Type: model.ChannelTypeDirect,
		Name: model.GetDMNameFromIds(userA, userB),
	}
	cmA := model.ChannelMember{
		UserId: userA,
		SchemeUser: true,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	cmB := model.ChannelMember{
		UserId: userB,
		SchemeUser: true,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	newChannel, err := ss.Channel().SaveDirectChannel(ctx, &channel, &cmA, &cmB)
	require.NoError(t, err)

	aBlockB := model.UserBlockUser{
		UserId:    userA,
		BlockedId: userB,
	}
	_, err = ss.Blocklist().SaveUserBlockUser(&aBlockB)
	require.NoError(t, err)

	channelMembers, err := ss.Channel().GetMembersByIds(newChannel.Id, []string{userA, userB})
	require.NoError(t, err)
	for _, cm := range(channelMembers){
		assert.False(t, cm.SchemeUser)
	}

	bBlockA := model.UserBlockUser{
		UserId:    userB,
		BlockedId: userA,
	}
	_, err = ss.Blocklist().SaveUserBlockUser(&bBlockA)
	require.NoError(t, err)
	channelMembers, err = ss.Channel().GetMembersByIds(newChannel.Id, []string{userA, userB})
	require.NoError(t, err)
	for _, cm := range(channelMembers){
		assert.False(t, cm.SchemeUser)
	}

	err = ss.Blocklist().DeleteUserBlockUser(userB, userA, true, true)
	require.NoError(t, err)
	channelMembers, err = ss.Channel().GetMembersByIds(newChannel.Id, []string{userA, userB})
	require.NoError(t, err)
	for _, cm := range(channelMembers){
		assert.False(t, cm.SchemeUser)
	}

	err = ss.Blocklist().DeleteUserBlockUser(userA, userB, true, true)
	require.NoError(t, err)
	channelMembers, err = ss.Channel().GetMembersByIds(newChannel.Id, []string{userA, userB})
	require.NoError(t, err)
	for _, cm := range(channelMembers){
		assert.True(t, cm.SchemeUser)
	}
}

func testListTeamBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	user1 := "biggouyyyyyyyyyyyyyyyyyyyy"
	user2 := "biggouyyyyyyyyyyyyyyyyyyyb"
	user3 := "biggouyyyyyyyyyyyyyyyyyyyn"
	user4 := "biggouyyyyyyyyyyyyyyyyyyyd"
	Team1 := "biggocyyyyyyyyyyyyyyyyyyyy"
	Team2 := "biggocyyyyyyyyyyyyyyyyyyyb"
	Team3 := "biggocyyyyyyyyyyyyyyyyyyyn"
	createAt := model.GetMillis()
	sampleData := model.TeamBlockUserList{}
	sampleData = append(sampleData, &model.TeamBlockUser{TeamId: Team1, BlockedId: user2, CreateBy: user1, CreateAt: createAt})
	sampleData = append(sampleData, &model.TeamBlockUser{TeamId: Team1, BlockedId: user3, CreateBy: user1, CreateAt: createAt})
	sampleData = append(sampleData, &model.TeamBlockUser{TeamId: Team2, BlockedId: user1, CreateBy: user2, CreateAt: createAt})
	sampleData = append(sampleData, &model.TeamBlockUser{TeamId: Team2, BlockedId: user3, CreateBy: user2, CreateAt: createAt})
	for _, TeamBlockUser := range(sampleData) {
		_, err := ss.Blocklist().SaveTeamBlockUser(TeamBlockUser)
		require.NoError(t, err)
	}

	var err error
	var retCBUL *model.TeamBlockUserList
	retCBUL, err = ss.Blocklist().ListTeamBlockUsers(Team1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsers(Team2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsers(Team3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsersByBlockedUser(user1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsersByBlockedUser(user2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsersByBlockedUser(user3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListTeamBlockUsersByBlockedUser(user4)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)
}

func testListChannelBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	user1 := "biggouyyyyyyyyyyyyyyyyyyyy"
	user2 := "biggouyyyyyyyyyyyyyyyyyyyb"
	user3 := "biggouyyyyyyyyyyyyyyyyyyyn"
	user4 := "biggouyyyyyyyyyyyyyyyyyyyd"
	channel1 := "biggocyyyyyyyyyyyyyyyyyyyy"
	channel2 := "biggocyyyyyyyyyyyyyyyyyyyb"
	channel3 := "biggocyyyyyyyyyyyyyyyyyyyn"
	createAt := model.GetMillis()
	sampleData := model.ChannelBlockUserList{}
	sampleData = append(sampleData, &model.ChannelBlockUser{ChannelId: channel1, BlockedId: user2, CreateBy: user1, CreateAt: createAt})
	sampleData = append(sampleData, &model.ChannelBlockUser{ChannelId: channel1, BlockedId: user3, CreateBy: user1, CreateAt: createAt})
	sampleData = append(sampleData, &model.ChannelBlockUser{ChannelId: channel2, BlockedId: user1, CreateBy: user2, CreateAt: createAt})
	sampleData = append(sampleData, &model.ChannelBlockUser{ChannelId: channel2, BlockedId: user3, CreateBy: user2, CreateAt: createAt})
	for _, channelBlockUser := range(sampleData) {
		_, err := ss.Blocklist().SaveChannelBlockUser(channelBlockUser)
		require.NoError(t, err)
	}

	var err error
	var retCBUL *model.ChannelBlockUserList
	retCBUL, err = ss.Blocklist().ListChannelBlockUsers(channel1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsers(channel2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsers(channel3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsersByBlockedUser(user1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsersByBlockedUser(user2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsersByBlockedUser(user3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListChannelBlockUsersByBlockedUser(user4)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)
}

func testListUserBlockUser(t *testing.T, _ request.CTX, ss store.Store) {
	user1 := "biggouyyyyyyyyyyyyyyyyyyyy"
	user2 := "biggouyyyyyyyyyyyyyyyyyyyb"
	user3 := "biggouyyyyyyyyyyyyyyyyyyyn"
	user4 := "biggouyyyyyyyyyyyyyyyyyyyd"
	createAt := model.GetMillis()
	sampleData := model.UserBlockUserList{}
	sampleData = append(sampleData, &model.UserBlockUser{UserId: user1, BlockedId: user2, CreateAt: createAt})
	sampleData = append(sampleData, &model.UserBlockUser{UserId: user1, BlockedId: user3, CreateAt: createAt})
	sampleData = append(sampleData, &model.UserBlockUser{UserId: user2, BlockedId: user1, CreateAt: createAt})
	sampleData = append(sampleData, &model.UserBlockUser{UserId: user2, BlockedId: user3, CreateAt: createAt})
	for _, userBlockUser := range(sampleData) {
		_, err := ss.Blocklist().SaveUserBlockUser(userBlockUser)
		require.NoError(t, err)
	}

	var err error
	var retCBUL *model.UserBlockUserList
	retCBUL, err = ss.Blocklist().ListUserBlockUsers(user1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListUserBlockUsers(user2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListUserBlockUsers(user3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)

	retCBUL, err = ss.Blocklist().ListUserBlockUsersByBlockedUser(user1)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListUserBlockUsersByBlockedUser(user2)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 1)

	retCBUL, err = ss.Blocklist().ListUserBlockUsersByBlockedUser(user3)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 2)

	retCBUL, err = ss.Blocklist().ListUserBlockUsersByBlockedUser(user4)
	require.NoError(t, err)
	assert.Len(t, *retCBUL, 0)
}

func testGetTeamBlockUserByEmail(t *testing.T, ctx request.CTX, ss store.Store) {
	user := &model.User{
		Email: MakeEmail(),
	}
	user, err := ss.User().Save(ctx, user)
	require.NoError(t, err)
	team1 := model.NewId()
	createAt := model.GetMillis()
	_, err = ss.Blocklist().SaveTeamBlockUser( &model.TeamBlockUser{TeamId: team1, BlockedId: user.Id, CreateBy: model.NewId(), CreateAt: createAt})
	require.NoError(t, err)

	tbu, err := ss.Blocklist().GetTeamBlockUserByEmail(team1, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Id, tbu.BlockedId)
}

func testGetChannelBlockUserByEmail(t *testing.T, ctx request.CTX, ss store.Store) {
	user := &model.User{
		Email: MakeEmail(),
	}
	user, err := ss.User().Save(ctx, user)
	require.NoError(t, err)
	channelId := model.NewId()
	createAt := model.GetMillis()
	_, err = ss.Blocklist().SaveChannelBlockUser( &model.ChannelBlockUser{ChannelId: channelId, BlockedId: user.Id, CreateBy: model.NewId(), CreateAt: createAt})

	cbu, err := ss.Blocklist().GetChannelBlockUserByEmail(channelId, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Id, cbu.BlockedId)
}
