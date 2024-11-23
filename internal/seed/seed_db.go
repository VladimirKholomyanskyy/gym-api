package seed

import (
	"log"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
)

func SeedExercises(repo *repository.ExerciseRepository) {
	var exercises []models.Exercise
	var err error
	exercises, err = repo.FindAll()
	if err != nil {
		log.Fatalf("Failed to count exercises: %v", err)
	}

	if len(exercises) == 0 {
		// Define exercises to populate
		exercises := []models.Exercise{
			{Name: "Bench Press", PrimaryMuscle: "Chest", SecondaryMuscle: []string{"Triceps", "Shoulders"}, Equipment: "Barbell", Description: "A compound chest exercise."},
			{Name: "Squat", PrimaryMuscle: "Legs", SecondaryMuscle: []string{"Glutes", "Lower Back"}, Equipment: "Barbell", Description: "A compound leg exercise."},
			{Name: "Deadlift", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Hamstrings", "Glutes"}, Equipment: "Barbell", Description: "A compound back and leg exercise."},
			{Name: "Pull-Up", PrimaryMuscle: "Back", SecondaryMuscle: []string{"Biceps"}, Equipment: "Bodyweight", Description: "A compound upper body exercise."},
		}

		// Insert exercises
		for _, exercise := range exercises {
			if err := repo.Create(&exercise); err != nil {
				log.Fatalf("Failed to seed exercises: %v", err)
			}
		}

		log.Println("Seeded exercises table with initial data.")
	} else {
		log.Println("Exercises table already populated. Skipping seed.")
	}
}
