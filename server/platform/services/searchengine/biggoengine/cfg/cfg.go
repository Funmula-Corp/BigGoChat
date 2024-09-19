package cfg

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

type ConfigKey string

const (
	PluginName = "com.biggo.search-engine"

	keyEnableAutocomplete ConfigKey = "enable_autocompletion"
	keyEnableIndexer      ConfigKey = "enable_indexer"
	keyEnableSearch       ConfigKey = "enable_search"

	keySerchServiceHost ConfigKey = "search_service_host"
)

func EnableAutocomplete(cfg *model.Config) bool {
	return getBool(cfg, keyEnableAutocomplete, false)
}

func EnableIndexer(cfg *model.Config) bool {
	return getBool(cfg, keyEnableIndexer, false)
}

func EnableSearch(cfg *model.Config) bool {
	return getBool(cfg, keyEnableSearch, false)
}

func SearchServiceHost(cfg *model.Config) string {
	return getString(cfg, keySerchServiceHost, "")
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

/*
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
*/

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
