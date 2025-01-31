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
import DeleteDialog from "./common/DeleteDialog";

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
          <Stack gap={4}>
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
            {description && (
              <Card.Description fontSize="md" color="gray.300">
                {description}
              </Card.Description>
            )}
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
              <DrawerHeader color="neon.500">
                Manage Training Program
              </DrawerHeader>
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
                      "This action cannot be undone. This will permanently delete your training program."
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

export default TrainingProgramCard;
