// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"database/sql"
	"encoding/base64"
	"testing"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"github.com/stretchr/testify/require"
)

func TestGetRemoteClusters(t *testing.T) {
	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()
		rcs, resp, err := th.SystemAdminClient.GetRemoteClusters(context.Background(), 0, 999999, model.RemoteClusterQueryFilter{})
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rcs)
	})

	th := setupForSharedChannels(t)
	defer th.TearDown()

	newRCs := []*model.RemoteCluster{
		{
			RemoteId:    model.NewId(),
			Name:        "remote1",
			SiteURL:     "http://example1.com",
			CreatorId:   th.SystemAdminUser.Id,
			Token:       model.NewId(),
			RemoteToken: model.NewId(),
		},
		{
			RemoteId:  model.NewId(),
			Name:      "remote2",
			SiteURL:   "http://example2.com",
			CreatorId: th.SystemAdminUser.Id,
		},
		{
			RemoteId:  model.NewId(),
			Name:      "remote3",
			SiteURL:   "http://example3.com",
			CreatorId: th.SystemAdminUser.Id,
			PluginID:  model.NewId(),
		},
	}

	for _, rc := range newRCs {
		_, appErr := th.App.AddRemoteCluster(rc)
		require.Nil(t, appErr)
	}

	t.Run("The returned data should be sanitized", func(t *testing.T) {
		rcs, resp, err := th.SystemAdminClient.GetRemoteClusters(context.Background(), 0, 999999, model.RemoteClusterQueryFilter{})
		CheckOKStatus(t, resp)
		require.NoError(t, err)
		require.Contains(t, rcs[0].Name, "remote")
		require.Zero(t, rcs[0].Token)
		require.Zero(t, rcs[0].RemoteToken)
	})

	testCases := []struct {
		Name               string
		Client             *model.Client4
		Page               int
		PerPage            int
		Filter             model.RemoteClusterQueryFilter
		ExpectedStatusCode int
		ExpectedError      bool
		ExpectedNames      []string
	}{
		{
			Name:               "Should reject if the user has not sufficient permissions",
			Client:             th.Client,
			Page:               0,
			PerPage:            999999,
			Filter:             model.RemoteClusterQueryFilter{},
			ExpectedStatusCode: 403,
			ExpectedError:      true,
			ExpectedNames:      []string{},
		},
		{
			Name:               "Should return all remote clusters",
			Client:             th.SystemAdminClient,
			Page:               0,
			PerPage:            999999,
			Filter:             model.RemoteClusterQueryFilter{},
			ExpectedStatusCode: 200,
			ExpectedError:      false,
			ExpectedNames:      []string{"remote1", "remote2", "remote3"},
		},
		{
			Name:               "Should return all remote clusters but those belonging to plugins",
			Client:             th.SystemAdminClient,
			Page:               0,
			PerPage:            999999,
			Filter:             model.RemoteClusterQueryFilter{ExcludePlugins: true},
			ExpectedStatusCode: 200,
			ExpectedError:      false,
			ExpectedNames:      []string{"remote1", "remote2"},
		},
		{
			Name:               "Should return only remote clusters belonging to plugins",
			Client:             th.SystemAdminClient,
			Page:               0,
			PerPage:            999999,
			Filter:             model.RemoteClusterQueryFilter{OnlyPlugins: true},
			ExpectedStatusCode: 200,
			ExpectedError:      false,
			ExpectedNames:      []string{"remote3"},
		},
		{
			Name:               "Should work as a paginated endpoint",
			Client:             th.SystemAdminClient,
			Page:               1,
			PerPage:            1,
			Filter:             model.RemoteClusterQueryFilter{},
			ExpectedStatusCode: 200,
			ExpectedError:      false,
			ExpectedNames:      []string{"remote2"},
		},
		{
			Name:               "Should return an empty set with a successful status",
			Client:             th.SystemAdminClient,
			Page:               0,
			PerPage:            999999,
			Filter:             model.RemoteClusterQueryFilter{InChannel: model.NewId()},
			ExpectedStatusCode: 200,
			ExpectedError:      false,
			ExpectedNames:      []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			rcs, resp, err := tc.Client.GetRemoteClusters(context.Background(), tc.Page, tc.PerPage, tc.Filter)
			checkHTTPStatus(t, resp, tc.ExpectedStatusCode)
			if tc.ExpectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Len(t, rcs, len(tc.ExpectedNames))
			names := []string{}
			for _, rc := range rcs {
				names = append(names, rc.Name)
			}
			require.ElementsMatch(t, tc.ExpectedNames, names)
		})
	}
}

