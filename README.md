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
- **PostgreSQL (Neon)**: Base de datos serverless en la nube
- **pgx/v5**: Driver PostgreSQL para Go
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

### 📍 Resumen de Endpoints

**Autenticación:**
- `POST   /api/auth/login`                            - Login y obtener token JWT

**Usuarios:**
- `POST   /api/users`                                 - Crear usuario (público)
- `GET    /api/users?id=`                             - Obtener usuario por ID (público)
- `GET    /api/users/me`                              - Obtener perfil autenticado (requiere token)
- `GET    /api/users/list`                            - Listar usuarios (admin)

**Citas:**
- `POST   /api/appointments`                          - Crear cita [requiere service_id] (autenticado)
- `GET    /api/appointments/my`                       - Mis citas (autenticado)
- `GET    /api/appointments/doctor`                   - Citas del doctor (doctor)
- `PUT    /api/appointments/cancel`                   - Cancelar cita (autenticado)

**Servicios Médicos:**
- `POST   /api/services/create`                       - Crear servicio (admin)
- `GET    /api/services`                              - Listar servicios activos (público)
- `POST   /api/services/assign`                       - Asignar servicio a doctor (admin)
- `GET    /api/services/doctors?service_id=`          - Doctores que ofrecen servicio (público)
- `GET    /api/services/available-slots?doctor_id=&service_id=&date=` - Horarios disponibles (público)

**Horarios Personalizados:**
- `POST   /api/schedules`                             - Crear horario (admin)
- `GET    /api/schedules/doctor/{id}`                 - Ver horarios de doctor (público)
- `DELETE /api/schedules/{id}`                        - Eliminar horario (admin)

**Analytics & Dashboard:**
- `GET    /api/analytics/dashboard`                   - Resumen del dashboard (admin)
- `GET    /api/analytics/revenue`                     - Estadísticas de ingresos (admin)
- `GET    /api/analytics/top-doctors?limit=10`        - Top doctores (admin)
- `GET    /api/analytics/top-services?limit=10`       - Top servicios (admin)

**Total:** 29 endpoints (25 previos + 4 analytics)

---

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

## 🏥 Citas Médicas (Appointments)

### Crear Cita

```
POST /api/appointments
```

Crea una nueva cita médica. **Requiere autenticación JWT.** El paciente es identificado automáticamente desde el token.

**Headers requeridos:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "doctor_id": "550e8400-e29b-41d4-a716-446655440000",
  "appointment_date": "2025-01-20",
  "appointment_time": "10:30",
  "reason": "Consulta general y revisión de exámenes"
}
```

**Response (201 Created):**

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440003",
  "patient_id": "660e8400-e29b-41d4-a716-446655440002",
  "doctor_id": "550e8400-e29b-41d4-a716-446655440000",
  "appointment_date": "2025-01-20",
  "appointment_time": "10:30",
  "status": "pending",
  "reason": "Consulta general y revisión de exámenes",
  "created_at": "2025-01-15T12:00:00Z"
}
```

**Errores posibles:**

- `400 Bad Request`: Datos inválidos (doctor no existe, fecha/hora inválida, doctor no disponible)
- `401 Unauthorized`: Token inválido o no proporcionado

### Obtener Mis Citas (Paciente)

```
GET /api/appointments/my
```

Obtiene todas las citas del paciente autenticado. **Requiere autenticación JWT.**

**Headers requeridos:**

```
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440003",
    "patient_id": "660e8400-e29b-41d4-a716-446655440002",
    "doctor_id": "550e8400-e29b-41d4-a716-446655440000",
    "appointment_date": "2025-01-20",
    "appointment_time": "10:30",
    "status": "pending",
    "reason": "Consulta general",
    "notes": "",
    "created_at": "2025-01-15T12:00:00Z"
  }
]
```

**Errores posibles:**

- `401 Unauthorized`: Token inválido o no proporcionado
- `500 Internal Server Error`: Error del servidor

### Obtener Citas del Doctor

```
GET /api/appointments/doctor
```

Obtiene todas las citas del doctor autenticado. **Requiere autenticación JWT y rol de doctor.**

**Headers requeridos:**

```
Authorization: Bearer <token>
```

**Response (200 OK):**

```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440003",
    "patient_id": "660e8400-e29b-41d4-a716-446655440002",
    "doctor_id": "550e8400-e29b-41d4-a716-446655440000",
    "appointment_date": "2025-01-20",
    "appointment_time": "10:30",
    "status": "pending",
    "reason": "Consulta general",
    "notes": "",
    "created_at": "2025-01-15T12:00:00Z"
  }
]
```

**Errores posibles:**

- `401 Unauthorized`: Token inválido o no proporcionado
- `403 Forbidden`: Usuario no tiene rol de doctor
- `500 Internal Server Error`: Error del servidor

### Cancelar Cita

```
PUT /api/appointments/cancel
```

Cancela una cita existente. **Requiere autenticación JWT.** Solo el paciente, el doctor involucrado o un admin pueden cancelar una cita.

