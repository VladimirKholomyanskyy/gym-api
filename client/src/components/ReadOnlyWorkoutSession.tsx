import { For, Heading, Stack, Text } from "@chakra-ui/react";
import CompletedExerciseLog, {
  CompletedExerciseLogItem,
} from "./CompletedExerciseLog";
import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";

const ReadOnlyWorkoutSession = ({
  session,
  logs,
}: {
  session: WorkoutSessionResponse;
  logs: LogExerciseResponse[];
}) => {
  const computedLogItems =
    session?.workoutSnapshot?.workoutExercises?.map((exercise) => {
      const found: CompletedExerciseLogItem[] = logs
        .filter((e) => exercise.id === e.id)
        .map((log) => ({
          setNumber: log.setNumber,
          repsCompleted: log.repsCompleted,
          weightCompleted: log.weightUsed,
        }))
        .sort((a, b) => a.setNumber - b.setNumber);
      return {
        exerciseName: exercise.exercise?.name,
        items: found,
      };
    }) || [];

  return (
    <Stack>
      <Heading>{session.workoutSnapshot.name}</Heading>
      <Text>Started at: {session.startedAt}</Text>
      <Text>Completed at: {session.completedAt}</Text>
      <For each={computedLogItems}>
        {(item, index) => (
          <CompletedExerciseLog
            key={index}
            exerciseName={item.exerciseName ? item.exerciseName : ""}
            items={item.items}
          />
        )}
      </For>
    </Stack>
  );
};

export default ReadOnlyWorkoutSession;
