-- Drop indexes
DROP INDEX IF EXISTS idx_images_checksum;
DROP INDEX IF EXISTS idx_images_status;
DROP INDEX IF EXISTS idx_images_created_at;

-- Drop table
DROP TABLE IF EXISTS images;
