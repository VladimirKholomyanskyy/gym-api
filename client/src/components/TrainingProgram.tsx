import {
  createTrainingProgram,
  deleteTrainingProgram,
  getAllTrainingPrograms,
  updateTrainingProgram,
} from "@/api/trainingPrograms";
import { TrainingProgram } from "@/types/api";
import {
  Box,
  Text,
  Flex,
  Input,
  Spinner,
  VStack,
  Stack,
} from "@chakra-ui/react";
import { toaster } from "@/components/ui/toaster";
import { Button } from "./ui/button";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "react-oidc-context";
import { setupAxiosInterceptors } from "@/api/apiClient";
import {
  DialogActionTrigger,
  DialogBody,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Field } from "./ui/field";
import TrainingProgramCard from "./TrainingProgramCard";

const TrainingPrograms: React.FC = () => {
  const auth = useAuth();
  const [programs, setPrograms] = useState<TrainingProgram[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [newProgram, setNewProgram] = useState({ name: "", description: "" });
  const ref = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const loadPrograms = async () => {
      console.log("load training programs");
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
  }, [auth.user]);

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
  const handleDeleteProgram = async (id: number) => {
    try {
      await deleteTrainingProgram(id);
      setPrograms((prev) => prev.filter((program) => program.id !== id));
    } catch (error) {
      console.log("Failed to delete program");
    }
  };

  const handleUpdateProgram = async (
    id: number,
    name: string,
    description: string
  ) => {
    const currentProgram = programs.find((program) => program.id === id);

    // Check if the name or description is actually changed
    if (
      currentProgram &&
      currentProgram.name === name.trim() &&
      currentProgram.description === description.trim()
    ) {
      console.log("No changes detected. Skipping update.");
      return; // Exit the function without performing an update
    }
    try {
      await updateTrainingProgram(id, { name, description });
      setPrograms((prev) => prev.filter((program) => program.id !== id));
    } catch (error) {
      console.log("Failed to update program");
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
          <VStack gap={4} mb={6} align="stretch">
            {programs.map((program) => (
              <TrainingProgramCard
                key={program.id}
                id={program.id}
                name={program.name}
                description={program.description}
                onDelete={handleDeleteProgram}
                onUpdate={handleUpdateProgram}
              />
            ))}
            <DialogRoot initialFocusEl={() => ref.current}>
              <DialogTrigger asChild>
                <Button
                  bottom="16px"
                  right="16px"
                  position="fixed"
                  variant="outline"
                  background="green"
                >
                  New
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>New Training Program</DialogTitle>
                </DialogHeader>
                <DialogBody pb="4">
                  <Stack gap="4">
                    <Field label="Name">
                      <Input
                        ref={ref}
                        placeholder="Training program name"
                        value={newProgram.name}
                        onChange={(e) =>
                          setNewProgram((prev) => ({
                            ...prev,
                            name: e.target.value,
                          }))
                        }
                      />
                    </Field>
                    <Field label="Description">
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
                    </Field>
                  </Stack>
                </DialogBody>
                <DialogFooter>
                  <DialogActionTrigger asChild>
                    <Button variant="outline">Cancel</Button>
                  </DialogActionTrigger>
                  <DialogActionTrigger asChild>
                    <Button onClick={handleAddProgram}>Save</Button>
                  </DialogActionTrigger>
                </DialogFooter>
              </DialogContent>
            </DialogRoot>
          </VStack>
        </>
      )}
    </Box>
  );
};

export default TrainingPrograms;
