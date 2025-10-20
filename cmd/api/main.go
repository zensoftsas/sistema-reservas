package main

import (
	"fmt"
	"log"
	"net/http"

	httpDelivery "version-1-0/internal/delivery/http"
	"version-1-0/internal/delivery/http/handler"
	"version-1-0/internal/repository/sqlite"
	"version-1-0/internal/usecase/analytics"
	"version-1-0/internal/usecase/appointment"
	"version-1-0/internal/usecase/auth"
	"version-1-0/internal/usecase/doctor"
	"version-1-0/internal/usecase/schedule"
	"version-1-0/internal/usecase/service"
	"version-1-0/internal/usecase/user"
	"version-1-0/pkg/email"
	"version-1-0/pkg/reminder"

	"version-1-0/pkg/config"
)

// @title           Sistema de Reservas - Clínica Internacional API
// @version         1.0
// @description     API completa para gestión de clínica médica con sistema de citas, servicios, horarios personalizados y analytics.
// @description
// @description     Características principales:
// @description     - Autenticación JWT con roles (admin, doctor, patient)
// @description     - Sistema de citas médicas con validaciones
// @description     - Servicios médicos con precios y duraciones
// @description     - Horarios personalizados por doctor
// @description     - Cálculo dinámico de slots disponibles
// @description     - Analytics y reportes para administradores
// @description     - Notificaciones por email (SendGrid)
// @description     - Recordatorios automáticos
// @description
// @termsOfService  http://swagger.io/terms/

// @contact.name   Zensoft Team
// @contact.url    https://github.com/zensoftsas/sistema-reservas
// @contact.email  zensoftsas@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Authentication
// @tag.description Endpoints de autenticación y login

// @tag.name Users
// @tag.description Gestión de usuarios (admin, doctor, patient)

// @tag.name Appointments
// @tag.description Sistema de citas médicas

// @tag.name Services
// @tag.description Gestión de servicios/consultas médicas

// @tag.name Schedules
// @tag.description Horarios personalizados de doctores

// @tag.name Analytics
// @tag.description Reportes y estadísticas (solo admin)

// @tag.name Doctors
// @tag.description Búsqueda de doctores por especialidad

