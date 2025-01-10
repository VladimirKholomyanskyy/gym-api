import {
  StartWorkoutSessionRequest,
  StartWorkoutSessionResponse,
  WSWorkoutSessionResponse,
} from "@/types/api";
import apiClient from "./apiClient";

export const createWorkoutSession = async (
  data: StartWorkoutSessionRequest
): Promise<StartWorkoutSessionResponse> => {
  console.log(`session=${data.workout_id}`);
  const response = await apiClient.post<StartWorkoutSessionResponse>(
    "/workout-sessions",
    data
  );
  return response.data;
};

export const finishWorkoutSession = async (
  sessionId: number
): Promise<StartWorkoutSessionResponse> => {
  const response = await apiClient.post<StartWorkoutSessionResponse>(
    `/workout-sessions/${sessionId}/finish`
  );
  return response.data;
};

export const getWorkoutSession = async (
  sessionId: number
): Promise<WSWorkoutSessionResponse> => {
  const response = await apiClient.get<WSWorkoutSessionResponse>(
    `/workout-sessions/${sessionId}`
  );
  return response.data;
};

export const getWorkoutSessions = async (): Promise<
  WSWorkoutSessionResponse[]
> => {
  const response = await apiClient.get<WSWorkoutSessionResponse[]>(
    `/workout-sessions`
  );
  return response.data;
};
