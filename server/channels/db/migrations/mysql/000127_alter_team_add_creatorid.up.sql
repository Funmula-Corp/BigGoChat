SET @preparedStatement = (SELECT IF (
  (
    SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE table_schema = DATABASE() 
    AND table_name = 'Teams' 
    AND column_name = 'CreatorId'
  ) = 0, 
  'ALTER TABLE Teams ADD COLUMN CreatorId VARCHAR(26);', 
  'select 0;'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Teams'
        AND table_schema = DATABASE()
        AND column_name = 'Description'
        AND column_type != 'text'
    ) > 0,
    'ALTER TABLE Teams MODIFY COLUMN Description text;',
    'SELECT 1'
));
PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

CREATE PROCEDURE Migrate_CreatorId ()
BEGIN
	IF (
		SELECT COUNT(*)
		FROM Teams
		WHERE CreatorId IS NULL
	) > 0 THEN
		UPDATE
			Teams
			INNER JOIN (
				SELECT
                    Teams.Email Email,
					Users.Id userId
				FROM
					Teams
					LEFT JOIN Users ON Teams.Email = Users.Email
				WHERE
					Teams.CreatorId is NULL
				GROUP BY
					Teams.Id) AS q ON q.Email = Teams.Email
				SET
					CreatorId = userId
				WHERE
					CreatorId IS NULL;
		UPDATE Teams SET CreatorId='' WHERE CreatorId IS NULL;
	END IF;
END;
	CALL Migrate_CreatorId ();
	DROP PROCEDURE IF EXISTS Migrate_CreatorId;