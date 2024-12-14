import { Card, IconButton, Stack } from "@chakra-ui/react";
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

export interface ExerciseCardProps {
  exercise: string;
  sets: number;
  reps: number;
}

const ExerciseCard = ({ exercise, sets, reps }: ExerciseCardProps) => {
  return (
    <Card.Root>
      <Card.Body>
        <Stack gap={4}>
          <Card.Title>{exercise}</Card.Title>
          <Card.Description>{`Sets: ${sets} Reps: ${reps}`}</Card.Description>
        </Stack>
      </Card.Body>
      <Card.Footer
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <Button colorScheme="blue">View / Edit</Button>
        <IconButton colorScheme="gray" aria-label="Edit name">
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
              <Button colorScheme="red">Delete</Button>
            </DialogFooter>
            <DialogCloseTrigger />
          </DialogContent>
        </DialogRoot>
      </Card.Footer>
    </Card.Root>
  );
};

export default ExerciseCard;
