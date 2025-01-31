import { Exercise } from "@/api/models";
import { VStack, Box, Text, List } from "@chakra-ui/react";

interface DescriptionTabProps {
  exercise: Exercise;
}

const DescriptionTab = ({ exercise }: DescriptionTabProps) => {
  return (
    <VStack
      gap={4}
      align="start"
      p={4}
      background="blackAlpha.800"
      borderRadius="lg"
      boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
    >
      <Box>
        <Text
          fontSize="xl"
          fontWeight="bold"
          color="neon.400"
          textShadow="0 0 10px cyan"
          textAlign="left"
        >
          Description
        </Text>
        <Text mt={2} color="gray.300" textAlign="left">
          {exercise.description ||
            "No description available for this exercise."}
        </Text>
      </Box>

      <Box>
        <Text
          fontSize="xl"
          fontWeight="bold"
          color="neon.400"
          textShadow="0 0 10px cyan"
          textAlign="left"
        >
          Primary Muscle
        </Text>
        <Text mt={2} color="gray.300" textAlign="left">
          {exercise.primaryMuscle || "Not specified"}
        </Text>
      </Box>

      <Box>
        <Text
          fontSize="xl"
          fontWeight="bold"
          color="neon.400"
          textShadow="0 0 10px cyan"
          textAlign="left"
        >
          Secondary Muscles
        </Text>
        {exercise.secondaryMuscle && exercise.secondaryMuscle.length > 0 ? (
          <List.Root mt={2} gap={1} color="gray.300" textAlign="left">
            {exercise.secondaryMuscle.map((muscle, index) => (
              <List.Item key={index}>{muscle}</List.Item>
            ))}
          </List.Root>
        ) : (
          <Text mt={2} color="gray.300" textAlign="left">
            None
          </Text>
        )}
      </Box>

      <Box>
        <Text
          fontSize="xl"
          fontWeight="bold"
          color="neon.400"
          textShadow="0 0 10px cyan"
          textAlign="left"
        >
          Equipment
        </Text>
        <Text mt={2} color="gray.300" textAlign="left">
          {exercise.equipment || "No equipment required"}
        </Text>
      </Box>
    </VStack>
  );
};

export default DescriptionTab;
