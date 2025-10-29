-- Add error_message column to images table
ALTER TABLE images ADD COLUMN IF NOT EXISTS error_message TEXT;
