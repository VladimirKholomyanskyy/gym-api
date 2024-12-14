import {
  Card,
  IconButton,
  Input,
  Stack,
  useDisclosure,
} from "@chakra-ui/react";
import { useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "./ui/button";
import { FaCheck, FaEdit, FaTrash } from "react-icons/fa";
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
  onUpdate,
}: TrainingProgramCardProps) => {
  const [isEditing, setIsEditing] = useState(false);
  const [newName, setNewName] = useState(name);
  const [newDescription, setNewDescription] = useState(description);
  const { onOpen, onClose } = useDisclosure();
  const navigate = useNavigate();

  const handleNavigate = () => {
    navigate(`/training-programs/${id}/workouts`);
  };

  const handleSave = () => {
    onUpdate(id, newName, newDescription); // Update the program's name and description
    setIsEditing(false);
  };

  const handleDelete = () => {
    onDelete(id); // Call the parent function to delete the program
    onClose();
  };

  return (
    <Card.Root>
      <Card.Body>
        <Stack gap="4">
          {isEditing ? (
            <>
              <Input
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="Edit name"
              />
              <Input
                value={newDescription}
                onChange={(e) => setNewDescription(e.target.value)}
                placeholder="Edit description"
              />
            </>
          ) : (
            <>
              <Card.Title
                mt="2"
                onClick={handleNavigate}
                cursor="pointer"
                _hover={{ color: "blue.500" }}
              >
                {name}
              </Card.Title>
              {description && (
                <Card.Description>{description}</Card.Description>
              )}
            </>
          )}
        </Stack>
      </Card.Body>
      <Card.Footer display="flex" justifyContent="space-between">
        {isEditing ? (
          <>
            <IconButton
              colorScheme="green"
              onClick={handleSave}
              aria-label="Save changes"
            >
              <FaCheck />
            </IconButton>
            <IconButton
              colorScheme="red"
              onClick={() => setIsEditing(false)}
              aria-label="Cancel editing"
            >
              <FaTrash />
            </IconButton>
          </>
        ) : (
          <>
            <Button colorScheme="blue" onClick={handleNavigate}>
              View / Edit Program
            </Button>
            <IconButton
              colorScheme="gray"
              onClick={() => setIsEditing(true)}
              aria-label="Edit name and description"
            >
              <FaEdit />
            </IconButton>
            <DialogRoot role="alertdialog">
              <DialogTrigger asChild>
                <IconButton
                  colorScheme="red"
                  onClick={onOpen}
                  aria-label="Delete Program"
                >
                  <FaTrash />
                </IconButton>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Are you sure?</DialogTitle>
                </DialogHeader>
                <DialogBody>
                  <p>
                    This action cannot be undone. This will permanently delete
                    your account and remove your data from our systems.
                  </p>
                </DialogBody>
                <DialogFooter>
                  <DialogActionTrigger asChild>
                    <Button variant="outline">Cancel</Button>
                  </DialogActionTrigger>
                  <Button colorPalette="red" onClick={handleDelete}>
                    Delete
                  </Button>
                </DialogFooter>
                <DialogCloseTrigger />
              </DialogContent>
            </DialogRoot>
          </>
        )}
      </Card.Footer>
    </Card.Root>
  );
};

export default TrainingProgramCard;