func TestCreateRemoteCluster(t *testing.T) {
	rcWithTeamAndPassword := &model.RemoteClusterWithPassword{
		RemoteCluster: &model.RemoteCluster{
			Name:    "remotecluster",
			SiteURL: "http://example.com",
			Token:   model.NewId(),
		},
		Password: "mysupersecret",
	}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		rcWithInvite, resp, err := th.SystemAdminClient.CreateRemoteCluster(context.Background(), rcWithTeamAndPassword)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rcWithInvite)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	t.Run("Should not work if the user doesn't have the right permissions", func(t *testing.T) {
		rcWithInvite, resp, err := th.Client.CreateRemoteCluster(context.Background(), rcWithTeamAndPassword)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rcWithInvite)
	})

	t.Run("Should not work if the siteURL is not set in the configuration", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "" })
		rcWithInvite, resp, err := th.SystemAdminClient.CreateRemoteCluster(context.Background(), rcWithTeamAndPassword)
		CheckUnprocessableEntityStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rcWithInvite)
	})

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "http://localhost:8065" })

	t.Run("Should enforce the presence of the password", func(t *testing.T) {
		// clean the password and check the response
		rcWithTeamAndPassword.Password = ""

		rcWithInvite, resp, err := th.SystemAdminClient.CreateRemoteCluster(context.Background(), rcWithTeamAndPassword)
		CheckBadRequestStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rcWithInvite)

		// reset password for the next tests
		rcWithTeamAndPassword.Password = "mysupersecret"
	})

	t.Run("Should return a sanitized remote cluster and its invite", func(t *testing.T) {
		rcWithInvite, resp, err := th.SystemAdminClient.CreateRemoteCluster(context.Background(), rcWithTeamAndPassword)
		CheckCreatedStatus(t, resp)
		require.NoError(t, err)
		require.Equal(t, rcWithTeamAndPassword.Name, rcWithInvite.RemoteCluster.Name)
		require.NotZero(t, rcWithInvite.Invite)
		require.Zero(t, rcWithInvite.RemoteCluster.Token)
		require.Zero(t, rcWithInvite.RemoteCluster.RemoteToken)

		rc, appErr := th.App.GetRemoteCluster(rcWithInvite.RemoteCluster.RemoteId)
		require.Nil(t, appErr)
		require.Equal(t, rcWithTeamAndPassword.Name, rc.Name)

		rci, appErr := th.App.DecryptRemoteClusterInvite(rcWithInvite.Invite, rcWithTeamAndPassword.Password)
		require.Nil(t, appErr)
		require.Equal(t, rc.RemoteId, rci.RemoteId)
		require.Equal(t, rc.RemoteToken, rci.Token)
		require.Equal(t, th.App.GetSiteURL(), rci.SiteURL)
	})
}

func TestRemoteClusterAcceptinvite(t *testing.T) {
	rcAcceptInvite := &model.RemoteClusterAcceptInvite{
		Name:     "remotecluster",
		Invite:   "myinvitecode",
		Password: "mysupersecret",
	}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		rc, resp, err := th.SystemAdminClient.RemoteClusterAcceptInvite(context.Background(), rcAcceptInvite)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rc)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	remoteId := model.NewId()
	invite := &model.RemoteClusterInvite{
		RemoteId: remoteId,
		SiteURL:  "http://localhost:8065",
		Token:    "token",
	}
	password := "mysupersecret"
	encrypted, err := invite.Encrypt(password)
	require.NoError(t, err)
	encoded := base64.URLEncoding.EncodeToString(encrypted)
	rcAcceptInvite.Invite = encoded

	t.Run("Should not work if the siteURL is not set in the configuration", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "" })
		rc, resp, err := th.SystemAdminClient.RemoteClusterAcceptInvite(context.Background(), rcAcceptInvite)
		CheckUnprocessableEntityStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rc)
	})

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "http://localhost:8065" })

	t.Run("should fail if the parameters are not valid", func(t *testing.T) {
		rcAcceptInvite.Name = ""
		defer func() { rcAcceptInvite.Name = "remotecluster" }()

		rc, resp, err := th.SystemAdminClient.RemoteClusterAcceptInvite(context.Background(), rcAcceptInvite)
		CheckBadRequestStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rc)
	})

	t.Run("should fail with the correct status code if the invite returns an app error", func(t *testing.T) {
		rcAcceptInvite.Invite = "malformedinvite"
		// reset the invite after
		defer func() { rcAcceptInvite.Invite = encoded }()

		rc, resp, err := th.SystemAdminClient.RemoteClusterAcceptInvite(context.Background(), rcAcceptInvite)
		CheckBadRequestStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rc)
	})

	t.Run("should not work if the user doesn't have the right permissions", func(t *testing.T) {
		rc, resp, err := th.Client.RemoteClusterAcceptInvite(context.Background(), rcAcceptInvite)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, rc)
	})

	t.Run("should return a sanitized remote cluster if the action succeeds", func(t *testing.T) {
		t.Skip("Requires server2server communication: ToBeImplemented")
	})
}

