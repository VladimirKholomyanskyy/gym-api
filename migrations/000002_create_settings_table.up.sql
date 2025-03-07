CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    measurement_units ENUM('metric', 'imperial') NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    notifications_enabled BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
