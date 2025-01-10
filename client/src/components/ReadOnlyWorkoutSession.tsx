import { LogExerciseResponse, WSWorkoutSessionResponse } from "@/types/api";
import { For, Heading, Stack, Text } from "@chakra-ui/react";
import CompletedExerciseLog, {
  CompletedExerciseLogItem,
} from "./CompletedExerciseLog";

const ReadOnlyWorkoutSession = ({
  session,
  logs,
}: {
  session: WSWorkoutSessionResponse;
  logs: LogExerciseResponse[];
}) => {
  const computedLogItems =
    session?.workout_snapshot?.Exercises.map((exercise) => {
      const found: CompletedExerciseLogItem[] = logs
        .filter((e) => exercise.ExerciseID === e.exercise_id)
        .map((log) => ({
          setNumber: log.set_number,
          repsCompleted: log.reps_completed,
          weightCompleted: log.weight_used,
        }))
        .sort((a, b) => a.setNumber - b.setNumber);
      return {
        exerciseName: exercise.Exercise.Name,
        items: found,
      };
    }) || [];

  return (
    <Stack>
      <Heading>{session.workout_snapshot.Name}</Heading>
      <Text>Started at: {session.started_at}</Text>
      <Text>Completed at: {session.completed_at}</Text>
      <For each={computedLogItems}>
        {(item, index) => (
          <CompletedExerciseLog
            key={index}
            exerciseName={item.exerciseName}
            items={item.items}
          />
        )}
      </For>
    </Stack>
  );
};

export default ReadOnlyWorkoutSession;
