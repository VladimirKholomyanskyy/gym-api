CREATE TABLE exercise_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id), 
    set_number INT NOT NULL CHECK (set_number > 0),
    reps INT NOT NULL CHECK (reps >= 0),
    weight DECIMAL(5, 2) NOT NULL CHECK (weight >= 0),
    logged_at TIMESTAMP NOT NULL
);

