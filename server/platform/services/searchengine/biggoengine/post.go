package biggoengine

import (
	"bytes"
	"context"
	"encoding/json"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/platform/services/searchengine/biggoengine/clients"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (be *BiggoEngine) DeleteChannelPosts(rctx request.CTX, channelID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeletePost(post *model.Post) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeleteUserPosts(rctx request.CTX, userID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexPost(post *model.Post, teamId string) (aErr *model.AppError) {
	return be.IndexPostsBulk([]*model.PostForIndexing{{
		Post:   *post.Clone(),
		TeamId: teamId,
	}})
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

	for _, post := range posts {
		// potentially improve with bulk CALL via UNWIND []post
		if _, err = clients.GraphQuery(`
			MERGE (p:post{post_id:$post_id})
			MERGE (c:channel{channel_id:$channel_id})
			MERGE (p)-[:in_channel]->(c)
		`, map[string]interface{}{
			"post_id":    post.Id,
			"channel_id": post.ChannelId,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
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
	}
	return
}

func (be *BiggoEngine) SearchPosts(channels model.ChannelList, searchParams []*model.SearchParams, page, perPage int) (postIds []string, matches model.PostSearchMatches, aErr *model.AppError) {
	mlog.Debug("BiggoIndexer", mlog.Any("channels", channels), mlog.Any("searchParams", searchParams), mlog.Int("page", page), mlog.Int("perPage", perPage))
	var (
		err error
		res *neo4j.EagerResult
	)
	if res, err = clients.GraphQuery(`
		CALL apoc.es.query($es_address, $es_index, '_doc', null, {
			fields: ['_id'],
			query: {
				bool: {
					must: [
						{
							match_phrase: {
								message: $term
							}
						},
						{
							exists: {
								field: 'type'
							}
						}
					],
					must_not: [
						{
							terms: {
								type: [
									'system_join_team',
									'system_join_channel'
								]
							}
						}
					]
				} 
			}, from: 0, size: $size
		}) YIELD value
		UNWIND value.hits.hits AS hit
		RETURN hit._id AS id;
	`, map[string]interface{}{
		"es_address": "http://172.17.0.1:9200",
		"es_index":   EsPostIndex,
		"term":       searchParams[0].Terms,
		"size":       perPage,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
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
