package common

import (
	"errors"

	"git.biggo.com/Funmula/BigGoChat/server/public/pluginapi"
)

var ErrNotFound = errors.New("not found")

type KVStore interface {
	Set(key string, value interface{}, options ...pluginapi.KVSetOption) (bool, error)
	Get(key string, o interface{}) error
	Delete(key string) error
	DeleteAll() error
	ListKeys(page, count int, options ...pluginapi.ListKeysOption) ([]string, error)
}
