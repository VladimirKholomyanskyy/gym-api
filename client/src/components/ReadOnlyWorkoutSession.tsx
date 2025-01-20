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
    <Stack>
      <Box mb={7} mt={7}>
        <Heading size="2xl" fontWeight="bold" textAlign="center">
          {session?.workoutSnapshot.name}
        </Heading>
        <Heading size="lg">
          Started at: {formatDateTime(session?.startedAt)}
        </Heading>
        <Heading size="lg">
          Completed at: {formatDateTime(session?.startedAt)}
        </Heading>
      </Box>
      <VStack align="stretch" width="100%" paddingLeft="8" paddingRight="8">
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
    </Stack>
  );
};

export default ReadOnlyWorkoutSession;