**Headers requeridos:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "appointment_id": "770e8400-e29b-41d4-a716-446655440003",
  "reason": "Tengo un compromiso urgente"
}
```

**Response (204 No Content)**

**Errores posibles:**

- `400 Bad Request`: Cita ya cancelada o datos inválidos
- `401 Unauthorized`: Token inválido o no proporcionado
- `403 Forbidden`: Usuario no tiene permisos para cancelar esta cita
- `404 Not Found`: Cita no encontrada

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

---

## 🎨 Frontend

El frontend de este proyecto está en un repositorio separado:

👉 **[Sistema Reservas - Frontend React](https://github.com/zensoftsas/sistema-reservas-frontend)**

---

## 📚 Documentación Adicional

- **[BACKEND_REFERENCE.md](./BACKEND_REFERENCE.md)** - Referencia rápida de endpoints y configuración
- **[README.md](./README.md)** - Documentación principal (este archivo)

---

## 🏥 Sistema de Servicios Médicos

### Descripción
El sistema permite gestionar servicios/consultas médicas con diferentes duraciones y precios. Cada servicio puede ser ofrecido por múltiples doctores (relación many-to-many), y los pacientes pueden reservar citas seleccionando el servicio deseado.

### Flujo de Reserva de Citas
```
1. Paciente selecciona SERVICIO
   GET /api/services
   → Lista de servicios disponibles (Consulta General, Cardiológica, etc.)

2. Sistema muestra DOCTORES que ofrecen ese servicio
   GET /api/services/doctors?service_id={id}
   → Lista de doctores disponibles para el servicio

3. Paciente selecciona DOCTOR y FECHA

4. Sistema calcula HORARIOS DISPONIBLES
   GET /api/services/available-slots?doctor_id={id}&service_id={id}&date=YYYY-MM-DD
   → Slots de tiempo basados en duración del servicio
   → Marcados como disponibles/ocupados

5. Paciente selecciona HORARIO y confirma
   POST /api/appointments
   → Cita creada con service_id
   → Duración automática del servicio
```

---

## 📅 Sistema de Horarios Personalizados

### Descripción
Cada doctor puede tener horarios personalizados por día de la semana. Esto permite:
- Horarios diferentes cada día
- Múltiples bloques horarios por día (ej: mañana + tarde)
- Días no laborables (sin schedules = no genera slots)
- Validación automática de conflictos

### Configuración de Horarios

Los horarios se configuran por día de la semana (monday-sunday) con:
- `start_time`: Hora de inicio (formato HH:MM)
- `end_time`: Hora de fin (formato HH:MM)
- `slot_duration`: Duración de cada slot en minutos

**Ejemplo de configuración:**
```
Dra. García - Lunes:
  - Bloque 1: 09:00-13:00 (mañana)
  - Bloque 2: 15:00-18:00 (tarde)

Dra. García - Viernes:
  - Bloque 1: 08:00-12:00

Dra. García - Miércoles:
  - Sin schedules (no trabaja)
```

### Endpoints de Horarios

#### 1. Crear Horario (Admin)
```bash
POST /api/schedules
Authorization: Bearer {admin-token}

Request:
{
  "doctor_id": "user-uuid",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "13:00",
  "slot_duration": 30
}

Response (201):
{
  "id": "schedule-uuid",
  "doctor_id": "doctor-real-id",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "13:00",
  "slot_duration": 30,
  "is_active": true,
  "created_at": "2025-10-14T13:41:47Z"
}
```

#### 2. Ver Horarios de un Doctor (Público)
```bash
GET /api/schedules/doctor/{user-id}

Response (200):
[
  {
    "id": "uuid",
    "day_of_week": "monday",
    "start_time": "09:00",
    "end_time": "13:00",
    "slot_duration": 30,
    "is_active": true
  },
  {
    "id": "uuid",
    "day_of_week": "monday",
    "start_time": "15:00",
    "end_time": "18:00",
    "slot_duration": 30,
    "is_active": true
  }
]
```

#### 3. Eliminar Horario (Admin)
```bash
DELETE /api/schedules/{schedule-id}
Authorization: Bearer {admin-token}

Response (200):
{
  "message": "Schedule deleted successfully"
}
```

### Integración con Slots Disponibles

**GetAvailableSlots ahora usa horarios reales:**
```bash
GET /api/services/available-slots?doctor_id={id}&service_id={id}&date=2025-10-20

# Si es Lunes (con 2 bloques: 09:00-13:00 y 15:00-18:00):
[
  {"time": "09:00", "available": true},
  {"time": "09:45", "available": true},
  {"time": "10:30", "available": true},
  {"time": "11:15", "available": true},
  {"time": "12:00", "available": true},
  {"time": "12:45", "available": true},
  // GAP - no genera slots 13:00-15:00
  {"time": "15:00", "available": true},
  {"time": "15:45", "available": true},
  {"time": "16:30", "available": true},
  {"time": "17:15", "available": true}
]

