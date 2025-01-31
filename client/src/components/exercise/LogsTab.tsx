import { ExerciseLogsApi, LogExerciseResponse } from "@/api";
import { apiConfig } from "@/api/apiConfig";
import { formatDateTime } from "@/utils/dateUtils";
import { Box, VStack, Text } from "@chakra-ui/react";
import { useEffect, useState } from "react";

interface LogsTabProps {
  exerciseId: string;
}

const LogsTab = ({ exerciseId }: LogsTabProps) => {
  const exerciseLogsApi = new ExerciseLogsApi(apiConfig);
  const [exerciseLogs, setExerciseLogs] = useState<LogExerciseResponse[]>([]);

  useEffect(() => {
    const fetchExercise = async () => {
      try {
        const response = await exerciseLogsApi.listExerciseLogs(
          undefined,
          exerciseId
        );
        setExerciseLogs(response.data);
      } catch (error) {
        console.log("Error", error);
      }
    };
    fetchExercise();
  }, [exerciseId]);

  return (
    <VStack gap={4} align="stretch" p={4}>
      {exerciseLogs.length === 0 ? (
        <Text color="gray.400" fontStyle="italic" textAlign="center">
          No logs available for this exercise.
        </Text>
      ) : (
        exerciseLogs.map((item, index) => (
          <Box
            key={index}
            p={4}
            border="1px solid"
            borderColor="neon.400"
            background="blackAlpha.800"
            boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
            borderRadius="md"
            transition="transform 0.2s ease-in-out"
            _hover={{
              transform: "scale(1.02)",
              boxShadow: "0 0 15px rgba(0, 255, 255, 1)",
            }}
          >
            <Text
              fontWeight="bold"
              color="neon.300"
              textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
            >
              {formatDateTime(item.loggedAt)}
            </Text>
            <Text color="gray.300" fontSize="lg">
              {item.repsCompleted}x{item.weightUsed}kg
            </Text>
          </Box>
        ))
      )}
    </VStack>
  );
};

export default LogsTab;
