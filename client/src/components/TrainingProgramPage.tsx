import React, { useCallback, useEffect, useRef, useState } from "react";
import { Box, VStack, Heading, Flex, Spinner, Input } from "@chakra-ui/react";

import { useNavigate, useParams } from "react-router";
import WorkoutCard from "./WorkoutCard";
import { Button } from "./ui/button";
import { Field } from "./ui/field";
import {
  DrawerActionTrigger,
  DrawerBackdrop,
  DrawerBody,
  DrawerCloseTrigger,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerRoot,
  DrawerTrigger,
} from "./ui/drawer";
import { toaster } from "./ui/toaster";
import { TrainingProgramsApi } from "@/api/apis/training-programs-api";
import { TrainingProgram } from "@/api/models/training-program";
import { apiConfig } from "@/api/apiConfig";
import { ExercisesApi, WorkoutExercisesApi, WorkoutsApi } from "@/api";

interface WorkoutCardProps {
  workoutId: string;
  name: string;
  exercises: string[];
}

const TrainingProgramPage: React.FC = () => {
  const [program, setProgram] = useState<TrainingProgram | null>(null);
  const [workoutExerciseCard, setWorkoutExercisesCard] = useState<
    WorkoutCardProps[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [newWorkout, setNewWorkout] = useState({ name: "" });
  const [newName, setNewName] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const { programId } = useParams();
  const navigate = useNavigate();
  const ref = useRef<HTMLInputElement>(null);
  const trainingProgramApi = new TrainingProgramsApi(apiConfig);
  const workoutApi = new WorkoutsApi(apiConfig);
  const workoutExercisesApi = new WorkoutExercisesApi(apiConfig);
  const exerciseApi = new ExercisesApi(apiConfig);

  if (!programId) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [programData, workoutsData] = await Promise.all([
          trainingProgramApi.getTrainingProgramById(programId),
          workoutApi.listWorkoutsForProgram(programId),
        ]);

        setProgram(programData.data);
        setNewName(programData.data.name);
        setNewDescription(
          programData.data.description ? programData.data.description : ""
        );

        const workoutCards = await Promise.all(
          workoutsData.data.map(async (workout) => {
            const workoutExercises =
              await workoutExercisesApi.listWorkoutExercises(workout.id);

            const exerciseNames = await Promise.all(
              workoutExercises.data.map(async (we) => {
                const exercise = await exerciseApi.getExerciseById(
                  we.exerciseId
                );
                return exercise.data.name;
              })
            );

            return {
              workoutId: workout.id,
              name: workout.name,
              exercises: exerciseNames.filter(Boolean),
            };
          })
        );
        setWorkoutExercisesCard(workoutCards);
      } catch (error) {
        console.error("Failed to load data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [programId]);

  const handleCreateWorkout = useCallback(async () => {
    if (!programId || !newWorkout.name) {
      return;
    }

    try {
      setLoading(true);
      const created = await workoutApi.addWorkoutToProgram(programId, {
        name: newWorkout.name,
      });

      setWorkoutExercisesCard((prev) => [
        ...prev,
        { workoutId: created.data.id, name: created.data.name, exercises: [] },
      ]);
      setNewWorkout({ name: "" });
    } catch (error: any) {
      console.error("Error creating workout:", error);
    } finally {
      setLoading(false);
    }
  }, [programId, newWorkout.name]);

  const handleDeleteWorkout = useCallback(
    async (workoutId: string) => {
      try {
        await workoutApi.deleteWorkout(programId!, workoutId);
        setWorkoutExercisesCard((prev) =>
          prev.filter((card) => card.workoutId !== workoutId)
        );
      } catch (error) {
        console.error("Error deleting workout:", error);
      }
    },
    [programId]
  );

  const handleUpdateProgram = async () => {
    try {
      const program = await trainingProgramApi.updateTrainingProgram(
        programId,
        {
          name: newName,
          description: newDescription,
        }
      );
      setProgram(program.data);
      toaster.create({
        title: "Training program updated.",
        type: "success",
        duration: 3000,
      });
    } catch (error) {
      toaster.create({
        title: "Failed to update training program.",
        description: "Please try again later.",
        type: "error",
        duration: 5000,
      });
    }
  };
  return (
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        {program?.name}
      </Heading>
      <Heading size="lg">{program?.description}</Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
        {loading ? (
          <Flex justifyContent="center" alignItems="center" height="50vh">
            <Spinner size="xl" color="magenta.400" />
          </Flex>
        ) : (
          workoutExerciseCard.map((workout) => (
            <WorkoutCard
              key={workout.workoutId}
              workoutId={workout.workoutId}
              programId={programId}
              onDelete={handleDeleteWorkout}
              name={workout.name}
              exercises={workout.exercises}
            />
          ))
        )}
        <Flex justifyContent="space-between">
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button
                background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
                color="white"
                _hover={{
                  filter: "brightness(1.2)",
                  boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
                }}
                size="lg"
              >
                + Add Workout
              </Button>
            </DrawerTrigger>
            <DrawerContent ref={ref} background="blackAlpha.900">
              <DrawerCloseTrigger />
              <DrawerHeader color="magenta.400">Add a New Workout</DrawerHeader>
              <DrawerBody>
                <VStack gap={4}>
                  <Field label="Name">
                    <Input
                      ref={ref}
                      placeholder="Workout name"
                      value={newWorkout.name}
                      onChange={(e) => setNewWorkout({ name: e.target.value })}
                    />
                  </Field>
                </VStack>
              </DrawerBody>
              <DrawerFooter>
                <DrawerActionTrigger asChild>
                  <Button variant="outline">Cancel</Button>
                </DrawerActionTrigger>
                <Button
                  background="magenta.400"
                  color="black"
                  onClick={handleCreateWorkout}
                >
                  Save
                </Button>
              </DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button
                background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
                color="white"
                _hover={{
                  filter: "brightness(1.2)",
                  boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
                }}
                size="lg"
              >
                Edit Program
              </Button>
            </DrawerTrigger>
            <DrawerContent ref={ref}>
              <DrawerCloseTrigger />
              <DrawerHeader color="magenta.400">
                Edit Training Program
              </DrawerHeader>
              <DrawerBody>
                <VStack gap={4}>
                  <Field label="Name">
                    <Input
                      ref={ref}
                      placeholder="Training program name"
                      value={newName}
                      onChange={(e) => setNewName(e.target.value)}
                    />
                  </Field>
                  <Field label="Description">
                    <Input
                      ref={ref}
                      placeholder="Training program description"
                      value={newDescription}
                      onChange={(e) => setNewDescription(e.target.value)}
                    />
                  </Field>
                </VStack>
              </DrawerBody>
              <DrawerFooter>
                <DrawerActionTrigger asChild>
                  <Button variant="outline">Cancel</Button>
                </DrawerActionTrigger>
                <Button
                  background="magenta.400"
                  color="black"
                  onClick={handleUpdateProgram}
                >
                  Save
                </Button>
              </DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
        </Flex>
      </VStack>
    </Box>
  );
};

export default TrainingProgramPage;
