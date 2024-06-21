SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'Users' 
    AND column_name = 'Mobilephone'
  ) = 0, 
  'ALTER TABLE Users ADD COLUMN Mobilephone VARCHAR(24) NOT NULL DEFAULT \'\';', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
