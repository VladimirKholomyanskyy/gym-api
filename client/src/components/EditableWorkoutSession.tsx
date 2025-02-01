import { LogExerciseResponse, WorkoutSessionResponse } from "@/api/models";
import ExerciseLog, { ExerciseLogItem } from "./ExerciseLog";
import { Button } from "./ui/button";

import { Heading, Card, Flex, Box, VStack, IconButton } from "@chakra-ui/react";
import { useState } from "react";
import { ExerciseLogsApi, WorkoutSessionsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";
import { formatDateTime } from "@/utils/dateUtils";
import { FaChevronLeft, FaChevronRight } from "react-icons/fa6";

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
  const exerciseLogsApi = new ExerciseLogsApi(apiConfig);
  // Function to handle "Previous" button
  const handlePrevious = () => {
    setCurrentCardIndex((prevIndex) =>
      prevIndex > 0 ? prevIndex - 1 : computedLogItems.length - 1
    );
  };

  const handleComplete = async () => {
    try {
      await workoutSessionApi.completeWorkoutSession(session.id);
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
      await exerciseLogsApi.logExercise({
        workoutSessionId: session.id,
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
    session?.workoutSnapshot?.workoutExercises?.map((workoutExercise) => {
      const exLogs: ExerciseLogItem[] = Array.from(
        { length: workoutExercise.sets },
        (_, index) => {
          const found = logs?.find(
            (e) => workoutExercise.id === e.id && e.setNumber === index + 1
          );
          return {
            id: index + 1,
            prevReps: 3,
            prevWeight: 20,
            reqReps: workoutExercise.reps,
            currentReps: found?.repsCompleted,
            currentWeight: found?.weightUsed,
          };
        }
      );

      return {
        exerciseId: workoutExercise.exercise?.id,
        exerciseName: workoutExercise.exercise?.name,
        logs: exLogs,
      };
    }) || [];
  console.log("computedLogItems:", computedLogItems);
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
      <Heading size="lg">{formatDateTime(session?.startedAt)}</Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
        {computedLogItems.length > 0 &&
          computedLogItems[currentCardIndex]?.logs && (
            <Card.Root
              size="sm"
              width="100%"
              background="blackAlpha.800"
              borderRadius="md"
              boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
              p={2}
              _hover={{ boxShadow: "0 0 20px rgba(0, 255, 255, 1)" }}
            >
              <Card.Header>
                <Flex justify="space-between" align="center">
                  <IconButton
                    color="neon.500"
                    onClick={handlePrevious}
                    background="transparent"
                  >
                    <FaChevronLeft />
                  </IconButton>
                  <Card.Title
                    fontSize="xl"
                    fontWeight="bold"
                    color="neon.400"
                    textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
                    cursor="pointer"
                    _hover={{ color: "neon.300" }}
                  >
                    {computedLogItems[currentCardIndex].exerciseName}
                  </Card.Title>
                  <IconButton
                    color="neon.500"
                    background="transparent"
                    onClick={handleNext}
                  >
                    <FaChevronRight />
                  </IconButton>
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
                      computedLogItems[currentCardIndex].exerciseId as string,
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
        <Flex gap="4" justify="space-between">
          <Button
            background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
            color="white"
            _hover={{
              filter: "brightness(1.2)",
              boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
            }}
            size="lg"
          >
            Add Note
          </Button>
          <Button
            background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
            color="white"
            _hover={{
              filter: "brightness(1.2)",
              boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
            }}
            size="lg"
            onClick={handleComplete}
          >
            Finish
          </Button>
        </Flex>
      </VStack>
    </Box>
  );
};

export default EditableWorkoutSession;
