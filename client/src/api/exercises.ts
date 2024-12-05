// api/exercises.ts
import { User } from 'oidc-client-ts';
import { Exercise } from '../types/api'; // Define the type of an exercise object
import apiClient from './apiClient';


export const listExercises = async (user: User): Promise<Exercise[]> => {
        if (!user?.access_token) {
            throw new Error('User is not authenticated or token is missing');
        }
        console.log(`token=${user.access_token}`)
      const response = await apiClient.get<Exercise[]>('/exercises', {
        headers: {
          Authorization: `Bearer ${user.access_token}`,
        },
    });
      return response.data;
  };
