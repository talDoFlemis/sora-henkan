ALTER TABLE IF EXISTS images
    DROP COLUMN IF EXISTS transformations;

ALTER TABLE IF EXISTS images
    ADD COLUMN transformations JSONB DEFAULT '[]'::jsonb;

CREATE INDEX IF NOT EXISTS idx_images_transformations ON images USING gin(transformations);