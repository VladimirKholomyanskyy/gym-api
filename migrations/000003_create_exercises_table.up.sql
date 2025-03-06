CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    primary_muscle VARCHAR(255) NOT NULL,
    secondary_muscle TEXT[],
    equipment VARCHAR(255) NOT NULL,
    description TEXT
);