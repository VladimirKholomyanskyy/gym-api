CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    measurement_units ENUM('metric', 'imperial') NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    notifications_enabled BOOLEAN NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
