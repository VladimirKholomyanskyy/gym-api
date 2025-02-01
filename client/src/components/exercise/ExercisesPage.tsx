import { Exercise, ExercisesApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";
import { Box, For, Heading, VStack, Card } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

const ExercisesPage = () => {
  const exerciseApi = new ExercisesApi(apiConfig);
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const navigate = useNavigate();
  useEffect(() => {
    const fetchExercises = async () => {
      try {
        const response = await exerciseApi.listExercises();
        setExercises(response.data);
      } catch (error) {
        console.log(error);
      }
    };
    fetchExercises();
  }, []);

  const handleNavigate = (id: string) => {
    navigate(`/exercises/${id}`);
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
        Exercises
      </Heading>

      <VStack gap={6} align="stretch" width="100%" p={4}>
        <For each={exercises}>
          {(exercise, index) => (
            <Card.Root
              key={index}
              size="sm"
              width="100%"
              background="blackAlpha.800"
              borderRadius="md"
              boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
              p={2}
              _hover={{ boxShadow: "0 0 20px rgba(0, 255, 255, 1)" }}
            >
              <Card.Body color="fg.muted">
                <Card.Title
                  fontSize="xl"
                  fontWeight="bold"
                  color="neon.400"
                  textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
                  cursor="pointer"
                  onClick={() => handleNavigate(exercise.id)}
                  _hover={{ color: "neon.300" }}
                >
                  {exercise.name}
                </Card.Title>
                <Card.Description fontSize="md" color="gray.300">
                  {exercise.description}
                </Card.Description>
              </Card.Body>
            </Card.Root>
          )}
        </For>
      </VStack>
    </Box>
  );
};

export default ExercisesPage;
