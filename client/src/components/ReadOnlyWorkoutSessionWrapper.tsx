import { useState, useEffect } from "react";
import { useParams, Navigate, useNavigate } from "react-router";
import ReadOnlyWorkoutSession from "./ReadOnlyWorkoutSession";
import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";
import { ExerciseLogsApi, WorkoutSessionsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";

const ReadOnlyWorkoutSessionWrapper = () => {
  const { id } = useParams();
  const [workoutSession, setWorkoutSession] =
    useState<WorkoutSessionResponse | null>(null);
  const [exerciseLogs, setExerciseLogs] = useState<LogExerciseResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const workoutSessionApi = new WorkoutSessionsApi(apiConfig);
  const exerciseLogsApi = new ExerciseLogsApi(apiConfig);
  const navigate = useNavigate();
  if (!id) {
    navigate("/error");
    return null;
  }
  useEffect(() => {
    const fetchWorkoutSession = async () => {
      try {
        const session = await workoutSessionApi.getWorkoutSession(id);
        console.log("Fetched session:", session);
        setWorkoutSession(session.data);

        const exLogs = await exerciseLogsApi.listExerciseLogs(id);
        console.log("Fetched logs:", exLogs);
        setExerciseLogs(exLogs.data);
      } catch (error) {
        console.error("Error fetching workout session:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkoutSession();
  }, [id]);

  if (loading) return <div>Loading...</div>;

  console.log("Read only");
  if (workoutSession?.completedAt) {
    return (
      <ReadOnlyWorkoutSession session={workoutSession} logs={exerciseLogs} />
    );
  }
  console.log("Read only2");
  return <Navigate to={`/workout-sessions/${id}/edit`} replace />;
};

export default ReadOnlyWorkoutSessionWrapper;
