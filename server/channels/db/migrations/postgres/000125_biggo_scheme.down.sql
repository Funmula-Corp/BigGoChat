ALTER TABLE public.channelmembers DROP COLUMN IF EXISTS schemeverified;
ALTER TABLE public.schemes DROP COLUMN IF EXISTS defaultteamverifiedrole;
ALTER TABLE public.schemes DROP COLUMN IF EXISTS defaultchannelverifiedrole;
ALTER TABLE public.teammembers DROP COLUMN IF EXISTS schemeverified;

ALTER TABLE public.schemes DROP COLUMN IF EXISTS defaultteammoderatorrole;
ALTER TABLE public.teammembers DROP COLUMN IF EXISTS schememoderator;

DROP TABLE IF EXISTS teamblockusers;
