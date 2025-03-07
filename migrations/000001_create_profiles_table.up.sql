CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) NOT NULL UNIQUE,
    sex VARCHAR(10) CHECK (sex IN ('male', 'female', 'other', 'prefer_not_to_say')),
    birthday DATE,
    weight DECIMAL(5,2) CHECK (weight > 0 AND weight < 500), -- Always stored in kg, realistic upper limit
    height DECIMAL(5,2) CHECK (height > 0 AND height < 3), -- Always stored in meters, realistic upper limit
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add index to improve soft delete queries
CREATE INDEX idx_profiles_deleted_at ON profiles(deleted_at) WHERE deleted_at IS NULL;
