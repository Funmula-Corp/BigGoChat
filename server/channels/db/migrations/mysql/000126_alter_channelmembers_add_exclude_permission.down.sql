SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'ChannelMembers' 
    AND column_name = 'ExcludePermissions'
  ) > 0, 
  'ALTER TABLE ChannelMembers DROP COLUMN ExcludePermissions;',
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
