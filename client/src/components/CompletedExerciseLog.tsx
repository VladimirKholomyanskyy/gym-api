import { Card, Flex, Table } from "@chakra-ui/react";

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
    <Card.Root
      size="sm"
      width="100%"
      background="blackAlpha.800"
      borderRadius="md"
      boxShadow="0 0 10px rgba(0, 255, 255, 0.8)"
      p={2}
      _hover={{ boxShadow: "0 0 20px rgba(0, 255, 255, 1)" }}
    >
      <Card.Header
        fontSize="xl"
        fontWeight="bold"
        color="neon.400"
        textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
        cursor="pointer"
        _hover={{ color: "neon.300" }}
      >
        {exerciseName}
      </Card.Header>
      <Flex align="stretch" justify="space-between">
        <Card.Body>
          <Table.Root background="blackAlpha.800" size="sm">
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
                  <Table.Cell textAlign="end">
                    {item.weightCompleted}
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Root>
        </Card.Body>
        <Card.Footer />
      </Flex>
    </Card.Root>
  );
};

export default CompletedExerciseLog;
