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
import ExercisePage from "./components/exercise/ExercisePage";
import ExercisesPage from "./components/exercise/ExercisesPage";
import UserProfile from "./components/UserProfile";
const ProtectedDashboard = PrivateRoute(Dashboard);
const ProtectedTrainingProgramPage = PrivateRoute(TrainingProgramPage);
const ProtectedTrainingProgramsPage = PrivateRoute(TrainingProgramsPage);
const ProtectedWorkout = PrivateRoute(WorkoutPage);
const ProtectedEditableWorkoutSession = PrivateRoute(EditableWorkoutWrapper);
const ProtectedReadOnlyWorkoutSession = PrivateRoute(
  ReadOnlyWorkoutSessionWrapper
);
const ProtectedWorkoutSessions = PrivateRoute(WorkoutSessionsPage);
const ProtectedExercisePage = PrivateRoute(ExercisePage);
const ProtectedExercisesPage = PrivateRoute(ExercisesPage);
const ProtectedUserProfilePage = PrivateRoute(UserProfile);

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
      <Flex direction="column" height="100vh">
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
            <Route
              path="/exercises/:exerciseId"
              element={<ProtectedExercisePage />}
            />
            <Route path="/exercises" element={<ProtectedExercisesPage />} />
            <Route path="/profile" element={<ProtectedUserProfilePage />} />
          </Routes>
        </Box>
        <BottomNav />
      </Flex>
    </Router>
  );
}

export default App;
