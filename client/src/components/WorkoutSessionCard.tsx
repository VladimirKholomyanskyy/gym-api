import { formatDateTime } from "@/utils/dateUtils";
import { Card, LinkBox, LinkOverlay, Stack, Text } from "@chakra-ui/react";
import { Link } from "react-router";

export interface WorkoutSessionCardProps {
  workoutSessionId: string;
  workoutName: string;
  sessionStart: string;
  sessionCompleted?: string;
  shortDescription?: string;
}

const WorkoutSessionCard = ({
  workoutSessionId,
  workoutName,
  sessionStart,
  sessionCompleted,
  shortDescription,
}: WorkoutSessionCardProps) => {
  return (
    <LinkBox>
      <Card.Root size="sm" borderRadius="none">
        <Card.Header>
          <LinkOverlay asChild>
            <Link
              to={{
                pathname: `/workout-sessions/${workoutSessionId}/${
                  sessionCompleted ? "view" : "edit"
                }`,
              }}
            >
              {workoutName}
            </Link>
          </LinkOverlay>
        </Card.Header>

        <Card.Body>{shortDescription}</Card.Body>
        <Card.Footer
          display="flex"
          justifyContent="space-between"
          alignItems="center"
        >
          <Stack>
            <Text>Started: {formatDateTime(sessionStart)}</Text>
            {sessionCompleted && (
              <Text>Completed: {formatDateTime(sessionCompleted)}</Text>
            )}
          </Stack>
        </Card.Footer>
      </Card.Root>
    </LinkBox>
  );
};

export default WorkoutSessionCard;
