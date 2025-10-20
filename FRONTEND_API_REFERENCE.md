# 📘 API Reference - Sistema de Reservas Médicas

**Versión:** 1.0.0
**Base URL Desarrollo:** `http://localhost:8080`
**Base URL Producción:** `https://tu-dominio.com`

---

## 🚨 IMPORTANTE: Leer Esto Primero

### Arquitectura de IDs

El sistema usa **diferentes IDs para diferentes tablas**. Esto es crítico para evitar errores:

| Tabla | ID | Uso |
|-------|----|----|
| `users` | `user.id` | Login, autenticación, seleccionar doctores |
| `patients` | `patient.id` | Citas médicas (automático) |
| `doctors` | `doctor.id` | Citas médicas (automático) |
| `appointments` | `appointment.id` | Cancelar, confirmar, completar citas |

**Regla simple:**
- Para **seleccionar doctores y ver horarios**: usa `user.id`
- Para **cancelar/modificar citas**: usa `appointment.id`
- El backend hace todas las conversiones automáticamente

---

## 🔐 Autenticación

Todos los endpoints protegidos requieren:
```
Authorization: Bearer {token}
```

El token expira en **24 horas**.

---

## 📋 Endpoints

### 1. Login

```
POST /api/auth/login
```

**Request:**
```json
{
  "email": "usuario@example.com",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-10-21T10:00:00Z",
  "user": {
    "id": "user-uuid-123",
    "email": "usuario@example.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "patient",
    "is_active": true
  }
}
```

**Errores:**
- `400`: Email o password faltantes
- `401`: Credenciales incorrectas

---

### 2. Registrar Usuario

```
POST /api/users
```

**Request:**
```json
{
  "email": "nuevo@example.com",
  "password": "password123",
  "first_name": "María",
  "last_name": "García",
  "phone": "+51987654321",
  "role": "patient"
}
```

**Roles válidos:** `patient`, `doctor`, `admin`

**Response 201:**
```json
{
  "id": "user-uuid-789",
  "email": "nuevo@example.com",
  "first_name": "María",
  "last_name": "García",
  "role": "patient",
  "created_at": "2025-10-20T10:00:00Z"
}
```

⚠️ **Si role es `patient` o `doctor`, automáticamente se crea en su tabla correspondiente.**

**Errores:**
- `400`: Email ya existe, password muy corto, rol inválido

---

### 3. Mi Perfil

```
GET /api/users/me
Authorization: Bearer {token}
```

**Response 200:**
```json
{
  "id": "user-uuid-123",
  "email": "usuario@example.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "phone": "+51987654321",
  "role": "patient",
  "is_active": true
}
```

---

### 4. Listar Servicios

```
GET /api/services
```

**Response 200:**
```json
[
  {
    "id": "service-uuid-456",
    "name": "Consulta Cardiológica",
    "description": "Evaluación cardiovascular completa",
    "duration_minutes": 45,
    "price": 150.0,
    "is_active": true
  },
  {
    "id": "service-uuid-111",
    "name": "Consulta General",
    "description": "Consulta médica general",
    "duration_minutes": 30,
    "price": 80.0,
    "is_active": true
  }
]
```

---

### 5. Doctores que Ofrecen un Servicio

```
GET /api/services/doctors?service_id={uuid}
```

**Response 200:**
```json
[
  {
    "id": "user-uuid-doctor-123",
    "email": "dr.garcia@clinica.com",
    "first_name": "Dra. Ana",
    "last_name": "García",
    "phone": "+51987654321",
    "role": "doctor",
    "is_active": true
  }
]
```

⚠️ **El `id` retornado es el `user.id` del doctor (necesario para siguiente endpoint).**

---

### 6. Horarios Disponibles ⭐

```
GET /api/services/available-slots?doctor_id={user-uuid}&service_id={uuid}&date=YYYY-MM-DD
```

**Parámetros:**
- `doctor_id`: user.id del doctor (del endpoint anterior)
- `service_id`: ID del servicio
- `date`: Fecha en formato YYYY-MM-DD

**Response 200:**
```json
[
  {"time": "09:00", "available": true},
  {"time": "09:45", "available": true},
  {"time": "10:30", "available": false},
  {"time": "11:15", "available": true},
  {"time": "12:00", "available": true},
  {"time": "15:00", "available": true},
  {"time": "15:45", "available": true}
]
```

**Si el doctor no trabaja ese día, retorna:** `[]`

**Filtrar solo disponibles:**
```javascript
const available = slots.filter(s => s.available);
```

---

### 7. Crear Cita

```
POST /api/appointments
Authorization: Bearer {token}
```

