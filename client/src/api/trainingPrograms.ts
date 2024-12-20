import apiClient from "./apiClient";
import { CreateTrainingProgramRequest, TrainingProgram } from "../types/api";

export const createTrainingProgram = async (
  data: CreateTrainingProgramRequest
): Promise<TrainingProgram> => {
  const response = await apiClient.post<TrainingProgram>(
    "/training-programs",
    data
  );
  return response.data;
};

export const getAllTrainingPrograms = async (): Promise<TrainingProgram[]> => {
  const response = await apiClient.get<TrainingProgram[]>("/training-programs");
  return response.data;
};

export const getTrainingProgram = async (
  programId: number
): Promise<TrainingProgram> => {
  const response = await apiClient.get<TrainingProgram>(
    `/training-programs/${programId}`
  );
  return response.data;
};

export const deleteTrainingProgram = async (
  programId: number
): Promise<boolean> => {
  const response = await apiClient.delete(`/training-programs/${programId}`);
  if (response.status === 200) {
    return true;
  }
  return false;
};

export const updateTrainingProgram = async (
  programId: number,
  data: CreateTrainingProgramRequest
): Promise<TrainingProgram> => {
  console.log(`name=${data.name}`);
  console.log(`desc=${data.description}`);
  const response = await apiClient.patch(
    `/training-programs/${programId}`,
    data
  );

  return response.data;
};