# Si es Miércoles (sin schedules):
[]  // No trabaja ese día
```

### Validaciones Implementadas

Al crear un horario:
- ✅ `doctor_id` debe ser un doctor válido y activo
- ✅ `day_of_week` debe ser monday-sunday
- ✅ `start_time` y `end_time` en formato HH:MM
- ✅ `start_time` debe ser antes de `end_time`
- ✅ `slot_duration` debe ser mayor que 0
- ✅ No puede solaparse con otro horario del mismo día
- ✅ Usa `doctor.id` real (no `user.id`)

### Casos de Uso
```
Caso 1: Doctor con horario partido
- Lunes: 08:00-12:00 (mañana) + 14:00-18:00 (tarde)
- Sistema genera 2 grupos de slots separados

Caso 2: Doctor con días libres
- Lunes, Martes, Jueves: Tiene schedules
- Miércoles, Viernes: Sin schedules
- GetAvailableSlots retorna [] en días sin schedule

Caso 3: Prevención de conflictos
- Intento de crear 09:00-13:00 cuando ya existe 11:00-15:00
- Sistema rechaza por overlap
```

### Arquitectura
```
Admin configura schedule
    ↓
Guarda en tabla schedules (doctor_id, day_of_week, start_time, end_time)
    ↓
Paciente consulta slots disponibles
    ↓
GetAvailableSlots:
  1. Obtiene día de la semana de la fecha (monday, tuesday, etc.)
  2. Consulta schedules del doctor para ese día
  3. Si no hay schedules → retorna []
  4. Si hay schedules → genera slots solo en esos horarios
  5. Marca slots ocupados por citas existentes
    ↓
Retorna slots con disponibilidad real
```

---

## 📊 Sistema de Analytics y Dashboard

### Descripción

El sistema de analytics proporciona estadísticas y métricas del negocio para administradores, incluyendo:
- Resumen general del dashboard con KPIs principales
- Análisis de ingresos por servicio
- Rankings de doctores y servicios más utilizados
- Tasas de cancelación y métricas de rendimiento

**Acceso:** Solo administradores (requiere rol `admin`)

### Endpoints de Analytics

#### 1. Dashboard Summary

**Obtiene resumen general con métricas clave:**

```bash
GET /api/analytics/dashboard
Authorization: Bearer {admin-token}

Response (200):
{
  "total_appointments": 150,
  "pending_appointments": 20,
  "confirmed_appointments": 45,
  "completed_appointments": 75,
  "cancelled_appointments": 10,
  "total_patients": 80,
  "total_doctors": 12,
  "total_revenue": 12500.50,
  "cancellation_rate": 6.67  // Porcentaje
}
```

**Métricas incluidas:**
- Total de citas y su distribución por estado
- Total de pacientes y doctores activos
- Ingresos totales (suma de citas completadas)
- Tasa de cancelación en porcentaje

#### 2. Revenue Stats

**Análisis de ingresos agrupados por servicio:**

```bash
GET /api/analytics/revenue
Authorization: Bearer {admin-token}

Response (200):
[
  {
    "service_id": "uuid",
    "service_name": "Consulta Cardiológica",
    "total_citas": 35,
    "revenue": 5250.00
  },
  {
    "service_id": "uuid",
    "service_name": "Consulta General",
    "total_citas": 60,
    "revenue": 4800.00
  }
]
```

**Características:**
- Solo incluye citas completadas
- Ordenado por ingresos (mayor a menor)
- Muestra nombre del servicio, cantidad de citas e ingresos totales

#### 3. Top Doctors

**Ranking de doctores por número de citas:**

```bash
GET /api/analytics/top-doctors?limit=10
Authorization: Bearer {admin-token}

Response (200):
[
  {
    "doctor_id": "doctor-uuid",
    "doctor_name": "Doctor abc123...",  // Placeholder
    "total_appointments": 85,
    "completed_appointments": 78
  },
  {
    "doctor_id": "doctor-uuid",
    "doctor_name": "Doctor def456...",
    "total_appointments": 67,
    "completed_appointments": 62
  }
]
```

**Parámetros:**
- `limit` (query, opcional): Número de doctores a retornar (default: 10)

**Nota:** El `doctor_name` actualmente usa un placeholder. En producción se haría JOIN con la tabla users.

#### 4. Top Services

**Ranking de servicios más populares:**

```bash
GET /api/analytics/top-services?limit=10
Authorization: Bearer {admin-token}

Response (200):
[
  {
    "service_id": "uuid",
    "service_name": "Consulta General",
    "total_citas": 95
  },
  {
    "service_id": "uuid",
    "service_name": "Consulta Cardiológica",
    "total_citas": 78
  }
]
```

**Parámetros:**
- `limit` (query, opcional): Número de servicios a retornar (default: 10)

**Características:**
- Incluye todas las citas (no solo completadas)
- Ordenado por cantidad de citas (mayor a menor)

### Arquitectura del Sistema de Analytics

```
HTTP Request (Admin)
    ↓
