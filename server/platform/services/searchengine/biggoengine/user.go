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
	indexUsersBulkQuery string = `
		UNWIND $users AS kvp
			MERGE (u:user{user_id:kvp.user_id})
			WITH kvp, u
			UNWIND kvp.channel_ids AS channel_id
				MERGE (c:channel{channel_id:channel_id})
					ON CREATE SET c = {channel_id:channel_id}
					ON MATCH SET c = c
				MERGE (u)-[:channel_member]->(c)
			WITH kvp, u
			UNWIND kvp.team_ids AS team_id
				MERGE (t:team{team_id:team_id})
					ON CREATE SET t = {team_id:team_id}
					ON MATCH SET t = t
				MERGE (u)-[:team_member]->(t)
	`
)

func (be *BiggoEngine) DeleteUser(user *model.User) (aErr *model.AppError) {
	var (
		client *elasticsearch.Client
		err    error
		res    *esapi.Response
	)

	if client, err = clients.EsClient(); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "DeleteUser"), mlog.Err(err))
		return
	}

	update := map[string]interface{}{
		"doc": map[string]interface{}{
			"delete_at": user.DeleteAt,
		},
	}

	var buffer []byte
	if buffer, err = json.Marshal(update); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "DeleteUser"), mlog.Err(err))
		return
	}

	if res, err = client.Update(EsUserIndex, user.Id, bytes.NewReader(buffer)); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "DeleteUser"), mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "DeleteUser"), mlog.Err(errors.New(string(buffer))), mlog.Any("user", user))
		}
	}
	return
}

func (be *BiggoEngine) InitializeUserIndex() {}

func (be *BiggoEngine) IndexUser(rctx request.CTX, user *model.User, teamsIds, channelsIds []string) (aErr *model.AppError) {
	var mobilephone string
	if user.Mobilephone != nil {
		mobilephone = *user.Mobilephone
	}
	return be.IndexUsersBulk([]*model.UserForIndexing{{
		Id:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Roles:       user.Roles,
		CreateAt:    user.CreateAt,
		DeleteAt:    user.DeleteAt,
		TeamsIds:    teamsIds,
		ChannelsIds: channelsIds,
		Mobilephone: mobilephone,
	}})
}

func (be *BiggoEngine) IndexUsersBulk(users []*model.UserForIndexing) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	var index = be.GetEsUserIndex()
	if indexer, err = clients.EsBulkIndex(index); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "IndexUser"), mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	usersMap := []map[string]any{}
	for _, user := range users {
		var buffer []byte
		if buffer, err = json.Marshal(user); err != nil {
			mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "IndexUser"), mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: user.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      index,
		}); err != nil {
			mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "IndexUser"), mlog.Err(err))
			continue
		}

		usersMap = append(usersMap, map[string]any{
			"channel_ids": user.ChannelsIds,
			"team_ids":    user.TeamsIds,
			"user_id":     user.Id,
		})
	}

	if _, err = clients.GraphQuery(indexUsersBulkQuery, map[string]interface{}{
		"users": usersMap,
	}); err != nil {
		//aErr = model.NewAppError("BiggoEngine., mlog.String("component", "user"), mlog.String("func_name", )IndexUsersBulk", "engine.biggo.indexer.index_user_bulk.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "IndexUser"), mlog.Err(err))
	}
	return
}

func (be *BiggoEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)
	query := `
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			fields: ['_id'],
			query: {
				prefix: {
					username: {
						value: $term
					}
				}
			}, from: 0, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit
		OPTIONAL MATCH (:user{user_id:hit._id})-[r:channel_member]->(:channel{channel_id:$channel_id})
		RETURN hit._id as id, r IS NOT NULL AS in_channel
	`
	queryParams := map[string]any{
		"es_address": fmt.Sprintf("%s://%s:%.0f",
			cfg.ElasticsearchProtocol(be.config),
			cfg.ElasticsearchHost(be.config),
			cfg.ElasticsearchPort(be.config),
		),
		"es_index":   EsUserIndex,
		"channel_id": channelId,
		"term":       term,
		"size":       options.Limit,
	}
	if res, err = clients.GraphQuery(query, queryParams); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "SearchUsersInChannel"), mlog.Err(err), mlog.String("query", query), mlog.Any("query_params", queryParams))
		return
	}

	userInChannel = []string{}
	userNotInChannel = []string{}
	for _, record := range res.Records {
		entry := record.AsMap()
		if entry["in_channel"].(bool) {
			userInChannel = append(userInChannel, entry["id"].(string))
		} else {
			userNotInChannel = append(userNotInChannel, entry["id"].(string))
		}
	}
	return
}

func (be *BiggoEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (result []string, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)
	query := `
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			fields: ['_id'],
			query: {
				prefix: {
					username: {
						value: $term
					}
				}
			}, from: 0, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit
		RETURN hit._id as id
	`
	queryParams := map[string]any{
		"es_address": fmt.Sprintf("%s://%s:%.0f",
			cfg.ElasticsearchProtocol(be.config),
			cfg.ElasticsearchHost(be.config),
			cfg.ElasticsearchPort(be.config),
		),
		"es_index": EsUserIndex,
		"term":     term,
		"size":     options.Limit,
	}
	if res, err = clients.GraphQuery(query, queryParams); err != nil {
		mlog.Error("BiggoEngine", mlog.String("component", "user"), mlog.String("func_name", "SearchUsersInTeam"), mlog.Err(err), mlog.String("query", query), mlog.Any("query_params", queryParams))
		return
	}

	result = []string{}
	for _, record := range res.Records {
		entry := record.AsMap()
		result = append(result, entry["id"].(string))
	}
	return
}
