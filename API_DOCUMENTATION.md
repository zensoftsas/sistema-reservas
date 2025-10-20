# 📘 API Documentation - Sistema de Reservas Médicas

**Version:** 1.0.0
**Base URL:** `http://localhost:8080`
**Base URL Producción:** `https://tu-dominio.com`

Esta documentación está diseñada específicamente para el equipo de frontend. Cada endpoint incluye ejemplos exactos de request y response.

---

## 📋 Tabla de Contenidos

1. [Arquitectura de IDs](#arquitectura-de-ids)
2. [Autenticación](#autenticación)
3. [Usuarios](#usuarios)
4. [Citas Médicas](#citas-médicas)
5. [Servicios](#servicios)
6. [Horarios](#horarios)
7. [Analytics](#analytics)
8. [Códigos de Error](#códigos-de-error)

---

## 🔑 Arquitectura de IDs

**IMPORTANTE:** El sistema utiliza dos tipos de IDs para usuarios con roles:

### Tabla de IDs
```
┌─────────────┬──────────────┬─────────────────┬──────────────────┐
│ Tabla       │ Campo        │ Tipo            │ Uso              │
├─────────────┼──────────────┼─────────────────┼──────────────────┤
│ users       │ id           │ UUID (user.id)  │ Auth, perfiles   │
│ patients    │ id           │ UUID (patient.id│ Citas médicas    │
│ patients    │ user_id      │ UUID            │ FK → users.id    │
│ doctors     │ id           │ UUID (doctor.id)│ Citas, horarios  │
│ doctors     │ user_id      │ UUID            │ FK → users.id    │
│ appointments│ patient_id   │ UUID            │ FK → patients.id │
│ appointments│ doctor_id    │ UUID            │ FK → doctors.id  │
└─────────────┴──────────────┴─────────────────┴──────────────────┘
```

### Creación Automática de IDs

Cuando creas un usuario con rol `doctor` o `patient`:

```json
// REQUEST: POST /api/users
{
  "email": "doctor@example.com",
  "password": "password123",
  "first_name": "Carlos",
  "last_name": "Pérez",
  "role": "doctor"
}

// RESPONSE: 201 Created
{
  "id": "user-uuid-123",  // Este es el user.id
  "email": "doctor@example.com",
  "first_name": "Carlos",
  "last_name": "Pérez",
  "role": "doctor",
  "created_at": "2025-10-20T10:00:00Z"
}
```

**Lo que sucede automáticamente en el backend:**
1. Se crea el registro en `users` con `id = user-uuid-123`
2. Se crea automáticamente un registro en `doctors`:
   ```
   doctors.id = doctor-uuid-456
   doctors.user_id = user-uuid-123
   doctors.specialty = "Medicina General" (default)
   doctors.license_number = "LIC-uid123"
   ```

**Cuándo usar cada ID:**

| Operación | ID a usar | Ejemplo |
|-----------|-----------|---------|
| Login | `user.id` del response de login | `"user-uuid-123"` |
| Crear cita (como paciente) | Backend usa automáticamente tu `patient.id` | N/A (automático) |
| Seleccionar doctor | `doctor_id` del endpoint `/api/services/doctors?service_id=X` | `"user-uuid-789"` |
| Consultar horarios | `user.id` del doctor | `"user-uuid-789"` |
| Cancelar cita | `appointment.id` (del endpoint `/api/appointments/my`) | `"appointment-uuid-555"` |

---

## 🔐 Autenticación

### 1. Login

**Endpoint:** `POST /api/auth/login`
**Autenticación:** No requerida
**Descripción:** Obtiene un token JWT válido por 24 horas

#### Request
```json
{
  "email": "doctor@example.com",
  "password": "password123"
}
```

#### Response Success (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci11dWlkLTEyMyIsInJvbGUiOiJkb2N0b3IiLCJleHAiOjE3Mjk1MzI4MDB9.signature",
  "expires_at": "2025-10-21T10:00:00Z",
  "user": {
    "id": "user-uuid-123",
    "email": "doctor@example.com",
    "first_name": "Carlos",
    "last_name": "Pérez",
    "phone": "+51987654321",
    "role": "doctor",
    "is_active": true,
    "created_at": "2025-10-20T09:00:00Z"
  }
}
```

#### Errores Comunes
```json
// 400 Bad Request - Email o password faltantes
{
  "error": "Email and password are required"
}

// 401 Unauthorized - Credenciales incorrectas
{
  "error": "Invalid credentials"
}

// 401 Unauthorized - Usuario inactivo
{
  "error": "User is inactive"
}
```

#### Ejemplo cURL
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@example.com",
    "password": "password123"
  }'
```

#### Ejemplo JavaScript (fetch)
```javascript
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Login failed');
  }

  const data = await response.json();

  // Guardar token en localStorage
  localStorage.setItem('token', data.token);
  localStorage.setItem('user', JSON.stringify(data.user));

  return data;
};
```

---

## 👥 Usuarios

### 2. Crear Usuario (Registro)

**Endpoint:** `POST /api/users`
**Autenticación:** No requerida
**Descripción:** Crea un nuevo usuario. Si el rol es `doctor` o `patient`, automáticamente crea el registro en la tabla correspondiente.

#### Request
```json
{
  "email": "paciente@example.com",
  "password": "password123",
  "first_name": "María",
  "last_name": "García",
  "phone": "+51987654321",
  "role": "patient"
}
```

**Roles válidos:** `admin`, `doctor`, `patient`

**Validaciones:**
- `email`: requerido, único, formato email válido
- `password`: mínimo 8 caracteres
- `first_name`: requerido
- `last_name`: requerido
- `phone`: requerido
- `role`: debe ser `admin`, `doctor` o `patient`

#### Response Success (201 Created)
```json
{
  "id": "user-uuid-789",
  "email": "paciente@example.com",
  "first_name": "María",
  "last_name": "García",
  "role": "patient",
  "created_at": "2025-10-20T10:00:00Z"
}
```

**Nota:** Si `role = "patient"`, automáticamente se crea:
```
patients.id = "patient-uuid-abc"
patients.user_id = "user-uuid-789"
patients.document_number = "00000000" (placeholder, debe actualizarse)
patients.birthdate = hace 18 años (default)
```

Si `role = "doctor"`, automáticamente se crea:
```
doctors.id = "doctor-uuid-def"
doctors.user_id = "user-uuid-789"
doctors.specialty = "Medicina General"
doctors.license_number = "LIC-uid789"
```

#### Errores Comunes
```json
// 400 Bad Request - Email ya existe
{
  "error": "email already exists"
}

// 400 Bad Request - Password muy corto
{
  "error": "password must be at least 8 characters long"
}

// 400 Bad Request - Rol inválido
{
  "error": "invalid role: must be admin, doctor, or patient"
}

// 400 Bad Request - Campo requerido faltante
{
  "error": "first name is required"
}
```

#### Ejemplo JavaScript
```javascript
const register = async (userData) => {
  const response = await fetch('http://localhost:8080/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return await response.json();
};

// Uso
try {
  const newUser = await register({
    email: 'nuevo@example.com',
    password: 'password123',
    first_name: 'Juan',
    last_name: 'Pérez',
    phone: '+51999999999',
    role: 'patient'
  });
  console.log('Usuario creado:', newUser);
} catch (error) {
  console.error('Error:', error.message);
}
```

### 3. Obtener Perfil del Usuario Autenticado

**Endpoint:** `GET /api/users/me`
**Autenticación:** ✅ Requerida
**Descripción:** Obtiene el perfil completo del usuario actualmente autenticado

#### Request
```
GET /api/users/me
Authorization: Bearer {token}
```

#### Response Success (200 OK)
```json
{
  "id": "user-uuid-123",
  "email": "doctor@example.com",
  "first_name": "Carlos",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "doctor",
  "is_active": true,
  "created_at": "2025-10-20T09:00:00Z"
}
```

#### Errores
```json
// 401 Unauthorized - Token inválido o expirado
{
  "error": "Unauthorized"
}

// 404 Not Found - Usuario no existe
{
  "error": "User not found"
}
```

#### Ejemplo JavaScript
```javascript
const getMyProfile = async () => {
  const token = localStorage.getItem('token');

  const response = await fetch('http://localhost:8080/api/users/me', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch profile');
  }

  return await response.json();
};
```

### 4. Obtener Usuario por ID

**Endpoint:** `GET /api/users?id={uuid}`
**Autenticación:** No requerida
**Descripción:** Obtiene la información pública de un usuario específico

#### Request
```
GET /api/users?id=user-uuid-123
```

#### Response Success (200 OK)
```json
{
  "id": "user-uuid-123",
  "email": "doctor@example.com",
  "first_name": "Carlos",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "doctor",
  "is_active": true,
  "created_at": "2025-10-20T09:00:00Z"
}
```

#### Errores
```json
// 400 Bad Request - ID no proporcionado
{
  "error": "User ID is required"
}

// 404 Not Found
{
  "error": "User not found"
}
```

---

## 📅 Citas Médicas

### 5. Crear Cita

**Endpoint:** `POST /api/appointments`
**Autenticación:** ✅ Requerida (patient)
**Descripción:** Crea una nueva cita médica. El paciente se identifica automáticamente desde el token JWT.

#### Request
```json
{
  "doctor_id": "user-uuid-doctor-123",
  "service_id": "service-uuid-456",
  "appointment_date": "2025-10-25",
  "appointment_time": "14:30",
  "reason": "Consulta de control y revisión de exámenes"
}
```

**Headers:**
```
Authorization: Bearer {patient-token}
Content-Type: application/json
```

**Validaciones:**
- `doctor_id`: debe existir y ser un doctor activo
- `service_id`: debe existir y estar activo
- El doctor debe ofrecer ese servicio (relación `doctor_services`)
- El horario debe estar disponible (sin conflictos)
- `appointment_date`: formato YYYY-MM-DD
- `appointment_time`: formato HH:MM (24 horas)

#### Response Success (201 Created)
```json
{
  "id": "appointment-uuid-789",
  "patient_id": "patient-uuid-abc",
  "doctor_id": "doctor-uuid-def",
  "service_id": "service-uuid-456",
  "service_name": "Consulta Cardiológica",
  "scheduled_at": "2025-10-25T14:30:00Z",
  "duration": 45,
  "reason": "Consulta de control y revisión de exámenes",
  "status": "pending",
  "notes": "",
  "created_at": "2025-10-20T10:30:00Z",
  "updated_at": "2025-10-20T10:30:00Z"
}
```

**Estados de cita:**
- `pending`: Cita creada, esperando confirmación del doctor
- `confirmed`: Doctor confirmó la cita
- `completed`: Cita realizada
- `cancelled`: Cita cancelada

#### Errores Comunes
```json
// 400 Bad Request - service_id requerido
{
  "error": "service_id is required"
}

// 400 Bad Request - Formato de fecha inválido
{
  "error": "Invalid date or time format"
}

// 404 Not Found - Doctor no encontrado
{
  "error": "doctor not found"
}

// 404 Not Found - Servicio no encontrado
{
  "error": "service not found"
}

// 400 Bad Request - Doctor no ofrece el servicio
{
  "error": "doctor does not offer this service"
}

// 409 Conflict - Horario no disponible
{
  "error": "time slot is not available"
}
```

#### Ejemplo JavaScript
```javascript
const createAppointment = async (appointmentData) => {
  const token = localStorage.getItem('token');

  const response = await fetch('http://localhost:8080/api/appointments', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(appointmentData),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return await response.json();
};

// Uso
try {
  const appointment = await createAppointment({
    doctor_id: 'user-uuid-doctor-123',
    service_id: 'service-uuid-456',
    appointment_date: '2025-10-25',
    appointment_time: '14:30',
    reason: 'Consulta de control'
  });
  console.log('Cita creada:', appointment);
} catch (error) {
  console.error('Error:', error.message);
}
```

### 6. Obtener Mis Citas (Paciente)

**Endpoint:** `GET /api/appointments/my`
**Autenticación:** ✅ Requerida
**Descripción:** Obtiene todas las citas del paciente autenticado

#### Request
```
GET /api/appointments/my
Authorization: Bearer {patient-token}
```

#### Response Success (200 OK)
```json
[
  {
    "id": "appointment-uuid-789",
    "patient_id": "patient-uuid-abc",
    "doctor_id": "doctor-uuid-def",
    "service_id": "service-uuid-456",
    "service_name": "Consulta Cardiológica",
    "scheduled_at": "2025-10-25T14:30:00Z",
    "duration": 45,
    "reason": "Consulta de control",
    "status": "pending",
    "notes": "",
    "created_at": "2025-10-20T10:30:00Z",
    "updated_at": "2025-10-20T10:30:00Z"
  },
  {
    "id": "appointment-uuid-555",
    "patient_id": "patient-uuid-abc",
    "doctor_id": "doctor-uuid-ghi",
    "service_id": "service-uuid-111",
    "service_name": "Consulta General",
    "scheduled_at": "2025-10-22T10:00:00Z",
    "duration": 30,
    "reason": "Revisión anual",
    "status": "confirmed",
    "notes": "",
    "created_at": "2025-10-19T15:00:00Z",
    "updated_at": "2025-10-19T16:30:00Z"
  }
]
```

#### Errores
```json
// 401 Unauthorized
{
  "error": "Unauthorized"
}
```

#### Ejemplo JavaScript
```javascript
const getMyAppointments = async () => {
  const token = localStorage.getItem('token');

  const response = await fetch('http://localhost:8080/api/appointments/my', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch appointments');
  }

  return await response.json();
};
```

### 7. Obtener Citas del Doctor

**Endpoint:** `GET /api/appointments/doctor`
**Autenticación:** ✅ Requerida (doctor)
**Descripción:** Obtiene todas las citas del doctor autenticado

#### Request
```
GET /api/appointments/doctor
Authorization: Bearer {doctor-token}
```

#### Response Success (200 OK)
```json
[
  {
    "id": "appointment-uuid-789",
    "patient_id": "patient-uuid-abc",
    "doctor_id": "doctor-uuid-def",
    "service_id": "service-uuid-456",
    "service_name": "Consulta Cardiológica",
    "scheduled_at": "2025-10-25T14:30:00Z",
    "duration": 45,
    "reason": "Consulta de control",
    "status": "pending",
    "notes": "",
    "created_at": "2025-10-20T10:30:00Z",
    "updated_at": "2025-10-20T10:30:00Z"
  }
]
```

### 8. Cancelar Cita

**Endpoint:** `PUT /api/appointments/cancel?id={appointment-id}`
**Autenticación:** ✅ Requerida
**Descripción:** Cancela una cita. Solo el paciente, el doctor involucrado o un admin pueden cancelar.

**⚠️ CORRECCIÓN DE BUG:** Este endpoint ahora maneja correctamente la comparación de IDs entre `user.id` y `patient.id`/`doctor.id`.

#### Request
```
PUT /api/appointments/cancel?id=appointment-uuid-789
Authorization: Bearer {token}
Content-Type: application/json
```

```json
{
  "reason": "Tengo un compromiso urgente que no puedo posponer"
}
```

**Validaciones:**
- Solo el paciente de la cita puede cancelar
- Solo el doctor de la cita puede cancelar
- Los administradores pueden cancelar cualquier cita
- La cita no debe estar ya cancelada

#### Response Success (204 No Content)
```
(Sin body, solo código 204)
```

#### Errores Comunes
```json
// 400 Bad Request - ID no proporcionado
{
  "error": "Appointment ID is required"
}

// 404 Not Found
{
  "error": "appointment not found"
}

// 403 Forbidden - Sin permisos
{
  "error": "insufficient permissions to cancel this appointment"
}

// 400 Bad Request - Ya cancelada
{
  "error": "appointment is already cancelled"
}
```

#### Ejemplo JavaScript
```javascript
const cancelAppointment = async (appointmentId, reason) => {
  const token = localStorage.getItem('token');

  const response = await fetch(`http://localhost:8080/api/appointments/cancel?id=${appointmentId}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ reason }),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  // 204 No Content = éxito
  return true;
};

// Uso
try {
  await cancelAppointment('appointment-uuid-789', 'Tengo un compromiso urgente');
  console.log('Cita cancelada exitosamente');
} catch (error) {
  console.error('Error:', error.message);
}
```

### 9. Confirmar Cita (Doctor/Admin)

**Endpoint:** `PUT /api/appointments/confirm?id={appointment-id}`
**Autenticación:** ✅ Requerida (doctor o admin)
**Descripción:** Confirma una cita pendiente

#### Request
```
PUT /api/appointments/confirm?id=appointment-uuid-789
Authorization: Bearer {doctor-token}
```

#### Response Success (200 OK)
```json
{
  "id": "appointment-uuid-789",
  "patient_id": "patient-uuid-abc",
  "doctor_id": "doctor-uuid-def",
  "service_id": "service-uuid-456",
  "service_name": "Consulta Cardiológica",
  "scheduled_at": "2025-10-25T14:30:00Z",
  "duration": 45,
  "reason": "Consulta de control",
  "status": "confirmed",
  "notes": "",
  "created_at": "2025-10-20T10:30:00Z",
  "updated_at": "2025-10-20T11:00:00Z"
}
```

#### Errores
```json
// 400 Bad Request - No está en estado pending
{
  "error": "appointment is not in pending status"
}

// 403 Forbidden - Sin permisos
{
  "error": "insufficient permissions to confirm this appointment"
}
```

### 10. Completar Cita (Doctor/Admin)

**Endpoint:** `PUT /api/appointments/complete?id={appointment-id}`
**Autenticación:** ✅ Requerida (doctor o admin)
**Descripción:** Marca una cita confirmada como completada

#### Request
```
PUT /api/appointments/complete?id=appointment-uuid-789
Authorization: Bearer {doctor-token}
Content-Type: application/json
```

```json
{
  "notes": "Paciente en buen estado. Se recetó paracetamol. Control en 1 mes."
}
```

#### Response Success (200 OK)
```json
{
  "id": "appointment-uuid-789",
  "patient_id": "patient-uuid-abc",
  "doctor_id": "doctor-uuid-def",
  "service_id": "service-uuid-456",
  "service_name": "Consulta Cardiológica",
  "scheduled_at": "2025-10-25T14:30:00Z",
  "duration": 45,
  "reason": "Consulta de control",
  "status": "completed",
  "notes": "Paciente en buen estado. Se recetó paracetamol. Control en 1 mes.",
  "created_at": "2025-10-20T10:30:00Z",
  "updated_at": "2025-10-25T15:15:00Z"
}
```

---

## 🏥 Servicios

### 11. Listar Servicios Activos

**Endpoint:** `GET /api/services`
**Autenticación:** No requerida
**Descripción:** Obtiene todos los servicios médicos activos

#### Request
```
GET /api/services
```

#### Response Success (200 OK)
```json
[
  {
    "id": "service-uuid-456",
    "name": "Consulta Cardiológica",
    "description": "Evaluación cardiovascular completa con electrocardiograma",
    "duration_minutes": 45,
    "price": 150.0,
    "is_active": true,
    "created_at": "2025-10-15T09:00:00Z",
    "updated_at": "2025-10-15T09:00:00Z"
  },
  {
    "id": "service-uuid-111",
    "name": "Consulta General",
    "description": "Consulta médica general para diagnóstico y tratamiento",
    "duration_minutes": 30,
    "price": 80.0,
    "is_active": true,
    "created_at": "2025-10-15T09:00:00Z",
    "updated_at": "2025-10-15T09:00:00Z"
  }
]
```

#### Ejemplo JavaScript
```javascript
const getServices = async () => {
  const response = await fetch('http://localhost:8080/api/services');

  if (!response.ok) {
    throw new Error('Failed to fetch services');
  }

  return await response.json();
};

// Uso
const services = await getServices();
console.log('Servicios disponibles:', services);
```

### 12. Crear Servicio (Admin)

**Endpoint:** `POST /api/services/create`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Crea un nuevo servicio médico

#### Request
```json
{
  "name": "Consulta Dermatológica",
  "description": "Evaluación y tratamiento de problemas de la piel",
  "duration_minutes": 30,
  "price": 100.0
}
```

**Headers:**
```
Authorization: Bearer {admin-token}
Content-Type: application/json
```

#### Response Success (201 Created)
```json
{
  "id": "service-uuid-new",
  "name": "Consulta Dermatológica",
  "description": "Evaluación y tratamiento de problemas de la piel",
  "duration_minutes": 30,
  "price": 100.0,
  "is_active": true,
  "created_at": "2025-10-20T11:00:00Z",
  "updated_at": "2025-10-20T11:00:00Z"
}
```

### 13. Asignar Servicio a Doctor (Admin)

**Endpoint:** `POST /api/services/assign`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Asigna un servicio a un doctor (relación many-to-many)

#### Request
```json
{
  "doctor_id": "user-uuid-doctor-123",
  "service_id": "service-uuid-456"
}
```

**Headers:**
```
Authorization: Bearer {admin-token}
Content-Type: application/json
```

**Validaciones:**
- `doctor_id` debe ser un `user.id` de un usuario con rol `doctor`
- `service_id` debe existir y estar activo
- El backend automáticamente obtiene el `doctor.id` real desde `doctors.user_id`

#### Response Success (200 OK)
```json
{
  "message": "Service assigned to doctor successfully"
}
```

#### Errores
```json
// 400 Bad Request - Doctor no encontrado
{
  "error": "doctor not found"
}

// 400 Bad Request - Servicio no encontrado
{
  "error": "service not found"
}

// 400 Bad Request - Ya está asignado
{
  "error": "service already assigned to doctor"
}
```

### 14. Ver Doctores que Ofrecen un Servicio

**Endpoint:** `GET /api/services/doctors?service_id={uuid}`
**Autenticación:** No requerida
**Descripción:** Obtiene todos los doctores activos que ofrecen un servicio específico

#### Request
```
GET /api/services/doctors?service_id=service-uuid-456
```

#### Response Success (200 OK)
```json
[
  {
    "id": "user-uuid-doctor-123",
    "email": "dr.garcia@clinica.com",
    "first_name": "Dra. Ana",
    "last_name": "García",
    "phone": "+51987654321",
    "role": "doctor",
    "is_active": true,
    "created_at": "2025-10-10T08:00:00Z"
  },
  {
    "id": "user-uuid-doctor-456",
    "email": "dr.lopez@clinica.com",
    "first_name": "Dr. Carlos",
    "last_name": "López",
    "phone": "+51987654322",
    "role": "doctor",
    "is_active": true,
    "created_at": "2025-10-11T09:00:00Z"
  }
]
```

**Nota:** Los IDs retornados son `user.id`, que es lo que necesitas para consultar horarios disponibles.

#### Ejemplo JavaScript
```javascript
const getDoctorsByService = async (serviceId) => {
  const response = await fetch(`http://localhost:8080/api/services/doctors?service_id=${serviceId}`);

  if (!response.ok) {
    throw new Error('Failed to fetch doctors');
  }

  return await response.json();
};

// Uso
const doctors = await getDoctorsByService('service-uuid-456');
console.log('Doctores disponibles:', doctors);
```

### 15. Ver Horarios Disponibles ⭐

**Endpoint:** `GET /api/services/available-slots?doctor_id={user-uuid}&service_id={uuid}&date=YYYY-MM-DD`
**Autenticación:** No requerida
**Descripción:** Calcula los slots de tiempo disponibles para un doctor en una fecha específica

#### Request
```
GET /api/services/available-slots?doctor_id=user-uuid-doctor-123&service_id=service-uuid-456&date=2025-10-25
```

**Parámetros:**
- `doctor_id`: El `user.id` del doctor (obtenido del endpoint `/api/services/doctors`)
- `service_id`: El ID del servicio seleccionado
- `date`: Fecha en formato YYYY-MM-DD

#### Response Success (200 OK)
```json
[
  {
    "time": "09:00",
    "available": true
  },
  {
    "time": "09:45",
    "available": true
  },
  {
    "time": "10:30",
    "available": false
  },
  {
    "time": "11:15",
    "available": true
  },
  {
    "time": "12:00",
    "available": true
  },
  {
    "time": "15:00",
    "available": true
  },
  {
    "time": "15:45",
    "available": true
  },
  {
    "time": "16:30",
    "available": true
  },
  {
    "time": "17:15",
    "available": false
  }
]
```

**Cómo funciona:**
1. Obtiene el día de la semana de la fecha (ej: "monday")
2. Consulta los horarios configurados del doctor para ese día en la tabla `schedules`
3. Si no hay horarios configurados para ese día, retorna `[]` (el doctor no trabaja)
4. Si hay horarios, genera slots según la duración del servicio seleccionado
5. Marca como `available: false` los slots que tienen conflictos con citas existentes
6. Solo cuenta citas con estado `pending`, `confirmed` (no `cancelled`)

#### Ejemplo de Horario del Doctor
```
Lunes:
  - Bloque 1: 09:00-13:00 (slot_duration: 45 minutos)
  - Bloque 2: 15:00-18:00 (slot_duration: 45 minutos)
Martes:
  - Sin horarios configurados
```

**Request para Lunes:**
```
GET /api/services/available-slots?doctor_id=X&service_id=Y&date=2025-10-28
```

**Response:**
```json
[
  {"time": "09:00", "available": true},
  {"time": "09:45", "available": true},
  {"time": "10:30", "available": false},  // Ya tiene cita
  {"time": "11:15", "available": true},
  {"time": "12:00", "available": true},
  {"time": "12:45", "available": true},
  // GAP de 13:00 a 15:00
  {"time": "15:00", "available": true},
  {"time": "15:45", "available": true},
  {"time": "16:30", "available": true},
  {"time": "17:15", "available": true}
]
```

**Request para Martes:**
```
GET /api/services/available-slots?doctor_id=X&service_id=Y&date=2025-10-29
```

**Response:**
```json
[]
```

#### Ejemplo JavaScript
```javascript
const getAvailableSlots = async (doctorId, serviceId, date) => {
  const response = await fetch(
    `http://localhost:8080/api/services/available-slots?doctor_id=${doctorId}&service_id=${serviceId}&date=${date}`
  );

  if (!response.ok) {
    throw new Error('Failed to fetch available slots');
  }

  return await response.json();
};

// Uso
const slots = await getAvailableSlots('user-uuid-doctor-123', 'service-uuid-456', '2025-10-25');

// Filtrar solo disponibles
const availableSlots = slots.filter(slot => slot.available);
console.log('Horarios disponibles:', availableSlots);
```

---

## ⏰ Horarios

### 16. Crear Horario (Admin)

**Endpoint:** `POST /api/schedules`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Crea un bloque horario para un doctor en un día específico

#### Request
```json
{
  "doctor_id": "user-uuid-doctor-123",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "13:00",
  "slot_duration": 45
}
```

**Headers:**
```
Authorization: Bearer {admin-token}
Content-Type: application/json
```

**Validaciones:**
- `doctor_id`: Debe ser un `user.id` de un doctor activo
- `day_of_week`: Debe ser uno de: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`
- `start_time` y `end_time`: Formato HH:MM (00:00 - 23:59)
- `start_time` debe ser antes de `end_time`
- `slot_duration`: Entero mayor que 0 (en minutos)
- No puede solaparse con otro horario del mismo doctor en el mismo día

#### Response Success (201 Created)
```json
{
  "id": "schedule-uuid-789",
  "doctor_id": "doctor-uuid-def",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "13:00",
  "slot_duration": 45,
  "is_active": true,
  "created_at": "2025-10-20T11:30:00Z",
  "updated_at": "2025-10-20T11:30:00Z"
}
```

**Nota:** El `doctor_id` en el response es el `doctor.id` real (no el `user.id`).

#### Errores
```json
// 400 Bad Request - day_of_week inválido
{
  "error": "invalid day_of_week, must be monday-sunday"
}

// 400 Bad Request - Formato de hora inválido
{
  "error": "start_time must be in HH:MM format"
}

// 400 Bad Request - start_time >= end_time
{
  "error": "start_time must be before end_time"
}

// 400 Bad Request - Overlap con horario existente
{
  "error": "schedule overlaps with existing schedule"
}
```

#### Ejemplo: Configurar múltiples bloques
```javascript
// Doctor trabaja Lunes con horario partido (mañana + tarde)
await createSchedule({
  doctor_id: 'user-uuid-doctor-123',
  day_of_week: 'monday',
  start_time: '09:00',
  end_time: '13:00',
  slot_duration: 45
});

await createSchedule({
  doctor_id: 'user-uuid-doctor-123',
  day_of_week: 'monday',
  start_time: '15:00',
  end_time: '18:00',
  slot_duration: 45
});
```

### 17. Ver Horarios de un Doctor

**Endpoint:** `GET /api/schedules/doctor/{user-id}`
**Autenticación:** No requerida
**Descripción:** Obtiene todos los horarios configurados para un doctor

#### Request
```
GET /api/schedules/doctor/user-uuid-doctor-123
```

#### Response Success (200 OK)
```json
[
  {
    "id": "schedule-uuid-789",
    "doctor_id": "doctor-uuid-def",
    "day_of_week": "monday",
    "start_time": "09:00",
    "end_time": "13:00",
    "slot_duration": 45,
    "is_active": true,
    "created_at": "2025-10-20T11:30:00Z",
    "updated_at": "2025-10-20T11:30:00Z"
  },
  {
    "id": "schedule-uuid-790",
    "doctor_id": "doctor-uuid-def",
    "day_of_week": "monday",
    "start_time": "15:00",
    "end_time": "18:00",
    "slot_duration": 45,
    "is_active": true,
    "created_at": "2025-10-20T11:31:00Z",
    "updated_at": "2025-10-20T11:31:00Z"
  },
  {
    "id": "schedule-uuid-791",
    "doctor_id": "doctor-uuid-def",
    "day_of_week": "friday",
    "start_time": "08:00",
    "end_time": "12:00",
    "slot_duration": 30,
    "is_active": true,
    "created_at": "2025-10-20T11:32:00Z",
    "updated_at": "2025-10-20T11:32:00Z"
  }
]
```

#### Ejemplo JavaScript
```javascript
const getDoctorSchedules = async (doctorUserId) => {
  const response = await fetch(`http://localhost:8080/api/schedules/doctor/${doctorUserId}`);

  if (!response.ok) {
    throw new Error('Failed to fetch schedules');
  }

  return await response.json();
};

// Uso
const schedules = await getDoctorSchedules('user-uuid-doctor-123');

// Agrupar por día
const schedulesByDay = schedules.reduce((acc, schedule) => {
  if (!acc[schedule.day_of_week]) {
    acc[schedule.day_of_week] = [];
  }
  acc[schedule.day_of_week].push(schedule);
  return acc;
}, {});

console.log('Horarios por día:', schedulesByDay);
```

### 18. Eliminar Horario (Admin)

**Endpoint:** `DELETE /api/schedules/{schedule-id}`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Elimina un horario configurado

#### Request
```
DELETE /api/schedules/schedule-uuid-789
Authorization: Bearer {admin-token}
```

#### Response Success (200 OK)
```json
{
  "message": "Schedule deleted successfully"
}
```

#### Errores
```json
// 404 Not Found
{
  "error": "schedule not found"
}
```

---

## 📊 Analytics (Solo Admin)

### 19. Dashboard Summary

**Endpoint:** `GET /api/analytics/dashboard`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Obtiene resumen con KPIs principales del sistema

#### Request
```
GET /api/analytics/dashboard
Authorization: Bearer {admin-token}
```

#### Response Success (200 OK)
```json
{
  "total_appointments": 150,
  "pending_appointments": 20,
  "confirmed_appointments": 45,
  "completed_appointments": 75,
  "cancelled_appointments": 10,
  "total_patients": 80,
  "total_doctors": 12,
  "total_revenue": 12500.50,
  "cancellation_rate": 6.67
}
```

**Métricas:**
- `total_appointments`: Total de citas en el sistema
- `pending_appointments`: Citas pendientes de confirmación
- `confirmed_appointments`: Citas confirmadas
- `completed_appointments`: Citas realizadas
- `cancelled_appointments`: Citas canceladas
- `total_patients`: Número de pacientes activos
- `total_doctors`: Número de doctores activos
- `total_revenue`: Ingresos totales (suma de citas completadas)
- `cancellation_rate`: Porcentaje de cancelación

### 20. Revenue Stats

**Endpoint:** `GET /api/analytics/revenue`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Análisis de ingresos por servicio

#### Request
```
GET /api/analytics/revenue
Authorization: Bearer {admin-token}
```

#### Response Success (200 OK)
```json
[
  {
    "service_id": "service-uuid-456",
    "service_name": "Consulta Cardiológica",
    "total_citas": 35,
    "revenue": 5250.00
  },
  {
    "service_id": "service-uuid-111",
    "service_name": "Consulta General",
    "total_citas": 60,
    "revenue": 4800.00
  }
]
```

**Nota:** Solo incluye citas completadas, ordenadas por ingresos (mayor a menor).

### 21. Top Doctors

**Endpoint:** `GET /api/analytics/top-doctors?limit={n}`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Ranking de doctores por número de citas

#### Request
```
GET /api/analytics/top-doctors?limit=10
Authorization: Bearer {admin-token}
```

**Query params:**
- `limit` (opcional): Número de doctores a retornar (default: 10)

#### Response Success (200 OK)
```json
[
  {
    "doctor_id": "doctor-uuid-def",
    "doctor_name": "Doctor abc123...",
    "total_appointments": 85,
    "completed_appointments": 78
  },
  {
    "doctor_id": "doctor-uuid-ghi",
    "doctor_name": "Doctor def456...",
    "total_appointments": 67,
    "completed_appointments": 62
  }
]
```

### 22. Top Services

**Endpoint:** `GET /api/analytics/top-services?limit={n}`
**Autenticación:** ✅ Requerida (admin)
**Descripción:** Ranking de servicios más solicitados

#### Request
```
GET /api/analytics/top-services?limit=10
Authorization: Bearer {admin-token}
```

**Query params:**
- `limit` (opcional): Número de servicios a retornar (default: 10)

#### Response Success (200 OK)
```json
[
  {
    "service_id": "service-uuid-111",
    "service_name": "Consulta General",
    "total_citas": 95
  },
  {
    "service_id": "service-uuid-456",
    "service_name": "Consulta Cardiológica",
    "total_citas": 78
  }
]
```

---

## ⚠️ Códigos de Error

### Códigos HTTP
- `200 OK`: Solicitud exitosa
- `201 Created`: Recurso creado exitosamente
- `204 No Content`: Operación exitosa sin contenido
- `400 Bad Request`: Datos inválidos o faltantes
- `401 Unauthorized`: Token inválido, expirado o no proporcionado
- `403 Forbidden`: Sin permisos para realizar la operación
- `404 Not Found`: Recurso no encontrado
- `405 Method Not Allowed`: Método HTTP incorrecto
- `409 Conflict`: Conflicto (ej: horario ocupado)
- `500 Internal Server Error`: Error del servidor

### Formato de Errores
```json
{
  "error": "Descripción del error"
}
```

O simplemente texto plano:
```
Error message here
```

---

## 🔄 Flujo Completo de Reserva de Citas

### Paso 1: Login del Paciente
```javascript
const loginData = await login('paciente@example.com', 'password123');
// Guardar: loginData.token, loginData.user
```

### Paso 2: Ver Servicios Disponibles
```javascript
const services = await getServices();
// Mostrar lista de servicios al usuario
// Usuario selecciona serviceId
```

### Paso 3: Ver Doctores que Ofrecen el Servicio
```javascript
const doctors = await getDoctorsByService(selectedServiceId);
// Mostrar lista de doctores al usuario
// Usuario selecciona doctorId (user.id del doctor)
```

### Paso 4: Seleccionar Fecha
```javascript
// Usuario selecciona fecha desde un calendario
const selectedDate = '2025-10-25';
```

### Paso 5: Ver Horarios Disponibles
```javascript
const slots = await getAvailableSlots(selectedDoctorId, selectedServiceId, selectedDate);
const availableSlots = slots.filter(s => s.available);
// Mostrar horarios disponibles
// Usuario selecciona un horario
```

### Paso 6: Crear Cita
```javascript
const appointment = await createAppointment({
  doctor_id: selectedDoctorId,
  service_id: selectedServiceId,
  appointment_date: selectedDate,
  appointment_time: selectedTime, // ej: "14:30"
  reason: 'Consulta de control'
});
// Mostrar confirmación
```

### Paso 7: Ver Mis Citas
```javascript
const myAppointments = await getMyAppointments();
// Mostrar lista de citas del paciente
```

### Paso 8 (Opcional): Cancelar Cita
```javascript
await cancelAppointment(appointmentId, 'Tengo un compromiso urgente');
// Mostrar confirmación de cancelación
```

---

## 🛡️ Seguridad

### Headers de Autenticación
Todos los endpoints protegidos requieren:
```
Authorization: Bearer {token}
```

### Expiración de Token
- Los tokens JWT expiran en **24 horas**
- Después de expirar, el usuario debe hacer login nuevamente
- No hay refresh tokens implementados (por ahora)

### Validación de Tokens
```javascript
const isTokenValid = () => {
  const token = localStorage.getItem('token');
  if (!token) return false;

  try {
    // Decodificar el payload del JWT (sin verificar firma)
    const payload = JSON.parse(atob(token.split('.')[1]));
    const expiration = payload.exp * 1000; // Convertir a milisegundos
    return Date.now() < expiration;
  } catch {
    return false;
  }
};

// Uso
if (!isTokenValid()) {
  // Redirigir a login
  window.location.href = '/login';
}
```

### Manejo de Errores de Autenticación
```javascript
const apiCall = async (url, options) => {
  const response = await fetch(url, options);

  if (response.status === 401) {
    // Token inválido o expirado
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
    throw new Error('Session expired');
  }

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return await response.json();
};
```

---

## 📝 Notas Importantes

### 1. IDs vs User IDs
- **Login y perfiles:** Usa `user.id`
- **Seleccionar doctor:** Usa `user.id` (del endpoint `/api/services/doctors`)
- **Consultar horarios:** Usa `user.id` del doctor
- **Crear cita:** El backend automáticamente convierte `user.id` a `patient.id` y `doctor.id`
- **Cancelar cita:** Usa `appointment.id` (no user.id)

### 2. Creación Automática de Roles
- Al crear usuario con `role = "patient"`, automáticamente se crea en tabla `patients`
- Al crear usuario con `role = "doctor"`, automáticamente se crea en tabla `doctors`
- Al crear usuario con `role = "admin"`, solo se crea en tabla `users`

### 3. Horarios y Slots
- Los slots se calculan dinámicamente basados en:
  - Horarios configurados del doctor (`schedules` table)
  - Duración del servicio seleccionado
  - Citas existentes del doctor en esa fecha
- Si un doctor no tiene horarios configurados para un día, retorna `[]` (no trabaja)

### 4. Estados de Cita
```
pending → confirmed → completed
        ↘ cancelled
```

---

## 🚀 Checklist de Integración Frontend

- [ ] Implementar login y guardar token en localStorage
- [ ] Agregar header `Authorization: Bearer {token}` a requests protegidos
- [ ] Manejar errores 401 (token expirado)
- [ ] Implementar flujo de registro de usuarios
- [ ] Listar servicios disponibles
- [ ] Mostrar doctores por servicio
- [ ] Calendario para seleccionar fecha
- [ ] Mostrar slots disponibles con estados (disponible/ocupado)
- [ ] Crear cita con validaciones
- [ ] Listar citas del paciente
- [ ] Cancelar citas con confirmación
- [ ] Confirmar citas (vista doctor)
- [ ] Completar citas (vista doctor)
- [ ] Dashboard de analytics (vista admin)

---

**Última actualización:** 2025-10-20
**Autor:** Backend Team - Zensoft
**Contacto:** Para dudas sobre la API, consultar con el equipo de backend.
