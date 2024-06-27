CREATE TABLE IF NOT EXISTS ChannelBlockUsers (
    ChannelId varchar(26) NOT NULL,
    BlockedId varchar(26) NOT NULL,
    CreateAt bigint(20) DEFAULT 0,
    CreateBy varchar(26) NOT NULL,
    PRIMARY KEY (ChannelId, BlockedId),
    KEY idx_channelblockusers_channel_id (ChannelId),
    KEY idx_channelblockusers_blocked_id (BlockedId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS UserBlockUsers (
    UserId varchar(26) NOT NULL,
    BlockedId varchar(26) NOT NULL,
    CreateAt bigint(20) DEFAULT 0,
    PRIMARY KEY (UserId, BlockedId),
    KEY idx_userblockusers_user_id (UserId),
    KEY idx_userblockusers_blocked_id (BlockedId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
