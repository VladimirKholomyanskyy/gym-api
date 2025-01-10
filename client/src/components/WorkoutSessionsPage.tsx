import { getWorkoutSessions } from "@/api/workout-sessions";
import { useState, useEffect } from "react";
import WorkoutSessionCard, {
  WorkoutSessionCardProps,
} from "./WorkoutSessionCard";
import { For, Heading, Spinner, Stack } from "@chakra-ui/react";

const WorkoutSessionsPage = () => {
  const [workoutSessionsCardProps, setWorkoutSessionsCardProps] = useState<
    WorkoutSessionCardProps[]
  >([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchWorkoutSessions = async () => {
      try {
        const sessions = await getWorkoutSessions();
        const sessionProps = sessions.map((session) => {
          const shortDesc = session.workout_snapshot.Exercises.map(
            (e) => e.Exercise.Name
          ).join(",");
          return {
            workoutName: session.workout_snapshot.Name,
            sessionStart: session.started_at,
            sessionCompleted: session.completed_at,
            shortDescription: shortDesc,
          };
        });
        setWorkoutSessionsCardProps(sessionProps);
      } catch (error) {
        console.error("Read Error fetching workout session:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkoutSessions();
  }, []);

  if (loading) return <Spinner>Loading...</Spinner>;
  return (
    <Stack>
      <Heading>Workout Sessions</Heading>
      <For each={workoutSessionsCardProps}>
        {(item, index) => (
          <WorkoutSessionCard
            key={index}
            workoutName={item.workoutName}
            sessionStart={item.sessionStart}
            shortDescription={item.shortDescription}
            sessionCompleted={item.sessionCompleted}
          />
        )}
      </For>
    </Stack>
  );
};

export default WorkoutSessionsPage;
