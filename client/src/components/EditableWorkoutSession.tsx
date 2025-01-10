import { LogExerciseResponse, WSWorkoutSessionResponse } from "@/types/api";
import ExerciseLog, { ExerciseLogItem } from "./ExerciseLog";
import { Button } from "./ui/button";
import { createExerciseLog } from "@/api/exercise-logs";
import { finishWorkoutSession } from "@/api/workout-sessions";
import { Stack, Heading, Card, Flex } from "@chakra-ui/react";
import { useState } from "react";

const EditableWorkoutSession = ({
  session,
  logs,
}: {
  session: WSWorkoutSessionResponse;
  logs: LogExerciseResponse[];
}) => {
  const [currentCardIndex, setCurrentCardIndex] = useState(0);

  const handleNext = () => {
    setCurrentCardIndex((prevIndex) => {
      return prevIndex < computedLogItems.length - 1 ? prevIndex + 1 : 0;
    });
  };

  // Function to handle "Previous" button
  const handlePrevious = () => {
    setCurrentCardIndex((prevIndex) =>
      prevIndex > 0 ? prevIndex - 1 : computedLogItems.length - 1
    );
  };

  const handleComplete = async () => {
    try {
      await finishWorkoutSession(session.session_id);
    } catch (error) {
      console.log(error);
    }
  };
  const handleLog = async (
    exerciseId: number,
    setNumber: number,
    repsCompleted: number,
    weightUsed: number
  ) => {
    try {
      await createExerciseLog(session.session_id, {
        exercise_id: exerciseId,
        reps_completed: repsCompleted,
        set_number: setNumber,
        weight_used: weightUsed,
      });
    } catch (error) {
      console.log(error);
    }
  };

  const computedLogItems =
    session?.workout_snapshot?.Exercises.map((exercise) => {
      const exLogs: ExerciseLogItem[] = Array.from(
        { length: exercise.Sets },
        (_, index) => {
          const found = logs.find(
            (e) =>
              exercise.ExerciseID === e.exercise_id &&
              e.set_number === index + 1
          );
          return {
            id: index + 1,
            prevReps: 3,
            prevWeight: 20,
            reqReps: exercise.Reps,
            currentReps: found?.reps_completed,
            currentWeight: found?.weight_used,
          };
        }
      );

      return {
        exerciseId: exercise.ExerciseID,
        exerciseName: exercise.Exercise.Name,
        logs: exLogs,
      };
    }) || [];
  console.log("computedLogItems:", computedLogItems);
  return (
    <Stack>
      <Heading size="4xl">{session?.workout_snapshot.Name}</Heading>
      <Heading size="2xl">{session?.started_at}</Heading>
      <Button onClick={handleComplete}>Finish</Button>
      {computedLogItems.length > 0 &&
        computedLogItems[currentCardIndex]?.logs && (
          <Card.Root>
            <Card.Header>
              <Card.Title>
                {computedLogItems[currentCardIndex].exerciseName}
              </Card.Title>
              <Flex justify="space-between" w="300px">
                <Button onClick={handlePrevious} colorScheme="blue">
                  Previous
                </Button>
                <Button onClick={handleNext} colorScheme="blue">
                  Next
                </Button>
              </Flex>
            </Card.Header>
            <Card.Body>
              <ExerciseLog
                key={currentCardIndex}
                items={computedLogItems[currentCardIndex].logs}
                onLog={function (
                  setNumber: number,
                  repsCompleted: number,
                  weightUsed: number
                ): void {
                  handleLog(
                    computedLogItems[currentCardIndex].exerciseId,
                    setNumber,
                    repsCompleted,
                    weightUsed
                  );
                }}
              />
            </Card.Body>
            <Card.Footer></Card.Footer>
          </Card.Root>
        )}
      <Flex gap="4" justify="flex-start">
        <Button>Add Note</Button>
      </Flex>
    </Stack>
  );
};

export default EditableWorkoutSession;
