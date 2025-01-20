import { Group, Input, Table } from "@chakra-ui/react";
import { Button } from "./ui/button";
import { useEffect, useState } from "react";

export interface ExerciseLogItem {
  id: number;
  prevReps: number;
  prevWeight: number;
  reqReps: number;
  currentReps?: number;
  currentWeight?: number;
}
interface ExerciseLogProps {
  items: ExerciseLogItem[];
  onLog: (setNumber: number, repsCompleted: number, weight: number) => void;
}

const ExerciseLog = ({ items, onLog }: ExerciseLogProps) => {
  console.log("Items length:", items.length);

  const [inputValues, setInputValues] = useState(() =>
    items.map(() => ({ reps: "", weight: "" }))
  );

  useEffect(() => {
    const updatedValues = items.map(() => ({ reps: "", weight: "" }));
    setInputValues(updatedValues);
    console.log("Updated inputValues to match items:", updatedValues);
  }, [items]);

  console.log("Current inputValues:", inputValues);
  // Handler for input changes
  const handleInputChange = (
    index: number,
    field: "reps" | "weight",
    value: string
  ) => {
    console.log(`${field}=${value}`);
    setInputValues((prev) =>
      prev.map((input, i) =>
        i === index ? { ...input, [field]: value } : input
      )
    );
  };

  // Handler for log button click
  const handleLog = (index: number) => {
    const { reps, weight } = inputValues[index];
    const repsCompleted = parseInt(reps, 10) || 0;
    const weightValue = parseFloat(weight) || 0;

    // Call the onLog function with the set number, reps, and weight
    onLog(index + 1, repsCompleted, weightValue);
  };
  return (
    <Table.Root size="sm">
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeader>Set</Table.ColumnHeader>
          <Table.ColumnHeader>Prev</Table.ColumnHeader>
          <Table.ColumnHeader>Rep</Table.ColumnHeader>
          <Table.ColumnHeader>kg</Table.ColumnHeader>
          <Table.ColumnHeader textAlign="end">log</Table.ColumnHeader>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {items.map((item, index) => (
          <Table.Row key={item.id} bg="bg.subtle">
            <Table.Cell>{item.id}</Table.Cell>
            <Table.Cell>{item.prevReps + "x" + item.prevWeight}</Table.Cell>
            <Table.Cell>
              <Input
                placeholder={
                  item.currentReps
                    ? String(item.currentReps)
                    : String(item.reqReps)
                }
                size="xs"
                value={inputValues[index].reps}
                onChange={(e) =>
                  handleInputChange(index, "reps", e.target.value)
                }
              />
            </Table.Cell>
            <Table.Cell>
              <Group attached>
                <Input
                  placeholder={
                    item.currentWeight ? String(item.currentWeight) : ""
                  }
                  size="xs"
                  value={inputValues[index].weight}
                  onChange={(e) =>
                    handleInputChange(index, "weight", e.target.value)
                  }
                />
              </Group>
            </Table.Cell>
            <Table.Cell textAlign="end">
              <Button onClick={() => handleLog(index)}>Log</Button>
            </Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
};

export default ExerciseLog;
