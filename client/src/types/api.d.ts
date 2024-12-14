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
