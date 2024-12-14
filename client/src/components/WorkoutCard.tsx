import {
  Card,
  IconButton,
  Stack,
  Box,
  useDisclosure,
  Editable,
} from "@chakra-ui/react";
import { Button } from "./ui/button";
import { FaEdit, FaTrash } from "react-icons/fa";
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
import { useState } from "react";
import { useNavigate } from "react-router";
import { LuCheck, LuPencilLine, LuX } from "react-icons/lu";

export interface WorkoutCardProps {
  workoutId: number;
  programId: number;
  name: string;
  exercises: string[];
  onDelete: (programId: number, workoutId: number) => void;
  onUpdate: (id: number, newName: string) => void;
}

const WorkoutCard = ({
  workoutId,
  programId,
  name,
  exercises,
  onDelete,
  onUpdate,
}: WorkoutCardProps) => {
  const [isEditing, setIsEditing] = useState(false);
  const [newName, setNewName] = useState(name);
  const { onOpen, onClose } = useDisclosure();
  const navigate = useNavigate();

  const handleNavigate = () => {
    navigate(`/training-programs/${programId}/workouts/${workoutId}`);
  };

  const handleSave = () => {
    onUpdate(workoutId, newName); // Update the workout's name
    setIsEditing(false);
  };

  const handleDelete = () => {
    onDelete(programId, workoutId); // Call the parent function to delete the workout
    onClose();
  };

  const joinedExercises = exercises.join(", ");

  return (
    <Card.Root>
      <Card.Body>
        <Stack gap="4">
          <Box>
            <Editable.Root
              defaultValue={name}
              onValueChange={(e) => setNewName(e.value)}
            >
              <Editable.Preview />
              <Editable.Input />
              <Editable.Control>
                <Editable.EditTrigger asChild>
                  <IconButton variant="ghost" size="xs">
                    <LuPencilLine />
                  </IconButton>
                </Editable.EditTrigger>
                <Editable.CancelTrigger asChild>
                  <IconButton variant="outline" size="xs">
                    <LuX />
                  </IconButton>
                </Editable.CancelTrigger>
                <Editable.SubmitTrigger asChild>
                  <IconButton variant="outline" size="xs" onClick={handleSave}>
                    <LuCheck />
                  </IconButton>
                </Editable.SubmitTrigger>
              </Editable.Control>
            </Editable.Root>

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
          View / Edit Workout
        </Button>
        <IconButton
          colorScheme="gray"
          aria-label="Edit name"
          onClick={() => setIsEditing(true)}
        >
          <FaEdit />
        </IconButton>
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
              <Button colorScheme="red" onClick={handleDelete}>
                Delete
              </Button>
            </DialogFooter>
            <DialogCloseTrigger />
          </DialogContent>
        </DialogRoot>
      </Card.Footer>
    </Card.Root>
  );
};

export default WorkoutCard;
