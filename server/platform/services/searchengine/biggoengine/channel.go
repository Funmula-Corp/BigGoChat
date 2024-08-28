package biggoengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/cfg"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	indexChannelBulkQuery string = `
		UNWIND $channels AS kvp
			MERGE (c:channel{channel_id:kvp.channel_id})
			MERGE (t:team{team_id:kvp.team_id})
			MERGE (c)-[:in_team]->(t)
	`
)

func (be *BiggoEngine) DeleteChannel(channel *model.Channel) (aErr *model.AppError) {
	var (
		client *elasticsearch.Client
		err    error
		res    *esapi.Response
	)

	if client, err = clients.EsClient(); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	var buffer []byte
	if buffer, err = json.Marshal(channel); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	if res, err = client.Update(EsChannelIndex, channel.Id, bytes.NewBuffer(buffer)); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("channel", channel))
		}
	}
	return
}

func (be *BiggoEngine) IndexChannel(rctx request.CTX, channel *model.Channel, userIDs, teamMemberIDs []string) (aErr *model.AppError) {
	return be.IndexChannelsBulk([]*model.Channel{channel})
}

func (be *BiggoEngine) IndexChannelsBulk(channels []*model.Channel) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	if indexer, err = clients.EsBulkIndex(EsChannelIndex); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	channelsMap := []map[string]string{}
	for _, channel := range channels {
		var buffer []byte
		if buffer, err = json.Marshal(channel); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: channel.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      EsChannelIndex,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		channelsMap = append(channelsMap, map[string]string{
			"channel_id": channel.Id,
			"team_id":    channel.TeamId,
		})
	}

	if _, err = clients.GraphQuery(indexChannelBulkQuery, map[string]interface{}{
		"channels": channelsMap,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
	}
	return
}

func (be *BiggoEngine) SearchChannels(teamId, userID, term string, isGuest bool) (result []string, aErr *model.AppError) {
	mlog.Debug("BiggoIndexer", mlog.String("teamId", teamId), mlog.String("userID", userID), mlog.String("term", term), mlog.Bool("isGuest", isGuest))
	var (
		err error
		res *neo4j.EagerResult
	)
	if res, err = clients.GraphQuery(`
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			fields: ['_id'],
			query: {
				prefix: {
					display_name: {
						value: $term
					}
				}
			}, from: 0, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit
		RETURN hit._id as id
	`, map[string]interface{}{
		"es_address": fmt.Sprintf("%s://%s:%.0f",
			cfg.ElasticsearchProtocol(be.config),
			cfg.ElasticsearchHost(be.config),
			cfg.ElasticsearchPort(be.config),
		),
		"es_index": EsChannelIndex,
		"term":     term,
		"size":     25,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	result = []string{}
	for _, record := range res.Records {
		entry := record.AsMap()
		result = append(result, entry["id"].(string))
	}
	return
}
