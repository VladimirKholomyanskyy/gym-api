import {
  Box,
  Flex,
  Input,
  Spinner,
  VStack,
  useDisclosure,
  Heading,
} from "@chakra-ui/react";
import { toaster } from "@/components/ui/toaster";
import { Button } from "./ui/button";
import { useEffect, useRef, useState } from "react";
import { Field } from "./ui/field";
import TrainingProgramCard from "./TrainingProgramCard";
import {
  DrawerActionTrigger,
  DrawerBackdrop,
  DrawerBody,
  DrawerCloseTrigger,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerRoot,
  DrawerTrigger,
} from "./ui/drawer";
import { TrainingProgram, TrainingProgramsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";

const TrainingProgramsPage: React.FC = () => {
  const [programs, setPrograms] = useState<TrainingProgram[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [newProgram, setNewProgram] = useState({ name: "", description: "" });
  const ref = useRef<HTMLInputElement>(null);
  const { onClose } = useDisclosure();
  const trainingProgramApi = new TrainingProgramsApi(apiConfig);

  useEffect(() => {
    const loadPrograms = async () => {
      console.log("load training programs");
      setLoading(true); // Ensure loading is set to true when the function starts
      try {
        const response = await trainingProgramApi.listTrainingPrograms();
        setPrograms(response.data);
      } catch (error) {
        console.error("Error loading training programs:", error);
        toaster.create({
          title: "Failed to load training programs.",
          description: "Please try again later.",
          type: "error",
          duration: 5000,
        });
      } finally {
        setLoading(false); // Set loading to false regardless of success or failure
      }
    };

    loadPrograms();
  }, []);

  // Add a program
  const handleAddProgram = async () => {
    if (!newProgram.name.trim()) {
      toaster.create({
        title: "Program name is required.",
        type: "error",
        duration: 3000,
      });
      return;
    }
    try {
      const created = await trainingProgramApi.createTrainingProgram(
        newProgram
      );
      setPrograms((prev) => [...prev, created.data]);
      setNewProgram({ name: "", description: "" });
      toaster.create({
        title: "Training program created.",
        type: "success",
        duration: 3000,
      });
      onClose();
    } catch (error) {
      toaster.create({
        title: "Failed to create training program.",
        description: "Please try again later.",
        type: "error",
        duration: 5000,
      });
    }
  };

  const handleDeleteProgram = async (id: string) => {
    try {
      await trainingProgramApi.deleteTrainingProgram(id);
      setPrograms((prev) => prev.filter((program) => program.id !== id));
      toaster.create({
        title: "Training program deleted.",
        type: "success",
        duration: 3000,
      });
    } catch (error) {
      toaster.create({
        title: "Failed to delete training program.",
        description: "Please try again later.",
        type: "error",
        duration: 5000,
      });
    }
  };

  const handleUpdateProgram = async (
    id: string,
    name: string,
    description: string
  ) => {
    const currentProgram = programs.find((program) => program.id === id);

    if (
      currentProgram &&
      currentProgram.name === name.trim() &&
      currentProgram.description === description.trim()
    ) {
      return;
    }
    try {
      await trainingProgramApi.updateTrainingProgram(id, { name, description });
      setPrograms((prev) => prev.filter((program) => program.id !== id));
      toaster.create({
        title: "Training program updated.",
        type: "success",
        duration: 3000,
      });
    } catch (error) {
      toaster.create({
        title: "Failed to update training program.",
        description: "Please try again later.",
        type: "error",
        duration: 5000,
      });
    }
  };

  return (
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        Training Programs
      </Heading>
      {loading ? (
        <Flex justifyContent="center" alignItems="center" height="50vh">
          <Spinner size="xl" color="magenta.400" />
        </Flex>
      ) : (
        <VStack gap={6} align="stretch" width="100%" p={4}>
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
          <DrawerRoot placement="bottom">
            <DrawerBackdrop />
            <DrawerTrigger asChild>
              <Button
                background="linear-gradient(90deg, rgba(255,0,255,1) 0%, rgba(0,255,255,1) 100%)"
                color="white"
                _hover={{
                  filter: "brightness(1.2)",
                  boxShadow: "0 0 10px rgba(255, 0, 255, 0.8)",
                }}
                size="lg"
              >
                + Add Program
              </Button>
            </DrawerTrigger>
            <DrawerContent ref={ref} background="blackAlpha.900">
              <DrawerCloseTrigger />
              <DrawerHeader color="magenta.400">
                Add a New Training Program
              </DrawerHeader>
              <DrawerBody>
                <VStack gap={4}>
                  <Field label="Name">
                    <Input
                      ref={ref}
                      placeholder="Training program name"
                      value={newProgram.name}
                      onChange={(e) =>
                        setNewProgram({ ...newProgram, name: e.target.value })
                      }
                    />
                  </Field>
                  <Field label="Description">
                    <Input
                      placeholder="Description"
                      value={newProgram.description}
                      onChange={(e) =>
                        setNewProgram({
                          ...newProgram,
                          description: e.target.value,
                        })
                      }
                    />
                  </Field>
                </VStack>
              </DrawerBody>
              <DrawerFooter>
                <DrawerActionTrigger asChild>
                  <Button variant="outline">Cancel</Button>
                </DrawerActionTrigger>
                <DrawerActionTrigger asChild>
                  <Button
                    background="magenta.400"
                    color="black"
                    onClick={handleAddProgram}
                  >
                    Save
                  </Button>
                </DrawerActionTrigger>
              </DrawerFooter>
            </DrawerContent>
          </DrawerRoot>
        </VStack>
      )}
    </Box>
  );
};

export default TrainingProgramsPage;
