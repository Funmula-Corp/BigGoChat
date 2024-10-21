// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"testing"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest"
)

func TestSystemStore(t *testing.T) {
	StoreTest(t, storetest.TestSystemStore)
}
