import { getExerciseLogs } from "@/api/exercise-logs";
import { getWorkoutSession } from "@/api/workout-sessions";
import { WSWorkoutSessionResponse, LogExerciseResponse } from "@/types/api";
import { useState, useEffect } from "react";
import { useParams, Navigate } from "react-router";
import ReadOnlyWorkoutSession from "./ReadOnlyWorkoutSession";

const ReadOnlyWorkoutSessionWrapper = () => {
  const { id } = useParams();
  const [workoutSession, setWorkoutSession] =
    useState<WSWorkoutSessionResponse | null>(null);
  const [exerciseLogs, setExerciseLogs] = useState<LogExerciseResponse[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchWorkoutSession = async () => {
      try {
        const session = await getWorkoutSession(Number(id));
        console.log("read session=", session);
        setWorkoutSession(session);
        const exLogs = await getExerciseLogs(Number(id));
        setExerciseLogs(exLogs);
      } catch (error) {
        console.error("Read Error fetching workout session:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkoutSession();
  }, [id]);

  if (loading) return <div>Loading...</div>;

  console.log("Read only");
  if (workoutSession?.completed_at) {
    return (
      <ReadOnlyWorkoutSession session={workoutSession} logs={exerciseLogs} />
    );
  }
  console.log("Read only2");
  return <Navigate to={`/workout-sessions/${id}/edit`} replace />;
};

export default ReadOnlyWorkoutSessionWrapper;