func TestGenerateRemoteClusterInvite(t *testing.T) {
	password := "mysupersecret"

	newRC := &model.RemoteCluster{
		Name:    "remotecluster",
		SiteURL: "http://example.com",
		Token:   model.NewId(),
	}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		newRC.CreatorId = th.SystemAdminUser.Id

		rc, appErr := th.App.AddRemoteCluster(newRC)
		require.Nil(t, appErr)
		require.NotZero(t, rc.RemoteId)

		inviteCode, resp, err := th.SystemAdminClient.GenerateRemoteClusterInvite(context.Background(), rc.RemoteId, password)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Zero(t, inviteCode)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	newRC.CreatorId = th.SystemAdminUser.Id

	rc, appErr := th.App.AddRemoteCluster(newRC)
	require.Nil(t, appErr)
	require.NotZero(t, rc.RemoteId)

	t.Run("Should not work if the siteURL is not set in the configuration", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "" })
		inviteCode, resp, err := th.SystemAdminClient.GenerateRemoteClusterInvite(context.Background(), rc.RemoteId, password)
		CheckUnprocessableEntityStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, inviteCode)
	})

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "http://localhost:8065" })

	t.Run("Should not work if the user doesn't have the right permissions", func(t *testing.T) {
		inviteCode, resp, err := th.Client.GenerateRemoteClusterInvite(context.Background(), rc.RemoteId, password)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, inviteCode)
	})

	t.Run("should not work if the remote cluster doesn't exist", func(t *testing.T) {
		inviteCode, resp, err := th.SystemAdminClient.GenerateRemoteClusterInvite(context.Background(), model.NewId(), password)
		CheckNotFoundStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, inviteCode)
	})

	t.Run("should not work if the password has been provided", func(t *testing.T) {
		inviteCode, resp, err := th.SystemAdminClient.GenerateRemoteClusterInvite(context.Background(), rc.RemoteId, "")
		CheckBadRequestStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, inviteCode)
	})

	t.Run("should generate a valid invite code", func(t *testing.T) {
		inviteCode, resp, err := th.SystemAdminClient.GenerateRemoteClusterInvite(context.Background(), rc.RemoteId, password)
		CheckCreatedStatus(t, resp)
		require.NoError(t, err)
		require.NotEmpty(t, inviteCode)

		invite, appErr := th.App.DecryptRemoteClusterInvite(inviteCode, password)
		require.Nil(t, appErr)
		require.Equal(t, rc.RemoteId, invite.RemoteId)
		require.Equal(t, rc.Token, invite.Token)
	})
}

func TestGetRemoteCluster(t *testing.T) {
	newRC := &model.RemoteCluster{
		Name:    "remotecluster",
		SiteURL: "http://example.com",
		Token:   model.NewId(),
	}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		newRC.CreatorId = th.SystemAdminUser.Id

		rc, appErr := th.App.AddRemoteCluster(newRC)
		require.Nil(t, appErr)
		require.NotZero(t, rc.RemoteId)
		require.NotZero(t, rc.Token)

		fetchedRC, resp, err := th.SystemAdminClient.GetRemoteCluster(context.Background(), rc.RemoteId)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, fetchedRC)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	newRC.CreatorId = th.SystemAdminUser.Id

	rc, appErr := th.App.AddRemoteCluster(newRC)
	require.Nil(t, appErr)
	require.NotZero(t, rc.RemoteId)

	t.Run("Should not work if the user doesn't have the right permissions", func(t *testing.T) {
		fetchedRC, resp, err := th.Client.GetRemoteCluster(context.Background(), rc.RemoteId)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, fetchedRC)
	})

	t.Run("should return not found if the id doesn't exist", func(t *testing.T) {
		fetchedRC, resp, err := th.SystemAdminClient.GetRemoteCluster(context.Background(), model.NewId())
		CheckNotFoundStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, fetchedRC)
	})

	t.Run("should return a sanitized remote cluster", func(t *testing.T) {
		fetchedRC, resp, err := th.SystemAdminClient.GetRemoteCluster(context.Background(), rc.RemoteId)
		CheckOKStatus(t, resp)
		require.NoError(t, err)
		require.Equal(t, rc.RemoteId, fetchedRC.RemoteId)
		require.Empty(t, fetchedRC.Token)
	})
}

