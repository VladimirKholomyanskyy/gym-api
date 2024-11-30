import React, { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { Box, VStack, Input, Textarea, Text } from '@chakra-ui/react';
import { Button } from './ui/button';
import { useAuth } from 'react-oidc-context';
import { createTrainingProgram } from '../api/trainingPrograms';
import { CreateTrainingProgramRequest } from '../types/api';

const CreateTrainingProgram: React.FC = () => {
  const { user, isAuthenticated } = useAuth();
  const { control, handleSubmit } = useForm<CreateTrainingProgramRequest>();
  const [message, setMessage] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const onSubmit = async (data: CreateTrainingProgramRequest) => {
    if (!isAuthenticated || !user) return;

    setLoading(true);
    setMessage('');

    try {
      await createTrainingProgram(data);
      setMessage('Training Program Created Successfully!');
    } catch (error: any) {
      setMessage(error.response?.data?.message || 'Failed to create training program.');
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
          render={({ field }) => <Input {...field} placeholder="Program Name" required />}
        />
        <Controller
          name="description"
          control={control}
          defaultValue=""
          render={({ field }) => <Textarea {...field} placeholder="Program Description" required />}
        />
        <Button colorScheme="teal" type="submit" loading={loading}>
          Create Program
        </Button>
        {message && <Text>{message}</Text>}
      </VStack>
    </Box>
  );
};

export default CreateTrainingProgram;
