ALTER TABLE channelmembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultteamverifiedrole character varying(64) DEFAULT '';
ALTER TABLE schemes ADD COLUMN IF NOT EXISTS defaultchannelverifiedrole character varying(64) DEFAULT '';
ALTER TABLE teammembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
