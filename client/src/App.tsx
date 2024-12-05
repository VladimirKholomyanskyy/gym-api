import "./App.css";
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router";
import AuthPage from "./components/AuthPage";
import PrivateRoute from "./components/PrivateRoute";
import Dashboard from "./components/Dashboard";
import TrainingProgramFlow from "./components/TrainingProgramFlow";
import BottomNav from "./components/BottomNavigate";
import { Box, Flex } from "@chakra-ui/react";
import TrainingPrograms from "./components/TrainingProgram";

const ProtectedDashboard = PrivateRoute(Dashboard);
const ProtectedTrainingProgramFlow = PrivateRoute(TrainingProgramFlow);
const ProtectedTrainingProgram = PrivateRoute(TrainingPrograms);

function App() {
  return (
    <Router>
      <Flex direction="column" height="100vh" pb="60px">
        <Box flex="1" overflowY="auto">
          <Routes>
            <Route path="/" element={<AuthPage />} />
            <Route path="/dashboard" element={<ProtectedDashboard />} />
            <Route
              path="/training-programs"
              element={<ProtectedTrainingProgram />}
            />
          </Routes>
        </Box>
        <BottomNav />
      </Flex>
    </Router>
  );
}

export default App;
