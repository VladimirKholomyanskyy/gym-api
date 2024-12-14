import apiClient from "./apiClient";
import { CreateWorkoutRequest, Workout } from "../types/api";

export const createWorkout = async (
  programId: number,
  data: CreateWorkoutRequest
): Promise<Workout> => {
  const response = await apiClient.post<Workout>(
    `/training-programs/${programId}/workouts`,
    data
  );
  return response.data;
};

export const getAllWorkouts = async (programId: number): Promise<Workout[]> => {
  const response = await apiClient.get<Workout[]>(
    `/training-programs/${programId}/workouts`
  );
  return response.data;
};

export const getWorkout = async (
  programId: number,
  workoutId: number
): Promise<Workout> => {
  const response = await apiClient.get<Workout>(
    `/training-programs/${programId}/workouts/${workoutId}`
  );
  return response.data;
};

export const deleteWorkout = async (
  programId: number,
  workoutId: number
): Promise<Workout> => {
  const response = await apiClient.delete<Workout>(
    `/training-programs/${programId}/workouts/${workoutId}`
  );
  console.log(response);
  return response.data;
};

export const updateWorkout = async (
  programId: number,
  workoutId: number,
  data: CreateWorkoutRequest
): Promise<Workout> => {
  const response = await apiClient.put<Workout>(
    `/training-programs/${programId}/workouts/${workoutId}`,
    data
  );
  console.log(`response=${response}`);
  return response.data;
};
