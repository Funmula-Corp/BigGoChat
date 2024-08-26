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

func (be *BiggoEngine) DeleteUser(user *model.User) (aErr *model.AppError) {
	if _, err := clients.GraphQuery(`
		MATCH (u:user{user_id:user_id})
		MATCH (u)-[r]->()
		DELETE r,u
		RETURN COUNT(u);
	`, map[string]interface{}{
		"user_id": user.Id,
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

	if res, err = client.Delete(EsUserIndex, user.Id); err != nil {
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

	if indexer, err = clients.EsBulkIndex(EsUserIndex); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	for _, user := range users {
		// potentially improve with bulk CALL via UNWIND []user
		if _, err = clients.GraphQuery(`
			MERGE (u:user{user_id:$user_id})
			WITH u
			UNWIND $channel_ids AS channel_id
				MERGE (c:channel{channel_id:channel_id})
				MERGE (u)-[:channel_member]->(c)
			WITH u
			UNWIND $team_ids AS team_id
				MERGE (t:team{team_id:team_id})
				MERGE (u)-[:team_member]->(t)
		`, map[string]interface{}{
			"channel_ids": user.ChannelsIds,
			"team_ids":    user.TeamsIds,
			"user_id":     user.Id,
		}); err != nil {
			//aErr = model.NewAppError("BiggoIndexer.IndexUsersBulk", "engine.biggo.indexer.index_user_bulk.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		var buffer []byte
		if buffer, err = json.Marshal(user); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: user.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      EsUserIndex,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}
	}
	return
}

func (be *BiggoEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) (userInChannel []string, userNotInChannel []string, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)
	if res, err = clients.GraphQuery(`
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
	`, map[string]interface{}{
		"es_address": "http://172.17.0.1:9200",
		"es_index":   EsUserIndex,
		"channel_id": channelId,
		"term":       term,
		"size":       25,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
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
	if res, err = clients.GraphQuery(`
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
	`, map[string]interface{}{
		"es_address": "http://172.17.0.1:9200",
		"es_index":   EsUserIndex,
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
