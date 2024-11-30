import React, { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { Box, VStack, Input, Text } from '@chakra-ui/react';
import { Button } from './ui/button';
import { createWorkout } from '../api/workouts';
import { CreateWorkoutRequest } from '../types/api';

interface Props {
  programId: string;
}

const CreateWorkout: React.FC<Props> = ({ programId }) => {
  const { control, handleSubmit } = useForm<CreateWorkoutRequest>();
  const [message, setMessage] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const onSubmit = async (data: Omit<CreateWorkoutRequest, 'programId'>) => {
    setLoading(true);
    setMessage('');

    try {
      await createWorkout(programId, { ...data, programId });
      setMessage('Workout Created Successfully!');
    } catch (error: any) {
      setMessage(error.response?.data?.message || 'Failed to create workout.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box p={5} borderWidth={1} borderRadius="lg">
      <VStack as="form" gap={4} onSubmit={handleSubmit(onSubmit)}>
        <Controller
          name="name"
          control={control}
          defaultValue=""
          render={({ field }) => <Input {...field} placeholder="Workout Name" required />}
        />
        <Button colorScheme="teal" type="submit" loading={loading}>
          Create Workout
        </Button>
        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

export default CreateWorkout;
