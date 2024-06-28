SET @preparedStatement = (SELECT IF(
    EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND column_name = 'SchemeVerified'
    ) > 0,
    'ALTER TABLE ChannelMembers DROP COLUMN SchemeVerified;',
    'SELECT 1;'
));

PREPARE addColumnIfExists FROM @preparedStatement;
EXECUTE addColumnIfExists;
DEALLOCATE PREPARE addColumnIfExists;

SET @preparedStatement = (SELECT IF(
    EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultTeamVerifiedRole'
    ) > 0,
    'ALTER TABLE Schemes DROP COLUMN DefaultTeamVerifiedRole;',
    'SELECT 1;'
));

PREPARE addColumnIfExists FROM @preparedStatement;
EXECUTE addColumnIfExists;
DEALLOCATE PREPARE addColumnIfExists;

SET @preparedStatement = (SELECT IF(
    EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultChannelVerifiedRole'
    ) > 0,
    'ALTER TABLE Schemes DROP COLUMN DefaultChannelVerifiedRole;',
    'SELECT 1;'
));

PREPARE addColumnIfExists FROM @preparedStatement;
EXECUTE addColumnIfExists;
DEALLOCATE PREPARE addColumnIfExists;

SET @preparedStatement = (SELECT IF(
    EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'TeamMembers'
        AND table_schema = DATABASE()
        AND column_name = 'SchemeVerified'
    ) > 0,
    'ALTER TABLE TeamMembers DROP COLUMN SchemeVerified;',
    'SELECT 1;'
));

PREPARE addColumnIfExists FROM @preparedStatement;
EXECUTE addColumnIfExists;
DEALLOCATE PREPARE addColumnIfExists;

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
