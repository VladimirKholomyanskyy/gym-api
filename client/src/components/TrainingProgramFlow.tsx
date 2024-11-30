import React, { useEffect, useState } from "react";
import {
  Box,
  VStack,
  HStack,
  Input,
  Text,
  Separator,
  createListCollection,
} from "@chakra-ui/react";
import {
  SelectRoot,
  SelectTrigger,
  SelectValueText,
  SelectContent,
  SelectItem,
} from "./ui/select";
import { Button } from './ui/button';
import { listExercises } from "../api/exercises";
import { createTrainingProgram } from "../api/trainingPrograms";
import { createWorkout } from "../api/workouts";
import { addExerciseToWorkout } from "../api/workout-exercises";
import {
  CreateTrainingProgramRequest,
  CreateWorkoutRequest,
  AddExerciseRequest,
  Exercise,
} from "../types/api";

const TrainingProgramFlow: React.FC = () => {
  const [programName, setProgramName] = useState("");
  const [programDescription, setProgramDescription] = useState("");
  const [workoutName, setWorkoutName] = useState("");
  const [selectedExerciseId, setSelectedExerciseId] = useState("");
  const [sets, setSets] = useState<number>(0);
  const [reps, setReps] = useState<number>(0);
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [programId, setProgramId] = useState<string | null>(null);
  const [workoutId, setWorkoutId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const fetchAllExercises = async () => {
      const data = await listExercises();
      console.log(data)
      setExercises(data);
    };
    fetchAllExercises();
  }, []);

  const handleCreateTrainingProgram = async () => {
    setLoading(true);
    setMessage("");

    try {
      const request: CreateTrainingProgramRequest = {
        name: programName,
        description: programDescription,
      };
      const program = await createTrainingProgram(request);
      setProgramId(program.ID);
      setMessage("Training program created successfully!");
    } catch (error: any) {
      setMessage(error.response?.data?.message || "Failed to create program.");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateWorkout = async () => {
    if (!programId) {
      setMessage("Please create a training program first.");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const request: CreateWorkoutRequest = { name: workoutName };
      const workout = await createWorkout(programId, request);
      setWorkoutId(workout.id);
      setMessage("Workout created successfully!");
    } catch (error: any) {
      setMessage(error.response?.data?.message || "Failed to create workout.");
    } finally {
      setLoading(false);
    }
  };

  const handleAddExercise = async () => {
    if (!workoutId) {
      setMessage("Please create a workout first.");
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const request: AddExerciseRequest = {
        exerciseId: selectedExerciseId,
        workoutId,
        sets,
        reps,
      };
      await addExerciseToWorkout(request);
      setMessage("Exercise added successfully!");
    } catch (error: any) {
      setMessage(error.response?.data?.message || "Failed to add exercise.");
    } finally {
      setLoading(false);
    }
  };
  console.log(exercises)
  return (
    <Box p={5}>
      <VStack gap={4}>
        <Text fontSize="lg" fontWeight="bold">
          Training Program Flow
        </Text>

        <Input
          placeholder="Program Name"
          value={programName}
          onChange={(e) => setProgramName(e.target.value)}
        />
        <Input
          placeholder="Program Description"
          value={programDescription}
          onChange={(e) => setProgramDescription(e.target.value)}
        />
        <Button
          colorScheme="teal"
          onClick={handleCreateTrainingProgram}
          loading={loading}
        >
          Create Training Program
        </Button>

        <Separator />

        <Input
          placeholder="Workout Name"
          value={workoutName}
          onChange={(e) => setWorkoutName(e.target.value)}
        />
        <Button
          colorScheme="blue"
          onClick={handleCreateWorkout}
          loading={loading}
        >
          Create Workout
        </Button>

        <Separator />

        <SelectRoot collection={createListCollection({items: exercises})}>
          <SelectTrigger clearable>
            <SelectValueText placeholder="Select Exercise" />
          </SelectTrigger>
          <SelectContent>
            {Array.isArray(exercises) && exercises.map((exercise) => (
              <SelectItem
                key={exercise.ID}
                item={{ value: exercise.ID }}
                onClick={() => setSelectedExerciseId(exercise.ID)}
              >
                {exercise.Name}
              </SelectItem>
            ))}
          </SelectContent>
        </SelectRoot>

        <HStack gap={2}>
          <Input
            type="number"
            placeholder="Sets"
            value={sets}
            onChange={(e) => setSets(Number(e.target.value))}
          />
          <Input
            type="number"
            placeholder="Reps"
            value={reps}
            onChange={(e) => setReps(Number(e.target.value))}
          />
        </HStack>
        <Button
          colorScheme="purple"
          onClick={handleAddExercise}
          loading={loading}
        >
          Add Exercise
        </Button>

        {message && <Text color="red.500">{message}</Text>}
      </VStack>
    </Box>
  );
};

export default TrainingProgramFlow;
