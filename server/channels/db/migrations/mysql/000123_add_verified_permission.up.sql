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
    'ALTER TABLE Schemes ADD COLUMN DefaultTeamVerifiedRole varchar(64);'
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
    'ALTER TABLE Schemes ADD COLUMN DefaultChannelVerifiedRole varchar(64);'
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


INSERT INTO `Roles` VALUES ('biggoyyyyyyyyyyyyyyyyyyyyd','team_verified','authentication.roles.team_verified.name','authentication.roles.team_verified.description',1552023386683,1552023386683,0,' list_team_channels join_public_channels read_public_channel view_team create_public_channel manage_public_channel_properties delete_public_channel create_private_channel manage_private_channel_properties delete_private_channel invite_user add_user_to_team',1,1),('biggoyyyyyyyyyyyyyyyyyyyyr','channel_verified','authentication.roles.channel_verified.name','authentication.roles.channel_verified.description',1552023386587,1552023386587,0,' read_channel add_reaction remove_reaction manage_public_channel_members upload_file get_public_link create_post use_slash_commands manage_private_channel_members delete_post edit_post',1,1);
