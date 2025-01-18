import { Card, IconButton, Stack, Box } from "@chakra-ui/react";
import { Button } from "./ui/button";
import { FaTrash } from "react-icons/fa";
import {
  DialogActionTrigger,
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogRoot,
} from "./ui/dialog";
import { useNavigate } from "react-router";

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
    <Card.Root size="sm">
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
      >
        <Button colorScheme="blue" onClick={handleNavigate}>
          View Exercises
        </Button>
        <DialogRoot role="alertdialog">
          <DialogTrigger asChild>
            <IconButton colorScheme="red" aria-label="Delete Workout">
              <FaTrash />
            </IconButton>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Are you sure?</DialogTitle>
            </DialogHeader>
            <DialogBody>
              <p>
                This action cannot be undone. This will permanently delete this
                workout and its data.
              </p>
            </DialogBody>
            <DialogFooter>
              <DialogActionTrigger asChild>
                <Button variant="outline">Cancel</Button>
              </DialogActionTrigger>
              <DialogActionTrigger asChild>
                <Button colorScheme="red" onClick={handleDelete}>
                  Delete
                </Button>
              </DialogActionTrigger>
            </DialogFooter>
            <DialogCloseTrigger />
          </DialogContent>
        </DialogRoot>
      </Card.Footer>
    </Card.Root>
  );
};

export default WorkoutCard;
