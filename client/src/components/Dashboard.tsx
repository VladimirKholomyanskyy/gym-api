import React from "react";
import { Box, Heading, Text } from "@chakra-ui/react";

const Dashboard: React.FC = () => {
  return (
    <Box p={4}>
      <Heading as="h2" size="lg">Dashboard</Heading>
      <Text>Welcome to your dashboard! You are authenticated.</Text>
    </Box>
  );
};

export default Dashboard;
