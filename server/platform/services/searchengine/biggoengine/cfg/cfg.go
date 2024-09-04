package cfg

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

type ConfigKey string

const (
	PluginName = "com.biggo.search-engine"

	keyEnableAutocomplete ConfigKey = "enable_autocompletion"
	keyEnableIndexer      ConfigKey = "enable_indexer"
	keyEnableIndexing     ConfigKey = "enable_indexing"
	keyEnableSearch       ConfigKey = "enable_search"

	keyNeo4jProtocol       ConfigKey = "neo4j_protocol"
	keyNeo4jHost           ConfigKey = "neo4j_host"
	keyNeo4jPort           ConfigKey = "neo4j_port"
	keyNeo4jUseCredentials ConfigKey = "neo4j_use_credentials"
	keyNeo4jUsername       ConfigKey = "neo4j_username"
	keyNeo4jPassword       ConfigKey = "neo4j_password"

	keyElasticsearchProtocol       ConfigKey = "elasticsearch_protocol"
	keyElasticsearchHost           ConfigKey = "elasticsearch_host"
	keyElasticsearchPort           ConfigKey = "elasticsearch_port"
	keyElasticsearchUseCredentials ConfigKey = "elasticsearch_use_credentials"
	keyElasticsearchUsername       ConfigKey = "elasticsearch_username"
	keyElasticsearchPassword       ConfigKey = "elasticsearch_password"
)

func EnableAutocomplete(cfg *model.Config) bool {
	return getBool(cfg, keyEnableAutocomplete, false)
}

func EnableIndexer(cfg *model.Config) bool {
	return getBool(cfg, keyEnableIndexer, false)
}

func EnableIndexing(cfg *model.Config) bool {
	return getBool(cfg, keyEnableIndexing, false)
}

func EnableSearch(cfg *model.Config) bool {
	return getBool(cfg, keyEnableSearch, false)
}

func Neo4jProtocol(cfg *model.Config) string {
	return getString(cfg, keyNeo4jProtocol, "bolt")
}

func Neo4jHost(cfg *model.Config) string {
	return getString(cfg, keyNeo4jHost, "localhost")
}

func Neo4jPort(cfg *model.Config) float64 {
	return getFloat64(cfg, keyNeo4jPort, 7687)
}

func Neo4jUseCredentials(cfg *model.Config) bool {
	return getBool(cfg, keyNeo4jUseCredentials, false)
}

func Neo4jUsername(cfg *model.Config) string {
	return getString(cfg, keyNeo4jUsername, "")
}

func Neo4jPassword(cfg *model.Config) string {
	return getString(cfg, keyNeo4jPassword, "")
}

func ElasticsearchProtocol(cfg *model.Config) string {
	return getString(cfg, keyElasticsearchProtocol, "http")
}

func ElasticsearchHost(cfg *model.Config) string {
	return getString(cfg, keyElasticsearchHost, "localhost")
}

func ElasticsearchPort(cfg *model.Config) float64 {
	return getFloat64(cfg, keyElasticsearchPort, 9200)
}

func ElasticsearchUseCredentials(cfg *model.Config) bool {
	return getBool(cfg, keyElasticsearchUseCredentials, false)
}

func ElasticsearchUsername(cfg *model.Config) string {
	return getString(cfg, keyElasticsearchUsername, "")
}

func ElasticsearchPassword(cfg *model.Config) string {
	return getString(cfg, keyElasticsearchPassword, "")
}

func ElasticsearchIndexChannelSuffix(cfg *model.Config) string {
	return getString(cfg, "elasticsearch_index_channel_suffix", "0")
}

func ElasticsearchIndexFileSuffix(cfg *model.Config) string {
	return getString(cfg, "elasticsearch_index_file_suffix", "0")
}

func ElasticsearchIndexPostSuffix(cfg *model.Config) string {
	return getString(cfg, "elasticsearch_index_post_suffix", "0")
}

func ElasticsearchIndexUserSuffix(cfg *model.Config) string {
	return getString(cfg, "elasticsearch_index_user_suffix", "0")
}

func getBool(cfg *model.Config, key ConfigKey, fallback bool) (value bool) {
	value = fallback
	if cfg != nil {
		if settings, ok := cfg.PluginSettings.Plugins[PluginName]; ok {
			if enabled, found := settings[string(key)]; found {
				return enabled.(bool)
			}
		}
	}
	return
}

func getFloat64(cfg *model.Config, key ConfigKey, fallback float64) (value float64) {
	value = fallback
	if cfg != nil {
		if settings, ok := cfg.PluginSettings.Plugins[PluginName]; ok {
			if enabled, found := settings[string(key)]; found {
				return enabled.(float64)
			}
		}
	}
	return
}

func getString(cfg *model.Config, key ConfigKey, fallback string) (value string) {
	value = fallback
	if cfg != nil {
		if settings, ok := cfg.PluginSettings.Plugins[PluginName]; ok {
			if enabled, found := settings[string(key)]; found {
				return enabled.(string)
			}
		}
	}
	return
}
