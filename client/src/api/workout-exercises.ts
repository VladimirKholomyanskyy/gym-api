import apiClient from './apiClient';
import { AddExerciseRequest } from '../types/api';
import { User } from 'oidc-client-ts';

export const addExerciseToWorkout = async (user: User,
  data: AddExerciseRequest
): Promise<void> => {
  await apiClient.post('/workout-exercises', 
    data,
    {
      headers: {
        Authorization: `Bearer ${user.access_token}`,
      }
    },
);
};
