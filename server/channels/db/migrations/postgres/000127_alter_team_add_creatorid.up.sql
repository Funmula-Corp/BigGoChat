ALTER TABLE public.teams ALTER COLUMN description TYPE VARCHAR(1024);
ALTER TABLE public.teams ADD COLUMN IF NOT EXISTS creatorid VARCHAR(26);

DO $$
BEGIN
    IF (SELECT count(*) FROM teams WHERE creatorid IS NULL) > 0 THEN
        WITH tc AS (
            SELECT u.id, t.email FROM teams AS t LEFT JOIN users AS u ON u.email = t.email WHERE u.id IS NOT NULL
        ) UPDATE team AS t SET creatorid = tc.id FROM tc WHERE t.email = tc.email;

        UPDATE team SET creatorid = '' WHERE creatorid IS NULL;
    END IF;
END $$;