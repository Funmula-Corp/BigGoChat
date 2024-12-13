DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users'
        AND column_name = 'description'
    ) THEN
        ALTER TABLE Users ADD COLUMN Description text DEFAULT NULL;
    END IF;
END $$; 