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

export const deleteWorkoutExercise = async (
  workoutExerciseId: number
): Promise<boolean> => {
  const response = await apiClient.delete(
    `/workout-exercises/${workoutExerciseId}`
  );
  if (response.status === 200) {
    return true;
  }
  return false;
};

export const patchWorkoutExercise = async (
  workoutExerciseId: number,
  data: AddExerciseRequest
): Promise<WorkoutExercise> => {
  const response = await apiClient.patch(
    `/workout-exercises/${workoutExerciseId}`,
    data
  );
  return response.data;
};