**Request:**
```json
{
  "doctor_id": "user-uuid-doctor-123",
  "service_id": "service-uuid-456",
  "appointment_date": "2025-10-25",
  "appointment_time": "14:30",
  "reason": "Consulta de control"
}
```

⚠️ **Usa `user.id` del doctor, NO `doctor.id`. El backend hace la conversión.**

**Response 201:**
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
  "status": "pending",
  "created_at": "2025-10-20T10:30:00Z"
}
```

**Estados de cita:**
- `pending`: Creada, esperando confirmación
- `confirmed`: Confirmada por doctor
- `completed`: Realizada
- `cancelled`: Cancelada

**Errores:**
- `400`: service_id requerido, fecha/hora inválida
- `404`: Doctor o servicio no encontrado
- `400`: Doctor no ofrece ese servicio
- `409`: Horario no disponible

---

### 8. Mis Citas

```
GET /api/appointments/my
Authorization: Bearer {token}
```

**Response 200:**
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
    "created_at": "2025-10-20T10:30:00Z"
  }
]
```

---

### 9. Cancelar Cita

```
PUT /api/appointments/cancel?id={appointment-id}
Authorization: Bearer {token}
```

**Request:**
```json
{
  "reason": "Tengo un compromiso urgente"
}
```

⚠️ **Usa `appointment.id`, NO `user.id`.**

**Response:** `204 No Content`

**Errores:**
- `404`: Cita no encontrada
- `403`: Sin permisos para cancelar
- `400`: Cita ya cancelada

---

### 10. Citas del Doctor (Solo Doctor)

```
GET /api/appointments/doctor
Authorization: Bearer {doctor-token}
```

**Response 200:** Igual formato que "Mis Citas"

---

### 11. Confirmar Cita (Solo Doctor/Admin)

```
PUT /api/appointments/confirm?id={appointment-id}
Authorization: Bearer {doctor-token}
```

**Response 200:** Objeto de cita con `status: "confirmed"`

**Errores:**
- `400`: No está en estado pending
- `403`: Sin permisos

---

### 12. Completar Cita (Solo Doctor/Admin)

```
PUT /api/appointments/complete?id={appointment-id}
Authorization: Bearer {doctor-token}
```

**Request:**
```json
{
  "notes": "Paciente en buen estado. Control en 1 mes."
}
```

**Response 200:** Objeto de cita con `status: "completed"`

---

### 13. Ver Horarios de Doctor

```
GET /api/schedules/doctor/{user-id}
```

**Response 200:**
```json
[
  {
    "id": "schedule-uuid-789",
    "doctor_id": "doctor-uuid-def",
    "day_of_week": "monday",
    "start_time": "09:00",
    "end_time": "13:00",
    "slot_duration": 45,
    "is_active": true
  },
  {
    "id": "schedule-uuid-790",
    "doctor_id": "doctor-uuid-def",
    "day_of_week": "monday",
    "start_time": "15:00",
    "end_time": "18:00",
    "slot_duration": 45,
    "is_active": true
  }
]
```

---

### 14. Crear Horario (Solo Admin)

```
POST /api/schedules
Authorization: Bearer {admin-token}
```

**Request:**
```json
{
  "doctor_id": "user-uuid-doctor-123",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "13:00",
  "slot_duration": 45
}
```

**Días válidos:** `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`

**Response 201:** Objeto de horario creado

---

### 15. Dashboard (Solo Admin)

```
GET /api/analytics/dashboard
Authorization: Bearer {admin-token}
```

**Response 200:**
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

---

### 16. Ingresos por Servicio (Solo Admin)

```
GET /api/analytics/revenue
Authorization: Bearer {admin-token}
```

**Response 200:**
```json
[
  {
    "service_id": "service-uuid-456",
    "service_name": "Consulta Cardiológica",
    "total_citas": 35,
    "revenue": 5250.00
  }
]
```

---

### 17. Top Doctores (Solo Admin)

```
GET /api/analytics/top-doctors?limit=10
Authorization: Bearer {admin-token}
```

**Response 200:**
```json
[
  {
    "doctor_id": "doctor-uuid-def",
    "doctor_name": "Doctor abc123...",
    "total_appointments": 85,
    "completed_appointments": 78
  }
]
```

---

### 18. Top Servicios (Solo Admin)

```
GET /api/analytics/top-services?limit=10
Authorization: Bearer {admin-token}
```

**Response 200:**
```json
[
  {
    "service_id": "service-uuid-111",
    "service_name": "Consulta General",
    "total_citas": 95
  }
]
```

---

## 🔄 Flujo Completo: Reservar Cita

