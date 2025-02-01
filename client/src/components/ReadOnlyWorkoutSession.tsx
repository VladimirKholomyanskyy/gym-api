import { Box, For, Heading, Stack, VStack } from "@chakra-ui/react";
import CompletedExerciseLog, {
  CompletedExerciseLogItem,
} from "./CompletedExerciseLog";
import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";
import { formatDateTime } from "@/utils/dateUtils";

const ReadOnlyWorkoutSession = ({
  session,
  logs,
}: {
  session: WorkoutSessionResponse;
  logs: LogExerciseResponse[];
}) => {
  const computedLogItems =
    session?.workoutSnapshot?.workoutExercises?.map((workoutExercise) => {
      const found: CompletedExerciseLogItem[] = logs
        .filter((e) => workoutExercise.exercise?.id === e.exerciseId)
        .map((log) => ({
          setNumber: log.setNumber,
          repsCompleted: log.repsCompleted,
          weightCompleted: log.weightUsed,
        }))
        .sort((a, b) => a.setNumber - b.setNumber);
      return {
        exerciseName: workoutExercise.exercise?.name,
        items: found,
      };
    }) || [];

  return (
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        {session?.workoutSnapshot.name}
      </Heading>
      <Heading size="lg">
        Started at: {formatDateTime(session?.startedAt)}
      </Heading>
      <Heading size="lg">
        Completed at: {formatDateTime(session?.startedAt)}
      </Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
        <For each={computedLogItems}>
          {(item, index) => (
            <CompletedExerciseLog
              key={index}
              exerciseName={item.exerciseName ? item.exerciseName : ""}
              items={item.items}
            />
          )}
        </For>
      </VStack>
    </Box>
  );
};

export default ReadOnlyWorkoutSession;
