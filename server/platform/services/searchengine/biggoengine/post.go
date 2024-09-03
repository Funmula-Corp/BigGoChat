package biggoengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/helper"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	deleteUserPostsQuery string = `{
		"query": {
			"match": {
				"user_id": "%s"
			}
		},
		"script": {
			"source": "ctx._source['delete_at'] = %dL",
			"lang": "painless"
		}
	}`
	deleteChannelPostsQuery string = `{
		"query": {
			"match": {
				"channel_id": "%s"
			}
		},
		"script": {
			"source": "ctx._source['delete_at'] = %dL",
			"lang": "painless"
		}
	}`
	indexPostBulkQuery string = `
		UNWIND $posts AS kvp
			MERGE (p:post{post_id:kvp.post_id})
				ON CREATE SET p = {post_id:kvp.post_id}
				ON MATCH SET p += {post_id:kvp.post_id}
			MERGE (c:channel{channel_id:kvp.channel_id})
				ON CREATE SET c = {channel_id:kvp.channel_id}
				ON MATCH SET c = c
			MERGE (u:user{user_id:kvp.user_id})
				ON CREATE SET u = {user_id:kvp.user_id}
				ON MATCH SET u = u
			MERGE (p)-[:in_channel]->(c)
			MERGE (p)-[:posted_by]->(u)
	`
)

func (be *BiggoEngine) DeleteChannelPosts(rctx request.CTX, channelID string) (aErr *model.AppError) {
	var (
		client *elasticsearch.Client
		err    error
		res    *esapi.Response
	)

	if client, err = clients.EsClient(); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	request := esapi.UpdateByQueryRequest{
		Index: []string{EsPostIndex}, Body: strings.NewReader(fmt.Sprintf(deleteChannelPostsQuery, channelID, model.GetMillis())),
	}

	if res, err = request.Do(context.Background(), client); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("channelID", channelID))
		}
	}
	return
}

func (be *BiggoEngine) DeletePost(post *model.Post) (aErr *model.AppError) {
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
	if buffer, err = json.Marshal(post); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	if res, err = client.Update(EsPostIndex, post.Id, bytes.NewBuffer(buffer)); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("post", post))
		}
	}
	return
}

func (be *BiggoEngine) DeleteUserPosts(rctx request.CTX, userID string) (aErr *model.AppError) {
	var (
		client *elasticsearch.Client
		err    error
		res    *esapi.Response
	)

	if client, err = clients.EsClient(); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	request := esapi.UpdateByQueryRequest{
		Index: []string{EsPostIndex}, Body: strings.NewReader(fmt.Sprintf(deleteUserPostsQuery, userID, model.GetMillis())),
	}

	if res, err = request.Do(context.Background(), client); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("userID", userID))
		}
	}
	return
}

func (be *BiggoEngine) InitializePostIndex() {
	if clients.CheckIndexExists(EsPostIndex) {
		return
	}

	settings := `{
		"settings": {
			"analysis": {
				"tokenizer": {
					"mm_search_index_analyzer": {
						"type": "pattern",
						"pattern": "[\\._/\\s]"
					}
				},
				"analyzer": {
					"mm_search_index_analyzer": {
						"type": "custom",
						"tokenizer": "mm_search_index_analyzer"
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"message": {
					"type": "text",
					"analyzer": "mm_search_index_analyzer"
				}
			}
		}
	}`

	client, _ := clients.EsClient()
	if res, err := client.Indices.Create(EsPostIndex,
		client.Indices.Create.WithBody(strings.NewReader(settings)),
	); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
	} else {
		defer res.Body.Close()
		if res.StatusCode > 400 {
			var buffer []byte
			if buffer, err = io.ReadAll(res.Body); err != nil {
				mlog.Error("BiggoIndexer", mlog.Err(err))
				return
			}
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))))
		}
	}
}

func (be *BiggoEngine) IndexPost(post *model.Post, teamId string) (aErr *model.AppError) {
	return be.IndexPostsBulk([]*model.PostForIndexing{{Post: *post.Clone(), TeamId: teamId}})
}

func (be *BiggoEngine) IndexPostsBulk(posts []*model.PostForIndexing) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	if indexer, err = clients.EsBulkIndex(EsPostIndex); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	postsMap := []map[string]string{}
	for _, post := range posts {
		// index only user / bot send messages - ignore all system messages -> [server/public/model/post.go] constants
		if post.Type != "" {
			continue
		}
		var buffer []byte
		if buffer, err = json.Marshal(post); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: post.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      EsPostIndex,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		postsMap = append(postsMap, map[string]string{
			"channel_id": post.ChannelId,
			"post_id":    post.Id,
			"user_id":    post.UserId,
		})
	}

	if _, err = clients.GraphQuery(indexPostBulkQuery, map[string]interface{}{
		"posts": postsMap,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
	}
	return
}

func (be *BiggoEngine) SearchPosts(channels model.ChannelList, searchParams []*model.SearchParams, page, perPage int) (postIds []string, matches model.PostSearchMatches, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)

	if len(searchParams) <= 0 {
		// no search parameters provided
		return
	}

	query, queryParams := helper.ComposeSearchParamsQuery(be.config, EsPostIndex, page, perPage, "message", searchParams[0])
	mlog.Debug("BiggoIndexer", mlog.String("posts_query", query), mlog.Any("posts_query_params", queryParams))
	if res, err = clients.GraphQuery(fmt.Sprintf(`%s RETURN hit._id AS id;`, query), queryParams); err != nil {
		mlog.Error("BiggoIndexer", mlog.String("function", "SearchPosts"), mlog.Err(err), mlog.String("query", query), mlog.Any("query_params", queryParams))
		return
	}

	postIds = []string{}
	matches = model.PostSearchMatches{}
	for _, record := range res.Records {
		entry := record.AsMap()
		matches[entry["id"].(string)] = []string{entry["id"].(string)}
		postIds = append(postIds, entry["id"].(string))
	}
	return
}
