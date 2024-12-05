export interface CreateTrainingProgramRequest {
  name: string;
  description: string;
}

export interface CreateWorkoutRequest {
  name: string;
}

export interface AddExerciseRequest {
  exercise_id: string;
  workout_id: string;
  sets: number;
  reps: number;
}

export interface Exercise {
  ID: string;
  Name: string;
  Description?: string;
}

export interface TrainingProgram {
  id: string;
  name: string;
  description?: string;
}

export interface Workout {
  ID: string;
  Name: string;
  TrainingProgramID: string;
}
