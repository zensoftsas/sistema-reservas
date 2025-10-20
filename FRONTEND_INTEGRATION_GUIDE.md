# 🎯 Guía de Integración Frontend - Sistema de Reservas

Esta guía te ayudará a integrar el frontend con el backend sin errores. **Lee esto primero antes de comenzar.**

---

## 🔴 CORRECCIONES CRÍTICAS REALIZADAS

### 1. Bug de Cancelación de Citas - CORREGIDO ✅

**Problema anterior:**
```javascript
// ❌ ESTO FALLABA ANTES
fetch('/api/appointments/cancel?id=appointment-123', {
  method: 'PUT',
  headers: { 'Authorization': 'Bearer token' },
  body: JSON.stringify({ reason: 'Motivo' })
});
// Error: "insufficient permissions to cancel this appointment"
```

**Causa del bug:**
El backend comparaba `user.id` con `patient.id` y `doctor.id`, que son IDs diferentes.

**Solución implementada:**
Ahora el backend convierte correctamente los IDs antes de comparar permisos.

**Cómo usarlo ahora:**
```javascript
// ✅ AHORA FUNCIONA CORRECTAMENTE
const cancelAppointment = async (appointmentId, reason) => {
  const token = localStorage.getItem('token');

  const response = await fetch(
    `http://localhost:8080/api/appointments/cancel?id=${appointmentId}`,
    {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ reason })
    }
  );

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return true; // 204 No Content = éxito
};
```

---

## 🗂️ Arquitectura de IDs - ENTENDER ESTO ES CRÍTICO

El sistema usa **diferentes tablas para diferentes propósitos**:

### Estructura de Tablas

```
users                    patients                 doctors
├── id (user.id)        ├── id (patient.id)     ├── id (doctor.id)
├── email               ├── user_id → users.id  ├── user_id → users.id
├── password_hash       ├── birthdate           ├── specialty
├── first_name          ├── document_number     ├── license_number
├── last_name           └── ...                 └── ...
├── role
└── ...

appointments
├── id (appointment.id)
├── patient_id → patients.id  ⚠️ NO es users.id
├── doctor_id → doctors.id    ⚠️ NO es users.id
├── service_id
└── ...
```

### Reglas de IDs

| Operación | ID a usar | De dónde viene |
|-----------|-----------|----------------|
| **Login** | `user.id` | Response de `/api/auth/login` → `user.id` |
| **Seleccionar doctor** | `user.id` | Response de `/api/services/doctors` → `[].id` |
| **Ver horarios** | `user.id` | Mismo que usaste para seleccionar doctor |
| **Crear cita** | `user.id` (doctor) | Backend convierte automáticamente |
| **Cancelar cita** | `appointment.id` | Response de `/api/appointments/my` → `[].id` |

### Ejemplo Práctico

```javascript
// 1. Usuario hace login
const loginResponse = await login('paciente@example.com', 'pass123');
const myUserId = loginResponse.user.id; // "user-uuid-123"

// 2. Backend automáticamente sabe que también existes en patients
// patient.id = "patient-uuid-abc"
// patient.user_id = "user-uuid-123"

// 3. Seleccionas un servicio y ves doctores disponibles
const doctors = await getDoctorsByService('service-uuid-456');
// doctors[0].id = "user-uuid-doctor-789"

// 4. Seleccionas fecha y ves horarios
const slots = await getAvailableSlots(
  'user-uuid-doctor-789',  // ✅ user.id del doctor
  'service-uuid-456',
  '2025-10-25'
);

// 5. Creas la cita
const appointment = await createAppointment({
  doctor_id: 'user-uuid-doctor-789',  // ✅ user.id del doctor
  service_id: 'service-uuid-456',
  appointment_date: '2025-10-25',
  appointment_time: '14:30',
  reason: 'Consulta'
});
// El backend convierte:
// - Tu user.id → patient.id automáticamente
// - doctor user.id → doctor.id automáticamente

