CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) NOT NULL UNIQUE,
    birthday DATE,
    weight DECIMAL(5,2) CHECK (weight > 0), -- Always stored in kg
    height DECIMAL(5,2) CHECK (height > 0),
    avatar_url TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
