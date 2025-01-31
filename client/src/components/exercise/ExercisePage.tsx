import { ExercisesApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";
import { Exercise } from "@/api/models";
import { Box, Flex, Heading, Link, Tabs } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import DescriptionTab from "./DescriptionTab";
import LogsTab from "./LogsTab";
import ExerciseProgress from "./ProgressTab";

const ExercisePage = () => {
  const { exerciseId } = useParams();
  const [loading, setLoading] = useState(true);
  const [exercise, setExercise] = useState<Exercise | null>(null);
  const navigate = useNavigate();
  const exerciseApi = new ExercisesApi(apiConfig);
  if (!exerciseId) {
    navigate("/error");
    return null;
  }
  useEffect(() => {
    const fetchExercise = async () => {
      try {
        const response = await exerciseApi.getExerciseById(exerciseId);
        setExercise(response.data);
      } catch (error) {
        console.log("Error", error);
      }
    };
    fetchExercise();
  }, [exerciseId]);
  return (
    <Box width="100%" minHeight="100vh" background="bg.subtle" p={6}>
      <Heading
        size="2xl"
        fontWeight="bold"
        textAlign="center"
        color="magenta.400"
        textShadow="0 0 10px rgba(255, 0, 255, 0.8)"
      >
        {exercise?.name}
      </Heading>

      <Tabs.Root fitted defaultValue="description">
        <Tabs.List>
          <Tabs.Trigger value="description" asChild>
            <Link
              fontWeight="bold"
              color="blue.300"
              textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
              cursor="pointer"
              _hover={{ color: "neon.300" }}
              unstyled
              href="#description"
            >
              Description
            </Link>
          </Tabs.Trigger>
          <Tabs.Trigger value="logs" asChild>
            <Link
              fontWeight="bold"
              color="blue.300"
              textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
              cursor="pointer"
              _hover={{ color: "neon.300" }}
              unstyled
              href="#logs"
            >
              Logs
            </Link>
          </Tabs.Trigger>
          <Tabs.Trigger value="progress" asChild>
            <Link
              fontWeight="bold"
              color="blue.300"
              textShadow="0 0 10px rgba(0, 255, 255, 0.8)"
              cursor="pointer"
              _hover={{ color: "neon.300" }}
              unstyled
              href="#progress"
            >
              Progress
            </Link>
          </Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="description">
          {exercise ? (
            <DescriptionTab exercise={exercise} />
          ) : (
            <Flex
              justifyContent="center"
              alignItems="center"
              height="50vh"
            ></Flex>
          )}
        </Tabs.Content>
        <Tabs.Content value="logs">
          <LogsTab exerciseId={exerciseId} />
        </Tabs.Content>
        <Tabs.Content value="progress">
          <ExerciseProgress exerciseId={exerciseId} />
        </Tabs.Content>
      </Tabs.Root>
    </Box>
  );
};

export default ExercisePage;