// 6. Obtienes tus citas
const myCitas = await getMyAppointments();
// myCitas[0].id = "appointment-uuid-555"
// myCitas[0].patient_id = "patient-uuid-abc"  ℹ️ Diferente de tu user.id
// myCitas[0].doctor_id = "doctor-uuid-def"    ℹ️ Diferente del user.id

// 7. Cancelas una cita
await cancelAppointment(
  'appointment-uuid-555',  // ✅ appointment.id
  'Tengo un compromiso'
);
```

---

## 🚀 Flujo Completo de Reserva

### Componente: Página de Servicios

```javascript
import { useState, useEffect } from 'react';

function ServicesPage() {
  const [services, setServices] = useState([]);
  const [selectedService, setSelectedService] = useState(null);

  useEffect(() => {
    const fetchServices = async () => {
      const response = await fetch('http://localhost:8080/api/services');
      const data = await response.json();
      setServices(data);
    };
    fetchServices();
  }, []);

  return (
    <div>
      <h1>Selecciona un Servicio</h1>
      <div className="services-grid">
        {services.map(service => (
          <ServiceCard
            key={service.id}
            service={service}
            onSelect={() => setSelectedService(service)}
          />
        ))}
      </div>

      {selectedService && (
        <DoctorSelection serviceId={selectedService.id} />
      )}
    </div>
  );
}
```

### Componente: Selección de Doctor

```javascript
function DoctorSelection({ serviceId }) {
  const [doctors, setDoctors] = useState([]);
  const [selectedDoctor, setSelectedDoctor] = useState(null);

  useEffect(() => {
    const fetchDoctors = async () => {
      const response = await fetch(
        `http://localhost:8080/api/services/doctors?service_id=${serviceId}`
      );
      const data = await response.json();
      setDoctors(data);
    };
    fetchDoctors();
  }, [serviceId]);

  return (
    <div>
      <h2>Selecciona un Doctor</h2>
      <div className="doctors-grid">
        {doctors.map(doctor => (
          <DoctorCard
            key={doctor.id}
            doctor={doctor}
            onSelect={() => setSelectedDoctor(doctor)}
          />
        ))}
      </div>

      {selectedDoctor && (
        <DateTimeSelection
          doctorId={selectedDoctor.id}  {/* ✅ user.id del doctor */}
          serviceId={serviceId}
        />
      )}
    </div>
  );
}
```

### Componente: Selección de Fecha y Hora

```javascript
import { useState } from 'react';

