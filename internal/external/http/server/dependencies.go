package server

import (
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	get_appointment_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_by_id"
	list_appointments_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/list_appointments"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
)

type Dependencies struct {
	Cache     cache.Cache
	DbService *persistence.DbService

	AppointmentRepository appointment_repository_contract.Repository

	CreateAppointmentUseCase  create_appointment_contract.UseCase
	GetAppointmentByIdUseCase get_appointment_by_id_contract.UseCase
	ListAppointmentsUseCase   list_appointments_contract.UseCase
}
