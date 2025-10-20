# 🎨 Documentación para el Equipo Frontend

## 📁 Archivos que Necesitas

Hemos preparado toda la documentación que necesitas para integrar el frontend con el backend sin errores.

---

## 📄 **ARCHIVO PRINCIPAL** → [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md)

**Este es el archivo que debes usar como referencia principal.**

### Contiene:

✅ **Todos los endpoints** con ejemplos exactos de request/response
✅ **Explicación de la arquitectura de IDs** (crítico para evitar errores)
✅ **Flujo completo de reserva de citas** paso a paso
✅ **Flujo de cancelación de citas** (bug corregido)
✅ **API Client helper** (código reutilizable)
✅ **Manejo de autenticación y tokens**
✅ **Códigos de error** con explicaciones
✅ **Notas sobre errores comunes** y cómo solucionarlos

### Formato:
- Ejemplos en JavaScript/fetch
- Request y Response JSON exactos
- Código listo para copiar y usar
- Explicaciones concisas

---

## 📚 Archivos Adicionales (Opcionales)

### 1. [FRONTEND_INTEGRATION_GUIDE.md](./FRONTEND_INTEGRATION_GUIDE.md)
**Guía detallada con componentes React completos**

Incluye:
- Componentes React listos para usar
- Código completo de:
  - Página de servicios
  - Selección de doctores
  - Calendario y horarios
  - Listado de citas
  - Cancelación de citas
- Manejo de estados
- Integración completa

**Usa este archivo si:** Necesitas ejemplos completos de implementación React.

---

### 2. [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
**Documentación exhaustiva de todos los endpoints**

Incluye:
- 22 endpoints documentados
- Ejemplos en cURL y JavaScript
- Explicación detallada de cada campo
- Validaciones y reglas de negocio
- Arquitectura completa del sistema

**Usa este archivo si:** Necesitas entender a profundidad cómo funciona la API.

---

## 🎯 Cómo Empezar

### 1. Lee primero: [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md)
Este archivo tiene todo lo que necesitas en formato conciso.

### 2. Implementa el flujo básico:

```javascript
// a) Login
const { token } = await login(email, password);

// b) Ver servicios
const services = await getServices();

// c) Ver doctores del servicio
const doctors = await getDoctorsByService(serviceId);

// d) Ver horarios disponibles
const slots = await getAvailableSlots(doctorId, serviceId, date);

// e) Crear cita
const appointment = await createAppointment({
  doctor_id: doctorId,
  service_id: serviceId,
  appointment_date: date,
  appointment_time: time,
  reason: 'Consulta'
});

// f) Ver mis citas
const myCitas = await getMyAppointments();

// g) Cancelar cita
await cancelAppointment(appointmentId, reason);
```

### 3. Copia el API Client helper

Del archivo [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md), copia la clase `ApiClient` que maneja:
- Autenticación automática
- Tokens expirados
- Errores
- Headers

---

## 🔑 Puntos Críticos

### 1. **Arquitectura de IDs** (Leer sección en FRONTEND_API_REFERENCE.md)

```
users.id        → Login, seleccionar doctores
patient.id      → Backend lo maneja automáticamente
doctor.id       → Backend lo maneja automáticamente
appointment.id  → Cancelar/modificar citas
```

**Regla simple:**
- Para seleccionar doctor → usa `user.id`
- Para cancelar cita → usa `appointment.id`

### 2. **Flujo de Reserva**

```
Servicios → Doctores → Fecha → Horarios → Crear Cita
```

Cada paso depende del anterior.

### 3. **Autenticación**

```javascript
// Guardar token después de login
localStorage.setItem('token', token);

// Enviar en cada request protegido
headers: {
  'Authorization': `Bearer ${token}`
}

// Manejar token expirado (401)
if (response.status === 401) {
  // Redirigir a login
  window.location.href = '/login';
}
```

---

## 🐛 Bug Crítico Corregido

### Problema anterior:
La cancelación de citas fallaba con "insufficient permissions".

### Solución:
✅ Ya está corregido en el backend.

