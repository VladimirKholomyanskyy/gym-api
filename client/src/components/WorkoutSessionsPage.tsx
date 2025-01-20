import { useState, useEffect } from "react";
import WorkoutSessionCard, {
  WorkoutSessionCardProps,
} from "./WorkoutSessionCard";
import { Box, For, Heading, Spinner, Stack, VStack } from "@chakra-ui/react";
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
      <Box mb={7} mt={7}>
        <Heading size="2xl" fontWeight="bold" textAlign="center">
          Workout Sessions
        </Heading>
      </Box>
      <VStack
        gap={6}
        align="stretch"
        width="100%"
        paddingLeft="8"
        paddingRight="8"
      >
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
    </Stack>
  );
};

export default WorkoutSessionsPage;
