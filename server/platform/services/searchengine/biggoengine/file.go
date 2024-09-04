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
	deletePostFilesQuery string = `{
		"query": {
			"match": {
				"post_id": "%s"
			}
		},
		"script": {
			"source": "ctx._source['delete_at'] = %dL",
			"lang": "painless"
		}
	}`
	deleteFilesBatchQuery string = `{
		"query": {
			"range": {
				"create_at": {
					"gte": %d
				}
			}
		}, "size": %d,
		"script": {
			"source": "ctx._source['delete_at'] = %dL",
			"lang": "painless"
		}
	}`
	deleteUserFilesQuery string = `{
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
	indexFilesBulkQuery string = `
		UNWIND $files AS kvp
			MERGE (f:file{file_id:kvp.file_id})
				ON CREATE SET f = {file_id:kvp.file_id}
				ON MATCH SET f += {file_id:kvp.file_id}
			MERGE (u:user{user_id:kvp.user_id})
				ON CREATE SET u = {user_id:kvp.user_id}
				ON MATCH SET u = u
			MERGE (p:post{post_id:kvp.post_id})
				ON CREATE SET p = {post_id:kvp.post_id}
				ON MATCH SET p = p
			MERGE (c:channel{channel_id:kvp.channel_id})
				ON CREATE SET c = {channel_id:kvp.channel_id}
				ON MATCH SET c = c
			MERGE (f)-[:in_channel]->(c)
			MERGE (f)-[:by_user]->(u)
			MERGE (f)-[:in_post]->(p)
	`
)

func (be *BiggoEngine) DeleteFile(fileID string) (aErr *model.AppError) {
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
	if buffer, err = json.Marshal(&model.FileInfo{
		Id: fileID, DeleteAt: model.GetMillis(),
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}

	if res, err = client.Update(EsFileIndex, fileID, bytes.NewBuffer(buffer)); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("fileID", fileID))
		}
	}
	return
}

func (be *BiggoEngine) DeleteFilesBatch(rctx request.CTX, endTime, limit int64) (aErr *model.AppError) {
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
		Index: []string{EsFileIndex}, Body: strings.NewReader(fmt.Sprintf(deleteFilesBatchQuery, endTime, limit, model.GetMillis())),
	}

	if res, err = request.Do(context.Background(), client); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("endTime", endTime), mlog.Any("limit", limit))
		}
	}
	return
}

func (be *BiggoEngine) DeletePostFiles(rctx request.CTX, postID string) (aErr *model.AppError) {
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
		Index: []string{EsFileIndex}, Body: strings.NewReader(fmt.Sprintf(deletePostFilesQuery, postID, model.GetMillis())),
	}

	if res, err = request.Do(context.Background(), client); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		if buffer, err := io.ReadAll(res.Body); err == nil {
			mlog.Error("BiggoIndexer", mlog.Err(errors.New(string(buffer))), mlog.Any("postID", postID))
		}
	}
	return
}

func (be *BiggoEngine) DeleteUserFiles(rctx request.CTX, userID string) (aErr *model.AppError) {
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
		Index: []string{EsFileIndex}, Body: strings.NewReader(fmt.Sprintf(deleteUserFilesQuery, userID, model.GetMillis())),
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

func (be *BiggoEngine) InitializeFilesIndex() {
	var index = be.GetEsFileIndex()
	if clients.CheckIndexExists(index) {
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
				"name": {
					"type": "text",
					"analyzer": "mm_search_index_analyzer"
				}
			}
		}
	}`

	client, _ := clients.EsClient()
	if res, err := client.Indices.Create(index,
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

func (be *BiggoEngine) IndexFile(file *model.FileInfo, channelId string) (aErr *model.AppError) {
	return be.IndexFilesBulk([]*model.FileForIndexing{{FileInfo: *file, ChannelId: channelId}})
}

func (be *BiggoEngine) IndexFilesBulk(files []*model.FileForIndexing) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	var index = be.GetEsFileIndex()
	if indexer, err = clients.EsBulkIndex(index); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
		return
	}
	defer indexer.Close(context.Background())

	filesMap := []map[string]string{}
	for _, file := range files {
		var buffer []byte
		if buffer, err = json.Marshal(file); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		if err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: file.Id,
			Body:       bytes.NewBuffer(buffer),
			Index:      index,
		}); err != nil {
			mlog.Error("BiggoIndexer", mlog.Err(err))
			continue
		}

		filesMap = append(filesMap, map[string]string{
			"channel_id": file.ChannelId,
			"file_id":    file.Id,
			"post_id":    file.PostId,
			"user_id":    file.CreatorId,
		})
	}

	if _, err = clients.GraphQuery(indexFilesBulkQuery, map[string]interface{}{
		"files": filesMap,
	}); err != nil {
		mlog.Error("BiggoIndexer", mlog.Err(err))
	}
	return
}

func (be *BiggoEngine) SearchFiles(channels model.ChannelList, searchParams []*model.SearchParams, page, perPage int) (result []string, aErr *model.AppError) {
	var (
		err error
		res *neo4j.EagerResult
	)

	// use index_* for search and index_1 for indexing
	// create analyzer and mapping for indices before feeding
	if len(searchParams) <= 0 {
		// no search parameters provided
		return
	}

	if searchParams[0].InChannels == nil {
		searchParams[0].InChannels = []string{}
	}
	for _, c := range channels {
		searchParams[0].InChannels = append(searchParams[0].InChannels, c.Id)
	}

	query, queryParams := helper.ComposeSearchParamsQuery(be.config, EsFileIndex, page, perPage, "name", searchParams[0])
	mlog.Debug("BiggoIndexer", mlog.String("files_query", query), mlog.Any("files_query_params", queryParams))
	if res, err = clients.GraphQuery(fmt.Sprintf(`%s RETURN hit._id AS id;`, query), queryParams); err != nil {
		mlog.Error("BiggoIndexer", mlog.String("function", "SearchFiles"), mlog.Err(err), mlog.String("query", query), mlog.Any("query_params", queryParams))
		return
	}

	result = []string{}
	for _, record := range res.Records {
		entry := record.AsMap()
		result = append(result, entry["id"].(string))
	}
	return
}
