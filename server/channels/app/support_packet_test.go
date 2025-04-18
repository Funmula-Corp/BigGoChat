// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"git.biggo.com/Funmula/BigGoChat/server/public/model"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/app/platform"
	"git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest/mocks"
	smocks "git.biggo.com/Funmula/BigGoChat/server/v8/channels/store/storetest/mocks"
	"git.biggo.com/Funmula/BigGoChat/server/v8/config"
	emocks "git.biggo.com/Funmula/BigGoChat/server/v8/einterfaces/mocks"
	fmocks "git.biggo.com/Funmula/BigGoChat/server/v8/platform/shared/filestore/mocks"
)

func TestCreatePluginsFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Happy path where we have a plugins file with no err
	fileData, err := th.App.createPluginsFile(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "plugins.json", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)

	// Turn off plugins so we can get an error
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.PluginSettings.Enable = false
	})

	// Plugins off in settings so no fileData and we get a warning instead
	fileData, err = th.App.createPluginsFile(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed to get plugin list for support package")
}

func TestGenerateSupportPacketYaml(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	licenseUsers := 100
	license := model.NewTestLicense("ldap")
	license.Features.Users = model.NewInt(licenseUsers)
	th.App.Srv().SetLicense(license)

	generateSupportPacket := func(t *testing.T) *model.SupportPacket {
		t.Helper()

		fileData, err := th.App.generateSupportPacketYaml(th.Context)
		require.NotNil(t, fileData)
		assert.Equal(t, "support_packet.yaml", fileData.Filename)
		assert.Positive(t, len(fileData.Body))
		assert.NoError(t, err)

		var packet model.SupportPacket
		require.NoError(t, yaml.Unmarshal(fileData.Body, &packet))
		require.NotNil(t, packet)
		return &packet
	}

	t.Run("Happy path", func(t *testing.T) {
		// Happy path where we have a support packet yaml file without any warnings
		packet := generateSupportPacket(t)

		/* Build information */
		assert.NotEmpty(t, packet.ServerOS)
		assert.NotEmpty(t, packet.ServerArchitecture)
		assert.Equal(t, model.CurrentVersion, packet.ServerVersion)
		// BuildHash is not present in tests

		/* DB */
		assert.NotEmpty(t, packet.DatabaseType)
		assert.NotEmpty(t, packet.DatabaseVersion)
		assert.NotEmpty(t, packet.DatabaseSchemaVersion)
		assert.Zero(t, packet.WebsocketConnections)
		assert.NotZero(t, packet.MasterDbConnections)
		assert.Zero(t, packet.ReplicaDbConnections)

		/* Cluster */
		assert.Empty(t, packet.ClusterID)

		/* File store */
		assert.Equal(t, "local", packet.FileDriver)
		assert.Equal(t, "OK", packet.FileStatus)

		/* LDAP */
		assert.Empty(t, packet.LdapVendorName)
		assert.Empty(t, packet.LdapVendorVersion)

		/* Elastic Search */
		assert.Empty(t, packet.ElasticServerVersion)
		assert.Empty(t, packet.ElasticServerPlugins)

		/* License */
		assert.Equal(t, "My awesome Company", packet.LicenseTo)
		assert.Equal(t, licenseUsers, packet.LicenseSupportedUsers)
		assert.Equal(t, false, packet.LicenseIsTrial)

		/* Server stats */
		assert.Equal(t, 4, packet.ActiveUsers) // from InitBasic()
		assert.Equal(t, 0, packet.DailyActiveUsers)
		assert.Equal(t, 0, packet.MonthlyActiveUsers)
		assert.Equal(t, 0, packet.InactiveUserCount)
		assert.Equal(t, 7, packet.TotalPosts)    // from InitBasic()
		assert.Equal(t, 3, packet.TotalChannels) // from InitBasic()
		assert.Equal(t, 1, packet.TotalTeams)    // from InitBasic()

		/* Jobs */
		assert.Empty(t, packet.DataRetentionJobs)
		assert.Empty(t, packet.MessageExportJobs)
		assert.Empty(t, packet.ElasticPostIndexingJobs)
		assert.Empty(t, packet.ElasticPostAggregationJobs)
		assert.Empty(t, packet.LdapSyncJobs)
		assert.Empty(t, packet.MigrationJobs)
	})

	t.Run("post count should be present if number of users extends AnalyticsSettings.MaxUsersForStatistics", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) {
			cfg.AnalyticsSettings.MaxUsersForStatistics = model.NewInt(1)
		})

		for i := 0; i < 5; i++ {
			p := th.CreatePost(th.BasicChannel)
			require.NotNil(t, p)
		}

		// InitBasic() already creats 7 posts
		// include two posts for user UserUnverified
		packet := generateSupportPacket(t)
		assert.Equal(t, 12, packet.TotalPosts)
	})

	t.Run("filestore fails", func(t *testing.T) {
		fb := &fmocks.FileBackend{}
		platform.SetFileStore(fb)(th.Server.Platform())
		fb.On("DriverName").Return("mock")
		fb.On("TestConnection").Return(errors.New("all broken"))

		packet := generateSupportPacket(t)

		assert.Equal(t, "mock", packet.FileDriver)
		assert.Equal(t, "FAIL: all broken", packet.FileStatus)
	})

	t.Run("no LDAP vendor info", func(t *testing.T) {
		ldapMock := &emocks.LdapInterface{}
		ldapMock.On(
			"GetVendorNameAndVendorVersion",
			mock.AnythingOfType("*request.Context"),
		).Return("", "", nil)
		th.App.Channels().Ldap = ldapMock

		packet := generateSupportPacket(t)

		assert.Equal(t, "unknown", packet.LdapVendorName)
		assert.Equal(t, "unknown", packet.LdapVendorVersion)
	})

	t.Run("found LDAP vendor info", func(t *testing.T) {
		ldapMock := &emocks.LdapInterface{}
		ldapMock.On(
			"GetVendorNameAndVendorVersion",
			mock.AnythingOfType("*request.Context"),
		).Return("some vendor", "v1.0.0", nil)
		th.App.Channels().Ldap = ldapMock

		packet := generateSupportPacket(t)

		assert.Equal(t, "some vendor", packet.LdapVendorName)
		assert.Equal(t, "v1.0.0", packet.LdapVendorVersion)
	})
}

