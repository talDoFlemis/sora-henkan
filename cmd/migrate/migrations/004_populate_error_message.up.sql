-- Populate null error_message values with empty string
UPDATE images SET error_message = '' WHERE error_message IS NULL;
