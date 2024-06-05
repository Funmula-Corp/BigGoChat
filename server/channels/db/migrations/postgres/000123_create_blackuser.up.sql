CREATE TABLE IF NOT EXISTS channelblockusers (
    channelid character varying(26),
    blockedid character varying(26),
    createat bigint DEFAULT 0,
    createby character varying(26),
    PRIMARY KEY (channelid, blockedid)
);

