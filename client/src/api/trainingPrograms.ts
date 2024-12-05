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
  console.log(response);
  return response.data;
};

export const deleteTrainingProgram = async (
  programId: string
): Promise<boolean> => {
  const response = await apiClient.delete(`/training-programs/${programId}`);
  if (response.status === 200) {
    return true;
  }
  return false;
};
