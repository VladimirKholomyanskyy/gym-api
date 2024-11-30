import apiClient from './apiClient';
import { AddExerciseRequest } from '../types/api';

export const addExerciseToWorkout = async (
  data: AddExerciseRequest
): Promise<void> => {
  await apiClient.post('/workout-exercises', data);
};
