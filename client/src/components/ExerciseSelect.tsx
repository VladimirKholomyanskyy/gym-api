import { Exercise } from "@/types/api";
import {
  SelectContent,
  SelectItem,
  SelectRoot,
  SelectTrigger,
  SelectValueText,
} from "./ui/select";
import { createListCollection } from "@chakra-ui/react/collection";

const ExerciseSelect: React.FC<{
  exercises: Exercise[];
  selectedExerciseId: string[];
  contentRef: React.RefObject<HTMLElement>;
  setSelectedExerciseId: (ids: string[]) => void;
}> = ({ exercises, selectedExerciseId, contentRef, setSelectedExerciseId }) => {
  const collection = createListCollection({
    items: exercises.map((exercise) => ({
      label: exercise.Name,
      value: exercise.ID.toString(),
    })),
  });

  return (
    <SelectRoot
      collection={collection}
      value={selectedExerciseId}
      onValueChange={({ value }) => {
        console.log("Selected Value:", value);
        setSelectedExerciseId(value);
      }}
    >
      <SelectTrigger clearable>
        <SelectValueText placeholder="Select Exercise" />
      </SelectTrigger>
      <SelectContent portalRef={contentRef}>
        {exercises.map((exercise) => (
          <SelectItem key={exercise.ID} item={String(exercise.ID)}>
            {exercise.Name}
          </SelectItem>
        ))}
      </SelectContent>
    </SelectRoot>
  );
};

export default ExerciseSelect;
