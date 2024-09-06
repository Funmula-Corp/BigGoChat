package biggoengine

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/searchlayer"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/searchtest"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/sqlstore"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/storetest"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/testlib"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
)

type BiggoEngineTestSuite struct {
	suite.Suite

	SQLSettings  *model.SqlSettings
	SQLStore     *sqlstore.SqlStore
	SearchEngine *searchengine.Broker
	Store        *searchlayer.SearchStore
	BiggoEngine  *BiggoEngine
	Context      request.CTX
}

func TestBiggoEngineTestSuite(t *testing.T) {
	suite.Run(t, &BiggoEngineTestSuite{
		Context: request.TestContext(t),
	})
}

func (s *BiggoEngineTestSuite) setupStore() {
	driverName := os.Getenv("MM_SQLSETTINGS_DRIVERNAME")
	if driverName == "" {
		driverName = model.DatabaseDriverPostgres
	}
	s.SQLSettings = storetest.MakeSqlSettings(driverName, false)

	var err error
	s.SQLStore, err = sqlstore.New(*s.SQLSettings, s.Context.Logger(), nil)
	if err != nil {
		s.Require().FailNow("Cannot initialize store: %s", err.Error())
	}

	config := &model.Config{}
	config.SetDefaults()
	config.SqlSettings.DisableDatabaseSearch = model.NewBool(true)
	config.PluginSettings.Plugins[cfg.PluginName] = map[string]any{
		"enable_autocompletion": true,
		"enable_indexer":        true,
		"enable_indexing":       true,
		"enable_search":         true,
		"neo4j_host":            "172.17.0.1",
		"elasticsearch_host":    "172.17.0.1",
	}
	s.SearchEngine = searchengine.NewBroker(config)
	s.Store = searchlayer.NewSearchLayer(&testlib.TestStore{Store: s.SQLStore}, s.SearchEngine, config)

	s.BiggoEngine = NewBiggoEngine(config)
	s.BiggoEngine.isSync.Store(true)
	s.SearchEngine.RegisterBiggoEngine(s.BiggoEngine)
	if err := s.BiggoEngine.Start(); err != nil {
		s.Require().FailNow("Cannot start BiggoEngine: %s", err.Error())
	}
}

func (s *BiggoEngineTestSuite) SetupSuite() {
	//s.setupIndexes()
	s.setupStore()
}

func (s *BiggoEngineTestSuite) TestBiggoSearchStoreTests() {
	return
	searchTestEngine := &searchtest.SearchTestEngine{
		Driver: searchtest.EngineBiggo,
	}

	s.Run("TestSearchChannelStore", func() {
		searchtest.TestSearchChannelStore(s.T(), s.Store, searchTestEngine)
	})

	s.Run("TestSearchUserStore", func() {
		searchtest.TestSearchUserStore(s.T(), s.Store, searchTestEngine)
	})

	s.Run("TestSearchPostStore", func() {
		searchtest.TestSearchPostStore(s.T(), s.Store, searchTestEngine)
	})

	s.Run("TestSearchFileInfoStore", func() {
		searchtest.TestSearchFileInfoStore(s.T(), s.Store, searchTestEngine)
	})
}
