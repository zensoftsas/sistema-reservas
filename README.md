# 🏥 Sistema de Reservas - Clínica Internacional

API REST para gestión de reservas médicas construida con Clean Architecture en Go.

## 📖 Descripción

Sistema backend para gestionar reservas de citas médicas, incluyendo administración de usuarios (pacientes, doctores, administradores), horarios de atención, y citas médicas.

## 🏗️ Arquitectura

El proyecto implementa **Clean Architecture** con las siguientes capas:

- **Domain Layer**: Entidades del negocio y reglas de dominio
- **Use Case Layer**: Lógica de aplicación y casos de uso
- **Repository Layer**: Interfaces de persistencia e implementaciones
- **Delivery Layer**: API HTTP REST

## 🛠️ Tecnologías

- **Go 1.24.3**: Lenguaje de programación
- **SQLite**: Base de datos embebida
- **bcrypt**: Hash seguro de contraseñas
- **UUID**: Generación de identificadores únicos
- **net/http**: Servidor HTTP estándar de Go
- **JWT (golang-jwt/jwt/v5)**: Autenticación basada en tokens
- **Middleware**: Logging, Panic Recovery, Autenticación JWT y Autorización por Roles

## 🚀 Cómo ejecutar

### Prerrequisitos
- Go 1.24.3 o superior instalado

### Instalación

1. Clonar el repositorio
```bash
git clone <repository-url>
cd version-1-0
```

2. Descargar dependencias
```bash
go mod download
```

3. Ejecutar el servidor
```bash
go run cmd/api/main.go
```

El servidor se iniciará en `http://localhost:8080`

## 📁 Estructura del Proyecto

```
version-1-0/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── domain/                  # Entidades del dominio
│   │   ├── user.go
│   │   ├── patient.go
│   │   ├── doctor.go
│   │   ├── appointment.go
│   │   └── schedule.go
│   ├── usecase/                 # Casos de uso
│   │   ├── user/
│   │   │   ├── dto.go
│   │   │   └── create_user.go
│   │   ├── appointment/
│   │   └── schedule/
│   ├── repository/              # Capa de persistencia
│   │   ├── interfaces.go
│   │   └── sqlite/
│   │       ├── connection.go
│   │       └── user_repository.go
│   └── delivery/                # Capa de entrega
│       └── http/
│           ├── router.go
│           ├── handler/
│           │   └── user_handler.go
│           └── middleware/
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── logger/
│   └── validator/
├── migrations/                  # Migraciones de BD (SQL)
├── docs/                        # Documentación
├── go.mod
├── go.sum
└── README.md
```

## 🔌 API Endpoints

### Health Check
```
GET /
```
Verifica que el servidor esté funcionando.

**Respuesta:**
```
Sistema de Reservas - API Running
```

### Crear Usuario
```
POST /api/users
```

Crea un nuevo usuario en el sistema.

**Request Body:**
```json
{
  "email": "doctor@clinica.com",
  "password": "password123",
  "first_name": "Dr. Carlos",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "doctor"
}
```

**Roles válidos:** `admin`, `doctor`, `patient`

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "doctor@clinica.com",
  "first_name": "Dr. Carlos",
  "last_name": "Pérez",
  "role": "doctor",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Errores posibles:**
- `400 Bad Request`: Datos inválidos o email ya existe
- `405 Method Not Allowed`: Método HTTP incorrecto

### Obtener Usuario por ID

```
GET /api/users?id=<uuid>
```

Obtiene la información de un usuario específico por su ID.

**Query Parameters:**
- `id` (required): UUID del usuario

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "doctor@clinica.com",
  "first_name": "Dr. Carlos",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "doctor",
  "is_active": true,
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Errores posibles:**
- `400 Bad Request`: ID no proporcionado o inválido
- `404 Not Found`: Usuario no encontrado
- `405 Method Not Allowed`: Método HTTP incorrecto

### Login (Autenticación)

```
POST /api/auth/login
```

Autentica un usuario y devuelve un token JWT para acceder a endpoints protegidos.

