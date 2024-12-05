import React, { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { Box, VStack, Input, Text } from "@chakra-ui/react";
import { Button } from "./ui/button";
import { addExerciseToWorkout } from "../api/workout-exercises";
import { AddExerciseRequest } from "../types/api";
import { useAuth } from "react-oidc-context";

interface Props {
  workoutId: string;
}

const AddExerciseToWorkout: React.FC<Props> = ({ workoutId }) => {
  const { control, handleSubmit } =
    useForm<Omit<AddExerciseRequest, "workoutId">>();
  const [message, setMessage] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const { user, isAuthenticated } = useAuth();
  if (!isAuthenticated || !user) return;

  const onSubmit = async (data: Omit<AddExerciseRequest, "workoutId">) => {
    setLoading(true);
    setMessage("");

    try {
      await addExerciseToWorkout(user, { ...data, workout_id: workoutId });
      setMessage("Exercise Added Successfully!");
    } catch (error: any) {
      setMessage(error.response?.data?.message || "Failed to add exercise.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box p={5} borderWidth={1} borderRadius="lg">
      <VStack as="form" gap={4} onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="exercise_id"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input {...field} placeholder="Exercise ID" required />
          )}
        />
        <Controller
          name="sets"
          control={control}
          defaultValue={0}
          render={({ field }) => (
            <Input {...field} type="number" placeholder="Sets" required />
          )}
        />
        <Controller
          name="reps"
          control={control}
          defaultValue={0}
          render={({ field }) => (
            <Input {...field} type="number" placeholder="Reps" required />
          )}
        />
        <Button colorScheme="teal" type="submit" loading={loading}>
          Add Exercise
        </Button>
        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

export default AddExerciseToWorkout;
