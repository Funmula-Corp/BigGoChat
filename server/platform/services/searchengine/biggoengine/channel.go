package biggoengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (be *BiggoEngine) DeleteChannel(channel *model.Channel) (aErr *model.AppError) {
	if _, err := clients.GraphQuery(`
		MATCH (c:channel{channel_id:$channel_id})
		MATCH (c)-[r1]->()
		MATCH ()-[r2]->(c)
		MATCH (c)-[r3]->()
		DELETE r1,r2,r3,c;
	`, map[string]interface{}{
		"channel_id": channel.Id,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	var (
		client *elasticsearch.Client
		err    error
		res    *esapi.Response
	)
	if client, err = clients.EsClient(); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	if res, err = client.Delete(EsChannelIndex, channel.Id); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))))
		}
		return
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

	for _, channel := range channels {
		// potentially improve with bulk CALL via UNWIND []channel
		if _, err = clients.GraphQuery(`
			MERGE (c:channel{channel_id:$channel_id})
			MERGE (t:team{team_id:$team_id})
			MERGE (c)-[:in_team]->(t)
		`, map[string]interface{}{
			"channel_id": channel.Id,
			"team_id":    channel.TeamId,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

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
		"es_address": "http://172.17.0.1:9200",
		"es_index":   EsChannelIndex,
		"term":       term,
		"size":       25,
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
