package main

import (
	"fmt"
	"log"
	"net/http"

	httpDelivery "version-1-0/internal/delivery/http"
	"version-1-0/internal/delivery/http/handler"
	"version-1-0/internal/repository/sqlite"
	"version-1-0/internal/usecase/auth"
	"version-1-0/internal/usecase/user"
)

func main() {
	fmt.Println("🏥 Sistema de Reservas - Clínica Internacional - API Server")
	fmt.Println("=============================================================\n")

	// Initialize SQLite database
	fmt.Println("📦 Inicializando base de datos...")
	db, err := sqlite.InitDB("./clinica.db")
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

	// Create use cases
	createUserUC := user.NewCreateUserUseCase(userRepo)
	getUserUC := user.NewGetUserUseCase(userRepo)
	listUsersUC := user.NewListUsersUseCase(userRepo)

	// Create auth use cases
	loginUC := auth.NewLoginUseCase(userRepo)

	// Create handlers
	userHandler := handler.NewUserHandler(createUserUC, getUserUC, listUsersUC)
	authHandler := handler.NewAuthHandler(loginUC)

	// Configure router
	router := httpDelivery.SetupRouter(userHandler, authHandler)

	// Configure HTTP server
	port := ":8080"
	fmt.Printf("🚀 Servidor HTTP iniciado en http://localhost%s\n", port)
	fmt.Println("📍 Endpoints disponibles:")
	fmt.Println("   GET  /                      - Health check")
	fmt.Println("   POST /api/users             - Crear usuario")
	fmt.Println("   GET  /api/users?id=<uuid>   - Obtener usuario por ID")
	fmt.Println("   POST /api/auth/login        - Login (obtener token)")
	fmt.Println("   GET  /api/users/me          - Obtener perfil (requiere token)")
	fmt.Println("   GET  /api/users/list        - Listar usuarios (solo admin)")
	fmt.Println("\n⏳ Presiona Ctrl+C para detener el servidor...\n")

	// Start HTTP server
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}