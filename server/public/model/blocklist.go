package model

type ChannelBlockUser struct {
	ChannelId string `json:"channel_id"`
	BlockedId string `json:"blocked_id"`
	CreateAt  int64  `json:"-"`
	CreateBy  string `json:"create_by"`
}

type ChannelBlockUserList []*ChannelBlockUser

func (o *ChannelBlockUser) PreSave() {
	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}
}

type UserBlockUser struct {
	UserId    string `json:"user_id"`
	BlockedId string `json:"blocked_id"`
	CreateAt  int64  `json:"-"`
}

type UserBlockUserList []*UserBlockUser

func (o *UserBlockUser) PreSave() {
	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}
}
