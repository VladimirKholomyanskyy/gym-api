import apiClient from './apiClient';
import { CreateWorkoutRequest, Workout } from '../types/api';

export const createWorkout = async (
  programId: string,
  data: CreateWorkoutRequest
): Promise<Workout> => {
  const response = await apiClient.post<Workout>(`/training-programs/${programId}/workouts`, data);
  return response.data;
};