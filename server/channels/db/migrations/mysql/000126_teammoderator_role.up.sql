SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'TeamMembers' 
    AND column_name = 'SchemeModerator'
  ) = 0, 
  'ALTER TABLE TeamMembers ADD COLUMN SchemeModerator boolean NOT NULL DEFAULT false;', 
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
  ) = 0, 
  'ALTER TABLE Schemes ADD COLUMN DefaultTeamModeratorRole VARCHAR(64) NOT NULL DEFAULT \'\';', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;
