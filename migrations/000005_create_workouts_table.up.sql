CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    training_program_id UUID NOT NULL REFERENCES training_programs(id),
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL CHECK (position > 0),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT unique_workout_position UNIQUE (training_program_id, position)
);
