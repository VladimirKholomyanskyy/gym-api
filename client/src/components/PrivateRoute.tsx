import React from "react";
import { withAuthenticationRequired } from "react-oidc-context";
import { Spinner, Text, VStack } from "@chakra-ui/react";

// This higher-order component allows any Component to be passed and protected
const PrivateRoute = (Component: React.FC) =>
  withAuthenticationRequired(Component, {
    OnRedirecting: () => (
      <VStack
        justify="center"
        align="center"
        height="100vh"
        bg="gray.50"
        gap={4}
      >
        <Spinner size="xl" color="teal.500" />
        <Text>Redirecting to the login page...</Text>
      </VStack>
    ),
  });

export default PrivateRoute;
