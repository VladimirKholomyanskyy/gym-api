import React, { useCallback, useEffect, useRef, useState } from "react";
import {
  Box,
  VStack,
  Heading,
  Flex,
  Spinner,
  Input,
  HStack,
  Separator,
} from "@chakra-ui/react";
import { getExercise } from "../api/exercises";
import {
  getTrainingProgram,
  updateTrainingProgram,
} from "../api/trainingPrograms";
import { createWorkout, deleteWorkout, getAllWorkouts } from "../api/workouts";
import { getAllWorkoutExercises } from "../api/workout-exercises";
import { TrainingProgram } from "../types/api";
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

interface WorkoutCardProps {
  workoutId: number;
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
  const programID = Number(programId);
  const navigate = useNavigate();
  const ref = useRef<HTMLInputElement>(null);

  if (!programID) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [programData, workoutsData] = await Promise.all([
          getTrainingProgram(programID),
          getAllWorkouts(programID),
        ]);

        setProgram(programData);
        setNewName(programData.name);
        setNewDescription(programData.description);

        const workoutCards = await Promise.all(
          workoutsData.map(async (workout) => {
            const workoutExercises = await getAllWorkoutExercises(workout.id);

            const exerciseNames = await Promise.all(
              workoutExercises.map(async (we) => {
                const exercise = await getExercise(we.exercise_id);
                return exercise.Name;
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
      const created = await createWorkout(programID, { name: newWorkout.name });

      setWorkoutExercisesCard((prev) => [
        ...prev,
        { workoutId: created.id, name: created.name, exercises: [] },
      ]);
      setNewWorkout({ name: "" });
    } catch (error: any) {
      console.error("Error creating workout:", error);
    } finally {
      setLoading(false);
    }
  }, [programId, newWorkout.name]);

  const handleDeleteWorkout = useCallback(
    async (workoutId: number) => {
      try {
        await deleteWorkout(programID!, workoutId);
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
      const program = await updateTrainingProgram(programID, {
        name: newName,
        description: newDescription,
      });
      setProgram(program);
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
    <Box p={5}>
      <VStack gap={4} minHeight="100%" align="stretch">
        <Heading size="4xl">{program?.name}</Heading>
        <Heading size="2xl">{program?.description}</Heading>
        <Separator />
        <HStack>
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button colorScheme="teal" aria-label="Add Exercise" size="lg">
                Add Workout
              </Button>
            </DrawerTrigger>
            <DrawerContent ref={ref}>
              <DrawerCloseTrigger />
              <DrawerHeader>Add a New Workout</DrawerHeader>

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
                <Button colorScheme="teal" onClick={handleCreateWorkout}>
                  Save
                </Button>
              </DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button colorScheme="teal" aria-label="Add Exercise" size="lg">
                Edit Program
              </Button>
            </DrawerTrigger>
            <DrawerContent ref={ref}>
              <DrawerCloseTrigger />
              <DrawerHeader>Edit Training Program</DrawerHeader>

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
                <Button colorScheme="teal" onClick={handleUpdateProgram}>
                  Save
                </Button>
              </DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
        </HStack>
        {loading ? (
          <Flex justifyContent="center" alignItems="center" height="100%">
            <Spinner />
          </Flex>
        ) : (
          workoutExerciseCard.map((workout) => (
            <WorkoutCard
              key={workout.workoutId}
              workoutId={workout.workoutId}
              programId={programID}
              onDelete={handleDeleteWorkout}
              name={workout.name}
              exercises={workout.exercises}
            />
          ))
        )}
      </VStack>
    </Box>
  );
};

export default TrainingProgramPage;
