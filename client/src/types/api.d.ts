export interface CreateTrainingProgramRequest {
    name: string;
    description: string;
  }
  
  export interface CreateWorkoutRequest {
    name: string;
  }
  
  export interface AddExerciseRequest {
    exerciseId: string;
    workoutId: string;
    sets: number;
    reps: number;
  }
  
  export interface Exercise {
    ID: string;
    Name: string;
    Description?: string;
  }

  export interface TrainingProgram {
    ID: string;
    Name: string;
    Description?: string;
  }

  export interface Workout {
    id: string;
    name: string;
    trainingProgramId: string;
  }
