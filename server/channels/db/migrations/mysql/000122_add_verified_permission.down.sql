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

DELETE FROM `Roles` WHERE name in ('team_verified', 'channel_verified');