AuthMiddleware → Valida JWT
    ↓
RequireRole("admin") → Verifica rol
    ↓
AnalyticsHandler → Maneja request
    ↓
AnalyticsUseCase → Lógica de negocio
    ↓
Repository → Consultas SQL con agregaciones
    ↓
Response (JSON con métricas)
```

### Consultas SQL Utilizadas

**Dashboard Summary:**
```sql
-- Conteo por estado
SELECT COUNT(*) FROM appointments WHERE status = ?

-- Ingresos totales
SELECT COALESCE(SUM(s.price), 0)
FROM appointments a
JOIN services s ON a.service_id = s.id
WHERE a.status = 'completed'

-- Conteo de usuarios por rol
SELECT COUNT(*) FROM users WHERE role = ? AND is_active = 1
```

**Revenue by Service:**
```sql
SELECT
    a.service_id,
    s.name,
    COUNT(*) as count,
    SUM(s.price) as revenue
FROM appointments a
JOIN services s ON a.service_id = s.id
WHERE a.status = 'completed'
GROUP BY a.service_id, s.name
ORDER BY revenue DESC
```

**Top Doctors:**
```sql
SELECT
    doctor_id,
    COUNT(*) as total,
    SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed
FROM appointments
GROUP BY doctor_id
ORDER BY total DESC
LIMIT ?
```

**Top Services:**
```sql
SELECT
    a.service_id,
    s.name,
    COUNT(*) as count
FROM appointments a
JOIN services s ON a.service_id = s.id
GROUP BY a.service_id, s.name
ORDER BY count DESC
LIMIT ?
```

### Casos de Uso

**Caso 1: Dashboard administrativo**
- El administrador accede al dashboard
- Sistema muestra KPIs principales en tiempo real
- Incluye gráficos de citas por estado y métricas financieras

**Caso 2: Análisis de ingresos**
- Administrador revisa qué servicios generan más ingresos
- Identifica servicios rentables vs. subutilizados
- Toma decisiones de pricing y marketing

**Caso 3: Evaluación de desempeño**
- Administrador consulta top doctores
- Identifica doctores con mayor demanda
- Planifica horarios y recursos según demanda

**Caso 4: Optimización de servicios**
- Administrador revisa servicios más solicitados
- Ajusta oferta de servicios según demanda real
- Asigna más doctores a servicios populares

### Validaciones y Seguridad

- ✅ Solo usuarios con rol `admin` pueden acceder
- ✅ Requiere autenticación JWT válida
- ✅ Límites configurables para rankings (default: 10, evita sobrecarga)
- ✅ Queries optimizadas con agregaciones SQL
- ✅ Solo datos agregados (no expone información sensible individual)

### Testing del Sistema

```bash
# 1. Login como administrador
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@clinica.com","password":"admin123"}'

# 2. Obtener dashboard summary
curl -X GET http://localhost:8080/api/analytics/dashboard \
  -H "Authorization: Bearer {admin-token}"

# 3. Ver ingresos por servicio
curl -X GET http://localhost:8080/api/analytics/revenue \
  -H "Authorization: Bearer {admin-token}"

# 4. Ver top 5 doctores
curl -X GET "http://localhost:8080/api/analytics/top-doctors?limit=5" \
  -H "Authorization: Bearer {admin-token}"

# 5. Ver top 10 servicios
curl -X GET "http://localhost:8080/api/analytics/top-services?limit=10" \
  -H "Authorization: Bearer {admin-token}"
```

### Estructura de Código

```
internal/
├── usecase/
│   └── analytics/
│       ├── dto.go                      # Estructuras de respuesta
│       ├── get_dashboard_summary.go    # KPIs principales
│       ├── get_revenue_stats.go        # Ingresos por servicio
│       ├── get_top_doctors.go          # Ranking de doctores
│       └── get_top_services.go         # Ranking de servicios
├── delivery/
│   └── http/
│       └── handler/
│           └── analytics_handler.go    # Handlers HTTP
└── repository/
    ├── interfaces.go                   # Métodos de analytics agregados
    └── sqlite/
        ├── appointment_repository.go   # Queries de analytics
        └── user_repository.go          # Conteos por rol
