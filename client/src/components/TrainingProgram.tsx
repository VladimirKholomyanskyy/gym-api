import {
  createTrainingProgram,
  deleteTrainingProgram,
  getAllTrainingPrograms,
} from "@/api/trainingPrograms";
import { TrainingProgram } from "@/types/api";
import {
  Box,
  Text,
  Flex,
  Input,
  List,
  ListItem,
  Spinner,
  VStack,
  IconButton,
} from "@chakra-ui/react";
import { Toaster, toaster } from "@/components/ui/toaster";
import { Button } from "./ui/button";
import { useEffect, useState } from "react";
import { useAuth } from "react-oidc-context";
import { setupAxiosInterceptors } from "@/api/apiClient";
import { FaTrash } from "react-icons/fa";

const TrainingPrograms: React.FC = () => {
  const auth = useAuth();
  const [programs, setPrograms] = useState<TrainingProgram[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [newProgram, setNewProgram] = useState({ name: "", description: "" });

  useEffect(() => {
    const loadPrograms = async () => {
      try {
        setupAxiosInterceptors(() => auth.user?.access_token || null);
        const data = await getAllTrainingPrograms();
        setPrograms(data);
      } catch (error) {
        toaster.create({ title: "Failed to load training programs" });
      } finally {
        setLoading(false);
      }
    };
    loadPrograms();
  }, [auth]);

  // Add a program
  const handleAddProgram = async () => {
    if (!newProgram.name.trim()) {
      return;
    }

    try {
      const created = await createTrainingProgram(newProgram);
      setPrograms((prev) => [...prev, created]);
      setNewProgram({ name: "", description: "" });
    } catch (error) {
      console.log("Failed to add program");
    }
  };

  // Delete a program
  const handleDeleteProgram = async (id: string) => {
    try {
      await deleteTrainingProgram(id);
      setPrograms((prev) => prev.filter((program) => program.id !== id));
    } catch (error) {
      console.log("Failed to delete program");
    }
  };
  return (
    <Box p={4} maxW="sm" mx="auto">
      <Text fontSize="xl" fontWeight="bold" textAlign="center" mb={4}>
        Training Programs
      </Text>

      {loading ? (
        <Flex justifyContent="center" alignItems="center" height="100vh">
          <Spinner />
        </Flex>
      ) : (
        <>
          <VStack gap={4} mb={6}>
            <Input
              placeholder="Program Name"
              value={newProgram.name}
              onChange={(e) =>
                setNewProgram((prev) => ({ ...prev, name: e.target.value }))
              }
            />
            <Input
              placeholder="Description"
              value={newProgram.description}
              onChange={(e) =>
                setNewProgram((prev) => ({
                  ...prev,
                  description: e.target.value,
                }))
              }
            />
            <Button colorScheme="teal" width="full" onClick={handleAddProgram}>
              Add Program
            </Button>
          </VStack>

          <List.Root gap={3}>
            {programs.map((program) => (
              <ListItem
                key={program.id}
                p={3}
                border="1px"
                borderColor="gray.200"
                borderRadius="md"
                display="flex"
                justifyContent="space-between"
                alignItems="center"
              >
                <Box>
                  <Text fontWeight="bold">{program.name}</Text>
                  <Text fontSize="sm" color="gray.600">
                    {program.description}
                  </Text>
                </Box>
                <IconButton
                  aria-label="Delete"
                  size="sm"
                  colorScheme="red"
                  onClick={() => handleDeleteProgram(program.id)}
                >
                  <FaTrash />
                </IconButton>
              </ListItem>
            ))}
          </List.Root>
        </>
      )}
    </Box>
  );
};

export default TrainingPrograms;
