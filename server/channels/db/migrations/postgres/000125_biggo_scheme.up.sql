ALTER TABLE public.channelmembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultteammoderatorrole varchar(64) DEFAULT ''::character varying NULL;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultteamverifiedrole varchar(64) DEFAULT ''::character;
ALTER TABLE public.schemes ADD COLUMN IF NOT EXISTS defaultchannelverifiedrole varchar(64) DEFAULT ''::character;
ALTER TABLE public.teammembers ADD COLUMN IF NOT EXISTS schemeverified boolean;
ALTER TABLE public.teammembers ADD COLUMN IF NOT EXISTS schememoderator bool DEFAULT false NULL;
