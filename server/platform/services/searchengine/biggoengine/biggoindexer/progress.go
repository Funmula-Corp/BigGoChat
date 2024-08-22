package biggoindexer

import (
	"time"
)

type IndexingProgress struct {
	Now            time.Time
	StartAtTime    int64
	EndAtTime      int64
	LastEntityTime int64

	TotalPostsCount int64
	DonePostsCount  int64
	DonePosts       bool
	LastPostID      string

	TotalFilesCount int64
	DoneFilesCount  int64
	DoneFiles       bool
	LastFileID      string

	TotalChannelsCount int64
	DoneChannelsCount  int64
	DoneChannels       bool
	LastChannelID      string

	TotalUsersCount int64
	DoneUsersCount  int64
	DoneUsers       bool
	LastUserID      string
}

func (ip *IndexingProgress) CurrentProgress() int64 {
	return (ip.DonePostsCount + ip.DoneChannelsCount + ip.DoneUsersCount + ip.DoneFilesCount) * 100 / (ip.TotalPostsCount + ip.TotalChannelsCount + ip.TotalUsersCount + ip.TotalFilesCount)
}

func (ip *IndexingProgress) IsDone() bool {
	return ip.DonePosts && ip.DoneChannels && ip.DoneUsers && ip.DoneFiles
}
