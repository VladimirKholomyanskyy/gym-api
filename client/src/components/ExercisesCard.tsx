import { Box, Card, Flex, IconButton, Stack, VStack } from "@chakra-ui/react";
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
import { NumberInputField, NumberInputRoot } from "./ui/number-input";
import { useState } from "react";
import { Exercise } from "@/api/models";
import { FaEllipsisVertical } from "react-icons/fa6";
import ConfirmationDialog from "./common/ConfirmationDialog";

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
    <Card.Root size="sm" width="100%" background="bg.error" borderRadius="none">
      <Flex align="stretch">
        <Card.Body>
          <Stack gap={4}>
            <Card.Title>{exercise}</Card.Title>
            <Card.Description>{`Sets: ${sets} Reps: ${reps}`}</Card.Description>
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
                    onDelete={onDelete}
                  />
                  <DrawerRoot placement="bottom">
                    <DrawerBackdrop />
                    <DrawerTrigger asChild>
                      <IconButton colorScheme="red">
                        <FaEdit /> Edit
                      </IconButton>
                    </DrawerTrigger>
                    <DrawerContent ref={contentRef}>
                      <DrawerCloseTrigger />
                      <DrawerHeader>Edit Exercise</DrawerHeader>
                      <DrawerBody>
                        <VStack gap={4}>
                          <ExerciseSelect
                            exercises={exercises}
                            defaultExerciseId={exerciseId.toString()}
                            contentRef={contentRef}
                            setSelectedExerciseId={setSelectedExerciseId}
                          />
                          <NumberInputRoot
                            defaultValue={sets.toString()}
                            min={1}
                            onValueChange={(e) => setSetsLocal(e.valueAsNumber)}
                          >
                            <NumberInputField placeholder="Number of Sets" />
                          </NumberInputRoot>

                          <NumberInputRoot
                            defaultValue={reps.toString()}
                            min={1}
                            onValueChange={(e) => setRepsLocal(e.valueAsNumber)}
                          >
                            <NumberInputField placeholder="Number of Reps" />
                          </NumberInputRoot>
                        </VStack>
                      </DrawerBody>
                      <DrawerFooter>
                        <DrawerActionTrigger asChild>
                          <Button variant="outline">Cancel</Button>
                        </DrawerActionTrigger>
                        <Button
                          colorScheme="teal"
                          onClick={() =>
                            onEdit(selectedExerciseId[0], repsLocal, setsLocal)
                          }
                        >
                          Save
                        </Button>
                      </DrawerFooter>
                    </DrawerContent>
                  </DrawerRoot>
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
