SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'Users' 
    AND column_name = 'mobilephone'
  ) = 0, 
  'ALTER TABLE Users ADD COLUMN mobilephone VARCHAR(24);', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
