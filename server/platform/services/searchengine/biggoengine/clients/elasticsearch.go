package clients

import (
	"encoding/json"
	"runtime"
	"time"

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

func EsClient() (client *elasticsearch.Client, err error) {
	if esClient == nil {
		if esClient, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{
				"http://localhost:9200",
			},

			APIKey:   "",
			Username: "",
			Password: "",

			MaxRetries: 5,
			RetryBackoff: func(attempt int) time.Duration {
				return time.Second * time.Duration(attempt)
			},
			RetryOnStatus: []int{
				429,
				502, 503, 504,
			},
		}); err != nil {
			return
		}
	}
	client = esClient
	return
}

func EsBulkIndex(index string) (indexer esutil.BulkIndexer, err error) {
	if _, err = EsClient(); err != nil {
		return
	}
	indexer, err = esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        esClient,
		Index:         index,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    10 * 1024 * 1024,
		FlushInterval: 30 * time.Second,
	})
	return
}
