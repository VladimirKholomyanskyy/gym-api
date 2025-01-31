package seed

import (
	"log"
	"math/rand"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training"
)

type DatabaseSeed struct {
	exerciseRepo        *training.ExerciseRepository
	workoutRepo         *training.WorkoutRepository
	trainingProgramRepo *training.TrainingProgramRepository
	workoutExerciseRepo *training.WorkoutExerciseRepository
	userRepository      *account.UserRepository
}

func NewDatabaseSeed(
	exerciseRepo *training.ExerciseRepository,
	workoutRepo *training.WorkoutRepository,
	trainingProgramRepo *training.TrainingProgramRepository,
	workoutExerciseRepo *training.WorkoutExerciseRepository,
	userRepository *account.UserRepository,
) *DatabaseSeed {
	return &DatabaseSeed{exerciseRepo: exerciseRepo, workoutRepo: workoutRepo, trainingProgramRepo: trainingProgramRepo, workoutExerciseRepo: workoutExerciseRepo, userRepository: userRepository}
}

func (d *DatabaseSeed) Seed() {
	exercises, err := d.exerciseRepo.FindAll()
	if err != nil {
		log.Fatalf("Failed to count exercises: %v", err)
	}
	if len(exercises) == 0 {
		user := account.User{ExternalID: "a5ce12b2-3d4d-439c-ac8d-cd5ca5d8ea33"}
		d.userRepository.CreateUser(&user)
		// Define exercises to populate
		exercises := []training.Exercise{
			{Name: "Bench Press", PrimaryMuscle: "Chest", SecondaryMuscle: []string{"Triceps", "Shoulders"}, Equipment: "Barbell", Description: "A compound chest exercise."},
			{Name: "Squat", PrimaryMuscle: "Legs", SecondaryMuscle: []string{"Glutes", "Lower Back"}, Equipment: "Barbell", Description: "A compound leg exercise."},
			{Name: "Deadlift", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Hamstrings", "Glutes"}, Equipment: "Barbell", Description: "A compound back and leg exercise."},
			{Name: "Pull-Up", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Biceps"}, Equipment: "Bodyweight", Description: "A compound upper body exercise."},
		}

		// Insert exercises
		var allExercises []training.Exercise
		for _, exercise := range exercises {
			if err := d.exerciseRepo.Create(&exercise); err != nil {
				log.Fatalf("Failed to seed exercises: %v", err)
			}
			allExercises = append(allExercises, exercise)
		}

		log.Println("Seeded exercises table with initial data.")
		userID := user.ID
		programs := []training.TrainingProgram{
			{Name: "Muscle growth", Description: "High volume workouts", UserID: userID},
			{Name: "Endurance", Description: "High reps count", UserID: userID},
			{Name: "Strength", Description: "Increase overall strength", UserID: userID},
			{Name: "Strength", Description: "Increase overall strength", UserID: userID},
		}
		var allPrograms []training.TrainingProgram
		for _, program := range programs {
			d.trainingProgramRepo.Create(&program)
			allPrograms = append(allPrograms, program)
		}

		log.Println("Seeded training programs table with initial data.")

		var allWorkouts []training.Workout
		for _, program := range allPrograms {
			workouts := []training.Workout{
				{Name: "Workout 1", TrainingProgramID: program.ID},
				{Name: "Workout 2", TrainingProgramID: program.ID},
				{Name: "Workout 3", TrainingProgramID: program.ID},
				{Name: "Workout 4", TrainingProgramID: program.ID},
			}
			for _, workout := range workouts {
				d.workoutRepo.Create(&workout)
				allWorkouts = append(allWorkouts, workout)
			}
		}
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		minReps := 6
		maxReps := 12
		minSets := 2
		maxSets := 4
		for _, workout := range allWorkouts {
			for _, exercise := range allExercises {
				d.workoutExerciseRepo.Create(&training.WorkoutExercise{
					WorkoutID:  workout.ID,
					ExerciseID: exercise.ID,
					Sets:       r.Intn(maxSets-minSets+1) + minSets,
					Reps:       r.Intn(maxReps-minReps+1) + minReps,
					Weight:     100})
			}
		}
		log.Println("Seeded workouts table with initial data.")

	} else {
		log.Println("Exercises table already populated. Skipping seed.")
	}
}
