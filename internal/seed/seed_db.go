package seed

import (
	"context"
	"log"

	"math/rand"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	trainingmodels "github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	trainingrepos "github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type DatabaseSeed struct {
	exerciseRepo        trainingrepos.ExerciseRepository
	workoutRepo         trainingrepos.WorkoutRepository
	trainingProgramRepo trainingrepos.TrainingProgramRepository
	workoutExerciseRepo trainingrepos.WorkoutExerciseRepository
	profileRepo         account.ProfileRepository
	settingsRepo        account.SettingRepository
}

func NewDatabaseSeed(
	exerciseRepo trainingrepos.ExerciseRepository,
	workoutRepo trainingrepos.WorkoutRepository,
	trainingProgramRepo trainingrepos.TrainingProgramRepository,
	workoutExerciseRepo trainingrepos.WorkoutExerciseRepository,
	profileRepo account.ProfileRepository,
	settingsRepo account.SettingRepository,

) *DatabaseSeed {
	return &DatabaseSeed{
		exerciseRepo:        exerciseRepo,
		workoutRepo:         workoutRepo,
		trainingProgramRepo: trainingProgramRepo,
		workoutExerciseRepo: workoutExerciseRepo,
		profileRepo:         profileRepo,
		settingsRepo:        settingsRepo,
	}
}

func (d *DatabaseSeed) Seed() {
	ctx := context.Background()
	_, totalCount, err := d.exerciseRepo.FindAll(ctx, 1, 100)
	if err != nil {
		log.Fatalf("Failed to count exercises: %v", err)
	}
	if totalCount == 0 {
		user := account.Profile{ExternalID: "a5ce12b2-3d4d-439c-ac8d-cd5ca5d8ea33"}
		d.profileRepo.Create(ctx, &user)
		// Define exercises to populate
		exercises := []trainingmodels.Exercise{
			{Name: "Bench Press", PrimaryMuscle: "Chest", SecondaryMuscle: []string{"Triceps", "Shoulders"}, Equipment: "Barbell", Description: "A compound chest exercise."},
			{Name: "Squat", PrimaryMuscle: "Legs", SecondaryMuscle: []string{"Glutes", "Lower Back"}, Equipment: "Barbell", Description: "A compound leg exercise."},
			{Name: "Deadlift", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Hamstrings", "Glutes"}, Equipment: "Barbell", Description: "A compound back and leg exercise."},
			{Name: "Pull-Up", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Biceps"}, Equipment: "Bodyweight", Description: "A compound upper body exercise."},
		}

		// Insert exercises
		var allExercises []trainingmodels.Exercise
		for _, exercise := range exercises {
			if err := d.exerciseRepo.Create(ctx, &exercise); err != nil {
				log.Fatalf("Failed to seed exercises: %v", err)
			}
			allExercises = append(allExercises, exercise)
		}

		log.Println("Seeded exercises table with initial data.")
		userID := user.ID
		programs := []trainingmodels.TrainingProgram{
			{Name: "Muscle growth", Description: "High volume workouts", ProfileID: userID},
			{Name: "Endurance", Description: "High reps count", ProfileID: userID},
			{Name: "Strength", Description: "Increase overall strength", ProfileID: userID},
			{Name: "Strength", Description: "Increase overall strength", ProfileID: userID},
		}
		var allPrograms []trainingmodels.TrainingProgram
		for _, program := range programs {
			d.trainingProgramRepo.Create(ctx, &program)
			allPrograms = append(allPrograms, program)
		}

		log.Println("Seeded training programs table with initial data.")

		var allWorkouts []trainingmodels.Workout
		for _, program := range allPrograms {
			workouts := []trainingmodels.Workout{
				{Name: "Workout 1", TrainingProgramID: program.ID},
				{Name: "Workout 2", TrainingProgramID: program.ID},
				{Name: "Workout 3", TrainingProgramID: program.ID},
				{Name: "Workout 4", TrainingProgramID: program.ID},
			}
			for _, workout := range workouts {
				d.workoutRepo.Create(ctx, &workout)
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
				d.workoutExerciseRepo.Create(ctx, &trainingmodels.WorkoutExercise{
					WorkoutID:  workout.ID,
					ExerciseID: exercise.ID,
					Sets:       r.Intn(maxSets-minSets+1) + minSets,
					Reps:       r.Intn(maxReps-minReps+1) + minReps,
				})
			}
		}
		log.Println("Seeded workouts table with initial data.")

	} else {
		log.Println("Exercises table already populated. Skipping seed.")
	}
}
