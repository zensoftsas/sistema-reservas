package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"version-1-0/internal/usecase/analytics"
)

// AnalyticsHandler handles HTTP requests for analytics endpoints
type AnalyticsHandler struct {
	getDashboardSummary *analytics.GetDashboardSummaryUseCase
	getRevenueStats     *analytics.GetRevenueStatsUseCase
	getTopDoctors       *analytics.GetTopDoctorsUseCase
	getTopServices      *analytics.GetTopServicesUseCase
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(
	getDashboardSummary *analytics.GetDashboardSummaryUseCase,
	getRevenueStats *analytics.GetRevenueStatsUseCase,
	getTopDoctors *analytics.GetTopDoctorsUseCase,
	getTopServices *analytics.GetTopServicesUseCase,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		getDashboardSummary: getDashboardSummary,
		getRevenueStats:     getRevenueStats,
		getTopDoctors:       getTopDoctors,
		getTopServices:      getTopServices,
	}
}

// GetDashboardSummary godoc
// @Summary      Dashboard summary
// @Description  Retorna resumen con KPIs principales del sistema (solo admin)
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.DashboardSummary
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Router       /api/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary, err := h.getDashboardSummary.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetRevenueStats godoc
// @Summary      Estadísticas de ingresos
// @Description  Retorna ingresos agrupados por servicio (solo admin)
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.RevenueByService
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Router       /api/analytics/revenue [get]
func (h *AnalyticsHandler) GetRevenueStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.getRevenueStats.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetTopDoctors godoc
// @Summary      Top doctores más solicitados
// @Description  Retorna ranking de doctores por cantidad de citas (solo admin)
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "Número de doctores a retornar" default(10)
// @Success      200  {array}   dto.TopDoctor
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/analytics/top-doctors [get]
func (h *AnalyticsHandler) GetTopDoctors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse limit from query params
	limit := 10 // Default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	doctors, err := h.getTopDoctors.Execute(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctors)
}

// GetTopServices godoc
// @Summary      Top servicios más populares
// @Description  Retorna ranking de servicios por demanda (solo admin)
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "Número de servicios a retornar" default(10)
// @Success      200  {array}   dto.TopService
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/analytics/top-services [get]
func (h *AnalyticsHandler) GetTopServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse limit from query params
	limit := 10 // Default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	services, err := h.getTopServices.Execute(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
