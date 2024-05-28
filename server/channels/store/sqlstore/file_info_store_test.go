// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/searchtest"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/storetest"
)

func TestFileInfoStore(t *testing.T) {
	StoreTestWithSqlStore(t, storetest.TestFileInfoStore)
}

func TestSearchFileInfoStore(t *testing.T) {
	StoreTestWithSearchTestEngine(t, searchtest.TestSearchFileInfoStore)
}
