import apiClient from './apiClient';
import { CreateWorkoutRequest, Workout } from '../types/api';
import { User } from 'oidc-client-ts';

export const createWorkout = async (user: User,
  programId: string,
  data: CreateWorkoutRequest
): Promise<Workout> => {
  const response = await apiClient.post<Workout>(`/training-programs/${programId}/workouts`, 
    data,
    {
      headers: {
        Authorization: `Bearer ${user.access_token}`,
      }
    },
);
  return response.data;
};