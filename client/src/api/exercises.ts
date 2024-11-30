// api/exercises.ts
import { Exercise } from '../types/api'; // Define the type of an exercise object
import apiClient from './apiClient';


export const listExercises = async (): Promise<Exercise[]> => {
      const response = await apiClient.get<Exercise[]>('/exercises');
      console.log(response.status)
      return response.data;
  };
