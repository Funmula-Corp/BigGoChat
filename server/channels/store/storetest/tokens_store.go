// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package storetest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/request"
)

func TestTokensStore(t *testing.T, rctx request.CTX, ss store.Store) {
	t.Run("TokensCleanup", func(t *testing.T) { testTokensCleanup(t, rctx, ss) })
}

func testTokensCleanup(t *testing.T, rctx request.CTX, ss store.Store) {
	now := model.GetMillis()

	for i := 0; i < 10; i++ {
		err := ss.Token().Save(&model.Token{
			Token:    model.NewRandomString(model.TokenSize),
			CreateAt: now - int64(i),
			Type:     model.TokenTypeOAuth,
			Extra:    "",
		})
		require.NoError(t, err)
	}

	tokens, err := ss.Token().GetAllTokensByType(model.TokenTypeOAuth)
	require.NoError(t, err)
	assert.Len(t, tokens, 10)

	ss.Token().Cleanup(now + int64(1))

	tokens, err = ss.Token().GetAllTokensByType(model.TokenTypeOAuth)
	require.NoError(t, err)
	assert.Len(t, tokens, 0)
}