func TestPatchRemoteCluster(t *testing.T) {
	newRC := &model.RemoteCluster{
		Name:        "remotecluster",
		DisplayName: "initialvalue",
		SiteURL:     "http://example.com",
		Token:       model.NewId(),
	}

	rcp := &model.RemoteClusterPatch{DisplayName: model.NewString("different value")}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		newRC.CreatorId = th.SystemAdminUser.Id

		rc, appErr := th.App.AddRemoteCluster(newRC)
		require.Nil(t, appErr)
		require.NotZero(t, rc.RemoteId)

		patchedRC, resp, err := th.SystemAdminClient.PatchRemoteCluster(context.Background(), rc.RemoteId, rcp)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, patchedRC)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	newRC.CreatorId = th.SystemAdminUser.Id

	rc, appErr := th.App.AddRemoteCluster(newRC)
	require.Nil(t, appErr)
	require.NotZero(t, rc.RemoteId)

	t.Run("Should not work if the user doesn't have the right permissions", func(t *testing.T) {
		patchedRC, resp, err := th.Client.PatchRemoteCluster(context.Background(), rc.RemoteId, rcp)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, patchedRC)
	})

	t.Run("should not work if the remote cluster is nonexistent", func(t *testing.T) {
		patchedRC, resp, err := th.SystemAdminClient.PatchRemoteCluster(context.Background(), model.NewId(), rcp)
		CheckNotFoundStatus(t, resp)
		require.Error(t, err)
		require.Empty(t, patchedRC)
	})

	t.Run("should correctly patch the remote cluster", func(t *testing.T) {
		rcp := &model.RemoteClusterPatch{DisplayName: model.NewString("patched!")}

		patchedRC, resp, err := th.SystemAdminClient.PatchRemoteCluster(context.Background(), rc.RemoteId, rcp)
		CheckOKStatus(t, resp)
		require.NoError(t, err)
		require.Equal(t, "patched!", patchedRC.DisplayName)
	})
}

func TestDeleteRemoteCluster(t *testing.T) {
	newRC := &model.RemoteCluster{
		Name:        "remotecluster",
		DisplayName: "initialvalue",
		SiteURL:     "http://example.com",
		Token:       model.NewId(),
	}

	t.Run("Should not work if the remote cluster service is not enabled", func(t *testing.T) {
		th := Setup(t)
		defer th.TearDown()

		newRC.CreatorId = th.SystemAdminUser.Id

		rc, appErr := th.App.AddRemoteCluster(newRC)
		require.Nil(t, appErr)
		require.NotZero(t, rc.RemoteId)

		resp, err := th.SystemAdminClient.DeleteRemoteCluster(context.Background(), rc.RemoteId)
		CheckNotImplementedStatus(t, resp)
		require.Error(t, err)
	})

	th := setupForSharedChannels(t).InitBasic()
	defer th.TearDown()

	newRC.CreatorId = th.SystemAdminUser.Id

	rc, appErr := th.App.AddRemoteCluster(newRC)
	require.Nil(t, appErr)
	require.NotZero(t, rc.RemoteId)

	t.Run("Should not work if the user doesn't have the right permissions", func(t *testing.T) {
		resp, err := th.Client.DeleteRemoteCluster(context.Background(), rc.RemoteId)
		CheckForbiddenStatus(t, resp)
		require.Error(t, err)
	})

	t.Run("should not work if the remote cluster is nonexistent", func(t *testing.T) {
		resp, err := th.SystemAdminClient.DeleteRemoteCluster(context.Background(), model.NewId())
		CheckNotFoundStatus(t, resp)
		require.Error(t, err)
	})

	t.Run("should correctly delete the remote cluster", func(t *testing.T) {
		resp, err := th.SystemAdminClient.DeleteRemoteCluster(context.Background(), rc.RemoteId)
		CheckNoContentStatus(t, resp)
		require.NoError(t, err)

		deletedRC, err := th.App.GetRemoteCluster(rc.RemoteId)
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.Empty(t, deletedRC)
	})

	t.Run("should return not found if the remote cluster is already deleted", func(t *testing.T) {
		resp, err := th.SystemAdminClient.DeleteRemoteCluster(context.Background(), rc.RemoteId)
		CheckNotFoundStatus(t, resp)
		require.Error(t, err)
	})
}
