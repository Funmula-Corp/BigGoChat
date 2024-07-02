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

	// they can still read posts in the dm channel
	posts, resp, err := client.GetPostsForChannel(context.Background(), dmChannel.Id, 0, 100, "", false, true)
	require.NoError(t, err)
	CheckOKStatus(t, resp)
	assert.Equal(t, 1, len(posts.Posts))
	assert.False(t, posts.HasNext)

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
	resp, err := th.SystemAdminClient.UpdateChannelMemberSchemeRoles(context.Background(), th.BasicChannel.Id, th.BasicUser.Id, &model.SchemeRoles{
		SchemeAdmin: true,
		SchemeVerified: true,
		SchemeUser: true,
	})
	require.NoError(t, err)
	CheckOKStatus(t, resp)
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

func TestChannelBlockUserPost(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client
	client2 := th.CreateClient()
	_, _, lErr := client2.Login(context.Background(), th.BasicUser2.Username, th.BasicUser2.Password)
	if lErr != nil {
		panic(lErr)
	}

	post := &model.Post{ChannelId: th.BasicChannel.Id, Message: "msg1"}
	_, resp, err := client.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: th.BasicChannel.Id, Message: "msg2 user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	_, resp, err = client.AddChannelBlockUser(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: th.BasicChannel.Id, Message: "msg3"}
	_, resp, err = client.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: th.BasicChannel.Id, Message: "msg4 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	// todo: block join
	_, resp, err = client2.AddChannelMember(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.Error(t, err)
	CheckForbiddenStatus(t, resp)

	_, resp, err = client.DeleteChannelBlockUser(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	CheckOKStatus(t, resp)

	_, resp, err = client2.AddChannelMember(context.Background(), th.BasicChannel.Id, th.BasicUser2.Id)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	post = &model.Post{ChannelId: th.BasicChannel.Id, Message: "msg5 from user2"}
	_, resp, err = client2.CreatePost(context.Background(), post)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)
}

// TODO: DM 被 bang 之後應該不能改 channel 的設定

func TestTeamBlockUserAddRemove(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	// Team Moderator can block team members
	schemeRole := model.SchemeRoles{
		SchemeAdmin:     false,
		SchemeUser:      true,
		SchemeGuest:     false,
		SchemeVerified:  true,
		SchemeModerator: true,
	}
	client := th.SystemAdminClient
	resp0, err0 := client.UpdateTeamMemberSchemeRoles(context.Background(), th.BasicTeam.Id, th.BasicUser.Id, &schemeRole)
	require.NoError(t, err0)
	CheckOKStatus(t, resp0)

	client = th.Client
	cbul1, resp1, err1 := client.ListTeamBlockUsers(context.Background(), th.BasicTeam.Id)
	require.NoError(t, err1)
	CheckOKStatus(t, resp1)
	assert.Len(t, *cbul1, 0)

	cbu2, resp2, err2 := client.AddTeamBlockUser(context.Background(), th.BasicTeam.Id, th.BasicUser2.Id)
	require.NoError(t, err2)
	CheckCreatedStatus(t, resp2)
	assert.Equal(t, cbu2.BlockedId, th.BasicUser2.Id)
	assert.Equal(t, cbu2.TeamId, th.BasicTeam.Id)
	assert.Equal(t, cbu2.CreateBy, th.BasicUser.Id)

	cbul3, resp3, err3 := client.ListTeamBlockUsers(context.Background(), th.BasicTeam.Id)
	require.NoError(t, err3)
	CheckOKStatus(t, resp3)
	assert.Len(t, *cbul3, 1)
	assert.Equal(t, *(*cbul3)[0], *cbu2)

	_, resp4, err4 := client.DeleteTeamBlockUser(context.Background(), th.BasicTeam.Id, th.BasicUser2.Id)
	require.NoError(t, err4)
	CheckOKStatus(t, resp4)

	cbul5, resp5, err5 := client.ListTeamBlockUsers(context.Background(), th.BasicTeam.Id)
	require.NoError(t, err5)
	CheckOKStatus(t, resp5)
	assert.Len(t, *cbul5, 0)

	// checking permissions
	// th.LoginBasic2()
	otherClient := th.CreateClient()
	otherUser := th.CreateUserWithClient(otherClient)
	th.LinkUserToTeam(otherUser, th.BasicTeam)
	_, _, lErr := otherClient.Login(context.Background(), otherUser.Username, otherUser.Password)
	if lErr != nil {
		panic(lErr)
	}

	_, resp6, err6 := otherClient.AddTeamBlockUser(context.Background(), th.BasicDeletedChannel.Id, th.BasicUser2.Id)
	require.Error(t, err6)
	CheckNotFoundStatus(t, resp6)

	_, resp7, err7 := otherClient.DeleteTeamBlockUser(context.Background(), th.BasicTeam.Id, th.BasicUser2.Id)
	require.Error(t, err7)
	CheckForbiddenStatus(t, resp7)

	_, resp8, err8 := otherClient.ListUserBlockUsers(context.Background(), th.BasicTeam.Id)
	require.Error(t, err8)
	CheckForbiddenStatus(t, resp8)
}

// put the user into team block list also remove he/she from team.
func TestTeamBlockUserChannel(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	sysAdmClient := th.SystemAdminClient
	client := th.Client

	channel0, resp0, err0 := client.GetChannelByNameForTeamName(th.Context.Context(), th.BasicChannel.Name, th.BasicTeam.Name, "")
	require.NoError(t, err0)
	CheckOKStatus(t, resp0)

	cbu1, resp1, err1 := sysAdmClient.AddTeamBlockUser(th.Context.Context(), th.BasicTeam.Id, th.BasicUser.Id)
	require.NoError(t, err1)
	CheckCreatedStatus(t, resp1)
	require.Equal(t, cbu1.BlockedId, th.BasicUser.Id)
	require.Equal(t, cbu1.TeamId, th.BasicTeam.Id)
	require.Equal(t, cbu1.CreateBy, th.SystemAdminUser.Id)

	boolTrue := true
	teamPatch := model.TeamPatch{
		AllowOpenInvite: &boolTrue,
	}
	team, patchResp, patchErr := sysAdmClient.PatchTeam(th.Context.Context(), th.BasicTeam.Id, &teamPatch)
	require.NoError(t, patchErr)
	CheckOKStatus(t, patchResp)
	require.True(t, team.AllowOpenInvite)

	t.Run("blocked user is removed from team", func(t *testing.T) {
		teams, resp2, err2 := client.GetTeamsForUser(th.Context.Context(), th.BasicUser.Id, "")
		require.NoError(t, err2)
		CheckOKStatus(t, resp2)
		for _, team := range teams {
			assert.NotEqual(t, team.Id, th.BasicTeam.Id)
		}
	})

	t.Run("blocked user is removed from channel", func(t *testing.T) {
		_, resp3, err3 := client.GetPostsForChannel(th.Context.Context(), channel0.Id, 0, 100, "", true, true)
		require.Error(t, err3)
		CheckForbiddenStatus(t, resp3)
	})

	t.Run("blocked user cannot join the team", func(t *testing.T) {
		_, respA, errA := client.AddTeamMember(th.Context.Context(), th.BasicTeam.Id, th.BasicUser.Id)
		require.Error(t, errA)
		CheckForbiddenStatus(t, respA)

		_, respA2, errA2 := client.AddTeamMemberFromInvite(th.Context.Context(), "", th.BasicTeam.InviteId)
		require.Error(t, errA2)
		CheckForbiddenStatus(t, respA2)

		//TODO AddTeamMemberFromToken
		// _, respA2, errA2 := client.AddTeamMemberFromInvite(th.Context.Context(), "", th.BasicTeam.InviteId)
		// require.Error(t, errA2)
		// CheckForbiddenStatus(t, respA2)

		_, respA3, errA3 := sysAdmClient.AddTeamMember(th.Context.Context(), th.BasicTeam.Id, th.BasicUser.Id)
		require.Error(t, errA3)
		CheckBadRequestStatus(t, respA3)

		_, respA4, errA4 := sysAdmClient.AddTeamMembers(th.Context.Context(), th.BasicTeam.Id, []string{th.BasicUser.Id,})
		require.Error(t, errA4)
		CheckBadRequestStatus(t, respA4)
	})

	t.Run("user removed from block list can join the team", func(t *testing.T) {
		_, resp5, err5 := sysAdmClient.DeleteTeamBlockUser(th.Context.Context(), th.BasicTeam.Id, th.BasicUser.Id)
		require.NoError(t, err5)
		CheckOKStatus(t, resp5)
		_, resp6, err6 := client.AddTeamMember(th.Context.Context(), th.BasicTeam.Id, th.BasicUser.Id)
		require.NoError(t, err6)
		CheckCreatedStatus(t, resp6)
	})
}
