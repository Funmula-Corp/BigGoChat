DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users'
        AND column_name = 'description'
    ) THEN
        ALTER TABLE Users DROP COLUMN Description;
    END IF;
END $$; 