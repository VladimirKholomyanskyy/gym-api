CREATE TABLE workout_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    workout_id INT NOT NULL REFERENCES workouts(id),
    snapshot JSONB NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP NULL
);