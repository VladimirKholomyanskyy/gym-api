import { setupAxiosInterceptors } from "@/api/apiClient";
import { getExercise, listExercises } from "@/api/exercises";
import {
  addExerciseToWorkout,
  getAllWorkoutExercises,
} from "@/api/workout-exercises";
import { getWorkout } from "@/api/workouts";
import { AddExerciseRequest, Exercise, Workout } from "@/types/api";
import {
  Box,
  Flex,
  Heading,
  IconButton,
  NumberInputRoot,
  Spinner,
  VStack,
} from "@chakra-ui/react";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "react-oidc-context";
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
import { FaPlus } from "react-icons/fa";
import ExerciseSelect from "./ExerciseSelect";

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
  const [message, setMessage] = useState("");
  const [selectedExerciseId, setSelectedExerciseId] = useState<string[]>([]);
  const [sets, setSets] = useState<number>(3);
  const [reps, setReps] = useState<number>(10);
  const auth = useAuth();
  const navigate = useNavigate();
  const contentRef = useRef<HTMLDivElement>(null);

  if (!programId || !workoutId) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        setupAxiosInterceptors(() => auth.user?.access_token || null);
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
        setMessage("Failed to load workout details.");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [auth.user?.access_token, programId, workoutId]);

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

  return (
    <Box p={5}>
      <VStack gap={4} minHeight="100%">
        <Heading size="4xl">{workout?.name}</Heading>
        {loading ? (
          <Flex justifyContent="center" alignItems="center" height="100%">
            <Spinner />
          </Flex>
        ) : (
          exercises.map((exercise) => (
            <ExerciseCard
              key={exercise.id}
              exercise={exercise.exerciseName}
              sets={exercise.sets}
              reps={exercise.reps}
            />
          ))
        )}
      </VStack>
      {/* Drawer */}
      <DrawerRoot placement="bottom">
        <DrawerBackdrop />
        <DrawerTrigger asChild>
          <IconButton
            colorScheme="teal"
            aria-label="Add Exercise"
            position="fixed"
            bottom={4}
            right={4}
            size="lg"
            borderRadius="full"
          >
            <FaPlus />
          </IconButton>
        </DrawerTrigger>
        <DrawerContent ref={contentRef}>
          <DrawerCloseTrigger />
          <DrawerHeader>Add a New Exercise</DrawerHeader>

          <DrawerBody>
            <VStack gap={4}>
              <ExerciseSelect
                exercises={exercisesSelect}
                contentRef={contentRef}
                selectedExerciseId={selectedExerciseId}
                setSelectedExerciseId={setSelectedExerciseId}
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
