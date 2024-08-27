package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var esClient *elasticsearch.Client = nil

type EnvelopeResponse struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []struct {
			ID         string          `json:"_id"`
			Source     json.RawMessage `json:"_source"`
			Highlights json.RawMessage `json:"highlight"`
			Sort       []interface{}   `json:"sort"`
		}
	}
}

func InitEsClient(config *model.Config) (err error) {
	if esClient == nil {
		esClient, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{
				fmt.Sprintf("%s://%s:%0.f",
					cfg.ElasticsearchProtocol(config),
					cfg.ElasticsearchHost(config),
					cfg.ElasticsearchPort(config),
				),
			},

			APIKey:   "",
			Username: cfg.ElasticsearchUsername(config),
			Password: cfg.ElasticsearchUsername(config),

			MaxRetries: 5,
			RetryBackoff: func(attempt int) time.Duration {
				return time.Second * time.Duration(attempt)
			},
			RetryOnStatus: []int{
				429, 502, 503, 504,
			},
		})
	}
	return
}

func EsClient() (client *elasticsearch.Client, err error) {
	if client = esClient; client == nil {
		err = errors.New("elasticsearch client not initialized")
	}
	return
}

func EsBulkIndex(index string) (indexer esutil.BulkIndexer, err error) {
	indexer, err = esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        esClient,
		Index:         index,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    10 * 1024 * 1024,
		FlushInterval: 30 * time.Second,
	})
	return
}
