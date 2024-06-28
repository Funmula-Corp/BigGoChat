SET @preparedStatement = (SELECT IF(
    EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND column_name = 'SchemeVerified'
    ) > 0,
    'SELECT 1;',
    'ALTER TABLE ChannelMembers ADD COLUMN SchemeVerified tinyint(4);'
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
    'SELECT 1;',
    'ALTER TABLE Schemes ADD COLUMN DefaultTeamVerifiedRole varchar(64) DEFAULT \'\';'
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
    'SELECT 1;',
    'ALTER TABLE Schemes ADD COLUMN DefaultChannelVerifiedRole varchar(64) DEFAULT \'\';'
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
    'SELECT 1;',
    'ALTER TABLE TeamMembers ADD COLUMN SchemeVerified tinyint(4);'
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
