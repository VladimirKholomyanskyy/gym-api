import { Flex, IconButton } from "@chakra-ui/react";
import { useNavigate } from "react-router";
import { FaHome, FaDumbbell, FaClipboardList, FaUser } from "react-icons/fa";

const BottomNav = () => {
  const navigate = useNavigate();

  return (
    <Flex gap="4" justify="center">
      <IconButton
        aria-label="Home"
        onClick={() => navigate("/home")}
        variant="ghost"
        color="white"
      >
        <FaHome />
      </IconButton>
      <IconButton
        aria-label="Training Programs"
        onClick={() => navigate("/training-programs")}
        variant="ghost"
        color="white"
      >
        <FaClipboardList />
      </IconButton>
      <IconButton
        aria-label="Workouts"
        onClick={() => navigate("/workout-sessions")}
        variant="ghost"
        color="white"
      >
        <FaDumbbell />
      </IconButton>
      <IconButton
        aria-label="Profile"
        onClick={() => navigate("/profile")}
        variant="ghost"
        color="white"
      >
        <FaUser />
      </IconButton>
    </Flex>
  );
};

export default BottomNav;