```

### Mejoras Futuras

- [ ] Gráficos de tendencias (citas por mes/semana)
- [ ] Análisis de horarios pico (peak hours)
- [ ] Razones de cancelación más comunes
- [ ] Tiempo promedio de espera
- [ ] Tasa de conversión pending → confirmed
- [ ] Exportar reportes a PDF/Excel
- [ ] Filtros por fecha (últimos 7 días, mes, año)
- [ ] Comparativas período actual vs. anterior

---

## 🐘 Migración a PostgreSQL + Neon

### Base de Datos en la Nube

El sistema migró de **SQLite** (base de datos embebida) a **PostgreSQL** con **Neon** como proveedor serverless en la nube.

**Neon** es una plataforma de PostgreSQL serverless que ofrece:
- Base de datos PostgreSQL totalmente administrada
- Escalamiento automático
- Branching de bases de datos (útil para desarrollo/staging)
- Hosting en AWS con alta disponibilidad
- Tier gratuito generoso para desarrollo

### ¿Por qué PostgreSQL + Neon?

**Ventajas sobre SQLite:**
- ✅ **Production-ready**: Diseñado para aplicaciones en producción
- ✅ **Concurrencia**: Soporta múltiples conexiones simultáneas
- ✅ **Tipos nativos**: BOOLEAN, TIMESTAMP, JSON, UUID nativos
- ✅ **Escalabilidad**: Crece con tu aplicación
- ✅ **Integridad**: Constraints y transactions robustas
- ✅ **Respaldos**: Backups automáticos y point-in-time recovery
- ✅ **Serverless**: No necesitas administrar infraestructura

**Neon específicamente:**
- 🚀 **Instant setup**: Base de datos lista en segundos
- 💰 **Free tier**: 0.5 GB storage, 1 proyecto
- 🌿 **Branching**: Crea copias de BD para testing
- 📊 **Dashboard**: Monitoreo visual de queries y rendimiento
- 🔒 **Seguridad**: SSL/TLS por defecto

### Configuración

**1. Crear cuenta en Neon:**
```bash
# Visita https://neon.tech
# Crear cuenta (GitHub/Google login disponible)
# Crear nuevo proyecto
```

**2. Obtener Connection String:**
```
Dashboard → Project → Connection Details → Connection String

Formato:
postgres://[user]:[password]@[host]/[database]?sslmode=require
```

**3. Configurar variable de entorno:**
```bash
# Linux/Mac
export DATABASE_URL="postgres://user:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require"

# Windows
set DATABASE_URL=postgres://user:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require

# .env file (recomendado)
DATABASE_URL=postgres://user:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
```

**4. Instalar driver pgx:**
```bash
go get github.com/jackc/pgx/v5/stdlib
```

### Migraciones

Las migraciones se ejecutan automáticamente al iniciar el servidor:

```go
// internal/repository/sqlite/connection.go ahora es postgresql
func InitDB() (*sql.DB, error) {
    connStr := os.Getenv("DATABASE_URL")
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        return nil, err
    }

    // Ejecuta migraciones
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return db, nil
}
```

**Archivos de migración:**
```
migrations/
├── 001_create_users.sql
├── 002_create_patients.sql
├── 003_create_doctors.sql
├── 004_create_appointments.sql
├── 005_create_schedules.sql
├── 006_create_services.sql
└── 007_create_doctor_services.sql
```

### Driver PostgreSQL

**pgx/v5** es el driver recomendado para PostgreSQL en Go:

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

// Conexión
db, err := sql.Open("pgx", connectionString)
```

**Ventajas de pgx:**
- ⚡ Alto rendimiento (más rápido que lib/pq)
- 🔧 Soporte completo de PostgreSQL
- 📦 Interfaz database/sql estándar
- 🛡️ Prepared statements automáticos
- 🔄 Connection pooling integrado

### Diferencias con SQLite

**Cambios en SQL:**

| SQLite | PostgreSQL |
|--------|------------|
| `?` placeholders | `$1, $2, $3` placeholders |
| `INTEGER` para bool | `BOOLEAN` nativo |
| `DATETIME` | `TIMESTAMP` nativo |
| `AUTOINCREMENT` | `SERIAL` o `IDENTITY` |
| Tipos flexibles | Tipos estrictos |

**Ejemplo de migración de query:**

```go
// ❌ SQLite
query := `
    INSERT INTO users (id, email, is_active, created_at)
    VALUES (?, ?, ?, ?)
`
db.Exec(query, id, email, 1, time.Now()) // is_active como INTEGER

// ✅ PostgreSQL
query := `
    INSERT INTO users (id, email, is_active, created_at)
    VALUES ($1, $2, $3, $4)
`
db.Exec(query, id, email, true, time.Now()) // is_active como BOOLEAN
```

**Tipos de datos actualizados:**

```sql
-- SQLite
is_active INTEGER DEFAULT 1
created_at DATETIME

-- PostgreSQL
is_active BOOLEAN DEFAULT true
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Conexión y Pooling

**Configuración de connection pool:**

```go
func InitDB() (*sql.DB, error) {
    db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)        // Máximo 25 conexiones abiertas
    db.SetMaxIdleConns(5)         // Mantener 5 conexiones idle
    db.SetConnMaxLifetime(5 * time.Minute)  // Reciclar cada 5 min

    // Verificar conexión
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("cannot ping database: %w", err)
    }

    return db, nil
}
```

**Valores recomendados para Neon free tier:**
- `MaxOpenConns`: 20-25 (Neon limita a ~100)
- `MaxIdleConns`: 5
- `ConnMaxLifetime`: 5 minutos

### Branching (Opcional)

**Crear branch para testing:**

```bash
# Desde Neon Dashboard
1. Ir a Branches
2. Click "Create Branch"
3. Nombrar (ej: "dev", "staging")
4. Copiar nueva connection string

