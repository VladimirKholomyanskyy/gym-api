import React, { useState, useEffect } from "react";
import { Box, Button, Text, HStack } from "@chakra-ui/react";
import {
  ProgressCircleRing,
  ProgressCircleRoot,
  ProgressCircleValueText,
} from "./ui/progress-circle";

const ExerciseTimer: React.FC = () => {
  const [timeLeft, setTimeLeft] = useState<number>(30); // Timer starts at 30 seconds
  const [isRunning, setIsRunning] = useState<boolean>(false);

  useEffect(() => {
    let timer: number | undefined;
    if (isRunning && timeLeft > 0) {
      timer = window.setInterval(() => {
        setTimeLeft((prev) => prev - 1);
      }, 1000);
    } else if (timeLeft === 0) {
      setIsRunning(false);
      clearInterval(timer);
    }
    return () => clearInterval(timer);
  }, [isRunning, timeLeft]);

  const handleStartPause = () => {
    setIsRunning((prev) => !prev);
  };

  const handleReset = () => {
    setIsRunning(false);
    setTimeLeft(30); // Reset to initial time
  };

  return (
    <Box textAlign="center" p={4}>
      <HStack>
        <Text fontSize="xl" mb={4}>
          Rest
        </Text>
        <ProgressCircleRoot value={(timeLeft / 30) * 100}>
          <ProgressCircleValueText>{timeLeft}s</ProgressCircleValueText>
          <ProgressCircleRing css={{ "--thickness": "2px" }} />
        </ProgressCircleRoot>

        <Button
          colorScheme={isRunning ? "red" : "green"}
          onClick={handleStartPause}
        >
          {isRunning ? "Pause" : "Start"}
        </Button>
        <Button onClick={handleReset}>Reset</Button>
      </HStack>
    </Box>
  );
};

export default ExerciseTimer;