function DateTimeSelection({ doctorId, serviceId }) {
  const [selectedDate, setSelectedDate] = useState('');
  const [slots, setSlots] = useState([]);
  const [selectedSlot, setSelectedSlot] = useState(null);

  const handleDateChange = async (date) => {
    setSelectedDate(date);

    // Obtener slots disponibles
    const response = await fetch(
      `http://localhost:8080/api/services/available-slots?` +
      `doctor_id=${doctorId}&service_id=${serviceId}&date=${date}`
    );
    const data = await response.json();
    setSlots(data);
  };

  const handleBooking = async () => {
    if (!selectedSlot) return;

    const token = localStorage.getItem('token');

    try {
      const response = await fetch('http://localhost:8080/api/appointments', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          doctor_id: doctorId,  // ✅ user.id del doctor
          service_id: serviceId,
          appointment_date: selectedDate,
          appointment_time: selectedSlot,
          reason: 'Consulta médica'
        })
      });

      if (!response.ok) {
        const error = await response.text();
        alert('Error: ' + error);
        return;
      }

      const appointment = await response.json();
      alert('Cita creada exitosamente!');

      // Redirigir a mis citas
      window.location.href = '/my-appointments';

    } catch (error) {
      console.error('Error:', error);
      alert('Error al crear la cita');
    }
  };

  return (
    <div>
      <h2>Selecciona Fecha y Hora</h2>

      <input
        type="date"
        min={new Date().toISOString().split('T')[0]}
        value={selectedDate}
        onChange={(e) => handleDateChange(e.target.value)}
      />

      {slots.length > 0 && (
        <div className="slots-grid">
          {slots
            .filter(slot => slot.available)
            .map(slot => (
              <button
                key={slot.time}
                className={selectedSlot === slot.time ? 'selected' : ''}
                onClick={() => setSelectedSlot(slot.time)}
              >
                {slot.time}
              </button>
            ))}
        </div>
      )}

      {selectedSlot && (
        <button onClick={handleBooking}>
          Confirmar Reserva
        </button>
      )}
    </div>
  );
}
```

### Componente: Mis Citas

```javascript
function MyAppointments() {
  const [appointments, setAppointments] = useState([]);

  useEffect(() => {
    const fetchAppointments = async () => {
      const token = localStorage.getItem('token');

      const response = await fetch('http://localhost:8080/api/appointments/my', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const data = await response.json();
      setAppointments(data);
    };

    fetchAppointments();
  }, []);

  const handleCancel = async (appointmentId) => {
    if (!confirm('¿Seguro que deseas cancelar esta cita?')) return;

    const reason = prompt('Motivo de cancelación:');
    if (!reason) return;

    const token = localStorage.getItem('token');

    try {
      const response = await fetch(
        `http://localhost:8080/api/appointments/cancel?id=${appointmentId}`,
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ reason })
        }
      );

      if (!response.ok) {
        const error = await response.text();
        alert('Error: ' + error);
        return;
      }

      alert('Cita cancelada exitosamente');

      // Recargar lista
      window.location.reload();

    } catch (error) {
      console.error('Error:', error);
      alert('Error al cancelar la cita');
    }
  };

  return (
    <div>
      <h1>Mis Citas</h1>

      <div className="appointments-list">
        {appointments.map(apt => (
          <div key={apt.id} className="appointment-card">
            <h3>{apt.service_name}</h3>
            <p>Fecha: {new Date(apt.scheduled_at).toLocaleString()}</p>
            <p>Duración: {apt.duration} minutos</p>
            <p>Estado: {apt.status}</p>

            {apt.status === 'pending' || apt.status === 'confirmed' ? (
              <button onClick={() => handleCancel(apt.id)}>
                Cancelar Cita
              </button>
            ) : null}
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## 🔐 Manejo de Autenticación

### Helper: API Client

```javascript
// utils/api.js
const API_BASE_URL = 'http://localhost:8080';

class ApiClient {
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

    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

    // Manejar error de autenticación
    if (response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      throw new Error('Session expired');
    }

    // Manejar otros errores
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    // Manejar 204 No Content
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

export const api = new ApiClient();
```

### Uso del API Client

```javascript
import { api } from './utils/api';

// Login
const loginData = await api.post('/api/auth/login', {
  email: 'paciente@example.com',
  password: 'password123'
});
localStorage.setItem('token', loginData.token);
localStorage.setItem('user', JSON.stringify(loginData.user));

// Obtener servicios
const services = await api.get('/api/services');

// Crear cita
const appointment = await api.post('/api/appointments', {
  doctor_id: 'user-uuid-doctor-123',
  service_id: 'service-uuid-456',
  appointment_date: '2025-10-25',
  appointment_time: '14:30',
  reason: 'Consulta'
});

// Cancelar cita
await api.put(`/api/appointments/cancel?id=${appointmentId}`, {
  reason: 'Motivo de cancelación'
});
```

---

## ⚠️ Errores Comunes y Soluciones

### Error 1: "insufficient permissions to cancel this appointment"

**Causa:** Este era el bug que se corrigió. Ahora ya no debería ocurrir.

**Solución:** Asegúrate de que:
- Estás enviando el token correcto
- El `appointment.id` es correcto
- El usuario logueado es el paciente o doctor de esa cita

### Error 2: "doctor does not offer this service"

**Causa:** Intentaste crear una cita con un doctor que no ofrece ese servicio.

**Solución:**
```javascript
// ❌ INCORRECTO
const doctor = doctors[0]; // Primer doctor de cualquier servicio
await createAppointment({ doctor_id: doctor.id, service_id: 'otro-servicio' });

// ✅ CORRECTO
const doctors = await getDoctorsByService(selectedServiceId);
const doctor = doctors[0]; // Doctor que SÍ ofrece el servicio
await createAppointment({ doctor_id: doctor.id, service_id: selectedServiceId });
```

### Error 3: "time slot is not available"

**Causa:** El horario seleccionado ya está ocupado.

**Solución:**
```javascript
// Siempre verificar disponibilidad antes de permitir selección
const slots = await getAvailableSlots(doctorId, serviceId, date);
const availableSlots = slots.filter(s => s.available);

// Mostrar solo los horarios disponibles
```

### Error 4: Token expirado (401)

**Causa:** El token JWT expiró (24 horas).

**Solución:**
```javascript
// Implementar verificación de expiración
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

// Verificar antes de cada request
if (isTokenExpired()) {
  // Redirigir a login
  window.location.href = '/login';
}
```

---

## 📋 Checklist de Integración

### Funcionalidades Básicas
- [ ] Login y almacenar token
- [ ] Logout y limpiar token
- [ ] Registro de nuevos usuarios (pacientes)
- [ ] Verificación de token expirado

### Reserva de Citas
- [ ] Listar servicios disponibles
- [ ] Filtrar doctores por servicio
- [ ] Calendario para seleccionar fecha
- [ ] Mostrar slots disponibles (solo los disponibles)
- [ ] Crear cita con validaciones
- [ ] Mostrar confirmación de cita creada

### Gestión de Citas
- [ ] Listar mis citas (paciente)
- [ ] Filtrar citas por estado (pending, confirmed, completed, cancelled)
- [ ] Cancelar cita con motivo
- [ ] Mostrar confirmación de cancelación

### Para Doctores
- [ ] Ver mis citas como doctor
- [ ] Confirmar citas pendientes
- [ ] Completar citas con notas
- [ ] Ver horarios configurados

### Para Administradores
- [ ] Dashboard con métricas
- [ ] Crear servicios
- [ ] Asignar servicios a doctores
- [ ] Crear horarios personalizados
- [ ] Ver analytics y reportes

---

## 🧪 Testing

### Endpoints para Probar Primero

```bash
# 1. Health Check
curl http://localhost:8080/

# 2. Crear usuario paciente
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User",
    "phone": "+51999999999",
    "role": "patient"
  }'

# 3. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 4. Listar servicios
curl http://localhost:8080/api/services

# 5. Ver doctores por servicio
curl "http://localhost:8080/api/services/doctors?service_id=SERVICE_ID"

# 6. Ver slots disponibles
curl "http://localhost:8080/api/services/available-slots?doctor_id=DOCTOR_USER_ID&service_id=SERVICE_ID&date=2025-10-25"

# 7. Crear cita (con token)
curl -X POST http://localhost:8080/api/appointments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "doctor_id": "DOCTOR_USER_ID",
    "service_id": "SERVICE_ID",
    "appointment_date": "2025-10-25",
    "appointment_time": "14:30",
    "reason": "Consulta de prueba"
  }'

# 8. Ver mis citas
curl http://localhost:8080/api/appointments/my \
  -H "Authorization: Bearer YOUR_TOKEN"

# 9. Cancelar cita
curl -X PUT "http://localhost:8080/api/appointments/cancel?id=APPOINTMENT_ID" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Prueba de cancelación"}'
```

---

## 📚 Documentación Completa

**Ver:** [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

Esta guía contiene:
- Todos los endpoints con ejemplos
- Request y response exactos
- Códigos de error
- Ejemplos en JavaScript
- Flujos completos

---

## 🆘 Soporte

**Dudas sobre la API:** Consultar con el equipo de backend
**Documentación:** Ver [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
**Swagger UI:** http://localhost:8080/swagger/index.html (cuando esté el servidor corriendo)

---

**Última actualización:** 2025-10-20
**Autor:** Backend Team - Zensoft
