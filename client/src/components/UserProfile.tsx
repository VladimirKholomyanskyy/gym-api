import React, { useState } from 'react';
import { useAuth } from 'react-oidc-context';
import axios from 'axios';
import {
  Box,
  Button,
  Input,
  Textarea,
  Text,
  VStack,
  HStack,
  Spinner,
} from '@chakra-ui/react';
import { useForm, Controller } from 'react-hook-form';

const CreateTrainingProgram = () => {
  const { user, isAuthenticated } = useAuth();
  const { control, handleSubmit, setValue } = useForm();
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const onSubmit = async (data: any) => {
    if (isAuthenticated && user) {
      const token = user.access_token; // Access the access token from user
      const programData = {
        name: data.name,
        description: data.description,
      };

      setLoading(true);

      try {
        const response = await axios.post(
          '/api/training-programs',
          programData,
          {
            headers: {
              Authorization: `Bearer ${token}`, // Add Bearer token to the headers
            },
          }
        );
        setMessage('Training Program Created Successfully!');
        console.log('Program created:', response.data);
      } catch (error) {
        console.error('Error creating training program:', error);
        setMessage('Failed to create training program.');
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box p={5} borderWidth={1} borderRadius="lg">
      <VStack gap={4} align="flex-start" as="form" onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="name"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input
              {...field}
              type="text"
              placeholder="Enter program name"
              required
            />
          )}
        />
        <Controller
          name="description"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Textarea
              {...field}
              placeholder="Enter program description"
              required
            />
          )}
        />

        <HStack gap={4}>
          <Button colorScheme="teal" type="submit" isLoading={loading}>
            {loading ? <Spinner size="sm" /> : 'Create Program'}
          </Button>
        </HStack>

        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

const CreateWorkout = ({ programId }: { programId: string }) => {
  const { user, isAuthenticated } = useAuth();
  const { control, handleSubmit, setValue } = useForm();
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const onSubmit = async (data: any) => {
    if (isAuthenticated && user) {
      const token = user.access_token; // Access the access token from user
      const workoutData = {
        name: data.name,
        programId,
      };

      setLoading(true);

      try {
        const response = await axios.post(
          `/api/training-programs/${programId}/workouts`,
          workoutData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        setMessage('Workout Created Successfully!');
        console.log('Workout created:', response.data);
      } catch (error) {
        console.error('Error creating workout:', error);
        setMessage('Failed to create workout.');
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box p={5} borderWidth={1} borderRadius="lg">
      <VStack gap={4} align="flex-start" as="form" onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="name"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input
              {...field}
              type="text"
              placeholder="Enter workout name"
              required
            />
          )}
        />

        <HStack gap={4}>
          <Button colorScheme="teal" type="submit"  isLoading={loading}>
            {loading ? <Spinner size="sm" /> : 'Create Workout'}
          </Button>
        </HStack>

        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

const AddExerciseToWorkout = ({ workoutId }: { workoutId: string }) => {
  const { user, isAuthenticated } = useAuth();
  const { control, handleSubmit } = useForm();
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const onSubmit = async (data: any) => {
    if (isAuthenticated && user) {
      const token = user.access_token;
      const exerciseData = {
        exerciseId: data.exerciseId,
        sets: data.sets,
        reps: data.reps,
      };

      setLoading(true);

      try {
        const response = await axios.post(
          '/api/workout-exercises',
          exerciseData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        setMessage('Exercise Added Successfully!');
        console.log('Exercise added:', response.data);
      } catch (error) {
        console.error('Error adding exercise to workout:', error);
        setMessage('Failed to add exercise to workout.');
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box p={5} borderWidth={1} borderRadius="lg">
      <VStack gap={4} align="flex-start" as="form" onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="exerciseId"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input
              {...field}
              type="text"
              placeholder="Enter exercise ID"
              required
            />
          )}
        />

        <Controller
          name="sets"
          control={control}
          defaultValue={0}
          render={({ field }) => (
            <Input
              {...field}
              type="number"
              placeholder="Enter number of sets"
              required
            />
          )}
        />

        <Controller
          name="reps"
          control={control}
          defaultValue={0}
          render={({ field }) => (
            <Input
              {...field}
              type="number"
              placeholder="Enter number of reps"
              required
            />
          )}
        />

        <HStack gap={4}>
          <Button colorScheme="teal" type="submit" isLoading={loading}>
            {loading ? <Spinner size="sm" /> : 'Add Exercise'}
          </Button>
        </HStack>

        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

const ProgramManager = () => {
  const { user, isAuthenticated } = useAuth();
  const [programId, setProgramId] = useState('');

  if (!isAuthenticated) {
    return <Text>Please log in to manage your programs.</Text>;
  }

  return (
    <Box p={5}>
      <CreateTrainingProgram />
      {programId && <CreateWorkout programId={programId} />}
      {programId && <AddExerciseToWorkout workoutId={programId} />}
    </Box>
  );
};

export default ProgramManager;
