CREATE TABLE set_logs (
    id SERIAL PRIMARY KEY,
    workout_log_id INT NOT NULL REFERENCES workout_logs(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id),
    set_number INT NOT NULL CHECK (set_number > 0), -- e.g., Set 1, Set 2, Set 3
    planned_reps INT NOT NULL CHECK (planned_reps >= 0), -- Planned reps
    actual_reps INT NOT NULL CHECK (actual_reps >= 0), -- Actual reps performed
    planned_weight DECIMAL(5, 2) NOT NULL CHECK (planned_weight >= 0), -- Planned weight
    actual_weight DECIMAL(5, 2) NOT NULL CHECK (actual_weight >= 0) -- Actual weight used
);
