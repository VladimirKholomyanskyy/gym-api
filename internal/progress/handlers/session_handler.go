package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	usecase "github.com/VladimirKholomyanskyy/gym-api/internal/progress/usecase"
)

type workoutSessionHandler struct {
	useCase usecase.WorkoutSessionUseCase
}

// NewWorkoutSessionsAPIService creates a default api service
func NewWorkoutSessionHandler(useCase usecase.WorkoutSessionUseCase) openapi.WorkoutSessionsAPIServicer {
	return &workoutSessionHandler{useCase: useCase}
}

func (h *workoutSessionHandler) GetWorkoutSession(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	session, err := h.useCase.GetByID(ctx, profileId, workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	snapshot, err := common.ConverWorkoutSnapshot(&session.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	if session.CompletedAt != nil {
		return openapi.Response(http.StatusOK, openapi.WorkoutSession{
			Id:              session.ID,
			StartedAt:       session.StartedAt,
			CompletedAt:     session.CompletedAt,
			WorkoutSnapshot: *snapshot,
		}), nil
	}
	return openapi.Response(http.StatusOK, openapi.WorkoutSession{
		Id:              session.ID,
		StartedAt:       session.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}

func (h *workoutSessionHandler) ListWorkoutSessions(ctx context.Context, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	sessions, totalCount, err := h.useCase.List(ctx, profileId, int(page), int(pageSize))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var convertedSessions []openapi.WorkoutSession
	for _, session := range sessions {
		snapshot, err := common.ConverWorkoutSnapshot(&session.Snapshot)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), nil
		}
		if session.CompletedAt != nil {
			convertedSessions = append(convertedSessions, openapi.WorkoutSession{
				Id:              session.ID,
				StartedAt:       session.StartedAt,
				CompletedAt:     session.CompletedAt,
				WorkoutSnapshot: *snapshot})

		} else {
			convertedSessions = append(convertedSessions, openapi.WorkoutSession{
				Id:              session.ID,
				StartedAt:       session.StartedAt,
				WorkoutSnapshot: *snapshot,
			})
		}

	}
	return openapi.Response(
		http.StatusOK,
		openapi.ListWorkoutSessions200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       convertedSessions}), nil

}

func (s *workoutSessionHandler) AddWorkoutSession(ctx context.Context, startWorkoutSessionRequest openapi.CreateWorkoutSessionRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	workoutSession, err := s.useCase.StartWorkout(ctx, profileId, startWorkoutSessionRequest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	snapshot, err := common.ConverWorkoutSnapshot(&workoutSession.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.WorkoutSession{
		Id:              workoutSession.ID,
		StartedAt:       workoutSession.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}

// FinishWorkoutSession - Mark a workout session as completed
func (h *workoutSessionHandler) CompleteWorkoutSession(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}

	workoutSession, err := h.useCase.CompleteWorkout(ctx, profileId, workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	snapshot, err := common.ConverWorkoutSnapshot(&workoutSession.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.WorkoutSession{
		Id:              workoutSession.ID,
		StartedAt:       workoutSession.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}
