import { Card, Table } from "@chakra-ui/react";

export interface CompletedExerciseLogItem {
  setNumber: number;
  repsCompleted: number;
  weightCompleted: number;
}
interface CompletedExerciseLogProps {
  exerciseName: string;
  items: CompletedExerciseLogItem[];
}

const CompletedExerciseLog = ({
  exerciseName,
  items,
}: CompletedExerciseLogProps) => {
  return (
    <Card.Root>
      <Card.Header>{exerciseName}</Card.Header>
      <Card.Body>
        <Table.Root size="sm">
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Set</Table.ColumnHeader>
              <Table.ColumnHeader>Reps</Table.ColumnHeader>
              <Table.ColumnHeader textAlign="end">kg</Table.ColumnHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {items.map((item) => (
              <Table.Row key={item.setNumber}>
                <Table.Cell>{item.setNumber}</Table.Cell>
                <Table.Cell>{item.repsCompleted}</Table.Cell>
                <Table.Cell textAlign="end">{item.weightCompleted}</Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      </Card.Body>
      <Card.Footer></Card.Footer>
    </Card.Root>
  );
};

export default CompletedExerciseLog;
