package pluginapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/public/plugin/plugintest"
	"git.biggo.com/Funmula/BigGoChat/server/public/pluginapi"
)

func TestStore(t *testing.T) {
	t.Run("master db singleton", func(t *testing.T) {
		api := &plugintest.API{}

		driver := &plugintest.Driver{}
		defer driver.AssertExpectations(t)
		driver.On("Conn", true).Return("test", nil)
		driver.On("ConnPing", "test").Return(nil)
		driver.On("ConnClose", "test").Return(nil)

		store := pluginapi.NewClient(api, driver).Store

		db1, err := store.GetMasterDB()
		require.NoError(t, err)
		require.NotNil(t, db1)

		db2, err := store.GetMasterDB()
		require.NoError(t, err)
		require.NotNil(t, db2)

		require.Same(t, db1, db2)
		require.NoError(t, store.Close())
	})

	t.Run("master db fallback", func(t *testing.T) {
		config := &model.Config{
			SqlSettings: model.SqlSettings{
				DriverName:                  model.NewString("ramsql"),
				DataSource:                  model.NewString("TestStore-master-db"),
				ConnMaxLifetimeMilliseconds: model.NewInt(2),
			},
		}

		driver := &plugintest.Driver{}
		defer driver.AssertExpectations(t)
		driver.On("Conn", true).Return("test", nil)
		driver.On("ConnPing", "test").Return(nil)
		driver.On("ConnClose", "test").Return(nil)

		api := &plugintest.API{}
		defer api.AssertExpectations(t)
		store := pluginapi.NewClient(api, driver).Store

		api.On("GetUnsanitizedConfig").Return(config)
		masterDB, err := store.GetMasterDB()
		require.NoError(t, err)
		require.NotNil(t, masterDB)

		// No replica is set up, should fallback to master
		replicaDB, err := store.GetReplicaDB()
		require.NoError(t, err)
		require.Same(t, replicaDB, masterDB)

		require.NoError(t, store.Close())
	})

	t.Run("replica db singleton", func(t *testing.T) {
		config := &model.Config{
			SqlSettings: model.SqlSettings{
				DriverName:                  model.NewString("ramsql"),
				DataSource:                  model.NewString("TestStore-master-db"),
				DataSourceReplicas:          []string{"TestStore-master-db"},
				ConnMaxLifetimeMilliseconds: model.NewInt(2),
			},
		}

		api := &plugintest.API{}
		defer api.AssertExpectations(t)
		api.On("GetUnsanitizedConfig").Return(config)

		driver := &plugintest.Driver{}
		defer driver.AssertExpectations(t)
		driver.On("Conn", false).Return("test", nil)
		driver.On("ConnPing", "test").Return(nil)
		driver.On("ConnClose", "test").Return(nil)

		store := pluginapi.NewClient(api, driver).Store

		db1, err := store.GetReplicaDB()
		require.NoError(t, err)
		require.NotNil(t, db1)

		db2, err := store.GetReplicaDB()
		require.NoError(t, err)
		require.NotNil(t, db2)

		require.Same(t, db1, db2)
		require.NoError(t, store.Close())
	})
}
