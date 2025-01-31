import { Text, HStack, Stack } from "@chakra-ui/react";
import { Avatar } from "./ui/avatar";
import { Link as ChakraLink } from "@chakra-ui/react";
import { Link } from "react-router";

const UserProfile = () => {
  return (
    <Stack gap="8" paddingLeft="8" paddingRight="8">
      <HStack gap="4">
        <Avatar
          name="Uchiha Sasuke"
          size="lg"
          src="https://cdn.myanimelist.net/r/84x124/images/characters/9/131317.webp?s=d4b03c7291407bde303bc0758047f6bd"
        />
        <Stack gap="0">
          <Text fontWeight="medium">Uchiha Sasuke</Text>
        </Stack>
      </HStack>
      <ChakraLink asChild>
        <Link to={`/workout-sessions`}>
          <Text fontWeight="bold">Goals</Text>
        </Link>
      </ChakraLink>
      <ChakraLink asChild>
        <Link to={`/workout-sessions`}>
          <Text fontWeight="bold">Workout History</Text>
        </Link>
      </ChakraLink>
      <ChakraLink asChild>
        <Link to={`/settings`}>
          <Text fontWeight="bold">Settings</Text>
        </Link>
      </ChakraLink>
    </Stack>
  );
};

export default UserProfile;
