import { ExerciseLogsApi } from "@/api";
import { apiConfig } from "@/api/apiConfig";
import { WeightPerDayResponseTotalWeightPerDayInner } from "@/api/models";
import { Box, Spinner, Text, VStack } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import {
  VictoryAxis,
  VictoryChart,
  VictoryLine,
  VictoryTheme,
  VictoryTooltip,
} from "victory";

interface ExerciseProgressProps {
  exerciseId: string;
  startDate?: string;
  endDate?: string;
}

const ExerciseProgress: React.FC<ExerciseProgressProps> = ({
  exerciseId,
  startDate,
  endDate,
}) => {
  const [data, setData] = useState<
    WeightPerDayResponseTotalWeightPerDayInner[]
  >([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const exerciseLogsApi = new ExerciseLogsApi(apiConfig);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response = await exerciseLogsApi.getWeightPerDay(
          exerciseId,
          startDate,
          endDate
        );
        const responseData = response.data?.totalWeightPerDay || [];
        setData(responseData);
        setError(null);
      } catch (err) {
        setError("Failed to fetch data. Please try again.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [exerciseId, startDate, endDate]);

  if (loading) {
    return (
      <Box textAlign="center" mt={4} color="neon.400">
        <Spinner size="lg" color="neon.400" />
      </Box>
    );
  }

  if (error) {
    return (
      <Box textAlign="center" mt={4}>
        <Text color="red.500" fontWeight="bold">
          {error}
        </Text>
      </Box>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Box textAlign="center" mt={4}>
        <Text color="gray.400" fontStyle="italic">
          No data available for the selected exercise.
        </Text>
      </Box>
    );
  }

  // Transform data for Victory
  const transformedData = data.map((item) => ({
    x: item.date || "",
    y: item.totalWeight || 0,
  }));

  return (
    <VStack gap={4} align="stretch" mt={4}>
      <Box
        p={6}
        border="1px solid"
        borderColor="neon.400"
        background="blackAlpha.900"
        boxShadow="0 0 15px rgba(0, 255, 255, 0.8)"
        borderRadius="lg"
      >
        <VictoryChart
          theme={VictoryTheme.material}
          domainPadding={20}
          animate={{ duration: 1500, easing: "elasticInOut" }}
          style={{
            parent: { maxWidth: "100%", overflow: "hidden" },
          }}
        >
          <VictoryAxis
            tickFormat={(tick) => new Date(tick).toLocaleDateString()}
            label="Date"
            style={{
              axis: { stroke: "#00fff0" },
              axisLabel: { fontSize: 14, padding: 30, fill: "#66ffee" },
              tickLabels: { fontSize: 10, angle: -45, fill: "#d4d4d8" },
              grid: { stroke: "rgba(0, 255, 255, 0.3)" },
            }}
          />
          <VictoryAxis
            dependentAxis
            label="Total Weight (kg)"
            style={{
              axis: { stroke: "#00fff0" },
              axisLabel: { fontSize: 14, padding: 40, fill: "#66ffee" },
              tickLabels: { fontSize: 12, fill: "#d4d4d8" },
              grid: { stroke: "rgba(0, 255, 255, 0.3)" },
            }}
          />
          <VictoryLine
            data={transformedData}
            style={{
              data: {
                stroke: "#00fff0",
                strokeWidth: 3,
                filter: "drop-shadow(0 0 10px cyan)",
              },
              labels: { fontSize: 12, fill: "#66ffee" },
            }}
            labels={({ datum }) => `${datum.y} kg`}
            labelComponent={
              <VictoryTooltip
                style={{ fill: "black" }}
                flyoutStyle={{ fill: "#66ffee" }}
              />
            }
          />
        </VictoryChart>
      </Box>
    </VStack>
  );
};

export default ExerciseProgress;