# Usar en desarrollo
export DATABASE_URL_DEV="postgres://...branch-dev..."
```

**Casos de uso:**
- Testing de migraciones sin afectar producción
- Desarrollo paralelo de features
- QA/Staging environment
- Rollback rápido

### Backup y Restore

**Backup automático (Neon):**
- Neon hace backups automáticos cada 24h
- Retención: 7 días (free tier)
- Point-in-time recovery disponible

**Backup manual con pg_dump:**

```bash
# Exportar toda la BD
pg_dump $DATABASE_URL > backup.sql

# Exportar solo esquema
pg_dump --schema-only $DATABASE_URL > schema.sql

# Exportar solo datos
pg_dump --data-only $DATABASE_URL > data.sql

# Restaurar
psql $DATABASE_URL < backup.sql
```

### Monitoreo

**Neon Dashboard ofrece:**
- 📊 Query performance metrics
- 🔍 Slow query log
- 💾 Storage usage
- 🔌 Active connections
- ⚡ Cache hit ratio

**Acceso:**
```
Dashboard → Your Project → Monitoring
```

**Métricas clave a monitorear:**
- Connection count (< límite de plan)
- Storage usage (< 0.5 GB en free tier)
- Query duration (identificar queries lentas)
- Error rate

### Deploy

**Variables de entorno en producción:**

```bash
# Railway / Render / Fly.io
DATABASE_URL=postgres://user:password@ep-xxx.aws.neon.tech/neondb?sslmode=require
JWT_SECRET=your-super-secret-key
PORT=8080
```

**Checklist de deploy:**
- ✅ DATABASE_URL configurado
- ✅ Migraciones probadas
- ✅ Connection pool configurado
- ✅ SSL/TLS habilitado (sslmode=require)
- ✅ Backups verificados
- ✅ Monitoreo activo

### Testing con PostgreSQL

**Opción 1: Usar Neon branch**
```bash
# Crear branch "test"
export DATABASE_URL_TEST="postgres://...branch-test..."

# Ejecutar tests
go test ./... -v
```

**Opción 2: PostgreSQL local**
```bash
# Docker
docker run --name postgres-test -e POSTGRES_PASSWORD=test -p 5432:5432 -d postgres:15

# Connection string local
export DATABASE_URL="postgres://postgres:test@localhost:5432/clinica_test?sslmode=disable"
```

**Opción 3: SQLite para tests unitarios**
```go
// Usar SQLite in-memory para tests rápidos
func setupTestDB() *sql.DB {
    db, _ := sql.Open("sqlite3", ":memory:")
    return db
}
```

### Ventajas para Producción

**Antes (SQLite):**
- ❌ Solo 1 conexión de escritura
- ❌ Archivo local (no escalable)
- ❌ Sin backups automáticos
- ❌ Limitaciones de tipos de datos
- ❌ No recomendado para producción

**Ahora (PostgreSQL + Neon):**
- ✅ Múltiples conexiones concurrentes
- ✅ Base de datos en la nube
- ✅ Backups automáticos
- ✅ Tipos de datos nativos
- ✅ Production-ready desde día 1
- ✅ Escalable horizontalmente
- ✅ SSL/TLS por defecto
- ✅ Monitoreo integrado

### Migración de Datos (SQLite → PostgreSQL)

Si tienes datos en SQLite que quieres migrar:

```bash
# 1. Exportar datos de SQLite
sqlite3 clinica.db .dump > dump.sql

# 2. Convertir a formato PostgreSQL
# (Reemplazar ? con $1, $2, etc.)
# (Convertir INTEGER a BOOLEAN donde aplique)
# (Convertir DATETIME a TIMESTAMP)

# 3. Importar a PostgreSQL
psql $DATABASE_URL < dump_converted.sql
```

**Herramientas útiles:**
- **pgloader**: Migración automática SQLite → PostgreSQL
- **DBeaver**: GUI para comparar esquemas
- **Neon CLI**: Gestión desde terminal

### Troubleshooting

**Error: "too many connections"**
```go
// Reducir MaxOpenConns
db.SetMaxOpenConns(10)
```

**Error: "connection timeout"**
```go
// Aumentar timeout
db.SetConnMaxLifetime(10 * time.Minute)
```

**Error: "SSL required"**
```bash
# Asegurar sslmode=require en connection string
DATABASE_URL="...?sslmode=require"
```

**Query lenta:**
```sql
-- Crear índices
CREATE INDEX idx_appointments_doctor ON appointments(doctor_id);
CREATE INDEX idx_appointments_date ON appointments(scheduled_at);
```

### Recursos

- **Neon Docs**: https://neon.tech/docs
- **pgx GitHub**: https://github.com/jackc/pgx
- **PostgreSQL Tutorial**: https://www.postgresql.org/docs/
- **Neon Discord**: Soporte comunitario

---

### Endpoints de Servicios

#### 1. Crear Servicio (Admin)
```bash
POST /api/services/create
Authorization: Bearer {admin-token}

