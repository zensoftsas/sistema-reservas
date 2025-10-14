package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"version-1-0/internal/usecase/schedule"
)

// ScheduleHandler handles HTTP requests for schedules
type ScheduleHandler struct {
	createScheduleUC *schedule.CreateScheduleUseCase
	getSchedulesUC   *schedule.GetDoctorSchedulesUseCase
	deleteScheduleUC *schedule.DeleteScheduleUseCase
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(
	createScheduleUC *schedule.CreateScheduleUseCase,
	getSchedulesUC *schedule.GetDoctorSchedulesUseCase,
	deleteScheduleUC *schedule.DeleteScheduleUseCase,
) *ScheduleHandler {
	return &ScheduleHandler{
		createScheduleUC: createScheduleUC,
		getSchedulesUC:   getSchedulesUC,
		deleteScheduleUC: deleteScheduleUC,
	}
}

// CreateScheduleRequest represents the request body for creating a schedule
type CreateScheduleRequest struct {
	DoctorID     string `json:"doctor_id"`
	DayOfWeek    string `json:"day_of_week"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	SlotDuration int    `json:"slot_duration"`
}

// CreateSchedule handles POST /api/schedules (admin only)
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.DoctorID == "" {
		http.Error(w, "doctor_id is required", http.StatusBadRequest)
		return
	}
	if req.DayOfWeek == "" {
		http.Error(w, "day_of_week is required", http.StatusBadRequest)
		return
	}
	if req.StartTime == "" {
		http.Error(w, "start_time is required", http.StatusBadRequest)
		return
	}
	if req.EndTime == "" {
		http.Error(w, "end_time is required", http.StatusBadRequest)
		return
	}
	if req.SlotDuration <= 0 {
		http.Error(w, "slot_duration must be greater than 0", http.StatusBadRequest)
		return
	}

	// Create schedule
	ctx := context.Background()
	schedule, err := h.createScheduleUC.Execute(
		ctx,
		req.DoctorID,
		req.DayOfWeek,
		req.StartTime,
		req.EndTime,
		req.SlotDuration,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

// GetDoctorSchedules handles GET /api/schedules/doctor/{id} (public)
func (h *ScheduleHandler) GetDoctorSchedules(w http.ResponseWriter, r *http.Request) {
	// Get doctor ID from URL path
	doctorID := r.PathValue("id")
	if doctorID == "" {
		http.Error(w, "Doctor ID is required", http.StatusBadRequest)
		return
	}

	// Get schedules
	ctx := context.Background()
	schedules, err := h.getSchedulesUC.Execute(ctx, doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedules)
}

// DeleteSchedule handles DELETE /api/schedules/{id} (admin only)
func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Get schedule ID from URL path
	scheduleID := r.PathValue("id")
	if scheduleID == "" {
		http.Error(w, "Schedule ID is required", http.StatusBadRequest)
		return
	}

	// Delete schedule
	ctx := context.Background()
	if err := h.deleteScheduleUC.Execute(ctx, scheduleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Schedule deleted successfully",
	})
}