func TestGenerateSupportPacket(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	})

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.LogSettings.FileLocation = dir
	})

	logLocation := config.GetLogFileLocation(dir)
	notificationsLogLocation := config.GetNotificationsLogFileLocation(dir)

	genMockLogFiles := func() {
		d1 := []byte("hello\ngo\n")
		genErr := os.WriteFile(logLocation, d1, 0777)
		require.NoError(t, genErr)
		genErr = os.WriteFile(notificationsLogLocation, d1, 0777)
		require.NoError(t, genErr)
	}
	genMockLogFiles()

	t.Run("generate support packet with logs", func(t *testing.T) {
		fileDatas := th.App.GenerateSupportPacket(th.Context, &model.SupportPacketOptions{
			IncludeLogs: true,
		})
		var rFileNames []string
		testFiles := []string{
			"support_packet.yaml",
			"metadata.yaml",
			"plugins.json",
			"sanitized_config.json",
			"mattermost.log",
			"notifications.log",
			"cpu.prof",
			"heap.prof",
			"goroutines",
		}
		for _, fileData := range fileDatas {
			require.NotNil(t, fileData)
			assert.Positive(t, len(fileData.Body))

			rFileNames = append(rFileNames, fileData.Filename)
		}
		assert.ElementsMatch(t, testFiles, rFileNames)
	})

	t.Run("generate support packet without logs", func(t *testing.T) {
		fileDatas := th.App.GenerateSupportPacket(th.Context, &model.SupportPacketOptions{
			IncludeLogs: false,
		})

		testFiles := []string{
			"support_packet.yaml",
			"metadata.yaml",
			"plugins.json",
			"sanitized_config.json",
			"cpu.prof",
			"heap.prof",
			"goroutines",
		}
		var rFileNames []string
		for _, fileData := range fileDatas {
			require.NotNil(t, fileData)
			assert.Positive(t, len(fileData.Body))

			rFileNames = append(rFileNames, fileData.Filename)
		}
		assert.ElementsMatch(t, testFiles, rFileNames)
	})

	t.Run("remove the log files and ensure that warning.txt file is generated", func(t *testing.T) {
		// Remove these two files and ensure that warning.txt file is generated
		err = os.Remove(logLocation)
		require.NoError(t, err)
		err = os.Remove(notificationsLogLocation)
		require.NoError(t, err)
		t.Cleanup(genMockLogFiles)

		fileDatas := th.App.GenerateSupportPacket(th.Context, &model.SupportPacketOptions{
			IncludeLogs: true,
		})
		testFiles := []string{
			"support_packet.yaml",
			"metadata.yaml",
			"plugins.json",
			"sanitized_config.json",
			"cpu.prof",
			"heap.prof",
			"warning.txt",
			"goroutines",
		}
		var rFileNames []string
		for _, fileData := range fileDatas {
			require.NotNil(t, fileData)
			assert.Positive(t, len(fileData.Body))

			rFileNames = append(rFileNames, fileData.Filename)
		}
		assert.ElementsMatch(t, testFiles, rFileNames)
	})

	t.Run("steps that generated an error should still return file data", func(t *testing.T) {
		mockStore := smocks.Store{}

		// Mock the post store to trigger an error
		ps := &smocks.PostStore{}
		ps.On("AnalyticsPostCount", &model.PostCountOptions{}).Return(int64(0), errors.New("all broken"))
		ps.On("ClearCaches")
		mockStore.On("Post").Return(ps)

		mockStore.On("User").Return(th.App.Srv().Store().User())
		mockStore.On("Channel").Return(th.App.Srv().Store().Channel())
		mockStore.On("Post").Return(th.App.Srv().Store().Post())
		mockStore.On("Team").Return(th.App.Srv().Store().Team())
		mockStore.On("Job").Return(th.App.Srv().Store().Job())
		mockStore.On("FileInfo").Return(th.App.Srv().Store().FileInfo())
		mockStore.On("Webhook").Return(th.App.Srv().Store().Webhook())
		mockStore.On("System").Return(th.App.Srv().Store().System())
		mockStore.On("License").Return(th.App.Srv().Store().License())
		mockStore.On("Close").Return(nil)
		mockStore.On("GetDBSchemaVersion").Return(1, nil)
		mockStore.On("GetDbVersion", false).Return("1.0.0", nil)

		clusterMockStore := &mocks.ClusterDiscoveryStore{}
		clusterMockStore.On("GetAll", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]*model.ClusterDiscovery{}, nil)
		mockStore.On("ClusterDiscovery").Return(clusterMockStore)

		th.App.Srv().SetStore(&mockStore)

		fileDatas := th.App.GenerateSupportPacket(th.Context, &model.SupportPacketOptions{
			IncludeLogs: false,
		})

		var rFileNames []string
		for _, fileData := range fileDatas {
			require.NotNil(t, fileData)
			assert.Positive(t, len(fileData.Body))

			rFileNames = append(rFileNames, fileData.Filename)
		}
		assert.Contains(t, rFileNames, "warning.txt")
		assert.Contains(t, rFileNames, "support_packet.yaml")
	})
}

