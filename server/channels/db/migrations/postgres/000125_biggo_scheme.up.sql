ALTER TABLE public.channelmembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultteammoderatorrole varchar(64) DEFAULT ''::character varying NULL;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultteamverifiedrole varchar(64) DEFAULT ''::character;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultchannelverifiedrole varchar(64) DEFAULT ''::character;
ALTER TABLE public.teammembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE public.teammembers ADD COLUMN IF NOT EXISTS schememoderator bool DEFAULT false NULL;

CREATE TABLE IF NOT EXISTS teamblockusers (
    teamid character varying(26),
    blockedid character varying(26),
    createby character varying(26),
    createat bigint DEFAULT 0,
    PRIMARY KEY (teamid, blockedid)
);
CREATE INDEX IF NOT EXISTS idx_teamblockusers_team_id ON teamblockusers USING hash (teamid);
CREATE INDEX IF NOT EXISTS idx_teamblockusers_blocked_id ON teamblockusers USING hash (blockedid);
