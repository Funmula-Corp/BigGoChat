package biggoengine

import (
	"fmt"
	"sync/atomic"
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

const (
	EngineName = "biggo"
	PluginName = "com.biggo.biggo-engine"

	PostIndex    = "post"
	FileIndex    = "file"
	UserIndex    = "user"
	ChannelIndex = "channel"

	EsChannelIndex string = "mm_biggoengine_channel"
	EsPostIndex    string = "mm_biggoengine_post"
	EsUserIndex    string = "mm_biggoengine_user"
)

func init() {}

func NewBiggoEngine(config *model.Config) *BiggoEngine {
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
	if settings, ok := be.config.PluginSettings.Plugins[PluginName]; ok {
		if enabled, found := settings["autocompletion_enabled"]; found {
			return enabled.(bool)
		}
	}
	return
}

func (be *BiggoEngine) IsChannelsIndexVerified() (result bool) {
	return
}

func (be *BiggoEngine) IsEnabled() (result bool) {
	if settings, ok := be.config.PluginSettings.Plugins[PluginName]; ok {
		if enabled, found := settings["enabled"]; found {
			return enabled.(bool)
		}
	}
	return
}

func (be *BiggoEngine) IsIndexingEnabled() (result bool) {
	if settings, ok := be.config.PluginSettings.Plugins[PluginName]; ok {
		if enabled, found := settings["indexing_enabled"]; found {
			return enabled.(bool)
		}
	}
	return
}

func (be *BiggoEngine) IsIndexingSync() (result bool) {
	return be.isSync.Load()
}

func (be *BiggoEngine) IsSearchEnabled() (result bool) {
	if settings, ok := be.config.PluginSettings.Plugins[PluginName]; ok {
		if enabled, found := settings["search_enabled"]; found {
			return enabled.(bool)
		}
	}
	return
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
	return nil
}

func (be *BiggoEngine) UpdateConfig(config *model.Config) {
	be.config = config
}
