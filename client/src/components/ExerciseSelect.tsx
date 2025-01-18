import { Exercise } from "@/api/models";
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
  defaultExerciseId?: string;
  contentRef: React.RefObject<HTMLDivElement>;
  setSelectedExerciseId: (ids: string[]) => void;
}> = ({ exercises, defaultExerciseId, contentRef, setSelectedExerciseId }) => {
  const collection = createListCollection({
    items: exercises,
    itemToString: (item) => item.name,
    itemToValue: (item) => item.id.toString(),
  });

  console.log(defaultExerciseId);
  return (
    <SelectRoot
      collection={collection}
      defaultValue={defaultExerciseId ? [defaultExerciseId] : []}
      onValueChange={({ value }) => {
        setSelectedExerciseId(value);
      }}
    >
      <SelectTrigger clearable>
        <SelectValueText placeholder="Select Exercise" />
      </SelectTrigger>
      <SelectContent portalRef={contentRef}>
        {exercises.map((exercise) => (
          <SelectItem key={exercise.id} item={exercise}>
            {exercise.name}
          </SelectItem>
        ))}
      </SelectContent>
    </SelectRoot>
  );
};

export default ExerciseSelect;
