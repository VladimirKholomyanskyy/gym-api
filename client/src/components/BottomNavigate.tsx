import { Flex, IconButton } from "@chakra-ui/react";
import { useNavigate } from "react-router";
import { FaHome, FaDumbbell, FaClipboardList, FaUser } from "react-icons/fa";

const BottomNav = () => {
  const navigate = useNavigate();

  return (
    <Flex gap="4" justify="center" padding="2" background="bg.subtle">
      <IconButton
        aria-label="Home"
        onClick={() => navigate("/home")}
        variant="ghost"
        color="pink.600"
      >
        <FaHome />
      </IconButton>
      <IconButton
        aria-label="Training Programs"
        onClick={() => navigate("/training-programs")}
        variant="ghost"
        color="pink.600"
      >
        <FaClipboardList />
      </IconButton>
      <IconButton
        aria-label="Workouts"
        onClick={() => navigate("/workout-sessions")}
        variant="ghost"
        color="pink.600"
      >
        <FaDumbbell />
      </IconButton>
      <IconButton
        aria-label="Profile"
        onClick={() => navigate("/profile")}
        variant="ghost"
        color="pink.600"
      >
        <FaUser />
      </IconButton>
    </Flex>
  );
};

export default BottomNav;
