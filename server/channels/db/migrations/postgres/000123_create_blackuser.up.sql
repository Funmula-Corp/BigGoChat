CREATE TABLE IF NOT EXISTS channelblockusers (
    channelid character varying(26),
    blockedid character varying(26),
    createat bigint DEFAULT 0,
    createby character varying(26),
    PRIMARY KEY (channelid, blockedid)
);
CREATE INDEX IF NOT EXISTS idx_channelblockusers_channel_id ON channelblockusers USING hash (channelid);
CREATE INDEX IF NOT EXISTS idx_channelblockusers_blocked_id ON channelblockusers USING hash (blockedid);

CREATE TABLE IF NOT EXISTS userblockusers (
    userid character varying(26),
    blockedid character varying(26),
    createat bigint DEFAULT 0,
    PRIMARY KEY (userid, blockedid)
);
CREATE INDEX IF NOT EXISTS idx_userblockusers_user_id ON userblockusers USING hash (userid);
CREATE INDEX IF NOT EXISTS idx_userblockusers_blocked_id ON userblockusers USING hash (blockedid);
