CREATE TABLE workout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    workout_id UUID NOT NULL REFERENCES workouts(id),
    snapshot JSONB NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NULL
);