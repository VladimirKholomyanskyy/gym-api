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
import ConfirmationDialog from "./common/ConfirmationDialog";

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
    <Card.Root size="sm" width="100%" background="bg.error" borderRadius="none">
      <Flex align="stretch">
        <Card.Body>
          <Stack gap="4">
            <Box>
              <Card.Title onClick={handleNavigate} cursor="pointer">
                {name}
              </Card.Title>
              {joinedExercises && (
                <Card.Description>{joinedExercises}</Card.Description>
              )}
            </Box>
          </Stack>
        </Card.Body>
        <Card.Footer
          display="flex"
          justifyContent="space-between"
          alignItems="center"
        ></Card.Footer>
        <Box>
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <IconButton height="100%" borderRadius="none">
                <FaEllipsisVertical />
              </IconButton>
            </DrawerTrigger>
            <DrawerContent>
              <DrawerCloseTrigger />
              <DrawerHeader>Add a New Training Program</DrawerHeader>
              <DrawerBody bg="bg.subtle">
                <Stack align="flex-start">
                  <ConfirmationDialog
                    message={
                      "This action cannot be undone. This will permanently delete your workout."
                    }
                    onDelete={handleDelete}
                  />
                  <IconButton
                    colorScheme="red"
                    aria-label="Delete Program"
                    onClick={handleNavigate}
                  >
                    <FaEdit /> Edit
                  </IconButton>
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
