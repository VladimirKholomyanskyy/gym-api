CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    primary_muscle VARCHAR(255) NOT NULL,
    secondary_muscle TEXT[],
    equipment VARCHAR(255) NOT NULL,
    description TEXT
);