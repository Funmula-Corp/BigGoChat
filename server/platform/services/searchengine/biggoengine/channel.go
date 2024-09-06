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
				ON CREATE SET c = {channel_id:kvp.channel_id,type:kvp.type}
				ON MATCH SET c += {channel_id:kvp.channel_id,type:kvp.type}
			MERGE (t:team{team_id:kvp.team_id})
				ON CREATE SET t = {team_id:kvp.team_id}
				ON MATCH SET t = t
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
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "DeleteChannel"), mlog.Err(err))
		return
	}

	update := map[string]interface{}{
		"doc": map[string]interface{}{
			"delete_at": channel.DeleteAt,
		},
	}

	var buffer []byte
	if buffer, err = json.Marshal(update); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "DeleteChannel"), mlog.Err(err))
		return
	}

	if res, err = client.Update(EsChannelIndex, channel.Id, bytes.NewReader(buffer)); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "DeleteChannel"), mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "DeleteChannel"), mlog.Err(errors.New(string(buffer))), mlog.Any("channel", channel))
		}
	}
	return
}

func (be *BiggoEngine) InitializeChannelIndex() {}

func (be *BiggoEngine) IndexChannel(rctx request.CTX, channel *model.Channel, userIDs, teamMemberIDs []string) (aErr *model.AppError) {
	return be.IndexChannelsBulk([]*model.Channel{channel})
}

func (be *BiggoEngine) IndexChannelsBulk(channels []*model.Channel) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	var index = be.GetEsChannelIndex()
	if indexer, err = clients.EsBulkIndex(index); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "IndexChannelsBulk"), mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	channelsMap := []map[string]string{}
	for _, channel := range channels {
		var buffer []byte
		if buffer, err = json.Marshal(channel); err != nil {
			mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "IndexChannelsBulk"), mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: channel.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      index,
		}); err != nil {
			mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "IndexChannelsBulk"), mlog.Err(err))
			continue
		}

		channelsMap = append(channelsMap, map[string]string{
			"channel_id": channel.Id,
			"team_id":    channel.TeamId,
			"type":       string(channel.Type),
		})
	}

	if _, err = clients.GraphQuery(indexChannelBulkQuery, map[string]interface{}{
		"channels": channelsMap,
	}); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "IndexChannelsBulk"), mlog.Err(err))
	}
	return
}

func (be *BiggoEngine) SearchChannels(teamId, userID, term string, isGuest bool) (result []string, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)
	query := `
		MATCH (u:user{user_id:$user_id})-[:channel_member]->(c:channel{type:'P'})
		WITH COLLECT(c.channel_id) AS channel_ids
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			fields: ['_id'],
			query: {
				bool: {
					should: [
						{
							bool: {
								must: [
									{ 
										match: {
											type: 'P'
										}
									},
									{ 
										prefix: {
											team_id: $team_id
										}
									},
									{
										terms: {
											id: channel_ids
										}
									},
									{
										prefix: {
											display_name: {
												value: $term
											}
										}
									}
								]
							}
						},
						{
							bool: {
								must: [
									{ 
										match: {
											type: 'O'
										}
									},
									{ 
										prefix: {
											team_id: $team_id
										}
									},
									{
										prefix: {
											display_name: {
												value: $term
											}
										}
									}
								]
							}
						}
					]
				}
			}, from: 0, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit
		RETURN hit._id as id
	`
	queryParams := map[string]interface{}{
		"es_address": fmt.Sprintf("%s://%s:%.0f",
			cfg.ElasticsearchProtocol(be.config),
			cfg.ElasticsearchHost(be.config),
			cfg.ElasticsearchPort(be.config),
		),
		"es_index": EsChannelIndex,
		"team_id":  teamId,
		"user_id":  userID,
		"term":     term,
		"size":     25,
	}
	if res, err = clients.GraphQuery(query, queryParams); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "channel"), mlog.String("func_name", "SearchChannels"), mlog.Err(err), mlog.String("query", query), mlog.Any("query_params", queryParams))
		return
	}

	result = []string{}
	for _, record := range res.Records {
		entry := record.AsMap()
		result = append(result, entry["id"].(string))
	}
	return
}
