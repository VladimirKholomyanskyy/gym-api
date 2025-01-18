import { useEffect, useState } from "react";
import { useParams, Navigate, useNavigate } from "react-router";
import EditableWorkoutSession from "./EditableWorkoutSession";

import { Flex, Spinner } from "@chakra-ui/react";
import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";
import { WorkoutSessionsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";

const EditableWorkoutWrapper = () => {
  const { id } = useParams();
  const [workoutSession, setWorkoutSession] =
    useState<WorkoutSessionResponse | null>(null);
  const [exerciseLogs, setExerciseLogs] = useState<LogExerciseResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);
  const workoutSessionApi = new WorkoutSessionsApi(apiConfig);
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
        const exLogs = await workoutSessionApi.listExerciseLogs(id);
        console.log("Fetched logs:", exLogs);
        setExerciseLogs(exLogs.data);
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

  if (workoutSession?.completedAt) {
    console.log("Session completed. Redirecting to view.");
    return <Navigate to={`/workout-sessions/${id}/view`} replace />;
  }

  return (
    <EditableWorkoutSession session={workoutSession} logs={exerciseLogs} />
  );
};

export default EditableWorkoutWrapper;
