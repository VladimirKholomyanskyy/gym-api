package training

import "gorm.io/gorm"

type TrainingProgramRepository struct {
	db *gorm.DB
}

func NewTrainingProgramRepository(db *gorm.DB) *TrainingProgramRepository {
	return &TrainingProgramRepository{db: db}
}

func (r *TrainingProgramRepository) Create(trainingProgram *TrainingProgram) error {
	return r.db.Create(trainingProgram).Error
}

func (r *TrainingProgramRepository) FindById(id uint) (*TrainingProgram, error) {
	var trainingProgram TrainingProgram
	err := r.db.First(&trainingProgram, id).Error
	if err != nil {
		return nil, err
	}
	return &trainingProgram, nil
}

func (r *TrainingProgramRepository) FindByUserId(userID uint) ([]TrainingProgram, error) {
	var trainingPrograms []TrainingProgram
	err := r.db.Where("user_id = ?", userID).Find(&trainingPrograms).Error
	if err != nil {
		return nil, err
	}
	return trainingPrograms, nil
}

func (r *TrainingProgramRepository) Update(trainingProgram *TrainingProgram) error {
	return r.db.Save(trainingProgram).Error
}

func (r *TrainingProgramRepository) Delete(program_id uint, user_id uint) error {
	return r.db.Where("user_id = ?", user_id).Delete(&TrainingProgram{}, program_id).Error
}

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Create(exercise *Exercise) error {
	return r.db.Create(exercise).Error
}
func (r *ExerciseRepository) FindAll() ([]Exercise, error) {
	var exercises []Exercise
	err := r.db.Find(&exercises).Error
	return exercises, err
}

func (r *ExerciseRepository) FindById(id uint) (*Exercise, error) {
	var exercise Exercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

func (r *ExerciseRepository) FindByPrimaryMuscle(primaryMuscle string) ([]Exercise, error) {
	var exercises []Exercise
	if err := r.db.Where("primary_muscle = ?", primaryMuscle).Find(&exercises).Error; err != nil {
		return nil, err
	}
	return exercises, nil
}

type WorkoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) Create(workout *Workout) error {
	return r.db.Create(workout).Error
}

func (r *WorkoutRepository) FindById(id uint) (*Workout, error) {
	var workout Workout
	err := r.db.Preload("Exercises.Exercise").First(&workout, id).Error
	return &workout, err
}

func (r *WorkoutRepository) FindByTrainingProgramId(id uint) ([]Workout, error) {
	var workouts []Workout
	err := r.db.Where("training_program_id = ?", id).
		Preload("Exercises.Exercise").
		Find(&workouts).Error
	return workouts, err
}

func (r *WorkoutRepository) Update(workout *Workout) error {
	return r.db.Save(workout).Error
}

func (r *WorkoutRepository) Delete(id uint) error {
	return r.db.Delete(&Workout{}, id).Error
}

type WorkoutExerciseRepository struct {
	db *gorm.DB
}

func NewWorkoutExerciseRepository(db *gorm.DB) *WorkoutExerciseRepository {
	return &WorkoutExerciseRepository{db: db}
}

func (r *WorkoutExerciseRepository) Create(exercise *WorkoutExercise) error {
	return r.db.Create(exercise).Error
}

func (r *WorkoutExerciseRepository) FindById(id uint) (*WorkoutExercise, error) {
	var exercise WorkoutExercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

func (r *WorkoutExerciseRepository) FindByWorkoutId(id uint) ([]WorkoutExercise, error) {
	var exercises []WorkoutExercise
	err := r.db.Where("workout_id = ?", id).Find(&exercises).Error
	return exercises, err
}

func (r *WorkoutExerciseRepository) FindByExerciseId(id uint) ([]WorkoutExercise, error) {
	var exercises []WorkoutExercise
	err := r.db.Where("exercise_id = ?", id).Find(&exercises).Error
	return exercises, err
}

func (r *WorkoutExerciseRepository) Update(workoutExercise *WorkoutExercise) error {
	return r.db.Save(workoutExercise).Error
}

func (r *WorkoutExerciseRepository) Delete(id uint) error {
	return r.db.Delete(&WorkoutExercise{}, id).Error
}
