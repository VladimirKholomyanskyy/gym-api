import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";
import ExerciseLog, { ExerciseLogItem } from "./ExerciseLog";
import { Button } from "./ui/button";

import { Stack, Heading, Card, Flex } from "@chakra-ui/react";
import { useState } from "react";
import { WorkoutSessionsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";

const EditableWorkoutSession = ({
  session,
  logs,
}: {
  session: WorkoutSessionResponse;
  logs: LogExerciseResponse[];
}) => {
  const [currentCardIndex, setCurrentCardIndex] = useState(0);

  const handleNext = () => {
    setCurrentCardIndex((prevIndex) => {
      return prevIndex < computedLogItems.length - 1 ? prevIndex + 1 : 0;
    });
  };
  const workoutSessionApi = new WorkoutSessionsApi(apiConfig);
  // Function to handle "Previous" button
  const handlePrevious = () => {
    setCurrentCardIndex((prevIndex) =>
      prevIndex > 0 ? prevIndex - 1 : computedLogItems.length - 1
    );
  };

  const handleComplete = async () => {
    try {
      await workoutSessionApi.finishWorkoutSession(session.id);
    } catch (error) {
      console.log(error);
    }
  };
  const handleLog = async (
    exerciseId: string,
    setNumber: number,
    repsCompleted: number,
    weightUsed: number
  ) => {
    try {
      await workoutSessionApi.logExercise(session.id, {
        exerciseId: exerciseId,
        repsCompleted: repsCompleted,
        setNumber: setNumber,
        weightUsed: weightUsed,
      });
    } catch (error) {
      console.log(error);
    }
  };

  const computedLogItems =
    session?.workoutSnapshot?.workoutExercises?.map((exercise) => {
      const exLogs: ExerciseLogItem[] = Array.from(
        { length: exercise.sets },
        (_, index) => {
          const found = logs.find(
            (e) => exercise.id === e.id && e.setNumber === index + 1
          );
          return {
            id: index + 1,
            prevReps: 3,
            prevWeight: 20,
            reqReps: exercise.reps,
            currentReps: found?.repsCompleted,
            currentWeight: found?.weightUsed,
          };
        }
      );

      return {
        exerciseId: exercise.id,
        exerciseName: exercise.exercise?.name,
        logs: exLogs,
      };
    }) || [];
  console.log("computedLogItems:", computedLogItems);
  return (
    <Stack>
      <Heading size="4xl">{session?.workoutSnapshot.name}</Heading>
      <Heading size="2xl">{session?.startedAt}</Heading>
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
