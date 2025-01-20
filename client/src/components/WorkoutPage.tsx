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
import {
  Exercise,
  WorkoutExerciseRequest,
  WorkoutResponse,
} from "@/api/models";
import {
  ExercisesApi,
  WorkoutExercisesApi,
  WorkoutsApi,
  WorkoutSessionsApi,
} from "@/api";
import { apiConfig } from "@/api/apiConfig";
import ConfirmationDialog from "./common/ConfirmationDialog";

interface ExerciseCardProps {
  id: string;
  exerciseName: string;
  exerciseId: string;
  sets: number;
  reps: number;
}

const WorkoutPage: React.FC = () => {
  const [exercises, setExercises] = useState<ExerciseCardProps[]>([]);
  const [exercisesSelect, setExercisesSelect] = useState<Exercise[]>([]);
  const [workout, setWorkout] = useState<WorkoutResponse | null>(null);
  const { programId } = useParams();
  const { workoutId } = useParams();
  const [loading, setLoading] = useState(true);
  const [selectedExerciseId, setSelectedExerciseId] = useState<string[]>([]);
  const [sets, setSets] = useState<number>(3);
  const [reps, setReps] = useState<number>(10);
  const navigate = useNavigate();
  const contentRef = useRef<HTMLDivElement>(null);
  const workoutApi = new WorkoutsApi(apiConfig);
  const workoutExercisesApi = new WorkoutExercisesApi(apiConfig);
  const exerciseApi = new ExercisesApi(apiConfig);
  const workoutSessionApi = new WorkoutSessionsApi(apiConfig);

  if (!programId || !workoutId) {
    navigate("/error");
    return null;
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        const workout = await workoutApi.getWorkoutForProgram(
          programId,
          workoutId
        );
        const workoutExercises = await workoutExercisesApi.listWorkoutExercises(
          workoutId
        );
        const exercises = await Promise.all(
          workoutExercises.data.map(async (workoutExercise) => {
            const exercise = await exerciseApi.getExerciseById(
              workoutExercise.exerciseId
            );

            return {
              id: workoutExercise.id,
              exerciseName: exercise.data.name,
              exerciseId: exercise.data.id,
              sets: workoutExercise.sets,
              reps: workoutExercise.reps,
            };
          })
        );
        const exercisesData = await exerciseApi.listExercises();
        setWorkout(workout.data);
        setExercises(exercises);
        setExercisesSelect(exercisesData.data);
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

    const requestData: WorkoutExerciseRequest = {
      exerciseId: selectedExerciseId[0],
      workoutId: workoutId,
      sets,
      reps,
    };

    try {
      await workoutExercisesApi.postWorkoutExercise(requestData);
      const workoutExercises = await workoutExercisesApi.listWorkoutExercises(
        workoutId
      );
      const updatedExercises = await Promise.all(
        workoutExercises.data.map(async (workoutExercise) => {
          const exercise = await exerciseApi.getExerciseById(
            workoutExercise.exerciseId
          );
          return {
            id: workoutExercise.id,
            exerciseName: exercise.data.name,
            exerciseId: exercise.data.id,
            sets: workoutExercise.sets,
            reps: workoutExercise.reps,
          };
        })
      );
      setExercises(updatedExercises);
    } catch (error) {}
  };

  const handleEdit = async (
    workoutExerciseId: string,
    exerciseId: string,
    sets: number,
    reps: number
  ) => {
    const requestData: WorkoutExerciseRequest = {
      exerciseId: exerciseId,
      workoutId: workoutId,
      sets,
      reps,
    };

    try {
      await workoutExercisesApi.patchWorkoutExercise(
        workoutExerciseId,
        requestData
      );
    } catch (error) {}
  };
  const handleDeleteWorkoutExercise = async (id: string) => {
    try {
      await workoutExercisesApi.deleteWorkoutExercise(id);
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

  const handleNavigate = async () => {
    const response = await workoutSessionApi.addWorkoutSession({
      workoutId: workoutId,
    });
    navigate(`/workout-sessions/${response.data.id}/edit`);
  };
  return (
    <Box>
      <Box mb={7} mt={7}>
        <Heading size="2xl" fontWeight="bold" textAlign="center">
          {workout?.name}
        </Heading>
      </Box>
      <VStack
        gap={6}
        align="stretch"
        width="100%"
        paddingLeft="8"
        paddingRight="8"
      >
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
                exerciseId: string,
                reps: number,
                sets: number
              ): void {
                handleEdit(exercise.id, exerciseId, sets, reps);
              }}
              exerciseId={exercise.exerciseId}
            />
          ))
        )}
        <Flex gap="4">
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button colorScheme="teal" aria-label="Add Exercise" size="lg">
                Add Exercise
              </Button>
            </DrawerTrigger>
            <ConfirmationDialog
              triggerLabel="Start Workout"
              triggerIcon={null}
              deleteLabel="Yes"
              cancelLabel="No"
              message={"Are you sure wnat to start workout"}
              onDelete={handleNavigate}
            />
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
        </Flex>
      </VStack>
    </Box>
  );
};

export default WorkoutPage;
