SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'TeamMembers' 
    AND column_name = 'SchemeModerator'
  ) > 0, 
  'ALTER TABLE TeamMembers DROP COLUMN SchemeModerator;', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'Schemes' 
    AND column_name = 'DefaultTeamModeratorRole'
  ) > 0, 
  'ALTER TABLE Schemes DROP COLUMN DefaultTeamModeratorRole;', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
