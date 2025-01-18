import { useState, useEffect } from "react";
import WorkoutSessionCard, {
  WorkoutSessionCardProps,
} from "./WorkoutSessionCard";
import { For, Heading, Spinner, Stack } from "@chakra-ui/react";
import { WorkoutSessionsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";

const WorkoutSessionsPage = () => {
  const [workoutSessionsCardProps, setWorkoutSessionsCardProps] = useState<
    WorkoutSessionCardProps[]
  >([]);
  const [loading, setLoading] = useState(true);
  const workoutSessionApi = new WorkoutSessionsApi(apiConfig);

  useEffect(() => {
    const fetchWorkoutSessions = async () => {
      try {
        const sessions = await workoutSessionApi.listWorkoutSessions();
        const sessionProps = sessions.data.map((session) => {
          const shortDesc = session.workoutSnapshot.workoutExercises
            ?.map((e) => e.exercise?.name)
            .join(",");
          return {
            workoutSessionId: session.id,
            workoutName: session.workoutSnapshot.name,
            sessionStart: session.startedAt,
            sessionCompleted: session.completedAt ? session.completedAt : "",
            shortDescription: shortDesc ? shortDesc : "",
          };
        });
        setWorkoutSessionsCardProps(sessionProps);
      } catch (error) {
        console.error("Read Error fetching workout session:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkoutSessions();
  }, []);

  if (loading) return <Spinner>Loading...</Spinner>;
  return (
    <Stack>
      <Heading>Workout Sessions</Heading>
      <For each={workoutSessionsCardProps}>
        {(item, index) => (
          <WorkoutSessionCard
            key={index}
            workoutSessionId={item.workoutSessionId}
            workoutName={item.workoutName}
            sessionStart={item.sessionStart}
            shortDescription={item.shortDescription}
            sessionCompleted={item.sessionCompleted}
          />
        )}
      </For>
    </Stack>
  );
};

export default WorkoutSessionsPage;