func main() {
	fmt.Println("🏥 Sistema de Reservas - Clínica Internacional - API Server")
	fmt.Println("=============================================================\n")

	// Load configuration
	cfg := config.LoadConfig()
	fmt.Printf("🔧 Configuración cargada:\n")
	fmt.Printf("   Puerto: %s\n", cfg.ServerPort)
	fmt.Printf("   Base de datos: %s\n", cfg.DatabaseURL)
	fmt.Printf("   JWT Expiration: %d horas\n\n", cfg.JWTExpirationHrs)

	// Initialize PostgreSQL database
	fmt.Println("📦 Inicializando base de datos...")
	db, err := sqlite.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer db.Close()
	fmt.Println("⏰ Servicio de recordatorios iniciado")

	// Create repositories
	userRepo := sqlite.NewSqliteUserRepository(db)
	doctorRepo := sqlite.NewSqliteDoctorRepository(db)
	patientRepo := sqlite.NewSqlitePatientRepository(db)
	appointmentRepo := sqlite.NewSqliteAppointmentRepository(db)
	serviceRepo := sqlite.NewSqliteServiceRepository(db)
	doctorServiceRepo := sqlite.NewSqliteDoctorServiceRepository(db)
	scheduleRepo := sqlite.NewSqliteScheduleRepository(db)

	// Create email service
	emailService := email.NewEmailService(
		cfg.SendGridAPIKey,
		cfg.SendGridFromEmail,
		cfg.SendGridFromName,
	)

	// Create reminder service
	reminderService := reminder.NewReminderService(appointmentRepo, userRepo, emailService)

	// Start reminder scheduler in background
	reminderService.Start()

	// Create use cases
	createUserUC := user.NewCreateUserUseCase(userRepo, doctorRepo, patientRepo)
	getUserUC := user.NewGetUserUseCase(userRepo)
	listUsersUC := user.NewListUsersUseCase(userRepo)
	updateUserUC := user.NewUpdateUserUseCase(userRepo)
	deleteUserUC := user.NewDeleteUserUseCase(userRepo)

	// Create appointment use cases
	createAppointmentUC := appointment.NewCreateAppointmentUseCase(appointmentRepo, userRepo, serviceRepo, doctorServiceRepo, emailService)
	getByPatientUC := appointment.NewGetAppointmentsByPatientUseCase(appointmentRepo, userRepo)
	getByDoctorUC := appointment.NewGetAppointmentsByDoctorUseCase(appointmentRepo, userRepo)
	cancelAppointmentUC := appointment.NewCancelAppointmentUseCase(appointmentRepo, userRepo, emailService)
	confirmAppointmentUC := appointment.NewConfirmAppointmentUseCase(appointmentRepo, userRepo, emailService)
	completeAppointmentUC := appointment.NewCompleteAppointmentUseCase(appointmentRepo, userRepo, emailService)
	getHistoryUC := appointment.NewGetPatientHistoryUseCase(appointmentRepo, userRepo)
	
	// Create doctor use cases
	searchDoctorsUC := doctor.NewSearchDoctorsUseCase(userRepo)

	// Create service use cases
	createServiceUC := service.NewCreateServiceUseCase(serviceRepo)
	listServicesUC := service.NewListServicesUseCase(serviceRepo)
	assignServiceToDoctorUC := service.NewAssignServiceToDoctorUseCase(doctorServiceRepo, serviceRepo, userRepo)
	getDoctorsByServiceUC := service.NewGetDoctorsByServiceUseCase(doctorServiceRepo, serviceRepo)
	getAvailableSlotsUC := service.NewGetAvailableSlotsUseCase(serviceRepo, appointmentRepo, userRepo, scheduleRepo)

	// Create auth use cases
	loginUC := auth.NewLoginUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpirationHrs)

	// Create schedule use cases
	createScheduleUC := schedule.NewCreateScheduleUseCase(scheduleRepo, userRepo)
	getSchedulesUC := schedule.NewGetDoctorSchedulesUseCase(scheduleRepo, userRepo)
	deleteScheduleUC := schedule.NewDeleteScheduleUseCase(scheduleRepo)

	// Create analytics use cases
	getDashboardSummaryUC := analytics.NewGetDashboardSummaryUseCase(appointmentRepo, userRepo)
	getRevenueStatsUC := analytics.NewGetRevenueStatsUseCase(appointmentRepo)
	getTopDoctorsUC := analytics.NewGetTopDoctorsUseCase(appointmentRepo, userRepo)
	getTopServicesUC := analytics.NewGetTopServicesUseCase(appointmentRepo)

	// Create handlers
	userHandler := handler.NewUserHandler(createUserUC, getUserUC, listUsersUC, updateUserUC, deleteUserUC)
	authHandler := handler.NewAuthHandler(loginUC)
	appointmentHandler := handler.NewAppointmentHandler(createAppointmentUC, getByPatientUC, getByDoctorUC, cancelAppointmentUC, confirmAppointmentUC, completeAppointmentUC, getHistoryUC)
	doctorHandler := handler.NewDoctorHandler(searchDoctorsUC)
	serviceHandler := handler.NewServiceHandler(createServiceUC, listServicesUC, assignServiceToDoctorUC, getDoctorsByServiceUC, getAvailableSlotsUC)
	scheduleHandler := handler.NewScheduleHandler(createScheduleUC, getSchedulesUC, deleteScheduleUC)
	analyticsHandler := handler.NewAnalyticsHandler(getDashboardSummaryUC, getRevenueStatsUC, getTopDoctorsUC, getTopServicesUC)

	// Configure router
	router := httpDelivery.SetupRouter(userHandler, authHandler, appointmentHandler, doctorHandler, serviceHandler, scheduleHandler, analyticsHandler, cfg.JWTSecret)

	// Configure HTTP server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Servidor HTTP iniciado en http://localhost:%s\n", cfg.ServerPort)
	fmt.Println("📍 Endpoints disponibles:")
	fmt.Println("   GET  /                      - Health check")
	fmt.Println("   GET  /swagger/              - Swagger API Documentation")
	fmt.Println("   POST /api/users             - Crear usuario")
	fmt.Println("   GET  /api/users?id=<uuid>   - Obtener usuario por ID")
	fmt.Println("   POST /api/auth/login        - Login (obtener token)")
	fmt.Println("   GET  /api/users/me          - Obtener perfil (requiere token)")
	fmt.Println("   GET  /api/users/list        - Listar usuarios (solo admin)")
	fmt.Println("   PUT    /api/users/{id}        - Actualizar usuario (admin o mismo user)")
	fmt.Println("   DELETE /api/users/delete?id=    - Eliminar usuario (solo admin)")
	fmt.Println("   POST   /api/appointments         - Crear cita (autenticado)")
	fmt.Println("   GET    /api/appointments/my      - Mis citas (autenticado)")
	fmt.Println("   GET    /api/appointments/doctor  - Citas doctor (solo doctor)")
	fmt.Println("   PUT    /api/appointments/cancel  - Cancelar cita (autenticado)")
	fmt.Println("   GET    /api/doctors/search?specialty= - Buscar doctores (público)")
	fmt.Println("   PUT    /api/appointments/confirm?id= - Confirmar cita (doctor/admin)")
	fmt.Println("   PUT    /api/appointments/complete?id= - Completar cita (doctor/admin)")
	fmt.Println("   GET    /api/appointments/history?patient_id= - Historial médico (autenticado)")
	fmt.Println("   POST   /api/services/create      - Crear servicio (solo admin)")
	fmt.Println("   GET    /api/services             - Listar servicios activos (público)")
	fmt.Println("   POST   /api/services/assign      - Asignar servicio a doctor (solo admin)")
	fmt.Println("   GET    /api/services/doctors?service_id= - Obtener doctores por servicio (público)")
	fmt.Println("   GET    /api/services/available-slots?doctor_id=&service_id=&date= - Obtener slots disponibles (público)")
	fmt.Println("   POST   /api/schedules            - Crear horario (admin)")
	fmt.Println("   GET    /api/schedules/doctor/{id} - Ver horarios de doctor (público)")
	fmt.Println("   DELETE /api/schedules/{id}       - Eliminar horario (admin)")
	fmt.Println("   GET    /api/analytics/dashboard  - Resumen del dashboard (solo admin)")
	fmt.Println("   GET    /api/analytics/revenue    - Estadísticas de ingresos (solo admin)")
	fmt.Println("   GET    /api/analytics/top-doctors?limit=10 - Top doctores (solo admin)")
	fmt.Println("   GET    /api/analytics/top-services?limit=10 - Top servicios (solo admin)")
	fmt.Println("\n⏳ Presiona Ctrl+C para detener el servidor...\n")

	// Start HTTP server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}