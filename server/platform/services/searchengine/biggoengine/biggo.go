package biggoengine

import (
	"fmt"
	"sync/atomic"
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"
)

const (
	EngineName = "biggo"

	ChannelVertex = "channel"
	FileVertex    = "file"
	PostVertex    = "post"
	UserVertex    = "user"

	EsChannelIndex = "mm_biggoengine_channel_*"
	EsFileIndex    = "mm_biggoengine_file_*"
	EsPostIndex    = "mm_biggoengine_post_*"
	EsUserIndex    = "mm_biggoengine_user_*"
)

func (be *BiggoEngine) GetEsChannelIndex() string {
	return fmt.Sprintf("mm_biggoengine_channel_%s", cfg.ElasticsearchIndexChannelSuffix(be.config))
}

func (be *BiggoEngine) GetEsFileIndex() string {
	return fmt.Sprintf("mm_biggoengine_file_%s", cfg.ElasticsearchIndexFileSuffix(be.config))
}

func (be *BiggoEngine) GetEsPostIndex() string {
	return fmt.Sprintf("mm_biggoengine_post_%s", cfg.ElasticsearchIndexPostSuffix(be.config))
}

func (be *BiggoEngine) GetEsUserIndex() string {
	return fmt.Sprintf("mm_biggoengine_user_%s", cfg.ElasticsearchIndexUserSuffix(be.config))
}

func NewBiggoEngine(config *model.Config) *BiggoEngine {
	clients.InitEsClient(config)
	clients.InitNeo4jClient(config)
	return &BiggoEngine{
		config: config,
	}
}

type BiggoEngine struct {
	config *model.Config

	isActive atomic.Bool
	isSync   atomic.Bool
}

func (be *BiggoEngine) DataRetentionDeleteIndexes(rctx request.CTX, cutoff time.Time) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) GetFullVersion() (result string) {
	return fmt.Sprintf("%d", be.GetVersion())
}

func (be *BiggoEngine) GetName() (result string) {
	return EngineName
}

func (be *BiggoEngine) GetPlugins() (result []string) {
	return []string{}
}

func (be *BiggoEngine) GetVersion() (result int) {
	return 0
}

func (be *BiggoEngine) IsActive() (result bool) {
	return be.isActive.Load()
}

func (be *BiggoEngine) IsAutocompletionEnabled() (result bool) {
	return cfg.EnableAutocomplete(be.config)
}

func (be *BiggoEngine) IsChannelsIndexVerified() (result bool) {
	return true
}

func (be *BiggoEngine) IsEnabled() (result bool) {
	return cfg.EnableIndexer(be.config)
}

func (be *BiggoEngine) IsIndexingEnabled() (result bool) {
	return cfg.EnableIndexing(be.config)
}

func (be *BiggoEngine) IsIndexingSync() (result bool) {
	return be.isSync.Load()
}

func (be *BiggoEngine) IsSearchEnabled() (result bool) {
	return cfg.EnableSearch(be.config)
}

func (be *BiggoEngine) PurgeIndexes(rctx request.CTX) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) PurgeIndexList(rctx request.CTX, indexes []string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) RefreshIndexes(rctx request.CTX) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) Start() (aErr *model.AppError) {
	be.isActive.Store(true)
	return
}

func (be *BiggoEngine) Stop() (aErr *model.AppError) {
	be.isActive.Store(false)
	return
}

func (be *BiggoEngine) TestConfig(rctx request.CTX, cfg *model.Config) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) UpdateConfig(config *model.Config) {
	be.config = config
}
