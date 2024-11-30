import apiClient from './apiClient';
import { CreateTrainingProgramRequest, TrainingProgram } from '../types/api';

export const createTrainingProgram = async (
  data: CreateTrainingProgramRequest
): Promise<TrainingProgram> => {
    const response = await apiClient.post<TrainingProgram>('/training-programs', data);
    return response.data;
};
