// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	mm_model "git.biggo.com/Funmula/BigGoChat/server/public/model"
)

func TestOnboardingSavesOrganizationName(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	err := th.App.CompleteOnboarding(th.Context, &mm_model.CompleteOnboardingRequest{
		Organization: "Mattermost In Tests",
	})
	require.Nil(t, err)
	defer func() {
		th.App.Srv().Store().System().PermanentDeleteByName(mm_model.SystemOrganizationName)
	}()

	sys, storeErr := th.App.Srv().Store().System().GetByName(mm_model.SystemOrganizationName)
	require.NoError(t, storeErr)
	require.Equal(t, "Mattermost In Tests", sys.Value)
}
