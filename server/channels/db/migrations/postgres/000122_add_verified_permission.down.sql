ALTER TABLE channelmembers DROP COLUMN IF EXISTS schemeverified;
ALTER TABLE schemes DROP COLUMN IF EXISTS defaultteamverifiedrole;
ALTER TABLE schemes DROP COLUMN IF EXISTS defaultchannelverifiedrole;
ALTER TABLE teammembers DROP COLUMN IF EXISTS schemeverified;

DELETE FROM public.roles WHERE name IN ('system_verified', 'team_verified', 'channel_verified');
