package http

import (
	"net/http"

	"version-1-0/internal/delivery/http/handler"
	"version-1-0/internal/delivery/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "version-1-0/docs"
)

// SetupRouter configures and returns the HTTP router with all application routes
func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, appointmentHandler *handler.AppointmentHandler, doctorHandler *handler.DoctorHandler, serviceHandler *handler.ServiceHandler, scheduleHandler *handler.ScheduleHandler, analyticsHandler *handler.AnalyticsHandler, jwtSecret string, allowedOrigins string) http.Handler {
	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Register user routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.Create(w, r)
		} else if r.Method == http.MethodGet {
			userHandler.GetByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register authentication routes
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// Register protected user routes
	protectedUserRoutes := http.HandlerFunc(userHandler.GetMe)
	protectedUserRoutesWithAuth := middleware.AuthMiddleware(jwtSecret)(protectedUserRoutes)
	mux.Handle("/api/users/me", protectedUserRoutesWithAuth)

	// Register admin-only routes
	// List users - requires admin role
	listUsersHandler := http.HandlerFunc(userHandler.List)
	listUsersWithAuth := middleware.AuthMiddleware(jwtSecret)(listUsersHandler)
	listUsersWithRole := middleware.RequireRole("admin")(listUsersWithAuth)
	mux.Handle("/api/users/list", listUsersWithRole)

	// Update user - requires authentication (admin can update anyone, users can update themselves)
	updateUserHandler := http.HandlerFunc(userHandler.Update)
	updateUserWithAuth := middleware.AuthMiddleware(jwtSecret)(updateUserHandler)
	mux.Handle("/api/users/", updateUserWithAuth)

	// Delete user endpoint will use query param: /api/users/delete?id=xxx
	deleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/users/delete" {
			userHandler.Delete(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	deleteWithRole := middleware.RequireRole("admin")(deleteHandler)
	deleteWithAuth := middleware.AuthMiddleware(jwtSecret)(deleteWithRole)
	mux.Handle("/api/users/delete", deleteWithAuth)

	// Appointment routes - require authentication
	// Create appointment - POST /api/appointments
	createAppointmentHandler := http.HandlerFunc(appointmentHandler.Create)
	createAppointmentWithAuth := middleware.AuthMiddleware(jwtSecret)(createAppointmentHandler)
	mux.Handle("/api/appointments", createAppointmentWithAuth)

	// Get my appointments - GET /api/appointments/my
	getMyAppointmentsHandler := http.HandlerFunc(appointmentHandler.GetMyAppointments)
	getMyAppointmentsWithAuth := middleware.AuthMiddleware(jwtSecret)(getMyAppointmentsHandler)
	mux.Handle("/api/appointments/my", getMyAppointmentsWithAuth)

	// Get doctor appointments - GET /api/appointments/doctor (requires doctor role)
	getDoctorAppointmentsHandler := http.HandlerFunc(appointmentHandler.GetDoctorAppointments)
	getDoctorAppointmentsWithRole := middleware.RequireRole("doctor")(getDoctorAppointmentsHandler)
	getDoctorAppointmentsWithAuth := middleware.AuthMiddleware(jwtSecret)(getDoctorAppointmentsWithRole)
	mux.Handle("/api/appointments/doctor", getDoctorAppointmentsWithAuth)

	// Cancel appointment - PUT /api/appointments/cancel
	cancelAppointmentHandler := http.HandlerFunc(appointmentHandler.Cancel)
	cancelAppointmentWithAuth := middleware.AuthMiddleware(jwtSecret)(cancelAppointmentHandler)
	mux.Handle("/api/appointments/cancel", cancelAppointmentWithAuth)

	// Confirm appointment - PUT /api/appointments/confirm?id=xxx (doctor or admin only)
	confirmAppointmentHandler := http.HandlerFunc(appointmentHandler.Confirm)
	confirmAppointmentWithAuth := middleware.AuthMiddleware(jwtSecret)(confirmAppointmentHandler)
	mux.Handle("/api/appointments/confirm", confirmAppointmentWithAuth)

	// Complete appointment - PUT /api/appointments/complete?id=xxx (doctor or admin only)
	completeAppointmentHandler := http.HandlerFunc(appointmentHandler.Complete)
	completeAppointmentWithAuth := middleware.AuthMiddleware(jwtSecret)(completeAppointmentHandler)
	mux.Handle("/api/appointments/complete", completeAppointmentWithAuth)

	// Get patient medical history - GET /api/appointments/history?patient_id=xxx
	getHistoryHandler := http.HandlerFunc(appointmentHandler.GetHistory)
	getHistoryWithAuth := middleware.AuthMiddleware(jwtSecret)(getHistoryHandler)
	mux.Handle("/api/appointments/history", getHistoryWithAuth)

	// Doctor routes - public search endpoint
	mux.HandleFunc("/api/doctors/search", doctorHandler.Search)

	// Service routes
	// Create service - POST /api/services (admin only)
	createServiceHandler := http.HandlerFunc(serviceHandler.Create)
	createServiceWithRole := middleware.RequireRole("admin")(createServiceHandler)
	createServiceWithAuth := middleware.AuthMiddleware(jwtSecret)(createServiceWithRole)
	mux.Handle("/api/services/create", createServiceWithAuth)

	// List services - GET /api/services (public)
	mux.HandleFunc("/api/services", serviceHandler.List)

	// Assign service to doctor - POST /api/services/assign (admin only)
	assignServiceHandler := http.HandlerFunc(serviceHandler.AssignToDoctor)
	assignServiceWithRole := middleware.RequireRole("admin")(assignServiceHandler)
	assignServiceWithAuth := middleware.AuthMiddleware(jwtSecret)(assignServiceWithRole)
	mux.Handle("/api/services/assign", assignServiceWithAuth)

	// Get doctors by service - GET /api/services/doctors?service_id=xxx (public)
	mux.HandleFunc("/api/services/doctors", serviceHandler.GetDoctorsByService)

	// Get available slots - GET /api/services/available-slots?doctor_id=xxx&service_id=yyy&date=YYYY-MM-DD (public)
	mux.HandleFunc("/api/services/available-slots", serviceHandler.GetAvailableSlots)

	// Schedule routes
	// Create schedule - POST /api/schedules (admin only)
	createScheduleHandler := http.HandlerFunc(scheduleHandler.CreateSchedule)
	createScheduleWithRole := middleware.RequireRole("admin")(createScheduleHandler)
	createScheduleWithAuth := middleware.AuthMiddleware(jwtSecret)(createScheduleWithRole)
	mux.Handle("/api/schedules", createScheduleWithAuth)

	// Get doctor schedules - GET /api/schedules/doctor/{id} (public)
	mux.HandleFunc("/api/schedules/doctor/{id}", scheduleHandler.GetDoctorSchedules)

	// Delete schedule - DELETE /api/schedules/{id} (admin only)
	deleteScheduleHandler := http.HandlerFunc(scheduleHandler.DeleteSchedule)
	deleteScheduleWithRole := middleware.RequireRole("admin")(deleteScheduleHandler)
	deleteScheduleWithAuth := middleware.AuthMiddleware(jwtSecret)(deleteScheduleWithRole)
	mux.Handle("DELETE /api/schedules/{id}", deleteScheduleWithAuth)

	// Analytics routes (admin only)
	// Dashboard summary - GET /api/analytics/dashboard
	dashboardHandler := http.HandlerFunc(analyticsHandler.GetDashboardSummary)
	dashboardWithRole := middleware.RequireRole("admin")(dashboardHandler)
	dashboardWithAuth := middleware.AuthMiddleware(jwtSecret)(dashboardWithRole)
	mux.Handle("/api/analytics/dashboard", dashboardWithAuth)

	// Revenue stats - GET /api/analytics/revenue
	revenueHandler := http.HandlerFunc(analyticsHandler.GetRevenueStats)
	revenueWithRole := middleware.RequireRole("admin")(revenueHandler)
	revenueWithAuth := middleware.AuthMiddleware(jwtSecret)(revenueWithRole)
	mux.Handle("/api/analytics/revenue", revenueWithAuth)

	// Top doctors - GET /api/analytics/top-doctors?limit=10
	topDoctorsHandler := http.HandlerFunc(analyticsHandler.GetTopDoctors)
	topDoctorsWithRole := middleware.RequireRole("admin")(topDoctorsHandler)
	topDoctorsWithAuth := middleware.AuthMiddleware(jwtSecret)(topDoctorsWithRole)
	mux.Handle("/api/analytics/top-doctors", topDoctorsWithAuth)

	// Top services - GET /api/analytics/top-services?limit=10
	topServicesHandler := http.HandlerFunc(analyticsHandler.GetTopServices)
	topServicesWithRole := middleware.RequireRole("admin")(topServicesHandler)
	topServicesWithAuth := middleware.AuthMiddleware(jwtSecret)(topServicesWithRole)
	mux.Handle("/api/analytics/top-services", topServicesWithAuth)

	// Swagger documentation endpoint
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Register health check route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sistema de Reservas - API Running"))
	})

	// Apply middlewares in order: CORS -> Recovery -> Logging -> Handlers
	// CORS middleware must be first to handle preflight requests
	withCORS := middleware.CORSMiddleware(allowedOrigins)(mux)

	// Recovery middleware wraps everything to catch panics
	withRecovery := middleware.RecoveryMiddleware(withCORS)

	// Logging middleware wraps the recovery middleware to log all requests
	withLogging := middleware.LoggingMiddleware(withRecovery)

	return withLogging
}
