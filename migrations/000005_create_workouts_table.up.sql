CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    training_program_id UUID NOT NULL REFERENCES training_programs(id),
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL CHECK (position > 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_workout_position UNIQUE (training_program_id, position)
);