**Request Body:**
```json
{
  "email": "doctor@clinica.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-01-16T10:30:00Z",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "doctor@clinica.com",
    "first_name": "Dr. Carlos",
    "last_name": "Pérez",
    "phone": "+51987654321",
    "role": "doctor",
    "is_active": true,
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

**Errores posibles:**
- `400 Bad Request`: Email o password faltantes
- `401 Unauthorized`: Credenciales inválidas o usuario inactivo
- `405 Method Not Allowed`: Método HTTP incorrecto

### Obtener Perfil del Usuario Autenticado

```
GET /api/users/me
```

Obtiene la información del usuario autenticado. **Requiere autenticación JWT.**

**Headers requeridos:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "doctor@clinica.com",
  "first_name": "Dr. Carlos",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "doctor",
  "is_active": true,
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Errores posibles:**
- `401 Unauthorized`: Token inválido, expirado o no proporcionado
- `404 Not Found`: Usuario no encontrado
- `405 Method Not Allowed`: Método HTTP incorrecto

### Listar Usuarios (Solo Admin)

```
GET /api/users/list
```

Obtiene una lista paginada de todos los usuarios del sistema. **Requiere autenticación JWT y rol de administrador.**

**Headers requeridos:**
```
Authorization: Bearer <token>
```

**Query parameters (opcionales):**
- `limit` (int): Número máximo de usuarios a retornar (default: 20, máximo: 100)
- `offset` (int): Número de usuarios a saltar para paginación (default: 0)

**Ejemplo:**
```
GET /api/users/list?limit=10&offset=0
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "admin@clinica.com",
      "first_name": "Admin",
      "last_name": "Sistema",
      "phone": "+51999999999",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-01-15T10:00:00Z"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "email": "doctor@clinica.com",
      "first_name": "Dr. Carlos",
      "last_name": "Pérez",
      "phone": "+51987654321",
      "role": "doctor",
      "is_active": true,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "total": 2,
  "limit": 20,
  "offset": 0,
  "has_more": false
}
```

**Errores posibles:**
- `401 Unauthorized`: Token inválido, expirado o no proporcionado
- `403 Forbidden`: Usuario no tiene rol de administrador
- `405 Method Not Allowed`: Método HTTP incorrecto
- `500 Internal Server Error`: Error del servidor

## 📊 Modelo de Datos

### User (Usuario)
- `id` (UUID): Identificador único
- `email` (string): Email único
- `password_hash` (string): Hash del password
- `first_name` (string): Nombre
- `last_name` (string): Apellido
- `phone` (string): Teléfono
- `role` (enum): admin | doctor | patient
- `is_active` (bool): Estado activo/inactivo
- `created_at`, `updated_at` (timestamp)

### Patient (Paciente)
- `id` (UUID): Identificador único
- `user_id` (UUID): Referencia al usuario
- `birthdate` (date): Fecha de nacimiento
- `document_type` (string): Tipo de documento
- `document_number` (string): Número de documento
- `address` (string): Dirección
- `emergency_contact_name` (string): Nombre contacto emergencia
- `emergency_contact_phone` (string): Teléfono emergencia
- `blood_type` (string, opcional): Tipo de sangre
- `allergies` ([]string, opcional): Alergias
- `created_at`, `updated_at` (timestamp)

### Doctor (Doctor)
- `id` (UUID): Identificador único
- `user_id` (UUID): Referencia al usuario
- `specialty` (string): Especialidad médica
- `license_number` (string): Número de licencia
- `years_of_experience` (int): Años de experiencia
- `education` (string, opcional): Educación
- `bio` (string, opcional): Biografía
- `consultation_fee` (float64): Tarifa consulta
- `is_available` (bool): Disponibilidad
- `created_at`, `updated_at` (timestamp)

### Appointment (Cita)
- `id` (UUID): Identificador único
- `patient_id` (UUID): Referencia al paciente
- `doctor_id` (UUID): Referencia al doctor
- `scheduled_at` (timestamp): Fecha/hora programada
- `duration` (int): Duración en minutos
- `reason` (string): Motivo de consulta
- `notes` (string): Notas adicionales
- `status` (enum): pending | confirmed | cancelled | completed
- `cancelled_at` (timestamp, opcional): Fecha de cancelación
- `cancellation_reason` (string, opcional): Motivo cancelación
- `created_at`, `updated_at` (timestamp)

### Schedule (Horario)
- `id` (UUID): Identificador único
- `doctor_id` (UUID): Referencia al doctor
- `day_of_week` (enum): monday | tuesday | wednesday | thursday | friday | saturday | sunday
- `start_time` (string): Hora inicio (formato "HH:MM")
- `end_time` (string): Hora fin (formato "HH:MM")
- `slot_duration` (int): Duración de slots en minutos
- `is_active` (bool): Estado activo/inactivo
- `created_at`, `updated_at` (timestamp)

## 🧪 Ejemplos de uso con curl

### Health Check
```bash
curl http://localhost:8080/
```

### Crear un Doctor
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@clinica.com",
    "password": "password123",
    "first_name": "Dr. Carlos",
    "last_name": "Pérez",
    "phone": "+51987654321",
    "role": "doctor"
  }'
```

### Crear un Paciente
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "paciente@example.com",
    "password": "password123",
    "first_name": "María",
    "last_name": "García",
    "phone": "+51912345678",
    "role": "patient"
  }'
```

### Crear un Administrador
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@clinica.com",
    "password": "admin123456",
    "first_name": "Juan",
    "last_name": "Rodríguez",
    "phone": "+51998765432",
    "role": "admin"
  }'
```