```javascript
// 1. Login
const { token, user } = await login(email, password);
localStorage.setItem('token', token);

// 2. Ver servicios
const services = await fetch('/api/services').then(r => r.json());

// 3. Usuario selecciona servicio
const selectedServiceId = services[0].id;

// 4. Ver doctores que ofrecen ese servicio
const doctors = await fetch(
  `/api/services/doctors?service_id=${selectedServiceId}`
).then(r => r.json());

// 5. Usuario selecciona doctor
const selectedDoctorId = doctors[0].id; // user.id del doctor

// 6. Usuario selecciona fecha
const selectedDate = '2025-10-25';

// 7. Ver horarios disponibles
const slots = await fetch(
  `/api/services/available-slots?doctor_id=${selectedDoctorId}&service_id=${selectedServiceId}&date=${selectedDate}`
).then(r => r.json());

// 8. Mostrar solo horarios disponibles
const availableSlots = slots.filter(s => s.available);

// 9. Usuario selecciona horario
const selectedTime = '14:30';

// 10. Crear cita
const appointment = await fetch('/api/appointments', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    doctor_id: selectedDoctorId,
    service_id: selectedServiceId,
    appointment_date: selectedDate,
    appointment_time: selectedTime,
    reason: 'Consulta de control'
  })
}).then(r => r.json());

console.log('Cita creada:', appointment);
```

---

## 🔄 Flujo: Cancelar Cita

```javascript
// 1. Obtener mis citas
const myAppointments = await fetch('/api/appointments/my', {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());

// 2. Seleccionar cita a cancelar
const appointmentToCancel = myAppointments[0];

// 3. Cancelar
await fetch(`/api/appointments/cancel?id=${appointmentToCancel.id}`, {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    reason: 'Tengo un compromiso urgente'
  })
});

console.log('Cita cancelada exitosamente');
```

---

## ⚠️ Códigos de Error

- `200 OK`: Éxito
- `201 Created`: Recurso creado
- `204 No Content`: Operación exitosa sin contenido
- `400 Bad Request`: Datos inválidos
- `401 Unauthorized`: Token inválido/expirado
- `403 Forbidden`: Sin permisos
- `404 Not Found`: Recurso no encontrado
- `409 Conflict`: Conflicto (ej: horario ocupado)
- `500 Internal Server Error`: Error del servidor

---

## 🛠️ Utilidades JavaScript

### API Client Helper

```javascript
class ApiClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem('token');

    const config = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      }
    };

    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, config);

    if (response.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
      throw new Error('Session expired');
    }

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    if (response.status === 204) {
      return null;
    }

    return await response.json();
  }

  get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  }

  post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  put(endpoint, data) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }

  delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
}

// Uso
const api = new ApiClient('http://localhost:8080');

// Login
const loginData = await api.post('/api/auth/login', {
  email: 'usuario@example.com',
  password: 'password123'
});

// Crear cita
const appointment = await api.post('/api/appointments', {
  doctor_id: 'user-uuid',
  service_id: 'service-uuid',
  appointment_date: '2025-10-25',
  appointment_time: '14:30',
  reason: 'Consulta'
});

// Cancelar cita
await api.put(`/api/appointments/cancel?id=${appointmentId}`, {
  reason: 'Motivo'
});
```

### Verificar Token Expirado

```javascript
const isTokenExpired = () => {
  const token = localStorage.getItem('token');
  if (!token) return true;

  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return Date.now() >= payload.exp * 1000;
  } catch {
    return true;
  }
};

// Uso
if (isTokenExpired()) {
  window.location.href = '/login';
}
```

---

## 📝 Notas Importantes

1. **Siempre usa `user.id` para seleccionar doctores y consultar horarios**
2. **Usa `appointment.id` para cancelar/modificar citas**
3. **El backend convierte automáticamente los IDs internamente**
4. **Los slots se filtran por `available: true`**
5. **El token expira en 24 horas**
6. **Solo se pueden cancelar citas `pending` o `confirmed`**
7. **Si un doctor no tiene horarios un día, retorna `[]`**

---

## 🆘 Errores Comunes

### "doctor does not offer this service"
**Solución:** Asegúrate de usar doctores del endpoint `/api/services/doctors?service_id=X`

### "time slot is not available"
**Solución:** Verifica que el slot tenga `available: true`

### "insufficient permissions to cancel this appointment"
**Solución:** Solo el paciente, doctor involucrado o admin puede cancelar

### 401 Unauthorized
**Solución:** Token expirado, haz login nuevamente

---

## 📞 Soporte

**Dudas:** Contactar al equipo de backend
**Swagger UI:** http://localhost:8080/swagger/index.html

---

**Última actualización:** 2025-10-20
**Versión:** 1.0.0
**Autor:** Backend Team - Zensoft
