-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    original_image_url TEXT NOT NULL,
    object_storage_image_key TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    transformed_image_key TEXT,
    checksum VARCHAR(64),
    transformations JSONB DEFAULT '[]'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_images_status ON images(status);
CREATE INDEX IF NOT EXISTS idx_images_checksum ON images(checksum);
CREATE INDEX IF NOT EXISTS idx_images_transformations ON images USING gin(transformations);
