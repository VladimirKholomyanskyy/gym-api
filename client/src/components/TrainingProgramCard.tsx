import { Card, IconButton, Stack } from "@chakra-ui/react";
import { useNavigate } from "react-router";
import { Button } from "./ui/button";
import { FaTrash } from "react-icons/fa";
import {
  DialogActionTrigger,
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";

interface TrainingProgramCardProps {
  id: number; // or number, depending on your program's ID type
  name: string;
  description: string;
  onDelete: (id: number) => void; // or number, depending on your ID type
  onUpdate: (id: number, newName: string, newDescription: string) => void;
}

const TrainingProgramCard = ({
  id,
  name,
  description,
  onDelete,
}: TrainingProgramCardProps) => {
  const navigate = useNavigate();

  const handleNavigate = () => {
    navigate(`/training-programs/${id}/workouts`);
  };

  const handleDelete = () => {
    onDelete(id); // Call the parent function to delete the program
  };

  return (
    <Card.Root size="sm" width="100%">
      <Card.Body>
        <Stack gap="4">
          <Card.Title
            mt="2"
            onClick={handleNavigate}
            cursor="pointer"
            _hover={{ color: "blue.500" }}
          >
            {name}
          </Card.Title>
          {description && <Card.Description>{description}</Card.Description>}
        </Stack>
      </Card.Body>
      <Card.Footer display="flex" justifyContent="space-between">
        <Button colorScheme="blue" onClick={handleNavigate}>
          View Workouts
        </Button>
        <DialogRoot role="alertdialog">
          <DialogTrigger asChild>
            <IconButton colorScheme="red" aria-label="Delete Program">
              <FaTrash />
            </IconButton>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Are you sure?</DialogTitle>
            </DialogHeader>
            <DialogBody>
              <p>
                This action cannot be undone. This will permanently delete your
                account and remove your data from our systems.
              </p>
            </DialogBody>
            <DialogFooter>
              <DialogActionTrigger asChild>
                <Button variant="outline">Cancel</Button>
              </DialogActionTrigger>
              <DialogActionTrigger asChild>
                <Button colorPalette="red" onClick={handleDelete}>
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

export default TrainingProgramCard;
