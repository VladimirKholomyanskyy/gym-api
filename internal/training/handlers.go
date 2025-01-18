package training

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
)

type ExerciseHandler struct {
	service *TrainingManager
}

func NewExerciseHandler(service *TrainingManager) *ExerciseHandler {
	return &ExerciseHandler{service: service}
}

func (h *ExerciseHandler) ListExercises(ctx context.Context) (openapi.ImplResponse, error) {

	exercises, err := h.service.GetAllExercises()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var convertedExercises []openapi.Exercise
	for _, exercise := range exercises {
		convertedExercises = append(convertedExercises,
			openapi.Exercise{
				Id:              fmt.Sprintf("%d", exercise.ID),
				Name:            exercise.Name,
				PrimaryMuscle:   exercise.PrimaryMuscle,
				SecondaryMuscle: exercise.SecondaryMuscle,
				Equipment:       exercise.Equipment,
				Description:     exercise.Description,
			})
	}

	return openapi.Response(http.StatusOK, convertedExercises), nil
}

func (h *ExerciseHandler) GetExerciseById(ctx context.Context, exerciseId string) (openapi.ImplResponse, error) {
	id, err := strconv.Atoi(exerciseId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	exercise, err := h.service.GetExercise(uint(id))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	if exercise == nil {
		return openapi.Response(http.StatusNotFound, nil), nil
	}
	return openapi.Response(http.StatusOK,
		openapi.Exercise{
			Id:              strconv.Itoa(id),
			Name:            exercise.Name,
			PrimaryMuscle:   exercise.PrimaryMuscle,
			SecondaryMuscle: exercise.SecondaryMuscle,
			Equipment:       exercise.Equipment,
			Description:     exercise.Description,
		}), nil
}

type TrainingProgramHandler struct {
	service *TrainingManager
}

func NewTrainingProgramHandler(service *TrainingManager) *TrainingProgramHandler {
	return &TrainingProgramHandler{service: service}
}
func (h *TrainingProgramHandler) ListTrainingPrograms(ctx context.Context) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	userPrograms, err := h.service.GetTrainingPrograms(userID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var userProgramsConverted []openapi.TrainingProgram
	for _, userProgram := range userPrograms {
		userProgramsConverted = append(userProgramsConverted,
			openapi.TrainingProgram{
				Id:          fmt.Sprintf("%d", userProgram.ID),
				Name:        userProgram.Name,
				Description: userProgram.Description,
			})
	}
	return openapi.Response(200, userProgramsConverted), nil
}

func (h *TrainingProgramHandler) CreateTrainingProgram(ctx context.Context, createTrainingProgramRequest openapi.CreateTrainingProgramRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	program, err := h.service.CreateTrainingProgram(createTrainingProgramRequest, userID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(201, openapi.TrainingProgram{Id: fmt.Sprintf("%d", program.ID), Name: program.Name, Description: program.Description}), nil
}

func (h *TrainingProgramHandler) GetTrainingProgramById(ctx context.Context, programId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	userProgram, err := h.service.GetTrainingProgram(userID, uint(programID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	if userProgram == nil {
		return openapi.Response(http.StatusNotFound, nil), nil
	}
	return openapi.Response(200, openapi.TrainingProgram{Id: fmt.Sprintf("%d", userProgram.ID), Name: userProgram.Name, Description: userProgram.Description}), nil
}

func (h *TrainingProgramHandler) DeleteTrainingProgram(ctx context.Context, programId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	err = h.service.DeleteTrainingProgram(userID, uint(programID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}

func (h *TrainingProgramHandler) UpdateTrainingProgram(ctx context.Context, programId string, createTrainingProgramRequest openapi.CreateTrainingProgramRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	userProgram, err := h.service.UpdateTrainingProgram(createTrainingProgramRequest, userID, uint(programID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusOK, openapi.TrainingProgram{Id: fmt.Sprintf("%d", userProgram.ID), Name: userProgram.Name, Description: userProgram.Description}), nil
}

func (h *TrainingProgramHandler) ListWorkoutsForProgram(ctx context.Context, programId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workouts, err := h.service.GetAllWorkouts(userID, uint(programID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var workoutsConverted []openapi.WorkoutResponse
	for _, workout := range workouts {
		workoutsConverted = append(workoutsConverted, openapi.WorkoutResponse{Id: fmt.Sprintf("%d", workout.ID), Name: workout.Name})
	}
	return openapi.Response(http.StatusOK, workoutsConverted), nil
}

func (h *TrainingProgramHandler) AddWorkoutToProgram(ctx context.Context, programId string, workoutRequest openapi.WorkoutRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}

	workout, err := h.service.AddWorkoutToProgram(workoutRequest, userID, uint(programID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.WorkoutResponse{Id: fmt.Sprintf("%d", workout.ID), Name: workout.Name}), nil
}

func (h *TrainingProgramHandler) GetWorkoutForProgram(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	workoutID, err := strconv.Atoi(workoutId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workout, err := h.service.GetWorkout(uint(workoutID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	if workout == nil {
		return openapi.Response(http.StatusNotFound, nil), nil
	}
	return openapi.Response(http.StatusOK, openapi.WorkoutResponse{Id: fmt.Sprintf("%d", workout.ID), Name: workout.Name}), nil
}

func (h *TrainingProgramHandler) UpdateWorkout(ctx context.Context, programId string, workoutId string, workoutRequest openapi.WorkoutRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workoutID, err := strconv.Atoi(workoutId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}

	workout, err := h.service.UpdateWorkout(workoutRequest, userID, uint(programID), uint(workoutID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusOK, openapi.WorkoutResponse{Id: fmt.Sprintf("%d", workout.ID), Name: workout.Name}), nil
}

func (h *TrainingProgramHandler) DeleteWorkout(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	programID, err := strconv.Atoi(programId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workoutID, err := strconv.Atoi(workoutId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	err = h.service.DeleteWorkout(userID, uint(programID), uint(workoutID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}

type WorkoutExerciseHandler struct {
	service *TrainingManager
}

func NewWorkoutExerciseHandler(service *TrainingManager) *WorkoutExerciseHandler {
	return &WorkoutExerciseHandler{service: service}
}

func (h *WorkoutExerciseHandler) ListWorkoutExercises(ctx context.Context, workoutId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	workoutID, err := strconv.Atoi(workoutId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workoutExercises, err := h.service.GetAllWorkoutExercisesByWorkout(userID, uint(workoutID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var workoutExercisesConverted []openapi.WorkoutExerciseResponse
	for _, workoutExercise := range workoutExercises {
		workoutExercisesConverted = append(workoutExercisesConverted,
			openapi.WorkoutExerciseResponse{
				Id:         fmt.Sprintf("%d", workoutExercise.ID),
				WorkoutId:  fmt.Sprintf("%d", workoutExercise.WorkoutID),
				ExerciseId: fmt.Sprintf("%d", workoutExercise.ExerciseID),
				Sets:       int32(workoutExercise.Sets),
				Reps:       int32(workoutExercise.Reps),
			})
	}
	return openapi.Response(http.StatusOK, workoutExercisesConverted), nil
}

func (h *WorkoutExerciseHandler) PostWorkoutExercise(ctx context.Context, workoutExerciseRequest openapi.WorkoutExerciseRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)

	workoutExercise, err := h.service.AddExerciseToWorkout(workoutExerciseRequest, userID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	workoutExerciseConverted := openapi.WorkoutExerciseResponse{
		Id:         fmt.Sprintf("%d", workoutExercise.ID),
		WorkoutId:  fmt.Sprintf("%d", workoutExercise.WorkoutID),
		ExerciseId: fmt.Sprintf("%d", workoutExercise.ExerciseID),
		Sets:       int32(workoutExercise.Sets),
		Reps:       int32(workoutExercise.Reps),
	}

	return openapi.Response(http.StatusCreated, workoutExerciseConverted), nil
}

func (h *WorkoutExerciseHandler) PatchWorkoutExercise(ctx context.Context, workoutExerciseId string, workoutExerciseRequest openapi.WorkoutExerciseRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	workoutExerciseID, err := strconv.Atoi(workoutExerciseId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workoutExercise, err := h.service.UpdateWorkoutExercise(workoutExerciseRequest, userID, uint(workoutExerciseID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusOK,
		openapi.WorkoutExerciseResponse{
			Id:         fmt.Sprintf("%d", workoutExercise.ID),
			WorkoutId:  fmt.Sprintf("%d", workoutExercise.WorkoutID),
			ExerciseId: fmt.Sprintf("%d", workoutExercise.ExerciseID),
			Sets:       int32(workoutExercise.Sets),
			Reps:       int32(workoutExercise.Reps),
		}), nil
}

func (h *WorkoutExerciseHandler) DeleteWorkoutExercise(ctx context.Context, workoutExerciseId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	workoutExerciseID, err := strconv.Atoi(workoutExerciseId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}

	err = h.service.DeleteWorkoutExercise(userID, uint(workoutExerciseID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}
