import {
  Box,
  Flex,
  Heading,
  Spinner,
  Stack,
  VStack,
  Text,
  NumberInputRoot,
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
import {
  DialogActionTrigger,
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Field } from "./ui/field";

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
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        {workout?.name}
      </Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
        {loading ? (
          <Flex justifyContent="center" alignItems="center" height="50vh">
            <Spinner size="xl" color="magenta.400" />
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
                + Add Exercise
              </Button>
            </DrawerTrigger>
            <DialogRoot role="alertdialog">
              <DialogTrigger asChild>
                <Button
                  background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
                  color="white"
                  _hover={{
                    filter: "brightness(1.2)",
                    boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
                  }}
                  size="lg"
                >
                  Start Workout
                </Button>
              </DialogTrigger>
              <DialogContent
                background="blackAlpha.900"
                border="1px solid"
                borderColor="neon.400"
                boxShadow="0 0 15px rgba(0, 255, 255, 0.8)"
                p={4}
              >
                <DialogHeader>
                  <DialogTitle
                    color="neon.400"
                    textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
                  >
                    Start Workout
                  </DialogTitle>
                </DialogHeader>
                <DialogBody>
                  <Text color="gray.300">
                    Are you sure wnat to start workout
                  </Text>
                </DialogBody>
                <DialogFooter>
                  <Stack direction="row" gap={4}>
                    <DialogActionTrigger asChild>
                      <Button
                        variant="outline"
                        borderColor="neon.400"
                        color="neon.400"
                        _hover={{ borderColor: "neon.300", color: "neon.300" }}
                      >
                        No
                      </Button>
                    </DialogActionTrigger>
                    <DialogActionTrigger asChild>
                      <Button
                        background="red.600"
                        color="white"
                        _hover={{
                          background: "red.400",
                          boxShadow: "0 0 15px red",
                        }}
                        onClick={handleNavigate}
                      >
                        Yes
                      </Button>
                    </DialogActionTrigger>
                  </Stack>
                </DialogFooter>
                <DialogCloseTrigger />
              </DialogContent>
            </DialogRoot>
            <DrawerContent ref={contentRef}>
              <DrawerCloseTrigger />
              <DrawerHeader color="magenta.400">
                Add a New Exercise
              </DrawerHeader>
              <DrawerBody>
                <VStack gap={4}>
                  <ExerciseSelect
                    exercises={exercisesSelect}
                    contentRef={contentRef}
                    setSelectedExerciseId={setSelectedExerciseId}
                    defaultExerciseId={""}
                  />
                  <Flex gap="4">
                    <Field label="Sets">
                      <NumberInputRoot
                        defaultValue="3"
                        min={1}
                        onValueChange={(e) => setSets(e.valueAsNumber)}
                      >
                        <NumberInputField placeholder="Number of Sets" />
                      </NumberInputRoot>
                    </Field>
                    <Field label="Reps">
                      <NumberInputRoot
                        defaultValue="10"
                        min={1}
                        onValueChange={(e) => setReps(e.valueAsNumber)}
                      >
                        <NumberInputField placeholder="Number of Reps" />
                      </NumberInputRoot>
                    </Field>
                  </Flex>
                </VStack>
              </DrawerBody>

              <DrawerFooter>
                <DrawerActionTrigger asChild>
                  <Button
                    variant="outline"
                    borderColor="neon.400"
                    color="neon.400"
                    _hover={{
                      borderColor: "neon.300",
                      color: "neon.300",
                    }}
                  >
                    Cancel
                  </Button>
                </DrawerActionTrigger>
                <Button
                  background="neon.400"
                  color="black"
                  onClick={handleSubmit}
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

export default WorkoutPage;