func TestGetNotificationsLog(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Disable notifications file setting in config so we should get an warning
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.NotificationLogSettings.EnableFile = false
	})

	fileData, err := th.App.getNotificationsLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "Unable to retrieve notifications.log because LogSettings: EnableFile is set to false")

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	})

	// Enable notifications file but point to an empty directory to get an error trying to read the file
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.NotificationLogSettings.EnableFile = true
		*cfg.LogSettings.FileLocation = dir
	})

	logLocation := config.GetNotificationsLogFileLocation(dir)

	// There is no notifications.log file yet, so this fails
	fileData, err = th.App.getNotificationsLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed read notifcation log file at path "+logLocation)

	// Happy path where we have file and no error
	d1 := []byte("hello\ngo\n")
	err = os.WriteFile(logLocation, d1, 0777)
	require.NoError(t, err)

	fileData, err = th.App.getNotificationsLog(th.Context)
	assert.NoError(t, err)
	require.NotNil(t, fileData)
	assert.Equal(t, "notifications.log", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
}

func TestGetMattermostLog(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// disable mattermost log file setting in config so we should get an warning
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.LogSettings.EnableFile = false
	})

	fileData, err := th.App.GetMattermostLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "Unable to retrieve mattermost.log because LogSettings: EnableFile is set to false")

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	})

	// Enable log file but point to an empty directory to get an error trying to read the file
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.LogSettings.EnableFile = true
		*cfg.LogSettings.FileLocation = dir
	})

	logLocation := config.GetLogFileLocation(dir)

	// There is no mattermost.log file yet, so this fails
	fileData, err = th.App.GetMattermostLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed read mattermost log file at path "+logLocation)

	// Happy path where we get a log file and no warning
	d1 := []byte("hello\ngo\n")
	err = os.WriteFile(logLocation, d1, 0777)
	require.NoError(t, err)

	fileData, err = th.App.GetMattermostLog(th.Context)
	require.NoError(t, err)
	require.NotNil(t, fileData)
	assert.Equal(t, "mattermost.log", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
}

func TestCreateSanitizedConfigFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Happy path where we have a sanitized config file with no err
	fileData, err := th.App.createSanitizedConfigFile(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "sanitized_config.json", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)
}

func TestCreateSupportPacketMetadata(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	t.Run("Happy path", func(t *testing.T) {
		fileData, err := th.App.createSupportPacketMetadata(th.Context)
		require.NoError(t, err)
		require.NotNil(t, fileData)
		assert.Equal(t, "metadata.yaml", fileData.Filename)
		assert.Positive(t, len(fileData.Body))

		metadate, err := model.ParsePacketMetadata(fileData.Body)
		assert.NoError(t, err)
		require.NotNil(t, metadate)
		assert.Equal(t, model.SupportPacketType, metadate.Type)
		assert.Equal(t, model.CurrentVersion, metadate.ServerVersion)
		assert.NotEmpty(t, metadate.ServerID)
		assert.NotEmpty(t, metadate.GeneratedAt)
	})
}
