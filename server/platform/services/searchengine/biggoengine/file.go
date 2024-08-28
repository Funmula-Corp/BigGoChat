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
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
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
			MERGE (u:user{user_id:kvp.user_id})
			MERGE (p:post{post_id:kvp.post_id})
			MERGE (c:channel{channel_id:kvp.channel_id})
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

func (be *BiggoEngine) IndexFile(file *model.FileInfo, channelId string) (aErr *model.AppError) {
	return be.IndexFilesBulk([]*model.FileForIndexing{{FileInfo: *file, ChannelId: channelId}})
}

func (be *BiggoEngine) IndexFilesBulk(files []*model.FileForIndexing) (aErr *model.AppError) {
	var (
		indexer esutil.BulkIndexer
		err     error
	)

	if indexer, err = clients.EsBulkIndex(EsPostIndex); err != nil {
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
			Index:      EsFileIndex,
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
	return
}
