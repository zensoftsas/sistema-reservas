package config

// Config representa la configuración de la aplicación
// Contiene todos los parámetros necesarios para inicializar y ejecutar el sistema
type Config struct {
	// Port es el puerto en el que el servidor escuchará las peticiones HTTP
	// Ejemplo: 8080, 3000
	Port int

	// DatabaseURL es la cadena de conexión a la base de datos
	// Formato: "postgres://usuario:contraseña@host:puerto/nombre_bd"
	DatabaseURL string

	// Debug indica si la aplicación se ejecuta en modo debug
	// Cuando es true, se mostrarán logs adicionales y mensajes de depuración
	Debug bool
}