Request:
{
  "name": "Consulta Cardiológica",
  "description": "Evaluación cardiovascular completa",
  "duration_minutes": 45,
  "price": 150.00
}

Response (201):
{
  "id": "uuid",
  "name": "Consulta Cardiológica",
  "duration_minutes": 45,
  "price": 150,
  "is_active": true,
  "created_at": "2025-10-14T11:15:44-05:00"
}
```

#### 2. Listar Servicios Activos (Público)
```bash
GET /api/services

Response (200):
[
  {
    "id": "uuid",
    "name": "Consulta General",
    "description": "Consulta médica general",
    "duration_minutes": 30,
    "price": 80,
    "is_active": true
  },
  ...
]
```

#### 3. Asignar Servicio a Doctor (Admin)
```bash
POST /api/services/assign
Authorization: Bearer {admin-token}

Request:
{
  "doctor_id": "user-uuid",
  "service_id": "service-uuid"
}

Response (200):
{
  "message": "Service assigned to doctor successfully"
}
```

#### 4. Ver Doctores que Ofrecen un Servicio (Público)
```bash
GET /api/services/doctors?service_id={uuid}

Response (200):
[
  {
    "id": "user-uuid",
    "email": "dr.garcia@clinica.com",
    "first_name": "Dr. Ana",
    "last_name": "García",
    "role": "doctor",
    "is_active": true
  },
  ...
]
```

#### 5. Ver Horarios Disponibles (Público) ⭐
```bash
GET /api/services/available-slots?doctor_id={user-uuid}&service_id={uuid}&date=2025-10-20

Response (200):
[
  {"time": "09:00", "available": true},
  {"time": "09:45", "available": true},
  {"time": "10:30", "available": false},  // Ocupado
  {"time": "11:15", "available": true},
  ...
]
```

### Creación de Citas con Servicios

**Endpoint actualizado:**
```bash
POST /api/appointments
Authorization: Bearer {patient-token}

Request:
{
  "doctor_id": "user-uuid",
  "service_id": "service-uuid",      // NUEVO - Obligatorio
  "appointment_date": "2025-10-20",
  "appointment_time": "10:30",
  "reason": "Consulta de seguimiento"
}

Response (201):
{
  "id": "appointment-uuid",
  "patient_id": "patient-uuid",
  "doctor_id": "doctor-real-id",
  "service_id": "service-uuid",
  "service_name": "Consulta Cardiológica",  // Nombre del servicio
  "scheduled_at": "2025-10-20T10:30:00Z",
  "duration": 45,                            // Del servicio automáticamente
  "reason": "Consulta de seguimiento",
  "status": "pending",
  "created_at": "2025-10-14T12:42:00Z"
}
```

### Validaciones Implementadas

Al crear una cita:
- ✅ `service_id` es obligatorio
- ✅ El servicio debe existir y estar activo
- ✅ El doctor debe ofrecer ese servicio (validación en BD)
- ✅ El horario debe estar disponible (sin conflictos)
- ✅ La duración se toma automáticamente del servicio
- ✅ Se usa el `doctor.id` real (no el `user.id`)

### Arquitectura de Base de Datos
```sql
-- Tabla de servicios
CREATE TABLE services (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL,
    price REAL NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Tabla de relación many-to-many
CREATE TABLE doctor_services (
    id TEXT PRIMARY KEY,
    doctor_id TEXT NOT NULL,           -- doctor.id (no user.id)
    service_id TEXT NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    FOREIGN KEY (service_id) REFERENCES services(id),
    UNIQUE(doctor_id, service_id)
);

-- Appointments actualizada
ALTER TABLE appointments ADD COLUMN service_id TEXT REFERENCES services(id);
```

### Algoritmo de Cálculo de Slots Disponibles
```
Input:
  - doctor_id (user_id del doctor)
  - service_id
  - date (YYYY-MM-DD)

Proceso:
  1. Obtener doctor.id real de la tabla doctors
  2. Obtener duración del servicio (ej: 45 minutos)
  3. Generar slots de 9:00 AM a 5:00 PM con esa duración
  4. Consultar citas existentes del doctor en esa fecha
  5. Para cada slot:
     - Verificar si hay overlap con citas existentes
     - Marcar como available=true/false
  6. Retornar lista de slots con disponibilidad

Output:
  Lista de TimeSlots con hora y estado de disponibilidad
```

### Casos de Uso Implementados
```go
// Servicios
- CreateServiceUseCase
- ListServicesUseCase
- AssignServiceToDoctorUseCase
- GetDoctorsByServiceUseCase
- GetAvailableSlotsUseCase  // Más complejo

// Citas (actualizado)
- CreateAppointmentUseCase  // Ahora incluye service_id
```

### Testing del Sistema
```bash
# 1. Login como admin
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@clinica.com","password":"admin123"}'

# 2. Crear servicio
curl -X POST http://localhost:8080/api/services/create \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Consulta General",
    "duration_minutes": 30,
    "price": 80
  }'

