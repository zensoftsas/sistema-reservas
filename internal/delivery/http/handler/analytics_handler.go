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
// @Success      200  {object}  analytics.DashboardSummary
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
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
// @Success      200  {array}   analytics.RevenueByService
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
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

// GetTopDoctors returns top doctors by appointment count
// GET /api/analytics/top-doctors?limit=10
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

// GetTopServices returns top services by usage count
// GET /api/analytics/top-services?limit=10
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
