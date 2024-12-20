import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router";
import AuthPage from "./components/AuthPage";
import PrivateRoute from "./components/PrivateRoute";
import Dashboard from "./components/Dashboard";
import BottomNav from "./components/BottomNavigate";
import { Box, Flex } from "@chakra-ui/react";
import TrainingProgramPage from "./components/TrainingProgramPage";
import WorkoutPage from "./components/WorkoutPage";
import { useAuth } from "react-oidc-context";
import React from "react";
import { setupAxiosInterceptors } from "./api/apiClient";
import TrainingProgramsPage from "./components/TrainingProgramsPage";

const ProtectedDashboard = PrivateRoute(Dashboard);
const ProtectedTrainingProgramPage = PrivateRoute(TrainingProgramPage);
const ProtectedTrainingProgramsPage = PrivateRoute(TrainingProgramsPage);
const ProtectedWorkout = PrivateRoute(WorkoutPage);

function App() {
  const auth = useAuth();

  React.useEffect(() => {
    // Set up the interceptor using the token from the auth context
    setupAxiosInterceptors(() => auth.user?.access_token || null);
  }, [auth]);
  return (
    <Router>
      <Flex direction="column" height="100vh" pb="60px">
        <Box flex="1" overflowY="auto">
          <Routes>
            <Route path="/" element={<AuthPage />} />
            <Route path="/dashboard" element={<ProtectedDashboard />} />
            <Route
              path="/training-programs"
              element={<ProtectedTrainingProgramsPage />}
            />
            <Route
              path="/training-programs/:programId/workouts"
              element={<ProtectedTrainingProgramPage />}
            />
            <Route
              path="/training-programs/:programId/workouts/:workoutId"
              element={<ProtectedWorkout />}
            />
          </Routes>
        </Box>
        <BottomNav />
      </Flex>
    </Router>
  );
}

export default App;
