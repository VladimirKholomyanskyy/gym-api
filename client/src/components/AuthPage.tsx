import React from "react";
import { Box, Button, Heading, Text, VStack, Spinner } from "@chakra-ui/react";
import { useAuth } from "react-oidc-context";

const AuthPage: React.FC = () => {
  const oidc = useAuth();

  if (oidc.isLoading) {
    return (
      <Box display="flex" alignItems="center" justifyContent="center" height="100vh">
        <Spinner size="xl" />
      </Box>
    );
  }

  return (
    <Box
      display="flex"
      alignItems="center"
      justifyContent="center"
      height="100vh"
      bg="gray.50"
      p={4}
    >
      <VStack
        gap={6}
        p={6}
        bg="white"
        boxShadow="md"
        borderRadius="lg"
        width={{ base: "90%", sm: "400px" }}
      >
        <Heading as="h1" size="lg" color="black">
          Gym App
        </Heading>
        {!oidc.isAuthenticated ? (
          <>
            <Text color="black">Login to track your workouts and progress!</Text>
            <Button
              colorScheme="teal"
              colorPalette="blue"
              onClick={() => oidc.signinRedirect()}
              width="full"
            >
              Login
            </Button>
          </>
        ) : (
          <>
            <Text color="black">Welcome, {oidc.user?.profile?.name || "User"}!</Text>
            <Button
              colorScheme="red"
              colorPalette="blue"
              onClick={() => oidc.signoutRedirect()}
              width="full"
            >
              Logout
            </Button>
          </>
        )}
      </VStack>
    </Box>
  );
};

export default AuthPage;
