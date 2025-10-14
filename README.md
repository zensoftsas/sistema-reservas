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
- **Sesión 8:** Horarios Personalizados + Múltiples Bloques + Días No Laborables

### Tecnologías Elegidas:
- **Go 1.21+** - Performance y concurrencia
- **SQLite** - Desarrollo rápido (migrar a PostgreSQL en producción)
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
