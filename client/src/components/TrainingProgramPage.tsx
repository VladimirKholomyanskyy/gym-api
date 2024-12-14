import React, { useCallback, useEffect, useRef, useState } from "react";
import {
  Box,
  VStack,
  Heading,
  Flex,
  Spinner,
  Stack,
  Input,
} from "@chakra-ui/react";
import { getExercise, listExercises } from "../api/exercises";
import { getTrainingProgram } from "../api/trainingPrograms";
import {
  createWorkout,
  deleteWorkout,
  getAllWorkouts,
  updateWorkout,
} from "../api/workouts";
import { getAllWorkoutExercises } from "../api/workout-exercises";
import { Exercise, TrainingProgram, Workout } from "../types/api";
import { useAuth } from "react-oidc-context";
import { setupAxiosInterceptors } from "@/api/apiClient";
import { useNavigate, useParams } from "react-router";
import WorkoutCard from "./WorkoutCard";
import {
  DialogActionTrigger,
  DialogBody,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Button } from "./ui/button";
import { Field } from "./ui/field";

interface WorkoutCardProps {
  workoutId: number;
  name: string;
  exercises: string[];
}

const TrainingProgramPage: React.FC = () => {
  const [program, setProgram] = useState<TrainingProgram | null>(null);
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [workoutExerciseCard, setWorkoutExercisesCard] = useState<
    WorkoutCardProps[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [newWorkout, setNewWorkout] = useState({ name: "" });

  const { programId } = useParams();
  const programID = Number(programId);
  const auth = useAuth();
  const navigate = useNavigate();
  const ref = useRef<HTMLInputElement>(null);

  if (!programID) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    setupAxiosInterceptors(() => auth.user?.access_token || null);

    const fetchData = async () => {
      try {
        const [exercisesData, programData, workoutsData] = await Promise.all([
          listExercises(),
          getTrainingProgram(programID),
          getAllWorkouts(programID),
        ]);

        setExercises(exercisesData);
        setProgram(programData);
        setWorkouts(workoutsData);

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
        setMessage("Failed to load training program details.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [auth.user?.access_token, programId]);

  const handleCreateWorkout = useCallback(async () => {
    if (!programId || !newWorkout.name) {
      setMessage("Please provide a name for the workout.");
      return;
    }

    try {
      setLoading(true);
      const created = await createWorkout(programID, { name: newWorkout.name });

      setWorkouts((prev) => [...prev, created]);
      setWorkoutExercisesCard((prev) => [
        ...prev,
        { workoutId: created.id, name: created.name, exercises: [] },
      ]);
      setNewWorkout({ name: "" });
      setMessage("Workout created successfully!");
    } catch (error: any) {
      console.error("Error creating workout:", error);
      setMessage(error.response?.data?.message || "Failed to create workout.");
    } finally {
      setLoading(false);
    }
  }, [programId, newWorkout.name]);

  const handleUpdateWorkout = useCallback(
    async (workoutId: number, name: string) => {
      try {
        setLoading(true);
        console.log(
          `workoutId=${workoutId}, programId=${programId}, name=${name}`
        );
        const updated = await updateWorkout(programID, workoutId, { name });

        setWorkouts((prev) =>
          prev.map((workout) =>
            workout.id === workoutId
              ? { ...workout, name: updated.name }
              : workout
          )
        );
        setWorkoutExercisesCard((prev) =>
          prev.map((card) =>
            card.workoutId === workoutId
              ? { ...card, name: updated.name }
              : card
          )
        );
        setMessage("Workout updated successfully!");
      } catch (error: any) {
        console.error("Error updating workout:", error);
        setMessage(
          error.response?.data?.message || "Failed to update workout."
        );
      } finally {
        setLoading(false);
      }
    },
    [programId]
  );

  const handleDeleteWorkout = useCallback(
    async (workoutId: number) => {
      try {
        await deleteWorkout(programID!, workoutId);
        setWorkouts((prev) =>
          prev.filter((workout) => workout.id !== workoutId)
        );
        setWorkoutExercisesCard((prev) =>
          prev.filter((card) => card.workoutId !== workoutId)
        );
      } catch (error) {
        console.error("Error deleting workout:", error);
      }
    },
    [programId]
  );

  return (
    <Box p={5}>
      <VStack gap={4} minHeight="100%">
        <Heading size="4xl">{program?.name}</Heading>
        <Heading size="2xl">{program?.description}</Heading>
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
              onUpdate={handleUpdateWorkout}
              name={workout.name}
              exercises={workout.exercises}
            />
          ))
        )}
        <DialogRoot initialFocusEl={() => ref.current}>
          <DialogTrigger asChild>
            <Button
              position="fixed"
              bottom="16px"
              right="16px"
              variant="outline"
              bg="green"
            >
              New
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>New Workout</DialogTitle>
            </DialogHeader>
            <DialogBody pb="4">
              <Stack gap="4">
                <Field label="Name">
                  <Input
                    ref={ref}
                    placeholder="Workout name"
                    value={newWorkout.name}
                    onChange={(e) => setNewWorkout({ name: e.target.value })}
                  />
                </Field>
              </Stack>
            </DialogBody>
            <DialogFooter>
              <DialogActionTrigger asChild>
                <Button variant="outline">Cancel</Button>
              </DialogActionTrigger>
              <DialogActionTrigger asChild>
                <Button onClick={handleCreateWorkout}>Save</Button>
              </DialogActionTrigger>
            </DialogFooter>
          </DialogContent>
        </DialogRoot>
      </VStack>
    </Box>
  );
};

export default TrainingProgramPage;
