CREATE TABLE IF NOT EXISTS ChannelBlockUsers (
    ChannelId varchar(26) NOT NULL,
    BlockedId varchar(26) NOT NULL,
    CreateAt bigint(20) DEFAULT 0,
    CreateBy varchar(26) NOT NULL,
    PRIMARY KEY (ChannelId, BlockedId),
    KEY idx_channelblackusers_channelid (ChannelId),
    KEY idx_channelblackusers_userid (BlockedId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
