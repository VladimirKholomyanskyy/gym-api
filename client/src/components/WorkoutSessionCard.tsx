import { Card, Stack, Text } from "@chakra-ui/react";

export interface WorkoutSessionCardProps {
  workoutName: string;
  sessionStart: string;
  sessionCompleted?: string;
  shortDescription: string;
}
const WorkoutSessionCard = ({
  workoutName,
  sessionStart,
  sessionCompleted,
  shortDescription,
}: WorkoutSessionCardProps) => {
  return (
    <Card.Root size="sm">
      <Card.Header>{workoutName}</Card.Header>
      <Card.Body>{shortDescription}</Card.Body>
      <Card.Footer
        display="flex"
        justifyContent="space-between"
        alignItems="center"
      >
        <Stack>
          <Text>Started: {sessionStart}</Text>
          {sessionCompleted && <Text>Completed: {sessionCompleted}</Text>}
        </Stack>
      </Card.Footer>
    </Card.Root>
  );
};

export default WorkoutSessionCard;
