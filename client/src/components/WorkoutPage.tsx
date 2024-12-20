import { getExercise, listExercises } from "@/api/exercises";
import {
  addExerciseToWorkout,
  deleteWorkoutExercise,
  getAllWorkoutExercises,
  patchWorkoutExercise,
} from "@/api/workout-exercises";
import { getWorkout } from "@/api/workouts";
import { AddExerciseRequest, Exercise, Workout } from "@/types/api";
import {
  Box,
  Flex,
  Heading,
  NumberInputRoot,
  Spinner,
  VStack,
} from "@chakra-ui/react";
import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router";
import ExerciseCard from "./ExercisesCard";
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
import { Button } from "./ui/button";
import { NumberInputField } from "./ui/number-input";
import ExerciseSelect from "./ExerciseSelect";
import { toaster } from "./ui/toaster";

interface ExerciseCardProps {
  id: number;
  exerciseName: string;
  exerciseId: number;
  sets: number;
  reps: number;
}

const WorkoutPage: React.FC = () => {
  const [exercises, setExercises] = useState<ExerciseCardProps[]>([]);
  const [exercisesSelect, setExercisesSelect] = useState<Exercise[]>([]);
  const [workout, setWorkout] = useState<Workout | null>(null);
  const { programId } = useParams();
  const { workoutId } = useParams();
  const programID = Number(programId);
  const workoutID = Number(workoutId);
  const [loading, setLoading] = useState(true);
  const [selectedExerciseId, setSelectedExerciseId] = useState<string[]>([]);
  const [sets, setSets] = useState<number>(3);
  const [reps, setReps] = useState<number>(10);
  const navigate = useNavigate();
  const contentRef = useRef<HTMLDivElement>(null);

  if (!programId || !workoutId) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        const workout = await getWorkout(programID, workoutID);
        const workoutExercises = await getAllWorkoutExercises(workoutID);
        const exercises = await Promise.all(
          workoutExercises.map(async (workoutExercise) => {
            const exercise = await getExercise(workoutExercise.exercise_id);

            return {
              id: workoutExercise.id,
              exerciseName: exercise.Name,
              exerciseId: exercise.ID,
              sets: workoutExercise.sets,
              reps: workoutExercise.reps,
            };
          })
        );
        const exercisesData = await listExercises();
        setWorkout(workout);
        setExercises(exercises);
        setExercisesSelect(exercisesData);
      } catch (error) {
        console.error("Failed to load data:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [programId, workoutId]);

  const handleSubmit = async () => {
    if (!selectedExerciseId) {
      return;
    }

    const requestData: AddExerciseRequest = {
      exercise_id: Number(selectedExerciseId),
      workout_id: Number(workoutId),
      sets,
      reps,
    };

    try {
      await addExerciseToWorkout(requestData);
      const workoutExercises = await getAllWorkoutExercises(workoutID);
      const updatedExercises = await Promise.all(
        workoutExercises.map(async (workoutExercise) => {
          const exercise = await getExercise(workoutExercise.exercise_id);
          return {
            id: workoutExercise.id,
            exerciseName: exercise.Name,
            exerciseId: exercise.ID,
            sets: workoutExercise.sets,
            reps: workoutExercise.reps,
          };
        })
      );
      setExercises(updatedExercises);
    } catch (error) {}
  };

  const handleEdit = async (
    workoutExerciseId: number,
    exerciseId: number,
    sets: number,
    reps: number
  ) => {
    const requestData: AddExerciseRequest = {
      exercise_id: Number(exerciseId),
      workout_id: Number(workoutId),
      sets,
      reps,
    };

    try {
      await patchWorkoutExercise(workoutExerciseId, requestData);
    } catch (error) {}
  };
  const handleDeleteWorkoutExercise = async (id: number) => {
    try {
      await deleteWorkoutExercise(id);
      toaster.create({
        title: "Removed exercise.",
        type: "success",
        duration: 3000,
      });
    } catch (error) {
      toaster.create({
        title: "Failed to remove exercise.",
        description: "Please try again later.",
        type: "error",
        duration: 5000,
      });
    }
  };
  return (
    <Box p={5}>
      <VStack gap={6} align="stretch">
        <Heading size="4xl">{workout?.name}</Heading>
        {loading ? (
          <Flex justifyContent="center" alignItems="center" height="50vh">
            <Spinner />
          </Flex>
        ) : (
          exercises.map((exercise) => (
            <ExerciseCard
              key={exercise.id}
              exercise={exercise.exerciseName}
              sets={exercise.sets}
              reps={exercise.reps}
              contentRef={contentRef}
              onDelete={() => handleDeleteWorkoutExercise(exercise.id)}
              exercises={exercisesSelect}
              onEdit={function (
                exerciseId: number,
                reps: number,
                sets: number
              ): void {
                handleEdit(exercise.id, exerciseId, sets, reps);
              }}
              exerciseId={exercise.exerciseId}
            />
          ))
        )}
      </VStack>
      {/* Drawer */}
      <DrawerRoot placement="bottom">
        <DrawerBackdrop />
        <DrawerTrigger asChild>
          <Button colorScheme="teal" aria-label="Add Exercise" size="lg">
            Add Exercise
          </Button>
        </DrawerTrigger>
        <DrawerContent ref={contentRef}>
          <DrawerCloseTrigger />
          <DrawerHeader>Add a New Exercise</DrawerHeader>

          <DrawerBody>
            <VStack gap={4}>
              <ExerciseSelect
                exercises={exercisesSelect}
                contentRef={contentRef}
                setSelectedExerciseId={setSelectedExerciseId}
                defaultExerciseId={""}
              />
              <NumberInputRoot
                defaultValue="3"
                min={1}
                onValueChange={(e) => setSets(e.valueAsNumber)}
              >
                <NumberInputField placeholder="Number of Sets" />
              </NumberInputRoot>

              <NumberInputRoot
                defaultValue="10"
                min={1}
                onValueChange={(e) => setReps(e.valueAsNumber)}
              >
                <NumberInputField placeholder="Number of Reps" />
              </NumberInputRoot>
            </VStack>
          </DrawerBody>

          <DrawerFooter>
            <DrawerActionTrigger asChild>
              <Button variant="outline">Cancel</Button>
            </DrawerActionTrigger>
            <Button colorScheme="teal" onClick={handleSubmit}>
              Save
            </Button>
          </DrawerFooter>
        </DrawerContent>
      </DrawerRoot>
    </Box>
  );
};

export default WorkoutPage;