# 3. Asignar servicio a doctor
curl -X POST http://localhost:8080/api/services/assign \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "doctor_id": "{user-uuid}",
    "service_id": "{service-uuid}"
  }'

# 4. Login como paciente
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"paciente@email.com","password":"password"}'

# 5. Ver servicios disponibles
curl http://localhost:8080/api/services

# 6. Ver doctores que ofrecen un servicio
curl "http://localhost:8080/api/services/doctors?service_id={uuid}"

# 7. Ver horarios disponibles
curl "http://localhost:8080/api/services/available-slots?doctor_id={uuid}&service_id={uuid}&date=2025-10-20"

# 8. Crear cita
curl -X POST http://localhost:8080/api/appointments \
  -H "Authorization: Bearer {patient-token}" \
  -H "Content-Type: application/json" \
  -d '{
    "doctor_id": "{uuid}",
    "service_id": "{uuid}",
    "appointment_date": "2025-10-20",
    "appointment_time": "10:30",
    "reason": "Consulta"
  }'

# 9. Verificar que slot ahora está ocupado
curl "http://localhost:8080/api/services/available-slots?doctor_id={uuid}&service_id={uuid}&date=2025-10-20"
```

---

## 🏗️ Decisiones de Arquitectura

### Sesiones de Desarrollo:
- **Sesión 1-2:** Setup inicial + CRUD Users + Auth JWT
- **Sesión 3:** UPDATE/DELETE Users + Validaciones
- **Sesión 4:** Sistema de Citas (Create, Get, Cancel)
- **Sesión 5:** Confirm/Complete Citas + Historial Médico + Búsqueda Doctores
- **Sesión 6:** Notificaciones Email (SendGrid) + Recordatorios Automáticos
- **Sesión 7:** Sistema de Servicios + Slots Disponibles + Integración Completa
- **Sesión 8:** Horarios Personalizados + Analytics Dashboard + Bug Fix Service ID
- **Sesión 9:** Migración a PostgreSQL + Neon + Production-Ready

### Tecnologías Elegidas:
- **Go 1.21+** - Performance y concurrencia
- **PostgreSQL (Neon)** - Base de datos serverless production-ready
- **pgx/v5** - Driver PostgreSQL de alto rendimiento
- **Clean Architecture** - Mantenibilidad y testabilidad
- **JWT** - Autenticación stateless
- **RBAC** - Control de acceso basado en roles
- **SendGrid** - Emails transaccionales
- **Goroutines** - Background workers (recordatorios)

### Problemas Resueltos:
- Búsqueda con acentos → Primeras 6 letras
- Domain modeling → ScheduledAt en lugar de date/time separados
- Email anti-spam → Single Sender Verification
- Recordatorios duplicados → Campos reminder_24h_sent, reminder_1h_sent
- **Sistema de servicios con many-to-many:** Relación doctor_services usando doctor.id real (no user.id)
- **Cálculo de slots disponibles:** Algoritmo que genera slots basados en duración del servicio y detecta conflictos
- **Validación de asignaciones:** Doctor debe ofrecer el servicio antes de crear cita
- **Duración automática:** La duración de la cita se obtiene del servicio, no es manual
- **Horarios personalizados por doctor:** Sistema de schedules con soporte para múltiples bloques horarios por día
- **Detección de overlap:** Validación automática de conflictos de horarios
- **Días no laborables:** GetAvailableSlots retorna [] si no hay schedules para ese día
- **Slots dinámicos:** Genera slots solo en horarios configurados, no en horario fijo
- **Analytics con SQL agregaciones:** Queries optimizadas con COUNT, SUM, GROUP BY, JOIN para métricas
- **Dashboard administrativo:** KPIs en tiempo real (citas, ingresos, tasas de cancelación)
- **Rankings dinámicos:** Top doctores y servicios con límites configurables
- **Seguridad en analytics:** Solo administradores acceden a métricas del negocio
- **Bug crítico de service_id:** AppointmentRepository.Create() no guardaba service_id en BD - Corregido agregando el campo al INSERT
- **Queries de analytics optimizadas:** Uso de COUNT, SUM, GROUP BY y JOIN para agregaciones eficientes
- **Métricas en tiempo real:** Dashboard calcula KPIs desde BD sin caché
- **Migración SQLite a PostgreSQL:** Actualización de todos los placeholders (? → $1, $2), tipos de datos (BOOLEAN, TIMESTAMP) y sintaxis SQL
- **Tipos nativos PostgreSQL:** BOOLEAN y TIMESTAMP sin conversiones manuales
- **Connection pooling:** Configurado con Neon para múltiples conexiones concurrentes
- **Production-ready database:** Migración completa a PostgreSQL serverless en AWS vía Neon

---

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es propiedad de Zensoft.

## 👥 Autores

- **Equipo de Desarrollo** - Zensoft

## 📞 Contacto

Para preguntas o soporte, contactar al equipo de desarrollo de Zensoft.
