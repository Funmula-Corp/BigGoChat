package biggoengine

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
)

func (be *BiggoEngine) DeleteFile(fileID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeleteFilesBatch(rctx request.CTX, endTime, limit int64) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeletePostFiles(rctx request.CTX, postID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) DeleteUserFiles(rctx request.CTX, userID string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) IndexFile(file *model.FileInfo, channelId string) (aErr *model.AppError) {
	return
}

func (be *BiggoEngine) SearchFiles(channels model.ChannelList, searchParams []*model.SearchParams, page, perPage int) (result []string, aErr *model.AppError) {
	return
}
