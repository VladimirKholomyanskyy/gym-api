CREATE TABLE workout_exercises (
    id SERIAL PRIMARY KEY,
    workout_id INT NOT NULL REFERENCES workouts(id),
    exercise_id INT NOT NULL REFERENCES exercises(id),
    sets INT NOT NULL CHECK (sets > 0),
    reps INT NOT NULL CHECK (reps > 0),
    weight DECIMAL(5, 2) NOT NULL CHECK (weight >= 0)
);
