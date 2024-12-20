import { Box, Card, IconButton, Stack, VStack } from "@chakra-ui/react";
import { Button } from "./ui/button";
import { FaTrash } from "react-icons/fa";
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
import { Exercise } from "@/types/api";
import { useState } from "react";

export interface ExerciseCardProps {
  exerciseId: number;
  exercise: string;
  exercises: Exercise[];
  sets: number;
  reps: number;
  contentRef: React.RefObject<HTMLDivElement>;
  onDelete: () => void;
  onEdit: (exerciseId: number, reps: number, sets: number) => void;
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
    <Card.Root size="sm">
      <Card.Body>
        <Stack gap={4}>
          <Box>
            <Card.Title>{exercise}</Card.Title>
            <Card.Description>{`Sets: ${sets} Reps: ${reps}`}</Card.Description>
          </Box>
        </Stack>
      </Card.Body>
      <Card.Footer
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <DrawerRoot placement="bottom">
          <DrawerBackdrop />
          <DrawerTrigger asChild>
            <Button colorScheme="teal" aria-label="Add Exercise" size="lg">
              Edit Exercise
            </Button>
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
                  onEdit(Number(selectedExerciseId), repsLocal, setsLocal)
                }
              >
                Save
              </Button>
            </DrawerFooter>
          </DrawerContent>
        </DrawerRoot>
        <DialogRoot role="alertdialog">
          <DialogTrigger asChild>
            <IconButton
              colorScheme="red"
              aria-label="Delete Workout"
              onClick={onDelete}
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
