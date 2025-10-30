-- Remove error_message column from images table
ALTER TABLE images DROP COLUMN IF EXISTS error_message;
