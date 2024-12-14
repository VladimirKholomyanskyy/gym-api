import apiClient from "./apiClient";
import { AddExerciseRequest, WorkoutExercise } from "../types/api";

export const addExerciseToWorkout = async (
  data: AddExerciseRequest
): Promise<WorkoutExercise> => {
  const response = await apiClient.post("/workout-exercises", data);
  return response.data;
};

export const getAllWorkoutExercises = async (
  workout_id: number
): Promise<WorkoutExercise[]> => {
  const response = await apiClient.get("/workout-exercises", {
    params: {
      parent: `workout/${workout_id}`,
    },
  });
  return response.data || [];
};