### Obtener Usuario por ID
```bash
# Reemplaza <user-id> con el UUID del usuario
curl http://localhost:8080/api/users?id=<user-id>

# Ejemplo con un UUID real:
curl http://localhost:8080/api/users?id=550e8400-e29b-41d4-a716-446655440000
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@clinica.com",
    "password": "password123"
  }'
```

### Obtener Perfil del Usuario Autenticado
```bash
# Primero haz login para obtener el token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"doctor@clinica.com","password":"password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Luego usa el token para acceder al perfil
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"

# O directamente con el token:
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Listar Usuarios (Solo Admin)
```bash
# Primero haz login con una cuenta de administrador
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@clinica.com","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Listar usuarios (sin parámetros - usa defaults: limit=20, offset=0)
curl http://localhost:8080/api/users/list \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Listar usuarios con paginación personalizada
curl "http://localhost:8080/api/users/list?limit=10&offset=0" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Segunda página (siguiente 10 usuarios)
curl "http://localhost:8080/api/users/list?limit=10&offset=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🔐 Seguridad

- **Autenticación JWT**: Tokens con expiración de 24 horas
- **Protección de endpoints**: Middleware de autenticación para rutas protegidas
- **Hashing de contraseñas**: bcrypt con costo 10
- **Validación de credenciales**: Verificación de email, password y estado del usuario
- Los passwords nunca se retornan en las respuestas JSON
- Validación de emails únicos a nivel de base de datos y aplicación
- Validación de roles permitidos
- **Header de autorización**: `Authorization: Bearer <token>`

## 📝 Validaciones

### Usuario
- Email: requerido y único
- Password: mínimo 8 caracteres
- FirstName, LastName: requeridos
- Phone: requerido
- Role: debe ser uno de: admin, doctor, patient

### Citas
- No se pueden agendar en el pasado
- Duración debe ser mayor a 0
- Cancelación requiere 24 horas de anticipación
- Solo se pueden confirmar citas en estado "pending"
- Solo se pueden completar citas en estado "confirmed"

### Horarios
- StartTime debe ser antes de EndTime
- Formato de tiempo: "HH:MM" (00:00 - 23:59)
- SlotDuration debe ser mayor a 0

## 🗄️ Base de Datos

El proyecto utiliza **SQLite** como base de datos embebida:
- Archivo: `clinica.db` (se crea automáticamente)
- Las migraciones se ejecutan al iniciar el servidor
- Soporte para FOREIGN KEY con CASCADE DELETE
- Índices optimizados para búsquedas frecuentes

## 🚧 Estado del Proyecto

**Versión actual:** 1.0.0-alpha

### ✅ Implementado
- [x] Clean Architecture con capas separadas
- [x] Entidades de dominio (User, Patient, Doctor, Appointment, Schedule)
- [x] Repositorio SQLite para Users
- [x] Use Case: Crear Usuario
- [x] Use Case: Obtener Usuario por ID
- [x] Use Case: Login con JWT
- [x] Use Case: Listar Usuarios con paginación
- [x] API REST: Endpoint POST /api/users (Crear usuario)
- [x] API REST: Endpoint GET /api/users?id=<uuid> (Obtener usuario)
- [x] API REST: Endpoint POST /api/auth/login (Autenticación)
- [x] API REST: Endpoint GET /api/users/me (Perfil autenticado)
- [x] API REST: Endpoint GET /api/users/list (Listar usuarios - solo admin)
- [x] Validaciones de negocio
- [x] Health check endpoint
- [x] Migraciones automáticas de base de datos
- [x] Middleware de logging (registra todas las requests)
- [x] Middleware de panic recovery (previene crashes del servidor)
- [x] Middleware de autenticación JWT (protege endpoints privados)
- [x] Middleware de autorización por roles (protege endpoints por rol)
- [x] Sistema de autenticación con tokens JWT (expiración 24h)
- [x] Sistema de paginación para listados

### 🔜 Pendiente
- [ ] Más endpoints CRUD (UPDATE, DELETE para usuarios)
- [ ] Gestión completa de pacientes
- [ ] Gestión completa de doctores
- [ ] Gestión de citas médicas
- [ ] Gestión de horarios
- [ ] Roles adicionales (doctor, patient) con permisos específicos
- [ ] Refresh tokens
- [ ] Configuración externa (JWT secret, database path, port)
- [ ] Tests unitarios
- [ ] Tests de integración
- [ ] Documentación OpenAPI/Swagger

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es propiedad de Clínica Internacional.

## 👥 Autores

- **Equipo de Desarrollo** - Clínica Internacional

## 📞 Contacto

Para preguntas o soporte, contactar al equipo de desarrollo de Clínica Internacional.
