import { listExercises } from "@/api/exercises";
import { getTrainingProgram } from "@/api/trainingPrograms";
import { getAllWorkouts } from "@/api/workouts";
import { TrainingProgram, Workout, Exercise } from "@/types/api";
import { useEffect, useState } from "react";

export const useProgramData = (programId: string) => {
  const [program, setProgram] = useState<TrainingProgram>();
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [exercises, setExercises] = useState<Exercise[]>([]);

  useEffect(() => {
    let isMounted = true;
    const fetchData = async () => {
      try {
        const programData = await getTrainingProgram(programId);
        const workoutData = await getAllWorkouts(programId);
        const exerciseData = await listExercises();
        if (isMounted) {
          setProgram(programData);
          setWorkouts(workoutData);
          setExercises(exerciseData);
        }
      } catch (error) {
        console.error("Error fetching program data:", error);
      }
    };

    fetchData();
    return () => {
      isMounted = false;
    };
  }, [programId]);

  return { program, workouts, exercises, setWorkouts };
};
