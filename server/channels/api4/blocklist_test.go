package api4

import (
	"context"
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserBlockUser(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client
	ubul0, resp0, err0 := client.ListUserBlockUsers(context.Background(), th.BasicUser.Id)
	require.NoError(t, err0)
	CheckOKStatus(t, resp0)
	assert.Len(t, *ubul0, 0)

	ubu1, resp1, err1 := client.AddUserBlockUser(context.Background(), th.BasicUser.Id, th.BasicUser2.Id)
	require.NoError(t, err1)
	CheckCreatedStatus(t, resp1)
	assert.Equal(t, ubu1.BlockedId, th.BasicUser2.Id)
	assert.Equal(t, ubu1.UserId, th.BasicUser.Id)

	ubul2, resp2, err2 := client.ListUserBlockUsers(context.Background(), th.BasicUser.Id)
	require.NoError(t, err2)
	CheckOKStatus(t, resp2)
	assert.Len(t, *ubul2, 1)
	assert.Equal(t, *(*ubul2)[0], *ubu1)

	_, resp3, err3 := client.DeleteUserBlockUser(context.Background(), th.BasicUser.Id, th.BasicUser2.Id)
	require.NoError(t, err3)
	CheckOKStatus(t, resp3)

	ubul4, resp4, err4 := client.ListUserBlockUsers(context.Background(), th.BasicUser.Id)
	require.NoError(t, err4)
	CheckOKStatus(t, resp4)
	assert.Len(t, *ubul4, 0)

	_, resp5, err5 := client.AddUserBlockUser(context.Background(), th.BasicUser2.Id, th.BasicUser.Id)
	require.Error(t, err5)
	CheckForbiddenStatus(t, resp5)

	_, resp6, err6 := client.DeleteUserBlockUser(context.Background(), th.BasicUser2.Id, th.BasicUser.Id)
	require.Error(t, err6)
	CheckForbiddenStatus(t, resp6)

	_, resp7, err7 := client.ListUserBlockUsers(context.Background(), th.BasicUser2.Id)
	require.Error(t, err7)
	CheckForbiddenStatus(t, resp7)
}

func TestUserBlockUserPost(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	ctx := context.Background()
	client := th.Client
	client2 := th.CreateClient()
	_, _, lErr := client2.Login(context.Background(), th.BasicUser2.Username, th.BasicUser2.Password)
	if lErr != nil {
		panic(lErr)
	}
	dmChannel, resp, err := client.CreateDirectChannel(ctx, th.BasicUser.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, dmChannel)

	post := &model.Post{ChannelId: dmChannel.Id, Message: "msg1"}
	_, resp2, err2 := client.CreatePost(context.Background(), post)
	require.NoError(t, err2)
	CheckCreatedStatus(t, resp2)

	// a block b
	_, resp, err = client.AddUserBlockUser(context.Background(), th.BasicUser.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg2"}
	_, resp, err = client.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg3 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	// TODO: test they can still read posts in the dm channel

	// b blocked a. they block each other
	_, resp, err = client2.AddUserBlockUser(context.Background(), th.BasicUser2.Id, th.BasicUser.Id)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg4"}
	_, resp, err = client.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg5 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	// a unblock b
	_, resp, err = client.DeleteUserBlockUser(context.Background(), th.BasicUser.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	CheckOKStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg6"}
	_, resp, err = client.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg7 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	// b unblock a, they unblocked each other
	_, resp, err = client2.DeleteUserBlockUser(context.Background(), th.BasicUser2.Id, th.BasicUser.Id)
	require.NoError(t, err)
	CheckOKStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg6"}
	_, resp, err = client.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: dmChannel.Id, Message: "msg7 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)
}

func TestChannelBlockUser(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client
	cbul0, resp0, err0 := client.ListChannelBlockUsers(context.Background(), th.BasicChannel.Id)
	require.NoError(t, err0)
	CheckOKStatus(t, resp0)
	assert.Len(t, *cbul0, 0)

	cbu1, resp1, err1 := client.AddChannelBlockUser(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.NoError(t, err1)
	CheckCreatedStatus(t, resp1)
	assert.Equal(t, cbu1.BlockedId, th.BasicUser2.Id)
	assert.Equal(t, cbu1.ChannelId, th.BasicChannel.Id)
	assert.Equal(t, cbu1.CreateBy, th.BasicUser.Id)

	cbul2, resp2, err2 := client.ListChannelBlockUsers(context.Background(), th.BasicChannel.Id)
	require.NoError(t, err2)
	CheckOKStatus(t, resp2)
	assert.Len(t, *cbul2, 1)
	assert.Equal(t, *(*cbul2)[0], *cbu1)

	_, resp3, err3 := client.DeleteChannelBlockUser(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.NoError(t, err3)
	CheckOKStatus(t, resp3)

	cbul4, resp4, err4 := client.ListChannelBlockUsers(context.Background(), th.BasicChannel.Id)
	require.NoError(t, err4)
	CheckOKStatus(t, resp4)
	assert.Len(t, *cbul4, 0)

	// checking permissions
	// th.LoginBasic2()
	otherClient := th.CreateClient()
	otherUser := th.CreateUserWithClient(otherClient)
	th.LinkUserToTeam(otherUser, th.BasicTeam)
	_, _, lErr := otherClient.Login(context.Background(), otherUser.Username, otherUser.Password)
	if lErr != nil {
		panic(lErr)
	}

	_, resp5, err5 := otherClient.AddChannelBlockUser(context.Background(), th.BasicDeletedChannel.Id, th.BasicUser2.Id)
	require.Error(t, err5)
	CheckForbiddenStatus(t, resp5)

	_, resp6, err6 := otherClient.DeleteChannelBlockUser(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.Error(t, err6)
	CheckForbiddenStatus(t, resp6)

	_, resp7, err7 := otherClient.ListUserBlockUsers(context.Background(), th.BasicChannel.Id)
	require.Error(t, err7)
	CheckForbiddenStatus(t, resp7)
}
