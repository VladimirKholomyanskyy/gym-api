import {
  Box,
  Card,
  Flex,
  IconButton,
  NumberInputRoot,
  Stack,
  VStack,
} from "@chakra-ui/react";
import { Button } from "./ui/button";
import { FaEdit } from "react-icons/fa";
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
import ExerciseSelect from "./ExerciseSelect";
import { NumberInputField } from "./ui/number-input";
import { useState } from "react";
import { Exercise } from "@/api/models";
import { FaEllipsisVertical } from "react-icons/fa6";
import DeleteDialog from "./common/DeleteDialog";
import { Field } from "./ui/field";

export interface ExerciseCardProps {
  exerciseId: string;
  exercise: string;
  exercises: Exercise[];
  sets: number;
  reps: number;
  contentRef: React.RefObject<HTMLDivElement>;
  onDelete: () => void;
  onEdit: (exerciseId: string, reps: number, sets: number) => void;
}

const ExerciseCard = ({
  exerciseId,
  exercise,
  exercises,
  sets,
  reps,
  contentRef,
  onDelete,
  onEdit,
}: ExerciseCardProps) => {
  const [setsLocal, setSetsLocal] = useState<number>(sets);
  const [repsLocal, setRepsLocal] = useState<number>(reps);
  const [selectedExerciseId, setSelectedExerciseId] = useState<string[]>([]);

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
              _hover={{ color: "neon.300" }}
            >
              {exercise}
            </Card.Title>
            <Card.Description
              fontSize="md"
              color="gray.300"
            >{`Sets: ${sets} Reps: ${reps}`}</Card.Description>
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
              <DrawerHeader color="neon.500">Add a New Exercise</DrawerHeader>
              <DrawerBody bg="bg.subtle">
                <Stack align="flex-start">
                  <DrawerRoot placement="bottom">
                    <DrawerBackdrop />
                    <DrawerTrigger asChild>
                      <IconButton
                        background="transparent"
                        color="neon.400"
                        _hover={{ color: "neon.300" }}
                        aria-label="Edit Program"
                      >
                        <FaEdit /> Edit
                      </IconButton>
                    </DrawerTrigger>
                    <DrawerContent ref={contentRef}>
                      <DrawerCloseTrigger />
                      <DrawerHeader color="magenta.400">
                        Edit Exercise
                      </DrawerHeader>
                      <DrawerBody>
                        <VStack gap={4}>
                          <ExerciseSelect
                            exercises={exercises}
                            defaultExerciseId={exerciseId.toString()}
                            contentRef={contentRef}
                            setSelectedExerciseId={setSelectedExerciseId}
                          />
                          <Flex gap="4">
                            <Field label="Sets">
                              <NumberInputRoot
                                defaultValue={sets.toString()}
                                min={1}
                                onValueChange={(e) =>
                                  setSetsLocal(e.valueAsNumber)
                                }
                              >
                                <NumberInputField placeholder="Number of Sets" />
                              </NumberInputRoot>
                            </Field>
                            <Field label="Reps">
                              <NumberInputRoot
                                defaultValue={reps.toString()}
                                min={1}
                                onValueChange={(e) =>
                                  setRepsLocal(e.valueAsNumber)
                                }
                              >
                                <NumberInputField placeholder="Number of Reps" />
                              </NumberInputRoot>
                            </Field>
                          </Flex>
                        </VStack>
                      </DrawerBody>
                      <DrawerFooter>
                        <DrawerActionTrigger asChild>
                          <Button
                            variant="outline"
                            borderColor="neon.400"
                            color="neon.400"
                            _hover={{
                              borderColor: "neon.300",
                              color: "neon.300",
                            }}
                          >
                            Cancel
                          </Button>
                        </DrawerActionTrigger>
                        <Button
                          background="neon.400"
                          color="black"
                          onClick={() =>
                            onEdit(selectedExerciseId[0], repsLocal, setsLocal)
                          }
                        >
                          Save
                        </Button>
                      </DrawerFooter>
                    </DrawerContent>
                  </DrawerRoot>
                  <DeleteDialog
                    message={
                      "This action cannot be undone. This will permanently remove the exercise from the workout."
                    }
                    onDelete={onDelete}
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

export default ExerciseCard;