### Cómo usar ahora:
```javascript
// Obtener appointment.id de "Mis Citas"
const appointments = await getMyAppointments();
const appointmentId = appointments[0].id;

// Cancelar
await cancelAppointment(appointmentId, 'Motivo de cancelación');
```

---

## 📊 Endpoints Más Usados

### Públicos (sin token):
- `GET /api/services` - Listar servicios
- `GET /api/services/doctors?service_id=X` - Doctores por servicio
- `GET /api/services/available-slots?...` - Horarios disponibles
- `POST /api/auth/login` - Login
- `POST /api/users` - Registro

### Protegidos (requieren token):
- `GET /api/users/me` - Mi perfil
- `POST /api/appointments` - Crear cita
- `GET /api/appointments/my` - Mis citas
- `PUT /api/appointments/cancel?id=X` - Cancelar cita
- `GET /api/appointments/doctor` - Citas del doctor (solo doctor)

### Admin (requieren token de admin):
- `GET /api/analytics/dashboard` - Dashboard
- `POST /api/services/create` - Crear servicio
- `POST /api/services/assign` - Asignar servicio a doctor
- `POST /api/schedules` - Crear horario

---

## 🔍 Testing Rápido

### 1. Verifica que el backend esté corriendo:
```bash
curl http://localhost:8080/
# Debe retornar: "Sistema de Reservas - API Running"
```

### 2. Prueba login:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 3. Lista servicios:
```bash
curl http://localhost:8080/api/services
```

---

## 📋 Checklist de Integración

### Funcionalidades Básicas
- [ ] Implementar login y guardar token
- [ ] Implementar logout y limpiar token
- [ ] Registro de usuarios
- [ ] Manejo de token expirado (401)

### Reserva de Citas
- [ ] Listar servicios
- [ ] Mostrar doctores por servicio
- [ ] Calendario para fecha
- [ ] Mostrar slots disponibles (filtrar `available: true`)
- [ ] Crear cita con validaciones
- [ ] Confirmación de cita creada

### Gestión de Citas
- [ ] Listar mis citas
- [ ] Filtrar por estado (pending, confirmed, completed, cancelled)
- [ ] Cancelar cita con motivo
- [ ] Confirmación de cancelación

### Vista Doctor
- [ ] Ver mis citas como doctor
- [ ] Confirmar citas pendientes
- [ ] Completar citas con notas

### Vista Admin
- [ ] Dashboard con métricas
- [ ] Gestión de servicios
- [ ] Asignar servicios a doctores
- [ ] Configurar horarios

---

## 🆘 ¿Tienes Dudas?

### 1. **Sobre endpoints y cómo usarlos**
→ Ver [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md)

### 2. **Sobre implementación con React**
→ Ver [FRONTEND_INTEGRATION_GUIDE.md](./FRONTEND_INTEGRATION_GUIDE.md)

### 3. **Sobre arquitectura y detalles técnicos**
→ Ver [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

### 4. **Error específico**
→ Buscar en sección "Errores Comunes" de [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md)

### 5. **Necesitas ejemplos de código**
→ Todos los archivos tienen código listo para copiar

---

## 🚀 URLs Importantes

### Desarrollo:
- **API Base URL:** `http://localhost:8080`
- **Swagger UI:** `http://localhost:8080/swagger/index.html`

### Producción:
- **API Base URL:** `https://tu-dominio.com` (configurar en frontend)

---

## 📞 Contacto

**Dudas sobre la API:** Equipo de Backend - Zensoft
**Documentación actualizada:** 2025-10-20

---

## ✅ Resumen

**Archivo principal:** [FRONTEND_API_REFERENCE.md](./FRONTEND_API_REFERENCE.md)

**Contiene todo lo que necesitas:**
- ✅ Endpoints con ejemplos
- ✅ Flujos completos
- ✅ Código JavaScript reutilizable
- ✅ Explicación de IDs
- ✅ Manejo de errores
- ✅ Autenticación

**Empieza leyendo ese archivo y luego implementa paso a paso.**

¡Éxito con la integración! 🎉
