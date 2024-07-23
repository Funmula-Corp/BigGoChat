// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/app"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/sqlstore"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/utils"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/config"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/i18n"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func initDBCommandContextCobra(command *cobra.Command, readOnlyConfigStore bool, options ...app.Option) (*app.App, error) {
	a, err := initDBCommandContext(getConfigDSN(command, config.GetEnvironment()), readOnlyConfigStore, options...)
	if err != nil {
		// Returning an error just prints the usage message, so actually panic
		panic(err)
	}

	a.InitPlugins(request.EmptyContext(a.Log()), *a.Config().PluginSettings.Directory, *a.Config().PluginSettings.ClientDirectory)
	a.DoAppMigrations()

	return a, nil
}

func InitDBCommandContextCobra(command *cobra.Command, options ...app.Option) (*app.App, error) {
	return initDBCommandContextCobra(command, true, options...)
}

func initDBCommandContext(configDSN string, readOnlyConfigStore bool, options ...app.Option) (*app.App, error) {
	if err := utils.TranslationsPreInit(); err != nil {
		return nil, err
	}
	model.AppErrorInit(i18n.T)

	// The option order is important as app.Config option reads app.StartMetrics option.
	options = append(options, app.Config(configDSN, readOnlyConfigStore, nil))
	s, err := app.NewServer(options...)
	if err != nil {
		return nil, err
	}

	a := app.New(app.ServerConnector(s.Channels()))

	if model.BuildEnterpriseReady == "true" {
		a.Srv().LoadLicense()
	}

	return a, nil
}

func initStoreCommandContextCobra(logger mlog.LoggerIFace, command *cobra.Command) (store.Store, error) {
	cfgDSN := getConfigDSN(command, config.GetEnvironment())
	cfgStore, err := config.NewStoreFromDSN(cfgDSN, true, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}

	config := cfgStore.Get()
	return sqlstore.New(config.SqlSettings, logger, nil)
}
