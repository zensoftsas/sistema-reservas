package service

import (
	"context"
	"errors"
	"fmt"

	"version-1-0/internal/repository"
)

// UnassignServiceFromDoctorUseCase handles removing a service from a doctor
type UnassignServiceFromDoctorUseCase struct {
	doctorServiceRepo repository.DoctorServiceRepository
	appointmentRepo   repository.AppointmentRepository
	userRepo          repository.UserRepository
}

// NewUnassignServiceFromDoctorUseCase creates a new instance
func NewUnassignServiceFromDoctorUseCase(
	doctorServiceRepo repository.DoctorServiceRepository,
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
) *UnassignServiceFromDoctorUseCase {
	return &UnassignServiceFromDoctorUseCase{
		doctorServiceRepo: doctorServiceRepo,
		appointmentRepo:   appointmentRepo,
		userRepo:          userRepo,
	}
}

// Execute removes a service assignment from a doctor
func (uc *UnassignServiceFromDoctorUseCase) Execute(ctx context.Context, userID, serviceID string) error {
	// Validate inputs
	if userID == "" {
		return errors.New("doctor ID is required")
	}

	if serviceID == "" {
		return errors.New("service ID is required")
	}

	// Get the actual doctor.id from doctors table
	doctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 1. Verify the assignment exists
	assignment, err := uc.doctorServiceRepo.FindByDoctorAndService(ctx, doctorID, serviceID)
	if err != nil {
		return err
	}

	if assignment == nil {
		return errors.New("el doctor no tiene este servicio asignado")
	}

	// 2. Check for future appointments with this service
	futureAppointments, err := uc.appointmentRepo.CountFutureAppointmentsByDoctorAndService(ctx, doctorID, serviceID)
	if err != nil {
		return err
	}

	if futureAppointments > 0 {
		return fmt.Errorf("no se puede desasignar el servicio porque el doctor tiene %d citas futuras programadas", futureAppointments)
	}

	// 3. Remove the assignment
	return uc.doctorServiceRepo.Remove(ctx, doctorID, serviceID)
}
