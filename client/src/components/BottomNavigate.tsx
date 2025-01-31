import { Box, Flex, IconButton } from "@chakra-ui/react";
import { useLocation, useNavigate } from "react-router";
import { FaHome, FaDumbbell, FaClipboardList, FaUser } from "react-icons/fa";
import { motion } from "framer-motion";

const MotionIconButton = motion(IconButton);

const BottomNav = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { label: "Home", icon: FaHome, route: "/home" },
    { label: "Training", icon: FaClipboardList, route: "/training-programs" },
    { label: "Exercises", icon: FaDumbbell, route: "/exercises" },
    { label: "Profile", icon: FaUser, route: "/profile" },
  ];

  return (
    <Flex
      justify="center"
      gap="2"
      p="3"
      bg="blackAlpha.900"
      boxShadow="0px -4px 10px rgba(0, 255, 255, 0.3)"
      borderTop="2px solid rgba(0, 255, 255, 0.6)"
      position="fixed"
      bottom="0"
      left="0"
      right="0"
      zIndex="1000"
      h="60px"
      align="center"
    >
      {navItems.map(({ label, icon: Icon, route }) => {
        const isActive = location.pathname === route;
        console.log("isActive:", isActive);

        return (
          <Flex
            key={label}
            direction="column"
            align="center"
            position="relative"
          >
            {/* Minimal Glowing Border Under Icon */}
            {isActive && (
              <Box
                position="absolute"
                bottom="-4px"
                w="60%"
                h="2px"
                bg="magenta"
                boxShadow="0px 0px 8px magenta"
                animation="pulse 1s ease-in-out infinite"
              />
            )}
            <MotionIconButton
              aria-label={label}
              onClick={() => navigate(route)}
              variant="ghost"
              color={isActive ? "magenta.400" : "cyan.400"}
              _hover={{
                color: "cyan.300", // Glowing hover effect
                textShadow: "0px 0px 8px cyan", // Soft glow on hover
              }}
              _active={{
                color: "magenta.400", // Glowing color on click
                textShadow: "0px 0px 12px magenta", // Intense glow on click
                scale: 0.95, // Slight shrink effect on click
              }}
              transition={{ type: "spring", stiffness: 300 }}
            >
              {" "}
              <Icon />
            </MotionIconButton>
          </Flex>
        );
      })}
    </Flex>
  );
};

export default BottomNav;
