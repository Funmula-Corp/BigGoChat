package biggoengine

import (
	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/request"
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
	return
}

func (be *BiggoEngine) SearchPosts(channels model.ChannelList, searchParams []*model.SearchParams, page, perPage int) (postIds []string, matches model.PostSearchMatches, aErr *model.AppError) {
	return
}
