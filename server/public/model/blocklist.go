package model

type TeamBlockUser struct {
	TeamId string `json:"channel_id"`
	BlockedId string `json:"blocked_id"`
	CreateAt  int64  `json:"-"`
	CreateBy  string `json:"create_by"`
}

type TeamBlockUserList []*TeamBlockUser

func (o *TeamBlockUser) PreSave() {
	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}
}

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

func (o *UserBlockUser) GetDMName() string {
	return GetDMNameFromIds(o.BlockedId, o.UserId)
}
