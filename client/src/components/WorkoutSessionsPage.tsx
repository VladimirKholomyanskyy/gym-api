import { useState, useEffect } from "react";
import WorkoutSessionCard, {
  WorkoutSessionCardProps,
} from "./WorkoutSessionCard";
import { Box, For, Heading, Spinner, VStack } from "@chakra-ui/react";
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
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        Workout Sessions
      </Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
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
      </VStack>
    </Box>
  );
};

export default WorkoutSessionsPage;
