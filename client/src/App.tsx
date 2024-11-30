import './App.css'
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router";
import AuthPage from "./components/AuthPage";
import PrivateRoute from "./components/PrivateRoute";
import Dashboard from "./components/Dashboard";
import TrainingProgramFlow from './components/TrainingProgramFlow';

const ProtectedDashboard = PrivateRoute(Dashboard);
const ProtectedTrainingProgramFlow = PrivateRoute(TrainingProgramFlow);

function App() {
    return (
      <Router>
        <Routes>
          <Route path="/" element={<AuthPage />} />
          <Route path="/dashboard" element={<ProtectedDashboard/>} />
          <Route path="/training-program" element={<ProtectedTrainingProgramFlow />} />
        </Routes>
      </Router>
    );
  }


export default App
