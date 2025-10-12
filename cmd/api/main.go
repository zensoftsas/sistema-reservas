package main

import (
	"fmt"
	"log"
	"net/http"

	httpDelivery "version-1-0/internal/delivery/http"
	"version-1-0/internal/delivery/http/handler"
	"version-1-0/internal/repository/sqlite"
	"version-1-0/internal/usecase/appointment"
	"version-1-0/internal/usecase/auth"
	"version-1-0/internal/usecase/doctor"
	"version-1-0/internal/usecase/user"
	"version-1-0/pkg/email"

	"version-1-0/pkg/config"
)

func main() {
	fmt.Println("🏥 Sistema de Reservas - Clínica Internacional - API Server")
	fmt.Println("=============================================================\n")

	// Load configuration
	cfg := config.LoadConfig()
	fmt.Printf("🔧 Configuración cargada:\n")
	fmt.Printf("   Puerto: %s\n", cfg.ServerPort)
	fmt.Printf("   Base de datos: %s\n", cfg.DatabasePath)
	fmt.Printf("   JWT Expiration: %d horas\n\n", cfg.JWTExpirationHours)

	// Initialize SQLite database
	fmt.Println("📦 Inicializando base de datos...")
	db, err := sqlite.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer db.Close()
	fmt.Println("✅ Base de datos inicializada")

	// Execute database migrations
	fmt.Println("🔧 Ejecutando migraciones...")
	err = sqlite.MigrateDB(db)
	if err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}
	fmt.Println("✅ Migraciones completadas\n")

	// Create repository
	userRepo := sqlite.NewSqliteUserRepository(db)

	// Create email service
	emailService := email.NewEmailService(
	cfg.SendGridAPIKey,
	cfg.SendGridFromEmail,
	cfg.SendGridFromName,
)
	// Create appointment repository
	appointmentRepo := sqlite.NewSqliteAppointmentRepository(db)

	// Create use cases
	createUserUC := user.NewCreateUserUseCase(userRepo)
	getUserUC := user.NewGetUserUseCase(userRepo)
	listUsersUC := user.NewListUsersUseCase(userRepo)
	updateUserUC := user.NewUpdateUserUseCase(userRepo)
	deleteUserUC := user.NewDeleteUserUseCase(userRepo)

	// Create appointment use cases
	createAppointmentUC := appointment.NewCreateAppointmentUseCase(appointmentRepo, userRepo, emailService)
	getByPatientUC := appointment.NewGetAppointmentsByPatientUseCase(appointmentRepo)
	getByDoctorUC := appointment.NewGetAppointmentsByDoctorUseCase(appointmentRepo)
	cancelAppointmentUC := appointment.NewCancelAppointmentUseCase(appointmentRepo, userRepo, emailService)
	confirmAppointmentUC := appointment.NewConfirmAppointmentUseCase(appointmentRepo, userRepo, emailService)
	completeAppointmentUC := appointment.NewCompleteAppointmentUseCase(appointmentRepo, userRepo, emailService)
	getHistoryUC := appointment.NewGetPatientHistoryUseCase(appointmentRepo, userRepo)
	
	// Create doctor use cases
	searchDoctorsUC := doctor.NewSearchDoctorsUseCase(userRepo)

	// Create auth use cases
	loginUC := auth.NewLoginUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpirationHours)

	// Create handlers
	userHandler := handler.NewUserHandler(createUserUC, getUserUC, listUsersUC, updateUserUC, deleteUserUC)
	authHandler := handler.NewAuthHandler(loginUC)
	// Create appointment handler
	appointmentHandler := handler.NewAppointmentHandler(createAppointmentUC, getByPatientUC, getByDoctorUC, cancelAppointmentUC, confirmAppointmentUC, completeAppointmentUC, getHistoryUC)
	// Create doctor handler
	doctorHandler := handler.NewDoctorHandler(searchDoctorsUC)

	// Configure router
	router := httpDelivery.SetupRouter(userHandler, authHandler, appointmentHandler, doctorHandler, cfg.JWTSecret)

	// Configure HTTP server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Servidor HTTP iniciado en http://localhost:%s\n", cfg.ServerPort)
	fmt.Println("📍 Endpoints disponibles:")
	fmt.Println("   GET  /                      - Health check")
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
	fmt.Println("\n⏳ Presiona Ctrl+C para detener el servidor...\n")

	// Start HTTP server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}