// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchlayer_test

import (
	"os"
	"sync"
	"testing"

	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/searchlayer"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/sqlstore"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/testlib"
	"git.biggo.com/Funmula/BigGoChat/server/v8/platform/services/searchengine"
	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/shared/mlog"
	"github.com/stretchr/testify/require"
)

// Test to verify race condition on UpdateConfig. The test must run with -race flag in order to verify
// that there is no race. Ref: (#MM-30868)
func TestUpdateConfigRace(t *testing.T) {
	logger := mlog.CreateTestLogger(t)

	driverName := os.Getenv("MM_SQLSETTINGS_DRIVERNAME")
	if driverName == "" {
		driverName = model.DatabaseDriverPostgres
	}
	settings := storetest.MakeSqlSettings(driverName, false)
	store, err := sqlstore.New(*settings, logger, nil)
	require.NoError(t, err)

	cfg := &model.Config{}
	cfg.SetDefaults()
	cfg.ClusterSettings.GossipPort = model.NewInt(9999)
	searchEngine := searchengine.NewBroker(cfg)
	layer := searchlayer.NewSearchLayer(&testlib.TestStore{Store: store}, searchEngine, cfg)
	var wg sync.WaitGroup

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			layer.UpdateConfig(cfg.Clone())
		}()
	}

	wg.Wait()
}
