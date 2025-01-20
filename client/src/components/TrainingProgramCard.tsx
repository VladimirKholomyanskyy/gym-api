import { Box, Card, Flex, IconButton, Stack } from "@chakra-ui/react";
import { useNavigate } from "react-router";
import { FaEdit } from "react-icons/fa";
import { FaEllipsisVertical } from "react-icons/fa6";
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
import ConfirmationDialog from "./common/ConfirmationDialog";

interface TrainingProgramCardProps {
  id: string;
  name: string;
  description?: string;
  onDelete: (id: string) => void;
  onUpdate: (id: string, newName: string, newDescription: string) => void;
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
    onDelete(id);
  };

  return (
    <Card.Root size="sm" width="100%" background="bg.error" borderRadius="none">
      <Flex align="stretch">
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
        <Card.Footer />
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
                      "This action cannot be undone. This will permanently delete your training program."
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

export default TrainingProgramCard;
