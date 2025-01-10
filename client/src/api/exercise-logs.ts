import { LogExerciseRequest, LogExerciseResponse } from "@/types/api";
import apiClient from "./apiClient";

export const createExerciseLog = async (
  sessionId: number,
  data: LogExerciseRequest
): Promise<LogExerciseResponse> => {
  const response = await apiClient.post<LogExerciseResponse>(
    `/workout-sessions/${sessionId}/logs`,
    data
  );
  return response.data;
};

export const getExerciseLog = async (
  sessionId: number,
  logId: number
): Promise<LogExerciseResponse> => {
  const response = await apiClient.get<LogExerciseResponse>(
    `/workout-sessions/${sessionId}/logs/${logId}`
  );
  return response.data;
};

export const getExerciseLogs = async (
  sessionId: number
): Promise<LogExerciseResponse[]> => {
  const response = await apiClient.get<LogExerciseResponse[]>(
    `/workout-sessions/${sessionId}/logs`
  );
  return response.data;
};
