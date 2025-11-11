-- Script de verificación para diagnosticar el problema con available-slots
-- Ejecutar en PostgreSQL

-- 1. Verificar si el user_id existe en la tabla users
SELECT
    'Usuario existe' as verificacion,
    id,
    email,
    first_name,
    last_name,
    role
FROM users
WHERE id = '6e117b8d-d769-4b09-9271-97a81b7529d2';

-- 2. Verificar si existe un doctor con ese user_id
SELECT
    'Doctor existe' as verificacion,
    d.id as doctor_id,
    d.user_id,
    u.first_name || ' ' || u.last_name as doctor_name
FROM doctors d
JOIN users u ON d.user_id = u.id
WHERE d.user_id = '6e117b8d-d769-4b09-9271-97a81b7529d2';

-- 3. Verificar si el servicio existe
SELECT
    'Servicio existe' as verificacion,
    id,
    name,
    duration_minutes,
    is_active
FROM services
WHERE id = 'ae432bad-15ff-467f-84fe-0e38744607b9';

-- 4. Verificar la cita que se quiere excluir
SELECT
    'Cita a excluir' as verificacion,
    a.id,
    a.patient_id,
    a.doctor_id,
    a.service_id,
    a.scheduled_at,
    a.status
FROM appointments a
WHERE a.id = '704e8b16-7763-4e93-b52c-f26e4e3e5e84';

-- 5. Listar todos los doctores disponibles para referencia
SELECT
    'Todos los doctores' as verificacion,
    d.id as doctor_id,
    d.user_id,
    u.first_name || ' ' || u.last_name as doctor_name,
    u.email
FROM doctors d
JOIN users u ON d.user_id = u.id
WHERE u.is_active = true
ORDER BY u.first_name;

-- 6. Verificar qué doctor_id se está usando en appointments
SELECT
    'Doctor en appointments' as verificacion,
    DISTINCT a.doctor_id,
    COUNT(*) as total_citas
FROM appointments a
GROUP BY a.doctor_id;
