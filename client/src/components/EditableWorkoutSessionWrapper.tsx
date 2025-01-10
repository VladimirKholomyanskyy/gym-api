import { useEffect, useState } from "react";
import { useParams, Navigate } from "react-router";
import EditableWorkoutSession from "./EditableWorkoutSession";
import { getWorkoutSession } from "@/api/workout-sessions";
import { LogExerciseResponse, WSWorkoutSessionResponse } from "@/types/api";
import { getExerciseLogs } from "@/api/exercise-logs";
import { Flex, Spinner } from "@chakra-ui/react";

const EditableWorkoutWrapper = () => {
  const { id } = useParams();
  const [workoutSession, setWorkoutSession] =
    useState<WSWorkoutSessionResponse | null>(null);
  const [exerciseLogs, setExerciseLogs] = useState<LogExerciseResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);

  useEffect(() => {
    const fetchWorkoutSession = async () => {
      try {
        const session = await getWorkoutSession(Number(id));
        console.log("Fetched session:", session);
        setWorkoutSession(session);
        const exLogs = await getExerciseLogs(Number(id));
        console.log("Fetched logs:", exLogs);
        setExerciseLogs(exLogs);
      } catch (error) {
        console.error("Error fetching workout session:", error);
        setFetchError(true);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkoutSession();
  }, [id]);

  if (loading)
    return (
      <Flex justifyContent="center" alignItems="center" height="100%">
        <Spinner />
      </Flex>
    );

  if (fetchError) {
    console.log("Redirecting to error due to fetch error.");
    return <Navigate to="/error" replace />;
  }

  if (!workoutSession) {
    console.log("Redirecting to error as workoutSession is null.");
    return <Navigate to="/error" replace />;
  }

  if (workoutSession?.completed_at) {
    console.log("Session completed. Redirecting to view.");
    return <Navigate to={`/workout-sessions/${id}/view`} replace />;
  }

  return (
    <EditableWorkoutSession session={workoutSession} logs={exerciseLogs} />
  );
};

export default EditableWorkoutWrapper;
