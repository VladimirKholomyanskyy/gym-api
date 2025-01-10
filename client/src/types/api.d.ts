export interface CreateTrainingProgramRequest {
  name: string;
  description: string;
}

export interface CreateWorkoutRequest {
  name: string;
}

export interface AddExerciseRequest {
  exercise_id: number;
  workout_id: number;
  sets: number;
  reps: number;
}

export interface Exercise {
  ID: number;
  Name: string;
  Description?: string;
}

export interface TrainingProgram {
  id: number;
  name: string;
  description: string;
}

export interface Workout {
  id: number;
  name: string;
}

export interface WorkoutExercise {
  id: number;
  workout_id: number;
  exercise_id: number;
  sets: number;
  reps: number;
}

// Request to start a workout session
export interface StartWorkoutSessionRequest {
  workout_id: number; // Corresponds to Go's `uint`
}

// Response after starting a workout session
export interface StartWorkoutSessionResponse {
  session_id: number; // Corresponds to Go's `uint`
  started_at: string; // ISO 8601 formatted date-time string
}

// Response for getting workout session logs
export interface GetWorkoutSessionLogsResponse {
  sessionId: number; // Corresponds to Go's `uint`
  workoutSnapshot: Record<string, any>; // Corresponds to Go's `map[string]interface{}`
  logs: WorkoutLogResponse[]; // Array of logs
  startedAt: string; // ISO 8601 formatted date-time string
  completedAt?: string | null; // Optional ISO 8601 date-time string (nullable)
}

// Response for an individual workout log entry
export interface WorkoutLogResponse {
  exerciseName: number; // Corresponds to Go's `uint`
  setNumber: number; // Corresponds to Go's `int`
  repsCompleted: number; // Corresponds to Go's `int`
  weightUsed: number; // Corresponds to Go's `float64`
  loggedAt: string; // ISO 8601 formatted date-time string
}

// Request to log an exercise
export interface LogExerciseRequest {
  exercise_id: number; // Corresponds to Go's `uint`, required
  set_number: number; // Corresponds to Go's `int`, required
  reps_completed: number; // Corresponds to Go's `int`, required
  weight_used: number; // Corresponds to Go's `float64`, required
}

// Response after logging an exercise
export interface LogExerciseResponse {
  log_id: number; // Corresponds to Go's `uint`
  exercise_id: number; // Corresponds to Go's `uint`
  set_number: number; // Corresponds to Go's `int`
  reps_completed: number; // Corresponds to Go's `int`
  weight_used: number; // Corresponds to Go's `float64`
  logged_at: string; // ISO 8601 formatted date-time string
}

// Exercise model
export interface WSExercise {
  ID: number;
  Name: string;
  PrimaryMuscle: string;
  SecondaryMuscle: string[];
  Equipment: string;
  Description: string;
}

// WorkoutExercise model
export interface WSWorkoutExercise {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  WorkoutID: number;
  ExerciseID: number;
  Sets: number;
  Reps: number;
  Weight: number;
  Exercise: WSExercise;
}

// WorkoutSnapshot model
export interface WSWorkoutSnapshot {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  Name: string;
  TrainingProgramID: number;
  Exercises: WSWorkoutExercise[];
}

// Main response model
export interface WSWorkoutSessionResponse {
  session_id: number;
  started_at: string;
  completed_at?: string;
  workout_snapshot: WSWorkoutSnapshot;
}
