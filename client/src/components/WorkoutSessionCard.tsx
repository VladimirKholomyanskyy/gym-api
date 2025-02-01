import { formatDateTime } from "@/utils/dateUtils";
import { Card, Flex, LinkBox, LinkOverlay, Stack } from "@chakra-ui/react";
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
          <Card.Body gap={4}>
            <Stack>
              <LinkOverlay asChild>
                <Link
                  to={{
                    pathname: `/workout-sessions/${workoutSessionId}/${
                      sessionCompleted ? "view" : "edit"
                    }`,
                  }}
                >
                  <Card.Title
                    fontSize="xl"
                    fontWeight="bold"
                    textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
                    cursor="pointer"
                    color="neon.400"
                    _hover={{ color: "neon.300" }}
                  >
                    {workoutName}
                  </Card.Title>
                </Link>
              </LinkOverlay>
              <Card.Description>{shortDescription}</Card.Description>
              <Card.Description fontSize="md" color="gray.300">
                Started: {formatDateTime(sessionStart)}
              </Card.Description>
              {sessionCompleted && (
                <Card.Description fontSize="md" color="gray.300">
                  Completed: {formatDateTime(sessionCompleted)}
                </Card.Description>
              )}
            </Stack>
          </Card.Body>
          <Card.Footer />
        </Flex>
      </Card.Root>
    </LinkBox>
  );
};

export default WorkoutSessionCard;
