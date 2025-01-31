import { Card, IconButton, Stack, Box, Flex } from "@chakra-ui/react";
import { FaEdit } from "react-icons/fa";
import { useNavigate } from "react-router";
import {
  DrawerBackdrop,
  DrawerBody,
  DrawerCloseTrigger,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerRoot,
  DrawerTrigger,
} from "./ui/drawer";
import { FaEllipsisVertical } from "react-icons/fa6";
import DeleteDialog from "./common/DeleteDialog";

export interface WorkoutCardProps {
  workoutId: string;
  programId: string;
  name: string;
  exercises: string[];
  onDelete: (programId: string, workoutId: string) => void;
}

const WorkoutCard = ({
  workoutId,
  programId,
  name,
  exercises,
  onDelete,
}: WorkoutCardProps) => {
  const navigate = useNavigate();

  const handleNavigate = () => {
    navigate(`/training-programs/${programId}/workouts/${workoutId}`);
  };

  const handleDelete = () => {
    onDelete(programId, workoutId); // Call the parent function to delete the workout
  };

  const joinedExercises = exercises.join(", ");

  return (
    <Card.Root
      size="sm"
      width="100%"
      background="blackAlpha.800"
      borderRadius="md"
      boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
      p={2}
      _hover={{ boxShadow: "0 0 20px rgba(0, 255, 255, 1)" }}
    >
      <Flex align="stretch" justify="space-between">
        <Card.Body>
          <Stack gap="4">
            <Box>
              <Card.Title
                fontSize="xl"
                fontWeight="bold"
                color="neon.400"
                textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
                cursor="pointer"
                onClick={handleNavigate}
                _hover={{ color: "neon.300" }}
              >
                {name}
              </Card.Title>
              {joinedExercises && (
                <Card.Description fontSize="md" color="gray.300">
                  {joinedExercises}
                </Card.Description>
              )}
            </Box>
          </Stack>
        </Card.Body>
        <Card.Footer />
        <Box>
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <IconButton
                height="100%"
                borderRadius="md"
                background="blackAlpha.600"
                color="white"
                _hover={{ background: "blackAlpha.400" }}
                aria-label="Options"
              >
                <FaEllipsisVertical />
              </IconButton>
            </DrawerTrigger>
            <DrawerContent>
              <DrawerCloseTrigger />
              <DrawerHeader color="neon.500">Add a New Workout</DrawerHeader>
              <DrawerBody bg="bg.subtle">
                <Stack align="flex-start">
                  <IconButton
                    background="transparent"
                    color="neon.400"
                    _hover={{ color: "neon.300" }}
                    aria-label="Edit Program"
                    onClick={handleNavigate}
                  >
                    <FaEdit /> Edit
                  </IconButton>
                  <DeleteDialog
                    message={
                      "This action cannot be undone. This will permanently delete your workout."
                    }
                    onDelete={handleDelete}
                  />
                </Stack>
              </DrawerBody>
              <DrawerFooter></DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
        </Box>
      </Flex>
    </Card.Root>
  );
};

export default WorkoutCard;
