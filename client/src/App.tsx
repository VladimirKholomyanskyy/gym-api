import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router";
import AuthPage from "./components/AuthPage";
import PrivateRoute from "./components/PrivateRoute";
import Dashboard from "./components/Dashboard";
import BottomNav from "./components/BottomNavigate";
import { Box, Flex, Spinner } from "@chakra-ui/react";
import TrainingProgramPage from "./components/TrainingProgramPage";
import WorkoutPage from "./components/WorkoutPage";
import { hasAuthParams, useAuth } from "react-oidc-context";
import React from "react";
import TrainingProgramsPage from "./components/TrainingProgramsPage";
import ReadOnlyWorkoutSessionWrapper from "./components/ReadOnlyWorkoutSessionWrapper";
import EditableWorkoutWrapper from "./components/EditableWorkoutSessionWrapper";
import WorkoutSessionsPage from "./components/WorkoutSessionsPage";

const ProtectedDashboard = PrivateRoute(Dashboard);
const ProtectedTrainingProgramPage = PrivateRoute(TrainingProgramPage);
const ProtectedTrainingProgramsPage = PrivateRoute(TrainingProgramsPage);
const ProtectedWorkout = PrivateRoute(WorkoutPage);
const ProtectedEditableWorkoutSession = PrivateRoute(EditableWorkoutWrapper);
const ProtectedReadOnlyWorkoutSession = PrivateRoute(
  ReadOnlyWorkoutSessionWrapper
);
const ProtectedWorkoutSessions = PrivateRoute(WorkoutSessionsPage);

function App() {
  const auth = useAuth();

  const [hasTriedSignin, setHasTriedSignin] = React.useState(false);

  // automatically sign-in
  React.useEffect(() => {
    if (
      !hasAuthParams() &&
      !auth.isAuthenticated &&
      !auth.activeNavigator &&
      !auth.isLoading &&
      !hasTriedSignin
    ) {
      auth.signinRedirect();
      setHasTriedSignin(true);
    }
  }, [auth, hasTriedSignin]);

  if (auth.isLoading) {
    return <Spinner>Signing you in/out...</Spinner>;
  }

  if (!auth.isAuthenticated) {
    return <div>Unable to log in</div>;
  }

  return (
    <Router>
      <Flex
        direction="column"
        height="100vh"
        className="dark"
        background="bg.subtle"
      >
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
            <Route
              path="/workout-sessions/:id/edit"
              element={<ProtectedEditableWorkoutSession />}
            />
            <Route
              path="/workout-sessions/:id/view"
              element={<ProtectedReadOnlyWorkoutSession />}
            />
            <Route
              path="/workout-sessions"
              element={<ProtectedWorkoutSessions />}
            />
          </Routes>
        </Box>
        <BottomNav />
      </Flex>
    </Router>
  );
}

export default App;
