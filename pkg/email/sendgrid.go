package email

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService handles sending emails via SendGrid
type EmailService struct {
	apiKey    string
	fromEmail string
	fromName  string
}

// NewEmailService creates a new email service instance
func NewEmailService(apiKey, fromEmail, fromName string) *EmailService {
	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// SendAppointmentCreated sends email when appointment is created
func (s *EmailService) SendAppointmentCreated(toEmail, patientName, doctorName, date, time string) error {
	subject := "Cita Médica Creada - Clinica Internacional"

	htmlContent := fmt.Sprintf(`
		<h2>Cita Médica Creada</h2>
		<p>Hola %s,</p>
		<p>Tu cita médica ha sido creada exitosamente.</p>
		<p><strong>Detalles:</strong></p>
		<ul>
			<li>Doctor: %s</li>
			<li>Fecha: %s</li>
			<li>Hora: %s</li>
			<li>Estado: Pendiente de confirmación</li>
		</ul>
		<p>El doctor confirmará tu cita pronto.</p>
		<p>Gracias,<br>Clinica Internacional</p>
	`, patientName, doctorName, date, time)

	return s.sendEmail(toEmail, subject, htmlContent)
}

// SendAppointmentConfirmed sends email when appointment is confirmed
func (s *EmailService) SendAppointmentConfirmed(toEmail, patientName, doctorName, date, time string) error {
	subject := "Cita Médica Confirmada - Clinica Internacional"

	htmlContent := fmt.Sprintf(`
		<h2>Cita Médica Confirmada</h2>
		<p>Hola %s,</p>
		<p>Tu cita médica ha sido confirmada por el doctor.</p>
		<p><strong>Detalles:</strong></p>
		<ul>
			<li>Doctor: %s</li>
			<li>Fecha: %s</li>
			<li>Hora: %s</li>
			<li>Estado: Confirmada</li>
		</ul>
		<p>Por favor, llega 15 minutos antes de tu cita.</p>
		<p>Gracias,<br>Clinica Internacional</p>
	`, patientName, doctorName, date, time)

	return s.sendEmail(toEmail, subject, htmlContent)
}

// SendAppointmentCancelled sends email when appointment is cancelled
func (s *EmailService) SendAppointmentCancelled(toEmail, patientName, doctorName, date, time, reason string) error {
	subject := "Cita Médica Cancelada - Clinica Internacional"

	htmlContent := fmt.Sprintf(`
		<h2>Cita Médica Cancelada</h2>
		<p>Hola %s,</p>
		<p>Tu cita médica ha sido cancelada.</p>
		<p><strong>Detalles:</strong></p>
		<ul>
			<li>Doctor: %s</li>
			<li>Fecha: %s</li>
			<li>Hora: %s</li>
			<li>Motivo: %s</li>
		</ul>
		<p>Puedes agendar una nueva cita cuando lo desees.</p>
		<p>Gracias,<br>Clinica Internacional</p>
	`, patientName, doctorName, date, time, reason)

	return s.sendEmail(toEmail, subject, htmlContent)
}

// SendAppointmentCompleted sends email when appointment is completed
func (s *EmailService) SendAppointmentCompleted(toEmail, patientName, doctorName, date string) error {
	subject := "Consulta Completada - Clinica Internacional"

	htmlContent := fmt.Sprintf(`
		<h2>Consulta Completada</h2>
		<p>Hola %s,</p>
		<p>Tu consulta médica ha sido completada.</p>
		<p><strong>Detalles:</strong></p>
		<ul>
			<li>Doctor: %s</li>
			<li>Fecha: %s</li>
		</ul>
		<p>Puedes ver el historial médico y las notas del doctor en tu perfil.</p>
		<p>Gracias por confiar en nosotros,<br>Clinica Internacional</p>
	`, patientName, doctorName, date)

	return s.sendEmail(toEmail, subject, htmlContent)
}

// sendEmail is the internal method that sends the email via SendGrid
func (s *EmailService) sendEmail(toEmail, subject, htmlContent string) error {
	// Check if API key is configured
	if s.apiKey == "" {
		log.Printf("SendGrid not configured, email not sent to %s: %s", toEmail, subject)
		return nil // Don't fail if SendGrid is not configured
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)

	if err != nil {
		log.Printf("Error sending email to %s: %v", toEmail, err)
		return err
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error (status %d) sending to %s: %s", response.StatusCode, toEmail, response.Body)
		return fmt.Errorf("sendgrid error: status %d", response.StatusCode)
	}

	log.Printf("Email sent successfully to %s: %s", toEmail, subject)
	return nil
}
