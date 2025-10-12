# 🔧 Backend API - Referencia Rápida

## 🚀 Iniciar Servidor
```bash
go run cmd/api/main.go
```

## 📝 Configuración (.env)
```env
SERVER_PORT=8080
DATABASE_PATH=clinica.db
JWT_SECRET=tu-secret-super-seguro
JWT_EXPIRATION_HOURS=48

# SendGrid
SENDGRID_API_KEY=SG.tu-api-key
SENDGRID_FROM_EMAIL=tu-email-verificado@gmail.com
SENDGRID_FROM_NAME=Clinica Internacional
```

## 📡 Endpoints Principales

### 🔐 Autenticación
**POST /api/auth/login**
```json
Request:
{
  "email": "zensoftsas@gmail.com",
  "password": "patient1234"
}

Response:
{
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "...",
    "role": "patient|doctor|admin",
    "first_name": "...",
    "last_name": "..."
  }
}
```

### 📅 Citas
**POST /api/appointments** (requiere auth)
```json
Request:
{
  "doctor_id": "uuid",
  "appointment_date": "2025-10-15",
  "appointment_time": "14:30",
  "reason": "Consulta general"
}
```

**GET /api/appointments/my** (requiere auth)
- Retorna citas del usuario autenticado

**PUT /api/appointments/confirm?id=uuid** (doctor/admin)
- Confirma una cita

**PUT /api/appointments/complete?id=uuid** (doctor/admin)
- Marca cita como completada

**PUT /api/appointments/cancel** (requiere auth)
```json
Request:
{
  "appointment_id": "uuid",
  "reason": "Motivo de cancelación"
}
```

### 👨‍⚕️ Doctores
**GET /api/doctors/search?specialty=cardio**
- Búsqueda por especialidad (primeras 6 letras)
- Público (no requiere auth)

### 👤 Usuarios
**GET /api/users/me** (requiere auth)
- Retorna perfil del usuario

**POST /api/users**
- Crear nuevo usuario

## 🧪 Tests Rápidos (cURL)

**Login como paciente:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"zensoftsas@gmail.com","password":"patient1234"}'
```

**Login como doctor:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dr.garcia@clinica.com","password":"doctor1234"}'
```

**Buscar doctores:**
```bash
curl http://localhost:8080/api/doctors/search?specialty=cardio
```

**Crear cita (con token):**
```bash
curl -X POST http://localhost:8080/api/appointments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU-TOKEN-AQUI" \
  -d '{"doctor_id":"470e338a-dcc8-4a80-8736-011ed38c74b7","appointment_date":"2025-10-15","appointment_time":"10:00","reason":"Consulta"}'
```

## 🏗️ Arquitectura

**Clean Architecture (4 capas)**
- **Domain:** Entidades del negocio
- **Use Cases:** Lógica de negocio
- **Repository:** Acceso a datos
- **Delivery:** HTTP handlers

### 📂 Estructura Principal
```
internal/
├── domain/          # Entidades
├── usecase/         # Casos de uso
├── delivery/http/   # Handlers y rutas
└── repository/      # Repositorios (SQLite)

pkg/
├── config/          # Configuración
├── email/           # SendGrid
└── reminder/        # Sistema de recordatorios
```

## 🔑 Roles y Permisos

- **patient:** Crear citas, ver sus citas, cancelar
- **doctor:** Confirmar/completar citas, ver agenda
- **admin:** Todos los permisos

## 📧 Notificaciones (SendGrid)

- Email al crear cita
- Email al confirmar cita
- Email al cancelar cita
- Email al completar cita
- Recordatorio 24h antes (automático)

## 🔄 Sistema de Recordatorios

- Cron job cada 10 minutos
- Envía recordatorio 24h antes
- Automático para citas confirmadas

## 📚 GitHub
https://github.com/zensoftsas/sistema-reservas

## 🎨 Frontend
[Sistema Reservas Frontend](https://github.com/zensoftsas/sistema-reservas-frontend) (próximamente)
